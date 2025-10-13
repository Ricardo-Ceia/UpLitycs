import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './Admin.css';

function Admin() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [userEmail, setUserEmail] = useState('');
  const [users, setUsers] = useState([]);
  const [stats, setStats] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterPlan, setFilterPlan] = useState('all');
  const navigate = useNavigate();

  useEffect(() => {
    checkSession();
  }, []);

  useEffect(() => {
    if (isAuthenticated) {
      fetchAdminData();
    }
  }, [isAuthenticated]);

  const checkSession = async () => {
    try {
      const response = await fetch('/api/admin/check-session', {
        credentials: 'include'
      });
      const data = await response.json();
      setIsAuthenticated(data.authenticated);
      if (data.authenticated) {
        setUserEmail(data.email);
      }
    } catch (error) {
      console.error('Session check error:', error);
      setIsAuthenticated(false);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogin = () => {
    // Redirect to Google OAuth login
    window.location.href = '/auth/google';
  };

  const handleLogout = async () => {
    try {
      await fetch('/auth/logout', {
        credentials: 'include'
      });
      setIsAuthenticated(false);
      setUsers([]);
      setStats(null);
      navigate('/');
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  const fetchAdminData = async () => {
    try {
      // Fetch users
      const usersResponse = await fetch('/api/admin/users', {
        credentials: 'include'
      });
      if (usersResponse.ok) {
        const usersData = await usersResponse.json();
        setUsers(usersData);
      }

      // Fetch stats
      const statsResponse = await fetch('/api/admin/stats', {
        credentials: 'include'
      });
      if (statsResponse.ok) {
        const statsData = await statsResponse.json();
        setStats(statsData);
      }
    } catch (error) {
      console.error('Error fetching admin data:', error);
    }
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getSubscriptionEndDate = (planStartedAt) => {
    const startDate = new Date(planStartedAt);
    const endDate = new Date(startDate);
    endDate.setMonth(endDate.getMonth() + 1); // Add 1 month
    return formatDate(endDate);
  };

  const getPlanBadgeClass = (plan) => {
    switch (plan) {
      case 'pro':
        return 'badge-pro';
      case 'business':
        return 'badge-business';
      default:
        return 'badge-free';
    }
  };

  const filteredUsers = users.filter(user => {
    const matchesSearch = 
      user.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
      user.email.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesPlan = filterPlan === 'all' || user.plan === filterPlan;
    
    return matchesSearch && matchesPlan;
  });

  if (isLoading) {
    return (
      <div className="admin-container">
        <div className="loading-spinner">Loading...</div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="admin-container">
        <div className="admin-login-box">
          <h1>🔐 Admin Access</h1>
          <p className="admin-subtitle">Restricted Area - Admin Only</p>
          <p className="admin-info">Only ricardoceia.sete@gmail.com can access this panel</p>
          
          <button onClick={handleLogin} className="admin-login-btn">
            Sign in with Google
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-dashboard">
      <header className="admin-header">
        <div className="admin-header-content">
          <div>
            <h1>📊 Admin Dashboard</h1>
            <p className="admin-user">Logged in as: {userEmail}</p>
          </div>
          <button onClick={handleLogout} className="admin-logout-btn">Logout</button>
        </div>
      </header>

      {stats && (
        <div className="admin-stats-grid">
          <div className="stat-card">
            <div className="stat-icon">👥</div>
            <div className="stat-content">
              <div className="stat-value">{stats.total_users}</div>
              <div className="stat-label">Total Users</div>
            </div>
          </div>
          
          <div className="stat-card">
            <div className="stat-icon">📱</div>
            <div className="stat-content">
              <div className="stat-value">{stats.total_apps}</div>
              <div className="stat-label">Total Monitors</div>
            </div>
          </div>
          
          <div className="stat-card">
            <div className="stat-icon">💳</div>
            <div className="stat-content">
              <div className="stat-value">{stats.active_subscribers}</div>
              <div className="stat-label">Active Subscribers</div>
            </div>
          </div>
          
          <div className="stat-card">
            <div className="stat-icon">📈</div>
            <div className="stat-content">
              <div className="stat-value">{stats.total_status_checks.toLocaleString()}</div>
              <div className="stat-label">Status Checks</div>
            </div>
          </div>
        </div>
      )}

      {stats && (
        <div className="plan-breakdown">
          <h3>Plan Distribution</h3>
          <div className="plan-stats">
            <div className="plan-stat">
              <span className="plan-label">Free:</span>
              <span className="plan-count">{stats.free_users}</span>
            </div>
            <div className="plan-stat">
              <span className="plan-label">Pro:</span>
              <span className="plan-count">{stats.pro_users}</span>
            </div>
            <div className="plan-stat">
              <span className="plan-label">Business:</span>
              <span className="plan-count">{stats.business_users}</span>
            </div>
          </div>
        </div>
      )}

      <div className="users-section">
        <div className="users-header">
          <h2>Users ({filteredUsers.length})</h2>
          <div className="users-filters">
            <input
              type="text"
              placeholder="Search users..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input"
            />
            <select
              value={filterPlan}
              onChange={(e) => setFilterPlan(e.target.value)}
              className="filter-select"
            >
              <option value="all">All Plans</option>
              <option value="free">Free</option>
              <option value="pro">Pro</option>
              <option value="business">Business</option>
            </select>
          </div>
        </div>

        <div className="table-container">
          <table className="users-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>User</th>
                <th>Email</th>
                <th>Plan</th>
                <th>Plan Started</th>
                <th>Subscription Ends</th>
                <th>Monitors</th>
                <th>Total Checks</th>
                <th>Stripe Customer</th>
                <th>Joined</th>
              </tr>
            </thead>
            <tbody>
              {filteredUsers.map(user => (
                <tr key={user.id}>
                  <td>{user.id}</td>
                  <td>
                    <div className="user-cell">
                      {user.avatar_url && (
                        <img src={user.avatar_url} alt={user.username} className="user-avatar" />
                      )}
                      <span>{user.username}</span>
                    </div>
                  </td>
                  <td>{user.email}</td>
                  <td>
                    <span className={`plan-badge ${getPlanBadgeClass(user.plan)}`}>
                      {user.plan}
                    </span>
                  </td>
                  <td>{formatDate(user.plan_started_at)}</td>
                  <td>
                    {user.stripe_subscription_id ? (
                      <span className="subscription-active">
                        {getSubscriptionEndDate(user.plan_started_at)}
                      </span>
                    ) : (
                      <span className="subscription-inactive">N/A</span>
                    )}
                  </td>
                  <td className="text-center">{user.app_count}</td>
                  <td className="text-center">{user.total_checks.toLocaleString()}</td>
                  <td>
                    {user.stripe_customer_id ? (
                      <code className="stripe-id">{user.stripe_customer_id.substring(0, 20)}...</code>
                    ) : (
                      <span className="no-stripe">—</span>
                    )}
                  </td>
                  <td>{formatDate(user.created_at)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default Admin;
