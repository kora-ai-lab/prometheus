import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { setToken } from '../lib/auth'

interface SettingsProps {
  visible: boolean
  onClose: () => void
}

export function Settings({ visible, onClose }: SettingsProps) {
  const [coreUrl, setCoreUrl] = useState('http://localhost:8080')
  const [apiToken, setApiToken] = useState('')

  useEffect(() => {
    if (visible) {
      setCoreUrl(localStorage.getItem('prometheus-core-url') || 'http://localhost:8080')
      setApiToken('')
    }
  }, [visible])

  const handleSave = () => {
    localStorage.setItem('prometheus-core-url', coreUrl)
    if (apiToken) setToken(apiToken)
    onClose()
  }

  return (
    <AnimatePresence>
      {visible && (
        <motion.div
          initial={{ opacity: 0, y: 10, scale: 0.98 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          exit={{ opacity: 0, y: -10, scale: 0.98 }}
          transition={{ duration: 0.2 }}
          className="execution-modal"
        >
          <div className="h-1 w-full rounded-t-lg bg-gradient-to-r from-purple-500 to-pink-600 mb-4" />

          <div className="flex items-center justify-between mb-4">
            <span className="text-sm font-medium text-purple-400">Settings</span>
            <button onClick={onClose} className="text-xs text-gray-500 hover:text-white transition-colors">
              Close
            </button>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-xs text-gray-400 mb-1">Core URL</label>
              <input
                type="text"
                value={coreUrl}
                onChange={(e) => setCoreUrl(e.target.value)}
                className="w-full bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-white placeholder-gray-500 outline-none focus:border-purple-500/50"
              />
            </div>

            <div>
              <label className="block text-xs text-gray-400 mb-1">API Token</label>
              <input
                type="password"
                value={apiToken}
                onChange={(e) => setApiToken(e.target.value)}
                placeholder="Enter new token..."
                className="w-full bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-white placeholder-gray-500 outline-none focus:border-purple-500/50"
              />
            </div>

            <button
              onClick={handleSave}
              className="w-full px-4 py-2 bg-purple-600 hover:bg-purple-500 text-white text-sm rounded-lg transition-colors"
            >
              Save
            </button>
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}
