package api

import (
	"context"
	"testing"

	"github.com/kora-ai-lab/prometheus/internal/task"
)

func BenchmarkTaskManager_Submit(b *testing.B) {
	m := NewTaskManager(func() *task.TaskDeps {
		return &task.TaskDeps{}
	})
	m.WithRunFn(func(ctx context.Context, tsk *task.Task, deps *task.TaskDeps) error {
		return nil
	})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.Submit("benchmark task goal")
	}
}

func BenchmarkTaskManager_GetStatus(b *testing.B) {
	m := NewTaskManager(func() *task.TaskDeps {
		return &task.TaskDeps{}
	})
	m.WithRunFn(func(ctx context.Context, tsk *task.Task, deps *task.TaskDeps) error {
		return nil
	})
	id := m.Submit("test goal")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.GetStatus(id)
	}
}

func BenchmarkTaskManager_ListActive(b *testing.B) {
	m := NewTaskManager(func() *task.TaskDeps {
		return &task.TaskDeps{}
	})
	m.WithRunFn(func(ctx context.Context, tsk *task.Task, deps *task.TaskDeps) error {
		return nil
	})
	for n := 0; n < 100; n++ {
		m.Submit("test goal")
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.ListActive()
	}
}
