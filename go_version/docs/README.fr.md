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

### Option 2: Installateur Windows (Recommandé)

Téléchargez et exécutez `prometheus-setup.exe` depuis [GitHub Releases](https://github.com/kora-ai-lab/prometheus/releases/latest). L'installateur va :
- Installer Prometheus dans Program Files
- L'ajouter automatiquement à votre PATH
- Créer des raccourcis dans le Menu Démarrer
- Créer un raccourci sur le bureau (optionnel)

### Option 3: PowerShell (Windows)

```powershell
irm https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.ps1 | iex
```

### Option 4: Télécharger l'exe directement (Windows)

Téléchargez `prometheus-windows-amd64.exe` depuis [GitHub Releases](https://github.com/kora-ai-lab/prometheus/releases/latest) et ajoutez-le à votre PATH.

> **Note :** Windows Defender peut détecter l'exécutable comme faux positif. C'est attendu pour les binaires non signés. Utilisez l'installateur MSI pour une meilleure compatibilité, ou ajoutez une exclusion si nécessaire.

### Option 5: Compiler depuis les sources

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

## Dépannage

### Faux positif Windows Defender

Si Windows Defender détecte Prometheus comme malware :

1. **Utilisez l'installateur Windows** - `prometheus-setup.exe` a une meilleure réputation que les exécutables bruts
2. **Ajoutez une exclusion :**
   - Ouvrez Sécurité Windows → Protection contre les virus et menaces → Gérer les paramètres
   - Faites défiler jusqu'à Exclusions → Ajouter ou supprimer des exclusions
   - Ajoutez le dossier d'installation (ex: `C:\Program Files\Prometheus`)
3. **Signalez à Microsoft :** Rapportez le faux positif sur https://www.microsoft.com/en-us/wdsi/filesubmission

### Problèmes d'installation

**"commande prometheus introuvable"**
- Fermez et rouvrez votre terminal après l'installation
- Ou ajoutez manuellement au PATH : `setx PATH "%PATH%;C:\Program Files\Prometheus"`

**Double-clic sur l'exe ferme immédiatement**
- Utilisez l'installateur MSI à la place
- Ou lancez depuis PowerShell : `.\prometheus-windows-amd64.exe --help`

**Interface web inaccessible**
- Vérifiez si le port 8080 est disponible : `netstat -ano | findstr :8080`
- Utilisez un autre port : `prometheus --web --port 9090`

## Fonctionnalités

- Génération de code par LLM
- Exécution de commandes
- Opérations sur fichiers
- Automatisation du navigateur
- Auto-mise à jour (`prometheus --update`)