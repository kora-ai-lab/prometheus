// Transitions Fluides pour Interface Temple
class TempleTransitions {
  constructor() {
    this.transitionDuration = 300;
    this.easing = 'cubic-bezier(0.4, 0, 0.2, 1)';
    this.init();
  }

  init() {
    this.setupZoneTransitions();
    this.setupStateTransitions();
    this.setupMicroInteractions();
    this.setupLoadingStates();
    this.setupErrorTransitions();
  }

  setupZoneTransitions() {
    // Transitions entre zones avec overlay
    this.zones = document.querySelectorAll('.zone-conversation, .zone-mondes, .zone-connaissance, .zone-inspiration');
    
    this.zones.forEach(zone => {
      zone.addEventListener('mouseenter', (e) => this.handleZoneHover(e, true));
      zone.addEventListener('mouseleave', (e) => this.handleZoneHover(e, false));
      zone.addEventListener('click', (e) => this.handleZoneClick(e));
    });
  }

  handleZoneHover(event, isEntering) {
    const zone = event.currentTarget;
    
    if (isEntering) {
      // Entrée dans la zone
      zone.style.transition = `all ${this.transitionDuration}ms ${this.easing}`;
      zone.style.transform = 'translateY(-2px)';
      zone.style.boxShadow = '0 8px 25px rgba(255, 215, 0, 0.3)';
      zone.style.borderColor = '#FFD700';
      
      // Animation du contenu
      this.animateContent(zone, 'enter');
    } else {
      // Sortie de la zone
      zone.style.transform = '';
      zone.style.boxShadow = '';
      zone.style.borderColor = '';
      
      // Animation du contenu
      this.animateContent(zone, 'leave');
    }
  }

  handleZoneClick(event) {
    const zone = event.currentTarget;
    
    // Ne pas transitionner si clic sur élément interactif
    if (event.target.closest('input, button, a, .nav-button')) return;
    
    // Transition de focus
    this.showTransitionOverlay(() => {
      this.focusZone(zone);
    });
  }

  focusZone(zone) {
    // Scroll smooth vers la zone
    zone.scrollIntoView({ 
      behavior: 'smooth', 
      block: 'center' 
    });

    // Highlight animation
    zone.style.transition = `all ${this.transitionDuration}ms ${this.easing}`;
    zone.style.transform = 'scale(1.02)';
    zone.style.boxShadow = '0 0 30px rgba(255, 215, 0, 0.5)';
    zone.style.borderColor = '#FFD700';

    // Pulse effect
    this.createPulseEffect(zone);

    setTimeout(() => {
      zone.style.transform = '';
      zone.style.boxShadow = '';
      zone.style.borderColor = '';
    }, 1000);
  }

  animateContent(zone, direction) {
    const elements = zone.querySelectorAll('.zone-header, .zone-content > *');
    
    elements.forEach((element, index) => {
      element.style.transition = `all ${this.transitionDuration}ms ${this.easing}`;
      element.style.transitionDelay = `${index * 50}ms`;
      
      if (direction === 'enter') {
        element.style.transform = 'translateX(5px)';
        element.style.opacity = '0.8';
      } else {
        element.style.transform = '';
        element.style.opacity = '';
      }
      
      setTimeout(() => {
        element.style.transform = '';
        element.style.opacity = '';
        element.style.transitionDelay = '';
      }, this.transitionDuration + (index * 50));
    });
  }

  setupStateTransitions() {
    // Transitions d'état de l'avatar
    const avatar = document.querySelector('.avatar-prometheus');
    if (avatar) {
      this.observeAvatarStates(avatar);
    }

    // Transitions de messages
    this.setupMessageTransitions();
    
    // Transitions de mondes
    this.setupWorldTransitions();
  }

  observeAvatarStates(avatar) {
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.type === 'attributes' && mutation.attributeName === 'class') {
          this.handleAvatarStateChange(avatar);
        }
      });
    });

    observer.observe(avatar, {
      attributes: true,
      attributeFilter: ['class']
    });
  }

  handleAvatarStateChange(avatar) {
    const currentClass = avatar.className;
    const states = ['écoute', 'réflexion', 'illumination', 'co-création'];
    const currentState = states.find(state => currentClass.includes(state));
    
    if (currentState) {
      this.showStateTransition(currentState);
    }
  }

  showStateTransition(state) {
    const transition = document.createElement('div');
    transition.className = 'state-transition';
    transition.innerHTML = `
      <div class="state-icon">${this.getStateIcon(state)}</div>
      <div class="state-name">${this.getStateName(state)}</div>
    `;
    
    document.body.appendChild(transition);

    setTimeout(() => {
      transition.classList.add('active');
    }, 10);

    setTimeout(() => {
      transition.classList.remove('active');
      setTimeout(() => transition.remove(), 300);
    }, 2000);
  }

  getStateIcon(state) {
    const icons = {
      'écoute': '👂',
      'réflexion': '🤔',
      'illumination': '💡',
      'co-création': '🤝'
    };
    return icons[state] || '🤖';
  }

  getStateName(state) {
    const names = {
      'écoute': 'Écoute',
      'réflexion': 'Réflexion',
      'illumination': 'Illumination',
      'co-création': 'Co-création'
    };
    return names[state] || state;
  }

  setupMessageTransitions() {
    const messagesContainer = document.querySelector('.chat-messages');
    if (!messagesContainer) return;

    // Observer les nouveaux messages
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        mutation.addedNodes.forEach((node) => {
          if (node.classList && node.classList.contains('message')) {
            this.animateMessageEntry(node);
          }
        });
      });
    });

    observer.observe(messagesContainer, {
      childList: true
    });
  }

  animateMessageEntry(message) {
    message.style.opacity = '0';
    message.style.transform = 'translateY(20px)';
    
    setTimeout(() => {
      message.style.transition = `all ${this.transitionDuration}ms ${this.easing}`;
      message.style.opacity = '1';
      message.style.transform = 'translateY(0)';
    }, 10);
  }

  setupWorldTransitions() {
    const worldOrbs = document.querySelectorAll('.monde-orb');
    
    worldOrbs.forEach(orb => {
      orb.addEventListener('click', (e) => this.handleWorldClick(e, orb));
      orb.addEventListener('mouseenter', (e) => this.handleWorldHover(e, orb, true));
      orb.addEventListener('mouseleave', (e) => this.handleWorldHover(e, orb, false));
    });
  }

  handleWorldClick(event, orb) {
    // Créer l'effet d'expansion
    this.createWorldExpansion(orb);
    
    // Transition vers la vue détaillée du monde
    setTimeout(() => {
      this.showTransitionOverlay(() => {
        this.navigateToWorld(orb);
      });
    }, 200);
  }

  handleWorldHover(event, orb, isEntering) {
    if (isEntering) {
      orb.style.transition = `all ${this.transitionDuration}ms ${this.easing}`;
      orb.style.transform = 'scale(1.05) rotate(5deg)';
      orb.style.boxShadow = '0 0 20px rgba(255, 215, 0, 0.6)';
      
      // Créer des particules
      this.createWorldParticles(orb);
    } else {
      orb.style.transform = '';
      orb.style.boxShadow = '';
    }
  }

  createWorldExpansion(orb) {
    const expansion = document.createElement('div');
    expansion.className = 'world-expansion';
    
    const rect = orb.getBoundingClientRect();
    expansion.style.left = rect.left + rect.width / 2 + 'px';
    expansion.style.top = rect.top + rect.height / 2 + 'px';
    
    document.body.appendChild(expansion);

    setTimeout(() => {
      expansion.classList.add('active');
    }, 10);

    setTimeout(() => {
      expansion.classList.remove('active');
      setTimeout(() => expansion.remove(), 300);
    }, 600);
  }

  createWorldParticles(orb) {
    const rect = orb.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;

    for (let i = 0; i < 8; i++) {
      const particle = document.createElement('div');
      particle.className = 'world-particle';
      
      const angle = (Math.PI * 2 * i) / 8;
      const distance = 30 + Math.random() * 20;
      
      particle.style.left = centerX + 'px';
      particle.style.top = centerY + 'px';
      particle.style.setProperty('--tx', `${Math.cos(angle) * distance}px`);
      particle.style.setProperty('--ty', `${Math.sin(angle) * distance}px`);
      
      document.body.appendChild(particle);
      
      setTimeout(() => particle.remove(), 1000);
    }
  }

  navigateToWorld(orb) {
    // Simulation de navigation vers le monde
    const worldName = orb.querySelector('.monde-name')?.textContent || 'Monde';
    console.log(`Navigation vers ${worldName}`);
    
    // Ici, on pourrait charger la vue détaillée du monde
  }

  setupMicroInteractions() {
    // Boutons avec ripple effect
    const buttons = document.querySelectorAll('button, .nav-button, .chat-send');
    buttons.forEach(button => {
      button.addEventListener('click', (e) => this.createRippleEffect(e, button));
    });

    // Cards avec hover effect
    const cards = document.querySelectorAll('.connaissance-card, .inspiration-card');
    cards.forEach(card => {
      card.addEventListener('mouseenter', (e) => this.handleCardHover(e, card, true));
      card.addEventListener('mouseleave', (e) => this.handleCardHover(e, card, false));
    });

    // Input avec focus effect
    const inputs = document.querySelectorAll('input, textarea');
    inputs.forEach(input => {
      input.addEventListener('focus', (e) => this.handleInputFocus(e, input, true));
      input.addEventListener('blur', (e) => this.handleInputFocus(e, input, false));
    });
  }

  createRippleEffect(event, button) {
    const ripple = document.createElement('div');
    ripple.className = 'ripple';
    
    const rect = button.getBoundingClientRect();
    const size = Math.max(rect.width, rect.height);
    const x = event.clientX - rect.left - size / 2;
    const y = event.clientY - rect.top - size / 2;
    
    ripple.style.width = ripple.style.height = size + 'px';
    ripple.style.left = x + 'px';
    ripple.style.top = y + 'px';
    
    button.appendChild(ripple);

    setTimeout(() => ripple.classList.add('active'), 10);
    setTimeout(() => ripple.remove(), 600);
  }

  handleCardHover(event, card, isEntering) {
    if (isEntering) {
      card.style.transition = `all ${this.transitionDuration}ms ${this.easing}`;
      card.style.transform = 'translateY(-4px)';
      card.style.boxShadow = '0 8px 25px rgba(0, 0, 0, 0.15)';
    } else {
      card.style.transform = '';
      card.style.boxShadow = '';
    }
  }

  handleInputFocus(event, input, isFocused) {
    if (isFocused) {
      input.style.transition = `all ${this.transitionDuration / 2}ms ${this.easing}`;
      input.style.transform = 'scale(1.02)';
      input.style.boxShadow = '0 0 0 2px rgba(255, 215, 0, 0.5)';
    } else {
      input.style.transform = '';
      input.style.boxShadow = '';
    }
  }

  setupLoadingStates() {
    // Loading skeleton screens
    this.createLoadingSkeletons();
    
    // Progress animations
    this.setupProgressAnimations();
  }

  createLoadingSkeletons() {
    const skeletonHTML = `
      <div class="skeleton-loader">
        <div class="skeleton-line"></div>
        <div class="skeleton-line short"></div>
        <div class="skeleton-line"></div>
      </div>
    `;

    // Ajouter des skeletons pour le chargement
    this.showLoadingState = (container) => {
      const skeleton = document.createElement('div');
      skeleton.innerHTML = skeletonHTML;
      container.appendChild(skeleton.firstElementChild);
    };

    this.hideLoadingState = (container) => {
      const skeleton = container.querySelector('.skeleton-loader');
      if (skeleton) skeleton.remove();
    };
  }

  setupProgressAnimations() {
    const progressBars = document.querySelectorAll('.connaissance-progress-bar');
    
    progressBars.forEach(bar => {
      const width = bar.style.width || '0%';
      bar.style.width = '0%';
      
      setTimeout(() => {
        bar.style.transition = `width 1s ${this.easing}`;
        bar.style.width = width;
      }, 100);
    });
  }

  setupErrorTransitions() {
    // Error shake animation
    this.showError = (element) => {
      element.style.animation = 'shake 0.5s ease-in-out';
      setTimeout(() => {
        element.style.animation = '';
      }, 500);
    };

    // Success pulse animation
    this.showSuccess = (element) => {
      element.style.animation = 'pulse-success 0.6s ease-in-out';
      setTimeout(() => {
        element.style.animation = '';
      }, 600);
    };
  }

  showTransitionOverlay(callback) {
    const overlay = document.createElement('div');
    overlay.className = 'transition-overlay';
    document.body.appendChild(overlay);

    setTimeout(() => {
      overlay.classList.add('active');
    }, 10);

    setTimeout(() => {
      if (callback) callback();
      
      setTimeout(() => {
        overlay.classList.remove('active');
        setTimeout(() => overlay.remove(), 300);
      }, 100);
    }, this.transitionDuration);
  }

  createPulseEffect(element) {
    const pulse = document.createElement('div');
    pulse.className = 'pulse-effect';
    
    const rect = element.getBoundingClientRect();
    pulse.style.left = rect.left + rect.width / 2 + 'px';
    pulse.style.top = rect.top + rect.height / 2 + 'px';
    
    document.body.appendChild(pulse);

    setTimeout(() => pulse.classList.add('active'), 10);
    setTimeout(() => pulse.remove(), 1000);
  }
}

// CSS pour les transitions
const transitionsCSS = `
.state-transition {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(26, 26, 46, 0.95);
  border: 2px solid #FFD700;
  border-radius: 12px;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  gap: 12px;
  z-index: 2000;
  opacity: 0;
  transform: translate(-50%, -50%) scale(0.8);
  transition: all 0.3s ease;
}

.state-transition.active {
  opacity: 1;
  transform: translate(-50%, -50%) scale(1);
}

.state-icon {
  font-size: 24px;
}

.state-name {
  font-size: 16px;
  font-weight: 600;
  color: #FFD700;
}

.world-expansion {
  position: fixed;
  width: 20px;
  height: 20px;
  border: 2px solid #FFD700;
  border-radius: 50%;
  transform: translate(-50%, -50%) scale(0);
  opacity: 1;
  transition: all 0.6s ease;
  z-index: 1000;
}

.world-expansion.active {
  transform: translate(-50%, -50%) scale(10);
  opacity: 0;
}

.world-particle {
  position: fixed;
  width: 4px;
  height: 4px;
  background: #FFD700;
  border-radius: 50%;
  transform: translate(-50%, -50%);
  opacity: 0;
  transition: all 1s ease;
  z-index: 999;
}

.world-particle:nth-child(odd) {
  background: #FF6B35;
}

.world-particle.active {
  opacity: 1;
  transform: translate(calc(-50% + var(--tx)), calc(-50% + var(--ty)));
}

.ripple {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  transform: scale(0);
  animation: ripple-animation 0.6s ease-out;
  pointer-events: none;
}

@keyframes ripple-animation {
  to {
    transform: scale(4);
    opacity: 0;
  }
}

.skeleton-loader {
  padding: 16px;
}

.skeleton-line {
  height: 12px;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.1));
  background-size: 200% 100%;
  animation: loading 1.5s ease-in-out infinite;
  border-radius: 4px;
  margin-bottom: 8px;
}

.skeleton-line.short {
  width: 60%;
}

@keyframes loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-2px); }
  20%, 40%, 60%, 80% { transform: translateX(2px); }
}

@keyframes pulse-success {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.pulse-effect {
  position: fixed;
  width: 100px;
  height: 100px;
  border: 2px solid #FFD700;
  border-radius: 50%;
  transform: translate(-50%, -50%) scale(0);
  opacity: 1;
  animation: pulse-animation 1s ease-out;
  pointer-events: none;
  z-index: 998;
}

@keyframes pulse-animation {
  0% {
    transform: translate(-50%, -50%) scale(0);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) scale(3);
    opacity: 0;
  }
}
`;

// Injecter CSS
const style = document.createElement('style');
style.textContent = transitionsCSS;
document.head.appendChild(style);

// Initialiser les transitions
document.addEventListener('DOMContentLoaded', () => {
  new TempleTransitions();
});

export default TempleTransitions;
