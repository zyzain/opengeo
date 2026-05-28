package mysql

import (
	"context"
	"database/sql"
	"time"

	"opengeo/service/brand/internal/domain"
)

type SnapshotRepository struct{ db *sql.DB }

func NewSnapshotRepository(db *sql.DB) *SnapshotRepository {
	return &SnapshotRepository{db: db}
}

func (r *SnapshotRepository) Save(ctx context.Context, s *domain.BrandSnapshot) error {
	query := `INSERT INTO brand_snapshots (brand_id, version, snapshot_data, change_log, created_by, created_at) VALUES (?,?,?,?,?,?)`
	result, err := r.db.ExecContext(ctx, query, s.BrandID, s.Version, s.SnapshotData, s.ChangeLog, s.CreatedBy, s.CreatedAt)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	s.ID = id
	return nil
}

func (r *SnapshotRepository) FindByID(ctx context.Context, tenantID, id int64) (*domain.BrandSnapshot, error) {
	query := `SELECT s.id, s.brand_id, s.version, s.snapshot_data, s.change_log, s.created_by, s.created_at FROM brand_snapshots s JOIN brands b ON s.brand_id=b.id WHERE s.id=? AND b.tenant_id=?`
	s := &domain.BrandSnapshot{}
	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(&s.ID, &s.BrandID, &s.Version, &s.SnapshotData, &s.ChangeLog, &s.CreatedBy, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *SnapshotRepository) List(ctx context.Context, tenantID, brandID int64, page, pageSize int32) ([]*domain.BrandSnapshot, int32, error) {
	var total int32
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM brand_snapshots s JOIN brands b ON s.brand_id=b.id WHERE s.brand_id=? AND b.tenant_id=?", brandID, tenantID).Scan(&total)

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, "SELECT s.id, s.brand_id, s.version, s.snapshot_data, s.change_log, s.created_by, s.created_at FROM brand_snapshots s JOIN brands b ON s.brand_id=b.id WHERE s.brand_id=? AND b.tenant_id=? ORDER BY s.created_at DESC LIMIT ? OFFSET ?", brandID, tenantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	snapshots := make([]*domain.BrandSnapshot, 0)
	for rows.Next() {
		s := &domain.BrandSnapshot{}
		rows.Scan(&s.ID, &s.BrandID, &s.Version, &s.SnapshotData, &s.ChangeLog, &s.CreatedBy, &s.CreatedAt)
		snapshots = append(snapshots, s)
	}
	return snapshots, total, nil
}
