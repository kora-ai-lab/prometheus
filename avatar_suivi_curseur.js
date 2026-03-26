// Suivi oculaire intelligent pour Avatar Prometheus
class AvatarEyeTracking {
  constructor(avatarElement) {
    this.avatar = avatarElement;
    this.œilGauche = avatarElement.querySelector('.œil-gauche');
    this.œilDroit = avatarElement.querySelector('.œil-droit');
    this.œilYeux = avatarElement.querySelector('.avatar-yeux');
    
    this.init();
  }
  
  init() {
    // Suivi du curseur
    document.addEventListener('mousemove', (e) => this.trackCursor(e));
    
    // Clignements naturels
    this.startBlinking();
    
    // Réactions aux interactions
    this.setupReactions();
  }
  
  trackCursor(event) {
    if (!this.œilGauche || !this.œilDroit) return;
    
    const rect = this.avatar.getBoundingClientRect();
    const avatarCenter = {
      x: rect.left + rect.width / 2,
      y: rect.top + rect.height / 2
    };
    
    // Calcul de l'angle vers le curseur
    const angle = Math.atan2(
      event.clientY - avatarCenter.y,
      event.clientX - avatarCenter.x
    );
    
    // Distance limitée pour un mouvement naturel
    const maxDistance = 2;
    const distance = Math.min(maxDistance, 2);
    
    // Mouvement des pupilles
    const moveX = Math.cos(angle) * distance;
    const moveY = Math.sin(angle) * distance;
    
    this.œilGauche.style.transform = `translate(${moveX}px, ${moveY}px)`;
    this.œilDroit.style.transform = `translate(${moveX}px, ${moveY}px)`;
  }
  
  startBlinking() {
    // Clignements aléatoires toutes les 3-5 secondes
    const blink = () => {
      if (this.œilYeux) {
        this.œilYeux.classList.add('clignement');
        setTimeout(() => {
          this.œilYeux.classList.remove('clignement');
        }, 150);
      }
      
      // Prochain clignement aléatoire
      const nextBlink = 3000 + Math.random() * 2000;
      setTimeout(blink, nextBlink);
    };
    
    setTimeout(blink, 2000);
  }
  
  setupReactions() {
    // Réaction au click
    this.avatar.addEventListener('click', () => {
      this.showReaction('click');
    });
    
    // Réaction au hover
    this.avatar.addEventListener('mouseenter', () => {
      this.showReaction('hover');
    });
    
    // Réaction quand l'utilisateur quitte
    this.avatar.addEventListener('mouseleave', () => {
      this.resetEyes();
    });
  }
  
  showReaction(type) {
    switch(type) {
      case 'click':
        this.createParticles();
        this.avatar.classList.add('succes');
        setTimeout(() => this.avatar.classList.remove('succes'), 600);
        break;
      case 'hover':
        this.dilatePupils();
        break;
    }
  }
  
  createParticles() {
    const particulesContainer = this.avatar.querySelector('.avatar-particules');
    if (!particulesContainer) return;
    
    // Créer 5 particules autour de l'avatar
    for (let i = 0; i < 5; i++) {
      const particule = document.createElement('div');
      particule.className = 'particule';
      particule.style.left = '50%';
      particule.style.top = '50%';
      
      // Direction aléatoire
      const angle = (Math.PI * 2 * i) / 5;
      const distance = 30 + Math.random() * 20;
      particule.style.setProperty('--tx', `${Math.cos(angle) * distance}px`);
      particule.style.setProperty('--ty', `${Math.sin(angle) * distance}px`);
      
      particulesContainer.appendChild(particule);
      
      // Nettoyer après l'animation
      setTimeout(() => particule.remove(), 2000);
    }
  }
  
  dilatePupils() {
    if (this.œilGauche && this.œilDroit) {
      this.œilGauche.style.transform = 'scale(1.3)';
      this.œilDroit.style.transform = 'scale(1.3)';
    }
  }
  
  resetEyes() {
    if (this.œilGauche && this.œilDroit) {
      this.œilGauche.style.transform = 'scale(1)';
      this.œilDroit.style.transform = 'scale(1)';
    }
  }
  
  // Changer l'état émotionnel avec transition
  setEmotionalState(newState) {
    const oldState = this.getCurrentState();
    
    // Transition en douceur
    this.avatar.style.transition = 'all 0.3s ease';
    
    // Retirer ancienne classe
    this.avatar.classList.remove(`avatar-${oldState}`);
    
    // Ajouter nouvelle classe
    this.avatar.classList.add(`avatar-${newState}`);
    
    // Réaction visuelle
    this.createParticles();
    
    setTimeout(() => {
      this.avatar.style.transition = '';
    }, 300);
  }
  
  getCurrentState() {
    const classes = this.avatar.className.split(' ');
    const stateClass = classes.find(cls => cls.startsWith('avatar-'));
    return stateClass ? stateClass.replace('avatar-', '') : 'écoute';
  }
}

// Initialisation automatique
document.addEventListener('DOMContentLoaded', () => {
  const avatars = document.querySelectorAll('.avatar-prometheus');
  avatars.forEach(avatar => new AvatarEyeTracking(avatar));
});

// Export pour React
export default AvatarEyeTracking;
