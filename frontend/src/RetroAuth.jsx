// filepath: /home/ricardo/go/src/uplytics/frontend/src/Auth.jsx
import React, { useState, useEffect } from 'react';
import './RetroAuth.css';

const RetroAuth = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const [currentTime, setCurrentTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  const handleGoogleAuth = (authType) => {
    setIsLoading(true);
    setTimeout(() => {
      window.location.href = `/auth/google?mode=${authType}`;
    }, 800);
  };

  const formatTime = (date) => {
    return date.toLocaleTimeString('en-US', { 
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };

  const formatDate = (date) => {
    return date.toLocaleDateString('en-US', {
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: '2-digit'
    });
  };

  return (
    <div className="retro-auth-container">
      <div className="retro-background">
        <div className="grid-overlay"></div>
        <div className="scan-lines"></div>
      </div>

      <div className="auth-system-bar">
        <div className="system-info">
          <span className="system-name">UPLYTICS OS</span>
          <span className="version">v2.1.3</span>
        </div>
        <div className="datetime">
          <div className="time">{formatTime(currentTime)}</div>
          <div className="date">{formatDate(currentTime)}</div>
        </div>
      </div>

      <div className="auth-main-content">
        <div className="auth-window">
          <div className="window-header">
            <div className="window-title">
              <div className="title-icon">🔐</div>
              <span>User Authentication System</span>
            </div>
            <div className="window-controls">
              <div className="control-btn minimize">_</div>
              <div className="control-btn maximize">□</div>
              <div className="control-btn close">×</div>
            </div>
          </div>

          <div className="window-content">
            <div className="auth-header">
              <div className="logo-section">
                <div className="retro-logo">UPLYTICS</div>
                <div className="tagline">Website Status Monitoring System</div>
              </div>
            </div>

            <div className="auth-tabs">
              <button 
                className={`tab-button ${isLogin ? 'active' : ''}`}
                onClick={() => setIsLogin(true)}
                disabled={isLoading}
              >
                <span className="tab-icon">👤</span>
                LOGIN
              </button>
              <button 
                className={`tab-button ${!isLogin ? 'active' : ''}`}
                onClick={() => setIsLogin(false)}
                disabled={isLoading}
              >
                <span className="tab-icon">✚</span>
                SIGN UP
              </button>
            </div>

            <div className="auth-form">
              <div className="form-section">
                <h2 className="section-title">
                  {isLogin ? 'WELCOME BACK' : 'CREATE ACCOUNT'}
                </h2>
                <p className="section-subtitle">
                  {isLogin 
                    ? 'Please authenticate to access your dashboard' 
                    : 'Join UPLYTICS to monitor your websites'
                  }
                </p>

                <div className="provider-section">
                  <div className="provider-label">AUTHENTICATION PROVIDER</div>
                  
                  <button 
                    className="google-auth-button"
                    onClick={() => handleGoogleAuth(isLogin ? 'login' : 'signup')}
                    disabled={isLoading}
                  >
                    <div className="button-content">
                      <div className="google-icon">G</div>
                      <span className="button-text">
                        {isLogin ? 'SIGN IN WITH GOOGLE' : 'SIGN UP WITH GOOGLE'}
                      </span>
                      {isLoading && <div className="loading-spinner"></div>}
                    </div>
                  </button>
                </div>

                <div className="auth-footer">
                  <div className="security-notice">
                    🔒 Secure OAuth 2.0 Authentication
                  </div>
                  <div className="terms-notice">
                    By continuing, you agree to our Terms of Service
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {isLoading && (
        <div className="loading-overlay">
          <div className="loading-box">
            <div className="loading-title">AUTHENTICATING...</div>
            <div className="progress-bar">
              <div className="progress-fill"></div>
            </div>
            <div className="loading-text">Connecting to Google OAuth...</div>
          </div>
        </div>
      )}
    </div>
  );
};

export default RetroAuth;