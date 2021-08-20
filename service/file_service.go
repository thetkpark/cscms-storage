package service

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
)

const MaxImageSize = 100 << 20

type FileService struct {
	*pb.UnimplementedFileServiceServer
	log hclog.Logger
}

func NewFileService(log hclog.Logger) *FileService {
	return &FileService{log: log}
}

func (s *FileService) UploadFile(stream pb.FileService_UploadFileServer) error {
	// Receive file info
	req, err := stream.Recv()
	if err != nil {
		return printAndReturnError(s.log, err, codes.Unknown, "unable to receive file info")
	}
	fileType := req.GetInfo().GetFileType()
	s.log.Info(fileType)

	// Create bytes buffer for incoming data
	fileData := bytes.Buffer{}
	fileSize := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.log.Info("no more data")
			break
		} else if err != nil {
			return printAndReturnError(s.log, err, codes.Unknown, "error receiving data")
		}

		chuck := req.GetFileData()
		fileSize += len(chuck)

		// Check MaxImageSize
		if fileSize > MaxImageSize {
			return printAndReturnError(s.log, nil, codes.InvalidArgument, "file too large")
		}

		// Write data chuck to bytes buffer
		_, err = fileData.Write(chuck)
		if err != nil {
			return printAndReturnError(s.log, err, codes.Internal, "cannot write data to fileData")
		}
	}

	// Create random filename
	fileName, err := uuid.NewRandom()
	if err != nil {
		return printAndReturnError(s.log, err, codes.Internal, "cannot generate uuid")
	}

	// Create new file on disk
	filePath := fmt.Sprintf("%s/%s", "files", fileName.String())
	file, err := os.Create(filePath)
	if err != nil {
		return printAndReturnError(s.log, err, codes.Internal, "cannot create new file on disk")
	}

	// Write fileData to created file
	_, err = fileData.WriteTo(file)
	if err != nil {
		return printAndReturnError(s.log, err, codes.Internal, "cannot write file data to disk")
	}

	// Create response and return
	res := &pb.UploadFileResponse{AccessCode: fileName.String()}
	err = stream.SendAndClose(res)
	if err != nil {
		return printAndReturnError(s.log, err, codes.Unknown, "cannot send response")
	}

	s.log.Debug("successfully write file", fileName, fileSize)
	return nil
}

func printAndReturnError(log hclog.Logger, err error, code codes.Code, msg string) error {
	log.Error(msg, err)
	return status.Errorf(code, msg)
}