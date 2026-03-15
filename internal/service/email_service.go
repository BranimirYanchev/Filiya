package service

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

var emailConfig *EmailConfig

// InitEmailConfig initializes email configuration from environment variables
func InitEmailConfig() {
	emailConfig = &EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     587, // Default to 587, override with SMTP_PORT if needed
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("SMTP_FROM_EMAIL"),
		FromName:     os.Getenv("SMTP_FROM_NAME"),
	}

	if port := os.Getenv("SMTP_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &emailConfig.SMTPPort)
	}
}

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string) error {
	if emailConfig == nil {
		return errors.New("email config not initialized")
	}

	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", emailConfig.FromName, emailConfig.FromEmail))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := mail.NewDialer(emailConfig.SMTPHost, emailConfig.SMTPPort, emailConfig.SMTPUser, emailConfig.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Error("Failed to send email: ", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info("Email sent successfully to: ", to)
	return nil
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(to, resetToken, resetURL string) error {
	subject := "Password Reset Request - Filia"

	// You can customize this URL based on your frontend
	if resetURL == "" {
		resetURL = os.Getenv("FRONTEND_URL") + "/reset-password?token=" + resetToken
	} else {
		resetURL = resetURL + "?token=" + resetToken
	}

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Password Reset Request</h2>
			<p>You have requested to reset your password for your Filia account.</p>
			<p>Click the link below to reset your password:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>Or copy this link: %s</p>
			<p>This link will expire in 1 hour.</p>
			<p>If you did not request this, please ignore this email.</p>
		</body>
		</html>
	`, resetURL, resetURL)

	return SendEmail(to, subject, body)
}

// SendFriendRequestEmail sends a friend request notification email
func SendFriendRequestEmail(to, senderName, acceptURL string) error {
	subject := fmt.Sprintf("New Friend Request from %s - Filia", senderName)

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>New Friend Request</h2>
			<p><strong>%s</strong> sent you a friend request on Filia.</p>
			<p><a href="%s">View Friend Request</a></p>
		</body>
		</html>
	`, senderName, acceptURL)

	return SendEmail(to, subject, body)
}
