package remote_file_storage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

func S3Read(pathToFile string, downloader *s3manager.Downloader) (*[]byte, error) {
	buffer := aws.WriteAtBuffer{}
	_, err := downloader.Download(&buffer, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(pathToFile),
	})
	if err != nil {
		return nil, err
	}

	readBytes := buffer.Bytes()
	return &readBytes, nil
}

func S3Write(pathToFile string, content *[]byte, uploader *s3manager.Uploader) (*s3manager.UploadOutput, error) {
	r := bytes.NewReader(*content)

	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(pathToFile),
		Body:   r,
	})
}
