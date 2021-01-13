//+build s3

package remote_file_storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"testing"
)

type server struct {
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
}

var s server

func (s *server) Init() {
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(os.Getenv("AWS_REGION"))))

	s.uploader = s3manager.NewUploader(sess)
	s.downloader = s3manager.NewDownloader(sess)
}

func setup() {
	_ = godotenv.Load("../.env")
	s.Init()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestS3Read(t *testing.T) {
	_, err := S3Read("gedcom/ITIS.ged", s.downloader)
	if err != nil {
		t.Errorf("failed to read from s3: %s", err)
	}
}

func TestS3Write(t *testing.T) {
	f, err := ioutil.ReadFile("../io/ITIS.ged")
	if err != nil {
		t.Errorf("failed to open file: %s", err)
	}

	_, err = S3Write("gedcom/ITIS.ged", &f, s.uploader)
	if err != nil {
		t.Errorf("failed to write to s3: %s", err)
	}
}
