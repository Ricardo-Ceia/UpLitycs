import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import RetroTerminalOnboarding from "./Onboard";
import Home from "./Home";
import Dashboard from "./Dashboard";
import RetroAuth from "./RetroAuth";
import ProtectedRoute from "./ProtectedRoute";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<RetroAuth />} />
        <Route path="/onboarding" element={
          <ProtectedRoute>
            <RetroTerminalOnboarding />
          </ProtectedRoute>
        } />
        <Route path="/dashboard" element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        } />
      </Routes>
    </Router>
  );
}

export default App;
