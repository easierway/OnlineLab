package onlinelab

import (
	"context"
	"errors"
	"time"
)

// NeverRefresh is to not start the config refreshing
const NeverRefresh time.Duration = -1

// ErrGetTreatment will be thrown when failing to get proper treatment
var ErrGetTreatment = errors.New("can not find the proper treatment")

// treatmentRange : the related treatment would be return
// when current slot number is in the range [MinSlotNum, MaxSlotNum]
type treatmentRange struct {
	TreatmentName string
	MinSlotNum    int
	MaxSlotNum    int
}

// TreatmentController is assign the volume to the proper treatment by the settings
type TreatmentController struct {
	config          Config
	treatmentRanges []treatmentRange
}

func (tc *TreatmentController) createTreatmentRange(name string,
	minSlot int, maxSlot int) treatmentRange {
	return treatmentRange{name, minSlot, maxSlot}
}

func (tc *TreatmentController) refreshConfig(context context.Context,
	refreshInterval time.Duration, configStorage ConfigStorage, labName string) {
	if refreshInterval <= 0 {
		return
	}
	for _ = range time.Tick(refreshInterval) {
		select {
		case <-context.Done():
		default:
			config, err := configStorage.GetConfig(labName)
			if err == nil {
				tc.config = config
				tc.calculateSlotRange()
			}
		}
	}
}

// CreateTreatmentController is to create a TreatmentController
func CreateTreatmentController(context context.Context, configStorage ConfigStorage,
	labName string, refreshInterval time.Duration) (*TreatmentController, error) {
	config, err := configStorage.GetConfig(labName)
	if err != nil {
		return nil, err
	}
	treatmentController := &TreatmentController{
		config: config,
	}
	go treatmentController.refreshConfig(context, refreshInterval, configStorage, labName)
	treatmentController.calculateSlotRange()
	return treatmentController, nil
}

func (tc *TreatmentController) calculateSlotRange() {
	offset := 0
	var treatmentRanges []treatmentRange
	for _, treatment := range tc.config.treatments {
		treatmentRanges = append(treatmentRanges, tc.createTreatmentRange(treatment.Name,
			offset, offset+treatment.VolumeProportion-1))
		offset = offset + treatment.VolumeProportion
	}
	tc.treatmentRanges = treatmentRanges
}

// GetNextTreatment is to get the treatment for the coming request
func (tc *TreatmentController) GetNextTreatment(requestID int) (string, error) {
	nextSlotNum := requestID % 100
	//	fmt.Printf("Next Slot Num: %d\n", nextSlotNum)
	for _, tr := range tc.treatmentRanges {
		//		fmt.Printf("treatmentRange: %v\n", tr)
		if nextSlotNum >= tr.MinSlotNum && nextSlotNum <= tr.MaxSlotNum {
			return tr.TreatmentName, nil
		}
	}
	return "", ErrGetTreatment
}
