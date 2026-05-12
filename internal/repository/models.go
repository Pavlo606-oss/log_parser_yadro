package repository

import "time"

type Log struct {
	ID         int64     `json:"id"`
	FileName   string    `json:"file_name"`
	FileType   string    `json:"file_type"`
	ImportedAt time.Time `json:"imported_at"`
}

type Node struct {
	ID              int64  `json:"id"`
	LogID           int64  `json:"log_id"`
	NodeGUID        string `json:"node_guid"`
	NodeDesc        string `json:"node_desc"`
	NumPorts        int64  `json:"num_ports"`
	NodeType        int64  `json:"node_type"`
	ClassVersion    int64  `json:"class_version"`
	BaseVersion     int64  `json:"base_version"`
	SystemImageGUID string `json:"system_image_guid"`
	PortGUID        string `json:"port_guid"`
}

type Port struct {
	ID            int64  `json:"id"`
	LogID         int64  `json:"log_id"`
	NodeID        int64  `json:"node_id"`
	NodeGUID      string `json:"node_guid"`
	PortGUID      string `json:"port_guid"`
	PortNum       int64  `json:"port_num"`
	LID           int64  `json:"lid"`
	PortPhyState  int64  `json:"port_phy_state"`
	PortState     int64  `json:"port_state"`
	LinkWidthActv int64  `json:"link_width_actv"`
	LinkSpeedActv int64  `json:"link_speed_actv"`
}

type Switch struct {
	ID                   int64  `json:"id"`
	LogID                int64  `json:"log_id"`
	NodeID               int64  `json:"node_id"`
	NodeGUID             string `json:"node_guid"`
	LinearFDBCap         int64  `json:"linear_fdb_cap"`
	RandomFDBCap         int64  `json:"random_fdb_cap"`
	MCastFDBCap          int64  `json:"mcast_fdb_cap"`
	LinearFDBTop         int64  `json:"linear_fdb_top"`
	DefPort              int64  `json:"def_port"`
	DefMCastPriPort      int64  `json:"def_mcast_pri_port"`
	DefMCastNotPriPort   int64  `json:"def_mcast_not_pri_port"`
	LifeTimeValue        int64  `json:"life_time_value"`
	PortStateChange      int64  `json:"port_state_change"`
	OptimizedSLVLMapping int64  `json:"optimized_slvl_mapping"`
	LidsPerPort          int64  `json:"lids_per_port"`
	PartEnfCap           int64  `json:"part_enf_cap"`
	InbEnfCap            int64  `json:"inb_enf_cap"`
	OutbEnfCap           int64  `json:"outb_enf_cap"`
	FilterRawInbCap      int64  `json:"filter_raw_inb_cap"`
	FilterRawOutbCap     int64  `json:"filter_raw_outb_cap"`
	ENP0                 int64  `json:"enp0"`
	MCastFDBTop          int64  `json:"mcast_fdb_top"`
}

type NodeInfo struct {
	ID           int64  `json:"id"`
	LogID        int64  `json:"log_id"`
	NodeID       int64  `json:"node_id"`
	NodeGUID     string `json:"node_guid"`
	SerialNumber string `json:"serial_number"`
	PartNumber   string `json:"part_number"`
	Revision     string `json:"revision"`
	ProductName  string `json:"product_name"`
}

type SharpANInfo struct {
	ID                     int64  `json:"id"`
	LogID                  int64  `json:"log_id"`
	SWGUID                 string `json:"sw_guid"`
	Endianness             int64  `json:"endianness"`
	EnableEndiannessPerJob int64  `json:"enable_endianness_per_job"`
	ReproducibilityDisable int64  `json:"reproducibility_disable"`
}

type NodeDetail struct {
	Node     *Node     `json:"node"`
	Switch   *Switch   `json:"switch,omitempty"`
	NodeInfo *NodeInfo `json:"node_info,omitempty"`
}

type LogMeta struct {
	ID         int64     `json:"id"`
	FileName   string    `json:"file_name"`
	FileType   string    `json:"file_type"`
	ImportedAt time.Time `json:"imported_at"`
	Status     string    `json:"status"`
	NodesCount int       `json:"nodes_count"`
	PortsCount int       `json:"ports_count"`
}
