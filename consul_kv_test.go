package onlinelab

import (
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestGetConfig(t *testing.T) {
	ccs, _ := NewConsulConfigStorage(&api.Config{})
	// local has not consul service
	if ccs == nil {
		ccs = &ConsulConfigStorage{}
	}

	config, _ := ccs.GetConfig("testLabNameNotExits")
	if config.Name != "" && len(config.treatments) != 0 {
		t.Error("config is not original value")
	}

	treatments := []Treatment{Treatment{"T1", 40}, Treatment{"T2", 60}}
	ccs.SetConfig(Config{"testLabName", treatments})
	config, _ = ccs.GetConfig("testLabName")
	if config.Name != "testLabName" {
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
