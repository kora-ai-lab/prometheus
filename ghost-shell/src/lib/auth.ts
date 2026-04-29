const TOKEN_KEY = 'prometheus-token'

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export async function initToken(): Promise<void> {
  try {
    const { readTextFile } = await import('@tauri-apps/plugin-fs')
    const { appDataDir, join } = await import('@tauri-apps/api/path')
    const baseDir = await appDataDir()
    const tokenPath = await join(baseDir, 'token.txt')
    const token = await readTextFile(tokenPath)
    if (token) {
      setToken(token.trim())
    }
  } catch {
    // Fallback to localStorage
  }
}
