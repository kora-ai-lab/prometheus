export interface Task {
  id: string
  goal: string
  status: 'running' | 'blocked' | 'done' | 'failed' | 'cancelled'
  progress: string
  result: string
  error: string
  createdAt: string
  updatedAt: string
}

export interface HealthResponse {
  status: string
  version: string
  uptime: string
}

export interface ExecuteResponse {
  task_id: string
}
