# 🐉 PROMETHEUS : MASTER PLAN (Souveraineté & Immunologie)

## 🚩 Vision & Principes Fondamentaux
- **Noyau** : Micro-Kernel "Liquide" en **Rust** (Performance, Sécurité Mémoire).
- **Philosophie** : `exec()` est la primitive centrale. Tout pouvoir est acquis par l'installation et le scripting.
- **Posture** : Immunologique. Le système ne se défend pas, il se scanne et se soigne en temps réel.
- **Souveraineté** : 100% Local-first, Mesh P2P, Zéro VPS, Zero-Dependence externe.
- **UX** : Interface "Invisible" (Tauri/Native) $\rightarrow$ Chat simple $\rightarrow$ Puissance brute.

---

## 🗺️ Feuille de Route (Les 5 Phases)

### ⚡ Phase 1 : The Spark (Le Noyau)
*L'objectif est de passer du texte à l'action.*
- **Focus** : Binaire Rust $\rightarrow$ LLM $\rightarrow$ Shell $\rightarrow$ Result.
- **Livrables** :
    - `config.yaml` et `prompt.md` (le système nerveux).
    - Interface `ModelProvider` (Ollama, GGUF local).
    - Boucle `Think` $\rightarrow$ `Execute` $\rightarrow$ `Observe`.
    - Exécution de tâches simples (ex: créer un dossier et un fichier).

### 🧬 Phase 2 : The Evolution (Capacités Dynamiques)
*L'objectif est qu'il apprenne à s'outiller seul.*
- **Focus** : Gestionnaire de compétences et Web Native.
- **Livrables** :
    - Système de stockage `/capabilities` (Scripts + Métadonnées).
    - Intégration **Chrome DevTools Protocol (CDP)** pour le web sans API.
    - Boucle "Besoin $\rightarrow$ Recherche $\rightarrow$ Installation $\rightarrow$ Mémorisation".
    - Capacité d'installer des simulateurs (Android/iOS) et SDKs.

### 🛡️ Phase 3 : The Immune System (L'Immunologie)
*L'objectif est la sécurité proactive et l'auto-guérison.*
- **Focus** : Intercepteur de sécurité et scanneurs.
- **Livrables** :
    - Intercepteur de commandes (Pré-exécution et Post-exécution).
    - Intégration de scanners (Semgrep, OWASP ZAP) pour le code produit et l'environnement.
    - Algorithme de détection de vulnérabilités (0-day) via LLM.
    - Module d'auto-patching et de nettoyage des traces sensibles.

### 🌐 Phase 4 : The Mesh (Réseau Souverain P2P)
*L'objectif est de transformer Prometheus en un réseau d'agents.*
- **Focus** : Tunneling chiffré et relais d'intentions.
- **Livrables** :
    - Tunnel P2P chiffré (type Wireguard/Tailscale) entre instances.
    - Système de "Relais d'Intents" (commander le PC depuis le Mobile).
    - Gestion des clés de confiance et identités cryptographiques.
    - Gestion des files d'attente asynchrones (exécution différée si machine offline).

### 🎨 Phase 5 : The Interface & UX (Omnicanal)
*L'objectif est la proximité absolue avec l'utilisateur.*
- **Focus** : Application légère et flux de confirmation.
- **Livrables** :
    - Application **Tauri** (PC/Linux) et **Wrapper Natif** (Mobile).
    - Interface de Chat minimaliste avec mode "Réflexion" (logs transparents).
    - Flux de confirmation pour les actions critiques (Paiements, Root, Suppressions).
    - Mode "Offline Total" optimisé pour 4GB de RAM.

---

## ⚙️ Stack Technique Consolidée
- **Core** : Rust + Tokio (Async).
- **LLM** : Ollama / Local GGUF / API (via ModelProvider).
- **Web** : Chrome DevTools Protocol.
- **Interface** : Tauri (Desktop) / Native (Mobile).
- **Réseau** : Tunneling P2P chiffré.
- **Config** : YAML + Markdown.
