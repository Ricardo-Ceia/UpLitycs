import "./App.css";
import React, { useState } from "react";
function Home() {
  return (
    <section className="homeBackground relative min-h-screen flex flex-col justify-center items-center px-6 text-center space-y-6 overflow-hidden">
        <h1 className="brand-logo">UPLYTICS</h1>
        <StatusLight/>
        <ShinyText text="Monitor Your Services. Catch Issues Before They Catch You." />
        <p className="text-textSecondary text-lg md:text-xl max-w-xl mt-4">
          Uplytics keeps your websites, APIs, and apps running smoothly with real-time alerts and a public status page.
        </p>
        <button className="retro-button mt-6">
          Get Started - Free
        </button>
    </section>
  );
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
function ShinyText({ text }) {
  return (
    <h1 className="text-4xl md:text-5xl font-bold text-primary leading-tight">
      {text.split("").map((c, index) => (
        <span
          key={index}
          style={{ animationDelay: `${index * 0.1}s` }}
          className="shiny-letter"
        >
          {c}
        </span>
      ))}
    </h1>
  );
}

export default Home;

