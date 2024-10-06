package handlers

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type CloudProvider struct {
	Name  string
	Mutex sync.Mutex
}

// HandlePhotoUpload handles the photo upload request.
// It reads the cloudProvider query parameter and the request body.
func HandlePhotoUpload(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	endpoint := "localhost:9000"
	accessKeyID := "4wg1w8OwGoBaKCHlpkxh"
	secretAccessKey := "jHB11ESetLV5ipNo2ZqLeHHrpm1ogQqW4tlumfLQ"
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	AWSClientSession, err := session.NewSession(&aws.Config{
		Region:           aws.String("eu-west-2"),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		DisableSSL:       aws.Bool(true),
	})
	if err != nil {
		logger.Error("Error connecting to minio as AWS", slog.Any("error", err))
		return
	}

	AWSClient := s3.New(AWSClientSession)

	queryParams := r.URL.Query()
	cloudProviderName := queryParams.Get("cloudProvider")
	bucketName := queryParams.Get("bucketName")
	fileName := queryParams.Get("fileName")

	logger.Info("Received request with cloudProvider '" + cloudProviderName + "'.")

	approvedProviders := []*CloudProvider{
		{Name: "aws"},
		{Name: "azure"},
		{Name: "gcp"},
	}

	cloudProvider := getCloudProvider(cloudProviderName, approvedProviders)

	if cloudProvider == nil {

		logger.Error("cloudProvider not recognised", slog.String("error", "cloudProvider '"+cloudProviderName+"' not recognised"), slog.Int("status", http.StatusBadRequest))
		return
	}

	cloudProvider.Mutex.Lock()
	defer cloudProvider.Mutex.Unlock()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		logger.Error("reading request body", slog.Int("status", http.StatusInternalServerError))
		return
	}
	defer r.Body.Close()

	fileSize := len(body)

	logger.Info("Received request of body size " + strconv.Itoa(fileSize))

	contentType := http.DetectContentType(body)

	if contentType != "image/png" {

		logger.Error("File format '"+contentType+"' is not valid.", slog.Int("status", http.StatusBadRequest), slog.String("fileFormat", contentType))
		return
	}
	info, err := AWSClient.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logger.Error("Error connecting to AWS (minio)", slog.String("fileName", fileName), slog.Any("minioError", err))
		return
	}
	logger.Info("Upload status", slog.Any("info", info))
	w.WriteHeader(http.StatusOK)

}

func getCloudProvider(cloudProvider string, approvedProviders []*CloudProvider) *CloudProvider {
	for _, p := range approvedProviders {
		if p.Name == cloudProvider {
			return p
		}
	}
	return nil
}
