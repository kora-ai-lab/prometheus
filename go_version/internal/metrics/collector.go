package metrics

import "sync/atomic"

type Collector struct {
	tasksStarted int64
	tasksDone    int64
	execs        int64
	llmCalls     int64
}

func New() *Collector {
	return &Collector{}
}

func (c *Collector) TaskStarted() { atomic.AddInt64(&c.tasksStarted, 1) }
func (c *Collector) TaskDone()    { atomic.AddInt64(&c.tasksDone, 1) }
func (c *Collector) Exec()        { atomic.AddInt64(&c.execs, 1) }
func (c *Collector) LLMCall()     { atomic.AddInt64(&c.llmCalls, 1) }

func (c *Collector) Snapshot() map[string]int64 {
	return map[string]int64{
		"tasks_started": atomic.LoadInt64(&c.tasksStarted),
		"tasks_done":    atomic.LoadInt64(&c.tasksDone),
		"execs":         atomic.LoadInt64(&c.execs),
		"llm_calls":     atomic.LoadInt64(&c.llmCalls),
	}
}
