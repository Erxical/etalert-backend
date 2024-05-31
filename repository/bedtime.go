package repository

type Bedtime struct {
	GoogleId  string `bson:"googleId"`
	SleepTime string `bson:"sleepTime"`
	WakeTime  string `bson:"wakeTime"`
}

type BedtimeRepository interface {
	InsertBedtime(bedtime *Bedtime) error
	GetBedtimeInfo(string) (*Bedtime, error)
	UpdateBedtime(str string, bedtime *Bedtime) error
}
