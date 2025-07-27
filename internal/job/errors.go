package job

import "errors"

var (
	// ErrJobAlreadyExists 任务已存在错误
	ErrJobAlreadyExists = errors.New("job already exists")

	// ErrJobNotFound 任务未找到错误
	ErrJobNotFound = errors.New("job not found")

	// ErrJobExecutionFailed 任务执行失败错误
	ErrJobExecutionFailed = errors.New("job execution failed")
)
