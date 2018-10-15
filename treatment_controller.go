package onlinelab

// treatmentRange : the related treatment would be return
// when current slot number is in the range [MinSlotNum, MaxSlotNum]
type treatmentRange struct {
	MinSlotNum    int
	MaxSlotNum    int
	TreatmentName string
}

// TreatmentController is to assign the requests to the right treatment
// by the settings
// type TreatmentController interface {
// 	GetNextTreatment(requestID int) Treatment
// 	Delete(context context.Context)
// }

type TreatmentController struct {
	config          Config
	treatmentRanges []treatmentRange
}

// func CreateTreatmentController(config Config,
// 	refreshInterval time.Duration) (TreatmentController, error) {
// 	return nil
// }

func (tc *TreatmentController) calculateSlotRange() {
	for treatment := range tc.config.treatments {

	}
}

// func (tc *TreatmentController) GetNextTreatment(requestID int) string {
// 	return nil
// }
