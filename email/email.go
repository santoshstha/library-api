package email

import (
    "context"
    "fmt"
    "net/smtp"
    "sync"
    "time"
    "library-api/logger"
)

type Email struct {
    To      string
    Subject string
    Body    string
    ID      string
}

type EmailStatus struct {
    ID     string
    Status string
    Time   time.Time
}

type BulkTask struct {
    ID         string
    Emails     map[string]string
    Total      int
    Completed  int
    Statuses   sync.Map
    InProgress bool
    Updates    chan EmailStatus // Add channel for real-time updates
}

type EmailService struct {
    queue    chan Email
    wg       sync.WaitGroup
    smtpHost string
    smtpPort string
    smtpUser string
    smtpPass string
    logger   *logger.AsyncLogger
    tasks    sync.Map
}

func NewEmailService(smtpHost, smtpPort, smtpUser, smtpPass string, workers int, logger *logger.AsyncLogger) *EmailService {
    service := &EmailService{
        queue:    make(chan Email, 1000),
        smtpHost: smtpHost,
        smtpPort: smtpPort,
        smtpUser: smtpUser,
        smtpPass: smtpPass,
        logger:   logger,
    }
    for i := 0; i < workers; i++ {
        service.wg.Add(1)
        go service.worker()
    }
    return service
}

func (s *EmailService) worker() {
    defer s.wg.Done()
    auth := smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)
    addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

    for email := range s.queue {
        msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", email.To, email.Subject, email.Body))
        err := smtp.SendMail(addr, auth, s.smtpUser, []string{email.To}, msg)
        status := EmailStatus{ID: email.ID, Time: time.Now()}
        if err != nil {
            s.logger.Log(fmt.Sprintf("Failed to send email to %s: %v", email.To, err))
            status.Status = "failed"
        } else {
            s.logger.Log(fmt.Sprintf("Email sent to %s", email.To))
            status.Status = "completed"
        }
        s.logger.Log(fmt.Sprintf("Updating status for %s to %s", email.ID, status.Status))
        s.updateTaskStatus(email.ID, status)
        time.Sleep(100 * time.Millisecond)
    }
}

func (s *EmailService) updateTaskStatus(emailID string, status EmailStatus) {
    s.tasks.Range(func(key, value interface{}) bool {
        task := value.(*BulkTask)
        if _, exists := task.Emails[status.ID]; exists {
            task.Statuses.Store(emailID, status)
            task.Updates <- status // Send real-time update
            if status.Status == "completed" || status.Status == "failed" {
                task.Completed++
                s.logger.Log(fmt.Sprintf("Task %s: Completed %d/%d", task.ID, task.Completed, task.Total))
                if task.Completed == task.Total {
                    task.InProgress = false
                    s.logger.Log(fmt.Sprintf("Task %s fully completed", task.ID))
                    close(task.Updates) // Close channel when done
                }
            }
        }
        return true
    })
}

func (s *EmailService) Send(to, subject, body, id string) {
    email := Email{To: to, Subject: subject, Body: body, ID: id}
    select {
    case s.queue <- email:
    default:
        s.logger.Log(fmt.Sprintf("Email queue full, dropping email to %s", to))
        s.updateTaskStatus(id, EmailStatus{ID: id, Status: "failed", Time: time.Now()})
    }
}

func (s *EmailService) SendBulk(recipients []string, subject, body string) string {
    taskID := fmt.Sprintf("task_%d", time.Now().UnixNano())
    ids := make(map[string]string)
    for i, to := range recipients {
        id := fmt.Sprintf("email_%d_%s", i, to)
        ids[to] = id
    }

    task := &BulkTask{
        ID:         taskID,
        Emails:     ids,
        Total:      len(recipients),
        Completed:  0,
        InProgress: true,
        Updates:    make(chan EmailStatus, len(recipients)), // Buffer for updates
    }
    s.tasks.Store(taskID, task)
    s.logger.Log(fmt.Sprintf("Created task %s with %d emails", taskID, task.Total))

    go func() {
        for to, id := range ids {
            s.Send(to, subject, body, id)
            task.Statuses.Store(id, EmailStatus{ID: id, Status: "in_progress", Time: time.Now()})
            task.Updates <- EmailStatus{ID: id, Status: "in_progress", Time: time.Now()}
        }
    }()

    return taskID
}

func (s *EmailService) GetTaskStatus(taskID string) (*BulkTask, bool) {
    task, ok := s.tasks.Load(taskID)
    if !ok {
        return nil, false
    }
    return task.(*BulkTask), true
}

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