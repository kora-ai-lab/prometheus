# Prometheus

Un runtime agent IA-first avec sécurité, observabilité et auto-évolution.

## Démarrage rapide

**macOS/Linux :**
```bash
curl -fsSL https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.sh | sh
prometheus "Votre objectif ici"
```

**Windows (PowerShell) :**
```powershell
irm https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.ps1 | iex
prometheus "Votre objectif ici"
```

## Installation

### Option 1: curl | sh (macOS/Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.sh | sh
```

### Option 2: PowerShell (Windows)

```powershell
irm https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.ps1 | iex
```

### Option 3: Télécharger l'exe directement (Windows)

Téléchargez `prometheus-windows-amd64.exe` depuis [GitHub Releases](https://github.com/kora-ai-lab/prometheus/releases/latest) et ajoutez-le à votre PATH.

### Option 4: Compiler depuis les sources

```bash
git clone https://github.com/kora-ai-lab/prometheus
cd prometheus/go_version
go build -ldflags="-s -w" -o prometheus ./cmd/prometheus
```

## Configuration

Définir la clé API :

```bash
export GROQ_API_KEY=sk-...
```

Ou utiliser l'interface web sur http://localhost:8080

## Utilisation

### CLI

```bash
prometheus "Créer une application hello world"
prometheus --web  # Démarrer l'interface web
```

### Interface web

1. Exécuter `prometheus --web`
2. Ouvrir http://localhost:8080
3. Entrer votre objectif

## Sécurité

- Exécution en sandbox
- Auto-confirmation pour les commandes risqués
- Secrets masqués dans les logs

## Fonctionnalités

- Génération de code par LLM
- Exécution de commandes
- Opérations sur fichiers
- Automatisation du navigateur
- Auto-mise à jour (`prometheus --update`)