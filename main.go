package main

import (
	"log"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/irumaru/encoding-control-server/config"
	"github.com/irumaru/encoding-control-server/scheduler"
)

type Runner struct {
	ID        int       `json:"ID" goam:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"CreatedAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"UpdatedAt" gorm:"autoUpdateTime"`
	Priority  int       `json:"Priority"`
	Status    string    `json:"Status" gorm:"default:Drop"`
	Name      string    `json:"Name"`
}

type RunnerRequest struct {
	Priority int    `json:"Priority"`
	Status   string `json:"Status"`
	Name     string `json:"Name"`
}

type Job struct {
	ID        int       `json:"ID" goam:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"CreatedAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"UpdatedAt" gorm:"autoUpdateTime"`
	Priority  int       `json:"Priority"`
	Status    string    `json:"Status" gorm:"default:Waiting"`
	Kind      string    `json:"Kind"`
	Option    string    `json:"Option"`
	RunnerID  int       `json:"RunnerId"`
	Name      string    `json:"Name"`
}

type JobRequest struct {
	Priority int    `json:"Priority"`
	Status   string `json:"Status"`
	Kind     string `json:"Kind"`
	Option   string `json:"Option"`
	Name     string `json:"Name"`
}

func dbInit() *gorm.DB {
	user := config.Cfg.Db_user
	password := config.Cfg.Db_password
	host := config.Cfg.Db_host
	port := config.Cfg.Db_port
	dbname := config.Cfg.Db_name

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

var db *gorm.DB

func main() {
	// LoadConfig
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("failed load config: %v", err)
	}

	// DB connect
	db = dbInit()

	db.AutoMigrate(&Runner{})
	db.AutoMigrate(&Job{})

	// Echo instance
	e := echo.New()

	// Scheduler
	scheduler.ControllerStart(db)

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	//
	// Runner API
	//
	e.GET("/api/v1/runner", func(c echo.Context) error {
		runners := []Runner{}
		if result := db.Find(&runners); result.Error != nil {
			log.Fatalf("failed query runner: %v", result.Error)
		}

		return c.JSON(200, runners)
	})
	e.POST("/api/v1/runner", func(c echo.Context) error {
		// Bind
		runnerReq := new(RunnerRequest)
		if err := c.Bind(runnerReq); err != nil {
			log.Fatalf("failed bind runner: %v", err)
		}

		// Create
		runner := Runner{}
		runner.Priority = runnerReq.Priority
		runner.Status = runnerReq.Status
		runner.Name = runnerReq.Name

		if result := db.Create(&runner); result.Error != nil || result.RowsAffected != 1 {
			log.Fatalf("failed create runner: %v", result.Error)
		}

		return c.JSON(200, runner)
	})
	e.POST("/api/v1/runner/:id", func(c echo.Context) error {
		// Bind
		runnerReq := new(RunnerRequest)
		if err := c.Bind(runnerReq); err != nil {
			log.Fatalf("failed bind runner: %v", err)
		}

		// Id
		idStr := c.Param("id")
		runnerId, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("failed convert id: %v", err)
		}

		// Create
		runner := Runner{}
		runner.ID = runnerId
		runner.Priority = runnerReq.Priority
		runner.Status = runnerReq.Status
		runner.Name = runnerReq.Name

		if result := db.Model(&runner).Updates(runner); result.Error != nil || result.RowsAffected != 1 {
			log.Fatalf("failed create runner: %v", result.Error)
		}
		return c.JSON(200, runner)
	})
	e.DELETE("/api/v1/runner/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("failed convert id: %v", err)
		}

		// Delete
		runner := Runner{}
		runner.ID = id

		if result := db.Delete(&runner); result.Error != nil || result.RowsAffected != 1 {
			log.Fatalf("failed delete runner: %v", result.Error)
		}

		return c.JSON(200, runner)
	})

	//
	// Job API
	//
	e.GET("/api/v1/job", func(c echo.Context) error {
		job := []Job{}
		if result := db.Find(&job); result.Error != nil {
			log.Fatalf("failed query runner: %v", result.Error)
		}

		return c.JSON(200, job)
	})
	e.POST("/api/v1/job", func(c echo.Context) error {
		// Bind
		jobReq := new(JobRequest)
		if err := c.Bind(jobReq); err != nil {
			log.Fatalf("failed bind runner: %v", err)
		}

		// Create
		job := Job{}
		job.Priority = jobReq.Priority
		job.Kind = jobReq.Kind
		job.Option = jobReq.Option
		job.Name = jobReq.Name

		if result := db.Create(&job); result.Error != nil || result.RowsAffected != 1 {
			log.Fatalf("failed create runner: %v", result.Error)
		}

		return c.JSON(200, job)
	})
	e.POST("/api/v1/job/:id", func(c echo.Context) error {
		// Bind
		jobReq := new(JobRequest)
		if err := c.Bind(jobReq); err != nil {
			log.Fatalf("failed bind runner: %v", err)
		}

		// Id
		idStr := c.Param("id")
		jobId, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("failed convert id: %v", err)
		}

		// Create
		job := Job{}
		job.ID = jobId
		job.Priority = jobReq.Priority
		job.Kind = jobReq.Kind
		job.Option = jobReq.Option
		job.Name = jobReq.Name

		if result := db.Model(&job).Updates(job); result.Error != nil || result.RowsAffected != 1 {
			log.Fatalf("failed create runner: %v", result.Error)
		}
		return c.JSON(200, job)
	})
	e.DELETE("/api/v1/job/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("failed convert id: %v", err)
		}

		// Delete
		job := Job{}
		job.ID = id

		if result := db.Delete(&job); result.Error != nil || result.RowsAffected != 1 {
			log.Fatalf("failed delete runner: %v", result.Error)
		}

		return c.JSON(200, job)
	})

	//
	// Virtual API
	//
	e.GET("/api/v1/runner/:id/job", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("failed convert id: %v", err)
		}

		job := Job{}
		if result := db.Where("runner_id = ?", id).Find(&job); result.Error != nil {
			log.Fatalf("failed query runner: %v", result.Error)
		}

		return c.JSON(200, job)
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
