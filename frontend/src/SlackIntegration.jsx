import React, { useState, useEffect } from 'react';
import { Plus, Trash2, Check, AlertCircle } from 'lucide-react';
import './SlackIntegration.css';

const SlackIntegration = ({ userPlan }) => {
  const [integration, setIntegration] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  // Only Pro and Business users can use Slack
  const canUseSlack = userPlan === 'pro' || userPlan === 'business';

  const fetchIntegration = async () => {
    const response = await fetch('/api/slack/integration', {
      credentials: 'include'
    });

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
      throw new Error(data.error || 'Failed to load Slack integration');
    }

    return data;
  };

  const handleStartAuth = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch('/api/slack/start-auth', {
        credentials: 'include'
      });

      const data = await response.json().catch(() => ({}));

      if (!response.ok) {
        throw new Error(data.error || 'Failed to start authentication');
      }
      
      // Redirect to Slack OAuth
      window.location.href = data.oauth_url;
    } catch (err) {
      console.error('Error starting Slack auth:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDisable = async () => {
    if (!window.confirm('Are you sure you want to disable Slack notifications?')) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const response = await fetch('/api/slack/disable', {
        method: 'POST',
        credentials: 'include'
      });

      const data = await response.json().catch(() => ({}));

      if (!response.ok) {
        throw new Error(data.error || 'Failed to disable Slack integration');
      }

      setIntegration(null);
      setSuccess(data.message || 'Slack integration disabled successfully');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      console.error('Error disabling Slack:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let isMounted = true;

    if (!canUseSlack) {
      setLoading(false);
      return () => {
        isMounted = false;
      };
    }

    const loadIntegration = async () => {
      setLoading(true);
      try {
        const data = await fetchIntegration();
        if (!isMounted) {
          return;
        }

        setIntegration(data.integration);
        setError(null);
      } catch (err) {
        if (!isMounted) {
          return;
        }

        console.error('Error fetching Slack integration:', err);
        setError(err.message);
        setIntegration(null);
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    loadIntegration();

    return () => {
      isMounted = false;
    };
  }, [canUseSlack]);

  if (loading) {
    return <div className="slack-integration loading">Loading...</div>;
  }

  if (!canUseSlack) {
    return (
      <div className="slack-integration plan-locked">
        <div className="locked-message">
          <AlertCircle size={24} />
          <div>
            <h3>Slack Integration</h3>
            <p>Available for Pro and Business plans</p>
            <small>Upgrade your plan to get Slack notifications for downtime alerts</small>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="slack-integration">
      <div className="slack-header">
        <div className="slack-title">
          <img src="https://a.slack-edge.com/80588/img/icons/ios_144.png" alt="Slack" className="slack-logo" />
          <div>
            <h3>Slack Integration</h3>
            <p>Get instant notifications for downtime in Slack</p>
          </div>
        </div>
      </div>

      {error && (
        <div className="alert alert-error">
          <AlertCircle size={16} />
          {error}
        </div>
      )}

      {success && (
        <div className="alert alert-success">
          <Check size={16} />
          {success}
        </div>
      )}

      {!integration ? (
        <div className="slack-content">
          <div className="slack-info">
            <h4>Connect your Slack workspace</h4>
            <p>When your monitored services go down, StatusFrame will send immediate alerts to your Slack channel.</p>
            
            <div className="slack-features">
              <div className="feature">
                <Check size={16} />
                <span>Instant downtime alerts</span>
              </div>
              <div className="feature">
                <Check size={16} />
                <span>Recovery notifications</span>
              </div>
              <div className="feature">
                <Check size={16} />
                <span>Customizable channels</span>
              </div>
            </div>
          </div>

          <button 
            className="btn-slack-connect"
            onClick={handleStartAuth}
            disabled={loading}
          >
            <Plus size={20} />
            Connect to Slack
          </button>
        </div>
      ) : (
        <div className="slack-connected">
          <div className="connected-info">
            <div className="status-badge">
              <Check size={16} />
              Connected
            </div>
            
            <div className="connection-details">
              <div className="detail-row">
                <span className="label">Workspace:</span>
                <span className="value">{integration.slack_team_name}</span>
              </div>
              
              <div className="detail-row">
                <span className="label">Channel:</span>
                <span className="value">#{integration.slack_channel_name}</span>
              </div>
              
              <div className="detail-row">
                <span className="label">Status:</span>
                <span className={`badge ${integration.is_enabled ? 'badge-active' : 'badge-inactive'}`}>
                  {integration.is_enabled ? 'Active' : 'Inactive'}
                </span>
              </div>
            </div>

            <div className="connection-timestamp">
              Connected on {new Date(integration.created_at).toLocaleDateString()}
            </div>
          </div>

          <button 
            className="btn-slack-disconnect"
            onClick={handleDisable}
            disabled={loading}
          >
            <Trash2 size={16} />
            Disconnect
          </button>
        </div>
      )}
    </div>
  );
};

export default SlackIntegration;
