package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"
)

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

func main() {
	// 模拟证书信息
	certInfo := &CertInfo{
		Domain:        "example.com",
		Issuer:        "Let's Encrypt Authority X3",
		NotBefore:     "2025-01-01 12:00:00",
		NotAfter:      "2025-08-03 12:00:00",
		DaysLeft:      7,
		CheckTime:     time.Now().Format("2006-01-02 15:04:05"),
		ServerInfo:    "test-server",
		ScriptPath:    "SSL证书监控服务",
		ConfigPath:    "/mnt/restart-life-api/configs/live.yaml",
		AlertClass:    "danger",
		AlertIcon:     "⚠️",
		AlertTitle:    "警告：SSL证书即将到期",
		AlertMessage:  "您的SSL证书将在 7 天后到期，请尽快安排续期工作。",
		BoxClass:      "danger-box",
		DaysLeftColor: "#f44336",
		IsUrgent:      true,
	}

	// 读取邮件模板
	tmplPath := "./template/ssl_cert_expiry.html"
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		fmt.Printf("failed to read email template: %v\n", err)
		return
	}

	// 解析模板
	tmpl, err := template.New("ssl_cert_expiry").Parse(string(tmplContent))
	if err != nil {
		fmt.Printf("failed to parse email template: %v\n", err)
		return
	}

	// 渲染模板
	var body bytes.Buffer
	if err := tmpl.Execute(&body, certInfo); err != nil {
		fmt.Printf("failed to render email template: %v\n", err)
		return
	}

	// 输出渲染结果到文件
	outputFile := "./test_output.html"
	if err := os.WriteFile(outputFile, body.Bytes(), 0644); err != nil {
		fmt.Printf("failed to write output file: %v\n", err)
		return
	}

	fmt.Printf("Template rendered successfully! Output saved to %s\n", outputFile)
	fmt.Printf("Template size: %d bytes\n", body.Len())
}
