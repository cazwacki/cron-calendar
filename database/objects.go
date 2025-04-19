package database

import "time"

type User struct {
	ID               string     `json:"username" binding:"required" gorm:"default:null"`
	Password         string     `json:"password" binding:"required" gorm:"not null; default:null"`
	Alias            *string    `json:"alias"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	SessionId        *string    `json:"session_id"`
	SessionExpiresAt *time.Time `json:"session_expires_at"`
}

type Session struct {
	ID        string    `json:"session_id" gorm:"default:null"`
	UserID    string    `json:"user_id" gorm:"not null; default:null"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Category struct {
	ID        string    `json:"category_id" gorm:"default:null"`
	UserID    string    `json:"user_id" gorm:"not null; default:null"`
	Name      string    `json:"name" binding:"required" gorm:"not null; default:null"`
	Desc      *string   `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Cron struct {
	ID         string    `json:"cron_id" gorm:"default:null"`
	UserID     string    `json:"user_id" gorm:"not null; default:null"`
	CategoryID *string   `json:"category_id"`
	Schedule   string    `json:"schedule" binding:"required" gorm:"not null; default:null"`
	Name       string    `json:"name" binding:"required" gorm:"not null; default:null"`
	Desc       *string   `json:"desc"`
	Enabled    bool      `json:"enabled" binding:"required" gorm:"not null; default:null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Task struct {
	ID          string     `json:"task_id" gorm:"default:null"`
	Name        string     `json:"name" binding:"required" gorm:"not null; default:null"`
	UserID      string     `json:"user_id" gorm:"not null; default:null"`
	CategoryID  *string    `json:"category_id"`
	CronID      *string    `json:"cron_id"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}
