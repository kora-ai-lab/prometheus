import type { Task, HealthResponse, ExecuteResponse } from './types'
import { getToken } from './auth'

const CORE_URL = import.meta.env.VITE_CORE_URL || 'http://localhost:8080'

async function request(path: string, options: RequestInit = {}): Promise<Response> {
  const headers: Record<string, string> = {
    ...(options.headers as Record<string, string>),
  }
  if (options.body) {
    headers['Content-Type'] = 'application/json'
  }
  return fetch(`${CORE_URL}${path}`, { ...options, headers })
}

export async function execute(goal: string): Promise<string> {
  const response = await request('/api/execute', {
    method: 'POST',
    body: JSON.stringify({ goal }),
    headers: { Authorization: `Bearer ${getToken()}` },
  })
  if (!response.ok) throw new Error(`Execute failed: ${response.statusText}`)
  const data: ExecuteResponse = await response.json()
  return data.task_id
}

export async function getTask(id: string): Promise<Task> {
  const response = await request(`/api/tasks/${id}`, {
    headers: { Authorization: `Bearer ${getToken()}` },
  })
  if (!response.ok) throw new Error(`Get task failed: ${response.statusText}`)
  return response.json()
}

export async function cancelTask(id: string): Promise<void> {
  const response = await request(`/api/tasks/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${getToken()}` },
  })
  if (!response.ok) throw new Error(`Cancel task failed: ${response.statusText}`)
}

export async function getHealth(): Promise<HealthResponse | null> {
  try {
    const response = await request('/api/health')
    if (!response.ok) return null
    return await response.json()
  } catch {
    return null
  }
}

export function streamTask(
  id: string,
  onEvent: (task: Task) => void,
  onDone: () => void
): AbortController {
  const controller = new AbortController()
  fetch(`${CORE_URL}/api/tasks/${id}/stream`, {
    headers: { Authorization: `Bearer ${getToken()}` },
    signal: controller.signal,
  }).then(async (response) => {
    const reader = response.body?.getReader()
    if (!reader) return
    const decoder = new TextDecoder()
    let buffer = ''
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''
      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = JSON.parse(line.slice(6)) as Task
          onEvent(data)
          if (['done', 'failed', 'cancelled'].includes(data.status)) {
            onDone()
            return
          }
        }
      }
    }
  }).catch(() => {})
  return controller
}
