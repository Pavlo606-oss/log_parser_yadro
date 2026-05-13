package repository

import (
	"context"
	"database/sql"
	"fmt"
	"repo/internal/parser"
)

func (r *Repository) CreateNodeInfos(ctx context.Context, q DBTX, logID int64, nodeIDs map[string]int64, infos []parser.NodeInfoRow) error {
	query := `
		INSERT INTO nodes_info (
			log_id,
			node_id,
			node_guid,
			serial_number,
			part_number,
			revision,
			product_name
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (log_id, node_guid) DO UPDATE SET
			node_id = EXCLUDED.node_id,
			serial_number = EXCLUDED.serial_number,
			part_number = EXCLUDED.part_number,
			revision = EXCLUDED.revision,
			product_name = EXCLUDED.product_name
	`

	for _, info := range infos {
		nodeID, ok := nodeIDs[info.NodeGUID]
		if !ok {
			return fmt.Errorf("create node_info: node id not found for node_guid=%q", info.NodeGUID)
		}

		_, err := q.ExecContext(
			ctx,
			query,
			logID,
			nodeID,
			info.NodeGUID,
			info.SerialNumber,
			info.PartNumber,
			info.Revision,
			info.ProductName,
		)
		if err != nil {
			return fmt.Errorf("create node_info node_guid=%q: %w", info.NodeGUID, err)
		}
	}

	return nil
}

func (r *Repository) GetNodeDetailByID(ctx context.Context, nodeID int64) (*NodeDetail, error) {
	query := `
		SELECT 
			n.id, n.log_id, n.node_guid, n.node_desc, n.num_ports, n.node_type,
			n.class_version, n.base_version, n.system_image_guid, n.port_guid,
			
			ni.id, ni.log_id, ni.node_id, ni.node_guid, ni.serial_number, 
			ni.part_number, ni.revision, ni.product_name,
			
			s.id, s.log_id, s.node_id, s.node_guid, s.linear_fdb_cap, s.random_fdb_cap,
			s.mcast_fdb_cap, s.linear_fdb_top, s.def_port, s.def_mcast_pri_port,
			s.def_mcast_not_pri_port, s.life_time_value, s.port_state_change,
			s.optimized_slvl_mapping, s.lids_per_port, s.part_enf_cap, s.inb_enf_cap,
			s.outb_enf_cap, s.filter_raw_inb_cap, s.filter_raw_outb_cap, s.enp0, s.mcast_fdb_top
		FROM nodes n
		LEFT JOIN nodes_info ni ON ni.node_id = n.id AND n.log_id = ni.log_id
		LEFT JOIN switches s ON s.node_id = n.id AND s.log_id = n.log_id
		WHERE n.id = $1
	`

	var node Node

	var niID sql.NullInt64
	var niLogID sql.NullInt64
	var niNodeID sql.NullInt64
	var niNodeGUID sql.NullString
	var niSerialNumber sql.NullString
	var niPartNumber sql.NullString
	var niRevision sql.NullString
	var niProductName sql.NullString

	var sID sql.NullInt64
	var sLogID sql.NullInt64
	var sNodeID sql.NullInt64
	var sNodeGUID sql.NullString
	var sLinearFDBCap sql.NullInt64
	var sRandomFDBCap sql.NullInt64
	var sMCastFDBCap sql.NullInt64
	var sLinearFDBTop sql.NullInt64
	var sDefPort sql.NullInt64
	var sDefMCastPriPort sql.NullInt64
	var sDefMCastNotPriPort sql.NullInt64
	var sLifeTimeValue sql.NullInt64
	var sPortStateChange sql.NullInt64
	var sOptimizedSLVLMapping sql.NullInt64
	var sLidsPerPort sql.NullInt64
	var sPartEnfCap sql.NullInt64
	var sInbEnfCap sql.NullInt64
	var sOutbEnfCap sql.NullInt64
	var sFilterRawInbCap sql.NullInt64
	var sFilterRawOutbCap sql.NullInt64
	var sENP0 sql.NullInt64
	var sMCastFDBTop sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, nodeID).Scan(
		&node.ID, &node.LogID, &node.NodeGUID, &node.NodeDesc, &node.NumPorts, &node.NodeType,
		&node.ClassVersion, &node.BaseVersion, &node.SystemImageGUID, &node.PortGUID,

		&niID, &niLogID, &niNodeID, &niNodeGUID, &niSerialNumber,
		&niPartNumber, &niRevision, &niProductName,

		&sID, &sLogID, &sNodeID, &sNodeGUID, &sLinearFDBCap, &sRandomFDBCap,
		&sMCastFDBCap, &sLinearFDBTop, &sDefPort, &sDefMCastPriPort,
		&sDefMCastNotPriPort, &sLifeTimeValue, &sPortStateChange,
		&sOptimizedSLVLMapping, &sLidsPerPort, &sPartEnfCap, &sInbEnfCap,
		&sOutbEnfCap, &sFilterRawInbCap, &sFilterRawOutbCap, &sENP0, &sMCastFDBTop,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("node with id=%d not found", nodeID)
		}
		return nil, fmt.Errorf("get node detail by id=%d: %w", nodeID, err)
	}

	detail := &NodeDetail{
		Node: &node,
	}

	if niID.Valid {
		detail.NodeInfo = &NodeInfo{
			ID:           niID.Int64,
			LogID:        niLogID.Int64,
			NodeID:       niNodeID.Int64,
			NodeGUID:     niNodeGUID.String,
			SerialNumber: niSerialNumber.String,
			PartNumber:   niPartNumber.String,
			Revision:     niRevision.String,
			ProductName:  niProductName.String,
		}
	}

	if sID.Valid {
		detail.Switch = &Switch{
			ID:                   sID.Int64,
			LogID:                sLogID.Int64,
			NodeID:               sNodeID.Int64,
			NodeGUID:             sNodeGUID.String,
			LinearFDBCap:         sLinearFDBCap.Int64,
			RandomFDBCap:         sRandomFDBCap.Int64,
			MCastFDBCap:          sMCastFDBCap.Int64,
			LinearFDBTop:         sLinearFDBTop.Int64,
			DefPort:              sDefPort.Int64,
			DefMCastPriPort:      sDefMCastPriPort.Int64,
			DefMCastNotPriPort:   sDefMCastNotPriPort.Int64,
			LifeTimeValue:        sLifeTimeValue.Int64,
			PortStateChange:      sPortStateChange.Int64,
			OptimizedSLVLMapping: sOptimizedSLVLMapping.Int64,
			LidsPerPort:          sLidsPerPort.Int64,
			PartEnfCap:           sPartEnfCap.Int64,
			InbEnfCap:            sInbEnfCap.Int64,
			OutbEnfCap:           sOutbEnfCap.Int64,
			FilterRawInbCap:      sFilterRawInbCap.Int64,
			FilterRawOutbCap:     sFilterRawOutbCap.Int64,
			ENP0:                 sENP0.Int64,
			MCastFDBTop:          sMCastFDBTop.Int64,
		}
	}

	return detail, nil
}
