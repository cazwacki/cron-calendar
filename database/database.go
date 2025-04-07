package database

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

var connection *gorm.DB
var databaseContext = log.With().Str("file", "database.go")

func getConnection() *gorm.DB {
	functionLogger := databaseContext.Str("function", "getConnection").Logger()
	functionLogger.Debug().Msg("invoked")

	if connection == nil {
		var err error
		connection, err = gorm.Open(sqlite.Open("cron-calendar.db"), &gorm.Config{})
		if err != nil {
			functionLogger.Err(err).Msg("failed to open a database connection; returning nil")
			return nil
		}
	}

	functionLogger.Debug().Msg("returning with success")
	return connection
}

func InitializeTables() error {
	functionLogger := databaseContext.Str("function", "InitializeTables").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	err := db.AutoMigrate(&User{}, &Session{}, &Category{}, &Cron{}, &Task{})
	if err != nil {
		functionLogger.Err(err).Msg("failed to initialize tables")
		return err
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func InsertUser(user User) error {
	functionLogger := databaseContext.Str("function", "InsertUser").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	result := db.Create(&user)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func GetUserById(userId string) (User, error) {
	functionLogger := databaseContext.Str("function", "GetUserById").Logger()
	functionLogger.Debug().Msg("invoked")

	var user User
	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return user, errors.New("failed to get database connection")
	}
	result := db.First(&user, "id = ?", userId)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return user, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return user, nil
}

func DeleteUserById(userId string) error {
	functionLogger := databaseContext.Str("function", "DeleteUserById").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	result := db.Where("id = ?", userId).Delete(&User{})
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func UpsertSession(session Session) error {
	functionLogger := databaseContext.Str("function", "UpsertSession").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&session)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func GetSessionById(sessionId string) (Session, error) {
	functionLogger := databaseContext.Str("function", "GetSessionById").Logger()
	functionLogger.Debug().Msg("invoked")

	var session Session
	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return session, errors.New("failed to get database connection")
	}
	result := db.First(&session, "id = ?", sessionId)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return session, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return session, nil
}

func UpsertCategory(category Category) error {
	functionLogger := databaseContext.Str("function", "UpsertCategory").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	// todo validate the user owns the existing category or fail
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&category)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func GetCategoryByUserAndId(userId string, categoryId string) (Category, error) {
	functionLogger := databaseContext.Str("function", "GetCategoryByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	var category Category
	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return category, errors.New("failed to get database connection")
	}
	result := db.Where("user_id = ?", userId).First(&category, "id = ?", categoryId)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return category, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return category, nil
}

func DeleteCategoryByUserAndId(userId string, categoryId string) error {
	functionLogger := databaseContext.Str("function", "DeleteCategoryByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	result := db.Where("user_id = ?", userId).Where("id = ?", categoryId).Delete(&Category{})
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func GetCategoriesByUserId(userId string) ([]Category, error) {
	functionLogger := databaseContext.Str("function", "GetCategoriesByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	var categories []Category
	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return categories, errors.New("failed to get database connection")
	}
	result := db.Where("user_id = ?", userId).Find(&categories)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return categories, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return categories, nil
}

func UpsertCron(cron Cron) error {
	functionLogger := databaseContext.Str("function", "UpsertCron").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	// todo validate the user owns the existing cron or fail
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&cron)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func GetCronByUserAndId(userId string, cronId string) (Cron, error) {
	functionLogger := databaseContext.Str("function", "GetCronByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	var cron Cron
	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return cron, errors.New("failed to get database connection")
	}
	result := db.Where("user_id = ?", userId).First(&cron, "id = ?", cronId)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return cron, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return cron, nil
}

func DeleteCronByUserAndId(userId string, cronId string) error {
	functionLogger := databaseContext.Str("function", "DeleteCronByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	result := db.Where("user_id = ?", userId).Where("id = ?", cronId).Delete(&Cron{})
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func GetCronsByUserId(userId string) ([]Cron, error) {
	functionLogger := databaseContext.Str("function", "GetCronsByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	var crons []Cron
	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return crons, errors.New("failed to get database connection")
	}
	result := db.Where("user_id = ?", userId).Find(&crons)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return crons, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return crons, nil
}

func GetAllCrons() ([]Cron, error) {
	functionLogger := databaseContext.Str("function", "GetAllCrons").Logger()
	functionLogger.Debug().Msg("invoked")

	var crons []Cron
	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return crons, errors.New("failed to get database connection")
	}
	result := db.Find(&crons)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return crons, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return crons, nil
}

func UpsertTask(task Task) error {
	functionLogger := databaseContext.Str("function", "UpsertTask").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	// todo validate the user owns the existing task or fail
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&task)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func GetTaskByUserAndId(userId string, taskId string) (Task, error) {
	functionLogger := databaseContext.Str("function", "GetTaskByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	var task Task
	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return task, errors.New("failed to get database connection")
	}
	result := db.Where("user_id = ?", userId).First(&task, "id = ?", taskId)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return task, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return task, nil
}

func DeleteTaskByUserAndId(userId string, taskId string) error {
	functionLogger := databaseContext.Str("function", "DeleteTaskByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	result := db.Where("user_id = ?", userId).Where("id = ?", taskId).Delete(&Task{})
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func GetTasksByUserId(userId string) ([]Task, error) {
	functionLogger := databaseContext.Str("function", "GetTasksByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	var tasks []Task
	db := getConnection()
	if db == nil {
		functionLogger.Error().Msg("failed to get database connection")
		return tasks, errors.New("failed to get database connection")
	}
	result := db.Where("user_id = ?", userId).Find(&tasks)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return tasks, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return tasks, nil
}
