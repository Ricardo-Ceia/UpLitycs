import React from 'react';
import { useNavigate } from 'react-router-dom';
import './UpgradeModal.css';

const UpgradeModal = ({ planInfo, onClose }) => {
  const navigate = useNavigate();

  const planFeatures = {
    free: {
      monitors: 1,
      interval: '5 minutes',
      features: ['1 monitor', '5-minute checks', 'Public status page', 'No alerts']
    },
    pro: {
      monitors: 10,
      interval: '1 minute',
      price: '$12',
      features: ['10 monitors', '1-minute checks', 'Email alerts', 'Custom domain', 'SSL monitoring']
    },
    business: {
      monitors: 50,
      interval: '30 seconds',
      price: '$29',
      features: ['50 monitors', '30-second checks', 'Email alerts', 'Webhook alerts', 'Priority support', 'API access']
    }
  };

  const currentPlan = planInfo?.plan || 'free';
  const usagePercent = planInfo?.plan_limit 
    ? Math.min((planInfo.app_count / planInfo.plan_limit) * 100, 100) 
    : 100;

  return (
    <div className="upgrade-modal-overlay" onClick={onClose}>
      <div className="upgrade-modal-container" onClick={(e) => e.stopPropagation()}>
        {/* Header */}
        <div className="upgrade-modal-header">
          <div className="glitch-wrapper">
            <h2 className="glitch-text" data-text="⚠️ UPGRADE REQUIRED">
              ⚠️ UPGRADE REQUIRED
            </h2>
          </div>
          <button className="modal-close-btn" onClick={onClose}>×</button>
        </div>

        {/* Usage Stats */}
        <div className="upgrade-stats">
          <div className="stat-card">
            <div className="stat-label">CURRENT PLAN</div>
            <div className="stat-value plan-badge">
              {currentPlan.toUpperCase()}
            </div>
          </div>
          
          <div className="stat-card">
            <div className="stat-label">MONITORS USED</div>
            <div className="stat-value">
              <span className="usage-current">{planInfo?.app_count || 0}</span>
              <span className="usage-separator">/</span>
              <span className="usage-limit">{planInfo?.plan_limit || 1}</span>
            </div>
          </div>
        </div>

        {/* Progress Bar */}
        <div className="upgrade-progress">
          <div className="progress-bar-container">
            <div 
              className="progress-bar-fill" 
              style={{ width: `${usagePercent}%` }}
            >
              <div className="progress-bar-glow"></div>
            </div>
          </div>
          <div className="progress-label">
            {usagePercent >= 100 ? '⚠️ LIMIT REACHED' : `${Math.round(usagePercent)}% CAPACITY`}
          </div>
        </div>

        {/* Message */}
        <div className="upgrade-message">
          <p>
            You've reached your <strong>{currentPlan.toUpperCase()}</strong> plan limit 
            of <strong>{planInfo?.plan_limit || 1}</strong> {planInfo?.plan_limit === 1 ? 'monitor' : 'monitors'}.
          </p>
          <p className="upgrade-cta">
            🚀 Upgrade now to unlock more monitors and advanced features!
          </p>
        </div>

        {/* Plan Options */}
        <div className="upgrade-plans">
          <div className="plan-card">
            <div className="plan-header pro-plan">
              <div className="plan-icon">⚡</div>
              <h3>PRO</h3>
              <div className="plan-price">{planFeatures.pro.price}<span>/month</span></div>
            </div>
            <div className="plan-features">
              {planFeatures.pro.features.map((feature, idx) => (
                <div key={idx} className="feature-item">
                  <span className="feature-check">✓</span>
                  <span>{feature}</span>
                </div>
              ))}
            </div>
            <div className="plan-highlight">
              Most Popular
            </div>
          </div>

          <div className="plan-card featured">
            <div className="plan-header business-plan">
              <div className="plan-icon">🚀</div>
              <h3>BUSINESS</h3>
              <div className="plan-price">{planFeatures.business.price}<span>/month</span></div>
            </div>
            <div className="plan-features">
              {planFeatures.business.features.map((feature, idx) => (
                <div key={idx} className="feature-item">
                  <span className="feature-check">✓</span>
                  <span>{feature}</span>
                </div>
              ))}
            </div>
            <div className="plan-highlight featured-badge">
              Best Value
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="upgrade-modal-footer">
          <button 
            className="btn-secondary" 
            onClick={onClose}
          >
            Maybe Later
          </button>
          <button 
            className="btn-primary btn-glow" 
            onClick={() => navigate('/pricing')}
          >
            <span>View All Plans</span>
            <span className="btn-arrow">→</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default UpgradeModal;
