package handlers

import (
	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage"
)

type EventHandler struct {
	log    logger.Logger
	strgPG storage.StoragePgI
}

func NewHandler(log logger.Logger, strgPG storage.StoragePgI) *EventHandler {
	return &EventHandler{
		log:    log,
		strgPG: strgPG,
	}
}
