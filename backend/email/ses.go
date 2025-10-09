package email

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SESClient struct {
	client *sesv2.Client
	sender string
}

// NewSESClient creates a new AWS SES client
func NewSESClient() (*SESClient, error) {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sender := os.Getenv("SES_SENDER_EMAIL")

	if region == "" || accessKey == "" || secretKey == "" || sender == "" {
		return nil, fmt.Errorf("missing AWS SES configuration in environment variables")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	return &SESClient{
		client: sesv2.NewFromConfig(cfg),
		sender: sender,
	}, nil
}

// AlertEmail represents the data for an alert email
type AlertEmail struct {
	AppName      string
	HealthURL    string
	StatusCode   int
	Status       string
	ErrorMessage string
	Timestamp    time.Time
	UserEmail    string
	Plan         string
}

// SendDowntimeAlert sends an email alert when a service goes down
func (s *SESClient) SendDowntimeAlert(alert AlertEmail) error {
	subject := fmt.Sprintf("🔴 Alert: %s is Down", alert.AppName)
	
	htmlBody := s.generateHTMLEmail(alert)
	textBody := s.generateTextEmail(alert)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(s.sender),
		Destination: &types.Destination{
			ToAddresses: []string{alert.UserEmail},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(htmlBody),
						Charset: aws.String("UTF-8"),
					},
					Text: &types.Content{
						Data:    aws.String(textBody),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("📧 Email sent successfully to %s (MessageId: %s)", alert.UserEmail, *result.MessageId)
	return nil
}

func (s *SESClient) generateHTMLEmail(alert AlertEmail) string {
	statusColor := "#ef4444" // red
	statusEmoji := "🔴"
	
	if alert.StatusCode >= 500 {
		statusColor = "#dc2626"
		statusEmoji = "🔴"
	} else if alert.StatusCode >= 400 {
		statusColor = "#f59e0b"
		statusEmoji = "🟡"
	}

	errorMsg := alert.ErrorMessage
	if errorMsg == "" {
		errorMsg = "Service returned an error status code"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Service Alert</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse; background-color: #f3f4f6;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; max-width: 100%%; border-collapse: collapse; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%); border-radius: 8px 8px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: bold;">UpLitycs</h1>
                            <p style="margin: 8px 0 0; color: #e0e7ff; font-size: 14px;">Service Monitoring Alert</p>
                        </td>
                    </tr>
                    
                    <!-- Alert Status -->
                    <tr>
                        <td style="padding: 30px 40px; text-align: center; background-color: %s; border-left: 4px solid %s;">
                            <h2 style="margin: 0; color: #ffffff; font-size: 24px; font-weight: bold;">%s Service Down Detected</h2>
                        </td>
                    </tr>
                    
                    <!-- Content -->
                    <tr>
                        <td style="padding: 30px 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.5;">
                                Your monitored service <strong>%s</strong> is experiencing downtime.
                            </p>
                            
                            <div style="background-color: #f9fafb; border-left: 4px solid %s; padding: 20px; margin: 20px 0; border-radius: 4px;">
                                <table role="presentation" style="width: 100%%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px; font-weight: 600;">Service Name:</td>
                                        <td style="padding: 8px 0; color: #111827; font-size: 14px; text-align: right;"><strong>%s</strong></td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px; font-weight: 600;">Health URL:</td>
                                        <td style="padding: 8px 0; color: #111827; font-size: 14px; text-align: right; word-break: break-all;"><code style="background-color: #e5e7eb; padding: 2px 6px; border-radius: 3px; font-size: 12px;">%s</code></td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px; font-weight: 600;">Status Code:</td>
                                        <td style="padding: 8px 0; color: #111827; font-size: 14px; text-align: right;"><strong style="color: %s;">%d (%s)</strong></td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px; font-weight: 600;">Time Detected:</td>
                                        <td style="padding: 8px 0; color: #111827; font-size: 14px; text-align: right;">%s</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px; font-weight: 600;">Error:</td>
                                        <td style="padding: 8px 0; color: #111827; font-size: 14px; text-align: right;">%s</td>
                                    </tr>
                                </table>
                            </div>
                            
                            <p style="margin: 20px 0; color: #374151; font-size: 14px; line-height: 1.5;">
                                <strong>What to do next:</strong>
                            </p>
                            <ul style="color: #374151; font-size: 14px; line-height: 1.8; margin: 10px 0; padding-left: 20px;">
                                <li>Check your service logs for errors</li>
                                <li>Verify your server is running and accessible</li>
                                <li>Check your network and firewall settings</li>
                                <li>Review recent deployments or configuration changes</li>
                            </ul>
                        </td>
                    </tr>
                    
                    <!-- CTA Button -->
                    <tr>
                        <td style="padding: 0 40px 30px; text-align: center;">
                            <a href="https://uplitycs.com/dashboard" style="display: inline-block; padding: 14px 32px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: 600; font-size: 16px;">View Dashboard</a>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 20px 40px; background-color: #f9fafb; border-top: 1px solid #e5e7eb; border-radius: 0 0 8px 8px;">
                            <p style="margin: 0; color: #6b7280; font-size: 12px; text-align: center; line-height: 1.5;">
                                You're receiving this email because you have alerts enabled for your <strong>%s plan</strong>.<br>
                                This is an automated alert from UpLitycs. Please do not reply to this email.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, statusColor, statusColor, statusEmoji, alert.AppName, statusColor, alert.AppName, alert.HealthURL, statusColor, alert.StatusCode, alert.Status, alert.Timestamp.Format("2006-01-02 15:04:05 MST"), errorMsg, alert.Plan)
}

func (s *SESClient) generateTextEmail(alert AlertEmail) string {
	errorMsg := alert.ErrorMessage
	if errorMsg == "" {
		errorMsg = "Service returned an error status code"
	}

	return fmt.Sprintf(`
UpLitycs - Service Monitoring Alert
====================================

🔴 SERVICE DOWN DETECTED

Your monitored service "%s" is experiencing downtime.

Service Details:
- Service Name: %s
- Health URL: %s
- Status Code: %d (%s)
- Time Detected: %s
- Error: %s

What to do next:
- Check your service logs for errors
- Verify your server is running and accessible
- Check your network and firewall settings
- Review recent deployments or configuration changes

View your dashboard: https://uplitycs.com/dashboard

---
You're receiving this email because you have alerts enabled for your %s plan.
This is an automated alert from UpLitycs. Please do not reply to this email.
`, alert.AppName, alert.AppName, alert.HealthURL, alert.StatusCode, alert.Status, alert.Timestamp.Format("2006-01-02 15:04:05 MST"), errorMsg, alert.Plan)
}
