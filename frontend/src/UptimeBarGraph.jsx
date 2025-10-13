import React from 'react';
import './UptimeBarGraph.css';

const UptimeBarGraph = ({ uptimeHistory, dataRetentionDays = 30 }) => {
  // Sort by date ascending (oldest to newest)
  const sortedHistory = [...uptimeHistory].sort((a, b) => 
    new Date(a.date) - new Date(b.date)
  );

  // Create a full array for the retention period including unmonitored days
  const today = new Date();
  const startDate = new Date(today);
  startDate.setDate(today.getDate() - (dataRetentionDays - 1)); // e.g., 6 days ago + today = 7 days

  const displayDays = [];
  for (let i = 0; i < dataRetentionDays; i++) {
    const currentDate = new Date(startDate);
    currentDate.setDate(startDate.getDate() + i);
    const dateString = currentDate.toISOString().split('T')[0];
    
    // Find if we have data for this date
    const dataForDate = sortedHistory.find(d => d.date.split('T')[0] === dateString);
    
    displayDays.push({
      date: dateString,
      uptime_percentage: dataForDate?.uptime_percentage || 0,
      total_checks: dataForDate?.total_checks || 0,
      successful_checks: dataForDate?.successful_checks || 0,
      isMonitored: !!dataForDate
    });
  }

  // If no data, show empty state
  if (displayDays.length === 0) {
    return (
      <div className="uptime-graph-container">
        <h3 className="graph-title">{dataRetentionDays}-Day Uptime History</h3>
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
          {dataRetentionDays}-Day Uptime History
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
          <div className="legend-item">
            <span className="legend-dot unmonitored"></span>
            <span className="legend-label">Not Monitored</span>
          </div>
        </div>
      </div>

      <div className="graph-content">
        <div className="graph-bars">
          {displayDays.map((day, index) => {
            const percentage = day.uptime_percentage || 0;
            const height = (percentage / maxPercentage) * 100;
            const date = new Date(day.date);
            const dayLabel = date.toLocaleDateString('en-US', { 
              month: 'short', 
              day: 'numeric' 
            });
            const dayOfMonth = date.getDate();
            const monthLabel = date.toLocaleDateString('en-US', { month: 'short' });
            
            // Check if this is the first day of a month or first bar
            const isStartOfMonth = dayOfMonth === 1 || index === 0;

            // Render differently for unmonitored days
            if (!day.isMonitored) {
              return (
                <div key={index} className="bar-wrapper">
                  <div className="bar-container">
                    <div 
                      className="bar bar-unmonitored"
                      title={`${dayLabel}: Not monitored yet`}
                    >
                      <span className="bar-icon">⚠</span>
                    </div>
                  </div>
                  <div className="bar-date">{dayOfMonth}</div>
                  {isStartOfMonth && (
                    <div className="bar-month">{monthLabel}</div>
                  )}
                </div>
              );
            }

            return (
              <div key={index} className="bar-wrapper">
                <div className="bar-container">
                  <div 
                    className="bar bar-monitored"
                    style={{
                      height: `${Math.max(height, 5)}%`,
                      backgroundColor: getBarColor(percentage),
                    }}
                    title={`${dayLabel}: ${percentage.toFixed(2)}% uptime (${day.successful_checks}/${day.total_checks} checks)`}
                  >
                    <span className="bar-percentage">{percentage >= 10 ? percentage.toFixed(0) + '%' : ''}</span>
                  </div>
                </div>
                <div className="bar-date">{dayOfMonth}</div>
                {isStartOfMonth && (
                  <div className="bar-month">{monthLabel}</div>
                )}
              </div>
            );
          })}
        </div>

        <div className="graph-stats">
          <div className="stat-card">
            <span className="stat-label">Average Uptime</span>
            <span className="stat-value">
              {(() => {
                const monitoredDays = displayDays.filter(d => d.isMonitored);
                if (monitoredDays.length === 0) return '0.00';
                return (monitoredDays.reduce((sum, day) => sum + (day.uptime_percentage || 0), 0) / monitoredDays.length).toFixed(2);
              })()}%
            </span>
          </div>
          <div className="stat-card">
            <span className="stat-label">Total Checks</span>
            <span className="stat-value">
              {displayDays.reduce((sum, day) => sum + (day.total_checks || 0), 0).toLocaleString()}
            </span>
          </div>
          <div className="stat-card">
            <span className="stat-label">Days Monitored</span>
            <span className="stat-value">{displayDays.filter(d => d.isMonitored).length} / {dataRetentionDays}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UptimeBarGraph;
