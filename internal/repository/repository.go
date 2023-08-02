package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
	"wb_l0/common"
	"wb_l0/config"
	"wb_l0/pkg/logger"
	"wb_l0/pkg/postgres"
)

const streamName = "my_stream"

type Repository struct {
	pool  *pgxpool.Pool
	cache map[int][]byte
}

func New(ctx context.Context, config config.PsqlStorage) *Repository {
	log := logger.GetLogger()
	pool, err := postgres.NewClient(ctx, config, 3)
	if err != nil {
		log.Fatal().Err(err).Msg("Can`t create psql client")
	}
	cache := make(map[int][]byte)
	return &Repository{pool: pool, cache: cache}
}

func (r *Repository) Get(id int) ([]byte, error) {
	bytes, ok := r.cache[id]
	if !ok {
		err := r.pool.QueryRow(context.Background(), "SELECT model FROM wb_table WHERE id=$1", id).Scan(&bytes)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		if err != nil {
			return nil, err
		}
	}
	return bytes, nil
}

func (r *Repository) CacheRecovery() error {
	rows, err := r.pool.Query(context.Background(), "SELECT id, model FROM wb_table")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var b []byte
		rows.Scan(&id, &b)
		r.cache[id] = b
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) ConsumeMessages() error {
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
		id := rand.Int()
		id %= 10000000 //maybe change
		log.Print(id, " smth")
		_, err := r.pool.Exec(ctx, "INSERT INTO wb_table VALUES ($1, $2)", id, msg.Data())
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}
	return nil
}

func Shutdown() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()
	js, _ := jetstream.New(nc)
	ctx := context.Background()
	js.DeleteStream(ctx, streamName)
}
