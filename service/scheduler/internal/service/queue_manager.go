package service

import (
	"fmt"
	"sync"
	"time"
)

// QueueManager 队列管理器
type QueueManager struct {
	mu          sync.RWMutex
	queues      map[string]*PriorityQueue // 多队列支持
	defaultQueue string
	maxWorkers  int
	workers     map[int64]*Worker
	preemptChan chan *QueueItem // 抢占信号通道
}

// Worker 工作协程
type Worker struct {
	ID       int64
	TaskID   int64
	IsActive bool
	StartedAt time.Time
}

// NewQueueManager 创建队列管理器
func NewQueueManager(maxWorkers int) *QueueManager {
	defaultQueue := "default"
	qm := &QueueManager{
		queues:       make(map[string]*PriorityQueue),
		defaultQueue: defaultQueue,
		maxWorkers:   maxWorkers,
		workers:      make(map[int64]*Worker),
		preemptChan:  make(chan *QueueItem, 100),
	}
	qm.queues[defaultQueue] = NewPriorityQueue()
	return qm
}

// ==================== 队列管理 ====================

// CreateQueue 创建队列
func (qm *QueueManager) CreateQueue(name string) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	if _, exists := qm.queues[name]; !exists {
		qm.queues[name] = NewPriorityQueue()
	}
}

// GetQueue 获取队列
func (qm *QueueManager) GetQueue(name string) *PriorityQueue {
	qm.mu.RLock()
	defer qm.mu.RUnlock()
	return qm.queues[name]
}

// DeleteQueue 删除队列
func (qm *QueueManager) DeleteQueue(name string) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	if name != qm.defaultQueue {
		delete(qm.queues, name)
	}
}

// ==================== 任务入队 ====================

// Enqueue 入队到默认队列
func (qm *QueueManager) Enqueue(item *QueueItem) {
	qm.EnqueueTo(qm.defaultQueue, item)
}

// EnqueueTo 入队到指定队列
func (qm *QueueManager) EnqueueTo(queueName string, item *QueueItem) {
	qm.mu.RLock()
	queue, exists := qm.queues[queueName]
	qm.mu.RUnlock()

	if !exists {
		qm.CreateQueue(queueName)
		qm.mu.RLock()
		queue = qm.queues[queueName]
		qm.mu.RUnlock()
	}

	// 如果是紧急任务，检查是否需要抢占
	if item.Priority >= 3 {
		qm.checkPreemption(item)
	}

	queue.Enqueue(item)
}

// ==================== 任务出队 ====================

// Dequeue 从默认队列出队
func (qm *QueueManager) Dequeue() *QueueItem {
	return qm.DequeueFrom(qm.defaultQueue)
}

// DequeueFrom 从指定队列出队
func (qm *QueueManager) DequeueFrom(queueName string) *QueueItem {
	qm.mu.RLock()
	queue, exists := qm.queues[queueName]
	qm.mu.RUnlock()

	if !exists {
		return nil
	}
	return queue.Dequeue()
}

// DequeueWithPriority 按优先级出队（紧急任务优先）
func (qm *QueueManager) DequeueWithPriority() *QueueItem {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	// 先检查紧急队列
	for _, queue := range qm.queues {
		if queue.Len() > 0 {
			peek := queue.Peek()
			if peek != nil && peek.Priority >= 3 {
				return queue.Dequeue()
			}
		}
	}

	// 再检查高优先级队列
	for _, queue := range qm.queues {
		if queue.Len() > 0 {
			peek := queue.Peek()
			if peek != nil && peek.Priority >= 2 {
				return queue.Dequeue()
			}
		}
	}

	// 最后从默认队列获取
	if queue, exists := qm.queues[qm.defaultQueue]; exists {
		return queue.Dequeue()
	}

	return nil
}

// ==================== 任务抢占 ====================

// checkPreemption 检查是否需要抢占
func (qm *QueueManager) checkPreemption(urgentItem *QueueItem) {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	// 查找正在执行的低优先级任务
	for taskID, worker := range qm.workers {
		if worker.IsActive && taskID > 0 {
			// 获取当前任务的优先级
			for _, queue := range qm.queues {
				if queue.Contains(taskID) {
					// 如果当前任务优先级低于紧急任务，触发抢占
					select {
					case qm.preemptChan <- urgentItem:
					default:
					}
					return
				}
			}
		}
	}
}

// GetPreemptChannel 获取抢占信号通道
func (qm *QueueManager) GetPreemptChannel() <-chan *QueueItem {
	return qm.preemptChan
}

// ==================== 工作协程管理 ====================

// AcquireWorker 获取工作协程
func (qm *QueueManager) AcquireWorker(taskID int64) (*Worker, error) {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if len(qm.workers) >= qm.maxWorkers {
		return nil, fmt.Errorf("no available workers")
	}

	worker := &Worker{
		ID:        taskID,
		TaskID:    taskID,
		IsActive:  true,
		StartedAt: time.Now(),
	}
	qm.workers[taskID] = worker
	return worker, nil
}

// ReleaseWorker 释放工作协程
func (qm *QueueManager) ReleaseWorker(taskID int64) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	delete(qm.workers, taskID)
}

// GetActiveWorkers 获取活跃工作协程
func (qm *QueueManager) GetActiveWorkers() []*Worker {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	workers := make([]*Worker, 0, len(qm.workers))
	for _, w := range qm.workers {
		if w.IsActive {
			workers = append(workers, w)
		}
	}
	return workers
}

// ==================== 队列统计 ====================

// QueueManagerStats 队列管理器统计
type QueueManagerStats struct {
	QueueCount    int         `json:"queue_count"`
	TotalItems    int         `json:"total_items"`
	ActiveWorkers int         `json:"active_workers"`
	MaxWorkers    int         `json:"max_workers"`
	QueueStats    map[string]*QueueStats `json:"queue_stats"`
}

// GetStats 获取统计信息
func (qm *QueueManager) GetStats() *QueueManagerStats {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	stats := &QueueManagerStats{
		QueueCount:    len(qm.queues),
		ActiveWorkers: len(qm.workers),
		MaxWorkers:    qm.maxWorkers,
		QueueStats:    make(map[string]*QueueStats),
	}

	for name, queue := range qm.queues {
		queueStats := queue.GetStats()
		stats.QueueStats[name] = queueStats
		stats.TotalItems += queueStats.Total
	}

	return stats
}

// ==================== 调度策略 ====================

// ScheduleStrategy 调度策略
type ScheduleStrategy string

const (
	StrategyFIFO     ScheduleStrategy = "fifo"     // 先进先出
	StrategyPriority ScheduleStrategy = "priority"  // 优先级调度
	StrategyFair     ScheduleStrategy = "fair"      // 公平调度（轮询）
)

// DequeueWithStrategy 按策略出队
func (qm *QueueManager) DequeueWithStrategy(strategy ScheduleStrategy) *QueueItem {
	switch strategy {
	case StrategyFIFO:
		return qm.Dequeue()
	case StrategyPriority:
		return qm.DequeueWithPriority()
	case StrategyFair:
		return qm.dequeueFair()
	default:
		return qm.Dequeue()
	}
}

// dequeueFair 公平调度（轮询各队列）
func (qm *QueueManager) dequeueFair() *QueueItem {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	// 轮询各队列
	for _, queue := range qm.queues {
		if queue.Len() > 0 {
			return queue.Dequeue()
		}
	}
	return nil
}

// ==================== 批量操作 ====================

// BatchEnqueue 批量入队
func (qm *QueueManager) BatchEnqueue(items []*QueueItem) {
	for _, item := range items {
		qm.Enqueue(item)
	}
}

// BatchDequeue 批量出队
func (qm *QueueManager) BatchDequeue(count int) []*QueueItem {
	result := make([]*QueueItem, 0, count)
	for i := 0; i < count; i++ {
		item := qm.Dequeue()
		if item == nil {
			break
		}
		result = append(result, item)
	}
	return result
}

// ==================== 队列维护 ====================

// CleanExpired 清理过期任务
func (qm *QueueManager) CleanExpired(maxAge time.Duration) int {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	cleaned := 0
	now := time.Now()

	for _, queue := range qm.queues {
		items := queue.GetAll()
		for _, item := range items {
			if now.Sub(item.CreatedAt) > maxAge {
				queue.Remove(item.TaskID)
				cleaned++
			}
		}
	}

	return cleaned
}

// Resize 调整工作协程数量
func (qm *QueueManager) Resize(newMaxWorkers int) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.maxWorkers = newMaxWorkers
}
