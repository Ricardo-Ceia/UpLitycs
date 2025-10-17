import React, { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Settings as SettingsIcon, LogOut, Lock, Shield, Bell, Slack as SlackIcon, AlertCircle, Check } from 'lucide-react';
import SlackIntegration from './SlackIntegration';
import './Settings.css';

const SettingsPage = () => {
  const [activeTab, setActiveTab] = useState('general');
  const [user, setUser] = useState(null);
  const [planInfo, setPlanInfo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    // Check for success message in URL
    const successParam = searchParams.get('success');
    if (successParam) {
      setSuccess(`✅ ${successParam} integration successful!`);
      setTimeout(() => setSuccess(null), 4000);
    }

    const errorParam = searchParams.get('error');
    if (errorParam) {
      setError(errorParam);
      setTimeout(() => setError(null), 4000);
    }

    // Set tab from URL if provided
    const tabParam = searchParams.get('tab');
    if (tabParam) {
      setActiveTab(tabParam);
    }

    fetchUserData();
  }, [searchParams]);

  const fetchUserData = async () => {
    try {
      setLoading(true);
      const [statusResponse, appsResponse] = await Promise.all([
        fetch('/api/user-status', { credentials: 'include' }),
        fetch('/api/user-apps', { credentials: 'include' })
      ]);

      if (!statusResponse.ok || !appsResponse.ok) {
        throw new Error('Failed to fetch user data');
      }

      const statusData = await statusResponse.json();
      const appsData = await appsResponse.json();

      setUser(statusData);
      setPlanInfo(appsData);
    } catch (err) {
      console.error('Error fetching user data:', err);
      setError('Failed to load settings');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await fetch('/auth/logout', { credentials: 'include' });
      navigate('/');
    } catch (err) {
      console.error('Error logging out:', err);
    }
  };

  if (loading) {
    return (
      <div className="settings-page loading">
        <div className="loader">Loading settings...</div>
      </div>
    );
  }

  if (!user || !planInfo) {
    return (
      <div className="settings-page error">
        <div className="error-message">{error || 'Failed to load settings'}</div>
      </div>
    );
  }

  return (
    <div className="settings-page">
      {error && (
        <div className="alert alert-error">
          <AlertCircle size={16} />
          {error}
        </div>
      )}

      {success && (
        <div className="alert alert-success">
          <Check size={16} />
          {success}
        </div>
      )}

      <div className="settings-container">
        {/* Header */}
        <div className="settings-header">
          <div className="header-content">
            <SettingsIcon size={32} />
            <div>
              <h1>Settings</h1>
              <p>Manage your account and preferences</p>
            </div>
          </div>

          <button className="btn-logout" onClick={handleLogout}>
            <LogOut size={18} />
            Logout
          </button>
        </div>

        <div className="settings-body">
          {/* Sidebar Tabs */}
          <aside className="settings-sidebar">
            <nav className="settings-tabs">
              <button
                className={`tab-btn ${activeTab === 'general' ? 'active' : ''}`}
                onClick={() => setActiveTab('general')}
              >
                <Shield size={18} />
                <span>General</span>
              </button>

              <button
                className={`tab-btn ${activeTab === 'integrations' ? 'active' : ''}`}
                onClick={() => setActiveTab('integrations')}
              >
                <SlackIcon size={18} />
                <span>Integrations</span>
                {planInfo.plan !== 'free' && (
                  <span className="tab-badge">Pro</span>
                )}
              </button>

              <button
                className={`tab-btn ${activeTab === 'notifications' ? 'active' : ''}`}
                onClick={() => setActiveTab('notifications')}
              >
                <Bell size={18} />
                <span>Notifications</span>
              </button>

              <button
                className={`tab-btn ${activeTab === 'security' ? 'active' : ''}`}
                onClick={() => setActiveTab('security')}
              >
                <Lock size={18} />
                <span>Security</span>
              </button>
            </nav>
          </aside>

          {/* Main Content */}
          <main className="settings-content">
            {/* General Tab */}
            {activeTab === 'general' && (
              <div className="settings-section">
                <h2>Account Settings</h2>

                <div className="setting-group">
                  <label>Avatar</label>
                  <div className="avatar-display">
                    {user.userName && (
                      <img 
                        src={user.userName} 
                        alt="User avatar" 
                        className="avatar-img"
                      />
                    )}
                    <div className="avatar-info">
                      <p className="avatar-name">Profile Picture</p>
                      <p className="avatar-hint">Managed by your OAuth provider</p>
                    </div>
                  </div>
                </div>

                <div className="setting-group">
                  <label>Plan Information</label>
                  <div className="plan-info-box">
                    <div className="plan-detail">
                      <span className="plan-label">Current Plan:</span>
                      <span className={`plan-badge plan-${planInfo.plan}`}>
                        {planInfo.plan.toUpperCase()}
                      </span>
                    </div>
                    <div className="plan-detail">
                      <span className="plan-label">Monitors:</span>
                      <span className="plan-value">
                        {planInfo.app_count} / {planInfo.plan_limit}
                      </span>
                    </div>
                    <div className="plan-detail">
                      <span className="plan-label">Status:</span>
                      <span className="plan-value status-active">Active</span>
                    </div>

                    <button 
                      className="btn-manage-plan"
                      onClick={() => window.location.href = '/pricing'}
                    >
                      Manage Plan
                    </button>
                  </div>
                </div>
              </div>
            )}

            {/* Integrations Tab */}
            {activeTab === 'integrations' && (
              <div className="settings-section">
                <h2>Integrations</h2>
                <p className="section-description">
                  Connect external services to receive notifications and automate workflows.
                </p>

                {planInfo.plan === 'free' ? (
                  <div className="integration-locked">
                    <Lock size={32} />
                    <h3>Upgrade to unlock integrations</h3>
                    <p>Integrations are available for Pro and Business plans</p>
                    <button 
                      className="btn-upgrade"
                      onClick={() => window.location.href = '/pricing'}
                    >
                      View Plans
                    </button>
                  </div>
                ) : (
                  <div className="integrations-list">
                    <SlackIntegration 
                      userPlan={planInfo.plan}
                    />
                  </div>
                )}
              </div>
            )}

            {/* Notifications Tab */}
            {activeTab === 'notifications' && (
              <div className="settings-section">
                <h2>Notifications</h2>
                <p className="section-description">
                  Control how and when you receive alerts about your services.
                </p>

                <div className="notification-settings">
                  <div className="notification-item">
                    <div className="notification-content">
                      <h3>Downtime Alerts</h3>
                      <p>Receive alerts when any of your monitored services go down</p>
                    </div>
                    <input type="checkbox" defaultChecked disabled />
                  </div>

                  <div className="notification-item">
                    <div className="notification-content">
                      <h3>Recovery Notifications</h3>
                      <p>Get notified when services come back online</p>
                    </div>
                    <input type="checkbox" defaultChecked disabled />
                  </div>

                  <div className="notification-item">
                    <div className="notification-content">
                      <h3>Weekly Summary</h3>
                      <p>Receive a weekly summary of your uptime metrics</p>
                    </div>
                    <input type="checkbox" disabled />
                  </div>

                  <p className="notification-hint">
                    💡 Enable integrations above to receive notifications via Slack, Email, and more.
                  </p>
                </div>
              </div>
            )}

            {/* Security Tab */}
            {activeTab === 'security' && (
              <div className="settings-section">
                <h2>Security & Privacy</h2>
                <p className="section-description">
                  Manage your account security and privacy settings.
                </p>

                <div className="security-settings">
                  <div className="security-item">
                    <div className="security-content">
                      <h3>Session Management</h3>
                      <p>You are currently logged in via Google OAuth</p>
                      <p className="session-info">
                        Last login: {new Date().toLocaleDateString()}
                      </p>
                    </div>
                  </div>

                  <div className="security-item">
                    <div className="security-content">
                      <h3>Two-Factor Authentication</h3>
                      <p>Add an extra layer of security to your account</p>
                    </div>
                    <button className="btn-secondary" disabled>
                      Coming Soon
                    </button>
                  </div>

                  <div className="security-item danger">
                    <div className="security-content">
                      <h3>Delete Account</h3>
                      <p>Permanently delete your account and all associated data</p>
                    </div>
                    <button className="btn-danger" disabled>
                      Coming Soon
                    </button>
                  </div>

                  <p className="security-notice">
                    🔒 Your data is encrypted and stored securely. OAuth tokens are never exposed.
                  </p>
                </div>
              </div>
            )}
          </main>
        </div>
      </div>
    </div>
  );
};

export default SettingsPage;
