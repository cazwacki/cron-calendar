package database

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	ID               string     `json:"username" binding:"required"`
	Password         string     `json:"password" binding:"required"`
	Alias            *string    `json:"alias"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	SessionId        *string    `json:"session_id"`
	SessionExpiresAt *time.Time `json:"session_expires_at"`
}

type Session struct {
	ID        string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Category struct {
	ID        string    `json:"category_id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name" binding:"required"`
	Desc      *string   `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Cron struct {
	ID         string    `json:"cron_id"`
	UserID     string    `json:"user_id"`
	CategoryID *string   `json:"category_id"`
	Schedule   string    `json:"schedule" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	Desc       *string   `json:"desc"`
	Enabled    bool      `json:"enabled" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Task struct {
	ID          string     `json:"task_id" binding:"required"`
	UserID      string     `json:"user_id"`
	CategoryID  *string    `json:"category_id"`
	CronID      *string    `json:"cron_id"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed" binding:"required"`
}

var db *gorm.DB

func getConnection() *gorm.DB {
	if db == nil {
		var err error
		db, err = gorm.Open(sqlite.Open("cron-calendar.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	}
	return db
}

func InitializeTables() {
	db := getConnection()
	db.AutoMigrate(&User{}, &Session{}, &Category{}, &Cron{}, &Task{})
}

func InsertUser(user User) error {
	db := getConnection()
	result := db.Create(&user)
	return result.Error
}

func GetUserById(userId string) (User, error) {
	db := getConnection()
	var user User
	result := db.First(&user, "id = ?", userId)
	return user, result.Error
}

func DeleteUserById(userId string) error {
	db := getConnection()
	result := db.Where("id = ?", userId).Delete(&User{})
	return result.Error
}

func UpsertSession(session Session) error {
	db := getConnection()
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&session)
	return result.Error
}

func GetSessionById(sessionId string) (Session, error) {
	db := getConnection()
	var session Session
	result := db.First(&session, "id = ?", sessionId)
	return session, result.Error
}

func UpsertCategory(category Category) error {
	db := getConnection()
	// todo validate the user owns the existing category or fail
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&category)
	return result.Error
}

func GetCategoryByUserAndId(userId string, categoryId string) (Category, error) {
	db := getConnection()
	var category Category
	result := db.Where("user_id = ?", userId).First(&category, "id = ?", categoryId)
	return category, result.Error
}

func DeleteCategoryByUserAndId(userId string, categoryId string) error {
	db := getConnection()
	result := db.Where("user_id = ?", userId).Where("id = ?", categoryId).Delete(&Category{})
	return result.Error
}

func GetCategoriesByUserId(userId string) ([]Category, error) {
	var categories []Category
	result := db.Where("user_id = ?", userId).Find(&categories)
	return categories, result.Error
}

func UpsertCron(cron Cron) error {
	db := getConnection()
	// todo validate the user owns the existing cron or fail
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&cron)
	return result.Error
}

func GetCronByUserAndId(userId string, cronId string) (Cron, error) {
	db := getConnection()
	var cron Cron
	result := db.Where("user_id = ?", userId).First(&cron, "id = ?", cronId)
	return cron, result.Error
}

func DeleteCronByUserAndId(userId string, cronId string) error {
	db := getConnection()
	result := db.Where("user_id = ?", userId).Where("id = ?", cronId).Delete(&Cron{})
	return result.Error
}

func GetCronsByUserId(userId string) ([]Cron, error) {
	var crons []Cron
	result := db.Where("user_id = ?", userId).Find(&crons)
	return crons, result.Error
}

func GetAllCrons() ([]Cron, error) {
	var crons []Cron
	result := db.Find(&crons)
	return crons, result.Error
}

func UpsertTask(task Task) error {
	db := getConnection()
	// todo validate the user owns the existing task or fail
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&task)
	return result.Error
}

func GetTaskByUserAndId(userId string, taskId string) (Task, error) {
	db := getConnection()
	var task Task
	result := db.Where("user_id = ?", userId).First(&task, "id = ?", taskId)
	return task, result.Error
}

func DeleteTaskByUserAndId(userId string, taskId string) error {
	db := getConnection()
	result := db.Where("user_id = ?", userId).Where("id = ?", taskId).Delete(&Task{})
	return result.Error
}

func GetTasksByUserId(userId string) ([]Task, error) {
	var tasks []Task
	result := db.Where("user_id = ?", userId).Find(&tasks)
	return tasks, result.Error
}
