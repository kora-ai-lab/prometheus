import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { getHealth } from '../lib/api'

export function HealthBanner() {
  const [isOnline, setIsOnline] = useState(true)

  useEffect(() => {
    let mounted = true
    const check = async () => {
      const health = await getHealth()
      if (mounted) setIsOnline(health !== null)
    }
    check()
    const id = setInterval(check, 5000)
    return () => { mounted = false; clearInterval(id) }
  }, [])

  return (
    <AnimatePresence>
      {!isOnline && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          className="offline-banner"
        >
          Core service offline
        </motion.div>
      )}
    </AnimatePresence>
  )
}
