(() => {
  const el = document.getElementById("app");
  if (!el) {
    return;
  }

  const stamp = new Date().toLocaleString();
  const p = document.createElement("p");
  p.textContent = `Local scaffold served at ${stamp}.`;
  el.appendChild(p);
})();

