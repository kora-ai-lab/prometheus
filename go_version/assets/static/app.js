const API = {
  async execute(goal) {
    const res = await fetch('/api/execute', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({goal})
    });
    return res.json();
  },

  async metrics() {
    const res = await fetch('/api/metrics');
    return res.json();
  },

  async settings() {
    const res = await fetch('/api/settings');
    return res.json();
  },

  async saveSettings(key, value) {
    const res = await fetch('/api/settings', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({key, value})
    });
    return res.json();
  }
};

document.addEventListener('DOMContentLoaded', () => {
  const navBtns = document.querySelectorAll('.nav-btn');
  const panels = document.querySelectorAll('.panel');

  navBtns.forEach(btn => {
    btn.addEventListener('click', () => {
      navBtns.forEach(b => b.classList.remove('active'));
      panels.forEach(p => p.classList.remove('active'));
      btn.classList.add('active');
      document.getElementById(`${btn.dataset.panel}-panel`).classList.add('active');
    });
  });

  const executeBtn = document.getElementById('execute-btn');
  const taskInput = document.getElementById('task-input');
  const output = document.getElementById('output');

  executeBtn.addEventListener('click', async () => {
    const goal = taskInput.value.trim();
    if (!goal) return;

    output.innerHTML = '<p class="loading">Executing...</p>';
    executeBtn.disabled = true;

    try {
      const result = await API.execute(goal);
      output.innerHTML = `<pre>${result.result || result.error}</pre>`;
    } catch (err) {
      output.innerHTML = `<p class="error">${err.message}</p>`;
    }

    executeBtn.disabled = false;
  });

  setInterval(async () => {
    try {
      const metrics = await API.metrics();
      document.getElementById('ram-usage').textContent =
        metrics.ram_usage ? `${metrics.ram_usage}MB` : '--';
      document.getElementById('task-count').textContent =
        `${metrics.tasks_done || 0}/${metrics.tasks_started || 0}`;
      document.getElementById('llm-calls').textContent = metrics.llm_calls || 0;
    } catch (e) {}
  }, 5000);

  const settingsForm = document.getElementById('settings-form');
  settingsForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(settingsForm);
    await API.saveSettings('model', formData.get('model'));
    alert('Settings saved');
  });
});