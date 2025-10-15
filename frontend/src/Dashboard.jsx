import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, TrendingUp, Activity, Clock, ExternalLink, Trash2, AlertCircle, Award, Copy, Check, X } from 'lucide-react';
import UpgradeModal from './UpgradeModal';
import PlanFeatures from './PlanFeatures';
import './Dashboard.css';

const Dashboard = () => {
  const [apps, setApps] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [planInfo, setPlanInfo] = useState({ plan: 'free', plan_limit: 1, app_count: 0 });
  const [showUpgradeModal, setShowUpgradeModal] = useState(false);
  const [deleteConfirm, setDeleteConfirm] = useState(null);
  const [showBadgeModal, setShowBadgeModal] = useState(null);
  const [badgePeriod, setBadgePeriod] = useState('24h');
  const [copiedBadge, setCopiedBadge] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchApps();
  }, []);

  const fetchApps = async () => {
    try {
      const response = await fetch('/api/user-apps', {
        credentials: 'include'
      });

      if (!response.ok) {
        throw new Error('Failed to fetch apps');
      }

      const data = await response.json();
      setApps(data.apps || []);
      setPlanInfo({
        plan: data.plan,
        plan_limit: data.plan_limit,
        app_count: data.app_count
      });
      setLoading(false);
    } catch (err) {
      console.error('Error fetching apps:', err);
      setError(err.message);
      setLoading(false);
    }
  };

  const handleAddApp = async () => {
    // Check if user can add more apps
    try {
      const response = await fetch('/api/check-plan-limit', {
        credentials: 'include'
      });
      const data = await response.json();

      if (data.can_add) {
        navigate('/onboarding');
      } else {
        setShowUpgradeModal(true);
      }
    } catch (err) {
      console.error('Error checking plan limit:', err);
    }
  };

  const handleDeleteApp = async (appId) => {
    try {
      const response = await fetch(`/api/apps/${appId}`, {
        method: 'DELETE',
        credentials: 'include'
      });

      if (!response.ok) {
        throw new Error('Failed to delete app');
      }

      // Refresh apps list
      fetchApps();
      setDeleteConfirm(null);
    } catch (err) {
      console.error('Error deleting app:', err);
      alert('Failed to delete app');
    }
  };

  const getStatusColor = (status) => {
    const statusMap = {
      'up': 'green',
      'degraded': 'yellow',
      'down': 'red',
      'error': 'red',
      'client_error': 'orange',
      'unknown': 'gray'
    };
    return statusMap[status] || 'gray';
  };

  const getStatusIcon = (status) => {
    if (status === 'up') return '✓';
    if (status === 'down' || status === 'error') return '✗';
    if (status === 'degraded') return '⚠';
    return '?';
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'Never';
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    
    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffMins < 1440) return `${Math.floor(diffMins / 60)}h ago`;
    return date.toLocaleDateString();
  };

  const getBadgeUrl = (slug, period = '24h') => {
    const baseUrl = window.location.origin;
    return `${baseUrl}/api/badge/${slug}?period=${period}`;
  };

  const getBadgeMarkdown = (slug, period = '24h') => {
    return `[![Uptime](${getBadgeUrl(slug, period)})](${window.location.origin}/status/${slug})`;
  };

  const getBadgeHtml = (slug, period = '24h') => {
    return `<a href="${window.location.origin}/status/${slug}"><img src="${getBadgeUrl(slug, period)}" alt="Uptime Badge" /></a>`;
  };

  const copyToClipboard = (text, type, appId) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopiedBadge(`${appId}-${type}`);
      setTimeout(() => setCopiedBadge(null), 2000);
    });
  };

  if (loading) {
    return (
      <div className="dashboard-container">
        <div className="loading-screen">
          <div className="loading-spinner"></div>
          <p>Loading your apps...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="dashboard-container">
        <div className="error-screen">
          <AlertCircle size={48} />
          <h2>Error Loading Dashboard</h2>
          <p>{error}</p>
          <button onClick={() => window.location.reload()}>Retry</button>
        </div>
      </div>
    );
  }

  return (
    <div className="dashboard-container">
      <div className="crt-overlay"></div>
      <div className="scan-lines"></div>

      {/* Header */}
      <header className="dashboard-header">
        <div className="header-content">
          <div className="brand-section">
            <h1 className="dashboard-title">
              <span className="glitch" data-text="STATUSFRAME">STATUSFRAME</span>
            </h1>
            <p className="dashboard-subtitle">Your Monitoring Dashboard</p>
          </div>

          <div className="header-actions">
            <div className="plan-badge">
              <span className="plan-icon">
                {planInfo.plan === 'free' && '🆓'}
                {planInfo.plan === 'pro' && '⚡'}
                {planInfo.plan === 'business' && '🚀'}
              </span>
              <span className="plan-name">{planInfo.plan.toUpperCase()}</span>
              <span className="plan-count">{planInfo.app_count}/{planInfo.plan_limit}</span>
            </div>
            
            {planInfo.plan === 'free' && (
              <button className="upgrade-btn" onClick={() => navigate('/pricing')}>
                <TrendingUp size={16} />
                Upgrade
              </button>
            )}
          </div>
        </div>
      </header>

      {/* Stats Overview */}
      <section className="stats-section">
        <div className="stat-card">
          <div className="stat-icon">
            <Activity />
          </div>
          <div className="stat-content">
            <span className="stat-label">Total Apps</span>
            <span className="stat-value">{apps.length}</span>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon status-green">
            <Activity />
          </div>
          <div className="stat-content">
            <span className="stat-label">Online</span>
            <span className="stat-value">
              {apps.filter(app => app.status === 'up').length}
            </span>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon status-red">
            <AlertCircle />
          </div>
          <div className="stat-content">
            <span className="stat-label">Issues</span>
            <span className="stat-value">
              {apps.filter(app => app.status === 'down' || app.status === 'error').length}
            </span>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon">
            <Clock />
          </div>
          <div className="stat-content">
            <span className="stat-label">Avg Uptime</span>
            <span className="stat-value">
              {apps.length > 0 
                ? (apps.reduce((sum, app) => sum + (app.uptime_24h || 0), 0) / apps.length).toFixed(1)
                : 0}%
            </span>
          </div>
        </div>
      </section>

      {/* Plan Features Section */}
      <PlanFeatures />

      {/* Apps Grid */}
      <section className="apps-section">
        <div className="section-header">
          <h2 className="section-title">Your Apps</h2>
          <button className="add-app-btn" onClick={handleAddApp}>
            <Plus size={20} />
            Add New App
          </button>
        </div>

        {apps.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">📊</div>
            <h3>No Apps Yet</h3>
            <p>Get started by adding your first monitoring app</p>
            <button className="cta-btn" onClick={handleAddApp}>
              <Plus size={20} />
              Add Your First App
            </button>
          </div>
        ) : (
          <div className="apps-grid">
            {apps.map((app) => (
              <div key={app.id} className={`app-card status-${getStatusColor(app.status)}`}>
                <div className="app-header">
                  <div className="app-info">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                      {app.logo_url && (
                        <img 
                          src={app.logo_url} 
                          alt={`${app.app_name} logo`}
                          style={{
                            width: '40px',
                            height: '40px',
                            objectFit: 'contain',
                            borderRadius: '4px'
                          }}
                          onError={(e) => {
                            e.target.style.display = 'none';
                          }}
                        />
                      )}
                      <div>
                        <h3 className="app-name">{app.app_name || 'Unnamed App'}</h3>
                        <span className="app-slug">/{app.slug || 'no-slug'}</span>
                      </div>
                    </div>
                  </div>
                  <div className={`status-indicator status-${getStatusColor(app.status)}`}>
                    {getStatusIcon(app.status)}
                  </div>
                </div>

                <div className="app-stats">
                  <div className="app-stat">
                    <span className="stat-label">Status</span>
                    <span className={`stat-value status-${getStatusColor(app.status)}`}>
                      {app.status}
                    </span>
                  </div>

                  <div className="app-stat">
                    <span className="stat-label">24h Uptime</span>
                    <span className="stat-value">{app.uptime_24h?.toFixed(1) || 0}%</span>
                  </div>

                  <div className="app-stat">
                    <span className="stat-label">Last Check</span>
                    <span className="stat-value">{formatDate(app.last_checked)}</span>
                  </div>
                </div>

                <div className="app-actions">
                  <button
                    className="action-btn view-btn"
                    onClick={() => navigate(`/status/${app.slug}`)}
                  >
                    <ExternalLink size={16} />
                    View Page
                  </button>
                  
                  <button
                    className="action-btn badge-btn"
                    onClick={() => setShowBadgeModal(app.id)}
                    title="Get Badge"
                  >
                    <Award size={16} />
                    Badge
                  </button>
                  
                  <button
                    className="action-btn delete-btn"
                    onClick={() => setDeleteConfirm(app.id)}
                  >
                    <Trash2 size={16} />
                  </button>
                </div>

                {/* Badge Modal */}
                {showBadgeModal === app.id && (
                  <div className="badge-modal-overlay" onClick={() => setShowBadgeModal(null)}>
                    <div className="badge-modal" onClick={(e) => e.stopPropagation()}>
                      <div className="modal-header">
                        <h3>Uptime Badge</h3>
                        <button className="close-btn" onClick={() => setShowBadgeModal(null)}>
                          <X size={20} />
                        </button>
                      </div>

                      <div className="modal-content">
                        <p className="modal-description">
                          Share your uptime status with these embeddable badges
                        </p>

                        {/* Period Selector */}
                        <div className="period-selector">
                          <label>Time Period:</label>
                          <div className="period-buttons">
                            {(() => {
                              // Determine available periods based on plan
                              const planPeriods = {
                                'free': ['24h', '7d'],
                                'pro': ['24h', '7d', '30d'],
                                'business': ['24h', '7d', '30d', '90d']
                              };
                              const availablePeriods = planPeriods[planInfo.plan] || ['24h', '7d'];
                              
                              return availablePeriods.map((period) => (
                                <button
                                  key={period}
                                  className={`period-btn ${badgePeriod === period ? 'active' : ''}`}
                                  onClick={() => setBadgePeriod(period)}
                                >
                                  {period}
                                </button>
                              ));
                            })()}
                          </div>
                          {planInfo.plan === 'free' && (
                            <p style={{ fontSize: '0.8rem', color: '#a0a8c0', marginTop: '0.5rem' }}>
                              💡 Upgrade to Pro for 30-day badges or Business for 90-day badges
                            </p>
                          )}
                          {planInfo.plan === 'pro' && (
                            <p style={{ fontSize: '0.8rem', color: '#a0a8c0', marginTop: '0.5rem' }}>
                              💡 Upgrade to Business for 90-day badges
                            </p>
                          )}
                        </div>

                        {/* Badge Preview */}
                        <div className="badge-preview">
                          <label>Preview:</label>
                          <div className="preview-box">
                            <img 
                              src={getBadgeUrl(app.slug, badgePeriod)} 
                              alt="Uptime Badge" 
                            />
                          </div>
                        </div>

                        {/* Markdown Code */}
                        <div className="code-section">
                          <label>Markdown (README.md, GitHub, etc.):</label>
                          <div className="code-box">
                            <code>{getBadgeMarkdown(app.slug, badgePeriod)}</code>
                            <button
                              className="copy-btn"
                              onClick={() => copyToClipboard(getBadgeMarkdown(app.slug, badgePeriod), 'markdown', app.id)}
                            >
                              {copiedBadge === `${app.id}-markdown` ? <Check size={16} /> : <Copy size={16} />}
                            </button>
                          </div>
                        </div>

                        {/* HTML Code */}
                        <div className="code-section">
                          <label>HTML:</label>
                          <div className="code-box">
                            <code>{getBadgeHtml(app.slug, badgePeriod)}</code>
                            <button
                              className="copy-btn"
                              onClick={() => copyToClipboard(getBadgeHtml(app.slug, badgePeriod), 'html', app.id)}
                            >
                              {copiedBadge === `${app.id}-html` ? <Check size={16} /> : <Copy size={16} />}
                            </button>
                          </div>
                        </div>

                        {/* Direct URL */}
                        <div className="code-section">
                          <label>Direct URL:</label>
                          <div className="code-box">
                            <code>{getBadgeUrl(app.slug, badgePeriod)}</code>
                            <button
                              className="copy-btn"
                              onClick={() => copyToClipboard(getBadgeUrl(app.slug, badgePeriod), 'url', app.id)}
                            >
                              {copiedBadge === `${app.id}-url` ? <Check size={16} /> : <Copy size={16} />}
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {deleteConfirm === app.id && (
                  <div className="delete-confirm">
                    <p>Delete this app?</p>
                    <div className="confirm-actions">
                      <button onClick={() => handleDeleteApp(app.id)}>Yes</button>
                      <button onClick={() => setDeleteConfirm(null)}>No</button>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </section>

      {/* Upgrade Modal */}
      {showUpgradeModal && (
        <UpgradeModal 
          planInfo={planInfo} 
          onClose={() => setShowUpgradeModal(false)} 
        />
      )}
    </div>
  );
};

export default Dashboard;
