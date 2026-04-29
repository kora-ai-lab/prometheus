const API_BASE = '';
let authToken = '';
let refreshInterval = null;

const elements = {
    goalInput: null,
    executeBtn: null,
    output: null,
    tasksList: null,
    healthBar: null
};

document.addEventListener('DOMContentLoaded', () => {
    elements.goalInput = document.getElementById('goal');
    elements.executeBtn = document.getElementById('execute-btn');
    elements.output = document.getElementById('output');
    elements.tasksList = document.getElementById('tasks-list');
    elements.healthBar = document.getElementById('health-bar');

    if (elements.goalInput) {
        elements.goalInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') execute();
        });
    }

    if (elements.executeBtn) {
        elements.executeBtn.addEventListener('click', execute);
    }

    loadAuthToken();
    checkHealth();
    refreshHealth();
});

function loadAuthToken() {
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token');
    if (token) {
        authToken = token;
        localStorage.setItem('prometheus_token', token);
    } else {
        authToken = localStorage.getItem('prometheus_token') || '';
    }
}

function getHeaders() {
    const headers = { 'Content-Type': 'application/json' };
    if (authToken) {
        headers['Authorization'] = `Bearer ${authToken}`;
    }
    return headers;
}

async function execute() {
    const goal = elements.goalInput.value.trim();
    if (!goal) return;

    setOutput('Submitting task...', false);
    setButtonLoading(true);

    try {
        const res = await fetch(`${API_BASE}/api/execute`, {
            method: 'POST',
            headers: getHeaders(),
            body: JSON.stringify({ Goal: goal })
        });

        if (!res.ok) {
            const err = await res.json().catch(() => ({ error: 'Request failed' }));
            setOutput(`Error: ${err.error}`, true);
            return;
        }

        const data = await res.json();
        setOutput(`Task submitted: ${data.task_id}\nMonitoring progress...`, false);
        streamTaskStatus(data.task_id);
    } catch (e) {
        setOutput(`Error: ${e.message}`, true);
    } finally {
        setButtonLoading(false);
    }
}

async function streamTaskStatus(taskId) {
    try {
        const res = await fetch(`${API_BASE}/api/tasks/${taskId}/stream`, {
            headers: getHeaders()
        });

        if (!res.ok) {
            pollTaskStatus(taskId);
            return;
        }

        const reader = res.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            buffer += decoder.decode(value, { stream: true });
            const lines = buffer.split('\n');
            buffer = lines.pop();

            for (const line of lines) {
                if (line.startsWith('data: ')) {
                    try {
                        const data = JSON.parse(line.slice(6));
                        setOutput(`Status: ${data.status}\nProgress: ${data.progress}`, false);

                        if (data.status === 'done' || data.status === 'failed' || data.status === 'cancelled') {
                            loadTasks();
                            return;
                        }
                    } catch {}
                }
            }
        }
    } catch {
        pollTaskStatus(taskId);
    }
}

async function pollTaskStatus(taskId) {
    refreshInterval = setInterval(async () => {
        try {
            const res = await fetch(`${API_BASE}/api/tasks/${taskId}`, {
                headers: getHeaders()
            });
            if (!res.ok) return;

            const task = await res.json();
            setOutput(`Status: ${task.status}\nProgress: ${task.progress}${task.result ? '\nResult: ' + task.result : ''}`, false);

            if (task.status === 'done' || task.status === 'failed' || task.status === 'cancelled') {
                clearInterval(refreshInterval);
                loadTasks();
            }
        } catch {}
    }, 1000);
}

async function loadTasks() {
    try {
        const res = await fetch(`${API_BASE}/api/tasks`, {
            headers: getHeaders()
        });
        if (!res.ok) return;

        const tasks = await res.json();
        renderTasks(tasks);
    } catch {}
}

function renderTasks(tasks) {
    if (!elements.tasksList) return;

    elements.tasksList.innerHTML = tasks.map(task => `
        <div class="task-item ${task.status}">
            <div>
                <strong>${task.goal || 'Unnamed task'}</strong>
                <div class="task-id">${task.id}</div>
            </div>
            <div class="task-progress">${task.status} - ${task.progress || ''}</div>
        </div>
    `).join('');
}

async function checkHealth() {
    try {
        const res = await fetch(`${API_BASE}/api/health`);
        if (!res.ok) return;

        const health = await res.json();
        if (elements.healthBar) {
            elements.healthBar.innerHTML = `
                <span>Status: ${health.status}</span>
                <span>Version: ${health.version}</span>
                <span>Uptime: ${health.uptime}</span>
            `;
        }
    } catch {}
}

function refreshHealth() {
    setInterval(checkHealth, 30000);
}

function setOutput(text, isError) {
    if (elements.output) {
        elements.output.textContent = text;
        elements.output.className = isError ? 'status-text error' : 'status-text';
    }
}

function setButtonLoading(loading) {
    if (elements.executeBtn) {
        elements.executeBtn.disabled = loading;
        elements.executeBtn.textContent = loading ? 'Executing...' : 'Execute';
    }
}
