package Models

import ( 
//	"fmt"
	"time" 
)

// Return all laps in system order by date
func GetAllLaps(u *[]Lap) (err error) {

	result := DB.Order("discovery_time desc").Find(u)
	return result.Error
}

// Return all laps in system order by date
func GetLastLap(u *Lap) (err error) {

        result := DB.Order("discovery_time desc").First(u)
	        return result.Error
		}

// Return last known lap
func GetLastRaceIDandTime(u *Lap) (lastLapRaceID uint, lastLapTime time.Time) {
	if DB.Order("discovery_time desc").First(u).Error == nil {
		lastLapRaceID = u.RaceID
		lastLapTime = u.DiscoveryTimePrepared
	}
	return
}

func GetLastLapNumberFromRaceByTagID(u *Lap, tagID string, raceID uint) (lastLapNumber uint) {
	if DB.Where("tag_id = ? AND race_id = ?", tagID, raceID).Order("discovery_time desc").First(u).Error == nil {
		lastLapNumber = u.LapNumber
	} else {
 		lastLapNumber = 0
	}
	return
}

func GetMyLastLapDataFromCurrentRace(u *Lap)  (err error) {
	result := DB.Where("tag_id = ? AND race_id = ?", u.TagID, u.RaceID).Order("discovery_time desc").First(u)
	return result.Error
}


// Get laps by tag ID
func GetAllLapsByTagId(u *[]Lap, tag_id string) (err error) {
	result := DB.Where("tag_id = ?" , tag_id).Order("discovery_time desc").Find(u)
	return result.Error
}

func AddNewLap(u *Lap) (err error) {
	if err = DB.Create(u).Error; err != nil {
		return err
	}

	return nil
}

func GetOneLap(u *Lap, lap_id string) (err error) {
	if err := DB.Where("id = ?", lap_id).First(u).Error; err != nil {
		return err
	}

	return nil
}

func PutOneLap(u *Lap) (err error) {
	DB.Save(u)
	return nil
}

func DeleteOneLap(u *Lap, lap_id string) (err error) {
	DB.Where("id = ?", lap_id).Delete(u)
	return nil
}
