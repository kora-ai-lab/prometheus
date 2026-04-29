//go:build integration

package api

import (
	"sync"
	"testing"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/task"
)

func TestTaskLifecycle(t *testing.T) {
	t.Parallel()
	m := NewTaskManager(func() *task.TaskDeps {
		return &task.TaskDeps{}
	})

	m.WithRunFn(func(ctx context.Context, t *task.Task, deps *task.TaskDeps) error {
		t.SetStatus(task.StatusDone)
		t.Result = "completed"
		return nil
	})

	id := m.Submit("test goal")
	time.Sleep(50 * time.Millisecond)

	status, err := m.GetStatus(id)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Status != task.StatusDone {
		t.Errorf("expected status done, got %s", status.Status)
	}
}

func TestTaskLifecycle_Progress(t *testing.T) {
	t.Parallel()
	m := NewTaskManager(func() *task.TaskDeps {
		return &task.TaskDeps{}
	})

	m.WithRunFn(func(ctx context.Context, t *task.Task, deps *task.TaskDeps) error {
		t.SetProgress("working...")
		time.Sleep(10 * time.Millisecond)
		t.SetStatus(task.StatusDone)
		return nil
	})

	id := m.Submit("test goal")
	time.Sleep(20 * time.Millisecond)

	status, _ := m.GetStatus(id)
	if status.Progress == "" {
		t.Error("expected progress to be set")
	}
}

func TestConcurrentTaskSubmission(t *testing.T) {
	t.Parallel()
	m := NewTaskManager(func() *task.TaskDeps {
		return &task.TaskDeps{}
	})

	m.WithRunFn(func(ctx context.Context, t *task.Task, deps *task.TaskDeps) error {
		t.SetStatus(task.StatusDone)
		return nil
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Submit("concurrent task")
		}()
	}
	wg.Wait()

	tasks := m.ListAll()
	if len(tasks) != 10 {
		t.Errorf("expected 10 tasks, got %d", len(tasks))
	}
}

func TestTaskCancel(t *testing.T) {
	t.Parallel()
	m := NewTaskManager(func() *task.TaskDeps {
		return &task.TaskDeps{}
	})

	m.WithRunFn(func(ctx context.Context, t *task.Task, deps *task.TaskDeps) error {
		<-ctx.Done()
		return ctx.Err()
	})

	id := m.Submit("long running task")
	time.Sleep(10 * time.Millisecond)

	err := m.Cancel(id)
	if err != nil {
		t.Errorf("Cancel failed: %v", err)
	}

	status, _ := m.GetStatus(id)
	if status.Status != task.StatusCancelled {
		t.Errorf("expected cancelled status, got %s", status.Status)
	}
}

func TestListActiveTasks(t *testing.T) {
	t.Parallel()
	m := NewTaskManager(func() *task.TaskDeps {
		return &task.TaskDeps{}
	})

	m.WithRunFn(func(ctx context.Context, t *task.Task, deps *task.TaskDeps) error {
		time.Sleep(100 * time.Millisecond)
		t.SetStatus(task.StatusDone)
		return nil
	})

	for i := 0; i < 5; i++ {
		m.Submit("active task")
	}
	time.Sleep(10 * time.Millisecond)

	active := m.ListActive()
	if len(active) != 5 {
		t.Errorf("expected 5 active tasks, got %d", len(active))
	}
}
