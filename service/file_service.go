package service

import (
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/pb"
)

type FileService struct {
	*pb.UnimplementedFileServiceServer
	log hclog.Logger
}

func NewFileService(log hclog.Logger) *FileService {
	return &FileService{log: log}
}

func (s *FileService) UploadFile(stream pb.FileService_UploadFileServer) error {
	return nil
}