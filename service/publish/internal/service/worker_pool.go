package service

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// PublishWorkerPool 发布任务工作池
type PublishWorkerPool struct {
	publishSvc  *PublishService
	concurrency int
	taskQueue   chan int64
	stopCh      chan struct{}
	wg          sync.WaitGroup
	rng         *rand.Rand
	mu          sync.Mutex
	minInterval time.Duration
	maxInterval time.Duration
}

// NewPublishWorkerPool 创建发布工作池
func NewPublishWorkerPool(publishSvc *PublishService, concurrency int) *PublishWorkerPool {
	if concurrency <= 0 {
		concurrency = 5
	}
	return &PublishWorkerPool{
		publishSvc:  publishSvc,
		concurrency: concurrency,
		taskQueue:   make(chan int64, concurrency*10),
		stopCh:      make(chan struct{}),
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
		minInterval: 500 * time.Millisecond,
		maxInterval: 3 * time.Second,
	}
}

// Start 启动工作池
func (p *PublishWorkerPool) Start() {
	for i := 0; i < p.concurrency; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// Stop 停止工作池（等待进行中的任务完成）
func (p *PublishWorkerPool) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}

// Submit 提交任务到队列
func (p *PublishWorkerPool) Submit(taskID int64) {
	select {
	case p.taskQueue <- taskID:
	case <-p.stopCh:
	}
}

// worker 工作协程
func (p *PublishWorkerPool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.stopCh:
			return
		case taskID, ok := <-p.taskQueue:
			if !ok {
				return
			}

			delay := p.getNextDelay()
			select {
			case <-p.stopCh:
				return
			case <-time.After(delay):
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			p.publishSvc.ExecutePublishTask(ctx, taskID, nil)
			cancel()
		}
	}
}

// getNextDelay 获取下一个随机延迟
func (p *PublishWorkerPool) getNextDelay() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()

	rangeNs := p.maxInterval.Nanoseconds() - p.minInterval.Nanoseconds()
	if rangeNs <= 0 {
		return p.minInterval
	}
	return p.minInterval + time.Duration(p.rng.Int63n(rangeNs))
}
