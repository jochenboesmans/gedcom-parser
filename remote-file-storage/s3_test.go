package remote_file_storage

import (
	"io/ioutil"
	"testing"
)

func TestS3Read(t *testing.T) {
	_, err := S3Read("ITIS.ged")
	if err != nil {
		t.Errorf("failed to read from s3: %s", err)
	}
}

func TestS3Write(t *testing.T) {
	f, err := ioutil.ReadFile("../io/ITIS.ged")
	if err != nil {
		t.Errorf("failed to open file: %s", err)
	}

	_, err = S3Write("ITIS.ged", f)
	if err != nil {
		t.Errorf("failed to write to s3: %s", err)
	}
}
