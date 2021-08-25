package Models

/**
* This package module have some methods for storage RFID labels
* and store them into database
*/

import (
	"../Packages/Config" // our packages
	"encoding/xml"
	"encoding/csv"
	"strings"
	"strconv"
	"bytes"
	"sync"
	"time"
	"log"
	"net"
	"fmt"
	"../Packages/Proxy"
)

// Buffer for new RFID requests
var laps []Lap

// Laps locker
var lapsLocker sync.Mutex

// Laps save into DB interval
var lapsSaveInterval int

// Check RFID mute timeout map
var rfidTimeoutMap map[string]time.Time

// Mute timeout duration (stored in .env)
var rfidLapMinimalTime int

// Check RFID mute timeout locker
var rfidTimeoutLocker sync.Mutex

// Start antenna listener
func StartAntennaListener(appAntennaListenerIp, rfidLapMinimalTimeString, lapsSaveIntervalString string, TIME_ZONE string, RACE_TIMEOUT_SEC int64) {

	if Config.PROXY_ACTIVE=="true" {
		fmt.Println("Started tcp proxy restream to", Config.PROXY_HOST,"and port:",Config.PROXY_PORT )
	}

	// Start buffer synchro with database
	go startSaveLapsBufferToDatabase(RACE_TIMEOUT_SEC)

	// Create RFID mute timeout
	rfidTimeoutMap = make(map[string]time.Time)

	// Prepare rfidLapMinimalTime
	rfidTimeout, rfidTimeoutErr := strconv.Atoi(rfidLapMinimalTimeString)
	if rfidTimeoutErr != nil {
		log.Panicln("Incorrect MINIMAL_LAP_TIME parameter in .env file")
	}
	rfidLapMinimalTime = int(rfidTimeout)

	// Prepare lapsSaveInterval
	lapsInterval, lapsIntervalErr := strconv.Atoi(lapsSaveIntervalString)
	if lapsIntervalErr != nil {
		log.Panicln("Incorrect LAPS_SAVE_INTERVAL parameter in .env file")
	}
	lapsSaveInterval = int(lapsInterval)

	// Start listener
	l, err := net.Listen("tcp", appAntennaListenerIp)
	if err != nil {
		log.Panicln("Can't start the antenna listener", err)
	}
	defer l.Close()

	// Listen new connections
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go newAntennaConnection(conn, TIME_ZONE)
	}
}

// Save laps buffer to database
func startSaveLapsBufferToDatabase(RACE_TIMEOUT_SEC int64) {
	for range time.Tick(time.Duration(lapsSaveInterval) * time.Second) {
		lapsLocker.Lock()
		var lapStruct Lap
		var currentRaceID, currentLapNumber uint
		lastRaceID, lastLapTime := GetLastRaceIDandTime(&lapStruct)
		if lastRaceID == 0 {
			currentRaceID = 1
		} else {
			currentRaceID = lastRaceID
			if (time.Now().UnixNano()/int64(time.Millisecond)-(RACE_TIMEOUT_SEC*1000) > lastLapTime.UnixNano()/int64(time.Millisecond)) {
				//last lap data was created more than 300 seconds ago
				//RaceID++ (create new race)
				currentRaceID = (lastRaceID+1)

			} else {
				//last lap data was created less than 300 seconds ago
				currentRaceID = lastRaceID
			}
		}


		// Save laps to database
		for _,lap := range laps {

			lastLapNumber := GetLastLapNumberFromRaceByTagID(lap.TagID, currentRaceID)
			if lastLapNumber == 0 {
				currentLapNumber = 1
			} else {
				currentLapNumber = lastLapNumber+1
			}
			lap.LapNumber=currentLapNumber
			lap.RaceID=currentRaceID
			fmt.Printf("Saved to db: %s, %d, %d\n", lap.TagID, lap.DiscoveryTime.UnixNano()/int64(time.Millisecond), lap.Antenna)
			if err := AddNewLap(&lap); err != nil {
				fmt.Println("Error. Lap not added to database")
			}
		}


		// Clear lap buffer
		var cL []Lap
		laps = cL
		lapsLocker.Unlock()

	}
}

// Add new lap to laps buffer (private func)
func addNewLapToLapsBuffer(lap Lap) {

	// Check minimal lap time (we save only laps grater than MINIMAL_LAP_TIME from .env file)

	if expiredTime, ok := rfidTimeoutMap[lap.TagID]; !ok {

		// First time for this TagID, save lap to buffer
		lapsLocker.Lock()
		laps = append(laps, lap)
		lapsLocker.Unlock()

		// Add new value to timeouts checker map
		setNewExpriredDataForRfidTag(lap.TagID)


	} else {

		// Check previous time
		tN := time.Now()
		if tN.After(expiredTime)  {

			// Time is over, save lap to buffer
			lapsLocker.Lock()
			laps = append(laps, lap)
			lapsLocker.Unlock()


			// Generate new expired time
			setNewExpriredDataForRfidTag(lap.TagID)


		} 
	}
}

// Set new expired date for rfid Tag
func setNewExpriredDataForRfidTag(tagID string) {

	newExpiredTime := time.Now().Add(time.Duration(rfidLapMinimalTime) * time.Second)
	rfidTimeoutLocker.Lock()
	rfidTimeoutMap[tagID] = newExpiredTime
	rfidTimeoutLocker.Unlock()

}

// New antenna connection (private func)
func newAntennaConnection(conn net.Conn, TIME_ZONE string) {

	defer conn.Close()

	// Read connection in lap
	for {
		buf := make([]byte, 8192)
		size, err := conn.Read(buf)
		if err == nil {
			data := buf[:size]
			var lap Lap
			err := xml.Unmarshal(data, &lap)

			// CSV data processing
			if err != nil {

				fmt.Println("Received data is not XML, trying CSV text...", err)

				//received data of type TEXT (parse TEXT).
				r := csv.NewReader(bytes.NewReader(data))
				r.Comma = ','
				r.FieldsPerRecord = 3
				CSV, err := r.Read()
				if err != nil {
					fmt.Println("Recived incorrect CSV data", err)
					continue
				}

				// Prepare antenna position
				antennaPosition, antennaErr := strconv.Atoi(strings.TrimSpace(CSV[2]))
				if antennaErr != nil {
					fmt.Println("Recived incorrect Antenna position value:", antennaErr)
					continue
				}

				// Prepare date
				fmt.Println(Config.TIME_ZONE)
				loc, loadLocErr := time.LoadLocation(Config.TIME_ZONE)
				if loadLocErr != nil {
					fmt.Println("time.LoadLocation(Config.TIME_ZONE) error:", loadLocErr)
					continue
				}

				xmlTimeFormat := `2006/01/02 15:04:05.000`
				discoveryTime, parseTimeErr := time.ParseInLocation(xmlTimeFormat, strings.TrimSpace(CSV[1]), loc)
				if parseTimeErr != nil {
					fmt.Println("Recived incorrect time from RFID reader:", parseTimeErr)
					continue
				}

				lap.DiscoveryTime = discoveryTime
				lap.TagID = strings.TrimSpace(CSV[0])
				lap.Antenna = uint8(antennaPosition)

				// XML data processing
			} else {

				// Prepare date
				loc, _ := time.LoadLocation(Config.TIME_ZONE)
				xmlTimeFormat := `2006/01/02 15:04:05.000`
				discoveryTime, err := time.ParseInLocation(xmlTimeFormat, lap.DiscoveryUnixTime, loc)



				//unixMillyTime:=discoveryTime.UnixNano()/int64(time.Millisecond)
				// If time is incorrect than skip them
				if err != nil {
					continue
				}

				lap.DiscoveryTime = discoveryTime
			}

			// Additional preparing for TagID
			lap.TagID = strings.ReplaceAll(lap.TagID, " ", "")

			//Debug all received data from RFID reader
			fmt.Printf("%s, %d, %d\n", lap.TagID, lap.DiscoveryTime.UnixNano()/int64(time.Millisecond), lap.Antenna)


			if Config.PROXY_ACTIVE=="true" {
				go Proxy.ProxyDataToMotosponder(lap.TagID, lap.DiscoveryTime.UnixNano()/int64(time.Millisecond), lap.Antenna )
			}
			// Add current Lap to Laps buffer
			go addNewLapToLapsBuffer(lap)
		}
	}
}
