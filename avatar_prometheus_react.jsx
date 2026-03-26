import React, { useState, useEffect } from 'react';
import './avatar_prometheus.css';

const AvatarPrometheus = ({ 
  state = 'écoute', 
  onStateChange,
  size = 'medium',
  interactive = true 
}) => {
  const [currentState, setCurrentState] = useState(state);
  const [isAnimating, setIsAnimating] = useState(false);

  const states = ['écoute', 'réflexion', 'illumination', 'co-création'];
  
  const handleClick = () => {
    if (!interactive || isAnimating) return;
    
    setIsAnimating(true);
    const currentIndex = states.indexOf(currentState);
    const nextState = states[(currentIndex + 1) % states.length];
    
    setCurrentState(nextState);
    onStateChange?.(nextState);
    
    setTimeout(() => setIsAnimating(false), 300);
  };

  const sizeClasses = {
    small: 'avatar-small',
    medium: 'avatar-medium', 
    large: 'avatar-large'
  };

  return (
    <div 
      className={`avatar-prometheus ${sizeClasses[size]} avatar-${currentState} ${interactive ? 'avatar-interactive' : ''}`}
      onClick={handleClick}
      role="img"
      aria-label={`Avatar Prometheus - état ${currentState}`}
      tabIndex={interactive ? 0 : -1}
      onKeyDown={(e) => e.key === 'Enter' && handleClick()}
    >
      <div className="avatar-cristaux">
        <div className="avatar-cœur">
          <div className="avatar-yeux">
            <div className="œil-gauche"></div>
            <div className="œil-droit"></div>
          </div>
          <div className="avatar-particules"></div>
        </div>
      </div>
      {interactive && (
        <div className="avatar-tooltip">
          Cliquez pour changer d'état
        </div>
      )}
    </div>
  );
};

export default AvatarPrometheus;
