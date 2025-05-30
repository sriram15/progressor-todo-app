package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sriram15/progressor-todo-app/internal/connection"
	"github.com/sriram15/progressor-todo-app/internal/database"
)

var (
	ErrorUnknown           = errors.New("unknown error")
	ErrNotFound            = errors.New("not found")
	ErrInvalidProject      = errors.New("invalid project")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidUpdate       = errors.New("invalid update")
	ErrCardTitleRequired   = errors.New("card title is required")
	ErrCardTrackingStarted = errors.New("card tracking already in progress")
	ErrCardTrackingStopped = errors.New("card tracking already stopped")
)

type CardStatus int

const (
	Todo CardStatus = iota
	Done
	Active
)

type UpdateCardParams struct {
	Title         string `json:"title"`
	EstimatedMins int    `json:"estimatedMins"`
	Description   string `json:"description"`
}

const userId = 1

type ICardService interface {
	GetAll(projectId uint, status CardStatus) ([]database.ListCardsRow, error)
	GetCardById(projectId uint, id uint) (database.GetCardRow, error)
	GetActiveTimeEntry(projectId uint, id uint) (database.TimeEntry, error)
	DeleteCard(projectId uint, id uint) error
	UpdateCard(projectId uint, id uint, updateCardParam UpdateCardParams) error
	UpdateCardStatus(projectId uint, id uint, status CardStatus) error
	AddCard(projectId uint, cardTitle string, estimatedMins uint) error
	StartCard(projectId uint, id uint) error
	StopCard(projectId uint, id uint) error
	Cleanup() error
}

type CardService struct {
	ctx                   context.Context
	projectService        IProjectService
	taskCompletionService ITaskCompletionService
	queries               *database.Queries
}

func NewCardService(projectService IProjectService, taskCompletionService ITaskCompletionService, queries *database.Queries) *CardService {
	return &CardService{
		ctx:                   context.Background(),
		projectService:        projectService,
		taskCompletionService: taskCompletionService,
		queries:               queries,
	}
}

func (c *CardService) GetAll(projectId uint, status CardStatus) ([]database.ListCardsRow, error) {

	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return []database.ListCardsRow{}, err
	}

	// TODO: Add validations
	// if status != Todo || status != Done {
	// 	return []database.ListCardsRow{}, ErrInvalidStatus
	// }

	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return nil, err
	}
	cards, err := c.queries.ListCards(c.ctx, db, database.ListCardsParams{Projectid: int64(projectId), Status: int64(status)})
	return cards, err
}

func (c *CardService) GetCardById(projectId uint, id uint) (database.GetCardRow, error) {
	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return database.GetCardRow{}, err
	}

	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return database.GetCardRow{}, err
	}
	card, err := c.queries.GetCard(c.ctx, db, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return database.GetCardRow{}, ErrNotFound
	}
	return card, nil
}

func (c *CardService) GetActiveTimeEntry(projectId uint, id uint) (database.TimeEntry, error) {

	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return database.TimeEntry{}, err
	}

	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return database.TimeEntry{}, err
	}
	timeEntry, err := c.queries.GetActiveTimeEntry(c.ctx, db, int64(id))
	if err != nil {
		return database.TimeEntry{}, err
	}

	return timeEntry, nil

}
func (c *CardService) DeleteCard(projectId uint, id uint) error {

	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return err
	}
	return c.queries.DeleteCard(c.ctx, db, database.DeleteCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
}

func (c *CardService) UpdateCard(projectId uint, id uint, updateCardParam UpdateCardParams) error {

	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return err
	}
	card, err := c.queries.GetCard(c.ctx, db, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return err
	}

	// TODO: Add update card validation
	if updateCardParam.Title == "" {
		return ErrInvalidUpdate
	}

	var description sql.NullString
	if updateCardParam.Description == "" {
		description = sql.NullString{Valid: true, String: ""}
	} else {
		description = sql.NullString{Valid: true, String: updateCardParam.Description}
	}

	return c.queries.UpdateCard(c.ctx, db, database.UpdateCardParams{
		Title:         updateCardParam.Title,
		Description:   description,
		ID:            card.CardID,
		Status:        card.Status,
		Trackedmins:   card.Trackedmins,
		Estimatedmins: int64(updateCardParam.EstimatedMins),
		Completedat:   card.Completedat,
	})

}
func (c *CardService) UpdateCardStatus(projectId uint, id uint, status CardStatus) error {
	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	card, err := c.queries.GetCard(c.ctx, tx, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return err
	}

	if status != Todo && status != Done {
		return ErrInvalidStatus
	}

	completedAt := sql.NullTime{}
	if status == Done {
		completedAt = sql.NullTime{Valid: true, Time: time.Now().UTC()}
	}

	err = c.queries.UpdateCard(c.ctx, tx, database.UpdateCardParams{
		Title:         card.Title,
		Description:   card.Description,
		ID:            card.CardID,
		Status:        int64(status),
		Trackedmins:   card.Trackedmins,
		Estimatedmins: card.Estimatedmins,
		Completedat:   completedAt,
	})
	if err != nil {
		return err
	}

	// If status is done, then also insert in TaskCompletions Table
	if status == Done {
		// Calculate base exp, time bonus exp and streak bonus exp (replace with your logic).
		baseExp := int64(10)
		timeBonusExp := int64(card.Trackedmins / 5) // TODO: Get from constant
		streakBonusExp := int64(0)                  // TODO

		_, err := c.queries.GetTaskCompletion(c.ctx, tx, database.GetTaskCompletionParams{
			Cardid: card.CardID,
			Userid: userId, // TODO: Get from session
		})

		// only create record if it does not exist
		if errors.Is(err, sql.ErrNoRows) {
			_, err = c.queries.CreateTaskCompletion(c.ctx, tx, database.CreateTaskCompletionParams{
				Cardid:         card.CardID,
				Userid:         userId,
				Baseexp:        baseExp,
				Timebonusexp:   timeBonusExp,
				Streakbonusexp: streakBonusExp,
				Totalexp:       baseExp + timeBonusExp + streakBonusExp,
			})
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (c *CardService) AddCard(projectId uint, cardTitle string, estimatedMins uint) error {

	if cardTitle == "" {
		return ErrCardTitleRequired
	}

	card := database.CreateCardParams{
		Title:         cardTitle,
		Status:        int64(Todo),
		Projectid:     int64(projectId),
		Estimatedmins: int64(estimatedMins),
	}
	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return err
	}

	return c.queries.CreateCard(c.ctx, db, card)
}

func (c *CardService) StartCard(projectId uint, id uint) error {
	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return err
	}
	card, err := c.queries.GetCard(c.ctx, db, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return err
	}

	// Check if card is already active
	if card.Isactive {
		return ErrCardTrackingStarted
	}

	// Check for other open cards which is currently in progress and stop the timer there
	activeCard, err := c.queries.GetActiveCard(c.ctx, db)

	// When the active card is empty. It will throw sql.ErrNoRows. If the err is not that, then return err
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else { // We have an active card already. Try to stop that first and return err if that fails
		err := c.StopCard(projectId, uint(activeCard.ID))
		if err != nil {
			return err
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = c.queries.UpdateCardActive(c.ctx, tx, database.UpdateCardActiveParams{
		ID:          int64(id),
		Isactive:    true,
		Trackedmins: card.Trackedmins,
	})
	if err != nil {
		return err
	}

	// Create a new Timeentry object and add it to the card
	currentStartTime := time.Now().UTC()
	_, err = c.queries.CreateTimeEntry(c.ctx, tx, database.CreateTimeEntryParams{
		Cardid:    int64(id),
		Starttime: currentStartTime,
		Endtime:   currentStartTime,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *CardService) StopCard(projectId uint, id uint) error {
	_, err := c.projectService.IsValidProject(projectId)
	if err != nil {
		return err
	}

	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return err
	}
	card, err := c.queries.GetCard(c.ctx, db, database.GetCardParams{
		ID:        int64(id),
		Projectid: int64(projectId),
	})
	if err != nil {
		return err
	}

	// Check if card is still active
	if !card.Isactive {
		return ErrCardTrackingStopped
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the active time entry
	activeTimeentry, err := c.queries.GetActiveTimeEntry(c.ctx, tx, int64(id))
	if err != nil {
		return err
	}

	currentEndTime := time.Now().UTC()
	duration := currentEndTime.Sub(activeTimeentry.Starttime).Minutes()
	err = c.queries.UpdateActiveTimeEntry(c.ctx, tx, database.UpdateActiveTimeEntryParams{
		ID:       activeTimeentry.ID,
		Endtime:  currentEndTime,
		Duration: int64(duration),
	})
	if err != nil {
		return err
	}

	err = c.queries.UpdateCardActive(c.ctx, tx, database.UpdateCardActiveParams{
		ID:          int64(id),
		Isactive:    false,
		Trackedmins: card.Trackedmins + int64(duration),
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *CardService) Cleanup() error {

	// Check for other open cards which is currently in progress and stop the timer there
	db, err := connection.GetOrReconnectDB()
	if err != nil {
		return err
	}

	activeCard, err := c.queries.GetActiveCard(c.ctx, db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}

	// Stop the active Card Now
	err = c.StopCard(uint(activeCard.Projectid), uint(activeCard.ID))
	if err != nil {
		return err
	}

	return nil
}
