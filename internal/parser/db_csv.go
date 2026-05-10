package parser

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ParseStatusCSV int

// Статусы парсинга
const (
	DefaultStatusCSV ParseStatusCSV = iota
	NodeParsingCSV
	PortParsingCSV
	SwitchParsingCSV
	NodeInfoParsingCSV
)

func ParseDBCSV(r io.Reader) (*DBCSVResult, error) {
	var status ParseStatusCSV
	var skipHeader bool
	var stringLen int
	result := &DBCSVResult{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch line {
		case "START_NODES":
			status = NodeParsingCSV
			skipHeader = true
			continue
		case "START_PORTS":
			status = PortParsingCSV
			skipHeader = true
			continue
		case "START_SWITCHES":
			status = SwitchParsingCSV
			skipHeader = true
			continue
		case "START_SYSTEM_GENERAL_INFORMATION":
			status = NodeInfoParsingCSV
			skipHeader = true
			continue
		case "END_NODES", "END_PORTS", "END_SWITCHES", "END_SYSTEM_GENERAL_INFORMATION":
			if status == DefaultStatusCSV {
				return nil, errors.New("invalid close section")
			}
			status = DefaultStatusCSV
			continue
		}

		if status == DefaultStatusCSV {
			continue
		}

		if skipHeader {
			recHeader, err := parseCSVLine(line)
			if err != nil {
				return nil, fmt.Errorf("parse header: %w", err)
			}
			stringLen = len(recHeader)
			skipHeader = false
			continue
		}

		switch status {
		case NodeParsingCSV:
			parseNode, err := parseNodes(line, stringLen)
			if err != nil {
				return nil, fmt.Errorf("parse node row: %w", err)
			}
			result.Nodes = append(result.Nodes, parseNode)
		case PortParsingCSV:
			parsePort, err := parsePorts(line, stringLen)
			if err != nil {
				return nil, fmt.Errorf("parse port row: %w", err)
			}
			result.Ports = append(result.Ports, parsePort)
		case SwitchParsingCSV:
			parseSwitch, err := parseSwitches(line, stringLen)
			if err != nil {
				return nil, fmt.Errorf("parse switch row: %w", err)
			}
			result.Switches = append(result.Switches, parseSwitch)
		case NodeInfoParsingCSV:
			parseInfo, err := parseNodeInfo(line, stringLen)
			if err != nil {
				return nil, fmt.Errorf("parse node info row: %w", err)
			}
			result.NodesInfo = append(result.NodesInfo, parseInfo)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan db_csv: %w", err)
	}

	return result, nil
}

func parseNodes(nodeLine string, lenString int) (NodeRow, error) {
	rec, err := parseCSVLine(nodeLine)
	if err != nil {
		return NodeRow{}, fmt.Errorf("parse node csv line: %w", err)
	}

	if len(rec) < lenString {
		return NodeRow{}, fmt.Errorf("invalid node row: expected at least %d fields, got %d", lenString, len(rec))
	}

	numPorts, err := strconv.Atoi(rec[1])
	if err != nil {
		return NodeRow{}, fmt.Errorf("invalid NumPorts %q: %w", rec[1], err)
	}

	nodeType, err := strconv.Atoi(rec[2])
	if err != nil {
		return NodeRow{}, fmt.Errorf("invalid NodeType %q: %w", rec[2], err)
	}

	classVersion, err := strconv.Atoi(rec[3])
	if err != nil {
		return NodeRow{}, fmt.Errorf("invalid ClassVersion %q: %w", rec[3], err)
	}

	baseVersion, err := strconv.Atoi(rec[4])
	if err != nil {
		return NodeRow{}, fmt.Errorf("invalid BaseVersion %q: %w", rec[4], err)
	}

	return NodeRow{
		NodeDesc:        rec[0],
		NumPorts:        numPorts,
		NodeType:        nodeType,
		ClassVersion:    classVersion,
		BaseVersion:     baseVersion,
		SystemImageGUID: rec[5],
		NodeGUID:        rec[6],
		PortGUID:        rec[7],
	}, nil
}

func parsePorts(portLine string, lenString int) (PortRow, error) {
	rec, err := parseCSVLine(portLine)
	if err != nil {
		return PortRow{}, fmt.Errorf("parse port csv line: %w", err)
	}

	if len(rec) < lenString {
		return PortRow{}, fmt.Errorf("invalid port row: expected at least %d fields, got %d", lenString, len(rec))
	}

	portNum, err := strconv.Atoi(rec[2])
	if err != nil {
		return PortRow{}, fmt.Errorf("invalid PortNum %q: %w", rec[2], err)
	}

	lid, err := strconv.Atoi(rec[6])
	if err != nil {
		return PortRow{}, fmt.Errorf("invalid LID %q: %w", rec[6], err)
	}

	linkWidthActv, err := strconv.Atoi(rec[10])
	if err != nil {
		return PortRow{}, fmt.Errorf("invalid LinkWidthActv %q: %w", rec[10], err)
	}

	linkSpeedActv, err := strconv.Atoi(rec[15])
	if err != nil {
		return PortRow{}, fmt.Errorf("invalid LinkSpeedActv %q: %w", rec[15], err)
	}

	portPhyState, err := strconv.Atoi(rec[19])
	if err != nil {
		return PortRow{}, fmt.Errorf("invalid PortPhyState %q: %w", rec[19], err)
	}

	portState, err := strconv.Atoi(rec[20])
	if err != nil {
		return PortRow{}, fmt.Errorf("invalid PortState %q: %w", rec[20], err)
	}

	return PortRow{
		NodeGUID:      rec[0],
		PortGUID:      rec[1],
		PortNum:       portNum,
		LID:           lid,
		PortPhyState:  portPhyState,
		PortState:     portState,
		LinkWidthActv: linkWidthActv,
		LinkSpeedActv: linkSpeedActv,
	}, nil
}

func parseSwitches(switchLine string, lenString int) (SwitchRow, error) {
	rec, err := parseCSVLine(switchLine)
	if err != nil {
		return SwitchRow{}, fmt.Errorf("parse switch csv line: %w", err)
	}

	if len(rec) < lenString {
		return SwitchRow{}, fmt.Errorf("invalid switch row: expected at least %d fields, got %d", lenString, len(rec))
	}

	var numParams []int
	for i := 1; i < lenString; i++ {
		numParam, err := strconv.Atoi(rec[i])
		if err != nil {
			return SwitchRow{}, fmt.Errorf("invalid numeric parameter: %q", rec[i])
		}
		numParams = append(numParams, numParam)
	}

	return SwitchRow{
		NodeGUID:             rec[0],
		LinearFDBCap:         numParams[0],
		RandomFDBCap:         numParams[1],
		MCastFDBCap:          numParams[2],
		LinearFDBTop:         numParams[3],
		DefPort:              numParams[4],
		DefMCastPriPort:      numParams[5],
		DefMCastNotPriPort:   numParams[6],
		LifeTimeValue:        numParams[7],
		PortStateChange:      numParams[8],
		OptimizedSLVLMapping: numParams[9],
		LidsPerPort:          numParams[10],
		PartEnfCap:           numParams[11],
		InbEnfCap:            numParams[12],
		OutbEnfCap:           numParams[13],
		FilterRawInbCap:      numParams[14],
		FilterRawOutbCap:     numParams[15],
		ENP0:                 numParams[16],
		MCastFDBTop:          numParams[17],
	}, nil
}

func parseNodeInfo(nodeInfoLine string, lenString int) (NodeInfoRow, error) {
	rec, err := parseCSVLine(nodeInfoLine)
	if err != nil {
		return NodeInfoRow{}, fmt.Errorf("parse node info csv line: %w", err)
	}

	if len(rec) < lenString {
		return NodeInfoRow{}, fmt.Errorf("invalid node info row: expected at least %d fields, got %d", lenString, len(rec))
	}

	return NodeInfoRow{
		NodeGUID:     rec[0],
		SerialNumber: rec[1],
		PartNumber:   rec[2],
		Revision:     rec[3],
		ProductName:  rec[4],
	}, nil
}

func parseCSVLine(line string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(line))
	r.Comma = ','
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1
	return r.Read()
}
