package service

import (
	"container/heap"
	"sync"
	"time"
)

// QueueItem 队列项
type QueueItem struct {
	TaskID      int64     `json:"task_id"`
	UserID      int64     `json:"user_id"`
	ContentID   int64     `json:"content_id"`
	ChannelID   int64     `json:"channel_id"`
	Priority    int32     `json:"priority"` // 0=低, 1=中, 2=高, 3=紧急
	ScheduledAt time.Time `json:"scheduled_at"`
	CreatedAt   time.Time `json:"created_at"`
	index       int       // heap内部索引
}

// priorityQueueHeap 堆实现（内部，不加锁）
type priorityQueueHeap struct {
	items   []*QueueItem
	indexes map[int64]*QueueItem
}

func (h *priorityQueueHeap) Len() int { return len(h.items) }

func (h *priorityQueueHeap) Less(i, j int) bool {
	if h.items[i].Priority != h.items[j].Priority {
		return h.items[i].Priority > h.items[j].Priority
	}
	return h.items[i].CreatedAt.Before(h.items[j].CreatedAt)
}

func (h *priorityQueueHeap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
	h.items[i].index = i
	h.items[j].index = j
}

func (h *priorityQueueHeap) Push(x interface{}) {
	item := x.(*QueueItem)
	item.index = len(h.items)
	h.items = append(h.items, item)
	h.indexes[item.TaskID] = item
}

func (h *priorityQueueHeap) Pop() interface{} {
	old := h.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	h.items = old[0 : n-1]
	delete(h.indexes, item.TaskID)
	return item
}

// PriorityQueue 优先级队列（线程安全）
type PriorityQueue struct {
	mu  sync.RWMutex
	heap *priorityQueueHeap
}

// NewPriorityQueue 创建优先级队列
func NewPriorityQueue() *PriorityQueue {
	h := &priorityQueueHeap{
		items:   make([]*QueueItem, 0),
		indexes: make(map[int64]*QueueItem),
	}
	heap.Init(h)
	return &PriorityQueue{heap: h}
}

// Len 返回队列长度
func (pq *PriorityQueue) Len() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return pq.heap.Len()
}

// Enqueue 入队
func (pq *PriorityQueue) Enqueue(item *QueueItem) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(pq.heap, item)
}

// Dequeue 出队（获取最高优先级任务）
func (pq *PriorityQueue) Dequeue() *QueueItem {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.heap.Len() == 0 {
		return nil
	}

	item := heap.Pop(pq.heap).(*QueueItem)
	return item
}

// Peek 查看队首元素（不移除）
func (pq *PriorityQueue) Peek() *QueueItem {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if pq.heap.Len() == 0 {
		return nil
	}
	return pq.heap.items[0]
}

// Contains 检查任务是否在队列中
func (pq *PriorityQueue) Contains(taskID int64) bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	_, exists := pq.heap.indexes[taskID]
	return exists
}

// Remove 移除指定任务
func (pq *PriorityQueue) Remove(taskID int64) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	item, exists := pq.heap.indexes[taskID]
	if !exists {
		return false
	}

	heap.Remove(pq.heap, item.index)
	return true
}

// UpdatePriority 更新任务优先级
func (pq *PriorityQueue) UpdatePriority(taskID int64, newPriority int32) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	item, exists := pq.heap.indexes[taskID]
	if !exists {
		return false
	}

	item.Priority = newPriority
	heap.Fix(pq.heap, item.index)
	return true
}

// GetAll 获取所有任务（按优先级排序）
func (pq *PriorityQueue) GetAll() []*QueueItem {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	result := make([]*QueueItem, len(pq.heap.items))
	copy(result, pq.heap.items)
	return result
}

// GetByPriority 获取指定优先级的任务
func (pq *PriorityQueue) GetByPriority(priority int32) []*QueueItem {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	result := make([]*QueueItem, 0)
	for _, item := range pq.heap.items {
		if item.Priority == priority {
			result = append(result, item)
		}
	}
	return result
}

// Clear 清空队列
func (pq *PriorityQueue) Clear() {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	pq.heap.items = make([]*QueueItem, 0)
	pq.heap.indexes = make(map[int64]*QueueItem)
	heap.Init(pq.heap)
}

// Stats 队列统计
type QueueStats struct {
	Total  int `json:"total"`
	Urgent int `json:"urgent"`
	High   int `json:"high"`
	Medium int `json:"medium"`
	Low    int `json:"low"`
}

// GetStats 获取队列统计
func (pq *PriorityQueue) GetStats() *QueueStats {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	stats := &QueueStats{
		Total: len(pq.heap.items),
	}

	for _, item := range pq.heap.items {
		switch item.Priority {
		case 3:
			stats.Urgent++
		case 2:
			stats.High++
		case 1:
			stats.Medium++
		case 0:
			stats.Low++
		}
	}

	return stats
}
