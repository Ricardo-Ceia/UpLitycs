import React, { useState, useEffect } from 'react';
import './Dashboard.css';
import { UNSAFE_DataWithResponseInit } from 'react-router-dom';

const Dashboard = () => {

  const [loading, setLoading] = useState(true);
  const [error ,setError] = useState(null);
  const [latestStatus, setLatestStatus] = useState(null);
  const [currentTime,setCurrentTime] = useState(new Date());
  const [interval,setInterval] = useState(null);
  
  const latestStatusData = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/latest-status');
      if(!response.ok){
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      setLatestStatus(data);
      setError(null);
    }catch(error){
      console.error('Error fetching latest status:',error);
      setError('Failed to load status data. Please try again later.');
    }finally{
      setLoading(false)
    }
  }
  useEffect(()=>{
    const timer = setInterval(()=>{
      setCurrentTime(new Date());
    },1000);
    return () => clearInterval(timer);
  },[])

  useEffect(()=>{
    latestStatusData();
    const interval = setInterval(latestStatusData,30000);
    return () => clearInterval(interval);
  })

  const getStatusColor = (statusCode) => {
    if(statusCode >= 200 && statusCode < 300) return 'Operational';
    if(statusCode >=300 && statusCode < 400) return 'Degraded';
    return 'down';
  }

  const formatTime = (dateString) => {
    return new Date(dateString).toLocaleTimeString('en-US',{
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  }

  if (loading && !latestStatus){
    return (
      <div className = "dashboard-container">
        <div className = "loading-state">
          <h2>Loading status data...</h2>
        </div>
      </div>
    );
  }

  if (error && !latestStatus){
    return (
      <div className="dashboard-container">
        <div className="error-state">
          <h2>Error loading status data</h2>
          <p>{error}</p>
          <button onClick={latestStatusData}>Retry</button>
        </div>
      </div>
    )
  }

  const statusCode = latestStatus?.Status_code || 0;
  const status = latestStatus?.Status || 'unknown';
  const lastChecked = latestStatus?.CheckedAt || new Date().toISOString();     

    return (
    <div className="dashboard-container">
      <div className="crt-overlay"></div>
      
      <header className="dashboard-header">
        <div className="header-content">
          <h1 className="dashboard-title">HOMEPAGE STATUS</h1>
          <div className="system-clock">
            <span className="clock-label">LAST CHECK</span>
            <span className="clock-time">{formatTime(checkedAt)}</span>
          </div>
        </div>
      </header>

      <section className="status-section">
        <div className="homepage-card">
          <div className="status-header">
            <div className="site-info">
              <h2 className="site-url">Homepage Monitor</h2>
              <div className={`status-indicator ${getStatusColor(statusCode)}`}>
                <div className={`status-led ${getStatusColor(statusCode)} active`}></div>
                <span className="status-text">{getStatusText(statusCode)}</span>
              </div>
            </div>
          </div>
          
          <div className="status-metrics">
            <div className="metric-row">
              <div className="metric">
                <span className="metric-label">Status</span>
                <span className="metric-value">{status}</span>
              </div>
              <div className="metric">
                <span className="metric-label">Status Code</span>
                <span className="metric-value">{statusCode}</span>
              </div>
            </div>
            
            <div className="metric-row">
              <div className="metric">
                <span className="metric-label">Last Checked</span>
                <span className="metric-value">{formatTime(checkedAt)}</span>
              </div>
              <div className="metric">
                <span className="metric-label">Current Time</span>
                <span className="metric-value">{currentTime.toLocaleTimeString()}</span>
              </div>
            </div>
          </div>

          {error && (
            <div className="error-banner">
              <span>⚠️ {error}</span>
              <button onClick={latestStatusData}>Retry</button>
            </div>
          )}
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
