package scheduler

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func ControllerStart(db *gorm.DB) {
	// Databaae init
	go func(db *gorm.DB) {
		for {
			Controller(db)
			time.Sleep(10 * time.Second)
		}
	}(db)
}

type Runner struct {
	ID        int       `json:"ID" goam:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"CreatedAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"UpdatedAt" gorm:"autoUpdateTime"`
	Priority  int       `json:"Priority"`
	Status    string    `json:"Status" gorm:"default:Drop"`
	Name      string    `json:"Name"`
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

func Controller(db *gorm.DB) {
	// priorityは数字が大きい方が優先度が高い

	// Schedule可能なRunnerを探す
	runner := []Runner{}
	if result := db.Raw("SELECT * FROM runners WHERE status = 'Ready' AND id NOT IN (SELECT runner_id FROM jobs WHERE status = 'Scheduled' OR status = 'Running') ORDER BY priority DESC LIMIT 1").
		Scan(&runner); result.Error != nil {
		log.Printf("Failed to get runner: %v", result.Error)
	}

	// 実行可能なRunnerがない場合は終了
	if len(runner) == 0 {
		return
	}

	// Schedule可能なJobを探し、実行する
	result := db.Model(&Job{}).Where("status = 'Waiting'").Order("priority DESC").Limit(1).Updates(Job{Status: "Scheduled", RunnerID: runner[0].ID})
	if result.Error != nil {
		log.Fatalf("Failed to update job: %v", result.Error)
	}

	// Schedule可能なRunnerはあるが、Schedule可能なJobがない場合は終了
	if result.RowsAffected == 0 {
		return
	}

	// JobをSchedule時はログに出力
	log.Printf("Scheduled %d job with an ID of %d", result.RowsAffected, runner[0].ID)
}
