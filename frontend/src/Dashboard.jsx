import React, { useState, useEffect } from 'react';
import './Dashboard.css';

const Dashboard = () => {
  const [homepageStatus, setHomepageStatus] = useState({
    url: 'https://example.com',
    status: 'operational',
    responseTime: '234ms',
    uptime: '99.95%',
    lastChecked: new Date().toLocaleTimeString(),
    statusCode: 200
  });

  const [currentTime, setCurrentTime] = useState(new Date());

  // Update time every second
  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  // Simulate status checks every 30 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      const statuses = ['operational', 'operational', 'operational', 'degraded', 'down'];
      const randomStatus = statuses[Math.floor(Math.random() * statuses.length)];
      
      setHomepageStatus(prev => ({
        ...prev,
        status: randomStatus,
        responseTime: randomStatus === 'down' ? 'Timeout' : `${Math.floor(Math.random() * 300) + 50}ms`,
        lastChecked: new Date().toLocaleTimeString(),
        statusCode: randomStatus === 'down' ? 0 : randomStatus === 'degraded' ? 500 : 200
      }));
    }, 30000);

    return () => clearInterval(interval);
  }, []);

  const getStatusColor = (status) => {
    switch (status) {
      case 'operational': return 'green';
      case 'degraded': return 'yellow';
      case 'down': return 'red';
      default: return 'gray';
    }
  };

  return (
    <div className="dashboard-container">
      {/* Retro CRT effect overlay */}
      <div className="crt-overlay"></div>
      
      {/* Header */}
      <header className="dashboard-header">
        <div className="header-content">
          <h1 className="dashboard-title">HOMEPAGE STATUS</h1>
          <div className="system-clock">
            <span className="clock-label">LAST CHECK</span>
            <span className="clock-time">{currentTime.toLocaleTimeString()}</span>
          </div>
        </div>
      </header>

      {/* Homepage Status Card */}
      <section className="status-section">
        <div className="homepage-card">
          <div className="status-header">
            <div className="site-info">
              <h2 className="site-url">{homepageStatus.url}</h2>
              <div className={`status-indicator ${getStatusColor(homepageStatus.status)}`}>
                <div className={`status-led ${getStatusColor(homepageStatus.status)} active`}></div>
                <span className="status-text">{homepageStatus.status.toUpperCase()}</span>
              </div>
            </div>
          </div>
          
          <div className="status-metrics">
            <div className="metric-row">
              <div className="metric">
                <span className="metric-label">Response Time</span>
                <span className="metric-value">{homepageStatus.responseTime}</span>
              </div>
              <div className="metric">
                <span className="metric-label">Status Code</span>
                <span className="metric-value">{homepageStatus.statusCode}</span>
              </div>
            </div>
            
            <div className="metric-row">
              <div className="metric">
                <span className="metric-label">Uptime</span>
                <span className="metric-value">{homepageStatus.uptime}</span>
              </div>
              <div className="metric">
                <span className="metric-label">Last Checked</span>
                <span className="metric-value">{homepageStatus.lastChecked}</span>
              </div>
            </div>
          </div>

          {/* Uptime Progress Bar */}
          <div className="uptime-bar">
            <div 
              className="uptime-fill"
              style={{ width: homepageStatus.uptime }}
            ></div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="dashboard-footer">
        <div className="footer-content">
          <div className="footer-text">
            &copy; 2025 UPLYTICS | Homepage Monitor
          </div>
          <div className="refresh-indicator">
            <span className="refresh-dot"></span>
            Auto-check: 30s
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Dashboard;
