import React, { useState, useEffect } from 'react';
import './Dashboard.css';

const Dashboard = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [latestStatus, setLatestStatus] = useState(null);
  const [currentTime, setCurrentTime] = useState(new Date());
  const [userTheme, setUserTheme] = useState('cyberpunk'); // default theme
  const [homepage, setHomepage] = useState('');

  const fetchUserData = async () => {
    try {
      const response = await fetch('/api/user-status', {
        credentials: 'include'
      });
      if (response.ok) {
        const data = await response.json();
        // Assuming the backend returns user data with theme
        setUserTheme(data.theme || 'cyberpunk');
        setHomepage(data.homepage || '');
      }
    } catch (error) {
      console.error('Error fetching user data:', error);
    }
  };

  const latestStatusData = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/latest-status', {
        credentials: 'include'
      });
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      setLatestStatus(data);
      setError(null);
    } catch (error) {
      console.error('Error fetching latest status:', error);
      setError('Failed to load status data. Please try again later.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUserData();
  }, []);

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    latestStatusData();
    const interval = setInterval(latestStatusData, 30000);
    return () => clearInterval(interval);
  }, []);

  const getStatusColor = (statusCode) => {
    if (statusCode >= 200 && statusCode < 300) return 'green';
    if (statusCode >= 300 && statusCode < 400) return 'yellow';
    return 'red';
  };

  const getStatusText = (statusCode) => {
    if (statusCode >= 200 && statusCode < 300) return 'OPERATIONAL';
    if (statusCode >= 300 && statusCode < 400) return 'DEGRADED';
    return 'DOWN';
  };

  const formatTime = (dateString) => {
    return new Date(dateString).toLocaleTimeString('en-US', {
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  };

  if (loading && !latestStatus) {
    return (
      <div className={`dashboard-container theme-${userTheme}`}>
        <div className="crt-overlay"></div>
        <div className="scan-lines"></div>
        <div className="loading-state">
          <div className="loading-spinner"></div>
          <h2>LOADING STATUS DATA...</h2>
        </div>
      </div>
    );
  }

  if (error && !latestStatus) {
    return (
      <div className={`dashboard-container theme-${userTheme}`}>
        <div className="crt-overlay"></div>
        <div className="scan-lines"></div>
        <div className="error-state">
          <div className="error-icon">⚠</div>
          <h2>ERROR LOADING STATUS DATA</h2>
          <p>{error}</p>
          <button onClick={latestStatusData} className="retry-button">
            RETRY CONNECTION
          </button>
        </div>
      </div>
    );
  }

  const statusCode = latestStatus?.Status_code || 0;
  const status = latestStatus?.Status || 'unknown';
  const checkedAt = latestStatus?.CheckedAt || new Date().toISOString();

  return (
    <div className={`dashboard-container theme-${userTheme}`}>
      <div className="crt-overlay"></div>
      <div className="scan-lines"></div>
      <div className="grid-background"></div>

      <header className="dashboard-header">
        <div className="header-content">
          <div className="brand-section">
            <h1 className="dashboard-brand">
              <span className="bracket">{"<"}</span>
              UPLYTICS
              <span className="bracket">{"/>"}</span>
            </h1>
            <span className="header-subtitle">STATUS DASHBOARD</span>
          </div>
          <div className="system-clock">
            <span className="clock-label">SYSTEM TIME</span>
            <span className="clock-time">{currentTime.toLocaleTimeString('en-US', { hour12: false })}</span>
          </div>
        </div>
      </header>

      <section className="status-section">
        <div className="status-header-bar">
          <h2 className="section-title">
            <span className="title-line">─────</span>
            MONITORED SERVICES
            <span className="title-line">─────</span>
          </h2>
        </div>

        <div className="homepage-card">
          <div className="card-header">
            <div className="site-info">
              <div className="site-url-container">
                <span className="url-label">[TARGET]</span>
                <h3 className="site-url">{homepage || 'Homepage'}</h3>
              </div>
              <div className={`status-badge status-${getStatusColor(statusCode)}`}>
                <div className={`status-led ${getStatusColor(statusCode)} active`}></div>
                <span className="status-text">{getStatusText(statusCode)}</span>
              </div>
            </div>
          </div>

          <div className="card-body">
            <div className="metrics-grid">
              <div className="metric-card">
                <div className="metric-icon">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
                  </svg>
                </div>
                <div className="metric-content">
                  <span className="metric-label">STATUS</span>
                  <span className="metric-value">{status}</span>
                </div>
              </div>

              <div className="metric-card">
                <div className="metric-icon">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
                  </svg>
                </div>
                <div className="metric-content">
                  <span className="metric-label">HTTP CODE</span>
                  <span className={`metric-value status-code-${getStatusColor(statusCode)}`}>
                    {statusCode}
                  </span>
                </div>
              </div>

              <div className="metric-card">
                <div className="metric-icon">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <circle cx="12" cy="12" r="10" />
                    <polyline points="12 6 12 12 16 14" />
                  </svg>
                </div>
                <div className="metric-content">
                  <span className="metric-label">LAST CHECK</span>
                  <span className="metric-value">{formatTime(checkedAt)}</span>
                </div>
              </div>

              <div className="metric-card">
                <div className="metric-icon">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                    <line x1="16" y1="2" x2="16" y2="6" />
                    <line x1="8" y1="2" x2="8" y2="6" />
                    <line x1="3" y1="10" x2="21" y2="10" />
                  </svg>
                </div>
                <div className="metric-content">
                  <span className="metric-label">DATE</span>
                  <span className="metric-value">{formatDate(checkedAt)}</span>
                </div>
              </div>
            </div>

            {error && (
              <div className="error-banner">
                <span className="error-icon">⚠</span>
                <span>{error}</span>
                <button onClick={latestStatusData} className="retry-small">RETRY</button>
              </div>
            )}
          </div>

          <div className="card-footer">
            <div className="footer-info">
              <span className="info-label">[AUTO-REFRESH]</span>
              <span className="info-value">30s interval</span>
            </div>
            <div className="footer-status">
              <span className="status-dot active"></span>
              <span>Monitoring Active</span>
            </div>
          </div>
        </div>
      </section>

      <footer className="dashboard-footer">
        <div className="footer-content">
          <div className="footer-text">
            © 2025 UPLYTICS | Monitoring System v2.1.3
          </div>
          <div className="footer-status">
            <span className="status-dot active"></span>
            <span>All Systems Operational</span>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Dashboard;