package jsonrpc

import (
	"bridge-allowance/config"
	"bridge-allowance/utils"
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
)

type RPCRequest struct {
	Method  string      `json:"method"`
	Jsonrpc string      `json:"jsonrpc"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type RPCHandler struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
}

func NewJsonRPCHandler(cfg *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *RPCHandler {
	return &RPCHandler{
		env:         cfg,
		logger:      logger,
		httpRequest: httpRequest,
	}
}

func (r *RPCHandler) CallJsonRPCMethod(request RPCRequest, Url string) ([]byte, error) {
	data, err := json.Marshal(request)
	if err != nil {
		r.logger.Errorf("Marshal: %v", err)
		return nil, err
	}
	respBody, err := r.httpRequest.PostRequest(Url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
