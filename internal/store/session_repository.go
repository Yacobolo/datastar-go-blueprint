package store

import (
	"context"

	"github.com/yacobolo/datastar-go-blueprint/internal/domain"
	"github.com/yacobolo/datastar-go-blueprint/internal/store/queries"
)

// SessionRepository is the concrete implementation of domain.SessionRepository.
// It wraps sqlc-generated queries and acts as a driven adapter in hexagonal architecture.
type SessionRepository struct {
	store *SQLiteStore
}

// Ensure SessionRepository implements domain.SessionRepository at compile time.
var _ domain.SessionRepository = (*SessionRepository)(nil)

// NewSessionRepository creates a new SessionRepository instance.
func NewSessionRepository(st *SQLiteStore) *SessionRepository {
	return &SessionRepository{store: st}
}

// GetSession retrieves a session by its ID.
func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (queries.Session, error) {
	return r.store.Queries().GetSession(ctx, sessionID)
}

// UpsertSession inserts or updates a session.
func (r *SessionRepository) UpsertSession(ctx context.Context, arg queries.UpsertSessionParams) error {
	return r.store.Queries().UpsertSession(ctx, arg)
}
