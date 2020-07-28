package remote_file_storage

import (
	"github.com/aws/aws-sdk-go/aws"
	session2 "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

func S3Write(r io.Reader) (*s3manager.UploadOutput, error) {
	session := session2.Must(session2.NewSession())

	uploader := s3manager.NewUploader(session)

	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("myBucketString"),
		Key:    aws.String("myKeyString"),
		Body:   r,
	})
}
