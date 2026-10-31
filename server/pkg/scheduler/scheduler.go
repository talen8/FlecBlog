package scheduler

import (
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"

	"flec_blog/pkg/logsetup"
)

// Job 定时任务接口
type Job interface {
	Name() string     // 任务名称
	Execute() error   // 执行任务
	Schedule() string // cron表达式
}

// SimpleJob 简单任务实现
type SimpleJob struct {
	name     string
	schedule string
	execute  func() error
}

// NewJob 创建一个简单任务
func NewJob(name, schedule string, execute func() error) Job {
	return &SimpleJob{
		name:     name,
		schedule: schedule,
		execute:  execute,
	}
}

func (j *SimpleJob) Name() string     { return j.name }
func (j *SimpleJob) Execute() error   { return j.execute() }
func (j *SimpleJob) Schedule() string { return j.schedule }

// Scheduler 调度器
type Scheduler struct {
	cron *cron.Cron
	jobs map[string]Job
}

// NewScheduler 创建调度器实例
func NewScheduler() *Scheduler {
	// 使用秒级别的cron表达式（6个字段）
	c := cron.New(cron.WithSeconds())

	return &Scheduler{
		cron: c,
		jobs: make(map[string]Job),
	}
}

// AddJob 添加定时任务
func (s *Scheduler) AddJob(job Job) error {
	name := job.Name()
	schedule := job.Schedule()

	// 包装任务执行函数，添加日志和错误处理
	wrappedFunc := func() {
		defer logsetup.Recover("定时任务", slog.String("job", name))
		startTime := time.Now()
		slog.Info("调度任务开始执行", slog.String("job", name))

		if err := job.Execute(); err != nil {
			slog.Error("调度任务执行失败", slog.String("job", name), slog.Any("err", err))
		} else {
			duration := time.Since(startTime)
			slog.Info("调度任务执行成功", slog.String("job", name), slog.Duration("cost", duration))
		}
	}

	// 添加到cron
	_, err := s.cron.AddFunc(schedule, wrappedFunc)
	if err != nil {
		return err
	}

	s.jobs[name] = job
	return nil
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	slog.Info("调度器已停止")
	ctx := s.cron.Stop()
	<-ctx.Done()
}

// GetJobs 获取所有已注册的任务
func (s *Scheduler) GetJobs() map[string]Job {
	jobs := make(map[string]Job, len(s.jobs))
	for name, job := range s.jobs {
		jobs[name] = job
	}
	return jobs
}
