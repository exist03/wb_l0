package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"wb_l0/common"
	"wb_l0/config"
	"wb_l0/pkg/logger"
	"wb_l0/pkg/postgres"
)

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

func (r *Repository) Create(id uint32, data []byte) error {
	_, err := r.pool.Exec(context.Background(), "INSERT INTO wb_table VALUES ($1, $2)", id, data)
	if err != nil {
		return err
	}
	return nil
}
