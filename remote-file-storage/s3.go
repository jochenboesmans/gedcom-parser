package remote_file_storage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	session2 "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func S3Read(pathToFile string) (*[]byte, error) {
	session := session2.Must(session2.NewSession(aws.NewConfig().WithRegion("eu-west-2")))

	downloader := s3manager.NewDownloader(session)

	buffer := aws.WriteAtBuffer{}
	_, err := downloader.Download(&buffer, &s3.GetObjectInput{
		Bucket: aws.String("jochen-gedcom"),
		Key:    aws.String(pathToFile),
	})
	if err != nil {
		return nil, err
	}

	readBytes := buffer.Bytes()
	return &readBytes, nil
}

func S3Write(pathToFile string, content *[]byte) (*s3manager.UploadOutput, error) {
	session := session2.Must(session2.NewSession(aws.NewConfig().WithRegion("eu-west-2")))

	uploader := s3manager.NewUploader(session)

	r := bytes.NewReader(*content)

	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("jochen-gedcom"),
		Key:    aws.String(pathToFile),
		Body:   r,
	})
}
