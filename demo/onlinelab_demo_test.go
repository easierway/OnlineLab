package demo

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/easierway/OnlineLab"
	"github.com/hashicorp/consul/api"
)

func createOnlinelab(labName string, t *testing.T) *onlinelab.TreatmentController {
	var err error
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
		labName, time.Second*1, refreshErrHandler)
	if err != nil {
		t.Error("failed to create treatment controller")
	}
	return treatmentControoler
}

// TestDemo: to run this demo, a local Consul agent is required
// And add the Configuration "HelloLab"  with the value "[{"Name":"T1","VolumeProportion":40},{"Name":"T2","VolumeProportion":60}]"
// to local Consul's KV storage.
func TestCodeSectionComparison(t *testing.T) {
	var (
		tn  string
		err error
	)
	treatmentControoler := createOnlinelab("HelloLab", t)
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

func TestDebugSectionControl(t *testing.T) {
	OutputDebugInfo := "T1"
	var (
		tn  string
		err error
	)
	treatmentControoler := createOnlinelab("HelloLab", t)
	for i := 1; i < 1000; i++ {
		tn, err = treatmentControoler.GetNextTreatment(i)
		if err != nil {
			t.Error(err)
		}
		if tn == OutputDebugInfo {
			t.Log("This is debug info.\n")
		}
	}
}

func TestProcessSpeedControl(t *testing.T) {
	var (
		tn  string
		err error
	)
	SlowDownFlag := "T1"
	treatmentControoler := createOnlinelab("HelloLab", t)

	cntT1 := 0
	cntT2 := 0
	for i := 1; i <= 10; i++ {
		tn, err = treatmentControoler.GetNextTreatment(i)
		if err != nil {
			t.Error(err)
			break
		}
		if tn == SlowDownFlag {
			time.Sleep(time.Millisecond * 10)
		}
	}
	t.Logf("T1 count: %d; T2 count: %d", cntT1, cntT2)
}

type ServiceClient interface {
	CallSomeService()
}

type ServiceClient1 struct {
}

type ServiceClient2 struct {
}

type ServiceClientProxy struct {
	service1            *ServiceClient1
	service2            *ServiceClient2
	treatmentControoler *onlinelab.TreatmentController
}

func (client *ServiceClient1) CallSomeService() {
	fmt.Println("Invoked T1 service.")
}

func (client *ServiceClient2) CallSomeService() {
	fmt.Println("Invoked T2 service.")
}

func (client *ServiceClientProxy) CallSomeService() {
	tn, err := client.treatmentControoler.GetNextTreatment(rand.Intn(100) + 1)
	if err != nil {
		panic("Failed to get next treatment.")
	}
	switch tn {
	case "T1":
		client.service1.CallSomeService()
	case "T2":
		client.service2.CallSomeService()
	}
}

func CreateServiceClientProxy(t *testing.T) *ServiceClientProxy {
	tc := createOnlinelab("HelloLab", t)
	return &ServiceClientProxy{
		&ServiceClient1{},
		&ServiceClient2{},
		tc,
	}
}

func TestServiceComparison(t *testing.T) {
	client := CreateServiceClientProxy(t)
	for i := 0; i < 10; i++ {
		client.CallSomeService()
	}
}
