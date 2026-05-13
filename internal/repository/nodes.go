package repository

import (
	"context"
	"database/sql"
	"fmt"
	"repo/internal/parser"
	"sort"
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

func (r *Repository) GetNodesTopology(ctx context.Context, logID int64) (*NodesTopology, error) {
	result := &NodesTopology{Items: make([]*NodeTopology, 0)}
	m := make(map[Node][]*Port)
	query := `
		SELECT 
			n.id,
    		n.log_id,
    		n.node_guid,
    		n.node_desc,
    		n.num_ports,
    		n.node_type,
    		n.class_version,
    		n.base_version,
    		n.system_image_guid,
    		n.port_guid,
			p.id,
    		p.log_id,
    		p.node_id,
    		p.node_guid,
    		p.port_guid,
    		p.port_num,
    		p.lid,
    		p.port_phy_state,
    		p.port_state,
    		p.link_width_actv,
    		p.link_speed_actv
			FROM nodes n  
		LEFT JOIN ports p ON n.id = p.node_id
		WHERE n.log_id = $1 `
	rows, err := r.db.QueryContext(ctx, query, logID)

	if err != nil {
		return nil, fmt.Errorf("get nodes by log_id = %d: %w", logID, err)
	}

	defer rows.Close()

	for rows.Next() {
		var portID sql.NullInt64
		var portLogID sql.NullInt64
		var portNodeID sql.NullInt64
		var portNodeGUID sql.NullString
		var portGUID sql.NullString
		var portNum sql.NullInt64
		var portLID sql.NullInt64
		var portPhyState sql.NullInt64
		var portState sql.NullInt64
		var portLinkWidthActv sql.NullInt64
		var portLinkSpeedActv sql.NullInt64

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
			&portID,
			&portLogID,
			&portNodeID,
			&portNodeGUID,
			&portGUID,
			&portNum,
			&portLID,
			&portPhyState,
			&portState,
			&portLinkWidthActv,
			&portLinkSpeedActv,
		)

		if err != nil {
			return nil, fmt.Errorf("scan node for log_id=%d: %w", logID, err)
		}

		if _, ok := m[node]; !ok {
			m[node] = make([]*Port, 0)
		}

		if portID.Valid {
			m[node] = append(
				m[node],
				&Port{
					ID:            portID.Int64,
					LogID:         portLogID.Int64,
					NodeID:        portNodeID.Int64,
					NodeGUID:      portNodeGUID.String,
					PortGUID:      portGUID.String,
					PortNum:       portNum.Int64,
					LID:           portLID.Int64,
					PortPhyState:  portPhyState.Int64,
					PortState:     portState.Int64,
					LinkWidthActv: portLinkWidthActv.Int64,
					LinkSpeedActv: portLinkSpeedActv.Int64,
				},
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nodes by log_id=%d: %w", logID, err)
	}

	nodesTopology := make([]*NodeTopology, 0)
	for node, ports := range m {
		nodeCopy := node
		nodesTopology = append(nodesTopology, &NodeTopology{Node: &nodeCopy, Ports: ports})
	}

	sort.Slice(nodesTopology, func(i, j int) bool {
		return nodesTopology[i].Node.ID < nodesTopology[j].Node.ID
	})

	result.Items = nodesTopology
	return result, nil
}
