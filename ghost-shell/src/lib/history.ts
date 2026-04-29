const HISTORY_KEY = 'prometheus-history'
const MAX_HISTORY = 20

export function addGoal(goal: string): void {
  const goals = getGoals()
  const filtered = goals.filter(g => g !== goal)
  filtered.unshift(goal)
  localStorage.setItem(HISTORY_KEY, JSON.stringify(filtered.slice(0, MAX_HISTORY)))
}

export function getGoals(): string[] {
  try {
    return JSON.parse(localStorage.getItem(HISTORY_KEY) || '[]')
  } catch {
    return []
  }
}
