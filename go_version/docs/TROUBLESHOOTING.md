# Troubleshooting

## Common Issues

### "API key not set"

Set your API key:
```bash
export GROQ_API_KEY=sk-...
```

### "Command timed out"

Increase timeout:
```bash
prometheus --timeout 300 "Your goal"
```

### "Permission denied"

Ensure your user has execute permissions:
```bash
chmod +x /usr/local/bin/prometheus
```

### "Module not found"

Update to latest version:
```bash
prometheus --update
```

### Web UI not loading

Check if port 8080 is available:
```bash
lsof -i :8080
```

Use a different port:
```bash
prometheus --web --port 9090
```

## Getting Help

- GitHub Issues: https://github.com/kora-ai-lab/prometheus/issues
- Documentation: https://docs.prometheus.ai