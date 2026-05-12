package parser

type DBCSVResult struct {
	Nodes     []NodeRow
	Ports     []PortRow
	Switches  []SwitchRow
	NodesInfo []NodeInfoRow
}

type NodeRow struct {
	NodeDesc        string
	NumPorts        int
	NodeType        int
	ClassVersion    int
	BaseVersion     int
	SystemImageGUID string
	NodeGUID        string
	PortGUID        string
}

type PortRow struct {
	NodeGUID      string
	PortGUID      string
	PortNum       int
	LID           int
	PortPhyState  int
	PortState     int
	LinkWidthActv int
	LinkSpeedActv int
}

type SwitchRow struct {
	NodeGUID             string
	LinearFDBCap         int
	RandomFDBCap         int
	MCastFDBCap          int
	LinearFDBTop         int
	DefPort              int
	DefMCastPriPort      int
	DefMCastNotPriPort   int
	LifeTimeValue        int
	PortStateChange      int
	OptimizedSLVLMapping int
	LidsPerPort          int
	PartEnfCap           int
	InbEnfCap            int
	OutbEnfCap           int
	FilterRawInbCap      int
	FilterRawOutbCap     int
	ENP0                 int
	MCastFDBTop          int
}

type NodeInfoRow struct {
	NodeGUID     string
	SerialNumber string
	PartNumber   string
	Revision     string
	ProductName  string
}

type SharpANInfoRow struct {
	SWGUID                 string
	Endianness             int
	EnableEndiannessPerJob int
	ReproducibilityDisable int
}
