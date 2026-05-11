package service

import (
	"context"
	"fmt"
	"io"
	"repo/internal/parser"
	"repo/internal/repository"
	"strings"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImportLog(ctx context.Context, r io.Reader, filename string) (int64, error) {
	var result int64

	switch {
	case strings.HasSuffix(filename, ".db_csv"):
		parseLog, err := parser.ParseDBCSV(r)
		if err != nil {
			return 0, fmt.Errorf("parse file=%s: %w", filename, err)
		}

		tx, err := s.repo.BeginTx(ctx)
		if err != nil {
			return 0, fmt.Errorf("begin tx: %w", err)
		}
		defer tx.Rollback()

		result, err = s.repo.CreateLog(ctx, tx, filename, ".db_csv")
		if err != nil {
			return 0, fmt.Errorf("create log: %w", err)
		}

		m, err := s.repo.CreateNodes(ctx, tx, result, parseLog.Nodes)
		if err != nil {
			return 0, fmt.Errorf("create nodes: %w", err)
		}

		err = s.repo.CreatePorts(ctx, tx, result, m, parseLog.Ports)
		if err != nil {
			return 0, fmt.Errorf("create ports: %w", err)
		}

		err = s.repo.CreateNodeInfos(ctx, tx, result, m, parseLog.NodesInfo)
		if err != nil {
			return 0, fmt.Errorf("create nodes_info: %w", err)
		}

		err = s.repo.CreateSwitches(ctx, tx, result, m, parseLog.Switches)
		if err != nil {
			return 0, fmt.Errorf("create switches: %w", err)
		}

		if err = tx.Commit(); err != nil {
			return 0, fmt.Errorf("commit: %w", err)
		}
		return result, nil
	case strings.HasSuffix(filename, ".sharp_an_info"):
		parseLog, err := parser.ParseSharpANInfo(r)
		if err != nil {
			return 0, fmt.Errorf("parse file=%s: %w", filename, err)
		}

		tx, err := s.repo.BeginTx(ctx)
		if err != nil {
			return 0, fmt.Errorf("begin tx: %w", err)
		}
		defer tx.Rollback()

		result, err = s.repo.CreateLog(ctx, tx, filename, ".sharp_an_info")
		if err != nil {
			return 0, fmt.Errorf("create log: %w", err)
		}

		err = s.repo.CreateSharpANInfos(ctx, tx, result, parseLog)
		if err != nil {
			return 0, fmt.Errorf("create sharp an infos: %w", err)
		}

		if err = tx.Commit(); err != nil {
			return 0, fmt.Errorf("commit: %w", err)
		}
		return result, nil
	default:
		return 0, fmt.Errorf("unsupported file type, filename: %s", filename)
	}
}

func (s *Service) GetNodes(ctx context.Context, logID int64) ([]repository.Node, error) {
	result, err := s.repo.GetNodesByLogID(ctx, logID)
	if err != nil {
		return nil, fmt.Errorf("get nodes: %w", err)
	}

	return result, nil
}

func (s *Service) GetNodeDetail(ctx context.Context, nodeID int64) (*repository.NodeDetail, error) {
	result, err := s.repo.GetNodeDetailByID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("get node detail: %w", err)
	}

	return result, nil
}

func (s *Service) GetPortsByNodeId(ctx context.Context, nodeID int64) ([]repository.Port, error) {
	result, err := s.repo.GetPortsByNodeID(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("get ports by node_id=%d: %w", nodeID, err)
	}

	return result, nil
}

func (s *Service) GetLogMetaByID(ctx context.Context, logID int64) (*repository.LogMeta, error) {
	result, err := s.repo.GetLogMeta(ctx, logID)
	if err != nil {
		return nil, fmt.Errorf("get log meta by log_id=%d: %w", logID, err)
	}
	return result, nil
}
