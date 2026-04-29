import { useState, useEffect } from 'react'
import { getHealth } from './api'
import type { HealthResponse } from './types'

export type HealthStatus = 'online' | 'offline'

export function useHealth() {
  const [health, setHealth] = useState<HealthStatus>('offline')
  const [healthData, setHealthData] = useState<HealthResponse | null>(null)

  useEffect(() => {
    const poll = async () => {
      const res = await getHealth()
      if (res) {
        setHealth('online')
        setHealthData(res)
      } else {
        setHealth('offline')
        setHealthData(null)
      }
    }
    poll()
    const interval = setInterval(poll, 5000)
    return () => clearInterval(interval)
  }, [])

  return { health, healthData }
}
