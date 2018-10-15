package onlinelab

import "context"

// Config is the configuration for each online test
type Config struct {
	Name       string
	treatments []Treatment
}

// Treatment is used to define the volume proportions of the test
type Treatment struct {
	Name             string
	volumeProportion float32
}

// ConfigStorage is to get the persisted OnlineLabConfig
type ConfigStorage interface {
	GetConfig() (Config, error)
}

// TreatmentController is to assign the requests to the right treatment
// by the settings
type TreatmentController interface {
	GetNextTreatment(requestID int) string
	Delete(context context.Context)
}
