package StatusAdapter

import statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"

type StatusHandler struct {
	statusUsecase statusUsecase.StatusUsecase
}

func NewStatusHandler(statusUsecase statusUsecase.StatusUsecase) *StatusHandler {
	return &StatusHandler{
		statusUsecase: statusUsecase,
	}
}
