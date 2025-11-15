package models

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Report struct {
    ID         int
    ReporterIP string
    ReportedIP string
    Reason     string
    Timestamp  string
}

type BannedIP struct {
    IP       string
    Reason   string
    BannedAt string
}

func InitDB() {
    os.MkdirAll("./data", 0755)

    var err error
    DB, err = sql.Open("sqlite3", "./data/app.db")
    if err != nil {
        log.Fatal(err)
    }

    _, err = DB.Exec(`
        PRAGMA foreign_keys = ON;
        PRAGMA journal_mode = WAL;
        PRAGMA synchronous = NORMAL;

        CREATE TABLE IF NOT EXISTS reports (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            reporter_ip TEXT NOT NULL,
            reported_ip TEXT NOT NULL,
            reason TEXT NOT NULL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS banned_ips (
            ip TEXT PRIMARY KEY,
            reason TEXT,
            banned_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS admins (
            id INTEGER PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL
        );
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Insert default admin
    hash, _ := HashPassword(os.Getenv("ADMIN_PASSWORD"))
    DB.Exec(`INSERT OR IGNORE INTO admins (username, password_hash) VALUES (?, ?)`,
        os.Getenv("ADMIN_USERNAME"), hash)
}

func IsBanned(ip string) bool {
    var count int
    DB.QueryRow("SELECT COUNT(*) FROM banned_ips WHERE ip = ?", ip).Scan(&count)
    return count > 0
}
