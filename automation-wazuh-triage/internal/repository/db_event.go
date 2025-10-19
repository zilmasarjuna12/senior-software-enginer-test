package repository

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/internal/entity"
	"automation-wazuh-triage/pkg/logger"
	"context"
	"database/sql"
)

type closedEventRepository struct {
	db *sql.DB
}

func NewClosedEventRepository(db *sql.DB) domain.ClosedEventRepository {
	return &closedEventRepository{
		db: db,
	}
}

func (r *closedEventRepository) SaveClosedEvent(ctx context.Context, closedEvent *entity.ClosedEvent) error {
	log := logger.WithRequestID(ctx)

	query := `
		INSERT INTO closed_events (event_id, rule_id, raw_event, reason, status, close_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		closedEvent.EventID,
		closedEvent.RuleID,
		closedEvent.RawEvent,
		closedEvent.Reason,
		closedEvent.Status,
		closedEvent.CloseAt,
	)

	if err != nil {
		log.WithError(err).Error("[repository - event - SaveClosedEvent]: Failed to save closed event")
		return err
	}

	log.Info("[repository - event - SaveClosedEvent]: Successfully saved closed event")
	return nil
}

func (r *closedEventRepository) FetchClosedEvents(ctx context.Context) ([]*entity.ClosedEvent, error) {
	log := logger.WithRequestID(ctx)

	query := `
		SELECT id, event_id, rule_id, raw_event, reason, status, close_at
		FROM closed_events
		ORDER BY close_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.WithError(err).Error("[repository - event - FetchClosedEvents]: Failed to fetch closed events")
		return nil, err
	}
	defer rows.Close()

	var closedEvents []*entity.ClosedEvent

	for rows.Next() {
		var event entity.ClosedEvent
		err := rows.Scan(
			&event.ID,
			&event.EventID,
			&event.RuleID,
			&event.RawEvent,
			&event.Reason,
			&event.Status,
			&event.CloseAt,
		)
		if err != nil {
			log.WithError(err).Error("[repository - event - FetchClosedEvents]: Failed to scan closed event")
			return nil, err
		}
		closedEvents = append(closedEvents, &event)
	}

	if err = rows.Err(); err != nil {
		log.WithError(err).Error("[repository - event - FetchClosedEvents]: Error iterating rows")
		return nil, err
	}

	log.WithField("count", len(closedEvents)).Info("[repository - event - FetchClosedEvents]: Successfully fetched closed events")
	return closedEvents, nil
}

func (r *closedEventRepository) FetchClosedEventByID(ctx context.Context, id string) (*entity.ClosedEvent, error) {
	log := logger.WithRequestID(ctx)

	query := `
		SELECT id, event_id, rule_id, raw_event, reason, status, close_at
		FROM closed_events
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var event entity.ClosedEvent
	err := row.Scan(
		&event.ID,
		&event.EventID,
		&event.RuleID,
		&event.RawEvent,
		&event.Reason,
		&event.Status,
		&event.CloseAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.WithField("id", id).Warn("[repository - event - FetchClosedEventByID]: Closed event not found")
			return nil, nil // Return nil to indicate not found
		}
		log.WithError(err).Error("[repository - event - FetchClosedEventByID]: Failed to fetch closed event by ID")
		return nil, err
	}

	log.WithField("id", id).Info("[repository - event - FetchClosedEventByID]: Successfully fetched closed event by ID")
	return &event, nil
}

func (r *closedEventRepository) FetchClosedEventByEventID(ctx context.Context, eventID string) (*entity.ClosedEvent, error) {
	log := logger.WithRequestID(ctx)

	query := `
		SELECT id, event_id, rule_id, raw_event, reason, status, close_at
		FROM closed_events
		WHERE event_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, eventID)

	var event entity.ClosedEvent
	err := row.Scan(
		&event.ID,
		&event.EventID,
		&event.RuleID,
		&event.RawEvent,
		&event.Reason,
		&event.Status,
		&event.CloseAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.WithField("event_id", eventID).Debug("[repository - event - FetchClosedEventByEventID]: Closed event not found")
			return nil, nil // Return nil to indicate not found
		}
		log.WithError(err).Error("[repository - event - FetchClosedEventByEventID]: Failed to fetch closed event by event ID")
		return nil, err
	}

	log.WithField("event_id", eventID).Info("[repository - event - FetchClosedEventByEventID]: Successfully fetched closed event by event ID")
	return &event, nil
}

func (r *closedEventRepository) UpdateClosedEventReason(ctx context.Context, id string, reason string) error {
	log := logger.WithRequestID(ctx)

	query := `
		UPDATE closed_events 
		SET reason = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, reason, id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("[repository - event - UpdateClosedEventReason]: Failed to update closed event reason")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).WithField("id", id).Error("[repository - event - UpdateClosedEventReason]: Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		log.WithField("id", id).Warn("[repository - event - UpdateClosedEventReason]: No closed event found with the given ID")
		return sql.ErrNoRows
	}

	log.WithField("id", id).WithField("reason", reason).Info("[repository - event - UpdateClosedEventReason]: Successfully updated closed event reason")
	return nil
}
