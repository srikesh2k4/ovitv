package main

import (
    "embed"
    "html/template"
    "log"
    "net/http"

    "github.com/srikesh2k4/ovitv/backend/handlers"
    "github.com/srikesh2k4/ovitv/backend/middleware"
    "github.com/srikesh2k4/ovitv/backend/models"
    "github.com/srikesh2k4/ovitv/backend/signaling"

    "github.com/joho/godotenv"
)

//go:embed templates/*
var adminTemplates embed.FS

func main() {

    // Load environment variables
    godotenv.Load()

    // Init SQLite DB
    models.InitDB()

    // Start WebSocket hub
    hub := signaling.NewHub()
    go hub.Run()

    // Load admin templates
    tmpl := template.Must(template.ParseFS(adminTemplates, "templates/*.html"))

    // Build middleware chain
    secure := middleware.Chain(
        middleware.RateLimit,
        middleware.CORS,
        middleware.SecurityHeaders,
    )

    // ------------------ API ROUTES ------------------

    http.Handle("/api/report", secure(http.HandlerFunc(handlers.ReportHandler)))

    http.Handle("/ws", secure(signaling.ServeWs(hub)))

    // ------------------ ADMIN ROUTES ------------------

    http.Handle("/admin/login", handlers.AdminLoginPage(tmpl))
    http.Handle("/admin/auth", http.HandlerFunc(handlers.AdminLogin))

    http.Handle("/admin/", middleware.AdminAuth(handlers.AdminDashboard(tmpl)))
    http.Handle("/admin/ban", middleware.AdminAuth(http.HandlerFunc(handlers.BanIP)))
    http.Handle("/admin/unban", middleware.AdminAuth(http.HandlerFunc(handlers.UnbanIP)))
    http.Handle("/admin/logout", http.HandlerFunc(handlers.AdminLogout))

    // Start server
    log.Println("Backend running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
