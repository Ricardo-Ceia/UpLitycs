import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, TrendingUp, Activity, Clock, ExternalLink, Trash2, AlertCircle } from 'lucide-react';
import './Dashboard.css';

const Dashboard = () => {
  const [apps, setApps] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [planInfo, setPlanInfo] = useState({ plan: 'free', plan_limit: 1, app_count: 0 });
  const [showUpgradeModal, setShowUpgradeModal] = useState(false);
  const [deleteConfirm, setDeleteConfirm] = useState(null);
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
              <span className="glitch" data-text="UPLYTICS">UPLYTICS</span>
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
                    <h3 className="app-name">{app.app_name || 'Unnamed App'}</h3>
                    <span className="app-slug">/{app.slug || 'no-slug'}</span>
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

      {/* Upgrade Modal */}
      {showUpgradeModal && (
        <div className="modal-overlay" onClick={() => setShowUpgradeModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2>⚠️ Upgrade Required</h2>
            </div>
            <div className="modal-body">
              <p>You've reached your {planInfo.plan.toUpperCase()} plan limit ({planInfo.plan_limit} {planInfo.plan_limit === 1 ? 'app' : 'apps'}).</p>
              <p>Upgrade to monitor more apps with faster checks!</p>
              
              <div className="upgrade-options">
                <div className="upgrade-option">
                  <h3>⚡ Pro Plan</h3>
                  <p className="price">$12/month</p>
                  <ul>
                    <li>10 monitors</li>
                    <li>1-minute checks</li>
                    <li>Unlimited alerts</li>
                  </ul>
                </div>
                <div className="upgrade-option highlighted">
                  <h3>🚀 Business Plan</h3>
                  <p className="price">$29/month</p>
                  <ul>
                    <li>50 monitors</li>
                    <li>30-second checks</li>
                    <li>Webhook alerts</li>
                  </ul>
                </div>
              </div>
            </div>
            <div className="modal-footer">
              <button className="modal-btn secondary" onClick={() => setShowUpgradeModal(false)}>
                Cancel
              </button>
              <button className="modal-btn primary" onClick={() => navigate('/pricing')}>
                View Pricing
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
