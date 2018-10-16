package onlinelab

import "testing"

func TestSetConfig(t *testing.T) {
	ccs, err := NewConsulConfigStorage("localhost:8500")
	if err != nil {
		t.Error(err)
	}
	treatments := []Treatment{Treatment{"T1", 40}, Treatment{"T2", 60}}
	if err := ccs.SetConfig(Config{"test", treatments}); err != nil {
		t.Error(err)
	}
}

func TestGetConfig(t *testing.T) {
	ccs, err := NewConsulConfigStorage("localhost:8500")
	if err != nil {
		t.Error(err)
	}
	config, err := ccs.GetConfig("test")
	if err != nil {
		t.Error(err)
	}
	if config.Name != "test" {
		t.Error("config name invalid")
	}
	if len(config.Treatments) != 2 {
		t.Error("config treatments invalid")
	}
	for _, tm := range config.Treatments {
		if tm.Name != "T1" && tm.Name != "T2" {
			t.Error("config treatment name invalid")
		}
		if tm.VolumeProportion != 40 && tm.VolumeProportion != 60 {
			t.Error("config treatment volumeProportion invalid")
		}
	}
}
