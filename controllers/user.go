package controllers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
    "library-api/models"
    "library-api/services"
    "github.com/dgrijalva/jwt-go"
)

type UserController struct {
    Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
    return &UserController{Service: service}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Email    string `json:"email"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    user, err := c.Service.CreateUser(req.Username, req.Password, req.Email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    user, err := c.Service.Login(req.Username, req.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (c *UserController) SendBulkEmails(w http.ResponseWriter, r *http.Request) {
    users, err := c.Service.Repo.DB.Find(&[]models.User{}).Rows()
    if err != nil {
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }
    defer users.Close()

    var emails []string
    for users.Next() {
        var user models.User
        c.Service.Repo.DB.ScanRows(users, &user)
        if user.Email != "" {
            emails = append(emails, user.Email)
        }
    }

    taskID := c.Service.EmailService.SendBulk(emails, "Library Update", "New books added!")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "task_id": taskID,
        "message": "Bulk emails queued",
    })
}

func (c *UserController) StreamBulkEmailStatus(w http.ResponseWriter, r *http.Request) {
    taskID := r.URL.Query().Get("task_id")
    if taskID == "" {
        http.Error(w, "Missing task_id", http.StatusBadRequest)
        return
    }

    task, ok := c.Service.EmailService.GetTaskStatus(taskID)
    if !ok {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    // Set SSE headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming not supported", http.StatusInternalServerError)
        return
    }

    // Stream updates
    type EmailStatusResponse struct {
        Email  string `json:"email"`
        ID     string `json:"id"`
        Status string `json:"status"`
        Time   string `json:"time"`
    }

    for status := range task.Updates {
        var email string
        for e, id := range task.Emails {
            if id == status.ID {
                email = e
                break
            }
        }
        resp := EmailStatusResponse{
            Email:  email,
            ID:     status.ID,
            Status: status.Status,
            Time:   status.Time.Format(time.RFC3339),
        }
        data, _ := json.Marshal(resp)
        fmt.Fprintf(w, "data: %s\n\n", data)
        flusher.Flush()
    }

    // Send final task completion event
    final := struct {
        TaskID     string `json:"task_id"`
        Total      int    `json:"total"`
        Completed  int    `json:"completed"`
        InProgress bool   `json:"in_progress"`
    }{
        TaskID:     task.ID,
        Total:      task.Total,
        Completed:  task.Completed,
        InProgress: task.InProgress,
    }
    data, _ := json.Marshal(final)
    fmt.Fprintf(w, "data: %s\n\n", data)
    flusher.Flush()
}

// Keep GetBulkEmailStatus for polling if needed
func (c *UserController) GetBulkEmailStatus(w http.ResponseWriter, r *http.Request) {
    taskID := r.URL.Query().Get("task_id")
    if taskID == "" {
        http.Error(w, "Missing task_id", http.StatusBadRequest)
        return
    }

    task, ok := c.Service.EmailService.GetTaskStatus(taskID)
    if !ok {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
	
	type EmailStatus struct {
        ID     string
        Status string
        Time   time.Time
    }

    type EmailStatusResponse struct {
        Email  string `json:"email"`
        ID     string `json:"id"`
        Status string `json:"status"`
        Time   string `json:"time"`
    }
    type TaskResponse struct {
        TaskID     string              `json:"task_id"`
        Total      int                 `json:"total"`
        Completed  int                 `json:"completed"`
        InProgress bool                `json:"in_progress"`
        Emails     []EmailStatusResponse `json:"emails"`
    }

    response := TaskResponse{
        TaskID:     task.ID,
        Total:      task.Total,
        Completed:  task.Completed,
        InProgress: task.InProgress,
        Emails:     []EmailStatusResponse{},
    }

    task.Statuses.Range(func(key, value interface{}) bool {
        status := value.(EmailStatus)
        var email string
        for e, id := range task.Emails {
            if id == status.ID {
                email = e
                break
            }
        }
        response.Emails = append(response.Emails, EmailStatusResponse{
            Email:  email,
            ID:     status.ID,
            Status: status.Status,
            Time:   status.Time.Format(time.RFC3339),
        })
        return true
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}