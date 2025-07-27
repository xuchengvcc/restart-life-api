package job

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

// SSLCertMonitorJob SSL证书监控任务
type SSLCertMonitorJob struct {
	logger       *logrus.Logger
	emailService services.EmailService
	emailConfig  *config.EmailConfig
	redisClient  *redis.Client
	domain       string
	interval     time.Duration
	alertDays    []int
}

// CertInfo 证书信息结构
type CertInfo struct {
	Domain     string
	Issuer     string
	NotBefore  string
	NotAfter   string
	DaysLeft   int
	CheckTime  string
	ServerInfo string
	ScriptPath string
	ConfigPath string

	// 邮件模板使用的字段
	AlertClass    string
	AlertIcon     string
	AlertTitle    string
	AlertMessage  string
	BoxClass      string
	DaysLeftColor string
	IsUrgent      bool
}

// NewSSLCertMonitorJob 创建SSL证书监控任务
func NewSSLCertMonitorJob(logger *logrus.Logger, emailService services.EmailService, emailConfig *config.EmailConfig, redisClient *redis.Client) *SSLCertMonitorJob {
	// 从配置中读取参数
	domain := viper.GetString("ssl_monitor.domain")
	if domain == "" {
		domain = "asecondchance.cn" // 默认值
	}

	intervalStr := viper.GetString("ssl_monitor.check_interval")
	if intervalStr == "" {
		intervalStr = "24h" // 默认每24小时检查一次
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		logger.WithError(err).Warn("invalid check_interval, using default 24h")
		interval = 24 * time.Hour
	}

	alertDays := viper.GetIntSlice("ssl_monitor.alert_days")
	if len(alertDays) == 0 {
		alertDays = []int{14, 7, 3} // 默认值
	}

	return &SSLCertMonitorJob{
		logger:       logger,
		emailService: emailService,
		emailConfig:  emailConfig,
		redisClient:  redisClient,
		domain:       domain,
		interval:     interval,
		alertDays:    alertDays,
	}
}

// GetName 获取任务名称
func (j *SSLCertMonitorJob) GetName() string {
	return "ssl-cert-monitor"
}

// GetInterval 获取执行间隔
func (j *SSLCertMonitorJob) GetInterval() time.Duration {
	return j.interval
}

// OnStart 任务启动时的回调
func (j *SSLCertMonitorJob) OnStart() error {
	j.logger.WithFields(logrus.Fields{
		"domain":     j.domain,
		"interval":   j.interval,
		"alert_days": j.alertDays,
		"pic_email":  j.emailConfig.GetPicEmail(),
	}).Info("SSL certificate monitor job starting")
	return nil
}

// OnStop 任务停止时的回调
func (j *SSLCertMonitorJob) OnStop() error {
	j.logger.Info("SSL certificate monitor job stopping")
	return nil
}

// getAlertKey 生成告警状态的Redis键
func (j *SSLCertMonitorJob) getAlertKey(domain string, alertDay int) string {
	return fmt.Sprintf("ssl_cert_alert:%s:%d_days", domain, alertDay)
}

// isAlertSent 检查是否已发送过告警
func (j *SSLCertMonitorJob) isAlertSent(domain string, alertDay int) bool {
	ctx := context.Background()
	key := j.getAlertKey(domain, alertDay)

	result, err := j.redisClient.Exists(ctx, key).Result()
	if err != nil {
		j.logger.WithError(err).Error("failed to check alert status in redis")
		return false // 出错时认为未发送，允许重新尝试
	}

	return result > 0
}

// setAlertSent 标记告警已发送
func (j *SSLCertMonitorJob) setAlertSent(domain string, alertDay int, daysLeft int) error {
	ctx := context.Background()
	key := j.getAlertKey(domain, alertDay)

	// 计算过期时间：比告警阈值多一天，确保不会重复发送
	// 例如：14天告警发送后，设置15天过期时间
	expiration := time.Duration(alertDay+1) * 24 * time.Hour

	// 如果剩余天数很少，设置较短的过期时间
	if daysLeft <= 1 {
		expiration = 2 * 24 * time.Hour // 2天
	} else if daysLeft <= 3 {
		expiration = 4 * 24 * time.Hour // 4天
	}

	value := fmt.Sprintf("sent_at:%d", time.Now().Unix())
	err := j.redisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		j.logger.WithError(err).Error("failed to set alert status in redis")
		return err
	}

	j.logger.WithFields(logrus.Fields{
		"key":        key,
		"expiration": expiration,
		"days_left":  daysLeft,
	}).Debug("alert status saved to redis")

	return nil
}

// cleanupExpiredAlerts 清理已过期的告警状态（可选，Redis会自动过期）
func (j *SSLCertMonitorJob) cleanupExpiredAlerts(domain string, currentDaysLeft int) {
	ctx := context.Background()

	// 清理比当前剩余天数更大的告警状态
	for _, alertDay := range j.alertDays {
		if currentDaysLeft > alertDay {
			key := j.getAlertKey(domain, alertDay)
			err := j.redisClient.Del(ctx, key).Err()
			if err != nil {
				j.logger.WithError(err).WithField("key", key).Debug("failed to cleanup alert status")
			} else {
				j.logger.WithField("key", key).Debug("cleaned up outdated alert status")
			}
		}
	}
}

// Execute 执行任务
func (j *SSLCertMonitorJob) Execute() error {
	j.logger.WithField("domain", j.domain).Info("checking SSL certificate")

	// 获取证书信息
	certInfo, err := j.getCertInfo(j.domain)
	if err != nil {
		j.logger.WithError(err).Error("failed to get certificate info")
		return err
	}

	j.logger.WithFields(logrus.Fields{
		"domain":    certInfo.Domain,
		"days_left": certInfo.DaysLeft,
		"issuer":    certInfo.Issuer,
	}).Info("certificate info retrieved")

	// 检查是否需要发送邮件
	shouldAlert := false
	alertReason := ""
	alertDay := 0

	for _, day := range j.alertDays {
		if certInfo.DaysLeft <= day {
			if !j.isAlertSent(j.domain, day) {
				shouldAlert = true
				alertDay = day
				alertReason = fmt.Sprintf("%d天内到期", day)
				break
			}
		}
	}

	if shouldAlert {
		j.logger.WithField("reason", alertReason).Info("sending SSL certificate alert email")

		if err := j.sendAlertEmail(certInfo); err != nil {
			j.logger.WithError(err).Error("failed to send alert email")
			return err
		}

		// 标记告警已发送
		if err := j.setAlertSent(j.domain, alertDay, certInfo.DaysLeft); err != nil {
			j.logger.WithError(err).Error("failed to save alert status")
			// 不返回错误，避免影响主流程
		}

		j.logger.WithField("recipient", j.emailConfig.GetPicEmail()).Info("SSL certificate alert email sent successfully")
	} else {
		j.logger.Info("no alert needed, certificate status is normal or alert already sent")
	}

	// 清理已过期的告警状态
	j.cleanupExpiredAlerts(j.domain, certInfo.DaysLeft)

	return nil
}

// getCertInfo 获取证书信息
func (j *SSLCertMonitorJob) getCertInfo(domain string) (*CertInfo, error) {
	// 连接到域名获取证书
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		ServerName: domain,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", domain, err)
	}
	defer conn.Close()

	// 获取证书链
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates found for %s", domain)
	}

	cert := certs[0]

	// 计算剩余天数
	now := time.Now()
	daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)

	// 获取系统信息
	hostname, _ := os.Hostname()

	certInfo := &CertInfo{
		Domain:     domain,
		Issuer:     cert.Issuer.CommonName,
		NotBefore:  cert.NotBefore.Format("2006-01-02 15:04:05"),
		NotAfter:   cert.NotAfter.Format("2006-01-02 15:04:05"),
		DaysLeft:   daysLeft,
		CheckTime:  now.Format("2006-01-02 15:04:05"),
		ServerInfo: hostname,
		ScriptPath: "SSL证书监控服务",
		ConfigPath: "/mnt/restart-life-api/configs/live.yaml",
	}

	// 设置警告级别
	if daysLeft <= 3 {
		certInfo.AlertClass = "critical"
		certInfo.AlertIcon = "🚨"
		certInfo.AlertTitle = "紧急：SSL证书即将到期！"
		certInfo.AlertMessage = fmt.Sprintf("您的SSL证书将在 %d 天后到期，请立即进行续期操作！", daysLeft)
		certInfo.BoxClass = "danger-box"
		certInfo.DaysLeftColor = "#d32f2f"
		certInfo.IsUrgent = true
	} else if daysLeft <= 7 {
		certInfo.AlertClass = "danger"
		certInfo.AlertIcon = "⚠️"
		certInfo.AlertTitle = "警告：SSL证书即将到期"
		certInfo.AlertMessage = fmt.Sprintf("您的SSL证书将在 %d 天后到期，请尽快安排续期工作。", daysLeft)
		certInfo.BoxClass = "danger-box"
		certInfo.DaysLeftColor = "#f44336"
		certInfo.IsUrgent = true
	} else if daysLeft <= 14 {
		certInfo.AlertClass = "warning"
		certInfo.AlertIcon = "📋"
		certInfo.AlertTitle = "提醒：SSL证书续期提醒"
		certInfo.AlertMessage = fmt.Sprintf("您的SSL证书将在 %d 天后到期，建议开始准备续期工作。", daysLeft)
		certInfo.BoxClass = "warning-box"
		certInfo.DaysLeftColor = "#ff9800"
		certInfo.IsUrgent = false
	}

	return certInfo, nil
}

// sendAlertEmail 发送告警邮件
func (j *SSLCertMonitorJob) sendAlertEmail(certInfo *CertInfo) error {
	// 读取邮件模板
	tmplPath := "/mnt/restart-life-api/template/ssl_cert_expiry.html"
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read email template: %v", err)
	}

	// 解析模板
	tmpl, err := template.New("ssl_cert_expiry").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	// 渲染模板
	var body bytes.Buffer
	if err := tmpl.Execute(&body, certInfo); err != nil {
		return fmt.Errorf("failed to render email template: %v", err)
	}

	// 构建邮件主题
	subject := fmt.Sprintf("【重要】SSL证书到期提醒 - %s (%d天后到期)", certInfo.Domain, certInfo.DaysLeft)

	// 发送邮件
	recipient := j.emailConfig.GetPicEmail()
	return j.emailService.SendNotificationEmail(recipient, subject, body.String())
}
