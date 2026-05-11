package repository

import (
	"context"
	"fmt"
	"repo/internal/parser"
)

func (r *Repository) CreatePorts(ctx context.Context, q DBTX, logID int64, nodeIDs map[string]int64, ports []parser.PortRow) error {
	query := `
		INSERT INTO ports (
			log_id,
			node_id,
			node_guid,
			port_guid,
			port_num,
			lid,
			port_phy_state,
			port_state,
			link_width_actv,
			link_speed_actv
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (log_id, node_guid, port_num) DO UPDATE SET
			node_id = EXCLUDED.node_id,
			port_guid = EXCLUDED.port_guid,
			lid = EXCLUDED.lid,
			port_phy_state = EXCLUDED.port_phy_state,
			port_state = EXCLUDED.port_state,
			link_width_actv = EXCLUDED.link_width_actv,
			link_speed_actv = EXCLUDED.link_speed_actv
	`

	for _, port := range ports {
		nodeID, ok := nodeIDs[port.NodeGUID]
		if !ok {
			return fmt.Errorf("create port: node id not found for node_guid=%q", port.NodeGUID)
		}

		_, err := q.ExecContext(
			ctx,
			query,
			logID,
			nodeID,
			port.NodeGUID,
			port.PortGUID,
			port.PortNum,
			port.LID,
			port.PortPhyState,
			port.PortState,
			port.LinkWidthActv,
			port.LinkSpeedActv,
		)
		if err != nil {
			return fmt.Errorf(
				"create port node_guid=%q port_num=%d: %w",
				port.NodeGUID,
				port.PortNum,
				err,
			)
		}
	}
	return nil
}

func (r *Repository) GetPortsByNodeID(ctx context.Context, nodeID int64) ([]Port, error) {
	result := make([]Port, 0)
	query := `
		SELECT 
			id,
    		log_id,
    		node_id,
    		node_guid,
    		port_guid,
    		port_num,
			lid,
    		port_phy_state,
    		port_state,
    		link_width_actv,
    		link_speed_actv
		FROM ports
		WHERE node_id = $1
		ORDER BY port_num`

	rows, err := r.db.QueryContext(ctx, query, nodeID)

	if err != nil {
		return nil, fmt.Errorf("get port by node_id = %d: %w", nodeID, err)
	}

	defer rows.Close()

	for rows.Next() {
		port := Port{}
		err := rows.Scan(
			&port.ID,
			&port.LogID,
			&port.NodeID,
			&port.NodeGUID,
			&port.PortGUID,
			&port.PortNum,
			&port.LID,
			&port.PortPhyState,
			&port.PortState,
			&port.LinkWidthActv,
			&port.LinkSpeedActv,
		)

		if err != nil {
			return nil, fmt.Errorf("scan port: %w", err)
		}
		result = append(result, port)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get port by node_id = %d: %w", nodeID, err)
	}

	return result, nil
}
