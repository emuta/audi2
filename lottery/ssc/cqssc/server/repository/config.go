package repository

import (
	"context"

	"audi/lottery/ssc/cqssc/server/model"
)

func (s *Repository) GetConfig(ctx context.Context) (*model.Config, error) {
	var cfg model.Config
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&cfg).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func (s *Repository) UpdateConfig(ctx context.Context, odds, comm float64) (*model.Config, error) {
	var cfg model.Config
	tx := s.db.Begin().Set("gorm:query_option", "for update")
	ch := make(chan error)
	go func() {
		defer close(ch)

		if err := s.db.Take(&cfg).Error; err != nil {
			ch <- err
			return
		}

		args := map[string]float64{
			"odds": odds,
			"comm": comm,
		}

		ch <- tx.Model(&cfg).Updates(args).Error
	}()

	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return &cfg, tx.Commit().Error
}
