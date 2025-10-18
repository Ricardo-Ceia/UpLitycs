import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, TrendingUp, Activity, Clock, ExternalLink, Trash2, AlertCircle, Award, Copy, Check, X, Shield, ShieldAlert, Settings, MessageCircle, Loader2 } from 'lucide-react';
import UpgradeModal from './UpgradeModal';
import PlanFeatures from './PlanFeatures';
import './Dashboard.css';

const BADGE_PERIODS = {
  free: ['24h', '7d'],
  pro: ['24h', '7d', '30d'],
  business: ['24h', '7d', '30d', '90d']
};

const Dashboard = () => {
  const [apps, setApps] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [planInfo, setPlanInfo] = useState({ plan: 'free', plan_limit: 1, app_count: 0 });
  const [showUpgradeModal, setShowUpgradeModal] = useState(false);
  const [deleteConfirm, setDeleteConfirm] = useState(null);
  const [showBadgeModal, setShowBadgeModal] = useState(null);
  const [badgePeriods, setBadgePeriods] = useState({}); // Store period per app ID
  const [copiedBadge, setCopiedBadge] = useState(null);
  const [discordStatus, setDiscordStatus] = useState({
    loading: false,
    connected: false,
    integration: null,
    error: null,
  });
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

          await loadDiscordStatus(data.plan);
      } else {
        setShowUpgradeModal(true);
      }
    } catch (err) {
      console.error('Error checking plan limit:', err);
    }
  };


      const loadDiscordStatus = async (plan) => {
        if (plan !== 'pro' && plan !== 'business') {
          setDiscordStatus({
            loading: false,
            connected: false,
            integration: null,
            error: null,
          });
          return;
        }

        try {
          setDiscordStatus((prev) => ({ ...prev, loading: true, error: null }));
          const response = await fetch('/api/discord/integration', {
            credentials: 'include',
          });

          const data = await response.json().catch(() => ({}));

          if (!response.ok) {
            throw new Error(data.error || 'Failed to load Discord integration');
          }

          setDiscordStatus({
            loading: false,
            connected: Boolean(data.integration),
            integration: data.integration,
            error: null,
          });
        } catch (err) {
          console.error('Error loading Discord integration:', err);
          setDiscordStatus({
            loading: false,
            connected: false,
            integration: null,
            error: err.message,
          });
        }
      };

      const startDiscordAuth = async () => {
        if (planInfo.plan !== 'pro' && planInfo.plan !== 'business') {
          navigate('/pricing');
          return;
        }

        try {
          setDiscordStatus((prev) => ({ ...prev, loading: true, error: null }));
          const response = await fetch('/api/discord/start-auth', {
            credentials: 'include',
          });

          const data = await response.json().catch(() => ({}));

          if (!response.ok) {
            throw new Error(data.error || 'Failed to start Discord authentication');
          }

          window.location.href = data.oauth_url;
        } catch (err) {
          console.error('Error starting Discord auth:', err);
          setDiscordStatus((prev) => ({ ...prev, loading: false, error: err.message }));
        }
      };

      const handleManageDiscord = () => {
        navigate('/settings?tab=integrations');
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

  const getSSLStatus = (app) => {
    if (!app.health_url?.startsWith('https://')) {
      return null; // Not HTTPS
    }
    
    if (!app.ssl_days_until_expiry) {
      return { status: 'checking', color: 'gray', icon: Shield, text: 'Checking...', tooltip: 'SSL certificate check in progress' };
    }

    const days = app.ssl_days_until_expiry;
    
    if (days <= 0) {
      return { 
        status: 'expired', 
        color: 'red', 
        icon: ShieldAlert, 
        text: 'EXPIRED',
        tooltip: `SSL certificate has expired! Renew immediately.`,
        badge: '🔴 CRITICAL'
      };
    } else if (days <= 7) {
      return { 
        status: 'critical', 
        color: 'red', 
        icon: ShieldAlert, 
        text: `${days}d left`,
        tooltip: `SSL expires in ${days} day${days === 1 ? '' : 's'}. Renew urgently!`,
        badge: '⚠️ URGENT'
      };
    } else if (days <= 30) {
      return { 
        status: 'warning', 
        color: 'yellow', 
        icon: ShieldAlert, 
        text: `${days}d left`,
        tooltip: `SSL expires in ${days} days. Consider renewing soon.`,
        badge: '⚡ SOON'
      };
    } else {
      return { 
        status: 'ok', 
        color: 'green', 
        icon: Shield, 
        text: `${days}d left`,
        tooltip: `SSL certificate valid for ${days} more days`,
        badge: '✓ VALID'
      };
    }
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

  const getBadgePeriod = (appId) => {
    return badgePeriods[appId] || '24h';
  };

  const setBadgePeriod = (appId, period) => {
    setBadgePeriods(prev => ({
      ...prev,
      [appId]: period
    }));
  };

  const activeBadgeApp = showBadgeModal !== null
    ? apps.find((app) => app.id === showBadgeModal)
    : null;

  const activeBadgePeriod = activeBadgeApp
    ? getBadgePeriod(activeBadgeApp.id)
    : '24h';

  const availableBadgePeriods = BADGE_PERIODS[planInfo.plan] || BADGE_PERIODS.free;

  // Lock body scroll when modal is open
  useEffect(() => {
    if (!activeBadgeApp) {
      return undefined;
    }

    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';

    return () => {
      document.body.style.overflow = previousOverflow;
    };
  }, [activeBadgeApp]);

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
            
            <button 
              className="settings-btn" 
              onClick={() => navigate('/settings')}
              title="Settings"
            >
              <Settings size={18} />
            </button>
            
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

      {/* Discord Shortcut */}
      <section className="discord-shortcut">
        <div className={`discord-card ${discordStatus.connected ? 'discord-connected' : 'discord-disconnected'}`}>
          <div className="discord-card-header">
            <div className="discord-icon-wrapper">
              <img 
                src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 127.14 96.36'%3E%3Cdefs%3E%3Cstyle%3E.a%7Bfill:%235865F2;%7D%3C/style%3E%3C/defs%3E%3Cpath class='a' d='M107.7,8.07A105.15,105.15,0,0,0,81.47,0a72.06,72.06,0,0,0-3.36,6.83A97.68,97.68,0,0,0,49,6.83,72.37,72.37,0,0,0,45.64,0A105.89,105.89,0,0,0,19.39,8.09C2.79,32.65-1.71,56.6.54,80.21h0A105.73,105.73,0,0,0,32.71,96.36,77.7,77.7,0,0,0,39.6,85.25a68.15,68.15,0,0,1-10.85-5.18c.91-.66,1.8-1.34,2.66-2a77.52,77.52,0,0,0,64.32,0c.87.71,1.76,1.39,2.66,2a68.68,68.68,0,0,1-10.87,5.22,77,77,0,0,0,6.89,11.1A105.73,105.73,0,0,0,126.6,80.22h0C129.24,52.84,122.09,29.11,107.7,8.07ZM42.45,65.69C36.66,65.69,31.6,60.55,31.6,54s5-11.75,10.89-11.75S53.3,47.55,53.3,54,48.25,65.69,42.45,65.69Zm42.24,0C78.91,65.69,73.86,60.55,73.86,54s5-11.75,10.89-11.75S95.55,47.55,95.55,54,90.5,65.69,84.69,65.69Z'/%3E%3C/svg%3E"
                alt="Discord"
                className="discord-logo"
              />
            </div>
            <div>
              <h3>Discord Alerts</h3>
              <p>
                {planInfo.plan === 'pro' || planInfo.plan === 'business'
                  ? discordStatus.connected
                    ? `Connected to #${discordStatus.integration?.channel_name || 'channel'}`
                    : 'Instant incident alerts in your Discord server'
                  : 'Upgrade to enable Discord notifications'}
              </p>
            </div>
          </div>

          {discordStatus.error && (
            <div className="discord-card-error">
              <AlertCircle size={16} />
              <span>{discordStatus.error}</span>
            </div>
          )}

          <div className="discord-card-actions">
            {planInfo.plan === 'pro' || planInfo.plan === 'business' ? (
              discordStatus.connected ? (
                <div className="discord-connected-row">
                  <span className="discord-status-badge discord-status-active">Connected</span>
                  <button className="discord-manage-btn" onClick={handleManageDiscord}>
                    Manage
                  </button>
                </div>
              ) : (
                <div className="discord-connect-row">
                  <span className="discord-status-badge discord-status-inactive">Not connected</span>
                  <button
                    className="discord-connect-btn"
                    onClick={startDiscordAuth}
                    disabled={discordStatus.loading}
                  >
                    {discordStatus.loading ? (
                      <>
                        <Loader2 className="spinner" size={16} />
                        Connecting...
                      </>
                    ) : (
                      'Connect Discord'
                    )}
                  </button>
                  <button className="discord-manage-btn" onClick={handleManageDiscord}>
                    Learn more
                  </button>
                </div>
              )
            ) : (
              <div className="discord-upgrade-row">
                <span className="discord-status-badge discord-status-locked">Locked</span>
                <button className="discord-upgrade-btn" onClick={() => navigate('/pricing')}>
                  Upgrade for Discord Alerts
                </button>
              </div>
            )}
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

                  {getSSLStatus(app) && (planInfo.plan === 'pro' || planInfo.plan === 'business') && (
                    <div className="app-stat ssl-stat">
                      <span className="stat-label">
                        <Shield size={12} style={{ display: 'inline', marginRight: '4px' }} />
                        SSL Certificate
                      </span>
                      <div className={`ssl-status-container ssl-${getSSLStatus(app).status}`}>
                        <div className={`ssl-badge ssl-badge-${getSSLStatus(app).color}`}>
                          {getSSLStatus(app).badge}
                        </div>
                        <span 
                          className={`stat-value ssl-status-${getSSLStatus(app).color}`}
                          title={getSSLStatus(app).tooltip}
                        >
                          {(() => {
                            const SSLIcon = getSSLStatus(app).icon;
                            return (
                              <span style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                                <SSLIcon size={16} />
                                <span className="ssl-text">{getSSLStatus(app).text}</span>
                              </span>
                            );
                          })()}
                        </span>
                      </div>
                      {app.ssl_issuer && (
                        <div className="ssl-issuer">
                          Issuer: {app.ssl_issuer}
                        </div>
                      )}
                    </div>
                  )}

                  {getSSLStatus(app) && planInfo.plan === 'free' && (
                    <div className="app-stat ssl-stat-locked">
                      <span className="stat-label">
                        <Shield size={12} style={{ display: 'inline', marginRight: '4px' }} />
                        SSL Certificate
                      </span>
                      <div className="ssl-locked-message">
                        <span className="lock-icon">🔒</span>
                        <span>Upgrade to Pro</span>
                      </div>
                    </div>
                  )}
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

      {activeBadgeApp && (
        <div className="badge-modal-overlay" onClick={() => setShowBadgeModal(null)}>
          <div className="badge-modal" onClick={(e) => e.stopPropagation()}>
            <button className="badge-close-btn" onClick={() => setShowBadgeModal(null)}>
              <X size={20} />
            </button>

            <div className="badge-modal-content">
              <div className="period-selector">
                <label>TIME PERIOD:</label>
                <div className="period-buttons">
                  {availableBadgePeriods.map((period) => (
                    <button
                      key={period}
                      className={`period-btn ${activeBadgePeriod === period ? 'active' : ''}`}
                      onClick={() => setBadgePeriod(activeBadgeApp.id, period)}
                    >
                      {period}
                    </button>
                  ))}
                </div>
                {(planInfo.plan === 'free' || planInfo.plan === 'pro') && (
                  <p className="upgrade-hint">
                    💡 Upgrade to Business for 90-day badges
                  </p>
                )}
              </div>

              <div className="badge-preview">
                <label>PREVIEW:</label>
                <div className="preview-box">
                  <img
                    key={`${activeBadgeApp.id}-${activeBadgePeriod}`}
                    src={getBadgeUrl(activeBadgeApp.slug, activeBadgePeriod)}
                    alt="Uptime Badge"
                  />
                </div>
              </div>

              <div className="code-section">
                <label>MARKDOWN (README.MD, GITHUB, ETC.):</label>
                <div className="code-box">
                  <code>{getBadgeMarkdown(activeBadgeApp.slug, activeBadgePeriod)}</code>
                  <button
                    className="copy-btn"
                    onClick={() => copyToClipboard(getBadgeMarkdown(activeBadgeApp.slug, activeBadgePeriod), 'markdown', activeBadgeApp.id)}
                  >
                    {copiedBadge === `${activeBadgeApp.id}-markdown` ? <Check size={16} /> : <Copy size={16} />}
                  </button>
                </div>
              </div>

              <div className="code-section">
                <label>HTML:</label>
                <div className="code-box">
                  <code>{getBadgeHtml(activeBadgeApp.slug, activeBadgePeriod)}</code>
                  <button
                    className="copy-btn"
                    onClick={() => copyToClipboard(getBadgeHtml(activeBadgeApp.slug, activeBadgePeriod), 'html', activeBadgeApp.id)}
                  >
                    {copiedBadge === `${activeBadgeApp.id}-html` ? <Check size={16} /> : <Copy size={16} />}
                  </button>
                </div>
              </div>

              <div className="code-section">
                <label>DIRECT URL:</label>
                <div className="code-box">
                  <code>{getBadgeUrl(activeBadgeApp.slug, activeBadgePeriod)}</code>
                  <button
                    className="copy-btn"
                    onClick={() => copyToClipboard(getBadgeUrl(activeBadgeApp.slug, activeBadgePeriod), 'url', activeBadgeApp.id)}
                  >
                    {copiedBadge === `${activeBadgeApp.id}-url` ? <Check size={16} /> : <Copy size={16} />}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

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
