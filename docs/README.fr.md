# Prometheus

Un runtime agent IA-first avec sécurité, observabilité et auto-évolution.

## Démarrage rapide

```bash
# Installer
curl -fsSL https://raw.githubusercontent.com/prometheus-dev/prometheus/main/scripts/install.sh | sh

# Exécuter
prometheus "Votre objectif ici"
```

## Installation

### Option 1: curl | sh

```bash
curl -fsSL https://raw.githubusercontent.com/prometheus-dev/prometheus/main/scripts/install.sh | sh
```

### Option 2: Compiler depuis les sources

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