import React from 'react';
import './UptimeBarGraph.css';

const UptimeBarGraph = ({ uptimeHistory }) => {
  // Sort by date ascending (oldest to newest)
  const sortedHistory = [...uptimeHistory].sort((a, b) => 
    new Date(a.date) - new Date(b.date)
  );

  // Get last 30 days
  const last30Days = sortedHistory.slice(-30);

  // If no data, show empty state
  if (last30Days.length === 0) {
    return (
      <div className="uptime-graph-container">
        <h3 className="graph-title">30-Day Uptime History</h3>
        <div className="graph-empty-state">
          <p>No uptime data available yet</p>
          <p className="empty-hint">Data will appear as health checks are performed</p>
        </div>
      </div>
    );
  }

  const getBarColor = (percentage) => {
    if (percentage >= 99) return 'var(--uptime-excellent)';
    if (percentage >= 95) return 'var(--uptime-good)';
    if (percentage >= 90) return 'var(--uptime-warning)';
    return 'var(--uptime-critical)';
  };

  const getBarLabel = (percentage) => {
    if (percentage >= 99) return 'Excellent';
    if (percentage >= 95) return 'Good';
    if (percentage >= 90) return 'Fair';
    return 'Critical';
  };

  const maxPercentage = 100;

  return (
    <div className="uptime-graph-container">
      <div className="graph-header">
        <h3 className="graph-title">
          <span className="graph-icon">📊</span>
          30-Day Uptime History
        </h3>
        <div className="graph-legend">
          <div className="legend-item">
            <span className="legend-dot excellent"></span>
            <span className="legend-label">≥99%</span>
          </div>
          <div className="legend-item">
            <span className="legend-dot good"></span>
            <span className="legend-label">95-99%</span>
          </div>
          <div className="legend-item">
            <span className="legend-dot warning"></span>
            <span className="legend-label">90-95%</span>
          </div>
          <div className="legend-item">
            <span className="legend-dot critical"></span>
            <span className="legend-label">&lt;90%</span>
          </div>
        </div>
      </div>

      <div className="graph-content">
        <div className="graph-bars">
          {last30Days.map((day, index) => {
            const percentage = day.uptime_percentage || 0;
            const height = (percentage / maxPercentage) * 100;
            const date = new Date(day.date);
            const dayLabel = date.toLocaleDateString('en-US', { 
              month: 'short', 
              day: 'numeric' 
            });

            return (
              <div key={index} className="bar-wrapper">
                <div className="bar-container">
                  <div 
                    className="bar"
                    style={{
                      height: `${Math.max(height, 2)}%`,
                      backgroundColor: getBarColor(percentage),
                    }}
                    title={`${dayLabel}: ${percentage.toFixed(2)}% uptime (${day.successful_checks}/${day.total_checks} checks)`}
                  >
                    <span className="bar-percentage">{percentage.toFixed(0)}%</span>
                  </div>
                </div>
                {index % 5 === 0 && (
                  <div className="bar-label">{dayLabel}</div>
                )}
              </div>
            );
          })}
        </div>

        <div className="graph-stats">
          <div className="stat-card">
            <span className="stat-label">Average Uptime</span>
            <span className="stat-value">
              {(last30Days.reduce((sum, day) => sum + (day.uptime_percentage || 0), 0) / last30Days.length).toFixed(2)}%
            </span>
          </div>
          <div className="stat-card">
            <span className="stat-label">Total Checks</span>
            <span className="stat-value">
              {last30Days.reduce((sum, day) => sum + (day.total_checks || 0), 0).toLocaleString()}
            </span>
          </div>
          <div className="stat-card">
            <span className="stat-label">Days Monitored</span>
            <span className="stat-value">{last30Days.length}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UptimeBarGraph;
