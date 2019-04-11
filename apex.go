package s32cs

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"

	apex "github.com/apex/go-apex"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ApexRun() {
	var reg *regexp.Regexp
	if r := os.Getenv("KEY_REGEXP"); r != "" {
		reg = regexp.MustCompile(r)
	} else {
		reg = nil
	}
	client := NewClient(session.New(), os.Getenv("ENDPOINT"), reg)

	handler := func(msg json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var event SQSEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			return nil, err
		}
		if event.QueueURL != "" {
			log.Printf("starting process sqs %s", event.QueueURL)
			return nil, client.ProcessSQS(event.QueueURL)
		}
		var s3event S3Event
		if err := json.Unmarshal(msg, &s3event); err != nil {
			return nil, err
		}
		log.Println("starting process s3 event:", s3event.String())
		if err := client.Process(s3event); err != nil {
			return nil, err
		}
		return true, nil
	}

	env := os.Getenv("AWS_EXECUTION_ENV")
	if strings.HasPrefix(env, "AWS_Lambda_nodejs") {
		// Apex node runtime (v0.x)
		apex.HandleFunc(handler)
	} else if strings.HasPrefix(env, "AWS_Lambda_go") {
		// Go native runtime
		lambda.Start(handler)
	} else {
		log.Printf("[error] Lambda execution environment %s is not supported", env)
	}
}
