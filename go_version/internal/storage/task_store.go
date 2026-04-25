package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/prometheus-dev/prometheus/internal/llm"
	"github.com/prometheus-dev/prometheus/internal/task"
)

type TaskStore interface {
	Save(*task.Task) error
	Load(id string) (*task.Task, error)
	Close() error
}

type Store struct {
	db *sql.DB
}

func Open(home string) (*Store, error) {
	db, err := sql.Open("sqlite", filepath.Join(home, "prometheus.db"))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA foreign_keys=ON",
	} {
		if _, err := db.ExecContext(context.Background(), pragma); err != nil {
			return nil, err
		}
	}

	s := &Store{db: db}
	if err := s.runMigrations(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) runMigrations() error {
	var currentVersion int
	s.db.QueryRow("SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&currentVersion)

	for _, migration := range migrations {
		if migration.Version <= currentVersion {
			continue
		}
		if _, err := s.db.Exec(migration.SQL); err != nil {
			return err
		}
		if _, err := s.db.Exec(
			"INSERT OR REPLACE INTO schema_version(version) VALUES(?)",
			migration.Version,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) Save(t *task.Task) error {
	contextJSON, err := json.Marshal(t.Context)
	if err != nil {
		return err
	}
	memoryJSON, err := json.Marshal(t.Memory)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		`INSERT INTO tasks (id, goal, status, context_json, memory_json, blocked_reason, retries, parse_errors, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   goal=excluded.goal,
		   status=excluded.status,
		   context_json=excluded.context_json,
		   memory_json=excluded.memory_json,
		   blocked_reason=excluded.blocked_reason,
		   retries=excluded.retries,
		   parse_errors=excluded.parse_errors,
		   updated_at=excluded.updated_at`,
		t.ID,
		t.Goal,
		t.Status,
		string(contextJSON),
		string(memoryJSON),
		t.BlockedReason,
		t.Retries,
		t.ParseErrors,
		t.CreatedAt.Format(time.RFC3339),
		t.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (s *Store) Load(id string) (*task.Task, error) {
	row := s.db.QueryRow(
		`SELECT id, goal, status, context_json, memory_json, blocked_reason, retries, parse_errors, created_at, updated_at
		 FROM tasks WHERE id = ?`,
		id,
	)

	var t task.Task
	var contextJSON string
	var memoryJSON string
	var createdAt string
	var updatedAt string
	if err := row.Scan(
		&t.ID,
		&t.Goal,
		&t.Status,
		&contextJSON,
		&memoryJSON,
		&t.BlockedReason,
		&t.Retries,
		&t.ParseErrors,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	if contextJSON != "" {
		if err := json.Unmarshal([]byte(contextJSON), &t.Context); err != nil {
			return nil, err
		}
	}
	if memoryJSON != "" {
		if err := json.Unmarshal([]byte(memoryJSON), &t.Memory); err != nil {
			return nil, err
		}
	}
	if t.Memory == nil {
		t.Memory = map[string]any{}
	}
	t.MaxRetries = 5
	t.MaxParseErrors = 3
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	if len(t.Context) == 0 {
		t.Context = []llm.Message{}
	}
	return &t, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
