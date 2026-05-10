package repository

import (
	"context"
	"fmt"
	"repo/internal/parser"
)

func (r *Repository) CreateSwitches(ctx context.Context, logID int64, nodeIDs map[string]int64, switches []parser.SwitchRow) error {
	query := `
		INSERT INTO switches (
			log_id,
			node_id,
			node_guid,
			linear_fdb_cap,
			random_fdb_cap,
			mcast_fdb_cap,
			linear_fdb_top,
			def_port,
			def_mcast_pri_port,
			def_mcast_not_pri_port,
			life_time_value,
			port_state_change,
			optimized_slvl_mapping,
			lids_per_port,
			part_enf_cap,
			inb_enf_cap,
			outb_enf_cap,
			filter_raw_inb_cap,
			filter_raw_outb_cap,
			enp0,
			mcast_fdb_top
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
		ON CONFLICT (log_id, node_guid) DO UPDATE SET
			node_id = EXCLUDED.node_id,
			linear_fdb_cap = EXCLUDED.linear_fdb_cap,
			random_fdb_cap = EXCLUDED.random_fdb_cap,
			mcast_fdb_cap = EXCLUDED.mcast_fdb_cap,
			linear_fdb_top = EXCLUDED.linear_fdb_top,
			def_port = EXCLUDED.def_port,
			def_mcast_pri_port = EXCLUDED.def_mcast_pri_port,
			def_mcast_not_pri_port = EXCLUDED.def_mcast_not_pri_port,
			life_time_value = EXCLUDED.life_time_value,
			port_state_change = EXCLUDED.port_state_change,
			optimized_slvl_mapping = EXCLUDED.optimized_slvl_mapping,
			lids_per_port = EXCLUDED.lids_per_port,
			part_enf_cap = EXCLUDED.part_enf_cap,
			inb_enf_cap = EXCLUDED.inb_enf_cap,
			outb_enf_cap = EXCLUDED.outb_enf_cap,
			filter_raw_inb_cap = EXCLUDED.filter_raw_inb_cap,
			filter_raw_outb_cap = EXCLUDED.filter_raw_outb_cap,
			enp0 = EXCLUDED.enp0,
			mcast_fdb_top = EXCLUDED.mcast_fdb_top
	`

	for _, sw := range switches {
		nodeID, ok := nodeIDs[sw.NodeGUID]
		if !ok {
			return fmt.Errorf("create switch: node id not found for node_guid=%q", sw.NodeGUID)
		}

		_, err := r.db.ExecContext(
			ctx,
			query,
			logID,
			nodeID,
			sw.NodeGUID,
			sw.LinearFDBCap,
			sw.RandomFDBCap,
			sw.MCastFDBCap,
			sw.LinearFDBTop,
			sw.DefPort,
			sw.DefMCastPriPort,
			sw.DefMCastNotPriPort,
			sw.LifeTimeValue,
			sw.PortStateChange,
			sw.OptimizedSLVLMapping,
			sw.LidsPerPort,
			sw.PartEnfCap,
			sw.InbEnfCap,
			sw.OutbEnfCap,
			sw.FilterRawInbCap,
			sw.FilterRawOutbCap,
			sw.ENP0,
			sw.MCastFDBTop,
		)
		if err != nil {
			return fmt.Errorf("create switch node_guid=%q: %w", sw.NodeGUID, err)
		}
	}

	return nil
}
