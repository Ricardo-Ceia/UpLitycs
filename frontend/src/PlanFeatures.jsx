import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { TrendingUp, Zap, Clock, Webhook, Globe, Shield, Code } from 'lucide-react';
import './PlanFeatures.css';

const PlanFeatures = () => {
  const [features, setFeatures] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchPlanFeatures();
  }, []);

  const fetchPlanFeatures = async () => {
    try {
      const response = await fetch('/api/plan-features', {
        credentials: 'include'
      });

      if (response.ok) {
        const data = await response.json();
        setFeatures(data);
      }
    } catch (err) {
      console.error('Error fetching plan features:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading || !features) {
    return null;
  }

  const planConfig = {
    free: {
      color: '#999',
      gradient: 'from-gray-600 to-gray-800',
      icon: '🆓'
    },
    pro: {
      color: '#00ffff',
      gradient: 'from-cyan-600 to-blue-600',
      icon: '⚡'
    },
    business: {
      color: '#ff00ff',
      gradient: 'from-purple-600 to-pink-600',
      icon: '🚀'
    }
  };

  const config = planConfig[features.plan] || planConfig.free;
  const usagePercent = (features.current_app_count / features.max_monitors) * 100;

  const formatInterval = (seconds) => {
    if (seconds >= 60) {
      return `${seconds / 60} minute${seconds > 60 ? 's' : ''}`;
    }
    return `${seconds} seconds`;
  };

  return (
    <div className="plan-features-container">
      {/* Current Plan Banner */}
      <div className={`plan-banner bg-gradient-to-r ${config.gradient}`}>
        <div className="plan-info">
          <span className="plan-icon">{config.icon}</span>
          <div>
            <h3 className="plan-name">{features.plan.toUpperCase()} PLAN</h3>
            <p className="plan-subtitle">
              {features.current_app_count} of {features.max_monitors} monitors used
            </p>
          </div>
        </div>
        
        {features.plan === 'free' && (
          <button 
            className="upgrade-btn-small"
            onClick={() => navigate('/pricing')}
          >
            <TrendingUp size={16} />
            Upgrade
          </button>
        )}
      </div>

      {/* Usage Progress */}
      <div className="usage-section">
        <div className="usage-header">
          <span className="usage-label">Monitor Usage</span>
          <span className="usage-count">
            {features.remaining_monitors} remaining
          </span>
        </div>
        <div className="usage-bar">
          <div 
            className="usage-fill"
            style={{ 
              width: `${usagePercent}%`,
              background: `linear-gradient(90deg, ${config.color}, ${config.color}aa)`
            }}
          />
        </div>
      </div>

      {/* Features Grid */}
      <div className="features-grid">
        <div className="feature-card">
          <div className="feature-icon">
            <Clock size={20} />
          </div>
          <div className="feature-content">
            <div className="feature-name">Check Interval</div>
            <div className="feature-value">
              {formatInterval(features.min_check_interval)}
            </div>
          </div>
        </div>

        <div className="feature-card">
          <div className="feature-icon">
            <Webhook size={20} />
          </div>
          <div className="feature-content">
            <div className="feature-name">Webhooks</div>
            <div className={`feature-value ${features.webhooks ? 'enabled' : 'disabled'}`}>
              {features.webhooks ? '✓ Enabled' : '✗ Disabled'}
            </div>
          </div>
        </div>

        <div className="feature-card">
          <div className="feature-icon">
            <Globe size={20} />
          </div>
          <div className="feature-content">
            <div className="feature-name">Custom Domain</div>
            <div className={`feature-value ${features.custom_domain ? 'enabled' : 'disabled'}`}>
              {features.custom_domain ? '✓ Enabled' : '✗ Disabled'}
            </div>
          </div>
        </div>

        <div className="feature-card">
          <div className="feature-icon">
            <Code size={20} />
          </div>
          <div className="feature-content">
            <div className="feature-name">API Access</div>
            <div className={`feature-value ${features.api_access ? 'enabled' : 'disabled'}`}>
              {features.api_access ? '✓ Enabled' : '✗ Disabled'}
            </div>
          </div>
        </div>

        <div className="feature-card">
          <div className="feature-icon">
            <Zap size={20} />
          </div>
          <div className="feature-content">
            <div className="feature-name">Alerts/Day</div>
            <div className="feature-value">
              {features.max_alerts_per_day}
            </div>
          </div>
        </div>
      </div>

      {/* Upgrade CTA for free users */}
      {features.plan === 'free' && (
        <div className="upgrade-cta">
          <p>🚀 Want more monitors and features?</p>
          <button 
            className="upgrade-btn-large"
            onClick={() => navigate('/pricing')}
          >
            View Upgrade Options
          </button>
        </div>
      )}
    </div>
  );
};

export default PlanFeatures;
