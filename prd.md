# 🔥 PROMETHEUS
## Product Requirements Document — Architecture Complète v2.0
### Souveraineté Numérique · Africa-First · Zéro Dépendance Externe

---

> **Ce document est la source de vérité absolue de Prometheus.**
> Chaque décision d'architecture, chaque phase, chaque ligne de code
> doit être cohérente avec ce qui est écrit ici.

---

## TABLE DES MATIÈRES

1. [Vision & Philosophie](#1-vision--philosophie)
2. [Critique & Corrections v2](#2-critique--corrections-v2)
3. [Architecture Globale](#3-architecture-globale)
4. [Le Noyau Liquide — Zéro Dépendance](#4-le-noyau-liquide--zéro-dépendance)
5. [Gestion du Contexte LLM](#5-gestion-du-contexte-llm)
6. [Système de Log — Mémoire Perpétuelle](#6-système-de-log--mémoire-perpétuelle)
7. [Capability Engine — Créer ce qui n'existe pas](#7-capability-engine--créer-ce-qui-nexiste-pas)
8. [Sécurité — Immunologie sans Sandbox externe](#8-sécurité--immunologie-sans-sandbox-externe)
9. [Terminal UX — Interface Humaine](#9-terminal-ux--interface-humaine)
10. [Stack Technique Finale](#10-stack-technique-finale)
11. [Phases de Développement](#11-phases-de-développement)
12. [Spécifications Non-Fonctionnelles](#12-spécifications-non-fonctionnelles)
13. [Distribution & Installation](#13-distribution--installation)

---

## 1. VISION & PHILOSOPHIE

### La Vision
Prometheus est un **runtime agentique à noyau liquide** : un binaire unique, autonome, sans dépendance externe, capable de se doter lui-même de tous les pouvoirs dont il a besoin pour accomplir n'importe quelle tâche — en local, offline, sur n'importe quel appareil africain ou mondial.

### Les 7 Principes Fondamentaux

**P1 — Un seul primitif de départ**
exec() est le seul pouvoir inné. Tout le reste — outils, capacités, sécurité, navigateur, simulateur — est acquis à la volée par Prometheus lui-même via ce primitif.

**P2 — Le Noyau est Liquide**
Le binaire Prometheus ne dépend d'aucun outil externe pour démarrer et évoluer. Il n'a pas besoin de Python, Node, Docker, ou même d'Ollama préinstallé. Il peut tout installer lui-même depuis zéro si nécessaire. Pas de sandbox système requis : la sécurité est implémentée nativement dans le noyau.

**P3 — La tâche ne meurt jamais**
Bloqué sur credentials manquants → il pause, demande, reprend exactement là où il s'est arrêté. Capability manquante → il l'installe, la mémorise, continue. Erreur → il lit, comprend, corrige, relance. La boucle tourne jusqu'à DONE.

**P4 — Contexte géré intelligemment**
Peu importe la taille de la context window du modèle utilisé (4K tokens pour un modèle léger, 1M pour un modèle cloud), Prometheus maintient une qualité constante par compaction progressive, résumé sémantique et mémoire à long terme externalisée dans SQLite.

**P5 — Mémoire perpétuelle compressée**
Prometheus se souvient de tout, pour toujours, sans jamais devenir lourd. Les logs vieux de plus de 7 jours sont compressés (zstd). Les sessions sont archivées. Tout est requêtable en langage naturel. "Qu'est-ce qu'on avait fait lundi dernier ?" reçoit une vraie réponse.

**P6 — Simplicité radicale pour l'utilisateur**
Zéro slash command. Zéro mode à activer. Zéro configuration requise au démarrage. Tu parles, il agit. La complexité est entièrement absorbée par la machine, jamais exposée à l'utilisateur.

**P7 — Africa-First Design**
- Fonctionne sur un Android 4GB RAM (Tecno, Infinix, Samsung A-series)
- Binaire unique < 20MB
- Mode offline total
- Réseau intermittent géré nativement
- Zéro abonnement cloud obligatoire
- Modèle communautaire (1 Raspberry Pi = Prometheus pour un village)

---

## 2. CRITIQUE & CORRECTIONS v2

### Ce qui était juste dans la v1
- Découpage en phases clair et progressif
- Philosophie immunologique forte
- CDP pour navigateur natif sans API externe
- Parallélisme par goroutines Go
- SQLite comme mémoire légère

### Failles corrigées dans v2

**Faille 1 — Sandbox externe = dépendance**
Docker comme sandbox crée une dépendance externe. Solution v2 : sandbox natif Go via seccomp + namespaces Linux côté noyau, avec fallback gracieux sur les systèmes ne le supportant pas (Android, Windows). Sur ces systèmes, la sécurité repose sur l'audit de commandes avant exécution + rate limiting.

**Faille 2 — Gestion de contexte absente**
La v1 ne spécifiait pas comment gérer les modèles légers avec 4K-32K tokens de contexte. Solution v2 : Context Manager avec 3 niveaux (Hot / Warm / Cold) et compaction automatique.

**Faille 3 — Logs sans stratégie d'archivage**
Des logs infinis deviennent inutilisables. Solution v2 : système de rotation/compression/archivage détaillé avec requêtabilité sémantique.

**Faille 4 — Capability vide si non trouvée**
Que se passe-t-il si une capability n'existe nulle part ? La v1 ne répondait pas. Solution v2 : le Capability Forge — Prometheus synthétise lui-même la capability via LLM si elle n'existe pas dans les registres.

**Faille 5 — Terminal non user-friendly**
stdin/stdout brut est hostile. Solution v2 : TUI (Terminal User Interface) native Go avec zones de dialogue, progression, logs, confirmations — sans dépendance Tauri ou Electron.

**Faille 6 — Rust proposé sans raison**
Rust = courbe d'apprentissage de 3 mois, binaires plus lourds sur ARM, cross-compilation plus complexe. Go compile nativement pour Android ARM64 en une commande. Go est le bon choix.

---

## 3. ARCHITECTURE GLOBALE

```
┌─────────────────────────────────────────────────────────────┐
│                    UTILISATEUR                              │
│            (Langage naturel, n'importe où)                  │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                   TUI / INTERFACE                           │
│   Terminal riche (bubbletea) · Web UI · Mobile WebView      │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                  LLM ADAPTER LAYER                          │
│   Ollama · Anthropic API · Google API · GGUF local direct   │
│   Interface unifiée : ModelProvider (text in → text out)    │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                   CONTEXT MANAGER                           │
│   Hot Buffer · Compaction · Résumé sémantique · Cold Store  │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                  CORE RUNTIME — LE NOYAU                    │
│                                                             │
│  ┌─────────────────┐  ┌──────────────────────────────────┐ │
│  │   TASK LOOP     │  │    SECURITY INTERCEPTOR          │ │
│  │                 │  │                                  │ │
│  │  Think (LLM)    │  │  Pre-exec audit                  │ │
│  │  Execute        │  │  Rate limiting (10 exec/s)       │ │
│  │  Observe        │  │  Dangerous ops confirmation      │ │
│  │  Decide         │  │  Native sandbox (seccomp)        │ │
│  │  Retry / Block  │  │  SAST inline (code produit)      │ │
│  │  Done           │  │  Auto-patch                      │ │
│  └────────┬────────┘  └──────────────────────────────────┘ │
│           │                                                 │
│  ┌────────▼────────┐  ┌──────────────────────────────────┐ │
│  │  ORCHESTRATOR   │  │    MEMORY MANAGER                │ │
│  │                 │  │                                  │ │
│  │  Multi-workers  │  │  Task state (SQLite)             │ │
│  │  goroutines     │  │  User prefs                      │ │
│  │  Task queue     │  │  Learned patterns                │ │
│  │  Coordination   │  │  Log perpétuel (zstd)            │ │
│  └────────┬────────┘  └──────────────────────────────────┘ │
│           │                                                 │
└───────────┼─────────────────────────────────────────────────┘
            │
┌───────────▼─────────────────────────────────────────────────┐
│                  CAPABILITY ENGINE                          │
│                                                             │
│  Discovery → Install → Forge (si absent) → Mémoriser        │
│                                                             │
│  ~/.prometheus/capabilities/                                │
│    git/          docker/        semgrep/                    │
│    chromium/     adb/           [forgées par LLM]/          │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. LE NOYAU LIQUIDE — ZÉRO DÉPENDANCE

### Le problème du sandbox externe

Docker comme sandbox = Prometheus dépend de Docker. Si Docker n'est pas installé, pas de sandbox. C'est une contradiction avec le principe "Zéro Dépendance Externe".

### La solution : Sandbox Natif Go

Prometheus implémente sa propre isolation directement dans le noyau Go, sans dépendance externe, via 3 mécanismes progressifs selon la plateforme :

#### Niveau 1 — Audit Pre-execution (toutes plateformes)
Avant toute exécution, le SecurityInterceptor analyse la commande :

```
Commande reçue → Tokenisation
              → Pattern matching (blacklist)
              → Analyse LLM (si commande ambiguë)
              → Score de risque 0-100
              → Score < 30 : exec direct
              → Score 30-70 : confirmation user
              → Score > 70 : refus + explication
```

Patterns dangereux détectés nativement (sans outil externe) :
- rm -rf / et variantes
- curl * | sh (exécution aveugle)
- chmod 777 / 
- Modification de /etc/passwd, /etc/hosts 
- Fork bombs: :(){ :|:& };: 
- Accès credentials système

#### Niveau 2 — Namespace Isolation (Linux / Android)
Sur Linux et Android, Prometheus peut créer des processus isolés sans Docker :

```go
// Isolation native via syscalls Linux - zéro dépendance
func (s *Sandbox) RunIsolated(cmd string) (*Result, error) {
    attr := &syscall.SysProcAttr{
        // Nouveau namespace réseau (pas d'accès réseau sans permission)
        Cloneflags: syscall.CLONE_NEWNET |
                    // Nouveau namespace processus (pas de kill des parents)
                    syscall.CLONE_NEWPID |
                    // Nouveau namespace filesystem
                    syscall.CLONE_NEWNS,
        // Limites mémoire et CPU
        Rlimit: []syscall.Rlimit{
            {Type: syscall.RLIMIT_AS, Cur: 512 * 1024 * 1024},  // 512MB max
            {Type: syscall.RLIMIT_CPU, Cur: 30},                  // 30s CPU max
        },
    }
    c := exec.Command("sh", "-c", cmd)
    c.SysProcAttr = attr
    return c.CombinedOutput()
}
```

#### Niveau 3 — Fallback Gracieux (Windows / Android sans root)
Quand l'isolation native n'est pas disponible, Prometheus bascule automatiquement sur :
- Répertoire de travail isolé (chroot simulé par changement de workdir)
- Variable d'environnement minimale (PATH contrôlé)
- Timeout strict sur toutes les commandes
- Log de tout ce qui est exécuté

#### Résultat
Prometheus est son propre sandbox. Aucun outil externe requis. Sur Android Termux : fonctionne parfaitement. Sur Linux complet : isolation maximale. Sur Windows : isolation partielle mais sécurisée. **Zero dépendance externe pour la sécurité.**

---

## 5. GESTION DU CONTEXTE LLM

### Le problème fondamental

Les modèles légers (Phi-3 mini, Gemma 2B, Llama 3.2 3B) ont des context windows de 4K à 32K tokens. Une session de développement de 2 heures peut facilement dépasser 100K tokens d'historique. Sans gestion intelligente, deux problèmes :
- **Troncature silencieuse** : le début de la conversation disparaît, le modèle "oublie" l'objectif
- **Dégradation de qualité** : trop de contexte non pertinent noie les informations importantes

### Architecture 3 niveaux : Hot / Warm / Cold

```
┌─────────────────────────────────────────────────────────────┐
│  HOT BUFFER — dans la context window du LLM                 │
│  Contenu : messages récents + system prompt + tâche active  │
│  Taille : 60% de la context window disponible               │
│  Format : messages bruts (role + content)                   │
│  Mise à jour : temps réel                                   │
└────────────────────────┬────────────────────────────────────┘
                         │ compaction tous les N messages
┌────────────────────────▼────────────────────────────────────┐
│  WARM STORE — SQLite en mémoire vive                        │
│  Contenu : résumés compactés des échanges précédents        │
│  Taille : ~10MB max en RAM                                  │
│  Format : JSON structuré (résumé + entités clés + décisions)│
│  Mise à jour : à chaque compaction                          │
└────────────────────────┬────────────────────────────────────┘
                         │ archivage après session
┌────────────────────────▼────────────────────────────────────┐
│  COLD STORE — SQLite sur disque                             │
│  Contenu : toutes les sessions passées compressées          │
│  Taille : illimitée (compression zstd 80-90% gain)          │
│  Format : sessions archivées + index sémantique             │
│  Accès : sur demande explicite                              │
└─────────────────────────────────────────────────────────────┘
```

### Algorithme de compaction

La compaction se déclenche automatiquement quand le Hot Buffer atteint 70% de la context window disponible.

```
TRIGGER: Hot Buffer > 70% context window

COMPACTION PROCESS:
1. Identifier les N derniers messages à garder intacts
   (N = messages depuis le dernier checkpoint, max 10)

2. Prendre les messages antérieurs (à compacter)

3. Envoyer au LLM avec ce prompt:
   "Résume cette conversation en conservant:
    - L'objectif principal de la tâche
    - Les décisions prises et leur justification
    - Les erreurs rencontrées et leurs solutions
    - L'état actuel du projet
    - Les credentials/configs mentionnés
    - Les prochaines étapes planifiées
    Format: JSON structuré, max 500 tokens"

4. Stocker le résumé dans Warm Store

5. Remplacer les messages compactés par:
   [CONTEXT_SUMMARY_REF: {id}] dans Hot Buffer

6. Le Hot Buffer reprend avec:
   - System prompt (toujours)
   - Résumé injecté comme message "context" 
   - N derniers messages intacts
   - Nouvelle tâche en cours
```

### Injection de contexte selon le modèle

```
DÉTECTION DE LA CONTEXT WINDOW:
- Ollama: API /api/show → context_length
- Anthropic: connu statiquement par modèle
- Google: connu statiquement par modèle
- GGUF local: lu depuis les métadonnées du fichier

STRATÉGIE PAR TAILLE:
< 8K tokens (très petit):
  → Hot: 40% (3K) — messages récents uniquement
  → Résumé ultra-condensé (200 tokens max)
  → Compaction tous les 5 messages

8K - 32K tokens (petit):
  → Hot: 60% — messages + contexte enrichi
  → Résumé standard (500 tokens)
  → Compaction tous les 15 messages

32K - 128K tokens (moyen):
  → Hot: 70% — historique complet récent
  → Résumé détaillé (1000 tokens)
  → Compaction tous les 30 messages

> 128K tokens (grand — Claude, GPT-4):
  → Hot: 80% — quasi pas de compaction nécessaire
  → Compaction tous les 100 messages
```

### Mode Découverte / Plug-and-play

Quand Prometheus démarre sur une machine inconnue, il lance automatiquement un **scan d'environnement** qui enrichit le contexte système :

```
SÉQUENCE DE DÉCOUVERTE (première fois et à chaque démarrage):

exec("uname -a")              → OS et kernel
exec("uname -m")              → Architecture (arm64, x86_64)
exec("nproc")                 → Nombre de CPUs
exec("free -m")               → RAM disponible
exec("df -h .")               → Espace disque
exec("which git python3 node npm docker chromium firefox adb")
                              → Outils présents
exec("ollama list 2>/dev/null") → Modèles LLM locaux disponibles
exec("ping -c 1 8.8.8.8 2>&1") → Internet disponible?
exec("cat /etc/os-release 2>/dev/null") → Distribution Linux

RÉSULTAT: EnvironmentProfile {
    os: "linux-android"
    arch: "arm64"
    ram_mb: 3840
    disk_gb: 45
    tools: ["git", "python3"]   // seuls ceux présents
    llm_available: ["phi3:mini"]
    internet: false
    package_manager: "apt"
}

Ce profil est:
1. Injecté dans le system prompt (condensé, ~100 tokens)
2. Stocké dans memory.db pour les sessions suivantes
3. Re-scanné au démarrage (changes détectés incrementalement)
```

---

## 6. SYSTÈME DE LOG — MÉMOIRE PERPÉTUELLE

### Objectif
Prometheus doit se souvenir de tout pour toujours, être requêtable en langage naturel, et ne jamais occuper plus de 50MB d'espace actif — même après 5 ans d'utilisation continue.

### Structure des logs

Chaque événement loggé est atomique et structuré :

```json
{
  "ts": "2025-01-15T14:32:07.123Z",
  "session": "sess_2025-01-15_001",
  "task_id": "task_abc123",
  "level": "exec",
  "event": {
    "type": "command_executed",
    "command": "git clone https://github.com/farmarket/app",
    "exit_code": 0,
    "duration_ms": 3420,
    "stdout_preview": "Cloning into 'app'...\nReceiving objects: 100%",
    "stderr": ""
  }
}
```

Types d'événements capturés :
- session_start / session_end 
- user_message — ce que l'utilisateur a dit
- llm_response — ce que le LLM a décidé (résumé, pas le texte brut)
- command_executed — commande + résultat + durée
- capability_installed — quelle capability installée
- task_blocked — pourquoi bloqué + ce qui a été demandé
- task_completed — objectif accompli + durée totale
- error_and_fix — erreur rencontrée + solution appliquée
- security_event — commande bloquée ou confirmée

### Cycle de vie du log

```
HEURE 0 à H+24:
  Fichier actif: ~/.prometheus/logs/2025-01-15.jsonl
  Format: JSON Lines (une ligne = un événement)
  Taille typique: 2-5MB par jour actif

FIN DE JOURNÉE (minuit automatique):
  1. Fermeture du fichier actif
  2. Génération du RÉSUMÉ JOURNALIER:
     → Demander au LLM: "Résume cette journée de travail:
        - Sur quoi avons-nous travaillé?
        - Qu'est-ce qui a été accompli?
        - Quelles difficultés rencontrées?
        - Quelles décisions importantes prises?
        Format: texte naturel, 200-400 mots"
     → Stocker le résumé: ~/.prometheus/summaries/2025-01-15.md
  3. Compression du log brut:
     → 2025-01-15.jsonl → 2025-01-15.jsonl.zst (gain ~85%)
     → Taille: 5MB → ~750KB
  4. Index mis à jour (SQLite)

APRÈS 7 JOURS:
  → Archives déplacées vers ~/.prometheus/archive/2025-01/
  → Index de recherche maintenu (full-text sur les résumés)
  → Les .jsonl.zst restent accessibles mais non chargés en RAM

APRÈS 90 JOURS:
  → Archives groupées par mois: 2025-01.tar.zst
  → Taille typique: 1 mois = 50-150MB brut → 3-10MB compressé
  → Les résumés journaliers (.md) restent non-compressés (petits)

APRÈS 1 AN:
  → Archive annuelle optionnelle
  → Résumé annuel généré par LLM
  → Résumés mensuels toujours accessibles
```

### Requêtabilité sémantique

La question "Sur quoi avions-nous travaillé lundi dernier ?" fonctionne ainsi :

```
PROCESSUS DE REQUÊTE TEMPORELLE:

1. Parser l'intention temporelle:
   "lundi dernier" → calculer la date (2025-01-13)

2. Charger le résumé journalier:
   ~/.prometheus/summaries/2025-01-13.md
   (toujours disponible, non-compressé, ~5KB)

3. Si besoin de détails → décompresser le log:
   ~/.prometheus/archive/2025-01/2025-01-13.jsonl.zst
   Décompression: ~100ms

4. Répondre avec le résumé + détails si demandé
```

### Espace occupé — Projections

```
UTILISATION QUOTIDIENNE (utilisateur actif 4h/jour):
  Log brut:        ~3MB/jour
  Log compressé:   ~450KB/jour
  Résumé:          ~10KB/jour (texte markdown)

PROJECTIONS CUMULÉES (compressé uniquement):
  1 semaine:       ~3MB
  1 mois:          ~14MB
  1 an:            ~170MB
  5 ans:           ~850MB

ESPACE ACTIF EN RAM (jamais plus que ça):
  Session courante:   ~5MB max
  Warm store SQLite:  ~10MB max
  Index de recherche: ~5MB max
  TOTAL ACTIF:        < 20MB RAM en permanence
```

### Format des résumés journaliers

```markdown
# Journal Prometheus — 2025-01-13 (Lundi)

## Résumé de la journée
Journée principalement dédiée au démarrage du projet FarMarket,
une application de vente directe agriculteur-acheteur pour l'Afrique
de l'Ouest. Structure de base établie, backend API opérationnel.

## Accompli aujourd'hui
- Création de la structure du projet (backend/ + mobile/ + docs/)
- Initialisation git + premier commit
- API REST Flask opérationnelle (endpoints: /products, /users, /orders)
- Base de données SQLite configurée avec schéma initial
- Mode offline implémenté pour les listings produits

## Difficultés rencontrées
- Credentials GitHub manquants pour le module auth → fournis par l'utilisateur
- Dépendance better-sqlite3 incompatible ARM64 → remplacée par sqlite3 natif

## Décisions importantes
- Architecture choisie: SQLite offline-first (pas PostgreSQL)
- Pas d'authentification OAuth pour la v1 (trop complexe, JWT suffit)

## État du projet à la fin de journée
Backend: 70% | Mobile: 10% | Tests: 0% | Docs: 30%

## Prochaines étapes prévues
- UI mobile React Native (écran listing + formulaire ajout produit)
- Intégration données prix FAO/GIEWS pour céréales
- Tests unitaires API

## Stats
Durée session: 4h12min | Commandes exécutées: 147 | Erreurs corrigées: 8
```

---

## 7. CAPABILITY ENGINE — CRÉER CE QUI N'EXISTE PAS

### Le cycle complet de gestion des capacités

```
BESOIN DÉTECTÉ PAR LE LLM
         │
         ▼
┌─────────────────────┐
│  1. LOCAL CACHE     │  Déjà installé sur cette machine ?
│  ~/.prometheus/     │  → OUI: utiliser directement
│  capabilities/      │  → NON: continuer
└─────────┬───────────┘
          │ NON
          ▼
┌─────────────────────┐
│  2. DISCOVERY       │  Chercher avec les outils natifs
│                     │  which, apt search, pip search,
│                     │  npm search, brew search, cargo search
│                     │  → TROUVÉ: installer + mémoriser
│                     │  → NON TROUVÉ: continuer
└─────────┬───────────┘
          │ NON TROUVÉ
          ▼
┌─────────────────────┐
│  3. WEB SEARCH      │  Si internet disponible:
│                     │  → Navigateur natif CDP
│                     │  → Chercher sur GitHub, PyPI, npm
│                     │  → Trouver l'outil adapté
│                     │  → Installer
│                     │  → TOUJOURS NON TROUVÉ: continuer
└─────────┬───────────┘
          │ TOUJOURS ABSENT
          ▼
┌─────────────────────┐
│  4. CAPABILITY      │  L'outil n'existe pas ou n'est pas
│     FORGE 🔥        │  adapté → Prometheus le crée lui-même
│                     │
│  "Je vais écrire    │  Processus:
│  le script moi-     │  a) LLM décrit la capability nécessaire
│  même"              │  b) LLM génère le code (Python/bash/Go)
│                     │  c) Tests automatiques
│                     │  d) Stockage dans capabilities/custom/
│                     │  e) Mémorisation avec métadonnées
└─────────────────────┘
```

### Structure d'une capability

Chaque capability installée ou forgée a le même format :

```
~/.prometheus/capabilities/
  git/
    meta.toml          # Métadonnées
    check.sh           # "Est-ce que cet outil est disponible?"
    install.sh         # "Comment l'installer sur différents OS"
    use_example.sh     # Exemples d'utilisation pour le LLM

  custom/
    farmarket-sync/
      meta.toml
      sync.py          # Script forgé par Prometheus lui-même
      test.sh
```

```toml
# meta.toml d'une capability
[capability]
name = "git"
version = "2.43.0"
type = "system"          # system | custom | forged
installed_at = "2025-01-15T10:00:00Z"
platform = "linux-arm64"
source = "apt"           # apt | brew | pip | npm | cargo | forged

[usage]
description = "Versioning de code source"
when_to_use = "Clone de repos, commit, push, gestion de branches"
install_cost_seconds = 45
size_mb = 12

[platforms]
linux = "apt install git"
macos = "brew install git"
android_termux = "pkg install git"
windows = "winget install git"
```

### Capability Forge — Créer l'inexistant

Le Forge s'active quand aucune solution existante n'est trouvée :

```
EXEMPLE: Prometheus a besoin d'un parser de prix agricoles FAO/GIEWS
         en format XML non-standard spécifique à l'Afrique de l'Ouest.
         Aucun package ne fait exactement ça.

FORGE PROCESS:

1. LLM génère la spécification:
   "Outil nécessaire: Parser XML FAO/GIEWS
    Input: URL ou fichier XML FAO
    Output: JSON {commodity, price, unit, market, date, country}
    Contraintes: offline-capable, retry réseau, cache local"

2. LLM génère le code:
   → Python script ~150 lignes
   → Tests unitaires inclus

3. Tests automatiques:
   exec("python test_fao_parser.py")
   → Si tests passent: capability forgée ✓
   → Si tests échouent: LLM corrige, retry jusqu'à 3 fois

4. Stockage:
   ~/.prometheus/capabilities/custom/fao-parser/
   → meta.toml (type = "forged")
   → fao_parser.py
   → test_fao_parser.py

5. Mémorisation:
   memory.db → "fao-parser: utilisé pour récupérer prix marchés Afrique Ouest"
   Disponible pour toutes les sessions futures
```

---

## 8. SÉCURITÉ — IMMUNOLOGIE SANS SANDBOX EXTERNE

### La posture immunologique

Prometheus ne "défend" pas. Il *scanne en continu* et *se soigne automatiquement*. C'est une différence fondamentale : un système de défense réagit aux attaques. Un système immunologique les anticipe et les neutralise avant qu'elles se produisent.

### Couche 1 — Intercepteur Pre-execution

```
TOUTE COMMANDE passe par l'intercepteur avant exec()

Score de risque calculé sur:
- Pattern matching (blacklist de 200+ patterns dangereux)
- Analyse contextuelle (la commande est-elle cohérente avec la tâche?)
- Historique (est-ce que des commandes similaires ont posé problème?)

Score 0-30: exec() direct
Score 31-70: log + exec() (audit silencieux)
Score 71-90: confirmation utilisateur requise
Score 91-100: refus + explication claire

Rate limiting: max 10 exec/seconde (protège contre fork bombs)
Timeout: toute commande > 5 minutes → kill + log
```

### Couche 2 — SAST sur code produit (aucun outil externe requis)

Prometheus implémente une analyse statique légère en Go, native, sans Semgrep ni outil externe :

```
À CHAQUE FICHIER DE CODE GÉNÉRÉ:

1. Analyse par pattern (Go natif, ~50ms):
   - SQL string concatenation → injection SQL
   - eval() / exec() sur input utilisateur → code injection
   - Credentials hardcodés (regex: password=, api_key=, token=)
   - HTTP sans TLS
   - Désactivation de vérification SSL

2. Analyse sémantique LLM (pour les cas ambigus, ~2s):
   "Ce code présente-t-il des risques de sécurité?
    Si oui, lesquels et comment les corriger?"

3. Auto-correction avant livraison:
   → Faille trouvée → LLM corrige → Re-analyse
   → Cycle jusqu'à clean

Si Semgrep est disponible (installé optionnellement):
   → Utilisé en complément pour les règles avancées
   → Mais jamais requis au démarrage
```

### Couche 3 — Scan environnement

```
AU DÉMARRAGE + HEBDOMADAIRE:

1. Ports ouverts: exec("ss -tlnp") ou exec("netstat -tlnp")
2. Services actifs: exec("systemctl list-units --type=service")
3. Permissions inhabituelles: exec("find / -perm -4000 2>/dev/null")
4. Variables d'environnement sensibles exposées
5. Fichiers de config avec credentials en clair

RÉSULTATS → analysés par LLM → rapport + corrections suggérées
```

### Couche 4 — DAST (si serveur web lancé)

```
QUAND Prometheus lance un serveur web:

1. Scan automatique (Go natif, sans ZAP):
   - Tester les headers HTTP manquants (HSTS, CSP, X-Frame-Options)
   - Tester endpoints sans auth
   - Tester upload de fichiers
   - Tester rate limiting absent

2. Si ZAP est disponible (installé optionnellement):
   → Scan complet

3. Corrections automatiques:
   → Headers manquants ajoutés
   → Endpoints sécurisés
   → Report final
```

---

## 9. TERMINAL UX — INTERFACE HUMAINE

### Principes UX

Le terminal Prometheus n'est PAS une interface développeur hostile. C'est une interface de conversation fluide qui révèle la puissance sous-jacente sans jamais l'imposer.

### Layout Terminal (bubbletea — Go natif)

```
┌─────────────────────────────────────────────────────────┐
│ PROMETHEUS  ●  phi3:mini  ●  offline  ●  FarMarket     │  ← Header (1 ligne)
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Toi  14:32                                             │
│  J'ai un projet — créer une app pour les agriculteurs   │
│  africains, vente directe sans intermédiaire            │
│                                                         │
│  Prometheus  14:32                                      │
│  Compris. Quelques questions rapides pour bien          │
│  démarrer :                                             │
│  • Plateforme : Android uniquement ou aussi web ?       │
│  • Connexion : les utilisateurs ont souvent la data ?   │
│                                                         │
│  Toi  14:33                                             │
│  Android priorité. Offline important, zones rurales.    │
│                                                         │
│  Prometheus  14:33                  [voir les logs ↓]  │
│  ████████████░░░░░  Création de la structure...  67%   │ ← Barre progression
│  ✓ mkdir farmarket/backend mobile docs                  │
│  ✓ git init + .gitignore                               │
│  ↻ npm init (en cours...)                               │
│                                                         │
├─────────────────────────────────────────────────────────┤
│ > _                                           [Ctrl+C]  │  ← Input (1 ligne)
└─────────────────────────────────────────────────────────┘
```

### Modes d'affichage

**Mode Conversation (défaut)**
Seuls les messages utilisateur/Prometheus sont visibles. Les commandes exec() tournent en arrière-plan. Une barre de progression discrète indique l'activité.

**Mode Transparent (Ctrl+L)**
Les logs de commandes s'affichent en temps réel. L'utilisateur voit exactement ce que Prometheus fait.

**Mode Silencieux (Ctrl+S)**
Prometheus travaille, notifie uniquement quand terminé ou bloqué.

### Gestion des blocages (la feature la plus importante)

```
Quand Prometheus est bloqué, il affiche ceci:

┌─────────────────────────────────────────────────────────┐
│ ⊙ PROMETHEUS A BESOIN D'UNE INFO                       │
│                                                         │
│ Je suis en train de cloner le repo privé               │
│ github.com/farmarket/auth-module                        │
│                                                         │
│ J'ai besoin de tes credentials GitHub.                  │
│                                                         │
│ Comment faire:                                          │
│ → Settings → Developer settings → Personal access      │
│   tokens → Generate new token (classic)                 │
│ → Coches: repo, read:org                               │
│                                                         │
│ Ton token: _________________________________            │
│                                                         │
│ [Entrer]  ou  [Passer cette étape]                     │
└─────────────────────────────────────────────────────────┘

Après ta réponse → Prometheus reprend exactement là où il était.
La tâche n'a pas été interrompue. Elle attendait.
```

### Confirmations pour actions critiques

```
┌─────────────────────────────────────────────────────────┐
│ ⚠ CONFIRMATION REQUISE                                 │
│                                                         │
│ Prometheus va exécuter:                                 │
│ > sudo apt install android-studio (1.2GB)              │
│                                                         │
│ Raison: Simulateur Android requis pour tester l'app    │
│ Durée estimée: ~8 minutes                               │
│ Espace requis: 2.4GB                                    │
│                                                         │
│ [Confirmer]  [Choisir autre chose]  [Annuler la tâche] │
└─────────────────────────────────────────────────────────┘
```

### Indicateurs visuels (ASCII, zéro dépendance)

```
⟳ En réflexion...
░░░████████████░  Exécution... 78%
✓ Terminé (4m 32s)
⊙ J'ai besoin d'une info
⚠ Confirmation requise
✗ Échec — voici ce qui s'est passé...

● Ollama local · offline         (vert = OK)
◌ Recherche modèle...            (jaune = en cours)
○ Aucun LLM disponible           (rouge = problème)
```

---

## 10. STACK TECHNIQUE FINALE

```toml
[runtime]
language = "Go 1.21+"
reason = """
  - Binaire unique 15-20MB, cross-compilation triviale
  - GOOS=android GOARCH=arm64 go build → ça marche
  - goroutines = parallélisme élégant sans complexité
  - stdlib complète (http, json, exec, sql)
  - Performances suffisantes (pas besoin de Rust pour ce cas)
  - Courbe d'apprentissage: 1 semaine vs 3 mois Rust
  - Communauté africaine Go > Rust (pragmatisme)
"""

[interface]
primary = "bubbletea"          # TUI native Go, zéro dépendance système
secondary = "http server Go"   # Web UI vanilla HTML/CSS/JS
mobile = "Android WebView"     # Wrapper natif sans Tauri

[llm_adapters]
unified_interface = "ModelProvider"  # Une interface, N implémentations
providers = [
  "ollama",           # Local, offline, défaut
  "anthropic",        # Claude API (optionnel)
  "google",           # Gemini API (optionnel)
  "gguf_direct",      # Chargement GGUF sans Ollama (llama.cpp embedded)
]

[storage]
tasks = "SQLite (tasks.db)"
memory = "SQLite (memory.db)"
logs = "JSON Lines + zstd compression"
capabilities = "TOML metadata + scripts"

[security]
sandbox = "Native Go (namespaces Linux) + graceful fallback"
sast = "Pattern matching Go natif + LLM pour cas ambigus"
dast = "Go natif + ZAP optionnel"
secrets = "Vault chiffré AES-256-GCM (Go stdlib)"

[parallelism]
model = "goroutines + channels"
max_workers = "runtime.NumCPU() * 2"
task_queue = "buffered channel"
ollama_parallel = "OLLAMA_NUM_PARALLEL=4 (auto-configuré)"

[configuration]
format = "TOML"
file = "~/.prometheus/prometheus.toml"
env_override = true  # Variables d'env surchargent le TOML

[dependencies_externes]
zero = true
note = """
  Prometheus démarre et fonctionne sans RIEN d'externe.
  Pas Python. Pas Node. Pas Docker. Pas Ollama préinstallé.
  Il peut tout installer lui-même si nécessaire.
  La seule dépendance est le binaire prometheus lui-même.
"""
```

---

## 11. PHASES DE DÉVELOPPEMENT

### Phase 0 — Proof of Concept (1 semaine)
**Goal:** Prouver que la boucle Think → Execute → Observe fonctionne.

**Livrables:**
- main.go — 400 lignes max
- LLM adapter Ollama uniquement
- Boucle de base sans persistance
- Test: "Crée un dossier projet/ avec un README.md"

**Critères de succès:**
- Binaire < 15MB
- Tâche simple réussie en < 30 secondes
- Aucun crash sur erreur shell

---

### Phase 1 — The Spark (3 semaines)
**Goal:** MVP fonctionnel, utilisable réellement.

**Livrables:**
- TUI bubbletea (conversation fluide)
- Task persistence SQLite
- Gestion Blocked state avec prompts clairs
- Retry logic sur erreurs
- LLM adapters: Ollama + Anthropic + Google
- Context Manager basique (Hot Buffer + compaction)
- Environment Discovery au démarrage
- Logs JSON Lines

**Critères de succès:**
- Binaire < 20MB
- Fonctionne sur Android Termux (testé)
- Task "Crée une API REST Python simple" réussit de bout en bout
- Reprise après interruption fonctionne

**Métriques:**
- RAM: < 200MB sans LLM
- Success rate tâches simples: > 80%
- Cold start: < 2 secondes

---

### Phase 2 — The Evolution (4 semaines)
**Goal:** Auto-installation de capacités + Context Manager complet.

**Livrables:**
- Capability Engine complet (Discovery → Install → Forge)
- Context Manager 3 niveaux (Hot/Warm/Cold)
- Compaction automatique adaptative
- Chrome DevTools Protocol (navigateur natif)
- Capability Forge (génération LLM si absent)
- Log rotation + compression zstd
- Résumés journaliers automatiques

**Critères de succès:**
- Installe git seul si manquant
- Scrape un site JavaScript sans API externe
- Fonctionne avec Phi-3 mini (4K context) aussi bien qu'avec Claude (200K)
- Logs compressés, résumés générés automatiquement

---

### Phase 3 — The Immune System (4 semaines)
**Goal:** Sécurité proactive, immunologie complète.

**Livrables:**
- Security Interceptor complet (score de risque, patterns, rate limit)
- SAST natif Go (pattern matching + LLM)
- Namespace isolation Linux (Niveau 2 sandbox)
- Scan environnement (ports, permissions, CVEs)
- DAST Go natif (headers, endpoints, auth)
- Vault credentials chiffré AES-256-GCM
- Auto-patching code vulnérable

**Critères de succès:**
- Détecte injection SQL dans son propre code et corrige
- Bloque fork bomb et commandes système dangereuses
- Vault credentials ne fuit pas en plain text

---

### Phase 4 — Observability & Performance (2 semaines)
**Goal:** Comprendre et optimiser.

**Livrables:**
- Structured logging complet
- Métriques internes (succès, latences, RAM)
- Requêtabilité sémantique des logs passés
- Profiling mémoire

**Critères de succès:**
- "Sur quoi avons-nous travaillé lundi?" → réponse correcte
- RAM profiling disponible
- Dashboard logs consultable

---

### Phase 5 — UX & Distribution (3 semaines)
**Goal:** Proximité maximale avec l'utilisateur.

**Livrables:**
- Web UI vanilla HTML/CSS/JS (serveur Go intégré)
- Android WebView wrapper (APK autonome)
- Installation one-liner (curl | sh — ironie assumée)
- Mode offline visuel clair
- Documentation utilisateur en français + anglais

**Critères de succès:**
- Un non-développeur peut démarrer Prometheus
- Web UI < 100KB non-compressé
- APK Android < 25MB

---

### Phase 6+ — Advanced (roadmap ouverte)
- Mesh P2P chiffré (commander son PC depuis son téléphone)
- Multi-agents orchestration avancée (plusieurs Prometheus coordonnés)
- Vision (analyse de screenshots pour debugger des bugs visuels)
- Voice interface (Whisper local)

---

## 12. SPÉCIFICATIONS NON-FONCTIONNELLES

### Performance

| Métrique | Target | Maximum |
|---|---|---|
| Binaire size | < 15MB | < 25MB |
| RAM sans LLM | < 150MB | < 300MB |
| RAM avec Phi-3 mini | < 4GB | < 6GB |
| Cold start | < 2s | < 5s |
| Think latency (Ollama local) | < 3s | < 10s |
| Think latency (API cloud) | < 1s | < 3s |
| Exec latency overhead | < 50ms | < 200ms |
| Log compressé par jour actif | ~500KB | ~2MB |
| Espace actif RAM (logs+mémoire) | < 20MB | < 50MB |

### Compatibilité

| Plateforme | Requis minimum | Priorité |
|---|---|---|
| Linux x86_64 | Kernel 3.10+ | P0 |
| Linux ARM64 | Kernel 4.0+ | P0 |
| Android ARM64 (Termux) | Android 10, 4GB RAM | P0 (Africa-first) |
| macOS ARM64 (Apple Silicon) | macOS 11+ | P0 |
| macOS x86_64 | macOS 10.15+ | P1 |
| Windows x86_64 | Windows 10+ | P1 |
| Raspberry Pi 4 (8GB) | Raspberry Pi OS 64-bit | P1 |

### Sécurité

| Principe | Implémentation |
|---|---|
| Zéro credentials en plain text | Vault AES-256-GCM Go stdlib |
| Least privilege | Jamais root par défaut |
| Rate limiting exec | Max 10/sec, configurable |
| Audit trail complet | Tous les exec() loggés |
| Sandbox natif | Namespaces Linux sans dépendance |
| SAST permanent | Chaque fichier généré scanné |

---

## 13. DISTRIBUTION & INSTALLATION

### Installation one-liner

```bash
# Linux / macOS / Android Termux
curl -L https://get.prometheus.dev | sh

# Ce script fait uniquement:
# 1. Détecter l'OS et l'architecture
# 2. Télécharger le bon binaire
# 3. Le placer dans ~/bin/prometheus
# 4. Ajouter ~/bin au PATH

# Démarrage immédiat
prometheus
```

### Configuration initiale

Au premier lancement, Prometheus crée sa structure et lance le scan d'environnement. Aucune configuration manuelle requise.

```
~/.prometheus/
  prometheus.toml    ← Configuration (générée automatiquement)
  tasks.db           ← État des tâches (SQLite)
  memory.db          ← Mémoire long terme (SQLite)
  capabilities/      ← Outils installés à la demande
  logs/
    2025-01-15.jsonl ← Log du jour (actif)
  summaries/
    2025-01-15.md    ← Résumé journalier (généré automatiquement)
  archive/
    2025-01/         ← Logs compressés du mois passé
  vault.enc          ← Credentials chiffrés AES-256-GCM
```

### Configuration minimale générée

```toml
# prometheus.toml (généré au premier lancement)
[llm]
provider = "ollama"          # Détecté automatiquement
model = "phi3:mini"          # Premier modèle trouvé
endpoint = "http://localhost:11434"

[security]
rate_limit_per_second = 10
dangerous_ops_confirmation = true

[memory]
compaction_threshold = 0.70  # Compacter à 70% de la context window
max_context_hot_ratio = 0.60
prune_summaries_after_days = 0  # Jamais (garder pour toujours)

[logs]
compress_after_days = 1
archive_after_days = 7
format = "jsonl"
compression = "zstd"

[capabilities]
auto_install = true
auto_forge_if_absent = true
```

---

## APPENDICE — Réponses aux questions ouvertes

**Q: Le terminal est-il user-friendly?**
Oui. bubbletea (Go natif, zéro dépendance) donne un terminal riche avec zones de dialogue distinctes, barres de progression, confirmation graphique, et séparation claire conversation/logs techniques. Ce n'est pas un terminal de développeur — c'est une interface de conversation qui tourne dans un terminal.

**Q: Comment on gère le contexte sur modèles légers?**
Context Manager adaptatif : il détecte la context window disponible, maintient un Hot Buffer à 60% de cette limite, et déclenche une compaction sémantique automatique à 70%. Le résumé compacté est ré-injecté comme contexte enrichi. Qualité constante quelle que soit la taille du modèle.

**Q: Le log peut-il couvrir plusieurs années?**
Oui. Compression zstd dès J+1 (gain 85%). Archivage mensuel. Résumés journaliers en markdown non-compressés (petits, rapides à lire). Projection : 5 ans d'utilisation active = ~850MB compressé total, < 20MB actif en RAM.

**Q: Que faire si une capability n'existe nulle part?**
Capability Forge : le LLM génère le script/outil nécessaire (Python, bash, Go), le teste automatiquement, et le stocke comme capability custom permanente. La prochaine fois, elle sera déjà disponible.

**Q: Sandbox sans Docker?**
Sandbox natif Go via syscalls Linux (namespaces + rlimits). Sur Android : fallback workdir isolé + timeout strict. Sur Windows : fallback audit pre-execution + rate limiting. Zéro dépendance externe pour la sécurité.

---

*Prometheus · v2.0 · Document de Référence*
*Conçu à Lomé · Pour l'Afrique · Utile partout*