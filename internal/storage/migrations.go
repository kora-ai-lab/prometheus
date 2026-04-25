package storage

type Migration struct {
	Version int
	SQL     string
}

var migrations = []Migration{
	{
		Version: 1,
		SQL: `
CREATE TABLE IF NOT EXISTS tasks (
  id TEXT PRIMARY KEY,
  goal TEXT NOT NULL,
  status TEXT NOT NULL,
  context_json TEXT,
  memory_json TEXT,
  blocked_reason TEXT,
  retries INTEGER NOT NULL DEFAULT 0,
  parse_errors INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS executions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id TEXT NOT NULL,
  command TEXT NOT NULL,
  stdout TEXT,
  stderr TEXT,
  exit_code INTEGER,
  duration_ms INTEGER,
  executed_at TEXT NOT NULL,
  FOREIGN KEY(task_id) REFERENCES tasks(id)
);
CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_exec_task ON executions(task_id);
`,
	},
	{
		Version: 2,
		SQL: `
CREATE TABLE IF NOT EXISTS user_prefs (
  key TEXT PRIMARY KEY,
  value TEXT,
  updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS learned_patterns (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  context TEXT,
  pattern TEXT,
  success_rate REAL DEFAULT 1.0,
  uses INTEGER DEFAULT 1,
  last_used TEXT
);
CREATE TABLE IF NOT EXISTS capabilities_cache (
  name TEXT PRIMARY KEY,
  installed INTEGER DEFAULT 0,
  version TEXT,
  path TEXT,
  metadata_json TEXT,
  installed_at TEXT
);
CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  date TEXT,
  summary TEXT,
  projects TEXT,
  technologies TEXT,
  created_at TEXT NOT NULL
);
`,
	},
	{
		Version: 3,
		SQL: `
CREATE TABLE IF NOT EXISTS log_index (
  date TEXT NOT NULL,
  month TEXT NOT NULL,
  summary_md TEXT,
  projects TEXT,
  technologies TEXT,
  goals TEXT,
  indexed_at DATETIME
);
CREATE VIRTUAL TABLE IF NOT EXISTS log_search
USING fts5(date, summary_md, projects, technologies, goals,
           content='log_index', content_rowid='rowid');
`,
	},
}
