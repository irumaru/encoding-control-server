package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

type Job struct {
	ID           uint64    `goam:"primaryKey;autoIncrement"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	Kind         string
	Priority     int8
	Status       string
	Client       string
	ExcutionNode string
	JobDetail    []JobDetail `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE;constraint:OnDelete:CASCADE"`
}

type JobDetail struct {
	ID        uint64 `goam:"primaryKey;autoIncrement"`
	JobID     uint64
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Status    string
	Detail    string
}

var db *gorm.DB

func main() {
	// DB connect
	db = dbInit()

	// DB migrate
	db.AutoMigrate(Job{})
	db.AutoMigrate(JobDetail{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/job/create", CreateJob)
	e.POST("/job/execution", ExecutionJob)
	e.GET("/job/get", GetJob)
	e.POST("/job/:id/update", UpdateJob)
	e.GET("/job/:id/get", GetJobDetails)
	e.DELETE("/job/:id/delete", DeleteJob)

	e.Logger.Fatal(e.Start(":8080"))
}

type JobPost struct {
	Kind     string `json:"kind" validate:"required"`
	Priority int8   `json:"priority" validate:"required"`
	Client   string `json:"client" validate:"required"`
}

func CreateJob(c echo.Context) error {
	post := new(JobPost)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, "Post error.")
	}
	if err := c.Validate(post); err != nil {
		return c.JSON(http.StatusBadRequest, "Json Validate error.")
	}

	job := Job{}
	job.Kind = post.Kind
	job.Priority = post.Priority
	job.Client = post.Client
	job.Status = "wait"

	if result := db.Create(&job); result.Error != nil {
		return c.JSON(http.StatusBadRequest, "Database error.")
	}

	return c.JSON(http.StatusOK, job)
}

type ExecutionPost struct {
	Kind         string `json:"kind" validate:"required"`
	ExcutionNode string `json:"execution_node" validate:"required"`
}

func ExecutionJob(c echo.Context) error {
	post := new(ExecutionPost)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, "Post error.")
	}
	if err := c.Validate(post); err != nil {
		return c.JSON(http.StatusBadRequest, "Json Validate error.")
	}

	job := Job{}

	result := db.Limit(1).Find(&job, `status = "wait" AND kind = ?`, post.Kind)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, "Database error.")
	}

	// 実行可能なジョブなし
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusOK, "No job.")
	}

	// ジョブ割り当て
	job.Status = "processing"
	job.ExcutionNode = post.ExcutionNode
	if result := db.Save(job); result.Error != nil {
		return c.JSON(http.StatusBadRequest, "Database error.")
	}

	return c.JSON(http.StatusOK, job)
}

func GetJob(c echo.Context) error {
	jobs := []Job{}
	if result := db.Find(&jobs); result.Error != nil {
		return c.JSON(http.StatusBadRequest, "Database error.")
	}

	return c.JSON(http.StatusOK, jobs)
}

type UpdatePost struct {
	Status string `json:"status" validate:"required"`
	Detail string `json:"detail" validate:"required"`
}

func UpdateJob(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID parse error.")
	}

	post := new(UpdatePost)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, "Post error.")
	}
	if err := c.Validate(post); err != nil {
		return c.JSON(http.StatusBadRequest, "Json Validate error.")
	}

	//詳細のアップデート
	jobD := JobDetail{}
	jobD.JobID = id
	jobD.Status = post.Status
	jobD.Detail = post.Detail
	if result := db.Create(&jobD); result.Error != nil {
		fmt.Println(result.Error)
		return c.JSON(http.StatusBadRequest, "Database error.")
	}

	//Jobのアップデート
	if post.Status == "Complete" || post.Status == "Failed" {
		if result := db.Model(&Job{}).Where("id = ?", id).Update("status", post.Status); result.Error != nil {
			return c.JSON(http.StatusBadRequest, "Database error.")
		}
	}

	return c.JSON(http.StatusOK, jobD)
}

func GetJobDetails(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID parse error.")
	}

	jobDs := []JobDetail{}
	if result := db.Find(&jobDs, "job_id = ?", id); result.Error != nil {
		return c.JSON(http.StatusBadRequest, "Database error.")
	}

	return c.JSON(http.StatusOK, jobDs)
}

func DeleteJob(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID parse error.")
	}

	result := db.Delete(&Job{}, id)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, "Database error.")
	}
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, "Not found target job.")
	}

	return c.JSON(http.StatusOK, "Delete OK.")
}

func dbInit() *gorm.DB {
	dsn := "test_user:test_password@tcp(db:3306)/test_db?charset=utf8mb4&parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
