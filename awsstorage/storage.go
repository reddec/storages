package awsstorage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"reddec/storages"
	"strings"
)

func New(bucket string, config *aws.Config) (storages.Storage, error) {
	s, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	return &storage{
		bucket:     bucket,
		client:     s3.New(s),
		downloader: s3manager.NewDownloader(s),
		uploader:   s3manager.NewUploader(s),
	}, nil
}

type storage struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	client     *s3.S3
	bucket     string
}

func (s *storage) Put(key []byte, data []byte) error {
	sKey := string(key)
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Body:   bytes.NewBuffer(data),
		Bucket: &s.bucket,
		Key:    &sKey,
	})
	return err
}

func (s *storage) Close() error {
	return nil
}

func (s *storage) Get(key []byte) ([]byte, error) {
	sKey := string(key)
	buffer := &aws.WriteAtBuffer{}
	_, err := s.downloader.Download(buffer, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &sKey,
	})
	if err != nil && strings.Contains(err.Error(), s3.ErrCodeNoSuchKey) {
		return nil, os.ErrNotExist
	}
	return buffer.Bytes(), err
}

func (s *storage) Del(key []byte) error {
	sKey := string(key)
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &sKey,
	})
	if err != nil && strings.Contains(err.Error(), s3.ErrCodeNoSuchKey) {
		return nil
	}
	return err
}

func (s *storage) Keys(handler func(key []byte) error) error {
	var err error
	reqErr := s.client.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &s.bucket,
	}, func(items *s3.ListObjectsOutput, lastPage bool) bool {
		for _, item := range items.Contents {
			if item.Key != nil {
				err = handler([]byte(*item.Key))
			}
		}
		return err == nil
	})
	if reqErr != nil {
		return reqErr
	}
	return err
}
