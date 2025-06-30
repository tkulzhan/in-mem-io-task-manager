package entity

import (
	"context"
	"fmt"
	"in-mem-io-task-manager/internal/infrastructure/logger"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusPending  TaskStatus = "pending"
	StatusRunning  TaskStatus = "running"
	StatusCanceled TaskStatus = "canceled"
	StatusFinished TaskStatus = "finished"
)

type DefaultTask struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description,omitempty"`
	Status         TaskStatus     `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	StartedAt      *time.Time     `json:"started_at"`
	FinishedAt     *time.Time     `json:"finished_at"`
	ProcessingTime *time.Duration `json:"processing_time"`
}

type DefaultTaskInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

func (t *DefaultTask) Execute(ctx context.Context, l *logger.Logger) error {
	startTime := time.Now()
	t.Status = StatusRunning
	t.StartedAt = &startTime

	randomDuration := time.Duration(rand.Intn(3)+3) * time.Minute
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(randomDuration)

	for {
		select {
		case <-ctx.Done():
			t.Status = StatusCanceled
			return ctx.Err()

		case <-timeout:
			finishTime := time.Now()
			t.Status = StatusFinished
			t.FinishedAt = &finishTime
			processingTime := finishTime.Sub(startTime)
			t.ProcessingTime = &processingTime
			return nil

		case <-ticker.C:
			l.Trace("Task is running", map[string]interface{}{
				"task_id":      t.GetID(),
				"title":        t.Title,
				"description":  t.Description,
				"elapsed_time": time.Since(startTime),
			})
		}
	}
}

func (t *DefaultTask) GetID() string {
	return t.ID
}

func NewDefaultTask(title, description string) *DefaultTask {
	return &DefaultTask{
		ID:          generateTaskID(),
		Title:       title,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}
}

func generateTaskID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return time.Now().Format("20060102150405") + "-" + fmt.Sprint(rand.Intn(100000))
	}
	return id.String()
}
