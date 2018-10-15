package onlinelab

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func createConfig(treatments []Treatment) Config {
	return Config{"Test", treatments}
}

type MockConfigStorage struct {
	config Config
}

func (cs *MockConfigStorage) GetConfig(labName string) (Config, error) {
	return cs.config, nil
}

func (cs *MockConfigStorage) SetTreatmentConfig(config Config) {
	cs.config = config
}

func testVolumeDividing(t *testing.T, tc *TreatmentController,
	expectedRate float64) {

	cntT1 := 0
	cntT2 := 0

	var (
		tn   string
		err1 error
	)
	for i := 1; i <= 1000000; i++ {
		tn, err1 = tc.GetNextTreatment(i)
		if err1 != nil {
			t.Error(err1)
			break
		}
		switch tn {
		case "T1":
			cntT1++
		case "T2":
			cntT2++
		}
	}
	fmt.Printf("cntT1 = %d, cntT2 = %d\n", cntT1, cntT2)
	actualRate := float64(cntT1) / (float64(cntT1) + float64(cntT2))
	if actualRate != expectedRate {
		t.Errorf("expected rate is %f, but actual is %f", expectedRate, actualRate)
	}
}

func TestVolumeDividingWithoutRefresh(t *testing.T) {
	mc := &MockConfigStorage{}
	treatments := []Treatment{Treatment{"T1", 20}, Treatment{"T2", 80}}
	mc.SetTreatmentConfig(Config{"test", treatments})
	ctx := context.Background()

	tc, err := CreateTreatmentController(ctx, mc, "test", NeverRefresh)
	if err != nil {
		t.Error(err)
	}
	testVolumeDividing(t, tc, 0.2)
}

func TestRefreshConfig(t *testing.T) {
	mc := &MockConfigStorage{}
	treatments := []Treatment{Treatment{"T1", 20}, Treatment{"T2", 80}}
	mc.SetTreatmentConfig(Config{"test", treatments})
	ctx := context.Background()
	tc, err := CreateTreatmentController(ctx, mc, "test", time.Second*1)
	if err != nil {
		t.Error(err)
	}
	testVolumeDividing(t, tc, 0.2)
	treatments = []Treatment{Treatment{"T1", 40}, Treatment{"T2", 60}}
	mc.SetTreatmentConfig(Config{"test", treatments})
	testVolumeDividing(t, tc, 0.2)
	time.Sleep(time.Second * 2)
	testVolumeDividing(t, tc, 0.4)

}
