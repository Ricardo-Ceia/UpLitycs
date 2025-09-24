import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import RetroTerminalOnboarding from "./Onboard";
import Home from "./Home";
import Dashboard from "./Dashboard";
import RetroAuth from "./RetroAuth"

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<RetroAuth />} />
        <Route path="/onboarding" element={<RetroTerminalOnboarding />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </Router>
  );
}

export default App;
