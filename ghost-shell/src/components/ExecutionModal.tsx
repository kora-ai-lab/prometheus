import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

interface ExecutionModalProps {
  task: { id: string | null; progress: string; status: string; result?: string; error?: string }
  onCancel: () => void
  onRetry: () => void
  onRunAgain: () => void
  onClose: () => void
}

const statusConfig: Record<string, { label: string; gradient: string; showProgress: boolean }> = {
  running: { label: 'Thinking...', gradient: 'from-cyan-500 to-purple-600', showProgress: true },
  blocked: { label: 'Blocked', gradient: 'from-yellow-500 to-orange-600', showProgress: false },
  done: { label: 'Done', gradient: 'from-green-500 to-emerald-600', showProgress: false },
  failed: { label: 'Failed', gradient: 'from-red-500 to-rose-600', showProgress: false },
  cancelled: { label: 'Cancelled', gradient: 'from-gray-500 to-gray-600', showProgress: false },
}

export function ExecutionModal({ task, onCancel, onRetry, onRunAgain, onClose }: ExecutionModalProps) {
  const [replyInput, setReplyInput] = useState('')

  const config = statusConfig[task.status] || statusConfig.running
  const isTerminal = ['done', 'failed', 'cancelled'].includes(task.status)

  return (
    <motion.div
      initial={{ opacity: 0, y: 10, scale: 0.98 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={{ opacity: 0, y: -10, scale: 0.98 }}
      transition={{ duration: 0.2 }}
      className="execution-modal"
      style={task.status === 'failed' ? { borderColor: 'rgba(255,60,60,0.3)' } : undefined}
    >
      <div className={`h-1 w-full rounded-t-lg bg-gradient-to-r ${config.gradient} mb-4`} />

      <div className="flex items-center justify-between mb-3">
        <span className={`text-sm font-medium ${
          task.status === 'failed' ? 'text-red-400' :
          task.status === 'done' ? 'text-green-400' :
          task.status === 'blocked' ? 'text-yellow-400' :
          task.status === 'cancelled' ? 'text-gray-500' :
          'text-cyan-400'
        }`}>
          {config.label}
        </span>
        {task.status === 'running' && (
          <button onClick={onCancel} className="text-xs text-gray-500 hover:text-white transition-colors">
            Cancel
          </button>
        )}
        {isTerminal && (
          <button onClick={onClose} className="text-xs text-gray-500 hover:text-white transition-colors">
            Close
          </button>
        )}
      </div>

      <AnimatePresence mode="wait">
        <motion.div
          key={task.status + task.progress}
          initial={{ opacity: 0, y: 5 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -5 }}
          transition={{ duration: 0.15 }}
        >
          {task.status === 'blocked' ? (
            <div className="space-y-3">
              <p className="text-sm text-gray-300">{task.progress}</p>
              <div className="flex gap-2">
                <input
                  type="text"
                  value={replyInput}
                  onChange={(e) => setReplyInput(e.target.value)}
                  placeholder="Type your reply..."
                  className="flex-1 bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-white placeholder-gray-500 outline-none focus:border-cyan-500/50"
                  onKeyDown={(e) => { if (e.key === 'Enter' && replyInput.trim()) { onClose() } }}
                />
                <button
                  onClick={onClose}
                  className="px-4 py-2 bg-cyan-600 hover:bg-cyan-500 text-white text-sm rounded-lg transition-colors"
                >
                  Reply
                </button>
              </div>
            </div>
          ) : task.status === 'done' ? (
            <div className="space-y-3">
              <p className="text-sm text-gray-300 whitespace-pre-wrap">{task.result || 'Task completed.'}</p>
              <div className="flex gap-2">
                <button
                  onClick={() => { navigator.clipboard?.writeText(task.result || '') }}
                  className="px-4 py-2 bg-white/5 hover:bg-white/10 text-white text-sm rounded-lg border border-white/10 transition-colors"
                >
                  Copy
                </button>
                <button
                  onClick={onRunAgain}
                  className="px-4 py-2 bg-cyan-600 hover:bg-cyan-500 text-white text-sm rounded-lg transition-colors"
                >
                  Run Again
                </button>
              </div>
            </div>
          ) : task.status === 'failed' ? (
            <div className="space-y-3">
              <p className="text-sm text-red-300">{task.error || 'An error occurred.'}</p>
              <button
                onClick={onRetry}
                className="px-4 py-2 bg-red-600 hover:bg-red-500 text-white text-sm rounded-lg transition-colors"
              >
                Retry
              </button>
            </div>
          ) : task.status === 'cancelled' ? (
            <p className="text-sm text-gray-500">Cancelled</p>
          ) : (
            <p className="text-sm text-gray-300">{task.progress || 'Thinking...'}</p>
          )}
        </motion.div>
      </AnimatePresence>

      {config.showProgress && (
        <div className="progress-bar mt-4">
          <motion.div
            className="progress-bar-fill"
            animate={{ x: ['-100%', '100%'] }}
            transition={{ duration: 1.5, repeat: Infinity, ease: 'linear' }}
            style={{ width: '50%' }}
          />
        </div>
      )}
    </motion.div>
  )
}
