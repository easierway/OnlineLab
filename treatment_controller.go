package onlinelab

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// NeverRefresh is to not start the config refreshing
const NeverRefresh time.Duration = -1

// ErrGetTreatment will be thrown when failing to get proper treatment
var ErrGetTreatment = errors.New("can not find the proper treatment")

// ErrTreatmentProportion will be thrown when the treatment proportion settings are wrong
var ErrTreatmentProportion = errors.New("sum of the proportion should be 100")

// treatmentRange : the related treatment would be return
// when current slot number is in the range [MinSlotNum, MaxSlotNum]
type treatmentRange struct {
	TreatmentName string
	MinSlotNum    int
	MaxSlotNum    int
}

// RefreshErrorHandlerFn is used to handle the configuration refreshing errors
type RefreshErrorHandlerFn func(err error)

// TreatmentController is assign the volume to the proper treatment by the settings
type TreatmentController struct {
	config          atomic.Value
	treatmentRanges atomic.Value
}

func (tc *TreatmentController) createTreatmentRange(name string,
	minSlot int, maxSlot int) treatmentRange {
	return treatmentRange{name, minSlot, maxSlot}
}

func (tc *TreatmentController) refreshConfig(context context.Context,
	refreshInterval time.Duration, configStorage ConfigStorage,
	labName string, refreshErrHandlerFn RefreshErrorHandlerFn) {
	if refreshInterval <= 0 {
		return
	}
	for _ = range time.Tick(refreshInterval) {
		select {
		case <-context.Done():
			return
		default:
			config, err := configStorage.GetConfig(labName)
			if err == nil {
				tc.config.Store(config)
				err = tc.calculateSlotRange()
			}
			if err != nil {
				refreshErrHandlerFn(err)
			}
		}
	}
}

// CreateTreatmentController is to create a TreatmentController
// configStorage: the implementations for different backend storages (e.g. Consul, K8S, Redis)
// labName: the name of your online lab, which relates to a treatment settings
// refreshInterval: the interval of refresh the settings from the backend storages
// refreshErrHandler: when failing to refreshing, the function will be invoked
func CreateTreatmentController(context context.Context, configStorage ConfigStorage,
	labName string, refreshInterval time.Duration,
	refreshErrHandler RefreshErrorHandlerFn) (*TreatmentController, error) {
	config, err := configStorage.GetConfig(labName)
	if err != nil {
		return nil, err
	}
	treatmentController := &TreatmentController{}
	treatmentController.config.Store(config)
	if err = treatmentController.validateConfig(); err != nil {
		return nil, err
	}
	go treatmentController.refreshConfig(context, refreshInterval, configStorage,
		labName, refreshErrHandler)
	treatmentController.calculateSlotRange()
	return treatmentController, nil
}

func (tc *TreatmentController) validateConfig() error {
	totalProportion := 0
	config := tc.config.Load().(*Config)
	for _, treatment := range config.treatments {
		totalProportion = totalProportion + treatment.VolumeProportion
	}
	if totalProportion != 100 {
		return ErrTreatmentProportion
	}
	return nil
}

func (tc *TreatmentController) calculateSlotRange() error {
	if err := tc.validateConfig(); err != nil {
		return err
	}
	offset := 0
	var treatmentRanges []treatmentRange
	config := tc.config.Load().(*Config)
	for _, treatment := range config.treatments {
		treatmentRanges = append(treatmentRanges, tc.createTreatmentRange(treatment.Name,
			offset, offset+treatment.VolumeProportion-1))
		offset = offset + treatment.VolumeProportion
	}
	tc.treatmentRanges.Store(&treatmentRanges)
	return nil
}

// GetNextTreatment is to get the treatment for the coming request
func (tc *TreatmentController) GetNextTreatment(requestID int) (string, error) {
	nextSlotNum := requestID % 100
	treatmentRanges := tc.treatmentRanges.Load().(*[]treatmentRange)
	for _, tr := range *treatmentRanges {
		if nextSlotNum >= tr.MinSlotNum && nextSlotNum <= tr.MaxSlotNum {
			return tr.TreatmentName, nil
		}
	}
	return "", ErrGetTreatment
}
