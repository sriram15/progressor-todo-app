package service

import (
	"context"
	"math"

	"github.com/sriram15/progressor-todo-app/internal/connection"
	"github.com/sriram15/progressor-todo-app/internal/database"
)

type GetStatsResult struct {
	WeekHrs  float64 `json:"weekHrs"`
	MonthHrs float64 `json:"monthHrs"`
	YearHrs  float64 `json:"yearHrs"`
}

type IProgressService interface {
	GetStats() (GetStatsResult, error)
	GetDailyTotalMinutes() ([]database.GetDailyTotalMinutesRow, error)
}

type ProgressService struct {
	ctx                   context.Context
	taskCompletionService ITaskCompletionService
	queries               *database.Queries
}

func NewProgressService(taskCompletionService ITaskCompletionService, queries *database.Queries) *ProgressService {
	return &ProgressService{
		ctx:                   context.Background(),
		taskCompletionService: taskCompletionService,
		queries:               queries,
	}
}

func (p *ProgressService) GetStats() (GetStatsResult, error) {
	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return GetStatsResult{}, err
	}

	weekMins, err := p.queries.AggregateWeekHours(p.ctx, db, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}
	monthMins, err := p.queries.AggregateMonthHours(p.ctx, db, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}

	yearMins, err := p.queries.AggregateYearHours(p.ctx, db, int64(1))
	if err != nil {
		return GetStatsResult{}, err
	}

	// Convert to hours from mins
	weekHours := math.Ceil(weekMins / 60.0)
	monthHours := math.Ceil(monthMins / 60.0)
	yearHours := math.Ceil(yearMins / 60.0)

	return GetStatsResult{
		WeekHrs:  weekHours,
		MonthHrs: monthHours,
		YearHrs:  yearHours,
	}, nil
}

func (p *ProgressService) GetDailyTotalMinutes() ([]database.GetDailyTotalMinutesRow, error) {
	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return nil, err
	}
	return p.queries.GetDailyTotalMinutes(p.ctx, db)
}

func (p *ProgressService) GetTotalExpForUser() (float64, error) {
	return p.taskCompletionService.TotalUserExp(userId)
}
