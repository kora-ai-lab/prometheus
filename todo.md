# Prometheus Master Roadmap & Detailed Todo List

## 🛠️ Phase 0: Proof of Concept (1 week)
- [x] **Core Loop Implementation**
    - [x] Define `Think` -> `Execute` -> `Observe` -> `Decide` loop logic
    - [x] Implement basic `ShellExecutor` for command execution
    - [x] Implement a simple `Observer` to capture stdout/stderr
    - [x] Implement a basic `DecisionEngine` to determine if the task is DONE or needs retry
- [x] **LLM Adapter Layer**
    - [x] Create `ModelProvider` interface (text in -> text out)
    - [x] Implement `OllamaProvider` using Ollama's API
    - [x] Set up basic system prompt for the agent
- [x] **Verification**
    - [x] Test Case: "Create a directory 'poc_test' with a file 'hello.txt' containing 'Hello Prometheus'"
    - [x] Profile binary size to ensure it's < 15MB
    - [x] Verify crash-resilience on shell errors

## ⚡ Phase 1: The Spark (3 weeks)
- [ ] **TUI Development**
    - [ ] Setup `bubbletea` framework
    - [ ] Implement conversation view (User/Agent messages)
    - [ ] Implement basic activity indicators (spinners, progress bars)
- [ ] **Persistence & Memory**
    - [ ] Setup SQLite for `tasks.db`
    - [ ] Implement task state machine (Pending, In Progress, Blocked, Completed)
    - [ ] Implement "Blocked" state handler (Pause loop -> Prompt user -> Resume)
- [ ] **Resilience & Context**
    - [ ] Implement automatic retry logic with exponential backoff
    - [ ] Create basic `ContextManager` (Hot Buffer)
    - [ ] Implement simple compaction (summarize old messages)
- [ ] **Infrastructure**
    - [ ] Implement `EnvironmentDiscovery` (uname, nproc, free, which tools...)
    - [ ] Setup JSON Lines logging to `~/.prometheus/logs/`
    - [ ] Expand LLM adapters: `AnthropicProvider`, `GoogleProvider`

## 🧬 Phase 2: The Evolution (4 weeks)
- [ ] **Capability Engine**
    - [ ] Implement `CapabilityManager` for discovery and metadata
    - [ ] Implement `INSTALL` flow (check availability -> execute install script)
    - [ ] Implement `FORGE` flow (LLM generates script -> tests script -> saves as capability)
- [ ] **Advanced Context Management**
    - [ ] Implement 3-level storage: Hot (RAM) -> Warm (SQLite RAM) -> Cold (SQLite Disk)
    - [ ] Implement adaptive compaction based on model context window size
- [ ] **Web Capabilities**
    - [ ] Integrate Chrome DevTools Protocol (CDP) for headfull/headless browsing
    - [ ] Implement basic scraping and interaction capabilities
- [ ] **Log System Optimization**
    - [ ] Implement `zstd` compression for logs older than 24h
    - [ ] Implement automated daily summary generation via LLM
    - [ ] Implement monthly archiving

## 🛡️ Phase 3: The Immune System (4 weeks)
- [ ] **Security Interceptor**
    - [ ] Build risk-scoring engine for commands (pattern matching + LLM analysis)
    - [ ] Implement rate limiting (max exec/sec)
    - [ ] Implement confirmation dialogs for high-risk commands
- [ ] **Native Security Tools**
    - [ ] Implement Go-native SAST (Scan generated code for SQL injection, hardcoded secrets)
    - [ ] Implement Linux Namespace isolation (Cloneflags: CLONE_NEWNET, CLONE_NEWPID, etc.)
    - [ ] Implement graceful fallback for Windows/Android
- [ ] **Environmental Defense**
    - [ ] Implement background environment scan (open ports, unusual permissions)
    - [ ] Implement native Go DAST for generated web servers
- [ ] **Secrets Management**
    - [ ] Implement AES-256-GCM encrypted Vault for API keys and credentials

## 📈 Phase 4: Observability & Performance (2 weeks)
- [ ] **Structured Observability**
    - [ ] Implement full structured logging with trace IDs
    - [ ] Build internal metrics collector (latency, RAM usage, task success rate)
- [ ] **Semantic Memory**
    - [ ] Implement semantic search over compressed logs using LLM summaries
    - [ ] Develop "Temporal Query" system ("What happened last Monday?")
- [ ] **Optimization**
    - [ ] Implement memory profiling and leak detection
    - [ ] Optimize binary size and startup time

## 🚢 Phase 5: UX & Distribution (3 weeks)
- [ ] **Multi-Surface UI**
    - [ ] Build Vanilla HTML/CSS/JS Web UI with integrated Go server
    - [ ] Create Android WebView wrapper (APK)
- [ ] **Distribution**
    - [ ] Create `curl | sh` installation one-liner
    - [ ] Implement automatic OS/Arch detection for binary delivery
- [ ] **Documentation & Finalization**
    - [ ] Write comprehensive user guides (EN/FR)
    - [ ] Final stress test and benchmark against PRD metrics
