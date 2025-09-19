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

  const [uptimeHistory, setUptimeHistory] = useState([]);
  const [currentTime, setCurrentTime] = useState(new Date());

  const generateUptimeHistory = () => {
    const days = [];
    const now = new Date();
    
    for (let i = 89; i >= 0; i--) {
      const date = new Date(now);
      date.setDate(date.getDate() - i);
      
      let uptimePercent;
      const randomOutage = Math.random();
      
      if (randomOutage < 0.02) { 
        uptimePercent = Math.random() * 50; 
      } else if (randomOutage < 0.08) { 
        uptimePercent = 70 + Math.random() * 29; 
      } else { 
        uptimePercent = 99 + Math.random() * 1; 
    }
      
      days.push({
        date: date,
        uptime: Math.round(uptimePercent * 100) / 100,
        status: uptimePercent >= 99 ? 'operational' : uptimePercent >= 90 ? 'degraded' : 'down'
      });
    }
    
    return days;
  };

  useEffect(() => {
    setUptimeHistory(generateUptimeHistory());
  }, []);

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);

    return () => clearInterval(timer);
  }, []);

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

  const formatDate = (date) => {
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric'
    });
  };

  const getUptimeColor = (uptime) => {
    if (uptime >= 99) return '#00FF7F';
    if (uptime >= 90) return '#FFB347';
         return '#FF6EC7'; 
  };

  const calculateOverallUptime = () => {
    if (uptimeHistory.length === 0) return 99.50;
    const total = uptimeHistory.reduce((sum, day) => sum + day.uptime, 0);
    return Math.round((total / uptimeHistory.length) * 100) / 100;
  };

  return (
    <div className="dashboard-container">
      <div className="crt-overlay"></div>
      
      <header className="dashboard-header">
        <div className="header-content">
          <h1 className="dashboard-title">HOMEPAGE STATUS</h1>
          <div className="system-clock">
            <span className="clock-label">LAST CHECK</span>
            <span className="clock-time">{currentTime.toLocaleTimeString()}</span>
          </div>
        </div>
      </header>

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
              <span className="metric-label">90-Day Uptime</span>
              <span className="metric-value">{calculateOverallUptime()}%</span>
            </div>
              <div className="metric">
                <span className="metric-label">Last Checked</span>
                <span className="metric-value">{homepageStatus.lastChecked}</span>
              </div>
            </div>
          </div>

          <div className="uptime-bar">
            <div 
              className="uptime-fill"
              style={{ width: homepageStatus.uptime }}
            ></div>
          </div>
        </div>
      </section>

      <section className="uptime-history-section">
        <div className="uptime-history-header">
          <h3 className="uptime-history-title">90-DAY UPTIME HISTORY</h3>
          <div className="uptime-legend">
            <div className="legend-item">
              <div className="legend-color operational"></div>
              <span>Operational</span>
            </div>
            <div className="legend-item">
              <div className="legend-color degraded"></div>
              <span>Degraded</span>
            </div>
            <div className="legend-item">
              <div className="legend-color down"></div>
              <span>Down</span>
            </div>
          </div>
        </div>
        
        <div className="uptime-history-chart">
          {uptimeHistory.map((day, index) => (
            <div
              key={index}
              className="uptime-bar"
              style={{ 
                backgroundColor: getUptimeColor(day.uptime),
                height: `${Math.max(day.uptime, 10)}%`
              }}
              title={`${formatDate(day.date)}: ${day.uptime}% uptime`}
            >
            </div>
          ))}
        </div>
        
        <div className="uptime-timeline">
          <span className="timeline-start">90 days ago</span>
          <span className="timeline-end">Today</span>
        </div>
      </section>

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
