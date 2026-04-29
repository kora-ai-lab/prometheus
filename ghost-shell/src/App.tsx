import { useState, useEffect, useRef } from 'react'
import { AnimatePresence } from 'framer-motion'
import { listen } from '@tauri-apps/api/event'
import { Omnibox } from './components/Omnibox'
import { ExecutionModal } from './components/ExecutionModal'
import { Settings } from './components/Settings'
import { HealthBanner } from './components/HealthBanner'
import { execute, streamTask, cancelTask, getTask } from './lib/api'
import { addGoal } from './lib/history'

interface TaskState {
  id: string | null
  progress: string
  status: string
  result?: string
  error?: string
}

function App() {
  const [visible, setVisible] = useState(false)
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [task, setTask] = useState<TaskState>({ id: null, progress: '', status: '' })
  const abortRef = useRef<AbortController | null>(null)
  const goalRef = useRef<string>('')

  useEffect(() => {
    const unlisten = listen('shortcut-triggered', () => {
      setVisible(v => !v)
    })
    return () => { unlisten.then(fn => fn()) }
  }, [])

  useEffect(() => {
    const unlisten = listen('settings-opened', () => {
      setSettingsOpen(true)
    })
    return () => { unlisten.then(fn => fn()) }
  }, [])

  const resetTask = () => {
    setTask({ id: null, progress: '', status: '' })
    goalRef.current = ''
  }

  const handleSubmit = async (goal: string) => {
    addGoal(goal)
    setVisible(false)
    goalRef.current = goal
    setTask({ id: 'working', progress: 'Initializing...', status: 'running' })
    try {
      const taskId = await execute(goal)
      setTask({ id: taskId, progress: 'Thinking...', status: 'running' })
      const controller = streamTask(
        taskId,
        (data) => {
          setTask(prev => ({
            ...prev,
            progress: data.progress,
            status: data.status,
          }))
        },
        async () => {
          try {
            const full = await getTask(taskId)
            setTask(prev => ({
              ...prev,
              progress: full.progress,
              status: full.status,
              result: full.result,
              error: full.error,
            }))
          } catch {
            setTask(prev => ({ ...prev, status: 'failed', error: 'Failed to fetch result' }))
          }
        }
      )
      abortRef.current = controller
    } catch {
      setTask({ id: null, progress: '', status: 'failed', error: 'Core service offline' })
    }
  }

  const handleCancel = async () => {
    if (task.id && task.id !== 'working') {
      abortRef.current?.abort()
      try { await cancelTask(task.id) } catch {}
    }
    resetTask()
  }

  const handleRetry = () => {
    if (goalRef.current) handleSubmit(goalRef.current)
  }

  const handleRunAgain = () => {
    if (goalRef.current) handleSubmit(goalRef.current)
  }

  const handleClose = () => {
    resetTask()
  }

  return (
    <div className="w-full h-full bg-transparent">
      <HealthBanner />
      <Omnibox visible={visible} onSubmit={handleSubmit} onClose={() => setVisible(false)} />

      <AnimatePresence>
        {task.id && (
          <ExecutionModal
            task={task}
            onCancel={handleCancel}
            onRetry={handleRetry}
            onRunAgain={handleRunAgain}
            onClose={handleClose}
          />
        )}
      </AnimatePresence>

      <Settings visible={settingsOpen} onClose={() => setSettingsOpen(false)} />
    </div>
  )
}

export default App
