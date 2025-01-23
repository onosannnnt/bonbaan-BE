package serviceAdapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	ServiceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

type ServiceHandler struct {
	ServiceUsecase ServiceUsecase.ServiceUsecase
}

func NewServiceHandler(ServiceUsecase ServiceUsecase.ServiceUsecase) *ServiceHandler {

	return &ServiceHandler{ServiceUsecase: ServiceUsecase}

}

func (h *ServiceHandler) CreateService(c *fiber.Ctx) error {
		var service Entities.Service

		if err := c.BodyParser(&service); err != nil {
			return utils.ResponseJSON(c, fiber.StatusBadRequest, "Please fill all the require fields", err, nil)
		}

		if err := h.ServiceUsecase.CreateService(&service); err != nil {
			return utils.ResponseJSON(c, fiber.StatusConflict, "this role already exists", err, nil)
		}

		return utils.ResponseJSON(c, fiber.StatusCreated, "success", nil, nil)
	
}


func (h *ServiceHandler) GetAll(c *fiber.Ctx) error {
    services, err := h.ServiceUsecase.GetAll()
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
    }
    return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, services)
}

func (h *ServiceHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	service, err := h.ServiceUsecase.GetByID(&id)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Success", nil, service)
}


func (h *ServiceHandler) UpdateService(c *fiber.Ctx) error {
    id := c.Params("id")
    var service Entities.Service

    if err := c.BodyParser(&service); err != nil {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err, nil)
    }

    // Convert the id to uuid.UUID
    uuidID, err := uuid.Parse(id)
    if err != nil {
        return utils.ResponseJSON(c, fiber.StatusBadRequest, "Invalid UUID format", err, nil)
    }

    service.ID = uuidID // Ensure the ID is set to the one from the URL

    if err := h.ServiceUsecase.UpdateService(&service); err != nil {
        return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Failed to update service", err, nil)
    }

    return utils.ResponseJSON(c, fiber.StatusOK, "Service updated successfully", nil, nil)
}



// func (h *ServiceHandler) GetServiceByID(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	service, err := h.usecase.GetServiceByID(uint(id))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(service)
// }

// func (h *ServiceHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	if err := h.usecase.DeleteService(uint(id)); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusNoContent)
// }

// func (h *ServiceHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
// 	var service Entities.Service
// 	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	if err := h.usecase.UpdateService(&service); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// }
