import React, { useState, useEffect, useRef } from 'react';
import { ChevronRight, Monitor, Zap, Sparkles, Mail, Check, Copy, Globe } from 'lucide-react';

const RetroOnboarding = () => {
  const [currentStep, setCurrentStep] = useState(0);
  const [selectedLanguage, setSelectedLanguage] = useState('javascript');
  const [selectedTheme, setSelectedTheme] = useState('');
  const [emailAlerts, setEmailAlerts] = useState('');
  const [homepageUrl, setHomepageUrl] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [typewriterText, setTypewriterText] = useState('');
  const [copied, setCopied] = useState(false);
  
  const codeExamples = {
    javascript: `// Health Check Endpoint - Node.js/Express
app.get('/health', (req, res) => {
  const health = {
    status: 'UP',
    timestamp: new Date().toISOString(),
    service: 'my-service',
    version: '1.0.0',
    uptime: process.uptime(),
    memory: process.memoryUsage(),
    environment: process.env.NODE_ENV
  };
  res.status(200).json(health);
});`,
    python: `# Health Check Endpoint - Python/Flask
@app.route('/health')
def health_check():
    health = {
        'status': 'UP',
        'timestamp': datetime.now().isoformat(),
        'service': 'my-service',
        'version': '1.0.0',
        'uptime': time.time() - start_time,
        'memory': psutil.virtual_memory().percent,
        'environment': os.environ.get('ENV')
    }
    return jsonify(health), 200`,
    go: `// Health Check Endpoint - Go/Gin
func HealthCheck(c *gin.Context) {
    health := gin.H{
        "status":      "UP",
        "timestamp":   time.Now().Format(time.RFC3339),
        "service":     "my-service",
        "version":     "1.0.0",
        "uptime":      time.Since(startTime).Seconds(),
        "memory":      getMemoryUsage(),
        "environment": os.Getenv("ENV"),
    }
    c.JSON(200, health)
}`,
    ruby: `# Health Check Endpoint - Ruby/Sinatra
get '/health' do
  content_type :json
  {
    status: 'UP',
    timestamp: Time.now.iso8601,
    service: 'my-service',
    version: '1.0.0',
    uptime: Time.now - START_TIME,
    memory: get_memory_usage,
    environment: ENV['RACK_ENV']
  }.to_json
end`
  };

  const themes = [
    {
      id: 'cyberpunk',
      name: 'Cyberpunk Theme',
      icon: <Zap className="w-8 h-8" />,
      preview: 'bg-gradient-to-br from-purple-900 via-pink-900 to-indigo-900',
      description: 'Synthwave vibes with neon colors',
      colors: ['#FF6EC7', '#8C52FF', '#00FFF7']
    },
    {
      id: 'matrix',
      name: 'Matrix Theme',
      icon: <Monitor className="w-8 h-8" />,
      preview: 'bg-gradient-to-br from-green-900 via-green-800 to-black',
      description: 'Green terminal matrix aesthetic',
      colors: ['#00ff41', '#008f11', '#003300']
    },
    {
      id: 'retro',
      name: 'Retro Theme',
      icon: <Monitor className="w-8 h-8" />,
      preview: 'bg-gradient-to-br from-orange-900 via-yellow-800 to-red-900',
      description: 'Vintage 80s computer vibes',
      colors: ['#ff6b35', '#f7931e', '#fdc500']
    },
    {
      id: 'minimal',
      name: 'Minimal Theme',
      icon: <Sparkles className="w-8 h-8" />,
      preview: 'bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900',
      description: 'Clean and minimalist design',
      colors: ['#ffffff', '#a0a0a0', '#606060']
    }
  ];

  useEffect(() => {
    const text = "INITIALIZING ONBOARDING SEQUENCE...";
    let index = 0;
    const timer = setInterval(() => {
      if (index <= text.length) {
        setTypewriterText(text.slice(0, index));
        index++;
      } else {
        clearInterval(timer);
      }
    }, 50);
    return () => clearInterval(timer);
  }, []);

  const copyCode = () => {
    navigator.clipboard.writeText(codeExamples[selectedLanguage]);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleSubmit = async () => {
    if (selectedTheme && emailAlerts && homepageUrl) {
      setIsSubmitting(true);
      
      // Prepare data for submission
      const onboardingData = {
        name: "User",
        homepage: homepageUrl,
        alerts: emailAlerts === 'yes' ? 'y' : 'n',
        theme: selectedTheme
      };
      
      try {
        const response = await fetch('/api/go-to-dashboard', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include',
          body: JSON.stringify(onboardingData)
        });
        
        if (response.ok) {
          setTimeout(() => {
            window.location.href = "/dashboard";
          }, 1500);
        } else {
          console.error('Onboarding failed');
          setIsSubmitting(false);
        }
      } catch (error) {
        console.error('Error:', error);
        setIsSubmitting(false);
      }
    }
  };

  const renderStep = () => {
    switch(currentStep) {
      case 0:
        return (
          <div className="space-y-6 animate-fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-cyan-400 mb-2" style={{fontFamily: "'Press Start 2P', monospace"}}>
                STEP 1: HEALTH ENDPOINT
              </h2>
              <p className="text-purple-400 text-sm">Configure your service health check</p>
            </div>

            <div className="bg-black/50 rounded-lg p-6 border border-purple-500/50">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                  <Globe className="w-5 h-5 text-cyan-400" />
                  <span className="text-cyan-400 text-sm font-mono">SELECT LANGUAGE</span>
                </div>
                <button
                  onClick={copyCode}
                  className="flex items-center gap-2 px-3 py-1 bg-purple-600/20 border border-purple-500 rounded text-xs text-purple-300 hover:bg-purple-600/30 transition-all"
                >
                  {copied ? <Check className="w-3 h-3" /> : <Copy className="w-3 h-3" />}
                  {copied ? 'COPIED!' : 'COPY'}
                </button>
              </div>

              <div className="flex gap-2 mb-4">
                {Object.keys(codeExamples).map((lang) => (
                  <button
                    key={lang}
                    onClick={() => setSelectedLanguage(lang)}
                    className={`px-4 py-2 rounded font-mono text-xs transition-all ${
                      selectedLanguage === lang
                        ? 'bg-cyan-500/30 border-2 border-cyan-400 text-cyan-300'
                        : 'bg-gray-800/50 border border-gray-600 text-gray-400 hover:border-purple-500'
                    }`}
                  >
                    {lang.toUpperCase()}
                  </button>
                ))}
              </div>

              <div className="bg-gray-950 rounded p-4 overflow-x-auto border border-gray-800">
                <pre className="text-xs font-mono">
                  <code className="text-green-400">{codeExamples[selectedLanguage]}</code>
                </pre>
              </div>

              <div className="mt-4 p-3 bg-gradient-to-r from-cyan-900/20 to-purple-900/20 rounded border border-cyan-500/30">
                <p className="text-xs text-cyan-300">
                  💡 <strong>TIP:</strong> This endpoint will be monitored every 30 seconds to track your service status
                </p>
              </div>
            </div>
          </div>
        );

      case 1:
        return (
          <div className="space-y-6 animate-fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-cyan-400 mb-2" style={{fontFamily: "'Press Start 2P', monospace"}}>
                STEP 2: ENTER YOUR SERVICE URL
              </h2>
              <p className="text-purple-400 text-sm">Provide the health check endpoint URL to monitor</p>
            </div>

            <div className="bg-black/50 rounded-lg p-8 border border-purple-500/50">
              <div className="flex justify-center mb-6">
                <div className="p-6 bg-gradient-to-br from-purple-600/20 to-cyan-600/20 rounded-full">
                  <Globe className="w-16 h-16 text-cyan-400" />
                </div>
              </div>

              <div className="space-y-4">
                <label className="block">
                  <span className="text-cyan-300 text-sm font-mono mb-2 block">HEALTH CHECK URL</span>
                  <input
                    type="url"
                    value={homepageUrl}
                    onChange={(e) => setHomepageUrl(e.target.value)}
                    placeholder="https://your-service.com/health"
                    className="w-full px-4 py-3 bg-gray-950 border-2 border-purple-500/50 rounded-lg text-cyan-300 font-mono focus:border-cyan-400 focus:outline-none focus:ring-2 focus:ring-cyan-400/20 transition-all"
                  />
                </label>

                <div className="p-4 bg-gradient-to-r from-cyan-900/20 to-purple-900/20 rounded border border-cyan-500/30">
                  <p className="text-xs text-cyan-300 mb-2">
                    <strong>Examples:</strong>
                  </p>
                  <ul className="text-xs text-gray-400 space-y-1 font-mono">
                    <li>• https://api.myapp.com/health</li>
                    <li>• https://myservice.herokuapp.com/status</li>
                    <li>• https://example.com/api/v1/health</li>
                  </ul>
                </div>

                {homepageUrl && (
                  <div className="p-3 bg-green-900/20 rounded border border-green-500/50">
                    <p className="text-xs text-green-400">
                      ✓ URL looks good! We'll check this endpoint every 30 seconds.
                    </p>
                  </div>
                )}
              </div>
            </div>
          </div>
        );

      case 2:
        return (
          <div className="space-y-6 animate-fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-cyan-400 mb-2" style={{fontFamily: "'Press Start 2P', monospace"}}>
                STEP 3: CHOOSE YOUR THEME
              </h2>
              <p className="text-purple-400 text-sm">Select a visual style for your status page</p>
            </div>

            <div className="grid gap-4">
              {themes.map((theme) => (
                <button
                  key={theme.id}
                  onClick={() => setSelectedTheme(theme.id)}
                  className={`group relative p-6 rounded-lg border-2 transition-all transform hover:scale-[1.02] ${
                    selectedTheme === theme.id
                      ? 'border-cyan-400 bg-cyan-400/10 shadow-lg shadow-cyan-400/30'
                      : 'border-purple-500/50 bg-black/30 hover:border-purple-400'
                  }`}
                >
                  <div className="flex items-center gap-4">
                    <div className={`p-4 rounded-lg ${theme.preview}`}>
                      <div className="text-white">{theme.icon}</div>
                    </div>
                    
                    <div className="flex-1 text-left">
                      <h3 className="text-lg font-bold text-cyan-300 mb-1">{theme.name}</h3>
                      <p className="text-sm text-gray-400">{theme.description}</p>
                      
                      <div className="flex gap-2 mt-3">
                        {theme.colors.map((color, i) => (
                          <div
                            key={i}
                            className="w-6 h-6 rounded-full border border-white/20"
                            style={{ backgroundColor: color, boxShadow: `0 0 10px ${color}50` }}
                          />
                        ))}
                      </div>
                    </div>

                    {selectedTheme === theme.id && (
                      <div className="absolute top-2 right-2">
                        <Check className="w-6 h-6 text-cyan-400" />
                      </div>
                    )}
                  </div>
                </button>
              ))}
            </div>
          </div>
        );

      case 3:
        return (
          <div className="space-y-6 animate-fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-cyan-400 mb-2" style={{fontFamily: "'Press Start 2P', monospace"}}>
                STEP 4: EMAIL ALERTS
              </h2>
              <p className="text-purple-400 text-sm">Get notified when your service goes down</p>
            </div>

            <div className="bg-black/50 rounded-lg p-8 border border-purple-500/50">
              <div className="flex justify-center mb-6">
                <div className="p-6 bg-gradient-to-br from-purple-600/20 to-cyan-600/20 rounded-full">
                  <Mail className="w-16 h-16 text-cyan-400" />
                </div>
              </div>

              <h3 className="text-center text-xl text-cyan-300 mb-6">
                Would you like to receive email alerts?
              </h3>

              <div className="grid grid-cols-2 gap-4">
                <button
                  onClick={() => setEmailAlerts('yes')}
                  className={`p-6 rounded-lg border-2 transition-all ${
                    emailAlerts === 'yes'
                      ? 'border-green-400 bg-green-400/10 shadow-lg shadow-green-400/30'
                      : 'border-gray-600 bg-gray-800/50 hover:border-green-500'
                  }`}
                >
                  <Check className="w-8 h-8 mx-auto mb-2 text-green-400" />
                  <span className="text-green-400 font-bold">YES</span>
                  <p className="text-xs text-gray-400 mt-2">Get instant notifications</p>
                </button>

                <button
                  onClick={() => setEmailAlerts('no')}
                  className={`p-6 rounded-lg border-2 transition-all ${
                    emailAlerts === 'no'
                      ? 'border-red-400 bg-red-400/10 shadow-lg shadow-red-400/30'
                      : 'border-gray-600 bg-gray-800/50 hover:border-red-500'
                  }`}
                >
                  <span className="block w-8 h-8 mx-auto mb-2 text-red-400 text-2xl font-bold">✕</span>
                  <span className="text-red-400 font-bold">NO</span>
                  <p className="text-xs text-gray-400 mt-2">Check dashboard manually</p>
                </button>
              </div>

              {emailAlerts === 'yes' && (
                <div className="mt-6 p-4 bg-gradient-to-r from-green-900/20 to-cyan-900/20 rounded border border-green-500/30">
                  <p className="text-sm text-green-300">
                    ✅ You'll receive alerts at your registered email when your service status changes
                  </p>
                </div>
              )}
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900 flex items-center justify-center p-4">
      {/* Animated background effects */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
    <div
        className="absolute inset-0 bg-[url('data:image/svg+xml,%3Csvg width=\'60\' height=\'60\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cdefs%3E%3Cpattern id=\'grid\' width=\'60\' height=\'60\' patternUnits=\'userSpaceOnUse\'%3E%3Cpath d=\'M 60 0 L 0 0 0 60\' fill=\'none\' stroke=\'rgba(139,92,246,0.1)\' stroke-width=\'1\'/%3E%3C/pattern%3E%3C/defs%3E%3Crect width=\'100%25\' height=\'100%25\' fill=\'url(%23grid)\'/%3E%3C/svg%3E')] opacity-50"
      />
      <div className="absolute top-0 left-1/4 w-96 h-96 bg-purple-500 rounded-full filter blur-[128px] opacity-20 animate-pulse" />
      <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-cyan-500 rounded-full filter blur-[128px] opacity-20 animate-pulse" />
    </div>

      <div className="relative z-10 w-full max-w-3xl">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold mb-2" style={{
            fontFamily: "'Press Start 2P', monospace",
            background: 'linear-gradient(45deg, #FF6EC7, #8C52FF, #00FFF7)',
            backgroundClip: 'text',
            WebkitBackgroundClip: 'text',
            color: 'transparent',
            textShadow: '0 0 30px rgba(255,110,199,0.5)'
          }}>
            UPLYTICS
          </h1>
          <div className="text-cyan-400 text-sm font-mono opacity-80">
            {typewriterText}
            <span className="inline-block w-2 h-4 bg-cyan-400 ml-1 animate-pulse" />
          </div>
        </div>

        {/* Main Container */}
        <div className="bg-black/40 backdrop-blur-md rounded-2xl border-2 border-purple-500/50 shadow-2xl shadow-purple-500/20 p-8">
          {/* Progress indicator */}
          <div className="flex items-center justify-between mb-8">
            {[0, 1, 2, 3].map((step) => (
              <div key={step} className="flex items-center flex-1">
                <div className={`w-10 h-10 rounded-full border-2 flex items-center justify-center font-mono transition-all ${
                  currentStep >= step 
                    ? 'bg-cyan-500/30 border-cyan-400 text-cyan-300 shadow-lg shadow-cyan-400/50' 
                    : 'border-gray-600 text-gray-500'
                }`}>
                  {currentStep > step ? '✓' : step + 1}
                </div>
                {step < 3 && (
                  <div className={`flex-1 h-0.5 mx-2 transition-all ${
                    currentStep > step ? 'bg-cyan-400' : 'bg-gray-700'
                  }`} />
                )}
              </div>
            ))}
          </div>

          {/* Step Content */}
          {renderStep()}

          {/* Navigation */}
          <div className="flex justify-between mt-8">
            <button
              onClick={() => setCurrentStep(Math.max(0, currentStep - 1))}
              disabled={currentStep === 0}
              className={`px-6 py-3 rounded-lg font-mono text-sm transition-all flex items-center gap-2 ${
                currentStep === 0
                  ? 'bg-gray-800 text-gray-600 cursor-not-allowed'
                  : 'bg-purple-600/30 border border-purple-500 text-purple-300 hover:bg-purple-600/50'
              }`}
            >
              ← BACK
            </button>

            {currentStep < 3 ? (
              <button
                onClick={() => setCurrentStep(currentStep + 1)}
                disabled={
                  (currentStep === 1 && !homepageUrl) ||
                  (currentStep === 2 && !selectedTheme)
                }
                className="px-6 py-3 rounded-lg font-mono text-sm bg-gradient-to-r from-cyan-500 to-purple-500 text-white hover:from-cyan-400 hover:to-purple-400 transition-all flex items-center gap-2 shadow-lg hover:shadow-cyan-400/50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                NEXT <ChevronRight className="w-4 h-4" />
              </button>
            ) : (
              <button
                onClick={handleSubmit}
                disabled={!selectedTheme || !emailAlerts || !homepageUrl || isSubmitting}
                className="px-8 py-3 rounded-lg font-mono text-sm bg-gradient-to-r from-green-500 to-cyan-500 text-white hover:from-green-400 hover:to-cyan-400 transition-all shadow-lg hover:shadow-green-400/50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? (
                  <span className="flex items-center gap-2">
                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    LAUNCHING...
                  </span>
                ) : (
                  'COMPLETE SETUP →'
                )}
              </button>
            )}
          </div>
        </div>

        {/* Status Bar */}
        <div className="mt-6 text-center">
          <div className="inline-flex items-center gap-3 px-4 py-2 bg-black/40 backdrop-blur-md rounded-full border border-purple-500/30">
            <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse shadow-lg shadow-green-400/50" />
            <span className="text-xs font-mono text-gray-400">SYSTEM READY</span>
            <span className="text-xs font-mono text-cyan-400">STEP {currentStep + 1} OF 4</span>
          </div>
        </div>
      </div>

      <style jsx>{`
        @import url('https://fonts.googleapis.com/css2?family=Press+Start+2P&display=swap');
        
        @keyframes fadeIn {
          from { opacity: 0; transform: translateY(10px); }
          to { opacity: 1; transform: translateY(0); }
        }
        
        .animate-fadeIn {
          animation: fadeIn 0.5s ease-out;
        }
      `}</style>
    </div>
  );
};

export default RetroOnboarding;