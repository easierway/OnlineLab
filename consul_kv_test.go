package onlinelab

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	ccs, _ := NewConsulConfigStorage("localhost:8500")
	// local has not consul service
	if ccs == nil {
		ccs = &ConsulConfigStorage{}
	}
	treatments := []Treatment{Treatment{"T1", 40}, Treatment{"T2", 60}}
	ccs.SetConfig(Config{"test", treatments})
	config, _ := ccs.GetConfig("test")
	if config.Name != "test" {
		t.Error("config name invalid")
	}
	if len(config.treatments) != 2 {
		t.Error("config treatments invalid")
	}
	for _, tm := range config.treatments {
		if tm.Name != "T1" && tm.Name != "T2" {
			t.Error("config treatment name invalid")
		}
		if tm.VolumeProportion != 40 && tm.VolumeProportion != 60 {
			t.Error("config treatment volumeProportion invalid")
		}
	}
}
