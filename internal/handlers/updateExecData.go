package handlers

import (
	"PLMS/internal/services"
	"gorm.io/gorm"
)

type UpdateExecDataHandler struct {
	db      *gorm.DB
	service *services.UpdateExcelDataService
}

func NewUpdateExecDataHandler(db *gorm.DB) *UpdateExecDataHandler {
	return &UpdateExecDataHandler{
		db:      db,
		service: services.NewUpdateExcelDataService(db),
	}
}

func (h *UpdateExecDataHandler) UpdateExecData(filePath string) error {
	return h.service.ImportExcelData(filePath)
}
