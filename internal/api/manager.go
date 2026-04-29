package api

import (
	"context"
	"errors"
	"sync"

	"github.com/kora-ai-lab/prometheus/internal/task"
)

type TaskManager struct {
	mu      sync.RWMutex
	tasks   map[string]*task.Task
	active  map[string]context.CancelFunc
	newDeps func() *task.TaskDeps
	runFn   func(context.Context, *task.Task, *task.TaskDeps) error
}

func NewTaskManager(depsFactory func() *task.TaskDeps) *TaskManager {
	return &TaskManager{
		tasks:   make(map[string]*task.Task),
		active:  make(map[string]context.CancelFunc),
		newDeps: depsFactory,
		runFn:   func(ctx context.Context, t *task.Task, deps *task.TaskDeps) error { return t.Run(ctx, deps) },
	}
}

func (m *TaskManager) WithRunFn(fn func(context.Context, *task.Task, *task.TaskDeps) error) {
	m.runFn = fn
}

func (m *TaskManager) Submit(goal string) string {
	t := task.New(goal)
	t.SetProgress("Initializing...")

	ctx, cancel := context.WithCancel(context.Background())

	m.mu.Lock()
	m.tasks[t.ID] = t
	m.active[t.ID] = cancel
	m.mu.Unlock()

	go func() {
		defer cancel()

		deps := m.newDeps()
		if err := m.runFn(ctx, t, deps); err != nil {
			m.mu.Lock()
			if t.Status == task.StatusRunning {
				t.Status = task.StatusFailed
				t.Error = err.Error()
				t.SetProgress("Failed: " + err.Error())
			}
			m.mu.Unlock()
		}

		m.mu.Lock()
		delete(m.active, t.ID)
		m.mu.Unlock()
	}()

	return t.ID
}

func (m *TaskManager) GetStatus(id string) (*task.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, ok := m.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}

	return t, nil
}

func (m *TaskManager) ListActive() []*task.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*task.Task
	for _, t := range m.tasks {
		if !isTerminal(t.Status) {
			result = append(result, t)
		}
	}
	return result
}

func (m *TaskManager) Cancel(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tasks[id]
	if !ok {
		return errors.New("task not found")
	}

	if isTerminal(t.Status) {
		return errors.New("task is already in terminal state")
	}

	cancel, ok := m.active[id]
	if !ok {
		return errors.New("task is not running")
	}

	cancel()
	t.Status = task.StatusCancelled
	t.SetProgress("Cancelled")
	delete(m.active, id)

	deps := m.newDeps()
	if err := deps.TaskStore.Save(t); err != nil {
		return err
	}

	return nil
}

func isTerminal(s task.TaskStatus) bool {
	return s == task.StatusDone || s == task.StatusFailed || s == task.StatusCancelled
}
