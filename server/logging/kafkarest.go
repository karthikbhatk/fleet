package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	CONTENT_TYPE_MIME = "application/vnd.kafka.json.v1+json"
	CONTENT_TYPE      = "Content-Type"
	TIMESTAMP         = "TimeStamp"
	URL_PUBLISH_TOPIC = "%s/topics/%s"
)

type KafkaRESTParams struct {
	KafkaProxyHost string
	KafkaTopic     string
	KafkaTimeout   int
}

type kafkaRESTProducer struct {
	client *http.Client
	URL    string
}

type kafkaRecords struct {
	Records []kafkaValue `json:"records"`
}

type kafkaValue struct {
	Value json.RawMessage `json:"value"`
}

func NewKafkaRESTWriter(p *KafkaRESTParams) (*kafkaRESTProducer, error) {
	return &kafkaRESTProducer{
		URL: fmt.Sprintf(URL_PUBLISH_TOPIC, p.KafkaProxyHost, p.KafkaTopic),
		client: &http.Client{
			Timeout: time.Duration(p.KafkaTimeout) * time.Second,
		},
	}, nil
}

func (l *kafkaRESTProducer) Write(ctx context.Context, logs []json.RawMessage) error {
	data := kafkaRecords{
		Records: make([]kafkaValue, len(logs)),
	}

	for i, log := range logs {
		data.Records[i] = kafkaValue{
			Value: log,
		}
	}

	output, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "kafka rest marshal")
	}

	resp, err := l.post(l.URL, bytes.NewBuffer(output))
	if err != nil {
		return errors.Wrap(err, "kafka rest post")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("Error posting to %s: %d. %s", l.URL, resp.StatusCode, string(body))
	}

	return nil
}

func (l *kafkaRESTProducer) post(url string, buf *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, errors.Wrap(err, "kafka rest new request")
	}

	now := float64(time.Now().UnixNano()) / float64(time.Second)
	req.Header.Set(CONTENT_TYPE, CONTENT_TYPE_MIME)
	req.Header.Set(TIMESTAMP, fmt.Sprintf("%f", now))

	return l.client.Do(req)
}
