import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import { useState, useEffect } from "react";
import RetroTerminalOnboarding from "./Onboard";
import Home from "./Home";
import StatusPage from "./StatusPage";
import RetroAuth from "./RetroAuth";
import ProtectedRoute from "./ProtectedRoute";

// Dashboard redirect component
function DashboardRedirect() {
  const [slug, setSlug] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUserSlug = async () => {
      try {
        const response = await fetch('/api/user-status', {
          credentials: 'include'
        });
        if (response.ok) {
          const data = await response.json();
          setSlug(data.slug);
        }
      } catch (error) {
        console.error('Error fetching user slug:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchUserSlug();
  }, []);

  if (loading) {
    return <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh', color: '#00ff41'}}>Loading...</div>;
  }

  if (slug) {
    return <Navigate to={`/status/${slug}`} replace />;
  }

  return <Navigate to="/onboarding" replace />;
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/auth" element={<RetroAuth />} />
        {/* Single unified status page - public for everyone, owners can control theme */}
        <Route path="/status/:slug" element={<StatusPage />} />
        <Route path="/onboarding" element={
          <ProtectedRoute>
            <RetroTerminalOnboarding />
          </ProtectedRoute>
        } />
        {/* Redirect /dashboard to user's status page */}
        <Route path="/dashboard" element={
          <ProtectedRoute>
            <DashboardRedirect />
          </ProtectedRoute>
        } />
      </Routes>
    </Router>
  );
}

export default App;