import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import RetroTerminalOnboarding from "./Onboard";
import Home from "./Home";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/onboarding" element={<RetroTerminalOnboarding />} />
      </Routes>
    </Router>
  );
}

export default App;
