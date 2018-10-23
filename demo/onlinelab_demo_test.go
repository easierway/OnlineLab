package demo

import (
	"context"
	"testing"
	"time"

	"github.com/easierway/OnlineLab"
	"github.com/hashicorp/consul/api"
)

// TestDemo: to run this demo, a local Consul agent is required
func TestDemo(t *testing.T) {
	var (
		tn  string
		err error
	)
	refreshErrHandler := func(err error) {
		if err != nil {
			t.Error(err)
		}
	}
	ctx := context.Background()
	configStorage, err := onlinelab.NewConsulConfigStorage(&api.Config{})
	if err != nil {
		t.Error("failed to create config storage")
	}
	treatmentControoler, err := onlinelab.CreateTreatmentController(ctx, configStorage,
		"HelloLab", time.Second*1, refreshErrHandler)
	if err != nil {
		t.Error("failed to create treatment controller")
	}

	cntT1 := 0
	cntT2 := 0
	for i := 1; i <= 1000000; i++ {
		tn, err = treatmentControoler.GetNextTreatment(i)
		if err != nil {
			t.Error(err)
			break
		}
		switch tn {
		case "T1":
			cntT1++
		case "T2":
			cntT2++
		}
	}
	t.Logf("T1 count: %d; T2 count: %d", cntT1, cntT2)
}
