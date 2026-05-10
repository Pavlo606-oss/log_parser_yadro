package repository

import (
	"context"
	"fmt"
	"repo/internal/parser"
)

func (r *Repository) CreateSharpANInfos(ctx context.Context, logID int64, infos []parser.SharpANInfoRow) error {
	query := `
		INSERT INTO sharp_an_info (
			log_id,
			sw_guid,
			endianness,
			enable_endianness_per_job,
			reproducibility_disable
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (log_id, sw_guid) DO UPDATE SET
			endianness = EXCLUDED.endianness,
			enable_endianness_per_job = EXCLUDED.enable_endianness_per_job,
			reproducibility_disable = EXCLUDED.reproducibility_disable
	`

	for _, info := range infos {
		_, err := r.db.ExecContext(
			ctx,
			query,
			logID,
			info.SWGUID,
			info.Endianness,
			info.EnableEndiannessPerJob,
			info.ReproducibilityDisable,
		)
		if err != nil {
			return fmt.Errorf("create sharp_an_info sw_guid=%q: %w", info.SWGUID, err)
		}
	}

	return nil
}
