package metrics

import "sync/atomic"

type Collector struct {
	tasksStarted int64
	tasksDone    int64
	execs        int64
	llmCalls     int64

	TaskLatency *Histogram
	LlmLatency  *Histogram
	ExecLatency *Histogram
}

func New() *Collector {
	return &Collector{
		TaskLatency: NewHistogram([]float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000}),
		LlmLatency:  NewHistogram([]float64{100, 250, 500, 1000, 2000, 5000, 10000}),
		ExecLatency: NewHistogram([]float64{10, 50, 100, 500, 1000, 5000}),
	}
}

func (c *Collector) TaskStarted() { atomic.AddInt64(&c.tasksStarted, 1) }
func (c *Collector) TaskDone()    { atomic.AddInt64(&c.tasksDone, 1) }
func (c *Collector) Exec()        { atomic.AddInt64(&c.execs, 1) }
func (c *Collector) LLMCall()     { atomic.AddInt64(&c.llmCalls, 1) }

func (c *Collector) RecordTaskLatency(ms float64)  { c.TaskLatency.Observe(ms) }
func (c *Collector) RecordLLMLatency(ms float64)   { c.LlmLatency.Observe(ms) }
func (c *Collector) RecordExecLatency(ms float64) { c.ExecLatency.Observe(ms) }

func (c *Collector) Snapshot() map[string]interface{} {
	return map[string]interface{}{
		"tasks_started": atomic.LoadInt64(&c.tasksStarted),
		"tasks_done":    atomic.LoadInt64(&c.tasksDone),
		"execs":         atomic.LoadInt64(&c.execs),
		"llm_calls":     atomic.LoadInt64(&c.llmCalls),
		"task_latency":  c.TaskLatency.Sum(),
		"llm_latency":   c.LlmLatency.Sum(),
		"exec_latency":  c.ExecLatency.Sum(),
	}
}