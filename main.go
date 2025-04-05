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

	"cron-calendar/database"
	"cron-calendar/handlers"
	"cron-calendar/middleware"
)

func main() {
	database.InitializeTables()
	go startTasker()
	startHttpService()
}

func startTasker() {
	cronHour := os.Getenv("LIST_GEN_HOUR")
	tasker := tasker.New(tasker.Option{
		Verbose: false,
	})
	tasker.Task(fmt.Sprintf("0 %s * * *", cronHour), func(ctx context.Context) (int, error) {
		crons, err := database.GetAllCrons()
		if err != nil {
			tasker.Log.Printf("error getting crons: %s", err.Error())
			return 1, err
		}
		for _, cron := range crons {
			completeExpression := fmt.Sprintf("0 0 0 %s", cron.Schedule)
			now := time.Now()
			nowTruncated := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			shouldCreateTask, err := gronx.New().IsDue(completeExpression, nowTruncated)
			if err != nil {
				fmt.Println("failed to check cron validity")
				// proper error logging
			}
			if shouldCreateTask {
				var newTask database.Task
				uuid, err := uuid.NewV7()
				if err != nil {
					fmt.Println("failed to generate new uuid")
					// proper error logging
				}
				newTask.ID = uuid.String()
				newTask.UserID = cron.UserID
				newTask.CategoryID = cron.CategoryID
				newTask.CronID = &cron.ID
				database.UpsertTask(newTask)
			}
		}
		return 0, err
	})
	tasker.Run()
}

func startHttpService() {
	router := gin.Default()

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
	router.Run()
}
