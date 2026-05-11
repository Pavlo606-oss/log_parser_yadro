package repository

import (
	"context"
	"fmt"
	"repo/internal/parser"
)

func (r *Repository) CreateNodes(ctx context.Context, q DBTX, logID int64, nodes []parser.NodeRow) (map[string]int64, error) {
	m := make(map[string]int64)
	for _, node := range nodes {
		var id int64
		query := `
    INSERT INTO nodes (log_id, node_guid, node_desc, num_ports, node_type, class_version, base_version, system_image_guid, port_guid)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    ON CONFLICT (log_id, node_guid) DO UPDATE SET
        node_desc = EXCLUDED.node_desc,
        num_ports = EXCLUDED.num_ports,
        node_type = EXCLUDED.node_type,
        class_version = EXCLUDED.class_version,
        base_version = EXCLUDED.base_version,
        system_image_guid = EXCLUDED.system_image_guid,
        port_guid = EXCLUDED.port_guid
    RETURNING id`
		row := q.QueryRowContext(ctx, query,
			logID,
			node.NodeGUID,
			node.NodeDesc,
			node.NumPorts,
			node.NodeType,
			node.ClassVersion,
			node.BaseVersion,
			node.SystemImageGUID,
			node.PortGUID,
		)
		if err := row.Scan(&id); err != nil {
			return nil, fmt.Errorf("create node = %q : %w", node.NodeGUID, err)
		}
		m[node.NodeGUID] = id
	}
	return m, nil
}

func (r *Repository) GetNodesByLogID(ctx context.Context, q DBTX, logID int64) ([]Node, error) {
	result := make([]Node, 0)
	query := `
		SELECT 
			id,
    		log_id,
    		node_guid,
    		node_desc,
    		num_ports,
    		node_type,
    		class_version,
    		base_version,
    		system_image_guid,
    		port_guid
			FROM nodes 
		WHERE log_id = $1 `
	rows, err := q.QueryContext(ctx, query, logID)

	if err != nil {
		return nil, fmt.Errorf("get nodes by log_id = %d: %w", logID, err)
	}

	defer rows.Close()

	for rows.Next() {
		node := Node{}
		err := rows.Scan(
			&node.ID,
			&node.LogID,
			&node.NodeGUID,
			&node.NodeDesc,
			&node.NumPorts,
			&node.NodeType,
			&node.ClassVersion,
			&node.BaseVersion,
			&node.SystemImageGUID,
			&node.PortGUID,
		)

		if err != nil {
			return nil, fmt.Errorf("scan node for log_id=%d: %w", logID, err)
		}
		result = append(result, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nodes by log_id=%d: %w", logID, err)
	}
	return result, nil
}
