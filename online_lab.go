package onlinelab

// Config is the configuration for each online test
// The sum of all treatments' the volume proportions should be 100
type Config struct {
	Name       string
	Treatments []Treatment
}

// Treatment is used to define the volume proportions of the test
type Treatment struct {
	Name string
	// volumeProportion is the proportion of incoming volume to the treatment
	VolumeProportion int
}

// ConfigStorage is to get the persisted OnlineLabConfig
type ConfigStorage interface {
	GetConfig(name string) (Config, error)
}
