package onlinelab

import (
	"context"
	"errors"
	"testing"
	"time"
)

const FailSecondTime = "FailSecondTime"

func createConfig(treatments []Treatment) Config {
	return Config{treatments}
}

type MockConfigStorage struct {
	cntFailSecondTime int
	config            Config
}

func (cs *MockConfigStorage) GetConfig(labName string) (Config, error) {
	if labName == FailSecondTime {
		cs.cntFailSecondTime++
		if cs.cntFailSecondTime == 2 {
			return cs.config, errors.New("failed to get config when refreshing")
		}
	}
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
	t.Logf("cntT1 = %d, cntT2 = %d\n", cntT1, cntT2)
	actualRate := float64(cntT1) / (float64(cntT1) + float64(cntT2))
	if actualRate != expectedRate {
		t.Errorf("expected rate is %f, but actual is %f", expectedRate, actualRate)
	}
}

func TestCreateControllerWithInvalidTreatmentSetting(t *testing.T) {
	mc := &MockConfigStorage{}
	treatments := []Treatment{Treatment{"T1", 20}, Treatment{"T2", 80},
		Treatment{"T3", 10}}
	mc.SetTreatmentConfig(Config{treatments})
	ctx := context.Background()
	_, err := CreateTreatmentController(ctx, mc, "test", NeverRefresh, nil)
	if err == nil {
		t.Error("invalid treament error is expected")
	}
}

func TestCreateControllerWithFailToRefresh(t *testing.T) {
	isFailedWhenRefreshing := false
	refreshErrHandler := func(err error) {
		isFailedWhenRefreshing = true
	}
	mc := &MockConfigStorage{}
	treatments := []Treatment{Treatment{"T1", 20}, Treatment{"T2", 80}}
	mc.SetTreatmentConfig(Config{treatments})
	ctx := context.Background()
	_, err := CreateTreatmentController(ctx, mc, FailSecondTime,
		time.Millisecond*100, refreshErrHandler)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 1)
	if !isFailedWhenRefreshing {
		t.Error("the error about refreshing is expected")
	}
}

func TestVolumeDividingWithoutRefresh(t *testing.T) {
	mc := &MockConfigStorage{}
	treatments := []Treatment{Treatment{"T1", 20}, Treatment{"T2", 80}}
	mc.SetTreatmentConfig(Config{treatments})
	ctx := context.Background()
	tc, err := CreateTreatmentController(ctx, mc, "test", NeverRefresh, nil)
	if err != nil {
		t.Error(err)
	}
	testVolumeDividing(t, tc, 0.2)
}

func TestRefreshConfig(t *testing.T) {
	mc := &MockConfigStorage{}
	treatments := []Treatment{Treatment{"T1", 20}, Treatment{"T2", 80}}
	mc.SetTreatmentConfig(Config{treatments})
	ctx := context.Background()
	tc, err := CreateTreatmentController(ctx, mc, "test", time.Second*1, nil)
	if err != nil {
		t.Error(err)
	}
	testVolumeDividing(t, tc, 0.2)
	treatments = []Treatment{Treatment{"T1", 40}, Treatment{"T2", 60}}
	mc.SetTreatmentConfig(Config{treatments})
	testVolumeDividing(t, tc, 0.2)
	time.Sleep(time.Second * 2)
	testVolumeDividing(t, tc, 0.4)
}
