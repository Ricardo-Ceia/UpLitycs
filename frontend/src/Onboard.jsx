import React, { useState, useEffect, useRef } from 'react';
import { ChevronRight, Monitor, Zap, Sparkles, Mail, Check, Copy, Globe } from 'lucide-react';

const RetroOnboarding = () => {
  const [currentStep, setCurrentStep] = useState(0);
  const [selectedLanguage, setSelectedLanguage] = useState('node-express');
  const [selectedTheme, setSelectedTheme] = useState('');
  const [emailAlerts, setEmailAlerts] = useState('');
  const [homepageUrl, setHomepageUrl] = useState('');
  const [appName, setAppName] = useState('');
  const [slug, setSlug] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [typewriterText, setTypewriterText] = useState('');
  const [copied, setCopied] = useState(false);
  
  const codeExamples = {
    'node-express': {
      name: 'Node.js + Express',
      description: 'Add this to your Express app',
      code: `// Health Check Endpoint - Node.js/Express
// 1. Install dependencies: npm install express
// 2. Add this route to your server file (e.g., app.js or server.js)

const express = require('express');
const app = express();

// Track when the server started
const startTime = Date.now();

// Health check endpoint
app.get('/health', (req, res) => {
  const healthData = {
    status: 'UP',
    timestamp: new Date().toISOString(),
    service: 'my-service',
    version: '1.0.0',
    uptime: Math.floor((Date.now() - startTime) / 1000), // seconds
    memory: {
      used: Math.round(process.memoryUsage().heapUsed / 1024 / 1024), // MB
      total: Math.round(process.memoryUsage().heapTotal / 1024 / 1024) // MB
    },
    environment: process.env.NODE_ENV || 'development'
  };
  
  res.status(200).json(healthData);
});

// Start your server
app.listen(3000, () => console.log('Server running on port 3000'));`,
      setup: [
        'Install Express: npm install express',
        'Create or open your server file (app.js or server.js)',
        'Copy the health check route into your file',
        'Make sure the endpoint is accessible at /health',
        'Test: curl http://localhost:3000/health'
      ]
    },
    'node-nextjs': {
      name: 'Next.js API Route',
      description: 'Create this in your Next.js app',
      code: `// Health Check - Next.js API Route
// 1. Create file: pages/api/health.js (Pages Router)
//    OR: app/api/health/route.js (App Router - see below)

// ----- For Pages Router (pages/api/health.js) -----
const startTime = Date.now();

export default function handler(req, res) {
  const healthData = {
    status: 'UP',
    timestamp: new Date().toISOString(),
    service: 'my-nextjs-app',
    version: '1.0.0',
    uptime: Math.floor((Date.now() - startTime) / 1000),
    environment: process.env.NODE_ENV
  };
  
  res.status(200).json(healthData);
}

// ----- For App Router (app/api/health/route.js) -----
const startTime = Date.now();

export async function GET(request) {
  const healthData = {
    status: 'UP',
    timestamp: new Date().toISOString(),
    service: 'my-nextjs-app',
    version: '1.0.0',
    uptime: Math.floor((Date.now() - startTime) / 1000),
    environment: process.env.NODE_ENV
  };
  
  return Response.json(healthData);
}`,
      setup: [
        'Pages Router: Create pages/api/health.js',
        'App Router: Create app/api/health/route.js',
        'Copy the appropriate code for your router type',
        'Deploy your app',
        'Test: https://yourdomain.com/api/health'
      ]
    },
    'python-flask': {
      name: 'Python + Flask',
      description: 'Add this to your Flask app',
      code: `# Health Check - Python/Flask
# 1. Install: pip install flask psutil
# 2. Add this to your Flask app (e.g., app.py)

from flask import Flask, jsonify
from datetime import datetime
import time
import psutil
import os

app = Flask(__name__)

# Track when the server started
start_time = time.time()

@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint for monitoring"""
    current_time = time.time()
    uptime_seconds = int(current_time - start_time)
    
    # Get memory usage
    process = psutil.Process()
    memory_info = process.memory_info()
    
    health_data = {
        'status': 'UP',
        'timestamp': datetime.now().isoformat(),
        'service': 'my-flask-app',
        'version': '1.0.0',
        'uptime': uptime_seconds,
        'memory': {
            'used_mb': round(memory_info.rss / 1024 / 1024, 2),
            'percent': psutil.virtual_memory().percent
        },
        'environment': os.environ.get('FLASK_ENV', 'production')
    }
    
    return jsonify(health_data), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)`,
      setup: [
        'Install dependencies: pip install flask psutil',
        'Create or open your Flask app file',
        'Add the health check route',
        'Run your app: python app.py',
        'Test: curl http://localhost:5000/health'
      ]
    },
    'python-fastapi': {
      name: 'Python + FastAPI',
      description: 'Add this to your FastAPI app',
      code: `# Health Check - Python/FastAPI
# 1. Install: pip install fastapi uvicorn psutil
# 2. Add this to your main.py

from fastapi import FastAPI
from datetime import datetime
from pydantic import BaseModel
import time
import psutil
import os

app = FastAPI()

# Track when the server started
start_time = time.time()

class HealthResponse(BaseModel):
    status: str
    timestamp: str
    service: str
    version: str
    uptime: int
    memory: dict
    environment: str

@app.get("/health", response_model=HealthResponse)
async def health_check():
    """
    Health check endpoint for monitoring service status
    Returns service health metrics including uptime and memory
    """
    current_time = time.time()
    uptime_seconds = int(current_time - start_time)
    
    # Get memory usage
    process = psutil.Process()
    memory_info = process.memory_info()
    
    return {
        "status": "UP",
        "timestamp": datetime.now().isoformat(),
        "service": "my-fastapi-app",
        "version": "1.0.0",
        "uptime": uptime_seconds,
        "memory": {
            "used_mb": round(memory_info.rss / 1024 / 1024, 2),
            "percent": psutil.virtual_memory().percent
        },
        "environment": os.environ.get("ENV", "production")
    }

# Run with: uvicorn main:app --host 0.0.0.0 --port 8000`,
      setup: [
        'Install: pip install fastapi uvicorn psutil',
        'Create or open your main.py',
        'Add the health check endpoint',
        'Run: uvicorn main:app --reload',
        'Test: curl http://localhost:8000/health'
      ]
    },
    'go-gin': {
      name: 'Go + Gin',
      description: 'Add this to your Gin application',
      code: `// Health Check - Go/Gin Framework
// 1. Install: go get -u github.com/gin-gonic/gin
// 2. Add this to your main.go

package main

import (
    "github.com/gin-gonic/gin"
    "runtime"
    "time"
    "os"
)

var startTime = time.Now()

type HealthResponse struct {
    Status      string            \`json:"status"\`
    Timestamp   string            \`json:"timestamp"\`
    Service     string            \`json:"service"\`
    Version     string            \`json:"version"\`
    Uptime      int64             \`json:"uptime"\`
    Memory      map[string]uint64 \`json:"memory"\`
    Environment string            \`json:"environment"\`
}

func HealthCheck(c *gin.Context) {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    
    health := HealthResponse{
        Status:    "UP",
        Timestamp: time.Now().Format(time.RFC3339),
        Service:   "my-go-service",
        Version:   "1.0.0",
        Uptime:    int64(time.Since(startTime).Seconds()),
        Memory: map[string]uint64{
            "alloc_mb":    memStats.Alloc / 1024 / 1024,
            "total_mb":    memStats.TotalAlloc / 1024 / 1024,
            "sys_mb":      memStats.Sys / 1024 / 1024,
        },
        Environment: getEnv("ENV", "production"),
    }
    
    c.JSON(200, health)
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

func main() {
    r := gin.Default()
    r.GET("/health", HealthCheck)
    r.Run(":8080") // listen on port 8080
}`,
      setup: [
        'Install Gin: go get -u github.com/gin-gonic/gin',
        'Create or open your main.go',
        'Add the health check handler',
        'Run: go run main.go',
        'Test: curl http://localhost:8080/health'
      ]
    },
    'go-standard': {
      name: 'Go Standard Library',
      description: 'Using only Go standard library',
      code: `// Health Check - Go Standard Library
// No external dependencies required!
// Add this to your main.go

package main

import (
    "encoding/json"
    "net/http"
    "runtime"
    "time"
    "os"
)

var startTime = time.Now()

type HealthResponse struct {
    Status      string            \`json:"status"\`
    Timestamp   string            \`json:"timestamp"\`
    Service     string            \`json:"service"\`
    Version     string            \`json:"version"\`
    Uptime      int64             \`json:"uptime"\`
    Memory      map[string]uint64 \`json:"memory"\`
    Environment string            \`json:"environment"\`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    // Only accept GET requests
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    
    health := HealthResponse{
        Status:    "UP",
        Timestamp: time.Now().Format(time.RFC3339),
        Service:   "my-go-service",
        Version:   "1.0.0",
        Uptime:    int64(time.Since(startTime).Seconds()),
        Memory: map[string]uint64{
            "alloc_mb": memStats.Alloc / 1024 / 1024,
            "total_mb": memStats.TotalAlloc / 1024 / 1024,
            "sys_mb":   memStats.Sys / 1024 / 1024,
        },
        Environment: getEnv("ENV", "production"),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

func main() {
    http.HandleFunc("/health", healthHandler)
    http.ListenAndServe(":8080", nil)
}`,
      setup: [
        'No dependencies needed - pure Go!',
        'Create or open your main.go',
        'Add the health check handler',
        'Run: go run main.go',
        'Test: curl http://localhost:8080/health'
      ]
    },
    'ruby-rails': {
      name: 'Ruby on Rails',
      description: 'Add this to your Rails app',
      code: `# Health Check - Ruby on Rails
# 1. Create a controller: rails generate controller Health
# 2. Add this to app/controllers/health_controller.rb

class HealthController < ApplicationController
  # Skip authentication for health checks
  skip_before_action :verify_authenticity_token
  skip_before_action :authenticate_user!, if: :devise_installed?
  
  # Track when the server started
  @@start_time = Time.now
  
  def check
    health_data = {
      status: 'UP',
      timestamp: Time.now.iso8601,
      service: 'my-rails-app',
      version: '1.0.0',
      uptime: (Time.now - @@start_time).to_i,
      memory: {
        used_mb: (\`ps -o rss= -p #{Process.pid}\`.to_i / 1024.0).round(2)
      },
      environment: Rails.env,
      rails_version: Rails::VERSION::STRING
    }
    
    render json: health_data, status: :ok
  end
  
  private
  
  def devise_installed?
    defined?(Devise)
  end
end

# 3. Add this to config/routes.rb:
# get '/health', to: 'health#check'`,
      setup: [
        'Generate controller: rails generate controller Health',
        'Add the code to app/controllers/health_controller.rb',
        'Add route to config/routes.rb: get \'/health\', to: \'health#check\'',
        'Restart server: rails server',
        'Test: curl http://localhost:3000/health'
      ]
    },
    'ruby-sinatra': {
      name: 'Ruby + Sinatra',
      description: 'Add this to your Sinatra app',
      code: `# Health Check - Ruby/Sinatra
# 1. Install: gem install sinatra
# 2. Add this to your app.rb

require 'sinatra'
require 'json'

# Track when the server started
START_TIME = Time.now

get '/health' do
  content_type :json
  
  # Get memory usage (works on Unix-like systems)
  memory_usage = begin
    \`ps -o rss= -p #{Process.pid}\`.to_i / 1024.0
  rescue
    0
  end
  
  health_data = {
    status: 'UP',
    timestamp: Time.now.iso8601,
    service: 'my-sinatra-app',
    version: '1.0.0',
    uptime: (Time.now - START_TIME).to_i,
    memory: {
      used_mb: memory_usage.round(2)
    },
    environment: ENV['RACK_ENV'] || 'development',
    ruby_version: RUBY_VERSION
  }
  
  health_data.to_json
end

# Run with: ruby app.rb
# Or with rackup: rackup -p 4567`,
      setup: [
        'Install Sinatra: gem install sinatra',
        'Create or open your app.rb',
        'Add the health check route',
        'Run: ruby app.rb',
        'Test: curl http://localhost:4567/health'
      ]
    },
    'php-laravel': {
      name: 'PHP + Laravel',
      description: 'Add this to your Laravel app',
      code: `<?php
// Health Check - Laravel
// 1. Create route in routes/api.php or routes/web.php
// 2. Create controller: php artisan make:controller HealthController

// ----- Add to routes/api.php -----
Route::get('/health', [HealthController::class, 'check']);

// ----- Create app/Http/Controllers/HealthController.php -----
namespace App\\Http\\Controllers;

use Illuminate\\Http\\Request;
use Illuminate\\Support\\Facades\\Cache;

class HealthController extends Controller
{
    /**
     * Health check endpoint for monitoring
     */
    public function check()
    {
        // Store start time in cache if not exists
        if (!Cache::has('app_start_time')) {
            Cache::forever('app_start_time', time());
        }
        
        $startTime = Cache::get('app_start_time');
        $uptime = time() - $startTime;
        
        $healthData = [
            'status' => 'UP',
            'timestamp' => now()->toIso8601String(),
            'service' => config('app.name', 'laravel-app'),
            'version' => '1.0.0',
            'uptime' => $uptime,
            'memory' => [
                'used_mb' => round(memory_get_usage(true) / 1024 / 1024, 2),
                'peak_mb' => round(memory_get_peak_usage(true) / 1024 / 1024, 2)
            ],
            'environment' => config('app.env'),
            'php_version' => PHP_VERSION,
            'laravel_version' => app()->version()
        ];
        
        return response()->json($healthData, 200);
    }
}`,
      setup: [
        'Create controller: php artisan make:controller HealthController',
        'Add the code to app/Http/Controllers/HealthController.php',
        'Add route to routes/api.php or routes/web.php',
        'Test: curl http://localhost:8000/api/health',
        'Or: php artisan serve then visit /api/health'
      ]
    },
    'dotnet': {
      name: '.NET / C#',
      description: 'Add this to your ASP.NET Core app',
      code: `// Health Check - ASP.NET Core / C#
// 1. Add to your Program.cs or Startup.cs

using Microsoft.AspNetCore.Mvc;
using System.Diagnostics;

// ----- For minimal API (Program.cs) -----
var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

var startTime = DateTime.UtcNow;
var process = Process.GetCurrentProcess();

app.MapGet("/health", () =>
{
    var uptime = (DateTime.UtcNow - startTime).TotalSeconds;
    
    var healthData = new
    {
        status = "UP",
        timestamp = DateTime.UtcNow.ToString("o"),
        service = "my-dotnet-service",
        version = "1.0.0",
        uptime = (int)uptime,
        memory = new
        {
            used_mb = process.WorkingSet64 / 1024 / 1024,
            private_mb = process.PrivateMemorySize64 / 1024 / 1024
        },
        environment = Environment.GetEnvironmentVariable("ASPNETCORE_ENVIRONMENT") ?? "Production",
        dotnet_version = Environment.Version.ToString()
    };
    
    return Results.Ok(healthData);
});

app.Run();

// ----- Or using Controller (HealthController.cs) -----
[ApiController]
[Route("[controller]")]
public class HealthController : ControllerBase
{
    private static readonly DateTime StartTime = DateTime.UtcNow;
    
    [HttpGet]
    public IActionResult Check()
    {
        var process = Process.GetCurrentProcess();
        var uptime = (DateTime.UtcNow - StartTime).TotalSeconds;
        
        var healthData = new
        {
            status = "UP",
            timestamp = DateTime.UtcNow.ToString("o"),
            service = "my-dotnet-service",
            version = "1.0.0",
            uptime = (int)uptime,
            memory = new
            {
                used_mb = process.WorkingSet64 / 1024 / 1024
            },
            environment = Environment.GetEnvironmentVariable("ASPNETCORE_ENVIRONMENT")
        };
        
        return Ok(healthData);
    }
}`,
      setup: [
        'Open Program.cs (minimal API) or create HealthController.cs',
        'Add the health check endpoint code',
        'Run: dotnet run',
        'Test: curl https://localhost:5001/health',
        'Deploy and update the URL with your domain'
      ]
    },
    'java-spring': {
      name: 'Java + Spring Boot',
      description: 'Add this to your Spring Boot app',
      code: `// Health Check - Java Spring Boot
// 1. Create a controller class
// 2. Add to src/main/java/com/yourapp/controller/HealthController.java

package com.yourapp.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import java.time.Instant;
import java.util.HashMap;
import java.util.Map;

@RestController
public class HealthController {
    
    private static final Instant START_TIME = Instant.now();
    private static final String VERSION = "1.0.0";
    
    @GetMapping("/health")
    public Map<String, Object> healthCheck() {
        Runtime runtime = Runtime.getRuntime();
        long uptime = Instant.now().getEpochSecond() - START_TIME.getEpochSecond();
        
        Map<String, Object> memory = new HashMap<>();
        memory.put("total_mb", runtime.totalMemory() / 1024 / 1024);
        memory.put("free_mb", runtime.freeMemory() / 1024 / 1024);
        memory.put("used_mb", (runtime.totalMemory() - runtime.freeMemory()) / 1024 / 1024);
        memory.put("max_mb", runtime.maxMemory() / 1024 / 1024);
        
        Map<String, Object> health = new HashMap<>();
        health.put("status", "UP");
        health.put("timestamp", Instant.now().toString());
        health.put("service", "my-spring-service");
        health.put("version", VERSION);
        health.put("uptime", uptime);
        health.put("memory", memory);
        health.put("environment", System.getenv("SPRING_PROFILES_ACTIVE"));
        health.put("java_version", System.getProperty("java.version"));
        
        return health;
    }
}

// Alternative: Use Spring Boot Actuator (recommended)
// 1. Add to pom.xml:
// <dependency>
//     <groupId>org.springframework.boot</groupId>
//     <artifactId>spring-boot-starter-actuator</artifactId>
// </dependency>
//
// 2. Add to application.properties:
// management.endpoints.web.exposure.include=health
// management.endpoint.health.show-details=always`,
      setup: [
        'Create HealthController.java in your controller package',
        'Add the health check method',
        'Run: ./mvnw spring-boot:run (Maven) or ./gradlew bootRun (Gradle)',
        'Test: curl http://localhost:8080/health',
        'OR use Spring Boot Actuator for built-in health checks'
      ]
    }
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
    const selectedExample = codeExamples[selectedLanguage];
    if (selectedExample) {
      navigator.clipboard.writeText(selectedExample.code);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleSubmit = async () => {
    if (selectedTheme && emailAlerts && homepageUrl && appName && slug) {
      setIsSubmitting(true);
      
      // Prepare data for submission
      const onboardingData = {
        name: "User",
        homepage: homepageUrl,
        alerts: emailAlerts === 'yes' ? 'y' : 'n',
        theme: selectedTheme,
        appName: appName,
        slug: slug
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
        const selectedExample = codeExamples[selectedLanguage] || Object.values(codeExamples)[0];
        
        return (
          <div className="space-y-6 animate-fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-cyan-400 mb-2" style={{fontFamily: "'Press Start 2P', monospace"}}>
                STEP 1: HEALTH ENDPOINT
              </h2>
              <p className="text-purple-400 text-sm">Add a health check endpoint to your service</p>
            </div>

            {/* Instructions Banner */}
            <div className="bg-gradient-to-r from-cyan-900/30 to-purple-900/30 rounded-lg p-6 border border-cyan-500/50">
              <h3 className="text-cyan-300 font-bold mb-3 flex items-center gap-2">
                <span className="text-xl">📋</span>
                What you need to do:
              </h3>
              <ol className="text-sm text-gray-300 space-y-2">
                <li className="flex gap-3">
                  <span className="text-cyan-400 font-bold">1.</span>
                  <span>Select your language/framework below</span>
                </li>
                <li className="flex gap-3">
                  <span className="text-cyan-400 font-bold">2.</span>
                  <span>Copy the code template provided</span>
                </li>
                <li className="flex gap-3">
                  <span className="text-cyan-400 font-bold">3.</span>
                  <span>Add it to your application following the setup instructions</span>
                </li>
                <li className="flex gap-3">
                  <span className="text-cyan-400 font-bold">4.</span>
                  <span>Deploy your changes and test the endpoint</span>
                </li>
                <li className="flex gap-3">
                  <span className="text-cyan-400 font-bold">5.</span>
                  <span>Come back here and enter your health check URL in the next step</span>
                </li>
              </ol>
            </div>

            <div className="bg-black/50 rounded-lg p-6 border border-purple-500/50">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                  <Globe className="w-5 h-5 text-cyan-400" />
                  <div>
                    <span className="text-cyan-400 text-sm font-mono block">SELECT YOUR FRAMEWORK</span>
                    <span className="text-purple-400 text-xs">{selectedExample.description}</span>
                  </div>
                </div>
                <button
                  onClick={copyCode}
                  className="flex items-center gap-2 px-4 py-2 bg-purple-600/20 border border-purple-500 rounded text-xs text-purple-300 hover:bg-purple-600/40 transition-all hover:scale-105"
                >
                  {copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
                  {copied ? 'COPIED!' : 'COPY CODE'}
                </button>
              </div>

              {/* Framework Selection */}
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2 mb-6">
                {Object.entries(codeExamples).map(([key, example]) => (
                  <button
                    key={key}
                    onClick={() => setSelectedLanguage(key)}
                    className={`px-3 py-3 rounded font-mono text-xs transition-all text-left ${
                      selectedLanguage === key
                        ? 'bg-cyan-500/30 border-2 border-cyan-400 text-cyan-300 shadow-lg shadow-cyan-400/30'
                        : 'bg-gray-800/50 border border-gray-600 text-gray-400 hover:border-purple-500 hover:text-gray-300'
                    }`}
                  >
                    <div className="font-bold">{example.name}</div>
                  </button>
                ))}
              </div>

              {/* Code Template */}
              <div className="bg-gray-950 rounded-lg p-4 overflow-x-auto border-2 border-gray-800 mb-4">
                <pre className="text-xs font-mono leading-relaxed">
                  <code className="text-green-400">{selectedExample.code}</code>
                </pre>
              </div>

              {/* Setup Instructions */}
              <div className="bg-gradient-to-r from-purple-900/30 to-pink-900/30 rounded-lg p-4 border border-purple-500/50 mb-4">
                <h4 className="text-purple-300 font-bold text-sm mb-3 flex items-center gap-2">
                  <span>⚙️</span> Setup Instructions for {selectedExample.name}:
                </h4>
                <ol className="space-y-2">
                  {selectedExample.setup.map((step, index) => (
                    <li key={index} className="flex gap-3 text-xs text-gray-300">
                      <span className="text-purple-400 font-bold min-w-[20px]">{index + 1}.</span>
                      <span className="font-mono">{step}</span>
                    </li>
                  ))}
                </ol>
              </div>

              {/* Important Tips */}
              <div className="space-y-2">
                <div className="p-3 bg-gradient-to-r from-cyan-900/20 to-blue-900/20 rounded border border-cyan-500/30">
                  <p className="text-xs text-cyan-300">
                    💡 <strong>TIP:</strong> We'll check this endpoint every 30 seconds. Make sure it returns JSON with a "status" field.
                  </p>
                </div>
                <div className="p-3 bg-gradient-to-r from-yellow-900/20 to-orange-900/20 rounded border border-yellow-500/30">
                  <p className="text-xs text-yellow-300">
                    ⚡ <strong>IMPORTANT:</strong> Your health endpoint must be publicly accessible (not behind authentication).
                  </p>
                </div>
                <div className="p-3 bg-gradient-to-r from-green-900/20 to-emerald-900/20 rounded border border-green-500/30">
                  <p className="text-xs text-green-300">
                    ✅ <strong>EXPECTED RESPONSE:</strong> JSON object with at minimum: <code className="bg-black/50 px-2 py-1 rounded">{`{"status": "UP"}`}</code>
                  </p>
                </div>
              </div>
            </div>
          </div>
        );

      case 1:
        const isValidUrl = homepageUrl && (homepageUrl.startsWith('http://') || homepageUrl.startsWith('https://'));
        
        return (
          <div className="space-y-6 animate-fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-cyan-400 mb-2" style={{fontFamily: "'Press Start 2P', monospace"}}>
                STEP 2: ENTER YOUR SERVICE URL
              </h2>
              <p className="text-purple-400 text-sm">Provide the health check endpoint URL to monitor</p>
            </div>

            {/* Instructions */}
            <div className="bg-gradient-to-r from-cyan-900/30 to-purple-900/30 rounded-lg p-5 border border-cyan-500/50">
              <h3 className="text-cyan-300 font-bold mb-2 flex items-center gap-2">
                <span className="text-xl">🔗</span>
                What is your health endpoint URL?
              </h3>
              <p className="text-sm text-gray-300">
                This is the full URL where you deployed your health check endpoint from Step 1. 
                It should be publicly accessible and return JSON with your service status.
              </p>
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
                  <p className="text-xs text-gray-400 mt-2">Must start with http:// or https://</p>
                </label>

                {/* Examples based on selected language */}
                <div className="p-4 bg-gradient-to-r from-cyan-900/20 to-purple-900/20 rounded border border-cyan-500/30">
                  <p className="text-xs text-cyan-300 mb-2">
                    <strong>📋 Common URL patterns for {codeExamples[selectedLanguage]?.name}:</strong>
                  </p>
                  <ul className="text-xs text-gray-400 space-y-1 font-mono">
                    {selectedLanguage.includes('node') && (
                      <>
                        <li>• https://myapp.herokuapp.com/health</li>
                        <li>• https://myapp.vercel.app/api/health</li>
                        <li>• https://api.myapp.com/health</li>
                      </>
                    )}
                    {selectedLanguage.includes('python') && (
                      <>
                        <li>• https://myapp.herokuapp.com/health</li>
                        <li>• https://myapp-xyz.fly.dev/health</li>
                        <li>• https://api.myapp.com/health</li>
                      </>
                    )}
                    {selectedLanguage.includes('go') && (
                      <>
                        <li>• https://myapp-xyz.fly.dev/health</li>
                        <li>• https://myservice.railway.app/health</li>
                        <li>• https://api.myapp.com/health</li>
                      </>
                    )}
                    {selectedLanguage.includes('ruby') && (
                      <>
                        <li>• https://myapp.herokuapp.com/health</li>
                        <li>• https://myapp-xyz.fly.dev/health</li>
                        <li>• https://api.myapp.com/health</li>
                      </>
                    )}
                    {selectedLanguage.includes('php') && (
                      <>
                        <li>• https://myapp.com/api/health</li>
                        <li>• https://api.myapp.com/health</li>
                        <li>• https://myapp.herokuapp.com/api/health</li>
                      </>
                    )}
                    {selectedLanguage.includes('dotnet') && (
                      <>
                        <li>• https://myapp.azurewebsites.net/health</li>
                        <li>• https://api.myapp.com/health</li>
                        <li>• https://myapp-xyz.fly.dev/health</li>
                      </>
                    )}
                    {selectedLanguage.includes('java') && (
                      <>
                        <li>• https://myapp.herokuapp.com/health</li>
                        <li>• https://myapp-xyz.fly.dev/health</li>
                        <li>• https://api.myapp.com/health</li>
                      </>
                    )}
                  </ul>
                </div>

                {/* Validation feedback */}
                {homepageUrl && !isValidUrl && (
                  <div className="p-3 bg-red-900/20 rounded border border-red-500/50">
                    <p className="text-xs text-red-400">
                      ⚠️ URL must start with http:// or https://
                    </p>
                  </div>
                )}

                {isValidUrl && (
                  <div className="p-3 bg-green-900/20 rounded border border-green-500/50">
                    <p className="text-xs text-green-400">
                      ✅ URL format looks good! We'll monitor this endpoint every 30 seconds.
                    </p>
                  </div>
                )}

                {/* Important notes */}
                <div className="p-4 bg-gradient-to-r from-yellow-900/20 to-orange-900/20 rounded border border-yellow-500/30">
                  <p className="text-xs text-yellow-300 mb-2">
                    <strong>⚡ Before proceeding, make sure:</strong>
                  </p>
                  <ul className="text-xs text-gray-300 space-y-1">
                    <li>✓ Your service is deployed and running</li>
                    <li>✓ The health endpoint is publicly accessible</li>
                    <li>✓ It returns JSON (not HTML or plain text)</li>
                    <li>✓ You can access it in your browser or with curl</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        );

      case 2:
        return (
          <div className="space-y-6 animate-fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-cyan-400 mb-2" style={{fontFamily: "'Press Start 2P', monospace"}}>
                STEP 3: NAME YOUR STATUS PAGE
              </h2>
              <p className="text-purple-400 text-sm">Configure your public status page details</p>
            </div>

            <div className="bg-black/50 rounded-lg p-8 border border-purple-500/50">
              <div className="space-y-6">
                <label className="block">
                  <span className="text-cyan-300 text-sm font-mono mb-2 block">APP NAME (Display Name)</span>
                  <input
                    type="text"
                    value={appName}
                    onChange={(e) => setAppName(e.target.value)}
                    placeholder="My Awesome App"
                    className="w-full px-4 py-3 bg-gray-950 border-2 border-purple-500/50 rounded-lg text-cyan-300 font-mono focus:border-cyan-400 focus:outline-none focus:ring-2 focus:ring-cyan-400/20 transition-all"
                  />
                  <p className="text-xs text-gray-400 mt-2">This will be shown on your public status page</p>
                </label>

                <label className="block">
                  <span className="text-cyan-300 text-sm font-mono mb-2 block">STATUS PAGE SLUG (URL)</span>
                  <div className="flex items-center gap-2">
                    <span className="text-gray-400 text-sm">uplytics.com/status/</span>
                    <input
                      type="text"
                      value={slug}
                      onChange={(e) => setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ''))}
                      placeholder="my-app"
                      className="flex-1 px-4 py-3 bg-gray-950 border-2 border-purple-500/50 rounded-lg text-cyan-300 font-mono focus:border-cyan-400 focus:outline-none focus:ring-2 focus:ring-cyan-400/20 transition-all"
                    />
                  </div>
                  <p className="text-xs text-gray-400 mt-2">Only lowercase letters, numbers, and hyphens allowed</p>
                </label>

                {slug && appName && (
                  <div className="p-4 bg-gradient-to-r from-cyan-900/20 to-purple-900/20 rounded border border-cyan-500/30">
                    <p className="text-xs text-cyan-300 mb-2">
                      <strong>Your public status page will be:</strong>
                    </p>
                    <div className="p-3 bg-black/50 rounded font-mono text-sm text-cyan-400">
                      {window.location.origin}/status/{slug}
                    </div>
                    <p className="text-xs text-gray-400 mt-2">
                      Anyone can visit this URL to check your service status
                    </p>
                  </div>
                )}
              </div>
            </div>
          </div>
        );

      case 3:
        return (
          <div className="space-y-6 animate-fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-cyan-400 mb-2" style={{fontFamily: "'Press Start 2P', monospace"}}>
                STEP 4: CHOOSE YOUR THEME
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

      case 4:
        return (
          <div className="space-y-6 animate-fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-cyan-400 mb-2" style={{fontFamily: "'Press Start 2P', monospace"}}>
                STEP 5: EMAIL ALERTS
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
            {[0, 1, 2, 3, 4].map((step) => (
              <div key={step} className="flex items-center flex-1">
                <div className={`w-10 h-10 rounded-full border-2 flex items-center justify-center font-mono transition-all ${
                  currentStep >= step 
                    ? 'bg-cyan-500/30 border-cyan-400 text-cyan-300 shadow-lg shadow-cyan-400/50' 
                    : 'border-gray-600 text-gray-500'
                }`}>
                  {currentStep > step ? '✓' : step + 1}
                </div>
                {step < 4 && (
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

            {currentStep < 4 ? (
              <button
                onClick={() => setCurrentStep(currentStep + 1)}
                disabled={
                  (currentStep === 1 && !homepageUrl) ||
                  (currentStep === 2 && (!appName || !slug)) ||
                  (currentStep === 3 && !selectedTheme)
                }
                className="px-6 py-3 rounded-lg font-mono text-sm bg-gradient-to-r from-cyan-500 to-purple-500 text-white hover:from-cyan-400 hover:to-purple-400 transition-all flex items-center gap-2 shadow-lg hover:shadow-cyan-400/50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                NEXT <ChevronRight className="w-4 h-4" />
              </button>
            ) : (
              <button
                onClick={handleSubmit}
                disabled={!selectedTheme || !emailAlerts || !homepageUrl || !appName || !slug || isSubmitting}
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
            <span className="text-xs font-mono text-cyan-400">STEP {currentStep + 1} OF 5</span>
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