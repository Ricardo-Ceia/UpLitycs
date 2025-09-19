import "./App.css";
import React, { useState, useEffect } from "react";

function Home() {
  return (
    <section className="homeBackground relative min-h-screen flex flex-col justify-center items-center px-6 text-center space-y-6 overflow-hidden">
        <h1 className="brand-logo">UPLYTICS</h1>
        <StatusLight/>
    <TypingText text="Monitor Your Services. Catch Issues Before They Catch You." speed={50} />
        <button className="retro-button mt-6" onClick={handleGetStarted}>
          Get Started - Free
        </button>
    </section>
  );
}

const handleGetStarted = () => {
  window.location.href = "/onboarding"
}

function StatusLight(){
  const [active, setActive] = useState("green");

  return (
      <div className="status-lights">
        {["red", "yellow", "green"].map(color => (
          <div
            key={color}
            className={`led ${color} ${active === color ? "active" : ""}`}
            onClick={() => setActive(color)}
          />
        ))}
      </div>
  )
}


function TypingText({ text, speed = 100 }) {
  const [displayedText, setDisplayedText] = useState("");
  const [index, setIndex] = useState(0);

  useEffect(() => {
    if (index < text.length) {
      const timeout = setTimeout(() => {
        setDisplayedText((prev) => prev + text[index]);
        setIndex(index + 1);
      }, speed);

      return () => clearTimeout(timeout);
    }
  }, [index, text, speed]);

  // NOTE: corrected class to match App.css .retro-typing
  return (
    <h2 className="retro-typing">
      {displayedText}
      <span className="vim-cursor">{index < text.length ? " " : ""}</span>
    </h2>
  );
}

export default Home;

