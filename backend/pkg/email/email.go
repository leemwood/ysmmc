package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/ysmmc/backend/internal/config"
)

type EmailService struct {
	host     string
	port     int
	user     string
	password string
	from     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		host:     config.AppConfig.SMTPHost,
		port:     config.AppConfig.SMTPPort,
		user:     config.AppConfig.SMTPUser,
		password: config.AppConfig.SMTPPassword,
		from:     config.AppConfig.SMTPFrom,
	}
}

type ResetPasswordData struct {
	ResetLink string
}

type WelcomeData struct {
	Username   string
	VerifyLink string
}

type ModelReviewData struct {
	Username  string
	ModelTitle string
	Status    string
	Reason    string
	ModelLink string
}

func (s *EmailService) Send(to, subject, body string) error {
	if s.host == "" || s.user == "" {
		return fmt.Errorf("SMTP not configured")
	}

	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	msg := fmt.Sprintf("From: %s\r\n", s.from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	msg += body

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: s.host}
		if err = client.StartTLS(config); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(s.user); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return client.Quit()
}

func (s *EmailService) renderTemplate(name string, data interface{}) (string, error) {
	templatePath := filepath.Join("templates", "emails", name+".html")
	
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (s *EmailService) SendResetPassword(to, resetLink string) error {
	data := ResetPasswordData{ResetLink: resetLink}
	body, err := s.renderTemplate("reset_password", data)
	if err != nil {
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>重置密码</h2>
				<p>请点击以下链接重置您的密码：</p>
				<p><a href="%s">%s</a></p>
				<p>此链接将在1小时后失效。</p>
			</body>
			</html>
		`, resetLink, resetLink)
	}

	return s.Send(to, "重置密码 - YSM模型站", body)
}

func (s *EmailService) SendWelcome(to, username, verifyLink string) error {
	data := WelcomeData{
		Username:   username,
		VerifyLink: verifyLink,
	}
	body, err := s.renderTemplate("welcome", data)
	if err != nil {
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>欢迎注册 YSM模型站</h2>
				<p>您好，%s！</p>
				<p>请点击以下链接验证您的邮箱：</p>
				<p><a href="%s">%s</a></p>
			</body>
			</html>
		`, username, verifyLink, verifyLink)
	}

	return s.Send(to, "欢迎注册 - YSM模型站", body)
}

func (s *EmailService) SendModelReview(to, username, modelTitle, status, reason, modelLink string) error {
	data := ModelReviewData{
		Username:   username,
		ModelTitle: modelTitle,
		Status:     status,
		Reason:     reason,
		ModelLink:  modelLink,
	}
	body, err := s.renderTemplate("model_review", data)
	if err != nil {
		statusText := "已通过"
		if status == "rejected" {
			statusText = "未通过"
		}
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>模型审核通知</h2>
				<p>您好，%s！</p>
				<p>您的模型「%s」审核%s。</p>
				<p><a href="%s">查看详情</a></p>
			</body>
			</html>
		`, username, modelTitle, statusText, modelLink)
	}

	subject := "模型审核通过 - YSM模型站"
	if status == "rejected" {
		subject = "模型审核未通过 - YSM模型站"
	}

	return s.Send(to, subject, body)
}

func (s *EmailService) IsConfigured() bool {
	return s.host != "" && s.user != "" && s.password != ""
}
