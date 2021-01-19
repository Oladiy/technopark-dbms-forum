package delivery

import (
	"encoding/json"
	"net/http"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/common/utils"
	databaseService "technopark-dbms-forum/internal/pkg/database_service"
)

type DatabaseServiceDelivery struct {
	DatabaseServiceRepository databaseService.Repository
}

func NewDatabaseServiceDelivery(databaseServiceRepository databaseService.Repository) *DatabaseServiceDelivery {
	return &DatabaseServiceDelivery {
		DatabaseServiceRepository: databaseServiceRepository,
	}
}

func (t *DatabaseServiceDelivery) GetStatus(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp, err := t.DatabaseServiceRepository.GetStatus()

	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)

		return
	}

	outputBuf, _ := json.Marshal(resp)
	_, _ = w.Write(outputBuf)
}

func (t *DatabaseServiceDelivery) ClearDatabase(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err := t.DatabaseServiceRepository.ClearDatabase()
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)

		return
	}
}