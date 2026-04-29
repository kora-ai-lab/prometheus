import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { listen } from '@tauri-apps/api/event'
import { Omnibox } from './components/Omnibox'
import { execute } from './lib/api'

interface TaskState {
  id: string | null
  progress: string
  status: string
}

function App() {
  const [visible, setVisible] = useState(false)
  const [task, setTask] = useState<TaskState>({ id: null, progress: '', status: '' })

  useEffect(() => {
    const unlisten = listen('shortcut-triggered', () => {
      setVisible(v => !v)
    })
    return () => { unlisten.then(fn => fn()) }
  }, [])

  const handleSubmit = async (goal: string) => {
    setVisible(false)
    setTask({ id: 'working', progress: 'Initializing...', status: 'running' })
    try {
      const taskId = await execute(goal)
      setTask({ id: taskId, progress: 'Thinking...', status: 'running' })
    } catch {
      setTask({ id: null, progress: 'Core service offline', status: 'failed' })
    }
  }

  const handleCloseModal = () => {
    setTask({ id: null, progress: '', status: '' })
  }

  return (
    <div className="w-full h-full bg-transparent">
      <Omnibox visible={visible} onSubmit={handleSubmit} onClose={() => setVisible(false)} />

      <AnimatePresence>
        {task.id && (
          <motion.div
            initial={{ opacity: 0, y: 10, scale: 0.98 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -10, scale: 0.98 }}
            transition={{ duration: 0.2 }}
            className="execution-modal"
          >
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm font-medium text-cyan-400">
                {task.status === 'running' ? 'Working...' : task.status === 'done' ? 'Done' : task.status}
              </span>
              {task.status === 'running' && (
                <button
                  onClick={handleCloseModal}
                  className="text-xs text-gray-500 hover:text-white transition-colors"
                >
                  Cancel
                </button>
              )}
            </div>
            <p className="text-sm text-gray-300">{task.progress}</p>
            {task.status === 'running' && (
              <div className="progress-bar">
                <div className="progress-bar-fill" />
              </div>
            )}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}

export default App
