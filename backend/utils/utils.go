package utils

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var statusMap = map[int]string{
	200: "up",
	201: "up",
	301: "redirect",
	302: "redirect",
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	404: "not found",
	500: "server error",
	502: "bad gateway",
	503: "service unavailable",
	504: "gateway timeout",
}

func MapStatusCode(code int) string {
	if status, ok := statusMap[code]; ok {
		return status
	}
	return "unknown"
}

func CheckUsername(username string) bool {

	if len(username) < 3 || len(username) > 20 {
		return false
	}
	return true
}

func CheckURLFormat(url string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}

func CheckAlerts(alerts string) bool {
	if alerts != "y" && alerts != "n" && alerts != "yes" && alerts != "no" {
		return false
	}
	return true
}

func getAWSConfig() (map[string]string, error) {
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	if bucketName == "" {
		log.Println("⚠️  AWS S3 bucket name not set in environment variables")
		return nil, os.ErrInvalid
	}
	if region == "" || accessKeyID == "" || secretAccessKey == "" {
		log.Println("⚠️  AWS credentials or region not set in environment variables")
		return nil, os.ErrInvalid
	}
	return map[string]string{
		"region":          region,
		"accessKeyID":     accessKeyID,
		"secretAccessKey": secretAccessKey,
	}, nil
}

func CreateAWSSession() (*session.Session, error) {
	awsConfig, err := getAWSConfig()
	if err != nil {
		return nil, err
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsConfig["region"]),
		Credentials: credentials.NewStaticCredentials(awsConfig["accessKeyID"], awsConfig["secretAccessKey"], ""),
	})
	if err != nil {
		log.Printf("❌ Failed to create AWS session: %v", err)
		return nil, err
	}
	return sess, nil
}

func UploadFileToS3(key, bucketName, prefix string, image []byte) error {
	sess, err := CreateAWSSession()

	if err != nil {
		return err
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(prefix + key),
		Body:          bytes.NewReader(image),
		ContentType:   aws.String(http.DetectContentType(image)),
		ContentLength: aws.Int64(int64(len(image))),
	})
	if err != nil {
		log.Printf("❌ Failed to upload file to S3: %v", err)
		return err
	}

	return nil
}
