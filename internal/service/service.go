package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"math/rand"
	"strconv"
	"time"
	"wb_l0/common"
	"wb_l0/internal"
	"wb_l0/pkg/logger"
)

const streamName = "my_stream"

type repository interface {
	Get(id int) ([]byte, error)
	Create(id uint32, data []byte) error
}
type Service struct {
	repository
}

func New(repository repository) *Service {
	return &Service{repository}
}

func (s *Service) Get(id string) ([]byte, error) {
	log := logger.GetLogger()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, common.ErrInvalidID
	}
	log.Debug().Msg(fmt.Sprintf("id == %d", idInt))
	bytes, err := s.repository.Get(idInt)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (s *Service) ConsumeMessages() error {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()
	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}
	ctx := context.Background()

	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"subject"},
	})
	if err != nil {
		return err
	}
	cons, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		InactiveThreshold: 5 * time.Minute,
		FilterSubject:     "subject",
	})
	fetchResult, _ := cons.Fetch(40)
	for msg := range fetchResult.Messages() {
		err := Validate(msg.Data())
		if err != nil {
			continue
		}
		id := rand.Uint32()
		err = s.repository.Create(id, msg.Data())
		if err != nil {
			return err
		}
	}
	return nil
}

func Validate(data []byte) error {
	m := &internal.Model{}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	err := dec.Decode(m)
	if err != nil {
		return err
	}
	return nil
}

func ShutdownStream() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()
	js, _ := jetstream.New(nc)
	ctx := context.Background()
	js.DeleteStream(ctx, streamName)
}
