package repository

type Bedtime struct {
	GoogleId  string `bson:"googleId"`
	SleepTime string `bson:"sleepTime"`
	WakeTime  string `bson:"wakeTime"`
}

type BedtimeUpdater struct {
	SleepTime string `bson:"sleepTime"`
	WakeTime  string `bson:"wakeTime"`
}

type BedtimeRepository interface {
	InsertBedtime(bedtime *Bedtime) error
	GetBedtimeInfo(string) (*Bedtime, error)
	UpdateBedtime(string, *BedtimeUpdater) error
}
