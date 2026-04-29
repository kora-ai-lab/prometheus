# Refactoring Plan: Technical Debt Cleanup (2026-04-29)

## Audit Summary

- **Total Go files audited**: 144
- **Large files (>200 lines)**: 17 files
- **TODO/FIXME comments**: 0 (none found)
- **Hardcoded values**: 6 categories
- **Error handling gaps**: 15+ instances
- **Duplicate patterns**: File-exists checks, initialization sequences

---

## 1. HARDCODED VALUES (High Priority)

### 1.1 Ollama Endpoint Constant
- **Files**: `internal/config/defaults.go:8,47`, `internal/llm/llama_provider.go:222`, `cmd/prometheus/main.go:13`
- **Issue**: `"http://127.0.0.1:11434"` repeated 4 times
- **Fix**: Extract to constant `DefaultOllamaEndpoint` in config package
- **Effort**: 15 min

### 1.2 Web UI Port Constant
- **Files**: `internal/config/defaults.go:37,78`, `cmd/prometheus/main.go:74,75,112`
- **Issue**: Port `8080` hardcoded in multiple places
- **Fix**: Extract to constant `DefaultWebPort` in config package
- **Effort**: 15 min

### 1.3 Context Window Size
- **Files**: `internal/llm/llama_provider.go:45`, `internal/api/manager.go:39`, `internal/capabilities/discovery.go:62`, `internal/llm/modelcatalog.go:31`, `cmd/prometheus/main.go` (test files too)
- **Issue**: `4096` repeated 7+ times
- **Fix**: Extract to constant `DefaultContextWindow` in llm package
- **Effort**: 20 min

### 1.4 HTTP Timeouts
- **Files**: `internal/llm/llama_provider.go:42,85`
- **Issue**: `120 * time.Second`, `30 * time.Second` hardcoded
- **Fix**: Extract to constants `DefaultHTTPTimeout`, `DefaultWaitTimeout`
- **Effort**: 10 min

### 1.5 Prometheus Home Directory
- **Files**: `internal/config/config.go:74,83-86`
- **Issue**: `".prometheus"` and subdirectory names hardcoded
- **Fix**: Extract to constants `PrometheusDirName`, `RuntimeDir`, `ModelsDir`, etc.
- **Effort**: 15 min

---

## 2. ERROR HANDLING GAPS (High Priority)

### 2.1 Ignored Errors in Llama Provider
- **File**: `internal/llm/llama_provider.go:103,207,208,288`
- **Issues**:
  - Line 103: `_ = resp.Body.Close()` - should handle error
  - Line 207: `_ = p.cmd.Process.Kill()` - should log error
  - Line 208: `_, _ = p.cmd.Process.Wait()` - should log error
  - Line 288: `_, _ = io.Copy(io.Discard, resp.Body)` - should handle error
- **Fix**: Add proper error handling or logging
- **Effort**: 30 min

### 2.2 Ignored Time Parse Errors
- **File**: `internal/logging/search.go:61,62`, `internal/task/task_store.go:153,154`
- **Issue**: `time.Parse` errors ignored with `_`
- **Fix**: Log warning or handle error appropriately
- **Effort**: 15 min

### 2.3 Ignored File Stat Errors
- **Files**: `internal/config/config.go:95`, `internal/vault/vault.go:57`, `cmd/prometheus/main.go:56,64,75`
- **Issue**: Pattern `if _, err := os.Stat(...); os.IsNotExist(err)` ignores stat error
- **Fix**: Check both stat error and existence
- **Effort**: 20 min

---

## 3. DUPLICATE PATTERNS (Medium Priority)

### 3.1 Service Initialization Duplication
- **Files**: `cmd/prometheus/main.go:40-95`, `internal/service/service_windows.go:76-130`
- **Issue**: Nearly identical initialization sequence (home, config, logger, storage, executor, etc.)
- **Fix**: Extract to `service.InitService(ctx) (*ServiceDeps, error)` function
- **Effort**: 1 hour

### 3.2 File Exists Check Pattern
- **Files**: 6+ locations using `if _, err := os.Stat(path); os.IsNotExist(err)`
- **Issue**: Repeated pattern for checking file existence
- **Fix**: Create `file.Exists(path) bool` helper
- **Effort**: 15 min

---

## 4. LARGE FILES (Medium Priority - Refactor to smaller units)

| File | Lines | Suggested Refactoring |
|------|-------|---------------------|
| `internal/config/config.go` | 358 | Split into config.go + defaults.go + paths.go |
| `internal/llm/llama_provider.go` | 346 | Extract server management to llama_server.go |
| `internal/logging/search.go` | 331 | Extract BM25/search algorithms to separate file |
| `internal/service/service_windows.go` | 317 | Already has service_unix.go, consider shared init |
| `internal/logging/summarizer.go` | 314 | Extract summary formatting to formatter.go |

---

## 5. UNUSED IMPORTS / CLEANUP (Low Priority)

### 5.1 Stub Client Ignored Parameters
- **File**: `internal/browser/cdp_client.go:9-33`
- **Issue**: Multiple `_ = param` in stub implementations
- **Fix**: Use proper no-op or remove stubs if unused
- **Effort**: 15 min

---

## Implementation Priority Order

1. **HIGH**: Extract hardcoded constants (1.1-1.5) - ~75 min
2. **HIGH**: Fix error handling gaps (2.1-2.3) - ~65 min
3. **MEDIUM**: Extract duplicate patterns (3.1-3.2) - ~75 min
4. **MEDIUM**: Refactor large files (4) - 3-4 hours
5. **LOW**: Cleanup unused code (5) - ~15 min

---

## Verification Commands

After each change:
```bash
go build ./...
go vet ./...
```

---

## Summary

- **Total issues found**: 25+
- **High priority**: 9 issues (constants + error handling)
- **Medium priority**: 6 issues (duplication + large files)
- **Low priority**: 1 issue (cleanup)

Estimated total effort: ~6 hours
