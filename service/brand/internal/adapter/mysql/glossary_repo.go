package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"opengeo/service/brand/internal/domain"
	"opengeo/service/brand/internal/port"
)

type GlossaryRepository struct{ db *sql.DB }

func NewGlossaryRepository(db *sql.DB) *GlossaryRepository {
	return &GlossaryRepository{db: db}
}

func (r *GlossaryRepository) Save(ctx context.Context, e *domain.GlossaryEntry) error {
	if e.ID == 0 {
		return r.create(ctx, e)
	}
	return r.update(ctx, e)
}

func (r *GlossaryRepository) create(ctx context.Context, e *domain.GlossaryEntry) error {
	aliasesJSON, _ := json.Marshal(e.Aliases)
	query := `INSERT INTO glossary_entries (brand_id, term, definition, category, aliases, context, is_forbidden, is_preferred, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?,?)`
	result, err := r.db.ExecContext(ctx, query, e.BrandID, e.Term, e.Definition, e.Category, string(aliasesJSON), e.Context, e.IsForbidden, e.IsPreferred, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	e.ID = id
	return nil
}

func (r *GlossaryRepository) update(ctx context.Context, e *domain.GlossaryEntry) error {
	aliasesJSON, _ := json.Marshal(e.Aliases)
	query := `UPDATE glossary_entries SET term=?, definition=?, category=?, aliases=?, context=?, is_forbidden=?, is_preferred=?, updated_at=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, e.Term, e.Definition, e.Category, string(aliasesJSON), e.Context, e.IsForbidden, e.IsPreferred, e.UpdatedAt, e.ID)
	return err
}

func (r *GlossaryRepository) FindByID(ctx context.Context, tenantID, id int64) (*domain.GlossaryEntry, error) {
	query := `SELECT g.id, g.brand_id, g.term, g.definition, g.category, g.aliases, g.context, g.is_forbidden, g.is_preferred, g.created_at, g.updated_at
		FROM glossary_entries g JOIN brands b ON g.brand_id=b.id WHERE g.id=? AND b.tenant_id=?`
	e := &domain.GlossaryEntry{}
	var aliasesJSON sql.NullString
	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(&e.ID, &e.BrandID, &e.Term, &e.Definition, &e.Category, &aliasesJSON, &e.Context, &e.IsForbidden, &e.IsPreferred, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if aliasesJSON.Valid {
		json.Unmarshal([]byte(aliasesJSON.String), &e.Aliases)
	}
	return e, nil
}

func (r *GlossaryRepository) List(ctx context.Context, tenantID, brandID int64, filter *port.GlossaryFilter) ([]*domain.GlossaryEntry, int32, error) {
	where := "WHERE g.brand_id=? AND b.tenant_id=?"
	args := []interface{}{brandID, tenantID}
	if filter.Category != "" {
		where += " AND g.category=?"
		args = append(args, filter.Category)
	}
	if filter.Keyword != "" {
		where += " AND (g.term LIKE ? OR g.definition LIKE ?)"
		kw := "%" + filter.Keyword + "%"
		args = append(args, kw, kw)
	}
	if filter.IsForbidden != nil {
		where += " AND g.is_forbidden=?"
		args = append(args, *filter.IsForbidden)
	}

	var total int32
	r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM glossary_entries g JOIN brands b ON g.brand_id=b.id %s", where), args...).Scan(&total)

	offset := (filter.Page - 1) * filter.PageSize
	query := fmt.Sprintf(`SELECT g.id, g.brand_id, g.term, g.definition, g.category, g.aliases, g.context, g.is_forbidden, g.is_preferred, g.created_at, g.updated_at FROM glossary_entries g JOIN brands b ON g.brand_id=b.id %s ORDER BY g.created_at DESC LIMIT ? OFFSET ?`, where)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	entries := make([]*domain.GlossaryEntry, 0)
	for rows.Next() {
		e := &domain.GlossaryEntry{}
		var aliasesJSON sql.NullString
		rows.Scan(&e.ID, &e.BrandID, &e.Term, &e.Definition, &e.Category, &aliasesJSON, &e.Context, &e.IsForbidden, &e.IsPreferred, &e.CreatedAt, &e.UpdatedAt)
		if aliasesJSON.Valid {
			json.Unmarshal([]byte(aliasesJSON.String), &e.Aliases)
		}
		entries = append(entries, e)
	}
	return entries, total, nil
}

func (r *GlossaryRepository) Delete(ctx context.Context, tenantID, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM glossary_entries WHERE id=? AND brand_id IN (SELECT id FROM brands WHERE tenant_id=?)", id, tenantID)
	return err
}

func (r *GlossaryRepository) BulkCreate(ctx context.Context, entries []*domain.GlossaryEntry) (int32, int32, error) {
	var imported, skipped int32
	for _, e := range entries {
		exists, _ := r.ExistsByTerm(ctx, 0, e.BrandID, e.Term)
		if exists {
			skipped++
			continue
		}
		if err := r.Save(ctx, e); err != nil {
			continue
		}
		imported++
	}
	return imported, skipped, nil
}

func (r *GlossaryRepository) ExistsByTerm(ctx context.Context, tenantID, brandID int64, term string) (bool, error) {
	var exists bool
	r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM glossary_entries WHERE brand_id=? AND term=?)", brandID, term).Scan(&exists)
	return exists, nil
}
