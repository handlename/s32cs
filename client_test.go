package s32cs_test

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fujiwara/s32cs"
)

var (
	sess     = session.New()
	endpoint = os.Getenv("CS_ENDPOINT")
)

func init() {
	s32cs.DEBUG = true
}

func TestUpload(t *testing.T) {
	f, err := os.Open(os.Getenv("TEST_SDF"))
	if err != nil {
		t.Skip("test sdf file is not specified by TEST_SDF env")
		return
	}
	client := s32cs.NewClient(sess, endpoint, nil)
	if err := client.Upload(f, ""); err != nil {
		t.Error(err)
	}
}

func TestProcess(t *testing.T) {
	var event s32cs.S3Event
	err := json.Unmarshal([]byte(os.Getenv("TEST_EVENT_JSON")), &event)
	if err != nil {
		t.Skip("TEST_EVENT_JSON invalid", err)
	}
	client := s32cs.NewClient(sess, endpoint, nil)
	if err := client.Process(event); err != nil {
		t.Error(err)
	}
}

func TestProcessRegexp(t *testing.T) {
	var event s32cs.S3Event
	err := json.Unmarshal([]byte(os.Getenv("TEST_EVENT_JSON")), &event)
	if err != nil {
		t.Skip("TEST_EVENT_JSON invalid", err)
	}
	client := s32cs.NewClient(
		sess,
		endpoint,
		regexp.MustCompile("^test/(.+?)/"),
	)
	if err := client.Process(event); err != nil {
		t.Error(err)
	}
}