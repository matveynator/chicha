package Models

import (
	"time"

	"gorm.io/gorm"
)

// Database locator in memory (GORM is calling by Models.DB)
var DB *gorm.DB

// Laps main data of the race
type Lap struct {
	//gorm.Model
	ID                    uint      `gorm:"primaryKey"`
	OwnerID               uint      `gorm:"index" json:"owner_id"`
	TagID                 string    `gorm:"char(80);index" json:"tag_id" xml:"TagID"`
	DiscoveryUnixTime     int64     `gorm:"index" json:"discovery_unix_time"`
	DiscoveryAverageUnixTime int64     `gorm:"index" json:"discovery_average_unix_time"`
	DiscoveryTime         string    `json:"-" xml:"DiscoveryTime"`
	DiscoveryTimePrepared time.Time `json:"discovery_time"`
	DiscoveryAverageTimePrepared time.Time `json:"discovery_average_time"`
	AverageResults        string       `json:"average_results"` 
	AverageResultsCount   uint       `gorm:"index" json:"average_results_count"`

	Antenna               uint8     `gorm:"index" json:"antenna" xml:"Antenna"`
	AntennaIP             string    `gorm:"char(128);index" json:"antenna_ip"`
	UpdatedAt             time.Time `json:"updated_at"`
	RaceID                uint      `gorm:"index" json:"race_id"`
	CurrentRacePosition   uint      `gorm:"index" json:"current_race_postition"`
	TimeBehindTheLeader   int64     `gorm:"index" json:"time_behind_the_leader"`
	LapNumber             int       `gorm:"index" json:"lap_number"`
	LapTime               int64     `gorm:"index" json:"lap_time"`
	LapPosition           uint      `gorm:"index" json:"lap_postition"`
	LapIsCurrent          int       `gorm:"index" json:"lap_is_current"`
	LapIsStrange          int       `gorm:"index" json:"lap_is_strange"`
	StageFinished         int       `gorm:"index" json:"stage_finished"`
	BestLapTime           int64     `gorm:"index" json:"best_lap_time"`
	BestLapNumber   	    int       `gorm:"index" json:"best_lap_number"`
	BestLapPosition       uint      `gorm:"index" json:"best_lap_postition"`
	RaceTotalTime         int64     `gorm:"index" json:"race_total_time"`
	BetterOrWorseLapTime  int64     `gorm:"index" json:"better_or_worse_lap_time"`
}

// Laps time labels table name
func (u *Lap) TableName() string {
	return "laps"
}
