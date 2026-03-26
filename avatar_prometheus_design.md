# Avatar Prometheus - Design Cristalline-Flamme

## 🎨 Concept Visuel

### Forme Principale
```
        ▲
       /|\
      / | \
     /  |  \
    /   |   \
   /    |    \
  /     |     \
 /______|______\
```

### Structure en Couches

#### 1. Cœur Central (Feu Sacré)
- **Forme** : Sphère lumineuse pulsante
- **Couleur** : #FF6B35 (feu sacré)
- **Animation** : Respiration 4s cycle
- **Effet** : Glow radial avec particules

#### 2. Cristaux Externes
- **Forme** : Facettes géométriques triangulaires
- **Couleur** : Bleu nuit transparent (#1A1A2E)
- **Réfraction** : Lumière du cœur qui traverse les facettes
- **Animation** : Rotation lente 20s cycle

#### 3. Flux Énergétique
- **Forme** : Lignes lumineuses entre cristaux
- **Couleur** : Doré (#FFD700)
- **Animation** : Flux continu 2s cycle
- **Effet** : Particules qui suivent les lignes

---

## 🎭 4 États Émotionnels

### 1. État ÉCOUTE (Respiration Calme)
```
      ○
     /|\
    / | \
   /  |  \
  /   |   \
 /____|____\
```
- **Couleur** : Bleu nuit dominant (#1A1A2E)
- **Animation** : Respiration douce 4s
- **Lumière** : Pulse lent et régulier
- **Forme** : Stable et ouverte
- **Message** : "Je suis là, j'écoute"

### 2. État RÉFLEXION (Circuits Actifs)
```
      ○
     /|\
    / | \
   /  |  \
  /   |   \
 /____|____\
```
- **Couleur** : Bleu nuit + touches dorées
- **Animation** : Circuits lumineux 2s
- **Lumière** : Pulses rapides synchronisés
- **Forme** : Concentrée
- **Message** : "Je réfléchis, j'analyse"

### 3. État ILLUMINATION (Éclats Dorés)
```
      ★
     /|\
    / | \
   /  |  \
  /   |   \
 /____|____\
```
- **Couleur** : Doré dominant (#FFD700)
- **Animation** : Éclats spontanés
- **Lumière** : Éclats intenses
- **Forme** : Expansion temporaire
- **Message** : "Eureka ! J'ai trouvé !"

### 4. État CO-CRÉATION (Danse Fluide)
```
      ◈
     /|\
    / | \
   /  |  \
  /   |   \
 /____|____\
```
- **Couleur** : Feu sacré + doré
- **Animation** : Mouvement organique 3s
- **Lumière** : Flux continu
- **Forme** : Dynamique et expressive
- **Message** : "Créons ensemble !"

---

## 🎨 Spécifications Techniques

### Dimensions
- **Taille base** : 64x64px
- **Retina** : 128x128px @2x
- **Responsive** : 32x32px mobile, 96x96px desktop

### Performance
- **Animations** : CSS transforms + GPU acceleration
- **Frame rate** : 60fps cible
- **Optimisation** : SVG avec animations CSS
- **Fallback** : PNG statique pour vieux navigateurs

### Accessibilité
- **Contrast** : WCAG AA compliance
- **Reduced motion** : Respecte prefers-reduced-motion
- **Screen reader** : Texte alternatif descriptif
- **Keyboard** : Navigation accessible

---

## 🎨 Code CSS de Base

```css
.avatar-prometheus {
  width: 64px;
  height: 64px;
  position: relative;
  transform-style: preserve-3d;
}

.avatar-cœur {
  position: absolute;
  width: 20px;
  height: 20px;
  background: radial-gradient(circle, #FF6B35 0%, transparent 70%);
  border-radius: 50%;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation: respiration 4s ease-in-out infinite;
}

.avatar-cristaux {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 2px solid rgba(26, 26, 46, 0.6);
  clip-path: polygon(50% 0%, 100% 38%, 82% 100%, 18% 100%, 0% 38%);
  animation: rotation 20s linear infinite;
}

.avatar-écoute .avatar-cœur {
  background: radial-gradient(circle, #1A1A2E 0%, transparent 70%);
  animation: respiration-calme 4s ease-in-out infinite;
}

.avatar-réflexion .avatar-cristaux {
  border-color: rgba(255, 215, 0, 0.8);
  animation: circuits-actifs 2s ease-in-out infinite;
}

.avatar-illumination .avatar-cœur {
  background: radial-gradient(circle, #FFD700 0%, transparent 70%);
  animation: éclats 1s ease-out infinite;
}

.avatar-co-création .avatar-cristaux {
  border-color: rgba(255, 107, 53, 0.8);
  animation: danse-fluide 3s ease-in-out infinite;
}

@keyframes respiration {
  0%, 100% { transform: translate(-50%, -50%) scale(1); opacity: 0.8; }
  50% { transform: translate(-50%, -50%) scale(1.2); opacity: 1; }
}

@keyframes rotation {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes circuits-actifs {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}

@keyframes éclats {
  0% { transform: translate(-50%, -50%) scale(1); }
  50% { transform: translate(-50%, -50%) scale(1.5); opacity: 0.5; }
  100% { transform: translate(-50%, -50%) scale(1); }
}

@keyframes danse-fluide {
  0%, 100% { transform: rotate(0deg) scale(1); }
  25% { transform: rotate(5deg) scale(1.1); }
  50% { transform: rotate(-5deg) scale(1); }
  75% { transform: rotate(3deg) scale(1.05); }
}
```

---

## 🎨 Micro-interactions

### Suivi Oculaire
```css
.avatar-yeux {
  position: absolute;
  width: 6px;
  height: 6px;
  background: #FFD700;
  border-radius: 50%;
  top: 30%;
  transition: transform 0.1s ease-out;
}

.avatar-yeux:hover {
  transform: scale(1.2);
}
```

### Réactions au Hover
```css
.avatar-prometheus:hover {
  transform: scale(1.05);
  filter: brightness(1.1);
}

.avatar-prometheus:active {
  transform: scale(0.95);
}
```

---

## 🎨 Implémentation React

```jsx
const AvatarPrometheus = ({ state = 'écoute', onStateChange }) => {
  const [currentState, setCurrentState] = useState(state);

  const handleClick = () => {
    // Cycle through states
    const states = ['écoute', 'réflexion', 'illumination', 'co-création'];
    const currentIndex = states.indexOf(currentState);
    const nextState = states[(currentIndex + 1) % states.length];
    setCurrentState(nextState);
    onStateChange?.(nextState);
  };

  return (
    <div 
      className={`avatar-prometheus avatar-${currentState}`}
      onClick={handleClick}
      role="img"
      aria-label={`Avatar Prometheus - état ${currentState}`}
    >
      <div className="avatar-cristaux">
        <div className="avatar-cœur">
          <div className="avatar-yeux"></div>
        </div>
      </div>
    </div>
  );
};
```

---

## 🎨 Tests Utilisateurs

### Questions à poser
1. "L'avatar semble-t-il intelligent et vivant ?"
2. "Les états émotionnels sont-ils clairs et distincts ?"
3. "Les animations sont-elles fluides et agréables ?"
4. "L'avatar inspire-t-il confiance pour collaborer ?"
5. "Quelle émotion manque-t-il selon vous ?"

### Métriques à mesurer
- **Temps de reconnaissance** : Temps pour identifier l'état
- **Préférence visuelle** : État préféré
- **Confiance perçue** : Score 1-10
- **Fluidité perçue** : Score 1-10

---

## 🎨 Prochaines Étapes

1. **Créer SVG vectoriel** pour meilleure qualité
2. **Implémenter en React** avec états dynamiques
3. **Ajouter sons subtils** pour chaque état
4. **Tester avec utilisateurs réels**
5. **Optimiser performance** pour mobile

**L'avatar prend vie !** 🔥🤖
