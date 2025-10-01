import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import './PublicStatusPage.css';
import { Activity, CheckCircle, XCircle, Clock, Globe } from 'lucide-react';

const PublicStatusPage = () => {
  const { slug } = useParams();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [statusData, setStatusData] = useState(null);
  const [currentTime, setCurrentTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    fetchStatusData();
    const interval = setInterval(fetchStatusData, 30000); // Refresh every 30 seconds
    return () => clearInterval(interval);
  }, [slug]);

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

  const getUptime = () => {
    // Calculate uptime percentage (you can enhance this with historical data)
    if (statusData?.latestStatus?.Status_code >= 200 && statusData?.latestStatus?.Status_code < 300) {
      return '99.9%';
    }
    return '0%';
  };

  if (loading && !statusData) {
    return (
      <div className={`public-status-container theme-${statusData?.user?.Theme || 'cyberpunk'}`}>
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
      <div className={`public-status-container theme-cyberpunk`}>
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

  const statusCode = statusData?.latestStatus?.Status_code || 0;
  const statusColor = getStatusColor(statusCode);
  const theme = statusData?.user?.Theme || 'cyberpunk';

  return (
    <div className={`public-status-container theme-${theme}`}>
      <div className="crt-overlay"></div>
      <div className="scan-lines"></div>
      <div className="grid-background"></div>

      {/* Header */}
      <header className="public-header">
        <div className="header-container">
          <div className="brand-info">
            <h1 className="app-name">
              {statusData?.user?.AppName || 'Service Status'}
            </h1>
            <p className="app-subtitle">System Status Monitor</p>
          </div>
          <div className="time-display">
            <Clock className="clock-icon" />
            <span>{currentTime.toLocaleTimeString('en-US', { hour12: false })}</span>
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
                {statusData?.latestStatus?.Status || 'Unknown'}
              </span>
            </div>
          </div>

          <div className="metric-box">
            <div className="metric-icon">
              <CheckCircle />
            </div>
            <div className="metric-content">
              <span className="metric-label">Uptime</span>
              <span className="metric-value">{getUptime()}</span>
            </div>
          </div>

          <div className="metric-box">
            <div className="metric-icon">
              <Clock />
            </div>
            <div className="metric-content">
              <span className="metric-label">Last Checked</span>
              <span className="metric-value metric-value-small">
                {formatTime(statusData?.latestStatus?.CheckedAt)}
              </span>
            </div>
          </div>

          <div className="metric-box">
            <div className="metric-icon">
              <Globe />
            </div>
            <div className="metric-content">
              <span className="metric-label">Endpoint</span>
              <span className="metric-value metric-value-small">
                {statusData?.user?.HealthUrl?.substring(0, 30)}...
              </span>
            </div>
          </div>
        </div>
      </section>

      {/* Additional Info */}
      <section className="info-section">
        <div className="info-card">
          <h3 className="info-title">About This Status Page</h3>
          <p className="info-text">
            This page shows the real-time status of our services. 
            We check the health of our systems every 30 seconds to ensure everything is running smoothly.
          </p>
          <div className="info-badge">
            <span className="badge-dot"></span>
            <span>Auto-refreshing every 30 seconds</span>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="public-footer">
        <div className="footer-container">
          <div className="footer-text">
            Powered by <span className="footer-brand">UPLYTICS</span>
          </div>
          <div className="footer-status">
            <span className={`footer-dot status-${statusColor}`}></span>
            <span>Updated {formatTime(statusData?.latestStatus?.CheckedAt)}</span>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default PublicStatusPage;
