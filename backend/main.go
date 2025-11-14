package main

import (
    "embed"
    "html/template"
    "log"
    "net/http"

    "ome-tv-pro/backend/handlers"
    "ome-tv-pro/backend/middleware"
    "ome-tv-pro/backend/models"
    "ome-tv-pro/backend/signaling"

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

    // Report endpoint
    http.Handle("/api/report", secure(http.HandlerFunc(handlers.ReportHandler)))

    // WebSocket endpoint
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
