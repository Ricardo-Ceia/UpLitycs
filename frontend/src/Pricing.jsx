import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Check, Zap, Rocket, Crown } from 'lucide-react';
import './Pricing.css';

const Pricing = () => {
  const [billingCycle, setBillingCycle] = useState('monthly'); // 'monthly' or 'yearly'
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  // Check if user is authenticated on mount
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await fetch('/auth/check-session', {
          credentials: 'include'
        });
        setIsAuthenticated(response.ok);
      } catch (err) {
        console.error('Error checking auth:', err);
        setIsAuthenticated(false);
      } finally {
        setIsLoading(false);
      }
    };
    checkAuth();
  }, []);

  const plans = [
    {
      name: 'Free',
      icon: '🆓',
      price: { monthly: 0, yearly: 0 },
      description: 'Perfect for trying out StatusFrame',
      features: [
        { text: '1 monitor', included: true },
        { text: '5-minute checks', included: true },
        { text: '7 days data retention', included: true },
        { text: 'Public status page', included: true },
        { text: 'Community support', included: true },
      ],
      cta: 'Get Started Free',
      highlight: false,
      tier: 'free'
    },
    {
      name: 'Pro',
      icon: '⚡',
      price: { monthly: 19, yearly: 190 },
      description: 'For serious projects',
      features: [
        { text: '25 monitors', included: true },
        { text: '1-minute checks', included: true },
        { text: '30 days data retention', included: true },
        { text: 'Public status pages', included: true },
        { text: 'Custom branding (logo)', included: true },
        { text: 'Response time graphs', included: true },
        { text: 'Priority support', included: true },
      ],
      cta: 'Start Pro',
      highlight: true,
      popular: true,
      tier: 'pro'
    },
    {
      name: 'Business',
      icon: '🚀',
      price: { monthly: 29, yearly: 290 },
      description: 'For teams and agencies',
      features: [
        { text: '100 monitors', included: true },
        { text: '30-second checks', included: true },
        { text: '90 days data retention', included: true },
        { text: 'Public status pages', included: true },
        { text: 'Custom branding (logo)', included: true },
        { text: 'Advanced analytics', included: true },
        { text: '24/7 priority support', included: true },
      ],
      cta: 'Start Business',
      highlight: false,
      tier: 'business'
    }
  ];

  const handleSelectPlan = async (tier) => {
    // If not authenticated, redirect to auth first
    if (!isAuthenticated) {
      navigate('/auth');
      return;
    }

    // User is authenticated
    if (tier === 'free') {
      // Free plan: go directly to onboarding
      navigate('/onboarding');
    } else {
      // Paid plan: create Stripe checkout session
      try {
        const response = await fetch('/api/create-checkout-session', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({
            plan: tier,
            billing_period: billingCycle
          })
        });

        if (response.ok) {
          const data = await response.json();
          // Redirect to Stripe checkout
          window.location.href = data.url;
        } else {
          console.error('Failed to create checkout session');
          alert('Failed to create checkout session. Please try again.');
        }
      } catch (err) {
        console.error('Error creating checkout session:', err);
        alert('An error occurred. Please try again.');
      }
    }
  };

  const calculateSavings = (plan) => {
    if (plan.price.monthly === 0) return null;
    const monthlyTotal = plan.price.monthly * 12;
    const yearlySavings = monthlyTotal - plan.price.yearly;
    return yearlySavings;
  };

  return (
    <div className="pricing-container">
      <div className="pricing-hero">
        <div className="glitch-text-container">
          <h1 className="pricing-title">
            <span className="glitch" data-text="Choose Your Plan">Choose Your Plan</span>
          </h1>
        </div>
        <p className="pricing-subtitle">
          Free plan forever. No credit card required to start.
        </p>

        {/* Billing Cycle Toggle */}
        <div className="billing-toggle">
          <button
            className={`toggle-btn ${billingCycle === 'monthly' ? 'active' : ''}`}
            onClick={() => setBillingCycle('monthly')}
          >
            Monthly
          </button>
          <button
            className={`toggle-btn ${billingCycle === 'yearly' ? 'active' : ''}`}
            onClick={() => setBillingCycle('yearly')}
          >
            Yearly
            <span className="save-badge">Save 17%</span>
          </button>
        </div>
      </div>

      {/* Pricing Cards */}
      <div className="pricing-grid">
        {plans.map((plan, index) => {
          const savings = calculateSavings(plan);
          const price = billingCycle === 'yearly' ? plan.price.yearly : plan.price.monthly;
          const priceDisplay = billingCycle === 'yearly' ? (price / 12).toFixed(0) : price;

          return (
            <div
              key={index}
              className={`pricing-card ${plan.highlight ? 'highlighted' : ''} ${plan.popular ? 'popular' : ''}`}
            >
              {plan.popular && (
                <div className="popular-badge">
                  <Zap size={14} />
                  Most Popular
                </div>
              )}

              <div className="plan-icon">{plan.icon}</div>
              <h3 className="plan-name">{plan.name}</h3>
              <p className="plan-description">{plan.description}</p>

              <div className="plan-price">
                <span className="currency">$</span>
                <span className="amount">{priceDisplay}</span>
                <span className="period">/month</span>
              </div>

              {billingCycle === 'yearly' && savings && (
                <div className="savings-text">
                  Save ${savings}/year
                </div>
              )}

              {billingCycle === 'yearly' && plan.price.yearly > 0 && (
                <div className="billed-yearly">
                  Billed ${plan.price.yearly}/year
                </div>
              )}

              <button
                className="plan-cta"
                onClick={() => handleSelectPlan(plan.tier)}
              >
                {plan.cta}
              </button>

              <div className="features-list">
                {plan.features.map((feature, idx) => (
                  <div key={idx} className={`feature-item ${feature.included ? 'included' : 'not-included'}`}>
                    <Check className="check-icon" size={16} />
                    <span>{feature.text}</span>
                  </div>
                ))}
              </div>
            </div>
          );
        })}
      </div>

      {/* FAQ Section */}
      <div className="pricing-faq">
        <h2 className="faq-title">Frequently Asked Questions</h2>
        <div className="faq-grid">
          <div className="faq-item">
            <h3>What's a monitor?</h3>
            <p>A monitor is a single endpoint you want to track (e.g., your API or website health check).</p>
          </div>
          <div className="faq-item">
            <h3>Can I upgrade later?</h3>
            <p>Yes! Upgrade or downgrade anytime. Changes take effect immediately with prorated billing.</p>
          </div>
        </div>
      </div>

      {/* CRT Scanlines Effect */}
      <div className="crt-scanlines"></div>
    </div>
  );
};

export default Pricing;
