import React, { useState, useEffect } from 'react';
import { Plus, Trash2, Check, AlertCircle } from 'lucide-react';
import './DiscordIntegration.css';

const DiscordIntegration = ({ userPlan }) => {
  const [integration, setIntegration] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  // Only Pro and Business users can use Discord
  const canUseDiscord = userPlan === 'pro' || userPlan === 'business';

  const fetchIntegration = async () => {
    const response = await fetch('/api/discord/integration', {
      credentials: 'include'
    });

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
      throw new Error(data.error || 'Failed to load Discord integration');
    }

    return data;
  };

  const handleStartAuth = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch('/api/discord/start-auth', {
        credentials: 'include'
      });

      const data = await response.json().catch(() => ({}));

      if (!response.ok) {
        throw new Error(data.error || 'Failed to start authentication');
      }
      
      // Redirect to Discord OAuth
      window.location.href = data.oauth_url;
    } catch (err) {
      console.error('Error starting Discord auth:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDisable = async () => {
    if (!window.confirm('Are you sure you want to disable Discord notifications?')) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const response = await fetch('/api/discord/disable', {
        method: 'POST',
        credentials: 'include'
      });

      const data = await response.json().catch(() => ({}));

      if (!response.ok) {
        throw new Error(data.error || 'Failed to disable Discord integration');
      }

      setIntegration(null);
      setSuccess(data.message || 'Discord integration disabled successfully');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      console.error('Error disabling Discord:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let isMounted = true;

    if (!canUseDiscord) {
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

        console.error('Error fetching Discord integration:', err);
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
  }, [canUseDiscord]);

  if (loading) {
    return <div className="discord-integration loading">Loading...</div>;
  }

  if (!canUseDiscord) {
    return (
      <div className="discord-integration plan-locked">
        <div className="locked-message">
          <AlertCircle size={24} />
          <div>
            <h3>Discord Integration</h3>
            <p>Available for Pro and Business plans</p>
            <small>Upgrade your plan to get Discord notifications for downtime alerts</small>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="discord-integration">
      <div className="discord-header">
        <div className="discord-title">
          <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 127.14 96.36'%3E%3Cdefs%3E%3Cstyle%3E.a%7Bfill:%235865F2;%7D%3C/style%3E%3C/defs%3E%3Cpath class='a' d='M107.7,8.07A105.15,105.15,0,0,0,81.47,0a72.06,72.06,0,0,0-3.36,6.83A97.68,97.68,0,0,0,49,6.83,72.37,72.37,0,0,0,45.64,0A105.89,105.89,0,0,0,19.39,8.09C2.79,32.65-1.71,56.6.54,80.21h0A105.73,105.73,0,0,0,32.71,96.36,77.7,77.7,0,0,0,39.6,85.25a68.15,68.15,0,0,1-10.85-5.18c.91-.66,1.8-1.34,2.66-2a77.52,77.52,0,0,0,64.32,0c.87.71,1.76,1.39,2.66,2a68.68,68.68,0,0,1-10.87,5.22,77,77,0,0,0,6.89,11.1A105.73,105.73,0,0,0,126.6,80.22h0C129.24,52.84,122.09,29.11,107.7,8.07ZM42.45,65.69C36.66,65.69,31.6,60.55,31.6,54s5-11.75,10.89-11.75S53.3,47.55,53.3,54,48.25,65.69,42.45,65.69Zm42.24,0C78.91,65.69,73.86,60.55,73.86,54s5-11.75,10.89-11.75S95.55,47.55,95.55,54,90.5,65.69,84.69,65.69Z'/%3E%3C/svg%3E" alt="Discord" className="discord-logo" />
          <div>
            <h3>Discord Integration</h3>
            <p>Get instant notifications for downtime in Discord</p>
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
        <div className="discord-content">
          <div className="discord-info">
            <h4>Connect your Discord account</h4>
            <p>When your monitored services go down, StatusFrame will send you instant alerts via Discord DMs.</p>
            
            <div className="discord-features">
              <div className="feature">
                <Check size={16} />
                <span>Instant downtime alerts via DM</span>
              </div>
              <div className="feature">
                <Check size={16} />
                <span>Recovery notifications</span>
              </div>
              <div className="feature">
                <Check size={16} />
                <span>No webhook setup needed</span>
              </div>
            </div>
          </div>

          <button 
            className="btn-discord-connect"
            onClick={handleStartAuth}
            disabled={loading}
          >
            <Plus size={20} />
            Connect Discord Account
          </button>
        </div>
      ) : (
        <div className="discord-connected">
          <div className="connected-info">
            <div className="status-badge">
              <Check size={16} />
              Connected
            </div>
            
            <div className="connection-details">
              <div className="detail-row">
                <span className="label">Discord User:</span>
                <span className="value">@{integration.discord_username}</span>
              </div>
              
              <div className="detail-row">
                <span className="label">Status:</span>
                <span className={`badge ${integration.is_enabled ? 'badge-active' : 'badge-inactive'}`}>
                  {integration.is_enabled ? 'Active - Receiving DMs' : 'Inactive'}
                </span>
              </div>
            </div>

            <div className="connection-timestamp">
              Connected on {new Date(integration.created_at).toLocaleDateString()}
            </div>
            
            <p style={{ fontSize: '0.9rem', color: '#888', marginTop: '1rem' }}>
              You'll receive Discord direct messages when your monitored services go down.
            </p>
          </div>

          <button 
            className="btn-discord-disconnect"
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

export default DiscordIntegration;
