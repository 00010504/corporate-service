package handlers

import (
	"context"
	"encoding/json"
	"genproto/common"

	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (e *EventHandler) Upsert(ctx context.Context, event *kafka.Message) error {

	var req common.UserCreatedModel

	if err := json.Unmarshal(event.Value, &req); err != nil {
		return err
	}

	e.log.Info("create user event", logger.Any("event", req))

	if err := e.strgPG.User().Upsert(&req); err != nil {
		return err
	}

	return nil

}
