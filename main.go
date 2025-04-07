package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/adhocore/gronx"
	"github.com/adhocore/gronx/pkg/tasker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"cron-calendar/database"
	"cron-calendar/handlers"
	"cron-calendar/middleware"
)

var fileContext = log.With().Str("file", "main.go")

func generateTasks(ctx context.Context) (int, error) {
	functionLogger := fileContext.Str("function", "generateTasks").Logger()
	functionLogger.Debug().Msg("invoked")

	crons, err := database.GetAllCrons()
	if err != nil {
		functionLogger.Err(err).Msg("failed to get crons from db")
		return 1, err
	}

	for _, cron := range crons {
		completeExpression := fmt.Sprintf("0 0 0 %s", cron.Schedule)
		now := time.Now()
		nowTruncated := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		taskDue, err := gronx.New().IsDue(completeExpression, nowTruncated)
		if err != nil {
			functionLogger.Warn().AnErr("err", err).Msgf("failed to evaluate cron %s's expression: %s; skipping", cron.ID, cron.Schedule)
			continue
		}
		shouldCreateTask := taskDue && cron.Enabled
		if shouldCreateTask {
			var newTask database.Task
			functionLogger.Debug().Msgf("cron %s is due for task creation", cron.ID)
			uuid, err := uuid.NewV7()
			if err != nil {
				functionLogger.Warn().AnErr("err", err).Msgf("failed to generate uuid for cron %s; skipping", cron.ID)
				continue
			}
			newTask.ID = uuid.String()
			newTask.Name = cron.Name
			newTask.UserID = cron.UserID
			newTask.CategoryID = cron.CategoryID
			newTask.CronID = &cron.ID
			fmt.Printf("%+v\n", newTask)
			if err := database.UpsertTask(newTask); err != nil {
				functionLogger.Err(err).Msg("failed to write back to database")
			}
		}
	}

	functionLogger.Debug().Msg("finished processing crons")
	return 0, err
}

func startTasker() {
	functionLogger := fileContext.Str("function", "startTasker").Logger()
	functionLogger.Debug().Msg("invoked")

	cronHour := os.Getenv("LIST_GEN_HOUR")
	tasker := tasker.New(tasker.Option{
		Verbose: false,
	})
	tasker.Task(fmt.Sprintf("0 %s * * *", cronHour), generateTasks)
	functionLogger.Debug().Msg("running task manager")
	tasker.Run()
	functionLogger.Debug().Msg("complete")
}

func startHttpService() {
	functionLogger := fileContext.Str("function", "startHttpService").Logger()
	functionLogger.Debug().Msg("invoked")

	router := gin.Default()
	router.Use(gin.CustomRecovery(middleware.GenerateErrorResponse))

	router.PUT("/user/register", handlers.RegisterUser)
	router.POST("/user/login", handlers.Login)

	authNeeded := router.Group("/")
	authNeeded.Use(middleware.ValidateAuthorization())
	{
		authNeeded.DELETE("/user/destroy", handlers.DestroyUser)
		// category management
		authNeeded.POST("/category", handlers.GetAllCategories)
		authNeeded.PUT("/category/:categoryId", handlers.PutCategory)
		authNeeded.GET("/category/:categoryId", handlers.GetCategory)
		authNeeded.DELETE("/category/:categoryId", handlers.DeleteCategory)
		// cron management
		authNeeded.POST("/cron", handlers.GetAllCrons)
		authNeeded.PUT("/cron/:cronId", handlers.PutCron)
		authNeeded.GET("/cron/:cronId", handlers.GetCron)
		authNeeded.DELETE("/cron/:cronId", handlers.DeleteCron)
		// task management
		authNeeded.POST("/task", handlers.GetAllTasks)
		authNeeded.PUT("/task/:taskId", handlers.PutTask)
		authNeeded.GET("/task/:taskId", handlers.GetTask)
		authNeeded.DELETE("/task/:taskId", handlers.DeleteTask)
	}
	functionLogger.Debug().Msg("listening")
	router.Run()
	functionLogger.Debug().Msg("complete")
}

func main() {
	functionLogger := fileContext.Str("function", "main").Logger()
	functionLogger.Debug().Msg("started")
	database.InitializeTables()
	go startTasker()
	startHttpService()
	functionLogger.Debug().Msg("complete")
}
