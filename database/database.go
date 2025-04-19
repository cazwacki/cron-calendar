package database

import (
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var databaseContext = log.With().Str("file", "database.go")

type Database struct {
	connection *gorm.DB
}

type CategoryDB interface {
	UpsertCategoryIfOwner(category Category) error
	GetCategoryByIdAndUserId(id string, userId string) (*Category, error)
	DeleteCategoryByIdAndUserId(id string, userId string) error
	SearchCategoriesByUserId(userId string) ([]Category, error)
}

type CronDB interface {
	UpsertCronIfOwner(cron Cron) error
	GetCronByIdAndUserId(id string, userId string) (*Cron, error)
	DeleteCronByIdAndUserId(id string, userId string) error
	SearchCronsByUserId(userId string) ([]Cron, error)
	GetAllCrons() ([]Cron, error)
}

type TaskDB interface {
	UpsertTaskIfOwner(task Task) error
	GetTaskByIdAndUserId(id string, userId string) (*Task, error)
	DeleteTaskByIdAndUserId(id string, userId string) error
	SearchTasksByUserId(userId string) ([]Task, error)
}

type UserDB interface {
	InsertUser(user User) error
	GetUserById(id string) (*User, error)
	DeleteUserById(id string) error
}

type SessionDB interface {
	UpsertSessionIfOwner(session Session) error
	GetSessionById(id string) (*Session, error)
}

func (d *Database) establishConnection() error {
	functionLogger := databaseContext.Str("function", "getConnection").Logger()
	functionLogger.Debug().Msg("invoked")

	if d.connection == nil {
		functionLogger.Debug().Msg("opening a new database connection")
		connection, err := gorm.Open(sqlite.Open("cron-calendar.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			functionLogger.Err(err).Msg("failed to open a database connection; returning nil")
			return err
		}
		functionLogger.Debug().Msg("new connection established")
		d.connection = connection
	} else {
		functionLogger.Debug().Msg("returning existing connection")
	}
	return nil
}

func (d *Database) InitializeTables() error {
	functionLogger := databaseContext.Str("function", "InitializeTables").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}
	err = d.connection.AutoMigrate(&User{}, &Session{}, &Category{}, &Cron{}, &Task{})
	if err != nil {
		functionLogger.Err(err).Msg("failed to initialize tables")
		return err
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) UpsertCategoryIfOwner(category Category) error {
	functionLogger := databaseContext.Str("function", "UpsertCategoryIfOwner").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	var existingCategory Category
	result := d.connection.First(&existingCategory, "id = ?", category.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}
	if existingCategory.UserID != "" && existingCategory.UserID != category.UserID {
		functionLogger.Warn().Msg("attempt to modify another user's category")
		return errors.New("constraint: attempt to modify another user's category")
	}

	result = d.connection.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&category)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) GetCategoryByIdAndUserId(id string, userId string) (*Category, error) {
	functionLogger := databaseContext.Str("function", "GetCategoryByIdAndUserId").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return nil, errors.New("failed to get database connection")
	}

	var category Category
	result := d.connection.Where("user_id = ?", userId).First(&category, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			functionLogger.Debug().Msg("no category found")
			return nil, nil
		}
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return nil, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return &category, nil
}

func (d *Database) DeleteCategoryByIdAndUserId(id string, userId string) error {
	functionLogger := databaseContext.Str("function", "DeleteCategoryByIdAndUserId").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	result := d.connection.Where("user_id = ?", userId).Where("id = ?", id).Delete(&Category{})
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) SearchCategoriesByUserId(userId string) ([]Category, error) {
	functionLogger := databaseContext.Str("function", "SearchCategoriesByUserId").Logger()
	functionLogger.Debug().Msg("invoked")

	var categories []Category
	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return categories, errors.New("failed to get database connection")
	}

	result := d.connection.Where("user_id = ?", userId).Find(&categories)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return categories, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return categories, nil
}

func (d *Database) UpsertCronIfOwner(cron Cron) error {
	functionLogger := databaseContext.Str("function", "UpsertCronIfOwner").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	var existingCron Cron
	result := d.connection.First(&existingCron, "id = ?", cron.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}
	if existingCron.UserID != "" && existingCron.UserID != cron.UserID {
		functionLogger.Warn().Msg("attempt to modify another user's cron")
		return errors.New("constraint: attempt to modify another user's cron")
	}

	result = d.connection.Save(&cron)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) GetCronByIdAndUserId(id string, userId string) (*Cron, error) {
	functionLogger := databaseContext.Str("function", "GetCronByIdAndUserId").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return nil, errors.New("failed to get database connection")
	}

	var cron Cron
	result := d.connection.Where("user_id = ?", userId).First(&cron, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			functionLogger.Debug().Msg("no cron found")
			return nil, nil
		}
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return nil, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return &cron, nil
}

func (d *Database) DeleteCronByIdAndUserId(id string, userId string) error {
	functionLogger := databaseContext.Str("function", "DeleteCronByIdAndUserId").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	result := d.connection.Where("user_id = ?", userId).Where("id = ?", id).Delete(&Cron{})
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) SearchCronsByUserId(userId string) ([]Cron, error) {
	functionLogger := databaseContext.Str("function", "GetCronsByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	var crons []Cron
	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return crons, errors.New("failed to get database connection")
	}

	result := d.connection.Where("user_id = ?", userId).Find(&crons)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return crons, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return crons, nil
}

func (d *Database) GetAllCrons() ([]Cron, error) {
	functionLogger := databaseContext.Str("function", "GetAllCrons").Logger()
	functionLogger.Debug().Msg("invoked")

	var crons []Cron
	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return crons, errors.New("failed to get database connection")
	}

	result := d.connection.Find(&crons)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return crons, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return crons, nil
}

func (d *Database) UpsertTaskIfOwner(task Task) error {
	functionLogger := databaseContext.Str("function", "UpsertTask").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	var existingTask Task
	result := d.connection.First(&existingTask, "id = ?", task.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}
	if existingTask.UserID != "" && existingTask.UserID != task.UserID {
		functionLogger.Warn().Msg("attempt to modify another user's task")
		return errors.New("constraint: attempt to modify another user's task")
	}

	result = d.connection.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&task)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) GetTaskByIdAndUserId(id string, userId string) (*Task, error) {
	functionLogger := databaseContext.Str("function", "GetTaskByUserAndId").Logger()
	functionLogger.Debug().Msg("invoked")

	var task Task
	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return nil, errors.New("failed to get database connection")
	}

	result := d.connection.Where("user_id = ?", userId).First(&task, "id = ?", taskId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			functionLogger.Debug().Msg("no task found")
			return nil, nil
		}
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return nil, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return &task, nil
}

func (d *Database) DeleteTaskByIdAndUserId(id string, userId string) error {
	functionLogger := databaseContext.Str("function", "DeleteTaskByIdAndUserId").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	result := d.connection.Where("user_id = ?", userId).Where("id = ?", taskId).Delete(&Task{})
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) GetTasksByUserId(userId string) ([]Task, error) {
	functionLogger := databaseContext.Str("function", "GetTasksByUserId").Logger()
	functionLogger.Debug().Msg("invoked")

	var tasks []Task
	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return tasks, errors.New("failed to get database connection")
	}

	result := d.connection.Where("user_id = ?", userId).Find(&tasks)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return tasks, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return tasks, nil
}

func (d *Database) InsertUser(user User) error {
	functionLogger := databaseContext.Str("function", "UpsertUserIfOwner").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	result := d.connection.Create(&user)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) GetUserById(userId string) (*User, error) {
	functionLogger := databaseContext.Str("function", "GetUserById").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return nil, errors.New("failed to get database connection")
	}

	var user User
	result := d.connection.First(&user, "id = ?", userId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			functionLogger.Debug().Msg("no user found")
			return nil, nil
		}
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return nil, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return &user, nil
}

func (d *Database) DeleteUserById(userId string) error {
	functionLogger := databaseContext.Str("function", "DeleteUserById").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	result := d.connection.Where("id = ?", userId).Delete(&User{})
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) UpsertSession(session Session) error {
	functionLogger := databaseContext.Str("function", "UpsertSession").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	result := d.connection.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&session)
	if result.Error != nil {
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return nil
}

func (d *Database) GetSessionById(sessionId string) (*Session, error) {
	functionLogger := databaseContext.Str("function", "GetSessionById").Logger()
	functionLogger.Debug().Msg("invoked")

	err := d.establishConnection()
	if err != nil {
		functionLogger.Error().Msg("failed to get database connection")
		return errors.New("failed to get database connection")
	}

	var session Session
	result := d.connection.First(&session, "id = ?", sessionId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			functionLogger.Debug().Msg("no session found")
			return nil, nil
		}
		functionLogger.Err(result.Error).Msg("failed to perform db operation")
		return nil, result.Error
	}

	functionLogger.Debug().Msg("returning with success")
	return &session, nil
}
