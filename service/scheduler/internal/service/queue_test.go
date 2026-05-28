package service

import (
	"sync"
	"testing"
	"time"
)

func TestPriorityQueue_Basic(t *testing.T) {
	pq := NewPriorityQueue()

	// 测试空队列
	if pq.Len() != 0 {
		t.Errorf("expected empty queue, got %d", pq.Len())
	}
	if pq.Dequeue() != nil {
		t.Error("expected nil from empty queue")
	}
	if pq.Peek() != nil {
		t.Error("expected nil from empty queue")
	}
}

func TestPriorityQueue_PriorityOrder(t *testing.T) {
	pq := NewPriorityQueue()

	// 添加不同优先级的任务
	pq.Enqueue(&QueueItem{TaskID: 1, Priority: 0, CreatedAt: time.Now()})
	pq.Enqueue(&QueueItem{TaskID: 2, Priority: 2, CreatedAt: time.Now()})
	pq.Enqueue(&QueueItem{TaskID: 3, Priority: 1, CreatedAt: time.Now()})
	pq.Enqueue(&QueueItem{TaskID: 4, Priority: 3, CreatedAt: time.Now()})

	if pq.Len() != 4 {
		t.Errorf("expected 4 items, got %d", pq.Len())
	}

	// 应该按优先级出队：3(紧急) -> 2(高) -> 1(中) -> 0(低)
	item := pq.Dequeue()
	if item.Priority != 3 {
		t.Errorf("expected priority 3, got %d", item.Priority)
	}

	item = pq.Dequeue()
	if item.Priority != 2 {
		t.Errorf("expected priority 2, got %d", item.Priority)
	}

	item = pq.Dequeue()
	if item.Priority != 1 {
		t.Errorf("expected priority 1, got %d", item.Priority)
	}

	item = pq.Dequeue()
	if item.Priority != 0 {
		t.Errorf("expected priority 0, got %d", item.Priority)
	}
}

func TestPriorityQueue_SamePriorityFIFO(t *testing.T) {
	pq := NewPriorityQueue()

	now := time.Now()
	pq.Enqueue(&QueueItem{TaskID: 1, Priority: 1, CreatedAt: now})
	pq.Enqueue(&QueueItem{TaskID: 2, Priority: 1, CreatedAt: now.Add(1 * time.Second)})
	pq.Enqueue(&QueueItem{TaskID: 3, Priority: 1, CreatedAt: now.Add(2 * time.Second)})

	// 同优先级按创建时间排序
	item := pq.Dequeue()
	if item.TaskID != 1 {
		t.Errorf("expected task 1, got %d", item.TaskID)
	}

	item = pq.Dequeue()
	if item.TaskID != 2 {
		t.Errorf("expected task 2, got %d", item.TaskID)
	}
}

func TestPriorityQueue_Contains(t *testing.T) {
	pq := NewPriorityQueue()

	pq.Enqueue(&QueueItem{TaskID: 100, Priority: 1})

	if !pq.Contains(100) {
		t.Error("expected to contain task 100")
	}
	if pq.Contains(200) {
		t.Error("should not contain task 200")
	}
}

func TestPriorityQueue_Remove(t *testing.T) {
	pq := NewPriorityQueue()

	pq.Enqueue(&QueueItem{TaskID: 1, Priority: 1})
	pq.Enqueue(&QueueItem{TaskID: 2, Priority: 2})
	pq.Enqueue(&QueueItem{TaskID: 3, Priority: 3})

	if !pq.Remove(2) {
		t.Error("expected to remove task 2")
	}
	if pq.Len() != 2 {
		t.Errorf("expected 2 items, got %d", pq.Len())
	}
	if pq.Contains(2) {
		t.Error("task 2 should be removed")
	}
}

func TestPriorityQueue_UpdatePriority(t *testing.T) {
	pq := NewPriorityQueue()

	pq.Enqueue(&QueueItem{TaskID: 1, Priority: 0})
	pq.Enqueue(&QueueItem{TaskID: 2, Priority: 1})
	pq.Enqueue(&QueueItem{TaskID: 3, Priority: 2})

	// 提升任务1的优先级到最高
	if !pq.UpdatePriority(1, 3) {
		t.Error("expected to update priority")
	}

	// 现在任务1应该排在最前面
	item := pq.Dequeue()
	if item.TaskID != 1 {
		t.Errorf("expected task 1, got %d", item.TaskID)
	}
	if item.Priority != 3 {
		t.Errorf("expected priority 3, got %d", item.Priority)
	}
}

func TestPriorityQueue_GetByPriority(t *testing.T) {
	pq := NewPriorityQueue()

	pq.Enqueue(&QueueItem{TaskID: 1, Priority: 0})
	pq.Enqueue(&QueueItem{TaskID: 2, Priority: 1})
	pq.Enqueue(&QueueItem{TaskID: 3, Priority: 1})
	pq.Enqueue(&QueueItem{TaskID: 4, Priority: 2})

	highPriority := pq.GetByPriority(1)
	if len(highPriority) != 2 {
		t.Errorf("expected 2 high priority items, got %d", len(highPriority))
	}
}

func TestPriorityQueue_Clear(t *testing.T) {
	pq := NewPriorityQueue()

	pq.Enqueue(&QueueItem{TaskID: 1, Priority: 1})
	pq.Enqueue(&QueueItem{TaskID: 2, Priority: 2})

	pq.Clear()

	if pq.Len() != 0 {
		t.Errorf("expected empty queue, got %d", pq.Len())
	}
}

func TestPriorityQueue_Stats(t *testing.T) {
	pq := NewPriorityQueue()

	pq.Enqueue(&QueueItem{TaskID: 1, Priority: 0})
	pq.Enqueue(&QueueItem{TaskID: 2, Priority: 1})
	pq.Enqueue(&QueueItem{TaskID: 3, Priority: 2})
	pq.Enqueue(&QueueItem{TaskID: 4, Priority: 3})
	pq.Enqueue(&QueueItem{TaskID: 5, Priority: 3})

	stats := pq.GetStats()
	if stats.Total != 5 {
		t.Errorf("expected 5 total, got %d", stats.Total)
	}
	if stats.Urgent != 2 {
		t.Errorf("expected 2 urgent, got %d", stats.Urgent)
	}
	if stats.High != 1 {
		t.Errorf("expected 1 high, got %d", stats.High)
	}
}

func TestQueueManager_Basic(t *testing.T) {
	qm := NewQueueManager(5)

	// 测试入队出队
	qm.Enqueue(&QueueItem{TaskID: 1, Priority: 1})
	qm.Enqueue(&QueueItem{TaskID: 2, Priority: 3})
	qm.Enqueue(&QueueItem{TaskID: 3, Priority: 2})

	item := qm.Dequeue()
	if item == nil {
		t.Fatal("expected non-nil item")
	}
	if item.Priority != 3 {
		t.Errorf("expected priority 3, got %d", item.Priority)
	}
}

func TestQueueManager_MultipleQueues(t *testing.T) {
	qm := NewQueueManager(5)

	qm.CreateQueue("urgent")
	qm.CreateQueue("normal")

	qm.EnqueueTo("urgent", &QueueItem{TaskID: 1, Priority: 3})
	qm.EnqueueTo("normal", &QueueItem{TaskID: 2, Priority: 1})

	item := qm.DequeueFrom("urgent")
	if item == nil || item.TaskID != 1 {
		t.Error("expected task 1 from urgent queue")
	}

	item = qm.DequeueFrom("normal")
	if item == nil || item.TaskID != 2 {
		t.Error("expected task 2 from normal queue")
	}
}

func TestQueueManager_DequeueWithPriority(t *testing.T) {
	qm := NewQueueManager(5)

	// 添加不同优先级的任务
	qm.Enqueue(&QueueItem{TaskID: 1, Priority: 0, CreatedAt: time.Now()})
	qm.Enqueue(&QueueItem{TaskID: 2, Priority: 1, CreatedAt: time.Now()})
	qm.Enqueue(&QueueItem{TaskID: 3, Priority: 3, CreatedAt: time.Now()})

	// 应该先出紧急任务
	item := qm.DequeueWithPriority()
	if item.Priority != 3 {
		t.Errorf("expected priority 3, got %d", item.Priority)
	}
}

func TestQueueManager_WorkerManagement(t *testing.T) {
	qm := NewQueueManager(2)

	// 获取工作协程
	w1, err := qm.AcquireWorker(100)
	if err != nil {
		t.Fatalf("acquire worker failed: %v", err)
	}
	if !w1.IsActive {
		t.Error("expected active worker")
	}

	w2, err := qm.AcquireWorker(200)
	if err != nil {
		t.Fatalf("acquire worker failed: %v", err)
	}
	_ = w2

	// 超过最大工作协程数
	_, err = qm.AcquireWorker(300)
	if err == nil {
		t.Error("expected error when exceeding max workers")
	}

	// 释放工作协程
	qm.ReleaseWorker(100)

	// 现在可以获取新的工作协程
	_, err = qm.AcquireWorker(300)
	if err != nil {
		t.Fatalf("acquire worker failed after release: %v", err)
	}
}

func TestQueueManager_Stats(t *testing.T) {
	qm := NewQueueManager(5)

	qm.Enqueue(&QueueItem{TaskID: 1, Priority: 1})
	qm.Enqueue(&QueueItem{TaskID: 2, Priority: 2})
	qm.AcquireWorker(1)

	stats := qm.GetStats()
	if stats.QueueCount != 1 {
		t.Errorf("expected 1 queue, got %d", stats.QueueCount)
	}
	if stats.TotalItems != 2 {
		t.Errorf("expected 2 items, got %d", stats.TotalItems)
	}
	if stats.ActiveWorkers != 1 {
		t.Errorf("expected 1 worker, got %d", stats.ActiveWorkers)
	}
}

func TestQueueManager_BatchOperations(t *testing.T) {
	qm := NewQueueManager(5)

	items := []*QueueItem{
		{TaskID: 1, Priority: 0},
		{TaskID: 2, Priority: 1},
		{TaskID: 3, Priority: 2},
	}

	qm.BatchEnqueue(items)

	result := qm.BatchDequeue(2)
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}

	// 第三个应该还在队列中
	if qm.Dequeue() == nil {
		t.Error("expected 1 item remaining")
	}
}

func TestQueueManager_CleanExpired(t *testing.T) {
	qm := NewQueueManager(5)

	qm.Enqueue(&QueueItem{TaskID: 1, Priority: 1, CreatedAt: time.Now().Add(-2 * time.Hour)})
	qm.Enqueue(&QueueItem{TaskID: 2, Priority: 1, CreatedAt: time.Now()})

	cleaned := qm.CleanExpired(1 * time.Hour)
	if cleaned != 1 {
		t.Errorf("expected 1 cleaned, got %d", cleaned)
	}

	if qm.Dequeue() == nil {
		t.Error("expected 1 item remaining")
	}
}

func TestQueueManager_Resize(t *testing.T) {
	qm := NewQueueManager(5)

	qm.Resize(10)

	stats := qm.GetStats()
	if stats.MaxWorkers != 10 {
		t.Errorf("expected max workers 10, got %d", stats.MaxWorkers)
	}
}

func TestPriorityQueue_Concurrent(t *testing.T) {
	pq := NewPriorityQueue()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			pq.Enqueue(&QueueItem{TaskID: int64(id), Priority: int32(id % 4)})
		}(i)
	}
	wg.Wait()

	if pq.Len() != 100 {
		t.Errorf("expected 100 items, got %d", pq.Len())
	}

	// 验证出队顺序
	lastPriority := int32(4)
	for pq.Len() > 0 {
		item := pq.Dequeue()
		if item.Priority > lastPriority {
			t.Error("items not in priority order")
		}
		lastPriority = item.Priority
	}
}

func TestQueueItem_Fields(t *testing.T) {
	item := &QueueItem{
		TaskID:      123,
		UserID:      456,
		ContentID:   789,
		ChannelID:   101,
		Priority:    2,
		ScheduledAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	if item.TaskID != 123 {
		t.Errorf("expected 123, got %d", item.TaskID)
	}
	if item.Priority != 2 {
		t.Errorf("expected 2, got %d", item.Priority)
	}
}
