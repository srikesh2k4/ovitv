package handlers

import (
    "html/template"
    "net/http"
    "ome-tv-pro/backend/middleware"
    "ome-tv-pro/backend/models"
    "ome-tv-pro/backend/utils"
)


const adminSession = "admin-session"

// ---------------------- LOGIN PAGE ----------------------

func AdminLoginPage(tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        // If already logged in -> redirect
        sess, _ := middleware.Store.Get(r, adminSession)
        if auth, ok := sess.Values["auth"].(bool); ok && auth {
            http.Redirect(w, r, "/admin/", http.StatusSeeOther)
            return
        }

        // Generate CSRF token
        csrf, _ := middleware.SetCSRF(w, r)

        err := tmpl.ExecuteTemplate(w, "login.html", map[string]string{
            "csrf": csrf,
        })

        if err != nil {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    }
}

// ---------------------- LOGIN POST ----------------------

func AdminLogin(w http.ResponseWriter, r *http.Request) {

    username := utils.Sanitize(r.FormValue("username"))
    password := r.FormValue("password")

    var hash string
    err := models.DB.QueryRow(
        "SELECT password_hash FROM admins WHERE username = ?",
        username,
    ).Scan(&hash)

    if err != nil || !models.CheckPassword(password, hash) {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    session, _ := middleware.Store.Get(r, adminSession)
    session.Values["auth"] = true
    session.Save(r, w)

    http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}

// ---------------------- DASHBOARD ----------------------

func AdminDashboard(tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        // Fetch reports
        rows, err := models.DB.Query(`
            SELECT id, reporter_ip, reported_ip, reason, timestamp 
            FROM reports 
            ORDER BY timestamp DESC
        `)
        if err != nil {
            http.Error(w, "DB error", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var reports []models.Report
        for rows.Next() {
            var rep models.Report
            if err := rows.Scan(&rep.ID, &rep.ReporterIP, &rep.ReportedIP, &rep.Reason, &rep.Timestamp); err == nil {
                reports = append(reports, rep)
            }
        }

        // Fetch bans
        rows2, err := models.DB.Query("SELECT ip, reason, banned_at FROM banned_ips")
        if err != nil {
            http.Error(w, "DB error", http.StatusInternalServerError)
            return
        }
        defer rows2.Close()

        var banned []models.BannedIP
        for rows2.Next() {
            var b models.BannedIP
            if err := rows2.Scan(&b.IP, &b.Reason, &b.BannedAt); err == nil {
                banned = append(banned, b)
            }
        }

        err = tmpl.ExecuteTemplate(w, "dashboard.html", struct {
            Reports []models.Report
            Banned  []models.BannedIP
        }{
            Reports: reports,
            Banned:  banned,
        })

        if err != nil {
            http.Error(w, "Template error", http.StatusInternalServerError)
        }
    }
}

// ---------------------- BAN / UNBAN ----------------------

func BanIP(w http.ResponseWriter, r *http.Request) {
    ip := utils.Sanitize(r.FormValue("ip"))
    reason := utils.Sanitize(r.FormValue("reason"))

    models.DB.Exec("INSERT OR IGNORE INTO banned_ips (ip, reason) VALUES (?, ?)", ip, reason)
    http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}

func UnbanIP(w http.ResponseWriter, r *http.Request) {
    ip := utils.Sanitize(r.URL.Query().Get("ip"))
    models.DB.Exec("DELETE FROM banned_ips WHERE ip = ?", ip)
    http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}

// ---------------------- LOGOUT ----------------------

func AdminLogout(w http.ResponseWriter, r *http.Request) {
    sess, _ := middleware.Store.Get(r, adminSession)
    delete(sess.Values, "auth")
    sess.Save(r, w)
    http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
