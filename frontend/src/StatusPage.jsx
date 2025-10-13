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
  const [userTheme, setUserTheme] = useState(null); // User's local theme preference
  const [responseTime, setResponseTime] = useState(null); // Real-time response time
  const [pingLoading, setPingLoading] = useState(false);
  const isValidSlug = slug && slug !== 'undefined';

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
    if (!isValidSlug) {
      setStatusData(null);
      setError('No status page specified');
      return;
    }

    fetchStatusData();
    const interval = setInterval(fetchStatusData, 30000); // Refresh every 30 seconds
    return () => clearInterval(interval);
  }, [slug, isValidSlug]);

  // Load user's theme preference from localStorage
  useEffect(() => {
    if (!isValidSlug) return;

    const savedTheme = localStorage.getItem(`theme-preference-${slug}`);
    if (savedTheme) {
      setUserTheme(savedTheme);
    }
  }, [slug, isValidSlug]);

  // Check ownership when both currentUserId and statusData are available
  useEffect(() => {
    if (currentUserId && statusData?.user_id) {
      setIsOwner(currentUserId === statusData.user_id);
    }
  }, [currentUserId, statusData]);

  // Close theme selector when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (showThemeSelector && !event.target.closest('.theme-selector-wrapper')) {
        console.log('Clicking outside, closing theme selector');
        setShowThemeSelector(false);
      }
    };

    document.addEventListener('click', handleClickOutside);
    return () => document.removeEventListener('click', handleClickOutside);
  }, [showThemeSelector]);

  const fetchStatusData = async () => {
    if (!isValidSlug) {
      setStatusData(null);
      setError('No status page specified');
      setLoading(false);
      return;
    }
    
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

  const fetchResponseTime = async () => {
    if (!isValidSlug) {
      setResponseTime(null);
      return;
    }

    try {
      setPingLoading(true);
      const response = await fetch(`/api/public/ping/${slug}`);
      
      if (response.ok) {
        const data = await response.json();
        setResponseTime(data.response_time);
      }
    } catch (err) {
      console.error('Error fetching response time:', err);
      setResponseTime(null);
    } finally {
      setPingLoading(false);
    }
  };

  // Fetch response time on initial load
  useEffect(() => {
    if (statusData && isValidSlug) {
      fetchResponseTime();
    }
  }, [statusData, slug, isValidSlug]);

  const handleThemeChange = async (newTheme, isPermanent = false) => {
    console.log('Theme change requested:', { newTheme, isPermanent, isOwner });
    
    // If owner is changing theme permanently
    if (isPermanent && isOwner) {
      try {
        const response = await fetch('/api/update-theme', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({ 
            theme: newTheme,
            slug: slug  // Add slug to identify which app to update
          }),
        });

        if (response.ok) {
          console.log('Owner theme updated successfully');
          // Refresh status data to get new theme
          await fetchStatusData();
          setShowThemeSelector(false);
          // Clear user preference since owner changed default
          localStorage.removeItem(`theme-preference-${slug}`);
          setUserTheme(null);
        } else {
          console.error('Failed to update theme');
        }
      } catch (err) {
        console.error('Error updating theme:', err);
      }
    } else {
      // Regular user changing theme locally
      console.log('Setting user theme locally:', newTheme);
      setUserTheme(newTheme);
      localStorage.setItem(`theme-preference-${slug}`, newTheme);
      setShowThemeSelector(false);
      console.log('User theme set, localStorage updated');
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
    const missingSlug = error === 'No status page specified';
    return (
      <div className={`status-container theme-cyberpunk`}>
        <div className="crt-overlay"></div>
        <div className="scan-lines"></div>
        <div className="error-container">
          <XCircle className="error-icon-large" />
          <h1>{missingSlug ? 'No Monitors Yet' : 'Status Page Not Found'}</h1>
          <p>{missingSlug ? 'Add your first app to generate a public status page.' : error}</p>
          <p className="error-hint">
            {missingSlug ? 'Head to your dashboard to create a monitor and choose a slug.' : 'Please check the URL and try again.'}
          </p>
        </div>
      </div>
    );
  }

  const statusCode = statusData?.status_code || 0;
  const statusColor = getStatusColor(statusCode);
  // Use user's local theme preference if set, otherwise use owner's theme
  const theme = userTheme || statusData?.theme || 'cyberpunk';

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
            {/* Theme selector available for everyone */}
            <div className="theme-selector-wrapper">
              <button 
                className="theme-toggle-btn"
                onClick={(e) => {
                  e.stopPropagation();
                  console.log('Theme button clicked, current state:', showThemeSelector);
                  setShowThemeSelector(!showThemeSelector);
                }}
                title={isOwner ? "Change Theme (Owner - changes for everyone)" : "Change Theme (for your view)"}
              >
                <Settings className="settings-icon" />
                <span className="btn-text">Theme</span>
              </button>
              {showThemeSelector && (
                <div className="theme-dropdown" onClick={(e) => e.stopPropagation()}>
                  <div className="dropdown-section-header">
                    {isOwner ? '👑 Choose Theme (applies to all)' : '🎨 Choose Your Theme'}
                  </div>
                  <button 
                    onClick={(e) => {
                      e.stopPropagation();
                      handleThemeChange('cyberpunk', isOwner);
                    }}
                    className={theme === 'cyberpunk' ? 'active' : ''}
                  >
                    <span className="theme-preview cyberpunk-preview"></span>
                    Cyberpunk {theme === 'cyberpunk' && '✓'}
                  </button>
                  <button 
                    onClick={(e) => {
                      e.stopPropagation();
                      handleThemeChange('matrix', isOwner);
                    }}
                    className={theme === 'matrix' ? 'active' : ''}
                  >
                    <span className="theme-preview matrix-preview"></span>
                    Matrix {theme === 'matrix' && '✓'}
                  </button>
                  <button 
                    onClick={(e) => {
                      e.stopPropagation();
                      handleThemeChange('retro', isOwner);
                    }}
                    className={theme === 'retro' ? 'active' : ''}
                  >
                    <span className="theme-preview retro-preview"></span>
                    Retro {theme === 'retro' && '✓'}
                  </button>
                  <button 
                    onClick={(e) => {
                      e.stopPropagation();
                      handleThemeChange('minimal', isOwner);
                    }}
                    className={theme === 'minimal' ? 'active' : ''}
                  >
                    <span className="theme-preview minimal-preview"></span>
                    Minimal {theme === 'minimal' && '✓'}
                  </button>
                  {userTheme && !isOwner && (
                    <>
                      <div className="dropdown-divider"></div>
                      <button 
                        onClick={(e) => {
                          e.stopPropagation();
                          localStorage.removeItem(`theme-preference-${slug}`);
                          setUserTheme(null);
                          setShowThemeSelector(false);
                        }}
                        className="reset-theme-btn"
                      >
                        ↺ Reset to Default
                      </button>
                    </>
                  )}
                </div>
              )}
            </div>
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
                {pingLoading ? (
                  'Checking...'
                ) : responseTime !== null ? (
                  `${responseTime} ms`
                ) : (
                  'N/A'
                )}
              </span>
              {!pingLoading && responseTime !== null && (
                <button 
                  onClick={fetchResponseTime}
                  className="refresh-ping-btn"
                  title="Refresh response time"
                >
                  ↻
                </button>
              )}
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

      {/* Uptime Bar Graph (7, 30, or 90 days based on plan) */}
      {statusData?.uptime_history && statusData.uptime_history.length > 0 && (
        <section className="graph-section">
          <UptimeBarGraph 
            uptimeHistory={statusData.uptime_history} 
            dataRetentionDays={statusData.data_retention_days || 30}
          />
        </section>
      )}

      {/* Footer */}
      <footer className="status-footer">
        <div className="footer-container">
          <div className="footer-text">
            Powered by <span className="footer-brand">STATUSFRAME</span>
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
