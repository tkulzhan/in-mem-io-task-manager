package service

import (
	"encoding/json"
	"in-mem-io-task-manager/internal/domain/entity"
	"in-mem-io-task-manager/internal/infrastructure/errors"
)

func parseDefaultTaskInput(data any) (entity.DefaultTaskInput, error) {
	var input entity.DefaultTaskInput

	raw, err := json.Marshal(data)
	if err != nil {
		return input, errors.NewBadRequestError("invalid data format")
	}

	if err := json.Unmarshal(raw, &input); err != nil {
		return input, errors.NewBadRequestError("invalid fields in task data")
	}

	if input.Title == "" {
		return input, errors.NewBadRequestError("title is required")
	}

	return input, nil
}

func deref(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
