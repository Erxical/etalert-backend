package service

type BedtimeInput struct {
	GoogleId  string `bson:"googleId"`
	SleepTime string `bson:"sleepTime"`
	WakeTime  string `bson:"wakeTime"`
}

type BedtimeUpdater struct {
	SleepTime string `bson:"sleepTime"`
	WakeTime  string `bson:"wakeTime"`
}

type BedtimeResponse struct {
	SleepTime string `bson:"sleepTime"`
	WakeTime  string `bson:"wakeTime"`
}

type BedtimeService interface {
	InsertBedtime(bedtime *BedtimeInput) error
	GetBedtimeInfo(string) (*BedtimeResponse, error)
	UpdateBedtime(str string, bedtime *BedtimeUpdater) error
}
