package email

import (
	"context"
	"fmt"
	"net/smtp"
	"sync"
	"time"
	"library-api/logger" // Correct import
)

// Email represents an email job
type Email struct {
	To      string
	Subject string
	Body    string
}

// EmailService manages email sending with a queue and workers
type EmailService struct {
	queue    chan Email
	wg       sync.WaitGroup
	smtpHost string
	smtpPort string
	smtpUser string
	smtpPass string
	logger   *logger.AsyncLogger // Instance field
}

// NewEmailService initializes the service with a worker pool and logger
func NewEmailService(smtpHost, smtpPort, smtpUser, smtpPass string, workers int, logger *logger.AsyncLogger) *EmailService {
	service := &EmailService{
		queue:    make(chan Email, 1000),
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		smtpUser: smtpUser,
		smtpPass: smtpPass,
		logger:   logger, // Store the passed logger
	}
	for i := 0; i < workers; i++ {
		service.wg.Add(1)
		go service.worker()
	}
	return service
}

// worker processes emails from the queue
func (s *EmailService) worker() {
	defer s.wg.Done()
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	for email := range s.queue {
		msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", email.To, email.Subject, email.Body))
		err := smtp.SendMail(addr, auth, s.smtpUser, []string{email.To}, msg)
		if err != nil {
			s.logger.Log(fmt.Sprintf("Failed to send email to %s: %v", email.To, err)) // Line 57
			continue
		}
		s.logger.Log(fmt.Sprintf("Email sent to %s", email.To)) // Line 60
		time.Sleep(100 * time.Millisecond)
	}
}

// Send queues an email to be sent
func (s *EmailService) Send(to, subject, body string) {
	select {
	case s.queue <- Email{To: to, Subject: subject, Body: body}:
	default:
		s.logger.Log(fmt.Sprintf("Email queue full, dropping email to %s", to)) // Line 72
	}
}

// SendBulk sends emails to multiple recipients
func (s *EmailService) SendBulk(recipients []string, subject, body string) {
	for _, to := range recipients {
		s.Send(to, subject, body)
	}
}

// Shutdown closes the queue and waits for workers to finish
func (s *EmailService) Shutdown(ctx context.Context) {
	close(s.queue)
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		s.logger.Log("Email service shutdown timed out")
	}
}