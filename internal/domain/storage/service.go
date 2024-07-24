package storage

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Service struct {
	MinioClient *minio.Client
	BucketName  string
}

func NewService() *Service {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	useSSL := true

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		MinioClient: minioClient,
		BucketName:  bucketName,
	}
}

func (s *Service) Find(objectName string) (*minio.Object, error) {
	ctx := context.Background()
	object, err := s.MinioClient.GetObject(ctx, s.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (s *Service) Upload(filename string, reader io.Reader, contentType string) (*minio.UploadInfo, error) {
	ctx := context.Background()

	result, err := s.MinioClient.PutObject(ctx, s.BucketName, filename, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (s *Service) Delete(filename string) error {
	ctx := context.Background()

	err := s.MinioClient.RemoveObject(ctx, s.BucketName, filename, minio.RemoveObjectOptions{
		GovernanceBypass: true,
	})
	if err != nil {
		return err
	}

	return nil
}
