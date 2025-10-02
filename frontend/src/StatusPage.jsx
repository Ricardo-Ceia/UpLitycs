import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import './StatusPage.css';
import UptimeBarGraph from './UptimeBarGraph';
import { Activity, CheckCircle, XCircle, Clock, Globe, Settings } from 'lucide-react';

const StatusPage = () => {
  const { slug } = useParams();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [statusData, setStatusData] = useState(null);
  const [currentTime, setCurrentTime] = useState(new Date());
  const [isOwner, setIsOwner] = useState(false);
  const [showThemeSelector, setShowThemeSelector] = useState(false);
  const [currentUserId, setCurrentUserId] = useState(null);

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  // Check if user is authenticated and owner
  useEffect(() => {
    const checkOwnership = async () => {
      try {
        const response = await fetch('/api/user-status', {
          credentials: 'include'
        });
        if (response.ok) {
          const userData = await response.json();
          setCurrentUserId(userData.userId);
        }
      } catch (err) {
        // Not authenticated, that's fine
        console.log('User not authenticated');
      }
    };
    checkOwnership();
  }, []);

  useEffect(() => {
    fetchStatusData();
    const interval = setInterval(fetchStatusData, 30000); // Refresh every 30 seconds
    return () => clearInterval(interval);
  }, [slug]);

  // Check ownership when both currentUserId and statusData are available
  useEffect(() => {
    if (currentUserId && statusData?.user_id) {
      setIsOwner(currentUserId === statusData.user_id);
    }
  }, [currentUserId, statusData]);

  const fetchStatusData = async () => {
    try {
      setLoading(true);
      const response = await fetch(`/api/public/status/${slug}`);
      
      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('Status page not found');
        }
        throw new Error('Failed to load status data');
      }
      
      const data = await response.json();
      setStatusData(data);
      setError(null);
    } catch (err) {
      console.error('Error fetching status:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleThemeChange = async (newTheme) => {
    try {
      const response = await fetch('/api/update-theme', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ theme: newTheme }),
      });

      if (response.ok) {
        // Refresh status data to get new theme
        await fetchStatusData();
        setShowThemeSelector(false);
      } else {
        console.error('Failed to update theme');
      }
    } catch (err) {
      console.error('Error updating theme:', err);
    }
  };

  const getStatusColor = (statusCode) => {
    if (statusCode >= 200 && statusCode < 300) return 'operational';
    if (statusCode >= 300 && statusCode < 400) return 'degraded';
    return 'down';
  };

  const getStatusText = (statusCode) => {
    if (statusCode >= 200 && statusCode < 300) return 'All Systems Operational';
    if (statusCode >= 300 && statusCode < 400) return 'Degraded Performance';
    return 'Service Down';
  };

  const getStatusIcon = (statusCode) => {
    if (statusCode >= 200 && statusCode < 300) {
      return <CheckCircle className="status-icon" />;
    }
    if (statusCode >= 300 && statusCode < 400) {
      return <Activity className="status-icon pulse" />;
    }
    return <XCircle className="status-icon" />;
  };

  const formatTime = (dateString) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    });
  };

  if (loading && !statusData) {
    return (
      <div className={`status-container theme-${statusData?.theme || 'cyberpunk'}`}>
        <div className="crt-overlay"></div>
        <div className="scan-lines"></div>
        <div className="loading-container">
          <div className="loading-spinner"></div>
          <p>Loading status page...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`status-container theme-cyberpunk`}>
        <div className="crt-overlay"></div>
        <div className="scan-lines"></div>
        <div className="error-container">
          <XCircle className="error-icon-large" />
          <h1>Status Page Not Found</h1>
          <p>{error}</p>
          <p className="error-hint">Please check the URL and try again</p>
        </div>
      </div>
    );
  }

  const statusCode = statusData?.status_code || 0;
  const statusColor = getStatusColor(statusCode);
  const theme = statusData?.theme || 'cyberpunk';

  return (
    <div className={`status-container theme-${theme}`}>
      <div className="crt-overlay"></div>
      <div className="scan-lines"></div>
      <div className="grid-background"></div>

      {/* Header */}
      <header className="status-header">
        <div className="header-container">
          <div className="brand-info">
            <h1 className="app-name">
              {statusData?.app_name || 'Service Status'}
            </h1>
            <p className="app-subtitle">System Status Monitor</p>
          </div>
          <div className="header-actions">
            <div className="time-display">
              <Clock className="clock-icon" />
              <span>{currentTime.toLocaleTimeString('en-US', { hour12: false })}</span>
            </div>
            {isOwner && (
              <div className="theme-selector-wrapper">
                <button 
                  className="theme-toggle-btn"
                  onClick={() => setShowThemeSelector(!showThemeSelector)}
                  title="Change Theme"
                >
                  <Settings className="settings-icon" />
                  <span className="btn-text">Theme</span>
                </button>
                {showThemeSelector && (
                  <div className="theme-dropdown">
                    <button onClick={() => handleThemeChange('cyberpunk')}>Cyberpunk</button>
                    <button onClick={() => handleThemeChange('matrix')}>Matrix</button>
                    <button onClick={() => handleThemeChange('retro')}>Retro</button>
                    <button onClick={() => handleThemeChange('minimal')}>Minimal</button>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </header>

      {/* Main Status */}
      <section className="main-status">
        <div className={`status-hero status-${statusColor}`}>
          <div className="status-icon-container">
            {getStatusIcon(statusCode)}
          </div>
          <h2 className="status-message">{getStatusText(statusCode)}</h2>
          <p className="status-detail">
            Current Status Code: <span className="status-code">{statusCode}</span>
          </p>
        </div>
      </section>

      {/* Metrics Grid */}
      <section className="metrics-section">
        <div className="metrics-grid">
          <div className="metric-box">
            <div className="metric-icon">
              <Activity />
            </div>
            <div className="metric-content">
              <span className="metric-label">Status</span>
              <span className={`metric-value status-${statusColor}`}>
                {statusData?.status || 'Unknown'}
              </span>
            </div>
          </div>

          <div className="metric-box">
            <div className="metric-icon">
              <CheckCircle />
            </div>
            <div className="metric-content">
              <span className="metric-label">24h Uptime</span>
              <span className="metric-value">{statusData?.uptime_24h?.toFixed(2) || 0}%</span>
            </div>
          </div>

          <div className="metric-box">
            <div className="metric-icon">
              <Clock />
            </div>
            <div className="metric-content">
              <span className="metric-label">Response Time</span>
              <span className="metric-value metric-value-small">
                {statusData?.response_time_ms || 0} ms
              </span>
            </div>
          </div>

          <div className="metric-box">
            <div className="metric-icon">
              <Globe />
            </div>
            <div className="metric-content">
              <span className="metric-label">Last Checked</span>
              <span className="metric-value metric-value-small">
                {formatTime(statusData?.checked_at)}
              </span>
            </div>
          </div>
        </div>
      </section>

      {/* 30-Day Uptime Bar Graph */}
      {statusData?.uptime_history && statusData.uptime_history.length > 0 && (
        <section className="graph-section">
          <UptimeBarGraph uptimeHistory={statusData.uptime_history} />
        </section>
      )}

      {/* Additional Info */}
      <section className="info-section">
        <div className="info-card">
          <h3 className="info-title">About This Status Page</h3>
          <p className="info-text">
            This page shows the real-time status of our services. 
            We check the health of our systems every 30 seconds to ensure everything is running smoothly.
          </p>
          <div className="info-details">
            <div className="info-item">
              <span className="info-label">Endpoint:</span>
              <span className="info-value">{statusData?.endpoint}</span>
            </div>
            <div className="info-item">
              <span className="info-label">Check Interval:</span>
              <span className="info-value">30 seconds</span>
            </div>
          </div>
          <div className="info-badge">
            <span className="badge-dot"></span>
            <span>Auto-refreshing every 30 seconds</span>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="status-footer">
        <div className="footer-container">
          <div className="footer-text">
            Powered by <span className="footer-brand">UPLYTICS</span>
          </div>
          <div className="footer-status">
            <span className={`footer-dot status-${statusColor}`}></span>
            <span>Updated {formatTime(statusData?.checked_at)}</span>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default StatusPage;
