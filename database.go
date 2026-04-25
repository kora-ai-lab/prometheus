package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

const (
	TaskStatusPending    string = "pending"
	TaskStatusRunning    string = "in_progress"
	TaskStatusBlocked   string = "blocked"
	TaskStatusCompleted string = "completed"
	TaskStatusFailed   string = "failed"
)

type DB struct {
	conn *sql.DB
}

func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			goal TEXT,
			status TEXT DEFAULT 'pending',
			retry_count INTEGER DEFAULT 0,
			error_message TEXT,
			blocked_reason TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			started_at DATETIME,
			completed_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task_id TEXT,
			step INTEGER,
			command TEXT,
			stdout TEXT,
			stderr TEXT,
			exit_code INTEGER,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(task_id) REFERENCES tasks(id)
		);`,
		`CREATE TABLE IF NOT EXISTS blocks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task_id TEXT,
			reason TEXT,
			user_input TEXT,
			resolved_at DATETIME,
			FOREIGN KEY(task_id) REFERENCES tasks(id)
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return nil, err
		}
	}

	return &DB{conn: db}, nil
}

func (db *DB) CreateTask(id, goal string) error {
	_, err := db.conn.Exec(
		"INSERT INTO tasks (id, goal, status) VALUES (?, ?, ?)",
		id, goal, TaskStatusPending,
	)
	return err
}

func (db *DB) GetTask(id string) (goal string, status string, retryCount int, err error) {
	err = db.conn.QueryRow(
		"SELECT goal, status, retry_count FROM tasks WHERE id = ?",
		id,
	).Scan(&goal, &status, &retryCount)
	return
}

func (db *DB) StartTask(id string) error {
	_, err := db.conn.Exec(
		"UPDATE tasks SET status = ?, started_at = ? WHERE id = ?",
		TaskStatusRunning, time.Now().Format(time.RFC3339), id,
	)
	return err
}

func (db *DB) UpdateTaskStatus(id, status string) error {
	_, err := db.conn.Exec(
		"UPDATE tasks SET status = ? WHERE id = ?",
		status, id,
	)
	return err
}

func (db *DB) CompleteTask(id string) error {
	_, err := db.conn.Exec(
		"UPDATE tasks SET status = ?, completed_at = ? WHERE id = ?",
		TaskStatusCompleted, time.Now().Format(time.RFC3339), id,
	)
	return err
}

func (db *DB) FailTask(id, errMsg string) error {
	_, err := db.conn.Exec(
		"UPDATE tasks SET status = ?, error_message = ?, completed_at = ? WHERE id = ?",
		TaskStatusFailed, errMsg, time.Now().Format(time.RFC3339), id,
	)
	return err
}

func (db *DB) BlockTask(id, reason string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE tasks SET status = ?, blocked_reason = ? WHERE id = ?",
		TaskStatusBlocked, reason, id,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO blocks (task_id, reason) VALUES (?, ?)",
		id, reason,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) UnblockTask(id, userInput string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE tasks SET status = ?, blocked_reason = NULL WHERE id = ?",
		TaskStatusRunning, id,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE blocks SET user_input = ?, resolved_at = ? WHERE task_id = ? AND resolved_at IS NULL",
		userInput, time.Now().Format(time.RFC3339), id,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) IncrementRetry(id string) error {
	_, err := db.conn.Exec(
		"UPDATE tasks SET retry_count = retry_count + 1 WHERE id = ?",
		id,
	)
	return err
}

func (db *DB) GetBlockReason(id string) (string, error) {
	var reason string
	err := db.conn.QueryRow(
		"SELECT blocked_reason FROM tasks WHERE id = ?",
		id,
	).Scan(&reason)
	return reason, err
}

func (db *DB) LogExecution(taskID string, step int, cmd, stdout, stderr string, exitCode int) error {
	_, err := db.conn.Exec(
		"INSERT INTO logs (task_id, step, command, stdout, stderr, exit_code) VALUES (?, ?, ?, ?, ?, ?)",
		taskID, step, cmd, stdout, stderr, exitCode,
	)
	return err
}

func (db *DB) GetTaskLogs(taskID string) ([]string, error) {
	rows, err := db.conn.Query(
		"SELECT command FROM logs WHERE task_id = ? ORDER BY step",
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []string
	for rows.Next() {
		var cmd string
		if err := rows.Scan(&cmd); err != nil {
			return nil, err
		}
		logs = append(logs, cmd)
	}
	return logs, nil
}

func (db *DB) GetTaskCount() (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	return count, err
}

func (db *DB) GetCompletedCount() (int, error) {
	var count int
	err := db.conn.QueryRow(
		"SELECT COUNT(*) FROM tasks WHERE status = ?",
		TaskStatusCompleted,
	).Scan(&count)
	return count, err
}

func (db *DB) Close() error {
	return db.conn.Close()
}