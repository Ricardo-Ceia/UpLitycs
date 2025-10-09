import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Check, Zap, Rocket, Crown } from 'lucide-react';
import './Pricing.css';

const Pricing = () => {
  const [billingCycle, setBillingCycle] = useState('monthly'); // 'monthly' or 'yearly'
  const navigate = useNavigate();

  const plans = [
    {
      name: 'Free',
      icon: '🆓',
      price: { monthly: 0, yearly: 0 },
      description: 'Perfect for trying out UpLitycs',
      features: [
        { text: '1 monitor', included: true },
        { text: '5-minute checks', included: true },
        { text: 'Public status page', included: true },
        { text: 'Email alerts', included: true },
      ],
      cta: 'Get Started Free',
      highlight: false,
      tier: 'free'
    },
    {
      name: 'Pro',
      icon: '⚡',
      price: { monthly: 12, yearly: 120 },
      description: 'For serious projects',
      features: [
        { text: '10 monitors', included: true },
        { text: '1-minute checks', included: true },
        { text: 'Custom slugs', included: true },
        { text: 'Unlimited alerts', included: true },
      ],
      cta: 'Pick Pro',
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
        { text: '50 monitors', included: true },
        { text: '30-second checks', included: true },
        { text: 'Webhook alerts', included: true },
        { text: 'Priority support', included: true },
      ],
      cta: 'Pick Business',
      highlight: false,
      tier: 'business'
    }
  ];

  const handleSelectPlan = (tier) => {
    if (tier === 'free') {
      navigate('/auth');
    } else {
      // Store selected plan in localStorage for after auth
      localStorage.setItem('selectedPlan', JSON.stringify({ tier, billingCycle }));
      navigate('/auth');
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
          Start with 7 days free. No credit card required.
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
