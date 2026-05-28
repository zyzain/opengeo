package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"opengeo/service/brand/internal/domain"
)

type BrandMetadataRepository struct{ db *sql.DB }

func NewBrandMetadataRepository(db *sql.DB) *BrandMetadataRepository {
	return &BrandMetadataRepository{db: db}
}

func (r *BrandMetadataRepository) Save(ctx context.Context, m *domain.BrandMetadata) error {
	viJSON, _ := json.Marshal(m.VIProfile)
	toneJSON, _ := json.Marshal(m.ToneProfile)
	audJSON, _ := json.Marshal(m.AudienceProfiles)
	compJSON, _ := json.Marshal(m.CompetitorList)

	query := `INSERT INTO brand_metadata (brand_id, vi_profile, tone_profile, audience_profiles, competitor_list, brand_values, unique_selling_points, schema_version, created_at, updated_at)
		VALUES (?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE vi_profile=VALUES(vi_profile), tone_profile=VALUES(tone_profile), audience_profiles=VALUES(audience_profiles), competitor_list=VALUES(competitor_list), brand_values=VALUES(brand_values), unique_selling_points=VALUES(unique_selling_points), updated_at=VALUES(updated_at)`

	valuesJSON, _ := json.Marshal(m.BrandValues)
	uspsJSON, _ := json.Marshal(m.UniqueSellingPoints)

	_, err := r.db.ExecContext(ctx, query, m.BrandID, string(viJSON), string(toneJSON), string(audJSON), string(compJSON), string(valuesJSON), string(uspsJSON), m.SchemaVersion, m.CreatedAt, m.UpdatedAt)
	return err
}

func (r *BrandMetadataRepository) FindByBrandID(ctx context.Context, tenantID, brandID int64) (*domain.BrandMetadata, error) {
	query := `SELECT m.brand_id, m.vi_profile, m.tone_profile, m.audience_profiles, m.competitor_list, m.brand_values, m.unique_selling_points, m.schema_version, m.created_at, m.updated_at
		FROM brand_metadata m JOIN brands b ON m.brand_id=b.id WHERE m.brand_id=? AND b.tenant_id=?`

	m := &domain.BrandMetadata{}
	var viJSON, toneJSON, audJSON, compJSON, valsJSON, uspsJSON sql.NullString

	err := r.db.QueryRowContext(ctx, query, brandID, tenantID).Scan(&m.BrandID, &viJSON, &toneJSON, &audJSON, &compJSON, &valsJSON, &uspsJSON, &m.SchemaVersion, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if viJSON.Valid {
		json.Unmarshal([]byte(viJSON.String), &m.VIProfile)
	}
	if toneJSON.Valid {
		json.Unmarshal([]byte(toneJSON.String), &m.ToneProfile)
	}
	if audJSON.Valid {
		json.Unmarshal([]byte(audJSON.String), &m.AudienceProfiles)
	}
	if compJSON.Valid {
		json.Unmarshal([]byte(compJSON.String), &m.CompetitorList)
	}
	if valsJSON.Valid {
		json.Unmarshal([]byte(valsJSON.String), &m.BrandValues)
	}
	if uspsJSON.Valid {
		json.Unmarshal([]byte(uspsJSON.String), &m.UniqueSellingPoints)
	}
	return m, nil
}

func (r *BrandMetadataRepository) Delete(ctx context.Context, tenantID, brandID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM brand_metadata WHERE brand_id=? AND brand_id IN (SELECT id FROM brands WHERE tenant_id=?)", brandID, tenantID)
	return err
}
