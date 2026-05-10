package repository

import "time"

type Log struct {
	Id         int64
	FileName   string
	FileType   string
	ImportedAt time.Time
}

type Node struct {
	ID              int64
	LogID           int64
	NodeGUID        string
	NodeDesc        string
	NumPorts        int64
	NodeType        int64
	ClassVersion    int64
	BaseVersion     int64
	SystemImageGUID string
	PortGUID        string
}

type Port struct {
	ID            int64
	LogID         int64
	NodeID        int64
	NodeGUID      string
	PortGUID      string
	PortNum       int64
	LID           int64
	PortPhyState  int64
	PortState     int64
	LinkWidthActv int64
	LinkSpeedActv int64
}

type Switch struct {
	ID                   int64
	LogID                int64
	NodeID               int64
	NodeGUID             string
	LinearFDBCap         int64
	RandomFDBCap         int64
	MCastFDBCap          int64
	LinearFDBTop         int64
	DefPort              int64
	DefMCastPriPort      int64
	DefMCastNotPriPort   int64
	LifeTimeValue        int64
	PortStateChange      int64
	OptimizedSLVLMapping int64
	LidsPerPort          int64
	PartEnfCap           int64
	InbEnfCap            int64
	OutbEnfCap           int64
	FilterRawInbCap      int64
	FilterRawOutbCap     int64
	ENP0                 int64
	MCastFDBTop          int64
}

type NodeInfo struct {
	ID           int64
	LogID        int64
	NodeID       int64
	NodeGUID     string
	SerialNumber string
	PartNumber   string
	Revision     string
	ProductName  string
}

type SharpANInfo struct {
	ID                     int64
	LogID                  int64
	SWGUID                 string
	Endianness             int64
	EnableEndiannessPerJob int64
	ReproducibilityDisable int64
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
