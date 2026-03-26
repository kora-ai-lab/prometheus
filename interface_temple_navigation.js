// Navigation Intelligente pour Interface Temple
class TempleNavigation {
  constructor() {
    this.currentZone = 'conversation';
    this.zones = ['conversation', 'mondes', 'connaissance', 'inspiration'];
    this.init();
  }

  init() {
    this.createNavigation();
    this.setupKeyboardShortcuts();
    this.setupZoneTransitions();
    this.setupAutoFocus();
    this.createBreadcrumbs();
  }

  createNavigation() {
    // Navigation flottante
    const nav = document.createElement('div');
    nav.className = 'temple-nav';
    nav.innerHTML = `
      <button class="nav-button active" data-zone="conversation" title="Conversation">
        💬
      </button>
      <button class="nav-button" data-zone="mondes" title="Mondes">
        🌍
      </button>
      <button class="nav-button" data-zone="connaissance" title="Connaissance">
        🧠
      </button>
      <button class="nav-button" data-zone="inspiration" title="Inspiration">
        💡
      </button>
    `;
    document.body.appendChild(nav);

    // Event listeners
    nav.querySelectorAll('.nav-button').forEach(button => {
      button.addEventListener('click', (e) => {
        const zone = e.target.dataset.zone;
        this.navigateToZone(zone);
      });
    });
  }

  navigateToZone(zone) {
    if (zone === this.currentZone) return;

    // Transition overlay
    this.showTransition();

    setTimeout(() => {
      // Mettre à jour les boutons
      document.querySelectorAll('.nav-button').forEach(btn => {
        btn.classList.toggle('active', btn.dataset.zone === zone);
      });

      // Focus sur la zone
      this.focusZone(zone);
      
      // Mettre à jour breadcrumbs
      this.updateBreadcrumbs(zone);

      // Cacher transition
      this.hideTransition();

      this.currentZone = zone;
    }, 300);
  }

  focusZone(zone) {
    const zoneElement = document.querySelector(`.zone-${zone}`);
    if (zoneElement) {
      // Scroll smooth vers la zone
      zoneElement.scrollIntoView({ 
        behavior: 'smooth', 
        block: 'center' 
      });

      // Highlight temporaire
      zoneElement.style.transition = 'all 0.3s ease';
      zoneElement.style.boxShadow = '0 0 30px rgba(255, 215, 0, 0.5)';
      zoneElement.style.transform = 'scale(1.02)';

      setTimeout(() => {
        zoneElement.style.boxShadow = '';
        zoneElement.style.transform = '';
      }, 1000);
    }
  }

  setupKeyboardShortcuts() {
    document.addEventListener('keydown', (e) => {
      // Alt + 1-4 pour naviguer
      if (e.altKey && e.key >= '1' && e.key <= '4') {
        const index = parseInt(e.key) - 1;
        this.navigateToZone(this.zones[index]);
        e.preventDefault();
      }

      // Tab pour cycle entre zones
      if (e.key === 'Tab' && !e.shiftKey && e.altKey) {
        const currentIndex = this.zones.indexOf(this.currentZone);
        const nextIndex = (currentIndex + 1) % this.zones.length;
        this.navigateToZone(this.zones[nextIndex]);
        e.preventDefault();
      }

      // Shift + Tab pour cycle inverse
      if (e.key === 'Tab' && e.shiftKey && e.altKey) {
        const currentIndex = this.zones.indexOf(this.currentZone);
        const prevIndex = (currentIndex - 1 + this.zones.length) % this.zones.length;
        this.navigateToZone(this.zones[prevIndex]);
        e.preventDefault();
      }

      // Ctrl + / pour afficher l'aide
      if (e.ctrlKey && e.key === '/') {
        this.showHelp();
        e.preventDefault();
      }
    });
  }

  setupZoneTransitions() {
    // Transitions entre zones
    this.zones.forEach(zone => {
      const zoneElement = document.querySelector(`.zone-${zone}`);
      if (zoneElement) {
        zoneElement.addEventListener('mouseenter', () => {
          this.previewZone(zone);
        });

        zoneElement.addEventListener('click', (e) => {
          if (!e.target.closest('input, button, .nav-button')) {
            this.navigateToZone(zone);
          }
        });
      }
    });
  }

  previewZone(zone) {
    const preview = document.createElement('div');
    preview.className = 'zone-preview';
    preview.innerHTML = `
      <div class="preview-content">
        <h3>${this.getZoneTitle(zone)}</h3>
        <p>${this.getZoneDescription(zone)}</p>
      </div>
    `;
    
    document.body.appendChild(preview);

    // Positionner près de la souris
    document.addEventListener('mousemove', (e) => {
      preview.style.left = e.clientX + 10 + 'px';
      preview.style.top = e.clientY + 10 + 'px';
    });

    // Retirer après 2 secondes
    setTimeout(() => preview.remove(), 2000);
  }

  getZoneTitle(zone) {
    const titles = {
      conversation: '💬 Conversation',
      mondes: '🌍 Mondes Actifs',
      connaissance: '🧠 Connaissance',
      inspiration: '💡 Inspiration'
    };
    return titles[zone] || zone;
  }

  getZoneDescription(zone) {
    const descriptions = {
      conversation: 'Dialogue avec Prometheus pour co-créer',
      mondes: 'Vos projets et écosystèmes connectés',
      connaissance: 'Apprentissages et évolution',
      inspiration: 'Idées et flux créatif'
    };
    return descriptions[zone] || '';
  }

  setupAutoFocus() {
    // Auto-focus intelligent selon l'activité
    const chatInput = document.querySelector('.chat-input');
    const mondesGallery = document.querySelector('.mondes-gallery');

    // Si l'utilisateur tape dans le chat, focus conversation
    if (chatInput) {
      chatInput.addEventListener('focus', () => {
        this.navigateToZone('conversation');
      });
    }

    // Si l'utilisateur interagit avec les mondes, focus mondes
    if (mondesGallery) {
      mondesGallery.addEventListener('click', (e) => {
        if (e.target.closest('.monde-orb')) {
          this.navigateToZone('mondes');
        }
      });
    }
  }

  createBreadcrumbs() {
    const breadcrumbs = document.createElement('div');
    breadcrumbs.className = 'temple-breadcrumbs';
    breadcrumbs.innerHTML = `
      <span class="breadcrumb-item active" data-zone="conversation">Conversation</span>
      <span class="breadcrumb-separator">›</span>
      <span class="breadcrumb-item" data-zone="mondes">Mondes</span>
      <span class="breadcrumb-separator">›</span>
      <span class="breadcrumb-item" data-zone="connaissance">Connaissance</span>
      <span class="breadcrumb-separator">›</span>
      <span class="breadcrumb-item" data-zone="inspiration">Inspiration</span>
    `;
    
    document.querySelector('.temple-creation').insertBefore(
      breadcrumbs, 
      document.querySelector('.temple-creation').firstChild
    );

    // Event listeners
    breadcrumbs.querySelectorAll('.breadcrumb-item').forEach(item => {
      item.addEventListener('click', (e) => {
        this.navigateToZone(e.target.dataset.zone);
      });
    });
  }

  updateBreadcrumbs(zone) {
    document.querySelectorAll('.breadcrumb-item').forEach(item => {
      item.classList.toggle('active', item.dataset.zone === zone);
    });
  }

  showTransition() {
    const overlay = document.createElement('div');
    overlay.className = 'transition-overlay active';
    document.body.appendChild(overlay);
  }

  hideTransition() {
    const overlay = document.querySelector('.transition-overlay');
    if (overlay) {
      overlay.classList.remove('active');
      setTimeout(() => overlay.remove(), 300);
    }
  }

  showHelp() {
    const help = document.createElement('div');
    help.className = 'navigation-help';
    help.innerHTML = `
      <div class="help-content">
        <h3>🎯 Navigation Rapide</h3>
        <ul>
          <li><kbd>Alt + 1</kbd> → Conversation</li>
          <li><kbd>Alt + 2</kbd> → Mondes</li>
          <li><kbd>Alt + 3</kbd> → Connaissance</li>
          <li><kbd>Alt + 4</kbd> → Inspiration</li>
          <li><kbd>Alt + Tab</kbd> → Zone suivante</li>
          <li><kbd>Alt + Shift + Tab</kbd> → Zone précédente</li>
          <li><kbd>Ctrl + /</kbd> → Aide</li>
        </ul>
        <p>Cliquez sur une zone pour la focaliser</p>
        <button class="help-close">Fermer</button>
      </div>
    `;
    
    document.body.appendChild(help);

    help.querySelector('.help-close').addEventListener('click', () => {
      help.remove();
    });

    // Fermer automatiquement après 10 secondes
    setTimeout(() => help.remove(), 10000);
  }

  // Navigation contextuelle
  navigateToContext(context) {
    switch(context) {
      case 'chat':
      case 'conversation':
        this.navigateToZone('conversation');
        // Focus input
        setTimeout(() => {
          const input = document.querySelector('.chat-input');
          if (input) input.focus();
        }, 500);
        break;
      
      case 'world':
      case 'monde':
      case 'project':
        this.navigateToZone('mondes');
        break;
      
      case 'knowledge':
      case 'learning':
      case 'stats':
        this.navigateToZone('connaissance');
        break;
      
      case 'inspiration':
      case 'ideas':
      case 'creative':
        this.navigateToZone('inspiration');
        break;
    }
  }

  // Smart navigation basée sur l'activité
  smartNavigation(activity) {
    const navigationMap = {
      'typing': 'conversation',
      'creating_world': 'mondes',
      'viewing_stats': 'connaissance',
      'browsing_ideas': 'inspiration',
      'coding': 'conversation', // Retourner au chat pour demander de l'aide
      'debugging': 'conversation'
    };

    const targetZone = navigationMap[activity];
    if (targetZone && targetZone !== this.currentZone) {
      this.navigateToZone(targetZone);
    }
  }
}

// CSS pour la navigation
const navigationCSS = `
.temple-breadcrumbs {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 12px;
}

.breadcrumb-item {
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  transition: color 0.2s ease;
}

.breadcrumb-item:hover,
.breadcrumb-item.active {
  color: #FFD700;
}

.breadcrumb-separator {
  color: rgba(255, 255, 255, 0.4);
}

.zone-preview {
  position: fixed;
  background: rgba(26, 26, 46, 0.95);
  border: 1px solid rgba(255, 215, 0, 0.3);
  border-radius: 8px;
  padding: 12px;
  font-size: 12px;
  z-index: 1000;
  pointer-events: none;
  max-width: 200px;
}

.preview-content h3 {
  margin: 0 0 4px 0;
  color: #FFD700;
  font-size: 14px;
}

.preview-content p {
  margin: 0;
  color: rgba(255, 255, 255, 0.8);
}

.navigation-help {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(26, 26, 46, 0.95);
  border: 1px solid rgba(255, 215, 0, 0.3);
  border-radius: 12px;
  padding: 24px;
  z-index: 2000;
  max-width: 400px;
}

.help-content h3 {
  margin: 0 0 16px 0;
  color: #FFD700;
}

.help-content ul {
  margin: 16px 0;
  padding: 0;
  list-style: none;
}

.help-content li {
  margin: 8px 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

kbd {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 4px;
  padding: 2px 6px;
  font-family: monospace;
  font-size: 11px;
}

.help-close {
  background: linear-gradient(135deg, #FF6B35, #FF8C42);
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  color: white;
  cursor: pointer;
  margin-top: 16px;
}
`;

// Injecter CSS
const style = document.createElement('style');
style.textContent = navigationCSS;
document.head.appendChild(style);

// Initialiser la navigation
document.addEventListener('DOMContentLoaded', () => {
  new TempleNavigation();
});

export default TempleNavigation;
