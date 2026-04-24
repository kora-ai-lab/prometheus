# 🔥 PROMETHEUS — SCAFFOLD COMPLET FINAL
## Version 5.0 — Tous les Gaps Fermés · Dépendances Exactes · Ordre Correct
### Zéro Dépendance Externe Après Lancement · Africa-First · Production-Ready

---

> **Conventions**
> `[ ]` à faire · `[~]` en cours · `[x]` terminé · `[!]` bloquant critique
> ⚡ = ordre obligatoire, ne pas sauter
> 🔁 = parallélisable avec d'autres tâches de la même phase
> 🧪 = critère de vérification obligatoire avant de continuer
> 👁 = implique la couche vision
> 🌐 = implique le contrôle navigateur
> `→` = dépend de (ne pas commencer sans que le prérequis soit vert)

> **Statut repo au 2026-04-24**
> `go build ./...` passe.
> Fondation P.0-P.2 complète : llama-server Windows AMD64 compilé + embed + build Go ~28MB.

---

## GAPS IDENTIFIÉS ET FERMÉS (par rapport aux versions précédentes)

```
GAP 1 : Embedding llama-server multi-plateforme → binaire 40MB+
         FIX : Build tags Go → chaque binaire n'embarque QUE son llama-server

GAP 2 : Playwright-go comme défaut → requiert Node.js 50MB + download navigateur
         FIX : CDP natif Go UNIQUEMENT par défaut · Playwright = optionnel via
               Capability Engine (installé à la demande, pas dans le binaire)

GAP 3 : Aucune pipeline CI pour compiler llama-server cross-platform
         FIX : GitHub Actions matrix build documentée + artefacts stockés

GAP 4 : Vision model non spécifiée (quel GGUF, quel endpoint)
         FIX : Même llama-server, GGUF vision-capable, endpoint /v1/chat/completions
               avec images en base64

GAP 5 : Dérivation clé vault non spécifiée → crypto correcte
         FIX : crypto/hkdf (Go stdlib, zéro dep) + machine-id stable

GAP 6 : Context window depuis le modèle non spécifiée
         FIX : llama-server expose /v1/models avec context_length

GAP 7 : Téléchargement modèle interruptible non spécifié
         FIX : HTTP Range requests + reprise sur le SHA256 check

GAP 8 : Goroutine leak detection absente
         FIX : go.uber.org/goleak dans TestMain

GAP 9 : Android Termux build tag non précisé
         FIX : GOOS=linux GOARCH=arm64 CGO_ENABLED=0 (Termux = Linux userland)

GAP 10 : Action "create" dans la task loop sans format précis
          FIX : Format JSON structuré avec chemin + contenu séparés

GAP 11 : Sauvegarde SQLite non précisée (quand ?)
          FIX : Sur chaque changement de Status + après chaque exec()

GAP 12 : Logs sécurité dans le log général → pas de séparation
          FIX : Fichier séparé ~/.prometheus/security/events.jsonl

GAP 13 : Download modèle sans barre de progression ni gestion d'interruption
          FIX : io.TeeReader + Range HTTP + SHA256 streaming

GAP 14 : Plusieurs tâches simultanées non spécifiées
          FIX : Une goroutine + ContextManager par tâche, SQLite WAL mode

GAP 15 : Rate limiter LLM absent (seulement exec)
          FIX : golang.org/x/time/rate sur les appels LLM aussi

GAP 16 : playwright-go dans go.mod principal → dépendance lourde au build
          FIX : go.mod secondaire dans internal/capabilities/playwright/

GAP 17 : Pas de spec pour arrêt propre (graceful shutdown)
          FIX : signal.NotifyContext + séquence de cleanup documentée

GAP 18 : Pas de spec pour la mise à jour de Prometheus lui-même
          FIX : prometheus update commande + vérification SHA256 du binaire

GAP 19 : Web UI serveur écoute sur 0.0.0.0 (risque sécurité)
          FIX : Écouter UNIQUEMENT sur 127.0.0.1 par défaut

GAP 20 : Aucune stratégie de test pour la vision (mock LLM visuel)
          FIX : Interface VisionProvider → mock injectée dans les tests
```

---

## TABLE DES MATIÈRES

```
PRÉ-REQUIS OUTILLAGE ............ P.0 — P.3
FONDATION ........................ F.1 LLM Embarqué
                                   F.2 System Prompt
                                   F.3 Vision
                                   F.4 Browser CDP
PHASE 0 — POC .................. T0.1 → T0.7
PHASE 1 — THE SPARK ............ T1.1 → T1.9
PHASE 2 — THE EVOLUTION ........ T2.1 → T2.6
PHASE 3 — IMMUNE SYSTEM ........ T3.1 → T3.5
PHASE 4 — OBSERVABILITY ........ T4.1 → T4.3
PHASE 5 — DISTRIBUTION ......... T5.1 → T5.4
PHASE 6 — ADVANCED ............. T6.1 → T6.5
CROSS-CUTTING ................... TC.1 → TC.4
```

---

## PRÉ-REQUIS OUTILLAGE — AVANT TOUT CODE

---

### P.0 — Machine de build ⚡

- [x] **P.0.1** Installer Go 1.22+ (pas 1.21 — 1.22 apporte des améliorations iter)
  ```bash
  # Linux/macOS
  curl -L https://go.dev/dl/go1.22.5.linux-amd64.tar.gz | tar -C /usr/local -xz
  export PATH=$PATH:/usr/local/go/bin

  # Vérification OBLIGATOIRE
  go version          # → go1.22.x
  go env CGO_ENABLED  # → noter la valeur par défaut
  ```

- [x] **P.0.2** Installer les outils de développement Go
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  go install go.uber.org/goleak/cmd/goleak@latest
  go install golang.org/x/vuln/cmd/govulncheck@latest
  go install github.com/goreleaser/goreleaser/v2@latest

  # Vérifications
  golangci-lint version   # → v1.x
  govulncheck --version   # → v1.x
  goreleaser --version    # → v2.x
  ```

- [x] **P.0.3** Installer cmake + ninja (pour compiler llama-server)
  ```bash
  # Ubuntu/Debian
  apt install cmake ninja-build build-essential -y

  # macOS
  brew install cmake ninja

  # Vérification
  cmake --version    # → 3.20+
  ninja --version    # → 1.11+
  ```

- [ ] **P.0.4** Installer Android NDK (pour cross-compiler llama-server ARM64)
  ```bash
  # Télécharger NDK r27b ou plus récent
  wget https://dl.google.com/android/repository/android-ndk-r27b-linux.zip
  unzip android-ndk-r27b-linux.zip -d $HOME/android-ndk
  export ANDROID_NDK=$HOME/android-ndk/android-ndk-r27b

  # Vérification
  $ANDROID_NDK/ndk-build --version   # → NDK version
  ls $ANDROID_NDK/build/cmake/android.toolchain.cmake  # doit exister
  ```

  🧪 Vérification P.0 :
  - `go version` → `go1.22.x`
  - `cmake --version` → `3.20+`
  - `$ANDROID_NDK/ndk-build --version` → version affichée
  - `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /dev/null ./...` dans un projet test → succès

---

### P.1 — Compiler llama-server pour toutes les plateformes ⚡

**Point clé (Gap 1 fermé)** : llama-server est un binaire C++ autonome.
On le compile une fois par plateforme, on le stocke dans `internal/llm/embedded/`.
Chaque binaire Go n'embarque QUE le llama-server de SA plateforme via build tags.

- [x] **P.1.1** Obtenir llama-server Windows AMD64 (pré-compilé depuis GitHub releases)
  ```bash
  # Télécharger la dernière version depuis llama.cpp releases
  wget https://github.com/ggml-org/llama.cpp/releases/latest/llama-b8914-bin-win-cpu-x64.zip
  unzip llama-b8914-bin-win-cpu-x64.zip
  cp llama-server.exe internal/llm/embedded/llama-server-windows-amd64.exe

  # Vérification
  ./llama-server.exe --version  # → version s'affiche
  ```
  ```bash
  git clone https://github.com/ggml-org/llama.cpp
  cd llama.cpp
  mkdir build-linux-amd64 && cd build-linux-amd64

  cmake .. \
    -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DBUILD_SHARED_LIBS=OFF \
    -DLLAMA_CURL=OFF \
    -DGGML_OPENMP=OFF \
    -DGGML_CUDA=OFF \
    -DGGML_METAL=OFF

  ninja llama-server

  # Vérification critique
  ldd bin/llama-server  # → "not a dynamic executable" OU seulement libm/libc (OK)
  ls -lh bin/llama-server  # → entre 4MB et 10MB
  ./bin/llama-server --version  # → version s'affiche
  ```

- [ ] **P.1.2** Compiler pour Linux ARM64 (cross-compilation depuis AMD64)
  ```bash
  cd llama.cpp
  mkdir build-linux-arm64 && cd build-linux-arm64

  # Installer le cross-compiler
  apt install gcc-aarch64-linux-gnu g++-aarch64-linux-gnu -y

  cmake .. \
    -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_C_COMPILER=aarch64-linux-gnu-gcc \
    -DCMAKE_CXX_COMPILER=aarch64-linux-gnu-g++ \
    -DCMAKE_SYSTEM_NAME=Linux \
    -DCMAKE_SYSTEM_PROCESSOR=aarch64 \
    -DBUILD_SHARED_LIBS=OFF \
    -DLLAMA_CURL=OFF \
    -DGGML_OPENMP=OFF \
    -DGGML_CUDA=OFF

  ninja llama-server

  # Vérification
  file bin/llama-server  # → "ELF 64-bit LSB executable, ARM aarch64"
  ls -lh bin/llama-server  # → entre 4MB et 10MB
  ```

- [ ] **P.1.3** Compiler pour Android ARM64 (ANDROID_PLATFORM=28 minimum)
  ```bash
  cd llama.cpp
  mkdir build-android-arm64 && cd build-android-arm64

  cmake .. \
    -G Ninja \
    -DCMAKE_TOOLCHAIN_FILE=$ANDROID_NDK/build/cmake/android.toolchain.cmake \
    -DANDROID_ABI=arm64-v8a \
    -DANDROID_PLATFORM=android-28 \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_C_FLAGS="-march=armv8-a" \
    -DCMAKE_CXX_FLAGS="-march=armv8-a" \
    -DBUILD_SHARED_LIBS=OFF \
    -DLLAMA_CURL=OFF \
    -DGGML_OPENMP=OFF \
    -DGGML_CUDA=OFF

  ninja llama-server

  # IMPORTANT : Ce binaire android fonctionne dans Termux (Linux userland)
  # car Termux expose un noyau Linux standard sous Android
  file bin/llama-server  # → "ELF 64-bit LSB executable, ARM aarch64"
  # NOTE : même binaire que linux-arm64 si Termux — vérifier au cas par cas
  ```

- [ ] **P.1.4** Compiler pour macOS ARM64 (sur machine Apple Silicon)
  ```bash
  cd llama.cpp
  mkdir build-darwin-arm64 && cd build-darwin-arm64

  cmake .. \
    -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DBUILD_SHARED_LIBS=OFF \
    -DLLAMA_CURL=OFF \
    -DGGML_METAL=OFF \
    -DGGML_OPENMP=OFF

  ninja llama-server

  # Vérification
  file bin/llama-server  # → "Mach-O 64-bit executable arm64"
  ls -lh bin/llama-server  # → entre 4MB et 10MB
  ```

- [ ] **P.1.5** Compiler pour macOS AMD64 (sur Intel ou cross depuis ARM64)
  ```bash
  # Si cross-compilation depuis Apple Silicon vers Intel :
  cmake .. \
    -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_OSX_ARCHITECTURES=x86_64 \
    -DBUILD_SHARED_LIBS=OFF \
    -DLLAMA_CURL=OFF \
    -DGGML_METAL=OFF

  ninja llama-server
  file bin/llama-server  # → "Mach-O 64-bit executable x86_64"
  ```

- [ ] **P.1.6** Compiler pour Windows AMD64 (cross depuis Linux avec MinGW)
  ```bash
  apt install mingw-w64 -y

  cmake .. \
    -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_TOOLCHAIN_FILE=../cmake/x86_64-w64-mingw32.cmake \
    -DBUILD_SHARED_LIBS=OFF \
    -DLLAMA_CURL=OFF

  ninja llama-server

  file bin/llama-server.exe  # → "PE32+ executable (console) x86-64"
  ls -lh bin/llama-server.exe  # → entre 4MB et 12MB
  ```

- [ ] **P.1.7** Copier les binaires compilés dans la structure du projet
  ```
  internal/llm/embedded/
    llama-server-linux-amd64       (compilé P.1.1)
    llama-server-linux-arm64       (compilé P.1.2)
    llama-server-android-arm64     (compilé P.1.3 — si différent de linux-arm64)
    llama-server-darwin-arm64      (compilé P.1.4)
    llama-server-darwin-amd64      (compilé P.1.5)
    llama-server-windows-amd64.exe (compilé P.1.6)
    checksums.sha256               (SHA256 de chaque binaire)
  ```

- [ ] **P.1.8** Générer les checksums SHA256
  ```bash
  cd internal/llm/embedded/
  sha256sum llama-server-* > checksums.sha256
  cat checksums.sha256  # vérifier que chaque ligne est correcte
  ```

  🧪 Vérification P.1 :
  - `file llama-server-linux-arm64` → ARM aarch64
  - `ldd llama-server-linux-amd64` → pas de `.so` dynamiques hors libm/libc
  - `ls -lh llama-server-*` → tous < 12MB
  - `sha256sum -c checksums.sha256` → tous OK

---

### P.2 — Configurer l'embedding par build tags ⚡ (Gap 1 fermé)

- [x] **P.2.1** Créer les fichiers d'embedding avec build tags
  ```go
  // internal/llm/embedded/embed_linux_amd64.go
  //go:build linux && amd64
  package embedded

  import _ "embed"

  //go:embed llama-server-linux-amd64
  var ServerBinary []byte

  const ServerName = "llama-server-linux-amd64"
  ```

  ```go
  // internal/llm/embedded/embed_linux_arm64.go
  //go:build linux && arm64
  package embedded

  import _ "embed"

  //go:embed llama-server-linux-arm64
  var ServerBinary []byte

  const ServerName = "llama-server-linux-arm64"
  ```

  ```go
  // internal/llm/embedded/embed_darwin_arm64.go
  //go:build darwin && arm64
  package embedded
  // ... idem
  ```

  ```go
  // internal/llm/embedded/embed_windows_amd64.go
  //go:build windows && amd64
  package embedded
  // ... idem pour .exe
  ```

  ```go
  // internal/llm/embedded/embed_fallback.go
  //go:build !((linux && amd64) || (linux && arm64) || (darwin && arm64) || (darwin && amd64) || (windows && amd64))
  package embedded

  var ServerBinary []byte  // vide sur plateformes non supportées
  const ServerName = ""
  ```

- [x] **P.2.2** Vérifier la taille des binaires produits
  ```bash
  # Binaire linux-amd64 : llama-server-linux-amd64 embarqué
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" -o bin/prometheus-linux-amd64 ./cmd/prometheus
  ls -lh bin/prometheus-linux-amd64
  # → entre 18MB et 28MB (Go runtime ~12MB + llama-server ~8MB)

  # Binaire linux-arm64 : llama-server-linux-arm64 embarqué
  CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags="-s -w" -o bin/prometheus-linux-arm64 ./cmd/prometheus
  ls -lh bin/prometheus-linux-arm64
  # → entre 16MB et 26MB
  ```

  🧪 Vérification P.2 :
  - `ls -lh bin/prometheus-linux-arm64` → < 30MB
  - `ls -lh bin/prometheus-linux-amd64` → < 30MB
  - Sur machine ARM64, llama-server extrait depuis le binaire → correct binaire ARM64

---

### P.3 — Pipeline CI pour builds automatisés ⚡

- [ ] **P.3.1** Créer `.github/workflows/build-llama-server.yml`
  ```yaml
  # Déclenché manuellement ou sur tag "llama-server-*"
  # Produit les binaires llama-server pour toutes les plateformes
  # Les stocke comme artefacts de release GitHub

  name: Build llama-server
  on:
    workflow_dispatch:
    push:
      tags: ['llama-server-*']

  jobs:
    build-linux:
      runs-on: ubuntu-22.04
      steps:
        - uses: actions/checkout@v4
          with:
            repository: ggml-org/llama.cpp
            ref: master

        - name: Install deps
          run: |
            apt-get update
            apt-get install -y cmake ninja-build \
              gcc-aarch64-linux-gnu g++-aarch64-linux-gnu

        - name: Build linux-amd64
          run: |
            mkdir build-amd64 && cd build-amd64
            cmake .. -G Ninja -DCMAKE_BUILD_TYPE=Release \
              -DBUILD_SHARED_LIBS=OFF -DLLAMA_CURL=OFF -DGGML_OPENMP=OFF
            ninja llama-server
            cp bin/llama-server ../llama-server-linux-amd64

        - name: Build linux-arm64
          run: |
            mkdir build-arm64 && cd build-arm64
            cmake .. -G Ninja -DCMAKE_BUILD_TYPE=Release \
              -DCMAKE_C_COMPILER=aarch64-linux-gnu-gcc \
              -DCMAKE_CXX_COMPILER=aarch64-linux-gnu-g++ \
              -DCMAKE_SYSTEM_NAME=Linux \
              -DCMAKE_SYSTEM_PROCESSOR=aarch64 \
              -DBUILD_SHARED_LIBS=OFF -DLLAMA_CURL=OFF -DGGML_OPENMP=OFF
            ninja llama-server
            cp bin/llama-server ../llama-server-linux-arm64

        - name: Upload artifacts
          uses: actions/upload-artifact@v4
          with:
            name: llama-server-linux
            path: llama-server-linux-*

    build-macos:
      runs-on: macos-14  # Apple Silicon
      steps:
        - uses: actions/checkout@v4
          with:
            repository: ggml-org/llama.cpp
            ref: master

        - name: Build darwin-arm64
          run: |
            mkdir build-arm64 && cd build-arm64
            cmake .. -G Ninja -DCMAKE_BUILD_TYPE=Release \
              -DBUILD_SHARED_LIBS=OFF -DLLAMA_CURL=OFF -DGGML_METAL=OFF
            ninja llama-server
            cp bin/llama-server ../llama-server-darwin-arm64

        - name: Upload artifacts
          uses: actions/upload-artifact@v4
          with:
            name: llama-server-macos
            path: llama-server-darwin-*
  ```

- [ ] **P.3.2** Créer `.github/workflows/ci.yml` (CI principal)
  ```yaml
  name: CI
  on: [push, pull_request]

  jobs:
    test:
      runs-on: ubuntu-22.04
      steps:
        - uses: actions/checkout@v4
        - uses: actions/setup-go@v5
          with:
            go-version: '1.22'

        # Placer des stubs pour les tests (pas les vrais binaires)
        - name: Create test stubs
          run: |
            mkdir -p internal/llm/embedded
            # Stubs vides pour que les build tags compilent
            echo "stub" > internal/llm/embedded/llama-server-linux-amd64
            echo "stub" > internal/llm/embedded/llama-server-linux-arm64
            echo "stub" > internal/llm/embedded/llama-server-darwin-arm64
            echo "stub" > internal/llm/embedded/llama-server-darwin-amd64
            echo "stub" > internal/llm/embedded/llama-server-windows-amd64.exe

        - name: Test
          run: CGO_ENABLED=0 go test -race -count=1 ./...

        - name: Lint
          run: golangci-lint run ./...

        - name: Vuln check
          run: govulncheck ./...

        - name: Build all platforms
          run: |
            CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build \
              -ldflags="-s -w" -o /dev/null ./cmd/prometheus
            CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build \
              -ldflags="-s -w" -o /dev/null ./cmd/prometheus
            CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build \
              -ldflags="-s -w" -o /dev/null ./cmd/prometheus
            CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
              -ldflags="-s -w" -o /dev/null ./cmd/prometheus
  ```

  🧪 Vérification P.3 :
  - GitHub Actions CI passe au vert sur un push
  - Artefacts `llama-server-*` téléchargeables depuis la page du workflow

---

## FONDATION — ARCHITECTURE CORE

---

### F.1 — LLM Embarqué · Architecture Exacte ⚡

**Stack LLM retenue après recherche :**

```
ARCHITECTURE DÉCIDÉE :
  llama-server (C++ statique, compilé P.1) embarqué via build tags (P.2)
  → extrait dans ~/.prometheus/runtime/llama-server au premier lancement
  → lancé comme subprocess sur un port local aléatoire
  → API OpenAI-compatible (/v1/chat/completions)
  → même instance partagée entre tous les workers

POURQUOI PAS yzma/purego :
  yzma (github.com/hybridgroup/yzma) utilise purego+ffi pour charger libllama.so
  à la runtime. Excellent concept mais :
  - Requiert libllama.so séparé (shared lib, ~8MB additionnel)
  - API plus complexe que HTTP (gestion mémoire FFI)
  - Moins de documentation disponible
  → llama-server subprocess est plus simple et plus robuste pour ce cas

POURQUOI PAS go-llama.cpp :
  - Requiert CGO → cross-compilation Android arm64 impossible sans NDK complet
  - go-skynet/go-llama.cpp non maintenu depuis 2023

MODÈLES PAR DÉFAUT (téléchargés une seule fois) :
  RAM < 3 GB → Qwen2.5-0.5B-Q4_K_M.gguf     (390MB)  minimal viable
  RAM 3-6 GB → Phi-3-mini-4k-Q4_K_M.gguf    (2.2GB)  ★ recommandé Africa
  RAM 6-12GB → Llama-3.2-3B-Q4_K_M.gguf     (2.0GB)  équilibré
  RAM > 12GB → Mistral-7B-v0.3-Q4_K_M.gguf  (4.1GB)  haute qualité

MODÈLES VISION (optionnels, téléchargés si besoin) :
  RAM < 4 GB → moondream2-Q4_K_M.gguf        (900MB)  basique
  RAM 4-8 GB → Phi-3-vision-Q4_K_M.gguf     (2.5GB)  ★ recommandé
  RAM > 8 GB → LLaVA-1.6-mistral-Q4_K_M.gguf(5.0GB)  excellent

SOURCE MODÈLES : Hugging Face (HTTPS direct, pas d'API key)
VÉRIFICATION   : SHA256 après chaque téléchargement
REPRISE        : HTTP Range requests si téléchargement interrompu
```

- [x] **F.1.1** ⚡ → P.2 Définir `internal/llm/embedded/extractor.go`
  ```go
  // Extrait le llama-server embarqué vers ~/.prometheus/runtime/
  // Vérifie le SHA256 avant d'écraser une version existante
  // chmod +x sur Linux/macOS
  // Ne ré-extrait pas si checksum identique (démarrage rapide)

  func ExtractServer() (path string, err error) {
      if len(embedded.ServerBinary) == 0 {
          return "", ErrPlatformNotSupported{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH}
      }
      dest := filepath.Join(prometheusHome(), "runtime", "llama-server")
      if runtime.GOOS == "windows" {
          dest += ".exe"
      }

      // Vérifier si déjà extrait avec le bon checksum
      if existingOK(dest, sha256OfEmbedded()) {
          return dest, nil
      }

      // Écrire
      if err := os.MkdirAll(filepath.Dir(dest), 0700); err != nil {
          return "", err
      }
      if err := os.WriteFile(dest, embedded.ServerBinary, 0700); err != nil {
          return "", err
      }
      return dest, nil
  }

  func sha256OfEmbedded() string {
      h := sha256.New()
      h.Write(embedded.ServerBinary)
      return hex.EncodeToString(h.Sum(nil))
  }

  func existingOK(path, expectedSHA256 string) bool {
      data, err := os.ReadFile(path)
      if err != nil { return false }
      h := sha256.New()
      h.Write(data)
      return hex.EncodeToString(h.Sum(nil)) == expectedSHA256
  }
  ```

- [x] **F.1.2** ⚡ Définir `internal/llm/modelcatalog.go`
  ```go
  type ModelEntry struct {
      Name          string
      Filename      string
      SizeBytes     int64
      MinRAMMb      int
      ContextWindow int
      URL           string    // Hugging Face direct HTTPS
      SHA256        string    // checksum du fichier .gguf
      Quality       string    // "minimal"|"recommended"|"balanced"|"high"
      IsVision      bool
  }

  var TextModels = []ModelEntry{
      {
          Name: "Qwen2.5 0.5B (minimal)",
          Filename: "qwen2.5-0.5b-q4_k_m.gguf",
          SizeBytes: 390_000_000,
          MinRAMMb: 1500, ContextWindow: 32768,
          URL: "https://huggingface.co/Qwen/Qwen2.5-0.5B-GGUF/resolve/main/qwen2.5-0.5b-q4_k_m.gguf",
          SHA256: "...", Quality: "minimal",
      },
      {
          Name: "Phi-3 Mini 4K (recommandé Africa)",
          Filename: "phi-3-mini-4k-instruct-q4_k_m.gguf",
          SizeBytes: 2_200_000_000,
          MinRAMMb: 3500, ContextWindow: 4096,
          URL: "https://huggingface.co/microsoft/Phi-3-mini-4k-instruct-gguf/resolve/main/Phi-3-mini-4k-instruct-q4.gguf",
          SHA256: "...", Quality: "recommended",
      },
      // ... Llama-3.2-3B, Mistral-7B
  }

  var VisionModels = []ModelEntry{
      {
          Name: "Moondream2 (minimal vision)",
          Filename: "moondream2-q4_k_m.gguf",
          SizeBytes: 900_000_000,
          MinRAMMb: 2000, IsVision: true,
          URL: "https://huggingface.co/vikhyatk/moondream2-gguf/resolve/main/...",
          SHA256: "...",
      },
      // ... Phi-3-vision, LLaVA-1.6
  }

  func SelectModel(ramMb int) *ModelEntry {
      switch {
      case ramMb < 3000: return &TextModels[0]  // Qwen minimal
      case ramMb < 6000: return &TextModels[1]  // Phi-3 mini ★
      case ramMb < 12000: return &TextModels[2] // Llama-3.2-3B
      default: return &TextModels[3]             // Mistral-7B
      }
  }
  ```

- [x] **F.1.3** ⚡ Définir `internal/llm/downloader.go` (Gap 7 fermé)
  ```go
  // Téléchargement avec :
  // 1. Barre de progression (io.TeeReader + compteur d'octets)
  // 2. Reprise (HTTP Range header si fichier partiel existe)
  // 3. Vérification SHA256 en streaming
  // 4. Fichier temporaire → rename atomique

  func Download(ctx context.Context, entry *ModelEntry, dest string,
                progress func(downloaded, total int64)) error {

      tmpPath := dest + ".tmp"
      var startByte int64

      // Vérifier si téléchargement partiel existe
      if fi, err := os.Stat(tmpPath); err == nil {
          startByte = fi.Size()
      }

      req, _ := http.NewRequestWithContext(ctx, "GET", entry.URL, nil)
      if startByte > 0 {
          req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startByte))
      }

      resp, err := http.DefaultClient.Do(req)
      if err != nil { return err }
      defer resp.Body.Close()

      f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
      if err != nil { return err }
      defer f.Close()

      h := sha256.New()
      // Si reprise, hasher les données déjà téléchargées
      if startByte > 0 {
          existing, _ := os.ReadFile(tmpPath)
          h.Write(existing[:startByte])
      }

      counter := &ProgressWriter{
          Writer:     io.MultiWriter(f, h),
          OnProgress: progress,
          Total:      entry.SizeBytes,
          Downloaded: startByte,
      }

      if _, err := io.Copy(counter, resp.Body); err != nil {
          return err // garder le .tmp pour reprise
      }

      // Vérifier SHA256
      got := hex.EncodeToString(h.Sum(nil))
      if got != entry.SHA256 {
          os.Remove(tmpPath)
          return fmt.Errorf("SHA256 mismatch: got %s, want %s", got, entry.SHA256)
      }

      // Rename atomique
      return os.Rename(tmpPath, dest)
  }
  ```

- [x] **F.1.4** ⚡ Définir `internal/llm/llama_provider.go`
  ```go
  type LocalLlamaProvider struct {
      serverPath   string
      modelPath    string
      visionPath   string       // optionnel
      textPort     int          // port aléatoire libre
      visionPort   int          // port aléatoire libre (si vision)
      textProc     *os.Process
      visionProc   *os.Process
      httpClient   *http.Client
      modelInfo    *ModelInfo
      mu           sync.Mutex
  }

  func NewLocalLlamaProvider(serverPath, modelPath, visionPath string) (*LocalLlamaProvider, error) {
      p := &LocalLlamaProvider{
          serverPath: serverPath,
          modelPath:  modelPath,
          visionPath: visionPath,
      }
      if err := p.startTextServer(); err != nil {
          return nil, err
      }
      if visionPath != "" {
          if err := p.startVisionServer(); err != nil {
              // Vision en échec → warning, pas d'erreur fatale
              log.Printf("WARN: vision model failed to start: %v", err)
          }
      }
      return p, nil
  }

  func (p *LocalLlamaProvider) startTextServer() error {
      port, err := freePort()
      if err != nil { return err }
      p.textPort = port

      // Calculer les paramètres selon la RAM disponible
      ramMb := availableRAMMb()
      threads := max(1, runtime.NumCPU()/2)
      ctxSize := p.modelInfo.ContextWindow  // lu depuis les métadonnées

      args := []string{
          "--model", p.modelPath,
          "--port", strconv.Itoa(port),
          "--host", "127.0.0.1",
          "--threads", strconv.Itoa(threads),
          "--ctx-size", strconv.Itoa(ctxSize),
          "--parallel", parallelCount(ramMb),
          "--log-disable",
          "--no-webui",
      }
      cmd := exec.Command(p.serverPath, args...)
      cmd.Stdout = io.Discard
      cmd.Stderr = io.Discard

      if err := cmd.Start(); err != nil { return err }
      p.textProc = cmd.Process

      // Attendre que le serveur soit prêt (healthcheck)
      return p.waitReady(port, 30*time.Second)
  }

  func (p *LocalLlamaProvider) waitReady(port int, timeout time.Duration) error {
      deadline := time.Now().Add(timeout)
      url := fmt.Sprintf("http://127.0.0.1:%d/health", port)
      for time.Now().Before(deadline) {
          resp, err := http.Get(url)
          if err == nil && resp.StatusCode == 200 {
              resp.Body.Close()
              return nil
          }
          time.Sleep(500 * time.Millisecond)
      }
      return fmt.Errorf("llama-server not ready after %v", timeout)
  }

  // ModelInfo() : récupère le context window depuis /v1/models (Gap 6 fermé)
  func (p *LocalLlamaProvider) ModelInfo() *ModelInfo {
      resp, _ := http.Get(fmt.Sprintf("http://127.0.0.1:%d/v1/models", p.textPort))
      // Parser la réponse pour extraire context_length
      // ...
      return p.modelInfo
  }

  // Close() : arrêt propre des processus (Gap 17 fermé partiellement)
  func (p *LocalLlamaProvider) Close() error {
      if p.textProc != nil { p.textProc.Signal(syscall.SIGTERM) }
      if p.visionProc != nil { p.visionProc.Signal(syscall.SIGTERM) }
      // Attendre max 5s, puis SIGKILL
      return nil
  }

  // HasVision() : true si visionPort > 0
  func (p *LocalLlamaProvider) HasVision() bool { return p.visionPort > 0 }
  ```

- [x] **F.1.5** ⚡ Interfaces `ModelProvider` et `VisionProvider`
  ```go
  // internal/llm/provider.go
  type Message struct {
      Role    string
      Content string
      Images  [][]byte  // PNG/JPEG base64 pour les messages vision
  }

  type Response struct {
      Content      string
      InputTokens  int
      OutputTokens int
  }

  type ModelInfo struct {
      Name          string
      ContextWindow int
      Provider      string
      HasVision     bool
  }

  type ModelProvider interface {
      Complete(ctx context.Context, messages []Message) (*Response, error)
      Stream(ctx context.Context, messages []Message, tokens chan<- string) error
      ModelInfo() *ModelInfo
      IsAvailable() bool
      HasVision() bool
      Close() error
  }

  // Providers cloud (Anthropic, Google) implémentent aussi HasVision() = true
  ```

- [x] **F.1.6** Autres providers : `OllamaProvider`, `AnthropicProvider`, `GoogleProvider`
  ```go
  // Tous implémentent ModelProvider
  // Tous pur Go (CGO_ENABLED=0), HTTP client standard

  // OllamaProvider : context window depuis GET /api/show {"name": model}
  //   → .modelinfo.context_length
  // AnthropicProvider : context window hardcodé par modèle
  //   haiku-4-5 → 200000, sonnet-4-6 → 200000
  // GoogleProvider : context window hardcodé
  //   gemini-2.0-flash → 1000000

  // Rate limiter LLM (Gap 15 fermé)
  // golang.org/x/time/rate sur les appels Complete()
  // Anthropic : défaut 50 req/min
  // Google : défaut 60 req/min
  // Local : pas de rate limit
  ```

- [x] **F.1.7** `ProviderFactory` — auto-détection hiérarchique
  ```go
  // internal/llm/factory.go
  func AutoDetect(cfg *config.LLMConfig) (ModelProvider, error) {
      // 1. Modèle local (toujours prioritaire)
      if modelExists(cfg.ModelPath) {
          serverPath, err := embedded.ExtractServer()
          if err != nil { goto tryOllama }
          return NewLocalLlamaProvider(serverPath, cfg.ModelPath, cfg.VisionModelPath)
      }

  tryOllama:
      // 2. Ollama si disponible
      if isOllamaRunning() {
          return NewOllamaProvider(cfg)
      }

      // 3. Anthropic
      if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
          return NewAnthropicProvider(cfg, key)
      }

      // 4. Google
      if key := os.Getenv("GOOGLE_API_KEY"); key != "" {
          return NewGoogleProvider(cfg, key)
      }

      // 5. Aucun LLM → guide setup
      return nil, ErrNoLLMAvailable{
          Hint: "Lancez 'prometheus setup' pour télécharger un modèle local",
      }
  }
  ```

- [x] **F.1.8** Premier lancement : `internal/llm/setup.go`
  ```go
  // Appelé si aucun modèle local trouvé
  func FirstRunSetup(env *discovery.EnvironmentProfile, ui UIWriter) error {
      model := SelectModel(env.RAMMb)

      ui.Print(fmt.Sprintf(`
  ╔══════════════════════════════════════════════════╗
  ║  PROMETHEUS — CONFIGURATION INITIALE            ║
  ╠══════════════════════════════════════════════════╣
  ║  Modèle recommandé : %s
  ║  Taille            : %s
  ║  Source            : Hugging Face (HTTPS)
  ╚══════════════════════════════════════════════════╝
  Continuer le téléchargement ? [O/n] `, model.Name, humanSize(model.SizeBytes)))

      if !readYesNo() { return ErrSetupCancelled }

      dest := filepath.Join(prometheusHome(), "models", model.Filename)
      return Download(context.Background(), model, dest, progressBar(ui))
  }
  ```

  🧪 Vérification F.1 :
  - Machine vierge : `./prometheus` → guide setup → modèle téléchargé → LLM répond
  - Téléchargement interrompu (Ctrl+C) → relancer → reprise depuis l'offset
  - SHA256 incorrect → rejet + message d'erreur clair
  - `./prometheus --provider anthropic` → utilise Anthropic sans télécharger
  - `LocalLlamaProvider.Close()` → process llama-server terminé (pas de zombie)

---

### F.2 — System Prompt Complet · L'ADN de Prometheus ⚡

- [x] **F.2.1** Créer `assets/prompts/system_v1.md` (embarqué via `//go:embed`)
  ```markdown
  ---
  prometheus_prompt_version: 1
  min_prometheus_version: 0.1.0
  ---

  Tu es Prometheus, un agent autonome sans limite de création.

  TU PEUX FAIRE :
  - Exécuter n'importe quelle commande shell (exec)
  - Installer n'importe quel outil existant (auto-install)
  - CRÉER ce qui n'existe pas : scripts, applications, APIs, protocoles,
    services, bases de données, interfaces — tout ce que peut exécuter
    un shell, tu peux le construire de zéro (create)
  - Contrôler un navigateur web : naviguer, cliquer, remplir, scraper (browser)
  - Voir les résultats visuels via screenshots (vision)
  - Te souvenir de toutes les sessions passées (memory)
  - Te donner de nouveaux pouvoirs en les installant ou en les créant (forge)

  RÈGLES ABSOLUES (jamais enfreintes) :
  1. Ne jamais t'arrêter sans raison valide.
     Bloqué sur une info → action=ask. Erreur → corriger, réessayer.
     Outil manquant → l'installer ou le créer, puis continuer.
  2. Toujours répondre en JSON strict. Zéro texte libre en dehors du JSON.
  3. Résultat visuel incertain → prendre un screenshot (action=vision).
  4. Commande dangereuse → mettre dangerous=true et expliquer.
  5. Jamais abandonner avant d'avoir atteint l'objectif ou épuisé toutes
     les alternatives possibles.
  ```

- [x] **F.2.2** ⚡ Format de réponse JSON (dans `assets/prompts/system_v1.md`)
  ```json
  FORMAT DE RÉPONSE OBLIGATOIRE (toujours ce JSON, jamais autre chose) :
  {
    "thinking": "raisonnement en 1-2 phrases maximum",
    "action": "exec|ask|browser|vision|create|done|error",
    "command": "commande shell complète (si action=exec)",
    "create_file": {
      "path": "chemin/vers/fichier.py",
      "content": "contenu complet du fichier"
    },
    "browser_action": "navigate|click|fill|submit|screenshot|get_html|eval_js|scroll|wait_for|get_cookies",
    "browser_args": {"url": "...", "selector": "...", "text": "...", "script": "..."},
    "vision_target": "browser|screen|file",
    "vision_file":   "chemin/vers/fichier.png (si vision_target=file)",
    "question": "question précise à l'utilisateur (si action=ask)",
    "dangerous": false,
    "why": "justification en une phrase"
  }

  RÈGLE create_file : créer le fichier directement via cet objet JSON,
  pas via exec("cat > fichier << 'EOF'..."). Cela évite les problèmes
  d'échappement et de caractères spéciaux.
  ```

- [x] **F.2.3** ⚡ `internal/prompt/builder.go` — assemblage adaptatif
  ```go
  // 5 blocs, budget adaptatif selon context window

  type PromptBuilder struct {
      systemPromptBase string     // Bloc A+E embarqué depuis assets/
      contextWindow    int        // tokens du modèle actif
      historyRatio     float64    // 0.60 = 60% pour l'historique
      env              *discovery.EnvironmentProfile
      capabilities     []string
      patterns         []string
  }

  func (pb *PromptBuilder) Build() string {
      // Budget pour le system prompt
      budget := int(float64(pb.contextWindow) * (1.0 - pb.historyRatio))

      prompt := pb.systemPromptBase   // Blocs A+E, jamais omis (~250 tokens)
      budget -= 250
      if budget <= 0 { return prompt }

      // Bloc B : Environnement (50 tokens)
      if budget > 50 {
          prompt += "\n\n" + pb.buildEnvBlock()
          budget -= 50
      }

      // Bloc C : Capacités (80 tokens max)
      if budget > 80 && len(pb.capabilities) > 0 {
          prompt += "\n" + pb.buildCapsBlock()
          budget -= 80
      }

      // Bloc D : Patterns appris (80 tokens max)
      if budget > 80 && len(pb.patterns) > 0 {
          prompt += "\n" + pb.buildPatternsBlock()
      }

      return prompt
  }

  func (pb *PromptBuilder) buildEnvBlock() string {
      return fmt.Sprintf("ENV:%s/%s RAM:%dMB CPU:%d Net:%s Pkg:%s Tools:%s",
          pb.env.OS, pb.env.Arch, pb.env.RAMMb, pb.env.CPUCores,
          boolStr(pb.env.Internet, "on", "off"),
          pb.env.PackageManager,
          strings.Join(pb.env.AvailableTools[:min(8, len(pb.env.AvailableTools))], ","),
      )
  }
  ```

- [x] **F.2.4** ⚡ `internal/prompt/parser.go` — parser robuste (Gap 10 fermé)
  ```go
  type Action struct {
      Thinking      string            `json:"thinking"`
      Action        string            `json:"action"`
      Command       string            `json:"command"`
      CreateFile    *CreateFileAction `json:"create_file,omitempty"`
      BrowserAction string            `json:"browser_action,omitempty"`
      BrowserArgs   map[string]string `json:"browser_args,omitempty"`
      VisionTarget  string            `json:"vision_target,omitempty"`
      VisionFile    string            `json:"vision_file,omitempty"`
      Question      string            `json:"question,omitempty"`
      Dangerous     bool              `json:"dangerous"`
      Why           string            `json:"why"`
  }

  type CreateFileAction struct {
      Path    string `json:"path"`
      Content string `json:"content"`
  }

  func ParseAction(raw string) (*Action, error) {
      // Extraire le premier bloc JSON valide
      // (le modèle peut ajouter du texte avant/après)
      jsonStr := extractFirstJSON(raw)
      if jsonStr == "" {
          return nil, &ParseError{Raw: raw, Msg: "aucun JSON trouvé"}
      }

      var action Action
      if err := json.Unmarshal([]byte(jsonStr), &action); err != nil {
          return nil, &ParseError{Raw: raw, ParseErr: err}
      }

      // Validation
      if action.Action == "" {
          return nil, &ParseError{Msg: "champ 'action' manquant"}
      }
      switch action.Action {
      case "exec":
          if action.Command == "" { return nil, &ParseError{Msg: "command vide"} }
      case "create":
          if action.CreateFile == nil || action.CreateFile.Path == "" {
              return nil, &ParseError{Msg: "create_file manquant"}
          }
      case "browser":
          if action.BrowserAction == "" {
              return nil, &ParseError{Msg: "browser_action manquant"}
          }
      case "ask":
          if action.Question == "" { return nil, &ParseError{Msg: "question vide"} }
      case "vision", "done", "error":
          // OK
      default:
          return nil, &ParseError{Msg: "action inconnue: " + action.Action}
      }

      return &action, nil
  }

  // extractFirstJSON : cherche le premier {...} complet dans le texte
  func extractFirstJSON(s string) string {
      start := strings.Index(s, "{")
      if start == -1 { return "" }
      depth := 0
      for i := start; i < len(s); i++ {
          switch s[i] {
          case '{': depth++
          case '}':
              depth--
              if depth == 0 { return s[start : i+1] }
          }
      }
      return ""
  }
  ```

- [x] **F.2.5** Versioning et rechargement du prompt utilisateur
  ```go
  // Charger dans l'ordre :
  // 1. ~/.prometheus/prompts/system_v1.md (si existe et version compatible)
  // 2. Sinon : version embarquée dans le binaire (assets/prompts/system_v1.md)

  func LoadSystemPrompt() string {
      userPrompt := filepath.Join(prometheusHome(), "prompts", "system_v1.md")
      if data, err := os.ReadFile(userPrompt); err == nil {
          if isCompatible(data) {
              return string(data)
          }
          log.Printf("WARN: prompt utilisateur incompatible, utilisation du défaut")
      }
      return defaultSystemPrompt  // //go:embed assets/prompts/system_v1.md
  }
  ```

  🧪 Vérification F.2 :
  - `EstimateTokens(Build())` < 500 sur Phi-3 mini (4096 tokens)
  - `EstimateTokens(Build())` < 800 sur Claude Sonnet (200K)
  - `ParseAction(`{"action":"exec","command":"ls"}`)` → Action{Action:"exec", Command:"ls"}
  - `ParseAction("Voici ma réponse : {\"action\":\"done\"}")` → Action{Action:"done"}
  - JSON invalide → ParseError → retry avec message de correction

---

### F.3 — Couche Vision · Prometheus Voit 👁 ⚡

- [x] **F.3.1** ⚡ Interface `VisionProvider` (séparée de ModelProvider)
  ```go
  // internal/vision/provider.go
  type VisionProvider interface {
      Analyze(ctx context.Context, imageBytes []byte, question string) (string, error)
      HasVision() bool
  }

  // Implémentations :
  // LocalVisionProvider  : llama-server avec vision model, port séparé
  // CloudVisionProvider  : Anthropic/Google (même provider que le texte)
  // NoOpVisionProvider   : mode dégradé (pas de crash, log "vision unavailable")
  ```

- [~] **F.3.2** `LocalVisionProvider` → llama-server avec modèle vision
  ```go
  // Même binaire llama-server, différent modèle GGUF
  // Port différent du provider texte
  // Endpoint : POST /v1/chat/completions
  // Message avec image :
  // {"role":"user","content":[
  //   {"type":"image_url","image_url":{"url":"data:image/png;base64,[BASE64]"}},
  //   {"type":"text","text":"[question]"}
  // ]}
  ```

- [~] **F.3.3** `ScreenCapture` — capture cross-platform sans dépendance externe
  ```go
  // internal/vision/capture.go
  func CaptureScreen(ctx context.Context, exec executor.Executor) ([]byte, error) {
      tmpFile := filepath.Join(os.TempDir(), "prometheus-screen.png")

      var cmd string
      switch runtime.GOOS {
      case "linux":
          // Essayer dans l'ordre jusqu'à succès
          for _, tool := range []string{"scrot", "import", "gnome-screenshot"} {
              if toolExists(tool) {
                  switch tool {
                  case "scrot": cmd = "scrot " + tmpFile
                  case "import": cmd = "import -window root " + tmpFile
                  case "gnome-screenshot": cmd = "gnome-screenshot -f " + tmpFile
                  }
                  break
              }
          }
          if cmd == "" {
              // Aucun outil → demander au Capability Engine d'installer scrot
              return nil, ErrScreenshotToolMissing{"linux", "scrot"}
          }
      case "darwin":
          cmd = "screencapture -x " + tmpFile  // natif macOS, zéro install
      case "windows":
          cmd = `powershell -command "Add-Type -AssemblyName System.Windows.Forms;` +
              `[System.Windows.Forms.Screen]::PrimaryScreen | ForEach-Object {` +
              `$bitmap = New-Object Drawing.Bitmap($_.Bounds.Width,$_.Bounds.Height);` +
              `$graphics = [Drawing.Graphics]::FromImage($bitmap);` +
              `$graphics.CopyFromScreen($_.Bounds.Location,[Drawing.Point]::Empty,$_.Bounds.Size);` +
              `$bitmap.Save('` + tmpFile + `')}"`
      }

      // Android/Termux : `screencap -p` (natif Android, zéro install)
      if isTermux() {
          cmd = "screencap -p " + tmpFile
      }

      result := exec.Execute(ctx, cmd, ExecOptions{Timeout: 10 * time.Second})
      if result.ExitCode != 0 {
          return nil, fmt.Errorf("screenshot failed: %s", result.Stderr)
      }

      data, err := os.ReadFile(tmpFile)
      os.Remove(tmpFile)  // nettoyer
      return data, err
  }
  ```

- [x] **F.3.4** Mode dégradé si vision indisponible
  ```go
  // NoOpVisionProvider : implémente VisionProvider
  func (n *NoOpVisionProvider) Analyze(_ context.Context, _ []byte, _ string) (string, error) {
      return "[vision non disponible — continuer sans analyse visuelle]", nil
  }
  func (n *NoOpVisionProvider) HasVision() bool { return false }
  ```

  🧪 Vérification F.3 :
  - `CaptureScreen()` sur macOS → PNG valide (screencapture natif)
  - `CaptureScreen()` sur Android Termux → PNG valide (screencap natif)
  - `LocalVisionProvider.Analyze(png, "Que vois-tu ?")` → description cohérente
  - `NoOpVisionProvider.Analyze(...)` → string "non disponible", pas d'erreur
  - Task ne crashe pas si vision non disponible (mode dégradé transparent)

---

### F.4 — Contrôle Navigateur · CDP Natif · Playwright Optionnel 🌐 ⚡

**Décision finale (Gap 2 fermé) :**

```
NIVEAU 1 — CDP natif Go (DANS LE BINAIRE, zéro dépendance externe)
  → nhooyr.io/websocket (pure Go, CGO_ENABLED=0)
  → Chromium/Chrome/Firefox détecté sur la machine
  → Si absent : Capability Engine installe Chromium
  → Couvre 90% des cas : navigate, get_html, click, fill, screenshot

NIVEAU 2 — Playwright-go (OPTIONNEL, installé via Capability Engine)
  playwright-go = pure Go mais télécharge un driver Node.js (~50MB)
  → Installé seulement si l'utilisateur en a besoin
  → go get github.com/playwright-community/playwright-go dans un module séparé
  → Couvre les 10% complexes : SPA, infinite scroll, fichiers, auth multi-step

PLAYWRIGHT N'EST PAS dans le go.mod principal.
Il est installé à la demande via Capability Engine si CDP est insuffisant.
```

- [x] **F.4.1** ⚡ Interface `BrowserClient`
  ```go
  // internal/browser/client.go
  type BrowserClient interface {
      Navigate(url string) error
      GetHTML() (string, error)
      GetText() (string, error)
      Click(selector string) error
      Fill(selector, value string) error
      Submit(selector string) error
      Screenshot() ([]byte, error)
      WaitForSelector(selector string, timeout time.Duration) error
      WaitForNavigation(timeout time.Duration) error
      EvalJS(script string) (interface{}, error)
      GetCookies() ([]*Cookie, error)
      SetCookie(c *Cookie) error
      ScrollDown(pixels int) error
      Close() error
  }
  ```

- [~] **F.4.2** ⚡ CDP Client — implémentation pure Go
  ```go
  // internal/browser/cdp_client.go
  // Dépendance : nhooyr.io/websocket (pure Go, CGO_ENABLED=0)
  // Protocole : WebSocket vers chrome --remote-debugging-port=PORT

  type CDPClient struct {
      port      int
      process   *os.Process
      conn      *websocket.Conn
      sessionID string
      mu        sync.Mutex
      nextID    int64
  }

  func LaunchCDP(browserPath string) (*CDPClient, error) {
      port, err := freePort()
      if err != nil { return nil, err }

      userDataDir, _ := os.MkdirTemp("", "prometheus-browser-")

      args := []string{
          "--headless=new",
          "--remote-debugging-port=" + strconv.Itoa(port),
          "--no-sandbox",
          "--disable-gpu",
          "--disable-dev-shm-usage",
          "--disable-extensions",
          "--user-data-dir=" + userDataDir,
          "about:blank",
      }

      cmd := exec.Command(browserPath, args...)
      cmd.Stdout = io.Discard
      cmd.Stderr = io.Discard
      if err := cmd.Start(); err != nil { return nil, err }

      c := &CDPClient{port: port, process: cmd.Process}
      if err := c.connect(); err != nil { return nil, err }
      return c, nil
  }

  // FindBrowser : cherche chromium/chrome sur la machine
  func FindBrowser() (string, error) {
      candidates := map[string][]string{
          "linux":   {"chromium", "chromium-browser", "google-chrome", "chrome"},
          "darwin":  {
              "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
              "/Applications/Chromium.app/Contents/MacOS/Chromium",
          },
          "windows": {
              os.Getenv("LOCALAPPDATA") + `\Google\Chrome\Application\chrome.exe`,
          },
      }
      for _, name := range candidates[runtime.GOOS] {
          if path, err := exec.LookPath(name); err == nil {
              return path, nil
          }
          // Pour les chemins absolus
          if _, err := os.Stat(name); err == nil {
              return name, nil
          }
      }
      return "", ErrBrowserNotFound{}
  }

  // Commandes CDP clés implémentées :
  // Page.navigate         → Navigate()
  // DOM.getDocument + DOM.getOuterHTML → GetHTML()
  // Input.dispatchMouseEvent → Click()
  // Input.insertText      → Fill()
  // Page.captureScreenshot → Screenshot() (retourne PNG bytes)
  // Runtime.evaluate      → EvalJS()
  // Network.getCookies    → GetCookies()
  // Page.handleJavaScriptDialog → accepter les popups auto
  ```

- [~] **F.4.3** `BrowserManager` — sélection automatique CDP / Playwright
  ```go
  // internal/browser/manager.go
  type BrowserManager struct {
      cdpClient        *CDPClient
      playwrightClient BrowserClient  // nil si non installé
      capEngine        *capabilities.Engine
      vision           vision.VisionProvider
  }

  func (bm *BrowserManager) GetClient(needsAdvanced bool) BrowserClient {
      if needsAdvanced && bm.playwrightClient != nil {
          return bm.playwrightClient
      }
      return bm.cdpClient  // toujours disponible
  }

  // Si Playwright nécessaire et non installé :
  func (bm *BrowserManager) EnsurePlaywright() error {
      return bm.capEngine.Ensure("playwright")
      // Capability Engine télécharge Node.js + playwright driver
      // Puis charge playwright-go dynamiquement
  }
  ```

- [~] **F.4.4** Vision intégrée aux actions browser (auto-screenshot)
  ```go
  // Après navigate, click, submit si VisionProvider disponible :
  func (bm *BrowserManager) Do(ctx context.Context, action *prompt.Action) string {
      client := bm.GetClient(action.BrowserNeedsAdvanced())
      var err error

      switch action.BrowserAction {
      case "navigate":
          err = client.Navigate(action.BrowserArgs["url"])
      case "click":
          err = client.Click(action.BrowserArgs["selector"])
      case "fill":
          err = client.Fill(action.BrowserArgs["selector"], action.BrowserArgs["text"])
      case "screenshot":
          img, e := client.Screenshot()
          if e != nil { return "ERROR: " + e.Error() }
          if bm.vision.HasVision() {
              analysis, _ := bm.vision.Analyze(ctx, img, "Décris ce que tu vois")
              return "[SCREENSHOT ANALYSÉ]\n" + analysis
          }
          // Sauvegarder le PNG
          path := saveScreenshot(img)
          return "Screenshot sauvegardé: " + path
      case "get_html":
          html, e := client.GetHTML()
          if e != nil { return "ERROR: " + e.Error() }
          return truncateMid(html, 50000)  // max 50KB
      case "eval_js":
          result, e := client.EvalJS(action.BrowserArgs["script"])
          if e != nil { return "ERROR: " + e.Error() }
          return fmt.Sprintf("%v", result)
      case "wait_for":
          timeout := parseDuration(action.BrowserArgs["timeout"], 10*time.Second)
          err = client.WaitForSelector(action.BrowserArgs["selector"], timeout)
      case "scroll":
          pixels, _ := strconv.Atoi(action.BrowserArgs["pixels"])
          err = client.ScrollDown(pixels)
      case "get_cookies":
          cookies, e := client.GetCookies()
          if e != nil { return "ERROR: " + e.Error() }
          data, _ := json.Marshal(cookies)
          return string(data)
      }

      if err != nil { return "ERROR: " + err.Error() }

      // Auto-screenshot après action significative (si vision dispo)
      if bm.vision.HasVision() && action.BrowserAction != "get_html" {
          img, e := client.Screenshot()
          if e == nil {
              analysis, _ := bm.vision.Analyze(ctx, img, "Résultat de l'action?")
              return "ACTION OK\n[VISION]\n" + analysis
          }
      }
      return "ACTION OK"
  }
  ```

  🧪 Vérification F.4 :
  - CDP : `Navigate("https://example.com")` + `GetHTML()` → HTML reçu avec le titre
  - CDP : `Screenshot()` → bytes PNG valide (testable avec `png.DecodeConfig`)
  - CDP : `Click("#submit")` → retour "ACTION OK" ou "ERROR: element not found"
  - Browser absent → `FindBrowser()` retourne `ErrBrowserNotFound` → Capability Engine déclenché
  - Browser fermé proprement : après `Close()`, pas de processus zombie

---

## PHASE 0 — PROOF OF CONCEPT (1 semaine)
**→ Prérequis : P.0, P.1, P.2, P.3, F.1, F.2 doivent être VERTS**
**Objectif : Boucle Think→Execute→Observe avec LLM embarqué.**

---

### T0.1 — Initialisation projet ⚡ → P.0

- [x] **T0.1.1** Structure de dossiers complète
  ```
  prometheus/
    cmd/prometheus/main.go
    internal/
      llm/
        embedded/
          embed_linux_amd64.go     ← build tags (P.2)
          embed_linux_arm64.go
          embed_darwin_arm64.go
          embed_darwin_amd64.go
          embed_windows_amd64.go
          embed_fallback.go
          extractor.go             ← F.1.1
          llama-server-*           ← binaires compilés (P.1)
          checksums.sha256
        llama_provider.go          ← F.1.4
        ollama_provider.go
        anthropic_provider.go
        google_provider.go
        factory.go                 ← F.1.7
        modelcatalog.go            ← F.1.2
        downloader.go              ← F.1.3
        setup.go                   ← F.1.8
        provider.go                ← F.1.5 interfaces
      executor/
        executor.go
        rate_limiter.go
      task/
        task.go
        loop.go
        deps.go                    ← TaskDeps struct
      prompt/
        builder.go                 ← F.2.3
        parser.go                  ← F.2.4
        estimator.go               ← EstimateTokens()
      config/
        config.go
        defaults.go
      discovery/
        discovery.go
      storage/
        migrations.go
        task_store.go
        memory_store.go
      logging/
        logger.go
        redact.go                  ← RedactSecrets()
      context/
        manager.go
      vision/
        provider.go                ← F.3.1
        capture.go                 ← F.3.3
        local_provider.go          ← F.3.2
        noop_provider.go           ← F.3.4
      browser/
        client.go                  ← F.4.1
        cdp_client.go              ← F.4.2
        manager.go                 ← F.4.3
      capabilities/
        engine.go
        registry.go
        forge.go
        store.go
      security/
        interceptor.go
        patterns.go
        sast.go
        sandbox.go
        env_scanner.go
      vault/
        vault.go
      ui/
        tui.go
        blocked.go
        webui.go
      metrics/
        collector.go
    assets/
      prompts/
        system_v1.md               ← F.2.1
      static/
        index.html
        style.css
        app.js
    go.mod
    go.sum
    Makefile
    .gitignore
    .golangci.yml
    goreleaser.yml
  ```

- [x] **T0.1.2** `go.mod` — dépendances Phase 0
  ```go
  module github.com/prometheus-dev/prometheus

  go 1.22

  require (
      // Config
      github.com/BurntSushi/toml v1.3.2

      // SQLite pure Go (CGO_ENABLED=0 compatible)
      modernc.org/sqlite v1.31.1

      // Rate limiter
      golang.org/x/time v0.5.0

      // Saisie sécurisée (passwords)
      golang.org/x/term v0.20.0

      // WebSocket pur Go (CDP browser)
      nhooyr.io/websocket v1.8.11

      // Compression zstd (logs)
      github.com/klauspost/compress v1.17.8

      // TUI (Phase 1)
      github.com/charmbracelet/bubbletea v0.27.0
      github.com/charmbracelet/lipgloss v1.0.0

      // Tests
      go.uber.org/goleak v1.3.0
  )

  // NOTE : playwright-go INTENTIONNELLEMENT ABSENT du go.mod principal
  // Il est installé via Capability Engine si nécessaire
  // NOTE : CGO_ENABLED=0 requis pour toutes les builds
  // Aucune dépendance CGO dans ce go.mod
  ```

- [x] **T0.1.3** `Makefile` complet
  ```makefile
  BINARY    := prometheus
  LDFLAGS   := -ldflags="-s -w"
  CGO       := CGO_ENABLED=0

  .PHONY: build build-all test lint size-check embed-servers clean

  build:
      $(CGO) go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/prometheus

  build-android:
      $(CGO) GOOS=linux GOARCH=arm64 go build $(LDFLAGS) \
          -o bin/$(BINARY)-linux-arm64 ./cmd/prometheus

  build-all:
      $(CGO) GOOS=linux   GOARCH=amd64  go build $(LDFLAGS) -o bin/$(BINARY)-linux-amd64 ./cmd/prometheus
      $(CGO) GOOS=linux   GOARCH=arm64  go build $(LDFLAGS) -o bin/$(BINARY)-linux-arm64 ./cmd/prometheus
      $(CGO) GOOS=darwin  GOARCH=arm64  go build $(LDFLAGS) -o bin/$(BINARY)-darwin-arm64 ./cmd/prometheus
      $(CGO) GOOS=darwin  GOARCH=amd64  go build $(LDFLAGS) -o bin/$(BINARY)-darwin-amd64 ./cmd/prometheus
      $(CGO) GOOS=windows GOARCH=amd64  go build $(LDFLAGS) -o bin/$(BINARY)-windows.exe ./cmd/prometheus

  test:
      $(CGO) go test -race -count=1 -timeout=120s ./...

  test-cover:
      $(CGO) go test -coverprofile=coverage.out -race ./...
      go tool cover -func=coverage.out | grep -E "total|<threshold"

  lint:
      golangci-lint run --timeout=5m ./...

  vuln:
      govulncheck ./...

  size-check:
      @$(MAKE) build-android
      @SIZE=$$(wc -c < bin/$(BINARY)-linux-arm64); \
       MAX=31457280; \
       if [ $$SIZE -gt $$MAX ]; then \
           echo "FAIL: binaire linux-arm64 trop lourd: $$SIZE > $$MAX bytes"; exit 1; \
       else echo "OK: linux-arm64 = $$SIZE bytes ($$(echo "$$SIZE / 1048576" | bc)MB)"; fi

  embed-servers:
      @echo "Copie des llama-server compilés dans internal/llm/embedded/"
      @for f in llama-server-linux-amd64 llama-server-linux-arm64 \
                llama-server-darwin-arm64 llama-server-darwin-amd64 \
                llama-server-windows-amd64.exe; do \
          if [ -f build/$$f ]; then cp build/$$f internal/llm/embedded/; \
          else echo "WARN: build/$$f absent"; fi; \
      done
      cd internal/llm/embedded && sha256sum llama-server-* > checksums.sha256

  clean:
      rm -rf bin/ coverage.out
  ```

- [x] **T0.1.4** `.gitignore`
  ```
  /bin/
  *.db *.db-shm *.db-wal
  *.enc
  *.gguf
  coverage.out
  /internal/llm/embedded/llama-server-*
  /internal/llm/embedded/checksums.sha256
  ```

  🧪 Vérification T0.1 :
  - `CGO_ENABLED=0 go build ./...` (avec stubs) → succès
  - `go vet ./...` → zéro erreur
  - `make size-check` → passe (nécessite les vrais binaires llama-server)

---

### T0.2 — Configuration ⚡ → T0.1

- [x] **T0.2.1** `internal/config/config.go`
  ```go
  type Config struct {
      LLM      LLMConfig      `toml:"llm"`
      Vision   VisionConfig   `toml:"vision"`
      Browser  BrowserConfig  `toml:"browser"`
      Security SecurityConfig `toml:"security"`
      Memory   MemoryConfig   `toml:"memory"`
      Logs     LogConfig      `toml:"logs"`
      UI       UIConfig       `toml:"ui"`
  }

  type LLMConfig struct {
      Provider        string `toml:"provider"`          // "local"|"ollama"|"anthropic"|"google"
      ModelPath       string `toml:"model_path"`        // chemin .gguf
      VisionModelPath string `toml:"vision_model_path"` // chemin .gguf vision
      ModelName       string `toml:"model_name"`        // pour ollama
      Endpoint        string `toml:"endpoint"`          // pour ollama/custom
      // API keys : JAMAIS dans TOML → ANTHROPIC_API_KEY, GOOGLE_API_KEY env vars
  }

  type SecurityConfig struct {
      RateLimitExecPerSec int      `toml:"rate_limit_per_second"`      // 10
      RateLimitLLMPerMin  int      `toml:"rate_limit_llm_per_min"`     // 60 (Gap 15)
      DangerousOpsConfirm bool     `toml:"dangerous_ops_confirmation"` // true
      SandboxEnabled      bool     `toml:"sandbox"`                    // false (auto-détecté)
  }

  type UIConfig struct {
      WebEnabled bool   `toml:"web_enabled"`  // false par défaut
      WebPort    int    `toml:"web_port"`      // 8080
      WebHost    string `toml:"web_host"`      // "127.0.0.1" (Gap 19)
  }
  ```

- [x] **T0.2.2** `prometheus.toml` par défaut (généré au premier lancement)
  ```toml
  [llm]
  provider = "local"
  # model_path sera défini par 'prometheus setup'

  [vision]
  enabled = true
  auto_capture = true

  [browser]
  enabled = true
  level = "cdp"   # "cdp" | "playwright" | "auto"
  headless = true
  timeout = 30

  [security]
  rate_limit_per_second = 10
  rate_limit_llm_per_min = 60
  dangerous_ops_confirmation = true
  sandbox = false

  [memory]
  compaction_threshold = 0.70

  [logs]
  compress_after_days = 1
  archive_after_days = 7

  [ui]
  web_enabled = false
  web_port = 8080
  web_host = "127.0.0.1"
  ```

  🧪 Vérification T0.2 :
  - Premier lancement → `~/.prometheus/prometheus.toml` créé
  - `PROMETHEUS_LLM_PROVIDER=anthropic ./prometheus` → override effectif
  - Config invalide → message d'erreur précis avec nom du champ

---

### T0.3 — Executor ⚡ → T0.1

- [x] **T0.3.1** `internal/executor/executor.go`
  ```go
  type ExecOptions struct {
      Timeout time.Duration  // défaut: 5min
      WorkDir string
      Env     []string       // variables d'env supplémentaires
  }

  type ExecResult struct {
      Command  string
      Stdout   string
      Stderr   string
      ExitCode int
      Duration time.Duration
      TimedOut bool
  }

  func Execute(ctx context.Context, command string, opts ExecOptions) *ExecResult {
      if opts.Timeout == 0 { opts.Timeout = 5 * time.Minute }

      ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
      defer cancel()

      var cmd *exec.Cmd
      switch runtime.GOOS {
      case "windows":
          cmd = exec.CommandContext(ctx, "cmd", "/C", command)
      default:
          cmd = exec.CommandContext(ctx, "sh", "-c", command)
      }

      if opts.WorkDir != "" { cmd.Dir = opts.WorkDir }
      cmd.Env = append(os.Environ(), opts.Env...)

      var stdout, stderr bytes.Buffer
      cmd.Stdout = &stdout
      cmd.Stderr = &stderr

      start := time.Now()
      err := cmd.Run()
      duration := time.Since(start)

      result := &ExecResult{
          Command:  command,
          Stdout:   truncateMid(stdout.String(), 50_000),
          Stderr:   truncateMid(stderr.String(), 20_000),
          Duration: duration,
          TimedOut: ctx.Err() == context.DeadlineExceeded,
      }

      if err != nil {
          if exitErr, ok := err.(*exec.ExitError); ok {
              result.ExitCode = exitErr.ExitCode()
          } else {
              result.ExitCode = -1
          }
      }
      return result
  }

  // truncateMid : garder les 40% du début + "...[N lignes tronquées]..." + 40% de la fin
  func truncateMid(s string, maxChars int) string {
      if len(s) <= maxChars { return s }
      half := maxChars * 2 / 5
      skipped := len(s) - 2*half
      lines := strings.Count(s[half:len(s)-half], "\n")
      return s[:half] + fmt.Sprintf("\n...[%d lignes tronquées]...\n", lines) + s[len(s)-half:]
  }
  ```

- [x] **T0.3.2** Rate limiter exec (golang.org/x/time/rate)
  ```go
  // internal/executor/rate_limiter.go
  type RateLimitedExecutor struct {
      limiter *rate.Limiter
  }

  func NewRateLimitedExecutor(maxPerSec int) *RateLimitedExecutor {
      return &RateLimitedExecutor{
          limiter: rate.NewLimiter(rate.Limit(maxPerSec), maxPerSec),
      }
  }

  func (r *RateLimitedExecutor) Execute(ctx context.Context, command string, opts ExecOptions) *ExecResult {
      // Attendre (pas rejeter) si rate limit atteint
      if err := r.limiter.Wait(ctx); err != nil {
          return &ExecResult{Command: command, ExitCode: -1,
              Stderr: "rate limit: " + err.Error()}
      }
      return Execute(ctx, command, opts)
  }
  ```

  🧪 Vérification T0.3 :
  - `Execute(ctx, "echo hello", {})` → Stdout="hello\n", ExitCode=0
  - `Execute(ctx, "sleep 10", {Timeout: 1s})` → TimedOut=true en ≤ 1.1s
  - Output 1MB → Stdout tronqué à 50KB avec indicateur "[N lignes tronquées]"
  - 20 exec en rafale → étalés sur ≥ 2s (rate 10/s)

---

### T0.4 — Environment Discovery ⚡ → T0.1

- [x] **T0.4.1** `internal/discovery/discovery.go`
  ```go
  type EnvironmentProfile struct {
      OS             string    // "linux" | "darwin" | "windows"
      Arch           string    // "amd64" | "arm64"
      Kernel         string
      RAMMb          int
      DiskGb         int
      CPUCores       int
      AvailableTools []string  // outils présents dans PATH
      LLMModels      []string  // si Ollama présent
      Internet       bool
      PackageManager string    // "apt" | "brew" | "pkg" | "dnf" | "winget"
      IsTermux       bool
      ScannedAt      time.Time
  }

  func Scan(ctx context.Context, exec executor.Executor) *EnvironmentProfile {
      p := &EnvironmentProfile{
          OS:        runtime.GOOS,
          Arch:      runtime.GOARCH,
          CPUCores:  runtime.NumCPU(),
          IsTermux:  isTermux(),
          ScannedAt: time.Now(),
      }
      p.RAMMb    = readRAM(ctx, exec)
      p.DiskGb   = readDisk(ctx, exec)
      p.Kernel   = readKernel(ctx, exec)
      p.AvailableTools = checkTools(ctx, exec, []string{
          "git", "python3", "python", "node", "npm", "docker",
          "curl", "wget", "chromium", "chromium-browser",
          "google-chrome", "firefox", "adb", "ffmpeg",
          "scrot", "import", "screencap",  // vision tools
      })
      p.PackageManager = detectPackageManager(ctx, exec)
      p.LLMModels      = detectOllamaModels(ctx, exec)
      p.Internet       = checkInternet(ctx, exec)
      return p
  }

  func isTermux() bool {
      prefix := os.Getenv("PREFIX")
      return strings.Contains(prefix, "com.termux")
  }
  ```

  🧪 Vérification T0.4 :
  - Android Termux → `IsTermux=true`, `PackageManager="pkg"`
  - Réseau coupé → `Internet=false` en < 5s
  - `AvailableTools` contient uniquement les outils réellement présents

---

### T0.5 — Task Loop ⚡ → T0.1, F.2

- [x] **T0.5.1** `internal/task/task.go` — types
  ```go
  type TaskStatus string
  const (
      StatusRunning   TaskStatus = "running"
      StatusBlocked   TaskStatus = "blocked"
      StatusDone      TaskStatus = "done"
      StatusFailed    TaskStatus = "failed"
      StatusCancelled TaskStatus = "cancelled"
  )

  type Task struct {
      ID             string
      Goal           string
      Status         TaskStatus
      Context        []llm.Message
      Memory         map[string]interface{}
      BlockedReason  string
      Retries        int
      MaxRetries     int    // défaut: 5
      ParseErrors    int    // retries parse JSON
      MaxParseErrors int    // défaut: 3
      CreatedAt      time.Time
      UpdatedAt      time.Time
  }
  ```

- [x] **T0.5.2** `internal/task/deps.go` — dépendances de la tâche
  ```go
  // TaskDeps centralise toutes les dépendances injectées dans la boucle
  // Facilite les tests (mock)
  type TaskDeps struct {
      Provider      llm.ModelProvider
      Executor      executor.Executor
      Vision        vision.VisionProvider
      Browser       *browser.BrowserManager
      PromptBuilder *prompt.Builder
      CapEngine     *capabilities.Engine
      Security      *security.Interceptor
      Logger        *logging.Logger
      TaskStore     storage.TaskStore
  }
  ```

- [x] **T0.5.3** ⚡ `internal/task/loop.go` — boucle principale
  ```go
  func (t *Task) Run(ctx context.Context, deps *TaskDeps) error {
      deps.Logger.LogTaskStart(t.ID, t.Goal)
      defer func() { deps.Logger.LogTaskEnd(t.ID, t.Status) }()

      for t.Status == StatusRunning {
          // 1. Construire les messages avec context manager
          messages := deps.PromptBuilder.BuildMessages(t.Context)

          // 2. Appeler le LLM
          start := time.Now()
          resp, err := deps.Provider.Complete(ctx, messages)
          deps.Logger.LogLLMCall(t.ID, resp, time.Since(start))
          if err != nil {
              t.handleLLMError(err)
              continue
          }

          // 3. Ajouter à l'historique
          t.Context = append(t.Context, llm.Message{Role:"assistant", Content:resp.Content})

          // 4. Parser l'action
          action, parseErr := prompt.ParseAction(resp.Content)
          if parseErr != nil {
              t.ParseErrors++
              if t.ParseErrors >= t.MaxParseErrors {
                  t.Status = StatusFailed
                  return nil
              }
              // Retry avec message de correction
              t.Context = append(t.Context, llm.Message{
                  Role: "user",
                  Content: fmt.Sprintf(
                      "Erreur: ta réponse n'était pas du JSON valide.\n"+
                      "Réponds UNIQUEMENT avec un objet JSON, rien d'autre.\n"+
                      "Erreur: %v\n"+
                      "Ta réponse (extrait): %s",
                      parseErr, truncateMid(resp.Content, 200)),
              })
              continue
          }
          t.ParseErrors = 0  // reset sur succès

          // 5. Dispatch
          var observation string
          switch action.Action {

          case "exec":
              // Vérification sécurité
              allowed, secErr := deps.Security.Allow(action.Command)
              if !allowed {
                  observation = "BLOQUÉ PAR SÉCURITÉ: " + secErr.Error()
                  t.Context = append(t.Context, llm.Message{Role:"user", Content:observation})
                  continue
              }
              result := deps.Executor.Execute(ctx, action.Command, executor.ExecOptions{})
              deps.Logger.LogExec(t.ID, result)
              observation = formatObservation(result)
              // Si "command not found" → Capability Engine
              if isCommandNotFound(result) {
                  tool := extractToolName(action.Command)
                  if ensureErr := deps.CapEngine.Ensure(ctx, tool); ensureErr == nil {
                      // Réessayer la même commande
                      result = deps.Executor.Execute(ctx, action.Command, executor.ExecOptions{})
                      observation = "[OUTIL INSTALLÉ]\n" + formatObservation(result)
                  }
              }
              t.Retries = 0

          case "create":
              // Créer un fichier directement depuis le JSON (Gap 10 fermé)
              cf := action.CreateFile
              if err := os.MkdirAll(filepath.Dir(cf.Path), 0755); err != nil {
                  observation = "ERROR mkdir: " + err.Error()
              } else if err := os.WriteFile(cf.Path, []byte(cf.Content), 0644); err != nil {
                  observation = "ERROR write: " + err.Error()
              } else {
                  observation = fmt.Sprintf("FILE_CREATED: %s (%d bytes)", cf.Path, len(cf.Content))
                  deps.Logger.LogFileCreated(t.ID, cf.Path)
              }

          case "browser":
              observation = deps.Browser.Do(ctx, action)
              deps.Logger.LogBrowserAction(t.ID, action.BrowserAction)

          case "vision":
              var img []byte
              var captureErr error
              if action.VisionTarget == "browser" {
                  img, captureErr = deps.Browser.Screenshot(ctx)
              } else if action.VisionTarget == "file" {
                  img, captureErr = os.ReadFile(action.VisionFile)
              } else {
                  img, captureErr = vision.CaptureScreen(ctx, deps.Executor)
              }
              if captureErr != nil {
                  observation = "VISION ERROR: " + captureErr.Error()
              } else {
                  analysis, _ := deps.Vision.Analyze(ctx, img, action.Why)
                  observation = "[VISION]\n" + analysis
              }
              deps.Logger.LogVisionCapture(t.ID, action.VisionTarget)

          case "ask":
              t.Status = StatusBlocked
              t.BlockedReason = action.Question
              deps.TaskStore.Save(t)  // sauvegarder l'état bloqué
              return nil

          case "done":
              t.Status = StatusDone
              deps.TaskStore.Save(t)
              return nil

          case "error":
              t.Retries++
              if t.Retries >= t.MaxRetries {
                  t.Status = StatusFailed
                  deps.TaskStore.Save(t)
                  return nil
              }
              observation = "Erreur notée. Réessaie avec une approche différente."
          }

          t.Context = append(t.Context, llm.Message{Role:"user", Content:observation})
          t.UpdatedAt = time.Now()
          deps.TaskStore.Save(t)  // sauvegarder après chaque étape (Gap 11)
      }
      return nil
  }
  ```

- [x] **T0.5.4** `task.Resume(answer string)` — reprise après blocage
  ```go
  func (t *Task) Resume(answer string) {
      t.Context = append(t.Context, llm.Message{
          Role:    "user",
          Content: "Réponse à ta question: " + answer,
      })
      t.Status = StatusRunning
      t.BlockedReason = ""
  }
  ```

  🧪 Vérification T0.5 :
  - `exec` → observation dans le contexte → loop continue
  - `create` → fichier créé sur disque → taille correcte
  - `ask` → StatusBlocked → `Resume("réponse")` → StatusRunning → StatusDone
  - 5 `error` → StatusFailed
  - 3 ParseErrors → StatusFailed
  - SQLite sauvé après chaque action (vérifier avec `sqlite3 tasks.db`)

---

### T0.6 — `main.go` Phase 0 ⚡ → T0.1 à T0.5, F.1

- [x] **T0.6.1** Séquence de démarrage complète
  ```go
  // cmd/prometheus/main.go
  func main() {
      // 0. Graceful shutdown (Gap 17 fermé)
      ctx, stop := signal.NotifyContext(context.Background(),
          os.Interrupt, syscall.SIGTERM)
      defer stop()

      // 1. Créer ~/.prometheus/ si absent
      prometheusHome := setupHome()

      // 2. Charger config
      cfg, err := config.Load(prometheusHome)
      exitOnError(err, "config")

      // 3. Logger (minimal pour Phase 0)
      logger := logging.New(prometheusHome)
      defer logger.Close()

      // 4. SQLite
      store, err := storage.Open(prometheusHome)
      exitOnError(err, "storage")
      defer store.Close()

      // 5. Environment Discovery
      execr := executor.NewRateLimitedExecutor(cfg.Security.RateLimitExecPerSec)
      env := discovery.Scan(ctx, execr)

      // 6. Extraire llama-server
      serverPath, err := embedded.ExtractServer()
      if err != nil {
          log.Printf("WARN: llama-server non disponible: %v", err)
      }

      // 7. Premier lancement ? → setup modèle
      if cfg.LLM.ModelPath == "" || !fileExists(cfg.LLM.ModelPath) {
          if err := llm.FirstRunSetup(env, os.Stdout); err != nil {
              exitOnError(err, "setup modèle")
          }
          // Recharger config (model_path mis à jour)
          cfg, _ = config.Load(prometheusHome)
      }

      // 8. Provider LLM
      provider, err := llm.AutoDetect(cfg.LLM, serverPath)
      exitOnError(err, "LLM")
      defer provider.Close()

      // 9. Vision (optionnel)
      var visionProvider vision.VisionProvider = &vision.NoOpVisionProvider{}
      if cfg.Vision.Enabled {
          visionProvider = llm.NewVisionProvider(provider, serverPath, cfg.LLM.VisionModelPath)
      }

      // 10. Security Interceptor
      secInterceptor := security.New(cfg.Security)

      // 11. Capability Engine
      capEngine := capabilities.NewEngine(execr, env, logger)

      // 12. Prompt Builder
      promptBuilder := prompt.NewBuilder(provider.ModelInfo(), env, capEngine)

      // 13. Browser Manager (CDP)
      browserMgr := browser.NewManager(capEngine, visionProvider)

      // 14. TaskDeps
      deps := &task.TaskDeps{
          Provider:      provider,
          Executor:      execr,
          Vision:        visionProvider,
          Browser:       browserMgr,
          PromptBuilder: promptBuilder,
          CapEngine:     capEngine,
          Security:      secInterceptor,
          Logger:        logger,
          TaskStore:     store,
      }

      // 15. Phase 0 : lire objectif depuis args/stdin, créer task, run
      goal := readGoal(os.Args[1:])
      t := task.New(goal)
      store.Save(t)

      for t.Status == task.StatusRunning || t.Status == task.StatusBlocked {
          if t.Status == task.StatusBlocked {
              fmt.Printf("\n⊙ PROMETHEUS A BESOIN D'UNE INFO:\n%s\n> ", t.BlockedReason)
              var answer string
              fmt.Scanln(&answer)
              t.Resume(answer)
          }
          t.Run(ctx, deps)
      }

      switch t.Status {
      case task.StatusDone:
          fmt.Println("\n✓ Terminé")
      case task.StatusFailed:
          fmt.Println("\n✗ Échec — voir les logs pour plus de détails")
          os.Exit(1)
      }
  }
  ```

  🧪 CRITÈRES DE SUCCÈS PHASE 0 :
  - [ ] Machine vierge + `./prometheus "Crée un fichier test.txt avec Bonjour"` → fichier créé
  - [ ] `./prometheus "Crée une API Python Flask simple"` → app.py créé et fonctionnel
  - [ ] `./prometheus "Écris le jeu Snake en Python"` → snake.py créé et exécutable
  - [ ] Réponse LLM non-JSON → retry → succès (pas de crash)
  - [ ] Binaire linux-arm64 : `ls -lh` → < 30MB
  - [ ] `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ./...` → succès
  - [ ] `go test ./...` → 100% verts
  - [ ] Graceful shutdown : Ctrl+C → task sauvée → relancer → reprend

---

## PHASE 1 — THE SPARK (3 semaines)
**→ Prérequis : Phase 0 VERTE**
**Objectif : MVP complet. TUI, persistence, vault, context manager.**

---

### T1.1 — Context Manager ⚡ → F.2.3

- [~] **T1.1.1** `internal/context/manager.go`
  ```go
  type Manager struct {
      hotBuffer     []llm.Message
      warmSummary   string         // JSON compacté des échanges précédents
      contextWindow int
      threshold     float64
      keepLast      int
      provider      llm.ModelProvider
  }

  func New(provider llm.ModelProvider) *Manager {
      info := provider.ModelInfo()
      m := &Manager{contextWindow: info.ContextWindow, provider: provider}
      m.configure()
      return m
  }

  func (m *Manager) configure() {
      switch {
      case m.contextWindow < 4000:  m.threshold=0.60; m.keepLast=5
      case m.contextWindow < 8000:  m.threshold=0.60; m.keepLast=5   // Phi-3 mini
      case m.contextWindow < 32000: m.threshold=0.65; m.keepLast=10
      case m.contextWindow < 128000:m.threshold=0.70; m.keepLast=20
      default:                       m.threshold=0.80; m.keepLast=50
      }
  }

  func (m *Manager) Add(msg llm.Message) {
      m.hotBuffer = append(m.hotBuffer, msg)
      if m.usageRatio() > m.threshold {
          m.compact()
      }
  }

  func (m *Manager) usageRatio() float64 {
      total := estimateTokens(m.warmSummary)
      for _, msg := range m.hotBuffer {
          total += estimateTokens(msg.Content)
      }
      return float64(total) / float64(m.contextWindow)
  }

  func (m *Manager) compact() {
      if len(m.hotBuffer) <= m.keepLast { return }

      toCompact := m.hotBuffer[:len(m.hotBuffer)-m.keepLast]

      // Construire le texte à compacter
      var sb strings.Builder
      for _, msg := range toCompact {
          sb.WriteString(msg.Role + ": " + msg.Content + "\n")
      }

      compactionPrompt := []llm.Message{
          {Role: "user", Content: fmt.Sprintf(
              `Résume en JSON compact (max 250 tokens) :
  {"goal":"objectif","done":["accompli"],"decisions":["décision"],
   "errors_fixed":["erreur→solution"],"state":"état actuel","next":["suite"]}
  JSON UNIQUEMENT. Conversation :
  %s`, sb.String())},
      }

      ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
      defer cancel()

      resp, err := m.provider.Complete(ctx, compactionPrompt)
      if err != nil { return }  // en cas d'échec, garder le buffer tel quel

      m.warmSummary = resp.Content
      m.hotBuffer = m.hotBuffer[len(m.hotBuffer)-m.keepLast:]
  }

  func (m *Manager) BuildMessages(systemPrompt string) []llm.Message {
      msgs := []llm.Message{{Role: "system", Content: systemPrompt}}
      if m.warmSummary != "" {
          msgs = append(msgs, llm.Message{
              Role:    "system",
              Content: "[CONTEXTE PRÉCÉDENT]\n" + m.warmSummary,
          })
      }
      return append(msgs, m.hotBuffer...)
  }
  ```

- [x] **T1.1.2** `estimateTokens()` — estimation rapide sans tokenizer
  ```go
  func estimateTokens(s string) int {
      if s == "" { return 0 }
      // Règle empirique : 1 token ≈ 4 chars (Latin) ≈ 2 chars (non-Latin)
      // Marge de sécurité × 1.15
      chars := len([]rune(s))
      return int(float64(chars) / 3.5 * 1.15)
  }
  ```

  🧪 Vérification T1.1 :
  - 200 messages sur Phi-3 mini (4096 tokens) → zéro overflow de context
  - Après compaction, `usageRatio()` < 0.50
  - Le résumé contient l'objectif original après 5 compactions

---

### T1.2 — Persistence SQLite ⚡ → T0.1

- [x] **T1.2.1** `modernc.org/sqlite` — CGO_ENABLED=0 compatible (CONFIRMÉ)

- [~] **T1.2.2** Migrations versionnées
  ```go
  // internal/storage/migrations.go
  var migrations = []Migration{
      {Version: 1, SQL: `
          CREATE TABLE IF NOT EXISTS tasks (
              id TEXT PRIMARY KEY,
              goal TEXT NOT NULL,
              status TEXT NOT NULL DEFAULT 'running',
              context_json TEXT,
              memory_json TEXT,
              blocked_reason TEXT,
              created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
              updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
          );
          CREATE TABLE IF NOT EXISTS executions (
              id INTEGER PRIMARY KEY AUTOINCREMENT,
              task_id TEXT NOT NULL,
              command TEXT NOT NULL,
              stdout TEXT, stderr TEXT,
              exit_code INTEGER, duration_ms INTEGER,
              executed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
              FOREIGN KEY(task_id) REFERENCES tasks(id)
          );
          CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY);
          CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
          CREATE INDEX IF NOT EXISTS idx_exec_task ON executions(task_id);
      `},
      {Version: 2, SQL: `
          CREATE TABLE IF NOT EXISTS user_prefs (
              key TEXT PRIMARY KEY, value TEXT, updated_at DATETIME);
          CREATE TABLE IF NOT EXISTS learned_patterns (
              id INTEGER PRIMARY KEY AUTOINCREMENT,
              context TEXT, pattern TEXT,
              success_rate REAL DEFAULT 1.0, uses INTEGER DEFAULT 1,
              last_used DATETIME);
          CREATE TABLE IF NOT EXISTS capabilities_cache (
              name TEXT PRIMARY KEY,
              installed BOOLEAN DEFAULT 0,
              version TEXT, path TEXT,
              metadata_json TEXT,
              installed_at DATETIME);
          CREATE TABLE IF NOT EXISTS sessions (
              id TEXT PRIMARY KEY, date DATE,
              summary TEXT, projects TEXT, technologies TEXT,
              created_at DATETIME);
          CREATE VIRTUAL TABLE IF NOT EXISTS session_search
          USING fts5(date, summary, projects, technologies);
      `},
  }
  ```

- [x] **T1.2.3** WAL mode pour concurrence (Gap 14)
  ```go
  // internal/storage/store.go
  func Open(dir string) (*Store, error) {
      db, err := sql.Open("sqlite", filepath.Join(dir, "prometheus.db"))
      if err != nil { return nil, err }

      // WAL mode pour écriture concurrent multi-goroutines
      db.Exec("PRAGMA journal_mode=WAL")
      db.Exec("PRAGMA synchronous=NORMAL")
      db.Exec("PRAGMA foreign_keys=ON")

      if err := runMigrations(db); err != nil { return nil, err }
      return &Store{db: db}, nil
  }
  ```

- [~] **T1.2.4** Interfaces TaskStore et MemoryStore avec implémentation SQLite

  🧪 Vérification T1.2 :
  - Save → tuer process → relancer → Load → état identique
  - 20 goroutines en écriture simultanée (WAL mode) → pas de corruption
  - `PRAGMA integrity_check` → "ok"

---

### T1.3 — Logging JSON Lines ⚡ → T0.1

- [~] **T1.3.1** `internal/logging/logger.go` — thread-safe, rotation journalière

- [x] **T1.3.2** `internal/logging/redact.go` — masquage des secrets (Gap 12 fermé)
  ```go
  var secretPatterns = []*regexp.Regexp{
      regexp.MustCompile(`(?i)(password|passwd|token|api[_-]?key|secret|bearer)\s*[=:]\s*\S+`),
      regexp.MustCompile(`ghp_[A-Za-z0-9]{36}`),
      regexp.MustCompile(`sk-[A-Za-z0-9]{48}`),
      regexp.MustCompile(`AIza[A-Za-z0-9_-]{35}`),  // Google API key
  }

  func RedactSecrets(s string) string {
      for _, p := range secretPatterns {
          s = p.ReplaceAllStringFunc(s, func(match string) string {
              parts := strings.SplitN(match, "=", 2)
              if len(parts) == 1 { return "[REDACTED]" }
              return parts[0] + "=[REDACTED]"
          })
      }
      return s
  }
  // Utilisé sur TOUT texte loggé — aucune exception
  ```

- [ ] **T1.3.3** Logs sécurité séparés (Gap 12 fermé)
  ```go
  // internal/logging/security_logger.go
  // ~/.prometheus/security/events.jsonl — fichier séparé
  // Événements : command_blocked, command_confirmed, sast_finding, port_detected

  type SecurityLogger struct {
      file *os.File
      mu   sync.Mutex
  }

  func (sl *SecurityLogger) LogBlockedCommand(cmd string, reasons []string) {
      sl.write("command_blocked", map[string]any{
          "command": RedactSecrets(cmd),  // masquer si credentials dans la commande
          "reasons": reasons,
      })
  }
  ```

  🧪 Vérification T1.3 :
  - `grep -r "ghp_\|sk-\|password=" ~/.prometheus/logs/` → aucun résultat
  - `python3 -m json.tool ~/.prometheus/logs/$(date +%Y-%m-%d).jsonl` → valide
  - `~/.prometheus/security/events.jsonl` créé sur commande bloquée

---

### T1.4 — Vault Credentials ⚡ → T0.1

- [~] **T1.4.1** AES-256-GCM avec HKDF (Go stdlib, CGO_ENABLED=0) (Gap 5 fermé)
  ```go
  // internal/vault/vault.go
  import (
      "crypto/aes"
      "crypto/cipher"
      "crypto/hkdf"     // golang.org/x/crypto/hkdf non requis — crypto/hkdf standard en Go 1.22+
      "crypto/rand"
      "crypto/sha256"
  )

  // Dérivation de clé stable par machine
  func deriveKey() ([]byte, error) {
      machineID := getMachineID()  // /etc/machine-id || hostname
      username := os.Getenv("USER")
      if username == "" { username = os.Getenv("USERNAME") }  // Windows

      // HKDF avec SHA256 (Go 1.22 stdlib)
      h := hkdf.New(sha256.New,
          []byte(machineID+":"+username),  // ikm (input key material)
          []byte("prometheus-vault-v1"),    // salt
          []byte("prometheus-vault-key"),   // info
      )
      key := make([]byte, 32)  // AES-256
      _, err := io.ReadFull(h, key)
      return key, err
  }

  func getMachineID() string {
      // Linux : /etc/machine-id
      if data, err := os.ReadFile("/etc/machine-id"); err == nil {
          return strings.TrimSpace(string(data))
      }
      // macOS : system_profiler
      // Windows : HKLM\SOFTWARE\Microsoft\Cryptography\MachineGuid
      // Fallback : hostname
      h, _ := os.Hostname()
      return h
  }

  // Format fichier : [salt:16][nonce:12][ciphertext+GCM-tag]
  // Chiffrement AES-256-GCM
  func (v *Vault) Set(key, value string) error { /* ... */ }
  func (v *Vault) Get(key string) (string, error) { /* ... */ }
  func (v *Vault) List() []string { /* retourne les noms, pas les valeurs */ }
  func (v *Vault) Delete(key string) error { /* ... */ }
  ```

- [ ] **T1.4.2** Saisie masquée avec `golang.org/x/term`
  ```go
  // internal/ui/secret_input.go
  func ReadSecret(prompt string) (string, error) {
      fmt.Print(prompt)
      b, err := term.ReadPassword(int(syscall.Stdin))
      fmt.Println()  // saut de ligne après saisie
      return string(b), err
  }
  ```

- [ ] **T1.4.3** Auto-réutilisation dans le Blocked handler
  ```go
  // Avant d'afficher la question à l'utilisateur :
  // Si question contient "token", "password", "key", "secret"
  //   → chercher dans vault
  //   → si trouvé → Resume automatiquement sans afficher la question
  ```

  🧪 Vérification T1.4 :
  - `hexdump ~/.prometheus/vault.enc | grep -i "ghp_"` → 0 résultats
  - Vault.enc copié sur autre machine → déchiffrement échoue (machine-id différent)
  - Deuxième session : credentials réutilisés silencieusement

---

### T1.5 — TUI Bubbletea ⚡ → T0.1

- [ ] **T1.5.1** Layout 3 zones (Header + Conversation + Input)
- [ ] **T1.5.2** Indicateurs ASCII : `⟳ ✓ ⊙ ⚠ ✗ ● ○`
- [ ] **T1.5.3** Mode transparent Ctrl+L (logs exec en temps réel)
- [ ] **T1.5.4** Écran blocage avec guide contextuel + saisie masquée
- [ ] **T1.5.5** Ctrl+C → confirmation avant quit + cleanup gracieux
- [ ] **T1.5.6** 👁 Indicateur vision : `[👁 Screenshot analysé]`
- [ ] **T1.5.7** 🌐 Indicateur browser : `[🌐 https://...]`
- [ ] **T1.5.8** Compatibilité TERM=dumb (fallback sans couleurs)

  🧪 Vérification T1.5 :
  - Terminal 80×24 → pas de coupure
  - Terminal 40×20 (mobile Termux) → dégradation gracieuse
  - TERM=dumb → affichage texte simple, pas de crash
  - Android Termux : testé physiquement

---

### T1.6 — Graceful Shutdown ⚡ → T0.6 (Gap 17 fermé)

- [ ] **T1.6.1** Séquence de cleanup complète
  ```go
  // cmd/prometheus/main.go — cleanup au signal.NotifyContext
  // Ordre de cleanup :
  // 1. Arrêter de recevoir de nouvelles tâches
  // 2. Attendre que la tâche courante atteigne StatusBlocked ou StatusDone
  //    (max 10s, puis force)
  // 3. Sauvegarder l'état (SQLite)
  // 4. Fermer le browser (CDP close)
  // 5. Arrêter llama-server (SIGTERM puis SIGKILL après 5s)
  // 6. Fermer SQLite proprement
  // 7. Fermer les fichiers de log
  // 8. Exit 0

  // Si Ctrl+C une deuxième fois → force quit immédiat (exit 1)
  ```

  🧪 Vérification T1.6 :
  - Ctrl+C pendant exec → task sauvée → relancer → reprend à la même étape
  - Ctrl+C × 2 → exit immédiat
  - Après shutdown : `pgrep llama-server` → 0 processus

---

### T1.7 — Tests E2E Phase 1 ⚡

- [ ] **T1.7.1** E2E 1 : Création complète
  ```
  Input : "Crée une API REST Flask avec GET /hello et GET /info système"
  Vérifie :
    app.py créé
    Flask installé (pip install flask)
    Serveur démarré en background
    curl http://localhost:5000/hello → {"message":"hello"}
    curl http://localhost:5000/info → JSON avec OS/RAM/CPU
  ```

- [ ] **T1.7.2** E2E 2 : Vision sur résultat web
  ```
  Input : "Crée une page HTML avec un bouton bleu 'Cliquez-moi'
           et prends un screenshot pour vérifier"
  Vérifie :
    index.html créé
    Serveur local lancé
    Screenshot pris via vision
    LLM confirme le bouton bleu visible
  ```

- [ ] **T1.7.3** E2E 3 : Browser CDP
  ```
  Input : "Va sur https://example.com et dis-moi le titre et le premier paragraphe"
  Vérifie :
    Browser lancé (CDP)
    Navigation réussie
    HTML parsé
    Titre "Example Domain" mentionné dans la réponse
  ```

- [ ] **T1.7.4** E2E 4 : Blocage + vault
  ```
  Input : "Clone https://github.com/user/private-repo"
  Vérifie :
    Erreur auth détectée
    Question token affichée (saisie masquée)
    Après token fourni → clone réussi
    Token sauvé dans vault
    Deuxième run → vault utilisé (pas de question)
  ```

- [ ] **T1.7.5** E2E 5 : Android Termux (test physique)
  ```
  Hardware : Android 10+, Termux, 4GB RAM, sans Ollama
  Input : "Écris un script Python qui calcule les 20 premiers nombres premiers"
  Vérifie :
    Modèle local démarré (Phi-3 mini ou Qwen)
    Script créé et exécuté
    Résultat correct affiché
    RAM totale pendant exec : < 4GB
  ```

  🧪 CRITÈRES DE SUCCÈS PHASE 1 :
  - [ ] E2E 1-5 : tous réussis
  - [ ] Android 4GB RAM sans Ollama → fonctionne
  - [ ] Context Manager : 200 messages sur Phi-3 mini → zéro overflow
  - [ ] Vault : `hexdump vault.enc | grep -i "ghp_"` → 0 résultats
  - [ ] Graceful shutdown : tâche sauvée, reprise possible
  - [ ] `go test -race ./...` → 100% verts
  - [ ] `goleak.VerifyNone(t)` → 0 goroutine leak
  - [ ] `make size-check` → passe (< 30MB)

---

## PHASE 2 — THE EVOLUTION (4 semaines)
**→ Prérequis : Phase 1 VERTE**

---

### T2.1 — Capability Engine ⚡

- [ ] **T2.1.1** `internal/capabilities/registry.go` — 60+ capabilities intégrées
  ```go
  // Catégories :
  // Dev      : git, python3, node, npm, go, rust, java, mvn, gradle
  // Web      : chromium, firefox, curl, wget, httpie
  // Data     : jq, sqlite3, redis-cli, csvkit, pandas (pip)
  // DevOps   : docker, kubectl, terraform, ansible-playbook
  // Security : semgrep, nmap, trivy
  // Media    : ffmpeg, imagemagick, yt-dlp, whisper (pip)
  // Mobile   : adb, apktool, dex2jar
  // Vision   : scrot, import (imagemagick), tesseract
  // Forge    : python3, bash (toujours présents pour forger)
  // Browser  : playwright (optionnel, heavy — Node.js driver)

  type Capability struct {
      Name        string
      CheckCmd    string
      InstallCmds map[string]string  // os → commande d'install
      Type        string  // "system"|"pip"|"npm"|"cargo"|"forged"
      SizeMb      int
      Description string
  }
  ```

- [ ] **T2.1.2** ⚡ `engine.Ensure(ctx, name)` — cycle complet
  ```
  1. capabilities_cache SQLite → déjà installé → return nil
  2. Registry intégré → recette trouvée → install → sauvegarder
  3. Package managers : apt search, pip search, npm search, brew search
  4. Internet (si dispo) : GitHub, PyPI, npm, crates.io
  5. Capability FORGE → rien trouvé → LLM génère
  ```

- [ ] **T2.1.3** Playwright comme capability spéciale
  ```go
  // entry playwright dans le registry :
  {
      Name: "playwright",
      CheckCmd: "playwright --version 2>/dev/null",
      InstallCmds: {
          // playwright-go s'installe lui-même via playwright.Install()
          "go": "go run github.com/playwright-community/playwright-go/cmd/playwright install",
      },
      Type: "go-module",
      SizeMb: 60,  // Node.js driver + Chromium
      Description: "Browser automation avancé (SPA, formulaires complexes)",
  }
  // Après installation : charger playwright-go dynamiquement
  ```

  🧪 Vérification T2.1 :
  - Machine sans git → `"git clone ..."` → git installé → clone réussi
  - Playwright installé → `BrowserManager` bascule sur playwright
  - Capability mémorisée → deuxième session → pas de réinstallation

---

### T2.2 — Capability Forge ⚡ → T2.1

- [ ] **T2.2.1** Cycle complet : spécifier → générer → tester (max 3×) → packager
- [ ] **T2.2.2** Langages préférés : Python > bash > Go
- [ ] **T2.2.3** Tests automatiques : syntaxe + exécution + test fonctionnel
- [ ] **T2.2.4** Storage : `~/.prometheus/capabilities/forged/[name]/`
- [ ] **T2.2.5** Bloc C du system prompt mis à jour après forge

  🧪 Vérification T2.2 :
  - Outil inexistant → forgé → tests passent → disponible sessions suivantes
  - `python3 -m py_compile script.py` → succès

---

### T2.3 — Logs : Compression + Résumés + Archivage 🔁

- [ ] **T2.3.1** Compression zstd (`klauspost/compress/zstd`, pur Go)
- [ ] **T2.3.2** Résumés journaliers automatiques en markdown
- [ ] **T2.3.3** Archivage mensuel
- [ ] **T2.3.4** FTS5 sur les résumés (requêtes sémantiques)
- [ ] **T2.3.5** Requêtes temporelles en langage naturel (Gap 14 fermé)

  🧪 Vérification T2.3 :
  - Log 5MB → compressé < 1MB
  - "Sur quoi lundi dernier ?" → résumé correct
  - 1 an d'usage simulé → < 200MB sur disque

---

### T2.4 — Vision Avancée 👁 🔁

- [ ] **T2.4.1** Auto-capture après events clés (serveur lancé, fichier HTML créé)
- [ ] **T2.4.2** PDF → images → analyse (via `pdftoppm` installé par Capability Engine)
- [ ] **T2.4.3** Comparaison maquette → résultat (itération automatique)

---

### T2.5 — Browser : Playwright Optionnel 🌐 🔁

- [ ] **T2.5.1** Playwright installé via Capability Engine (si CDP insuffisant)
- [ ] **T2.5.2** `PlaywrightClient` implémente `BrowserClient`
- [ ] **T2.5.3** Simulateur mobile ADB (applications mobile-only)
- [ ] **T2.5.4** Gestion CAPTCHA : signaler à l'utilisateur (action=ask)

---

### T2.6 — Mise à jour Prometheus (Gap 18 fermé) 🔁

- [ ] **T2.6.1** `prometheus update` — vérifier la dernière version
  ```go
  // Appel à GitHub Releases API
  // Comparer version actuelle (buildinfo) avec dernière release
  // Si nouvelle version disponible : proposer
  func CheckUpdate(ctx context.Context) (*ReleaseInfo, error)
  ```

- [ ] **T2.6.2** Téléchargement du nouveau binaire
  ```go
  // Télécharger le binaire pour la plateforme courante
  // Vérifier SHA256 (dans checksums.txt de la release)
  // Remplacer le binaire courant
  // Redémarrer si confirmé par l'utilisateur
  func Update(ctx context.Context, release *ReleaseInfo) error
  ```

---

## PHASE 3 — THE IMMUNE SYSTEM (4 semaines)
**→ Prérequis : Phase 2 VERTE**

---

### T3.1 — Security Interceptor ⚡

- [ ] **T3.1.1** 200+ patterns dangereux avec scoring 0-100
- [ ] **T3.1.2** Actions : log(0-30) / log+exec(31-70) / confirm(71-90) / block(91+)
- [ ] **T3.1.3** LLM pour cas ambigus (score 40-70)
- [ ] **T3.1.4** Rate limiter exec (`golang.org/x/time/rate`, 10/s)
- [ ] **T3.1.5** Rate limiter LLM (`golang.org/x/time/rate`, 60/min) (Gap 15 fermé)
- [ ] **T3.1.6** Logs sécurité dans fichier séparé (Gap 12 fermé)

---

### T3.2 — SAST Natif Go ⚡

- [ ] **T3.2.1** 100+ règles : SQL injection, XSS, secrets hardcodés, HTTP non sécurisé, eval()
- [ ] **T3.2.2** Scan automatique sur chaque fichier créé par `create` ou après exec()
- [ ] **T3.2.3** Auto-correction : finding → LLM corrige → re-scan → vert
- [ ] **T3.2.4** Semgrep optionnel via Capability Engine (complément, pas prérequis)

---

### T3.3 — Sandbox Natif ⚡

- [ ] **T3.3.1** Niveau 1 : isolation workdir (toutes plateformes)
- [ ] **T3.3.2** Niveau 2 : Linux namespaces (CLONE_NEWNET|CLONE_NEWPID|CLONE_NEWNS + rlimits)
- [ ] **T3.3.3** Auto-détection du niveau disponible
- [ ] **T3.3.4** Pas de Docker requis — sandbox natif

---

### T3.4 — Scan Environnement + DAST 🔁

- [ ] **T3.4.1** Scan hôte : ports, permissions, credentials en clair
- [ ] **T3.4.2** DAST Go natif : headers HTTP manquants, endpoints sans auth
- [ ] **T3.4.3** Auto-patching avec confirmation + re-scan

---

### T3.5 — Vault v2 🔁

- [ ] **T3.5.1** Option passphrase (PBKDF2) en plus de machine-id
- [ ] **T3.5.2** Expiration des credentials configurables
- [ ] **T3.5.3** Audit log des accès vault

---

## PHASE 4 — OBSERVABILITY & PERFORMANCE (2 semaines)
**→ Prérequis : Phase 3 VERTE**

---

### T4.1 — Métriques Internes 🔁

- [ ] **T4.1.1** `SessionMetrics` : tasks, execs, LLM tokens, vision, browser, compactions
- [ ] **T4.1.2** `prometheus metrics` → affichage terminal clair
- [ ] **T4.1.3** `prometheus --profile` → pprof Go stdlib (mémoire + goroutines)

---

### T4.2 — Requêtabilité Sémantique 🔁

- [ ] **T4.2.1** Parser 20+ formes temporelles
- [ ] **T4.2.2** FTS5 : texte + browser sessions + vision captures
- [ ] **T4.2.3** `"Montre les screenshots de la session farmarket"` → liste les PNG

---

### T4.3 — Auto-Test de Prometheus 🔁

- [ ] **T4.3.1** Suite de 20 tâches de référence intégrées au binaire
  ```
  prometheus selftest  → exécute les 20 tâches, mesure succès/temps/RAM
  ```
- [ ] **T4.3.2** Rapport de performance par version

---

## PHASE 5 — UX & DISTRIBUTION (3 semaines)
**→ Prérequis : Phase 4 VERTE**

---

### T5.1 — Web UI Vanilla 🔁

- [ ] **T5.1.1** Serveur HTTP Go intégré — UNIQUEMENT sur `127.0.0.1` (Gap 19 fermé)
- [ ] **T5.1.2** Embed `//go:embed assets/static` — < 100KB total
- [ ] **T5.1.3** WebSocket bidirectionnel (temps réel)
- [ ] **T5.1.4** 👁 Affichage inline des screenshots dans la conversation
- [ ] **T5.1.5** 🌐 Indicateur browser actif
- [ ] **T5.1.6** Responsive mobile (Chrome Android)

---

### T5.2 — Distribution ⚡

- [ ] **T5.2.1** `install.sh` universel (Linux/macOS/Android Termux/Windows WSL)
  ```bash
  #!/bin/sh
  # Détecter OS + ARCH → télécharger bon binaire → placer dans ~/bin/
  # Vérifier SHA256 avant d'installer
  # Ajouter ~/bin au PATH si absent
  ```

- [ ] **T5.2.2** GoReleaser + GitHub Actions → release automatique sur tag

- [ ] **T5.2.3** SHA256 checksums pour chaque binaire (fichier `checksums.txt`)

- [ ] **T5.2.4** Vérifier SHA256 dans `install.sh` avant d'installer le binaire

---

### T5.3 — Documentation

- [~] **T5.3.1** `README.md` (EN) + `README.fr.md` (FR) — installation en 3 lignes
- [~] **T5.3.2** `docs/guide-android.md` — Termux, step-by-step, avec captures
- [~] **T5.3.3** `docs/guide-raspberry-pi.md` — serveur village communautaire
- [~] **T5.3.4** `docs/system-prompt.md` — personnalisation de l'IA
- [~] **T5.3.5** `docs/browser-control.md` — CDP vs Playwright, exemples
- [~] **T5.3.6** `docs/vision.md` — modèles supportés, configuration
- [~] **T5.3.7** `docs/capabilities.md` — capabilities intégrées + forge

---

### T5.4 — `prometheus` CLI complète

- [~] **T5.4.1** Sous-commandes
  ```
  prometheus                → démarrer en mode interactif (TUI)
  prometheus setup          → configurer le modèle LLM
  prometheus update         → mettre à jour Prometheus
  prometheus metrics        → afficher les métriques de session
  prometheus logs [DATE]    → afficher/rechercher les logs
  prometheus vault list     → lister les credentials
  prometheus selftest       → tester l'installation
  prometheus --web          → démarrer avec Web UI
  prometheus --model PATH   → utiliser un modèle personnalisé
  prometheus --provider api → forcer un provider cloud
  ```

---

## PHASE 6 — ADVANCED (roadmap ouverte)

- [ ] **T6.1** Multi-agents parallèles (orchestrateur + N goroutines)
- [ ] **T6.2** 👁 OCR amélioré (tesseract via Capability Engine)
- [ ] **T6.3** Mesh P2P chiffré (commander PC depuis téléphone)
- [ ] **T6.4** Interface voix (Whisper via llama-server, modèle GGUF whisper)
- [ ] **T6.5** Marketplace de capabilities forgées

---

## CROSS-CUTTING — CONTINU

### TC.1 — Tests et Qualité

- [ ] Coverage > 80% tous packages : `go test -coverprofile=cov.out ./...`
- [ ] `goleak.VerifyNone(t)` dans TestMain de chaque package
- [ ] `golangci-lint run` → zéro warning
- [ ] `CGO_ENABLED=0 go build ./...` passe toujours (vérifié en CI)
- [ ] Fuzzing : `ParseAction`, `ExtractFirstJSON`, `EstimateTokens`, `RedactSecrets`
- [ ] Suite de régression : 20 tâches de référence avant chaque release (Gap 20)
  ```
  5 création simple    5 avec erreurs et retry
  5 browser CDP        5 vision screenshot
  ```
- [ ] Mocks pour tous les providers LLM, vision, browser (tests sans réseau)

### TC.2 — Sécurité Continue

- [ ] `govulncheck ./...` avant chaque release
- [ ] Audit dépendances mensuel (`go list -json -m all`)
- [ ] `grep -r "ghp_\|sk-\|password=" ~/.prometheus/` → 0 résultats (test auto)
- [ ] Aucune image de screenshot loggée en clair
- [ ] Mise à jour immédiate si CVE critique dans une dépendance

### TC.3 — Performance Continue

- [ ] Benchmark 5 tâches de référence à chaque release
- [ ] Alerte si régression > 20% sur durée, RAM, tokens
- [ ] `make size-check` → passe en CI (< 30MB linux-arm64)
- [ ] Cold start < 3s mesuré automatiquement

### TC.4 — .golangci.yml Minimum

```yaml
run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - ineffassign
    - unused
    - gofmt
    - goimports
    - gosec       # sécurité Go
    - bodyclose   # fuite de response body
    - noctx       # contexte manquant
    - prealloc    # performance slice
    - unconvert   # conversions inutiles

linters-settings:
  gosec:
    excludes:
      - G304  # File path provided as taint input (normal pour Prometheus)
```

---

## RÉCAPITULATIF DÉPENDANCES EXACTES

```
go.mod PRINCIPAL (Phase 0+) :
  github.com/BurntSushi/toml          v1.3.2   config TOML
  modernc.org/sqlite                  v1.31.1  SQLite pure Go (CGO=0)
  golang.org/x/time                   v0.5.0   rate.Limiter
  golang.org/x/term                   v0.20.0  saisie masquée
  nhooyr.io/websocket                 v1.8.11  WebSocket CDP (pure Go)
  github.com/klauspost/compress       v1.17.8  zstd logs
  github.com/charmbracelet/bubbletea  v0.27.0  TUI (Phase 1)
  github.com/charmbracelet/lipgloss   v1.0.0   styles TUI
  go.uber.org/goleak                  v1.3.0   goroutine leak tests

OUTILS BUILD (pas dans go.mod) :
  github.com/golangci/golangci-lint   → golangci-lint
  golang.org/x/vuln/cmd/govulncheck   → govulncheck
  github.com/goreleaser/goreleaser/v2 → goreleaser

OPTIONNEL VIA CAPABILITY ENGINE (pas dans go.mod principal) :
  github.com/playwright-community/playwright-go  browser avancé
  (télécharge aussi Node.js driver ~50MB au premier usage)

BUILD TOOLS EXTERNES (pas Go) :
  cmake + ninja         compiler llama-server
  android-ndk-r27b+     cross-compiler pour Android ARM64
  aarch64-linux-gnu-gcc cross-compiler Linux ARM64
  mingw-w64             cross-compiler Windows

CONTRAINTES ABSOLUES :
  CGO_ENABLED=0         toujours, pour toutes les builds Go
  GOOS=linux GOARCH=arm64 → fonctionne sur Android Termux (Linux userland)
  Web UI : écouter uniquement 127.0.0.1 (jamais 0.0.0.0)
  playwright-go : jamais dans le go.mod principal
```

---

## RÉCAPITULATIF CRITÈRES DE SUCCÈS FINAUX

### Zéro dépendance externe après lancement
- [ ] Machine vierge → `./prometheus setup` → modèle téléchargé → 100% offline ensuite
- [ ] Pas d'Ollama, pas de Docker, pas de Python, pas de Node.js requis au démarrage
- [ ] Capability Engine installe ce dont Prometheus a besoin en cours de route

### Fonctionnel
- [ ] Android Termux 4GB RAM → 100% offline
- [ ] `"Crée le jeu Snake en Python"` → code créé et fonctionnel
- [ ] `"Va sur ce site et remplis ce formulaire"` → fait sans API externe (CDP)
- [ ] `"Voici une maquette, reproduis-la"` → UI générée, comparée visuellement, corrigée
- [ ] `"Cet outil n'existe pas, crée-le"` → script forgé et fonctionnel
- [ ] `"Sur quoi lundi dernier ?"` → réponse depuis les archives compressées
- [ ] Credential demandé une fois → jamais redemandé (vault)
- [ ] Interruption Ctrl+C → reprise exacte à la même étape

### Performance
- [ ] Binaire linux-arm64 : < 30MB
- [ ] Cold start : < 3s
- [ ] RAM sans LLM : < 150MB
- [ ] Logs 1 an : < 200MB compressé
- [ ] Vision : analyse screenshot < 5s (modèle local)

### Sécurité
- [ ] `hexdump ~/.prometheus/vault.enc | grep -i "ghp_"` → 0 résultats
- [ ] Fork bomb dans sandbox → tué en < 30s
- [ ] SQL injection dans code généré → détectée et corrigée automatiquement
- [ ] Web UI : `ss -tlnp | grep 8080` → écoute UNIQUEMENT sur 127.0.0.1

### Distribution
- [ ] `curl -L install.sh | sh` → fonctionne sur Ubuntu, macOS, Android Termux
- [ ] SHA256 vérifié avant installation
- [ ] CI/CD → build + test + release automatisés sur tag

---

*Prometheus · Scaffold Complet Final v5.0*
*Tous les Gaps Fermés · Dépendances Exactes · Ordre Correct*
*Conçu à Lomé · Souveraineté Numérique Africaine · Zéro Dépendance Externe*
