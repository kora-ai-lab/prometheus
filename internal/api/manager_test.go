package api

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/task"
)

type mockTaskStore struct {
	mu     sync.Mutex
	saved  []*task.Task
	failOn bool
}

func (m *mockTaskStore) Save(t *task.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failOn {
		return nil
	}
	m.saved = append(m.saved, t)
	return nil
}

func (m *mockTaskStore) lastSaved() *task.Task {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.saved) == 0 {
		return nil
	}
	return m.saved[len(m.saved)-1]
}

func newDepsFactory(store *mockTaskStore) func() *task.TaskDeps {
	return func() *task.TaskDeps {
		return &task.TaskDeps{
			TaskStore: store,
		}
	}
}

func noopRun(ctx context.Context, t *task.Task, deps *task.TaskDeps) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestSubmitAndGetStatus(t *testing.T) {
	store := &mockTaskStore{}
	mgr := NewTaskManager(newDepsFactory(store))
	mgr.WithRunFn(noopRun)

	id := mgr.Submit("test goal")

	if id == "" {
		t.Fatal("expected non-empty task ID")
	}

	tk, err := mgr.GetStatus(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tk.Goal != "test goal" {
		t.Errorf("expected goal 'test goal', got %q", tk.Goal)
	}

	if tk.Status != task.StatusRunning {
		t.Errorf("expected status running, got %q", tk.Status)
	}

	if tk.Progress == "" {
		t.Error("expected progress to be set")
	}

	_, err = mgr.GetStatus("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestCancelFlow(t *testing.T) {
	store := &mockTaskStore{}
	mgr := NewTaskManager(newDepsFactory(store))
	mgr.WithRunFn(noopRun)

	id := mgr.Submit("cancel me")

	err := mgr.Cancel(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tk, err := mgr.GetStatus(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tk.Status != task.StatusCancelled {
		t.Errorf("expected status cancelled, got %q", tk.Status)
	}

	if tk.Progress != "Cancelled" {
		t.Errorf("expected progress 'Cancelled', got %q", tk.Progress)
	}

	err = mgr.Cancel("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}

	err = mgr.Cancel(id)
	if err == nil {
		t.Error("expected error when cancelling already cancelled task")
	}
}

func TestListActive(t *testing.T) {
	store := &mockTaskStore{}
	mgr := NewTaskManager(newDepsFactory(store))
	mgr.WithRunFn(noopRun)

	id1 := mgr.Submit("task 1")
	id2 := mgr.Submit("task 2")

	active := mgr.ListActive()
	if len(active) != 2 {
		t.Errorf("expected 2 active tasks, got %d", len(active))
	}

	err := mgr.Cancel(id1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	active = mgr.ListActive()
	if len(active) != 1 {
		t.Errorf("expected 1 active task, got %d", len(active))
	}

	if active[0].ID != id2 {
		t.Errorf("expected remaining task to be %q, got %q", id2, active[0].ID)
	}
}

func TestConcurrentSubmissions(t *testing.T) {
	store := &mockTaskStore{}
	mgr := NewTaskManager(newDepsFactory(store))
	mgr.WithRunFn(noopRun)

	const count = 50
	var wg sync.WaitGroup
	ids := make([]string, count)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ids[i] = mgr.Submit("concurrent task " + string(rune('A'+i)))
		}(i)
	}

	wg.Wait()

	for _, id := range ids {
		if id == "" {
			t.Error("expected non-empty task ID from concurrent submission")
		}

		tk, err := mgr.GetStatus(id)
		if err != nil {
			t.Errorf("unexpected error for task %s: %v", id, err)
		}

		if tk == nil {
			continue
		}

		if tk.Status != task.StatusRunning {
			t.Errorf("expected status running for task %s, got %q", id, tk.Status)
		}
	}

	active := mgr.ListActive()
	if len(active) != count {
		t.Errorf("expected %d active tasks, got %d", count, len(active))
	}
}

func TestSetProgress(t *testing.T) {
	tk := task.New("test")

	tk.SetProgress("Thinking...")
	if tk.Progress != "Thinking..." {
		t.Errorf("expected progress 'Thinking...', got %q", tk.Progress)
	}

	if !tk.UpdatedAt.After(tk.CreatedAt) && !tk.UpdatedAt.Equal(tk.CreatedAt) {
		t.Error("expected UpdatedAt to be updated")
	}
}

func TestTaskFields(t *testing.T) {
	tk := task.New("goal")

	if tk.Result != "" {
		t.Errorf("expected empty Result, got %q", tk.Result)
	}

	if tk.Progress == "" && tk.Status == task.StatusRunning {
		tk.SetProgress("Initializing...")
	}

	if tk.Error != "" {
		t.Errorf("expected empty Error, got %q", tk.Error)
	}
}

func TestTerminalStates(t *testing.T) {
	store := &mockTaskStore{}
	mgr := NewTaskManager(newDepsFactory(store))

	mgr.WithRunFn(func(ctx context.Context, tk *task.Task, deps *task.TaskDeps) error {
		tk.SetStatus(task.StatusDone)
		return nil
	})

	id := mgr.Submit("test")

	time.Sleep(50 * time.Millisecond)

	active := mgr.ListActive()
	for _, tk := range active {
		if tk.ID == id {
			t.Error("done task should not appear in active list")
		}
	}

	mgr2 := NewTaskManager(newDepsFactory(store))
	mgr2.WithRunFn(func(ctx context.Context, tk *task.Task, deps *task.TaskDeps) error {
		tk.SetStatus(task.StatusFailed)
		tk.Error = "test error"
		return nil
	})

	id2 := mgr2.Submit("test2")
	time.Sleep(50 * time.Millisecond)

	active = mgr2.ListActive()
	for _, tk := range active {
		if tk.ID == id2 {
			t.Error("failed task should not appear in active list")
		}
	}
}

func TestCancelSavesTask(t *testing.T) {
	store := &mockTaskStore{}
	mgr := NewTaskManager(newDepsFactory(store))
	mgr.WithRunFn(noopRun)

	id := mgr.Submit("save test")

	err := mgr.Cancel(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	saved := store.lastSaved()
	if saved == nil {
		t.Fatal("expected task to be saved on cancel")
	}

	if saved.Status != task.StatusCancelled {
		t.Errorf("expected saved task status cancelled, got %q", saved.Status)
	}
}
