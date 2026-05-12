package handler

type PostLogReq struct {
	Path string `json:"path"`
}

type PostLogResp struct {
	LogID int64 `json:"log_id"`
}
