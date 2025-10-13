import "./Home.css";
import React, { useState, useEffect } from "react";
import { Activity, Zap, Shield, Bell, Terminal, TrendingUp } from 'lucide-react';

function Home() {
  return (
    <div className="home-container">
      {/* CRT Effect Overlay */}
      <div className="crt-overlay"></div>
      <div className="scan-lines"></div>
      
      {/* Animated Background Grid */}
      <div className="grid-background"></div>
      
      {/* Floating Particles */}
      <div className="particles">
        {[...Array(15)].map((_, i) => (
          <div key={i} className="particle" style={{
            left: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 5}s`,
            animationDuration: `${5 + Math.random() * 10}s`
          }} />
        ))}
      </div>

      {/* Hero Section */}
      <section className="hero-section">
        <div className="hero-content">
          {/* Animated Logo */}
          <div className="logo-container">
            <h1 className="brand-logo-new">
              <span className="logo-bracket">{"<"}</span>
              STATUSFRAME
              <span className="logo-bracket">{"/>"}</span>
            </h1>
            <div className="logo-subtitle">
              <Terminal className="terminal-icon" />
              <span className="typing-text">STATUS MONITORING SYSTEM</span>
            </div>
          </div>

          {/* Status Lights */}
          <StatusLight />

          {/* Main Tagline */}
          <div className="tagline-container">
            <h2 className="main-tagline">
              Monitor Your Services.
              <br />
              <span className="tagline-highlight">Catch Issues Before They Catch You.</span>
            </h2>
          </div>

          {/* CTA Button */}
          <button className="cta-button-new" onClick={handleGetStarted}>
            <span className="button-glow"></span>
            <span className="button-text">
              <Zap className="button-icon" />
              START MONITORING - FREE
              <span className="button-arrow">→</span>
            </span>
          </button>

          {/* Trust Indicators */}
          <div className="trust-indicators">
            <div className="trust-item">
              <Activity className="trust-icon" />
              <span>1-100 Monitors</span>
            </div>
            <div className="trust-item">
              <Zap className="trust-icon" />
              <span>30s-5min Checks</span>
            </div>
            <div className="trust-item">
              <TrendingUp className="trust-icon" />
              <span>7-90 Days Retention</span>
            </div>
          </div>
        </div>

        {/* Video Demo Section */}
        <div className="video-demo-section">
          <h3 className="section-title">
            <span className="title-line">─────</span>
            SEE IT IN ACTION
            <span className="title-line">─────</span>
          </h3>
          <div className="video-container">
            <video 
              className="demo-video" 
              controls 
              preload="metadata"
              poster="/demo-poster.jpg"
              width="1920"
              height="1080"
            >
              <source src="/demo.mp4" type="video/mp4" />
              Your browser does not support the video tag.
            </video>
            <p className="video-subtext">Monitor setup in under 60 seconds</p>
          </div>
        </div>

        {/* Features Grid */}
        <div className="features-grid">
          <FeatureCard
            icon={<Shield />}
            title="SCALABLE MONITORING"
            description="1-100 monitors • Free, Pro & Business plans"
            status="FLEXIBLE"
          />
          <FeatureCard
            icon={<Zap />}
            title="SMART HEALTH CHECKS"
            description="5min to 30sec intervals • Plan-based frequency"
            status="ACTIVE"
          />
          <FeatureCard
            icon={<TrendingUp />}
            title="HISTORICAL DATA"
            description="7-90 days retention • Upgrade for more history"
            status="RUNNING"
          />
        </div>

        {/* How It Works Section */}
        <div className="how-it-works">
          <h3 className="section-title">
            <span className="title-line">─────</span>
            HOW IT WORKS
            <span className="title-line">─────</span>
          </h3>
          <div className="steps-container">
            <Step number="01" title="SIGN UP" description="Quick OAuth login with Google" />
            <Step number="02" title="CONFIGURE" description="Add your service health endpoint" />
            <Step number="03" title="MONITOR" description="Get real-time status updates" />
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="home-footer">
        <div className="footer-content">
          <div className="footer-text">
            © 2025 STATUSFRAME | Built for Developers
          </div>
          <div className="footer-center">
            <a 
              href="https://x.com/SalzDevs" 
              target="_blank" 
              rel="noopener noreferrer"
              className="footer-x-handle"
            >
              @SalzDevs
            </a>
          </div>
          <div className="footer-status">
            <span className="status-dot"></span>
            <span>All Systems Operational</span>
          </div>
        </div>
      </footer>
    </div>
  );
}

const handleGetStarted = () => {
  window.location.href = "/auth";
};

function StatusLight() {
  const [active, setActive] = useState("green");

  return (
    <div className="status-lights-new">
      {["red", "yellow", "green"].map(color => (
        <div
          key={color}
          className={`led-new ${color} ${active === color ? "active" : ""}`}
          onClick={() => setActive(color)}
        />
      ))}
    </div>
  );
}

function FeatureCard({ icon, title, description, status }) {
  return (
    <div className="feature-card">
      <div className="feature-header">
        <div className="feature-icon">{icon}</div>
        <div className="feature-status">{status}</div>
      </div>
      <h3 className="feature-title">{title}</h3>
      <p className="feature-description">{description}</p>
      <div className="feature-footer">
        <span className="feature-line"></span>
      </div>
    </div>
  );
}

function Step({ number, title, description }) {
  return (
    <div className="step-card">
      <div className="step-number">{number}</div>
      <h4 className="step-title">{title}</h4>
      <p className="step-description">{description}</p>
    </div>
  );
}

export default Home;