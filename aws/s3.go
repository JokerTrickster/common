package aws

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"mime/multipart"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
)

// ImgMeta stores metadata for different image types
type ImgMeta struct {
	Bucket     string
	Domain     string
	Path       string
	Width      int
	Height     int
	ExpireTime time.Duration
}

type ImgType uint8

const (
	ImgTypeFood     = ImgType(0)
	ImgTypeCategory = ImgType(1)
	ImgTypeProfile  = ImgType(2)
)

// Singleton metadata for images
var imgMeta = map[ImgType]ImgMeta{
	ImgTypeFood: {
		Bucket:     "dev-food-recommendation",
		Domain:     "dev-food-recommendation.s3.ap-northeast-2.amazonaws.com",
		Path:       "images",
		Width:      512,
		Height:     512,
		ExpireTime: 2 * time.Hour,
	},
	ImgTypeCategory: {
		Bucket:     "dev-food-recommendation",
		Domain:     "dev-food-recommendation.s3.ap-northeast-2.amazonaws.com",
		Path:       "category",
		Width:      62,
		Height:     62,
		ExpireTime: 2 * time.Hour,
	},
	ImgTypeProfile: {
		Bucket:     "dev-food-recommendation",
		Domain:     "dev-food-recommendation.s3.ap-northeast-2.amazonaws.com",
		Path:       "profiles",
		Width:      512,
		Height:     512,
		ExpireTime: 2 * time.Hour,
	},
}

type S3Service struct {
	service       *AWSService
	initOnce      sync.Once
	uploader      *manager.Uploader
	downloader    *manager.Downloader
	presignClient *s3.PresignClient
}

// GetS3Service creates a singleton instance of S3Service
func GetS3Service(region string) *S3Service {
	service := GetAWSService(region)
	return &S3Service{service: service}
}

// GetUploader initializes and returns the S3 uploader
func (s *S3Service) GetUploader() (*manager.Uploader, error) {
	var err error
	s.initOnce.Do(func() {
		uploader, e := s.service.GetUploader()
		if e != nil {
			err = e
			return
		}
		s.uploader = uploader
	})
	return s.uploader, err
}

// GetPresignClient initializes and returns the S3 presign client
func (s *S3Service) GetPresignClient() (*s3.PresignClient, error) {
	var err error
	s.initOnce.Do(func() {
		client, e := s.service.GetPresignClient()
		if e != nil {
			err = e
			return
		}
		s.presignClient = client
	})
	return s.presignClient, err
}

// UploadImage uploads an image to S3
func (s *S3Service) UploadImage(ctx context.Context, file *multipart.FileHeader, filename string, imgType ImgType) error {
	meta, ok := imgMeta[imgType]
	if !ok {
		return fmt.Errorf("invalid image type: %v", imgType)
	}

	imgBuffer, err := processImage(file, meta.Width, meta.Height, false)
	if err != nil {
		return err
	}

	uploader, err := s.GetUploader()
	if err != nil {
		return err
	}

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(meta.Bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s", meta.Path, filename)),
		Body:        imgBuffer,
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %v", err)
	}
	return nil
}

// GetSignedURL generates a presigned URL for an S3 object
func (s *S3Service) GetSignedURL(ctx context.Context, filename string, imgType ImgType) (string, error) {
	meta, ok := imgMeta[imgType]
	if !ok {
		return "", fmt.Errorf("invalid image type: %v", imgType)
	}

	client, err := s.GetPresignClient()
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s/%s", meta.Path, filename)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(meta.Bucket),
		Key:    aws.String(key),
	}

	result, err := client.PresignGetObject(ctx, presignParams, s3.WithPresignExpires(meta.ExpireTime))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}
	return result.URL, nil
}

// processImage handles image decoding, resizing, and encoding
func processImage(file *multipart.FileHeader, width, height int, crop bool) (*bytes.Buffer, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	if crop {
		img = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	} else {
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, img, imaging.PNG, imaging.PNGCompressionLevel(png.BestCompression)); err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}
	return buf, nil
}
