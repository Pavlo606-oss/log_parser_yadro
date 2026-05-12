package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"repo/internal/service"
	"strconv"
)

type Handler struct {
	s *service.Service
}

func (h *Handler) PostLog(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req PostLogReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "decode request body", http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}

	file, err := os.Open(req.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, os.ErrPermission) {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
		http.Error(w, "failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	result, err := h.s.ImportLog(r.Context(), file, req.Path)
	if err != nil {
		http.Error(w, "failed to import log", http.StatusInternalServerError)
		return
	}

	if err = writeJSON(w, http.StatusCreated, PostLogResp{LogID: result}); err != nil {
		return
	}
}

func (h *Handler) GetNodeTopology(w http.ResponseWriter, r *http.Request) {
	logIDStr := r.PathValue("log_id")
	logID, err := strconv.ParseInt(logIDStr, 10, 64)
	if err != nil || logID <= 0 {
		http.Error(w, "incorrect log_id", http.StatusBadRequest)
		return
	}

	result, err := h.s.GetNodes(r.Context(), logID)
	if err != nil {
		http.Error(w, "failed to get nodes", http.StatusInternalServerError)
		return
	}

	if err = writeJSON(w, http.StatusOK, result); err != nil {
		return
	}
}

func (h *Handler) GetNodeDetail(w http.ResponseWriter, r *http.Request) {
	nodeIDStr := r.PathValue("node_id")
	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)

	if err != nil || nodeID <= 0 {
		http.Error(w, "incorrect node_id", http.StatusBadRequest)
		return
	}

	result, err := h.s.GetNodeDetail(r.Context(), nodeID)
	if err != nil {
		http.Error(w, "failed to get node detail", http.StatusInternalServerError)
		return
	}

	if err = writeJSON(w, http.StatusOK, result); err != nil {
		return
	}
}

func (h *Handler) GetPorts(w http.ResponseWriter, r *http.Request) {
	nodeIDStr := r.PathValue("node_id")
	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)

	if err != nil || nodeID <= 0 {
		http.Error(w, "incorrect node_id", http.StatusBadRequest)
		return
	}

	result, err := h.s.GetPortsByNodeId(r.Context(), nodeID)
	if err != nil {
		http.Error(w, "failed to get ports", http.StatusInternalServerError)
		return
	}

	if err = writeJSON(w, http.StatusOK, result); err != nil {
		return
	}
}

func (h *Handler) GetLogMeta(w http.ResponseWriter, r *http.Request) {
	logIDStr := r.PathValue("log_id")
	logID, err := strconv.ParseInt(logIDStr, 10, 64)
	if err != nil || logID <= 0 {
		http.Error(w, "incorrect log_id", http.StatusBadRequest)
		return
	}

	result, err := h.s.GetLogMetaByID(r.Context(), logID)
	if err != nil {
		http.Error(w, "failed to get meta log", http.StatusInternalServerError)
		return
	}

	if err = writeJSON(w, http.StatusOK, result); err != nil {
		return
	}
}
