package handlers

import (
    "encoding/json"
    "errors"
    "io"
    "net/http"

    "github.com/srikesh2k4/ovitv/backend/models"
    "github.com/srikesh2k4/ovitv/backend/utils"
)


type ReportRequest struct {
    ReportedIP string `json:"reported_ip"`
    Reason     string `json:"reason"`
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Ensure application/json
    if r.Header.Get("Content-Type") != "application/json" {
        http.Error(w, "Invalid content type", http.StatusBadRequest)
        return
    }

    // Prevent huge body payload attacks
    r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB limit

    var data ReportRequest
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := json.Unmarshal(body, &data); err != nil {
        http.Error(w, "Malformed JSON", http.StatusBadRequest)
        return
    }

    // Validate input
    if err := validateReport(data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    reporter := utils.GetIP(r)
    reason := utils.Sanitize(data.Reason)
    reported := utils.Sanitize(data.ReportedIP)

    _, err = models.DB.Exec(
        "INSERT INTO reports (reporter_ip, reported_ip, reason) VALUES (?, ?, ?)",
        reporter, reported, reason,
    )

    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Return JSON success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
    })
}

func validateReport(data ReportRequest) error {
    if data.ReportedIP == "" {
        return errors.New("reported_ip is required")
    }
    if data.Reason == "" {
        return errors.New("reason is required")
    }
    if len(data.Reason) > 500 {
        return errors.New("reason too long")
    }
    return nil
}
