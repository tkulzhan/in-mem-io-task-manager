package entity

import (
	"context"
	"math/rand"
	"time"
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
	Status         TaskStatus     `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	StartedAt      *time.Time     `json:"started_at,omitempty"`
	FinishedAt     *time.Time     `json:"finished_at,omitempty"`
	ProcessingTime *time.Duration `json:"processing_time,omitempty"`
}

func (t *DefaultTask) Execute(ctx context.Context) {
	startTime := time.Now()
	t.Status = StatusRunning
	t.StartedAt = &startTime

	// Simulate task processing for 3 to 5 minutes
	randomDuration := time.Duration(rand.Intn(3)+3) * time.Minute
	time.Sleep(randomDuration)

	finishTime := time.Now()
	t.Status = StatusFinished
	t.FinishedAt = &finishTime
	processingTime := finishTime.Sub(startTime)

	t.ProcessingTime = &processingTime
}

func (t *DefaultTask) GetID() string {
	return t.ID
}