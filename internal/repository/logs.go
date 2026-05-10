package repository

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *Repository) CreateLog(ctx context.Context, fileName, fileType string) (int64, error) {
	var id int64
	query := "INSERT INTO logs(filename, file_type) VALUES ($1, $2) RETURNING id"
	row := r.db.QueryRowContext(ctx, query, fileName, fileType)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("create log filename=%q file_type=%q: %w", fileName, fileType, err)
	}
	return id, nil
}

func (r *Repository) GetLogMeta(ctx context.Context, logID int64) (*LogMeta, error) {
	query := `
        SELECT
            l.id,
            l.file_name,
            l.file_type,
            l.imported_at,
            (SELECT COUNT(*) FROM nodes WHERE log_id = l.id) AS nodes_count,
            (SELECT COUNT(*) FROM ports WHERE log_id = l.id) AS ports_count
        FROM logs l
        WHERE l.id = $1
    `

	var meta LogMeta
	err := r.db.QueryRowContext(ctx, query, logID).Scan(
		&meta.ID,
		&meta.FileName,
		&meta.FileType,
		&meta.ImportedAt,
		&meta.NodesCount,
		&meta.PortsCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("log with id=%d not found", logID)
		}
		return nil, fmt.Errorf("get log meta by id=%d: %w", logID, err)
	}

	if meta.NodesCount > 0 || meta.PortsCount > 0 {
		meta.Status = "processed"
	} else {
		meta.Status = "empty"
	}

	return &meta, nil
}
