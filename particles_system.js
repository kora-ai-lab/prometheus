// Particules créatives dynamiques
class ParticleSystem {
  constructor() {
    this.particles = [];
  }

  createParticle(x, y) {
    return { x, y, vx: Math.random() - 0.5, vy: Math.random() - 0.5 };
  }
}
