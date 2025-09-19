import React, { useState, useEffect, useRef } from 'react';
import './RetroTerminalOnboarding.css';

const TerminalOnboarding = () => {
  const steps = [
    {
      question: "$ uplytics init\n\nInitializing Uplytics setup...\n\nWhat's your name?",
      type: "text",
      key: "name",
      validation: (value) => value.trim().length > 0 ? null : "Name cannot be empty",
      placeholder: "Enter your name"
    },
    {
      question: "What's the URL of your Homepage?",
      type: "url",
      key: "homepage",
      validation: (value) => {
        const urlPattern = /^https?:\/\/.+\..+/;
        return urlPattern.test(value) ? null : "Please enter a valid URL (e.g., https://example.com)";
      },
      placeholder: "https://example.com"
    },
    {
      question: "Do you want to receive alerts? (y/n)",
      type: "boolean",
      key: "alerts",
      validation: (value) => {
        const normalized = value.toLowerCase().trim();
        return ['y', 'yes', 'n', 'no'].includes(normalized) ? null : "Please enter y/yes or n/no";
      },
      placeholder: "y/n"
    },
    {
      question: "Configuration complete!\n\nPress ENTER to finish the onboarding...",
      type: "complete",
      key: "finish"
    }
  ];

  const [stepIndex, setStepIndex] = useState(0);
  const [displayedText, setDisplayedText] = useState("");
  const [inputValue, setInputValue] = useState("");
  const [answers, setAnswers] = useState({});
  const [history, setHistory] = useState([]);
  const [error, setError] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const [showInput, setShowInput] = useState(false);
  const inputRef = useRef(null);
  const terminalRef = useRef(null);

  const scrollToBottom = () => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  };

  useEffect(() => {
    scrollToBottom();
  }, [history, displayedText]);

  useEffect(() => {
    setDisplayedText("");
    setShowInput(false);
    setError("");
    setIsTyping(false); // No typing animation
    
    const question = steps[stepIndex]?.question ?? "";
    
    // Show text instantly
    setDisplayedText(question);
    setShowInput(true);
    
    // Focus input immediately
    setTimeout(() => inputRef.current?.focus(), 0);
  }, [stepIndex]);

  const validateInput = (value) => {
    const step = steps[stepIndex];
    if (step.validation) {
      return step.validation(value);
    }
    return null;
  };

  const formatAnswer = (key, value) => {
    if (key === 'alerts') {
      const normalized = value.toLowerCase().trim();
      return ['y', 'yes'].includes(normalized) ? 'enabled' : 'disabled';
    }
    return value;
  };

  const handleEnter = (e) => {
    if (e.key === "Enter" && !isTyping) {
      const currentStep = steps[stepIndex];
      const trimmedValue = inputValue.trim();
      
      if (currentStep.type === "complete") {
        // Finish onboarding
        setHistory(prev => [
          ...prev,
          `> ${inputValue}`,
          "Launching Uplytics dashboard...",
          "Setup complete! 🚀",
          "",
          "Redirecting..."
        ]);
        console.log("Onboarding completed with data:", answers);
        setTimeout(() => {
          window.location.href = "/dashboard";
        }, 2000);
        return;
      }

      const validationError = validateInput(trimmedValue);
      
      if (validationError) {
        setError(validationError);
        return;
      }

      setError("");
      const formattedAnswer = formatAnswer(currentStep.key, trimmedValue);
      
      setAnswers((prev) => ({ ...prev, [currentStep.key]: trimmedValue }));
      
      // Add current question and answer to history
      setHistory(prev => [
        ...prev,
        steps[stepIndex]?.question || "",
        `> ${trimmedValue}`,
        `✓ ${currentStep.key}: ${formattedAnswer}`,
        ""
      ]);
      
      setInputValue("");
      
      if (stepIndex < steps.length - 1) {
        setStepIndex(stepIndex + 1); // Immediate step change
      }
    }
  };

  const getPromptSymbol = () => {
    if (steps[stepIndex]?.type === "complete") return "$ ";
    return "? ";
  };

  return (
    <div className="terminal-container">
      <div className="terminal-window">
        {/* Terminal Header */}
        <div className="terminal-header">
          <div className="terminal-controls">
            <div className="terminal-control close"></div>
            <div className="terminal-control minimize"></div>
            <div className="terminal-control maximize"></div>
          </div>
          <span className="terminal-title">uplytics@setup:~</span>
          <div className="terminal-step-indicator">
            Step {stepIndex + 1}/{steps.length}
          </div>
        </div>

        {/* Terminal Content */}
        <div ref={terminalRef} className="terminal-content">
          {/* History */}
          {history.map((line, index) => (
            <div 
              key={index} 
              className={`terminal-history-line ${
                line.startsWith('>') ? 'input' : 
                line.startsWith('✓') ? 'success' : 'normal'
              }`}
            >
              {line}
            </div>
          ))}

          {/* Current Question */}
          <div className="terminal-question" aria-live="polite">
            {displayedText}
          </div>

          {/* Error Message */}
          {error && (
            <div className="terminal-error">
              ⚠ {error}
            </div>
          )}

          {/* Input Line */}
          {showInput && (
            <div className="terminal-input-line">
              <span className="terminal-prompt">
                {getPromptSymbol()}
              </span>
              <div className="terminal-input-container">
                <input
                  ref={inputRef}
                  className="terminal-input"
                  type="text"
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  onKeyDown={handleEnter}
                  placeholder={steps[stepIndex]?.placeholder || ""}
                  autoComplete="off"
                />
              </div>
            </div>
          )}
        </div>

        {/* Progress Bar */}
        <div className="terminal-progress-bar">
          <div 
            className="terminal-progress-fill"
            style={{ width: `${((stepIndex + 1) / steps.length) * 100}%` }}
          ></div>
        </div>
      </div>
    </div>
  );
};

export default TerminalOnboarding;
