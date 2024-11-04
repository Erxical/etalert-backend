package repository

import "time"

type Counter struct {
	ID  string `bson:"_id,omitempty"`
	Seq int    `bson:"seq"`
}

type Schedule struct {
	Id              string    `bson:"_id,omitempty"`
	GoogleId        string    `bson:"googleId"`
	RoutineId       string    `bson:"routineId"`
	Name            string    `bson:"name"`
	Date            time.Time `bson:"date"`
	StartTime       string    `bson:"startTime"`
	EndTime         string    `bson:"endTime"`
	IsHaveEndTime   bool      `bson:"isHaveEndTime"`
	OriName         string    `bson:"oriName"`
	OriLatitude     float64   `bson:"oriLatitude"`
	OriLongitude    float64   `bson:"oriLongitude"`
	DestName        string    `bson:"destName"`
	DestLatitude    float64   `bson:"destLatitude"`
	DestLongitude   float64   `bson:"destLongitude"`
	GroupId         int       `bson:"groupId"`
	Transportation  string    `bson:"transportation"`
	Priority        int       `bson:"priority"`
	IsHaveLocation  bool      `bson:"isHaveLocation"`
	IsFirstSchedule bool      `bson:"isFirstSchedule"`
	IsTraveling     bool      `bson:"isTraveling"`
	IsUpdated       bool      `bson:"isUpdated"`

	Recurrence   string `bson:"recurrence"`
	RecurrenceId int    `bson:"recurrenceId"`
}

type TrafficResponse struct {
	Tm struct {
		ID  string `json:"@id"`
		Poi []struct {
			ID string `json:"id"`
			P  struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"p"`
			Ic  int    `json:"ic"`
			Ty  int    `json:"ty"`
			Cs  int    `json:"cs"`
			D   string `json:"d"`
			C   string `json:"c"`
			F   string `json:"f"`
			T   string `json:"t"`
			L   int    `json:"l"`
			Dl  int    `json:"dl"`
			R   string `json:"r"`
			Cbl struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"cbl"`
			Ctr struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"ctr"`
		} `json:"poi"`
	} `json:"tm"`
}

type Forecast struct {
	Summary   Summary    `json:"summary"`
	Waypoints []Waypoint `json:"waypoints"`
}

type Summary struct {
	IconCode int    `json:"iconCode"`
	Hazards  Hazard `json:"hazards"`
}

type Hazard struct {
	MaxHazardIndex int `json:"maxHazardIndex"`
}

type Waypoint struct {
	IconCode       int           `json:"iconCode"`
	ShortPhrase    string        `json:"shortPhrase"`
	IsDayTime      bool          `json:"isDayTime"`
	CloudCover     int           `json:"cloudCover"`
	Precipitation  Precipitation `json:"precipitation"`
	LightningCount int           `json:"lightningCount"`
	Hazards        Hazard        `json:"hazards"`
	Notifications  []string      `json:"notifications"`
}

type Precipitation struct {
	Dbz  float64 `json:"dbz"`
	Type string  `json:"type"`
}

type ScheduleRepository interface {
	GetTravelTime(oriLat string, oriLong string, destLat string, destLong string, mode string, depTime string) (string, error)
	GetTraffic(oriLat float64, oriLong float64, destLat float64, destLong float64) (TrafficResponse, error)
	GetWeather(oriLat string, oriLong string, destLat string, destLong string, depTime string) (Forecast, error)
	GetNextGroupId() (int, error)
	GetNextRecurrenceId() (int, error)
	CalculateNextRecurrenceDate(currentDate, recurrence string, count int) ([]string, error)
	BatchInsertSchedules(schedules []Schedule) error
	InsertSchedule(schedule *Schedule) error
	GetAllSchedules(gId string, date string) ([]*Schedule, error)
	GetScheduleById(id string) (*Schedule, error)
	GetSchedulesByGroupId(groupId int) ([]*Schedule, error)
	GetMainSchedulesByRecurrenceId(recurrenceId int, date string) ([]*Schedule, error)
	GetSchedulesByRecurrenceId(recurrenceId int, date string) ([]*Schedule, error)
	UpdateSchedule(id string, schedule *Schedule) error
	UpdateScheduleTime(id string, startTime string, endTime string) error
	DeleteSchedule(groupId int) error
	DeleteScheduleByRecurrenceId(recurrenceId int, date string) error
}
