package handler

import (
	"bytes"
	"context"
	"encoding/json"

	pb "github.com/AleGaliev/runtimemetrics/internal/proto"
	"github.com/AleGaliev/runtimemetrics/internal/proto/convertor"
)

type UserServer struct {
	pb.UnimplementedMetricsServer
	storage Storage
}

func NewUserServer(storage Storage) *UserServer {
	return &UserServer{storage: storage}
}

func (s *UserServer) UpdateMetrics(ctx context.Context, req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	resp := pb.UpdateMetricsResponse{}

	metrics := convertor.ConvertProtoSliceToMetricsSimple(req.GetMetrics())

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}

	if err = s.storage.BatchUpdateMetrics(bytes.NewReader(jsonData)); err != nil {
		return nil, err
	}

	return &resp, nil
}
