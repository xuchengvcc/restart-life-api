package job

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// JobManager 定时任务管理器
type JobManager struct {
	jobs   map[string]Job
	logger *logrus.Logger
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex
}

// Job 定时任务接口
type Job interface {
	// GetName 获取任务名称
	GetName() string

	// GetInterval 获取执行间隔
	GetInterval() time.Duration

	// Execute 执行任务
	Execute() error

	// OnStart 任务启动时的回调
	OnStart() error

	// OnStop 任务停止时的回调
	OnStop() error
}

// NewJobManager 创建新的任务管理器
func NewJobManager(logger *logrus.Logger) *JobManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &JobManager{
		jobs:   make(map[string]Job),
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddJob 添加任务
func (jm *JobManager) AddJob(job Job) error {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	name := job.GetName()
	if _, exists := jm.jobs[name]; exists {
		return ErrJobAlreadyExists
	}

	jm.jobs[name] = job
	jm.logger.WithField("job", name).Info("job added to manager")
	return nil
}

// RemoveJob 移除任务
func (jm *JobManager) RemoveJob(name string) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	if _, exists := jm.jobs[name]; exists {
		delete(jm.jobs, name)
		jm.logger.WithField("job", name).Info("job removed from manager")
	}
}

// StartJob 启动单个任务
func (jm *JobManager) StartJob(name string) error {
	jm.mu.RLock()
	job, exists := jm.jobs[name]
	jm.mu.RUnlock()

	if !exists {
		return ErrJobNotFound
	}

	// 调用任务的启动回调
	if err := job.OnStart(); err != nil {
		jm.logger.WithField("job", name).WithError(err).Error("job start callback failed")
		return err
	}

	jm.wg.Add(1)
	go jm.runJob(job)

	jm.logger.WithField("job", name).Info("job started")
	return nil
}

// StartAll 启动所有任务
func (jm *JobManager) StartAll() error {
	jm.mu.RLock()
	jobList := make([]Job, 0, len(jm.jobs))
	for _, job := range jm.jobs {
		jobList = append(jobList, job)
	}
	jm.mu.RUnlock()

	for _, job := range jobList {
		if err := jm.StartJob(job.GetName()); err != nil {
			jm.logger.WithField("job", job.GetName()).WithError(err).Error("failed to start job")
			return err
		}
	}

	jm.logger.WithField("count", len(jobList)).Info("all jobs started")
	return nil
}

// Stop 停止所有任务
func (jm *JobManager) Stop() {
	jm.logger.Info("stopping job manager...")

	// 停止所有任务
	jm.mu.RLock()
	for _, job := range jm.jobs {
		if err := job.OnStop(); err != nil {
			jm.logger.WithField("job", job.GetName()).WithError(err).Error("job stop callback failed")
		}
	}
	jm.mu.RUnlock()

	// 取消上下文
	jm.cancel()

	// 等待所有任务结束
	jm.wg.Wait()

	jm.logger.Info("job manager stopped")
}

// runJob 运行单个任务
func (jm *JobManager) runJob(job Job) {
	defer jm.wg.Done()

	name := job.GetName()
	interval := job.GetInterval()

	jm.logger.WithFields(logrus.Fields{
		"job":      name,
		"interval": interval,
	}).Info("job started running")

	// 立即执行一次
	if err := job.Execute(); err != nil {
		jm.logger.WithField("job", name).WithError(err).Error("job execution failed")
	}

	// 创建定时器
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-jm.ctx.Done():
			jm.logger.WithField("job", name).Info("job stopped due to context cancellation")
			return
		case <-ticker.C:
			if err := job.Execute(); err != nil {
				jm.logger.WithField("job", name).WithError(err).Error("job execution failed")
			}
		}
	}
}

// GetJobStatus 获取任务状态
func (jm *JobManager) GetJobStatus() map[string]interface{} {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	status := make(map[string]interface{})
	status["total_jobs"] = len(jm.jobs)
	status["running"] = jm.ctx.Err() == nil

	jobs := make([]map[string]interface{}, 0, len(jm.jobs))
	for name, job := range jm.jobs {
		jobInfo := map[string]interface{}{
			"name":     name,
			"interval": job.GetInterval().String(),
		}
		jobs = append(jobs, jobInfo)
	}
	status["jobs"] = jobs

	return status
}
