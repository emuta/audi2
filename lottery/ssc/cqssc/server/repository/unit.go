package repository

import (
	"context"

	"audi/lottery/ssc/cqssc/server/model"
)

func (s *Repository) GetUnit(ctx context.Context, id int64) (*model.Unit, error) {
	unit := model.Unit{Id: id}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&unit).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &unit, nil
}

func (s *Repository) FindUnit(ctx context.Context, args map[string]interface{}) (*[]model.Unit, error) {
	var units []model.Unit
	ch := make(chan error)
	go func() {
		defer close(ch)

		tx := s.db.Model(&model.Unit{})
		if id, ok := args["id"]; ok {
			tx.Where("id = ?", id)
		}

		if name, ok := args["name"]; ok {
			tx.Where("name = ?", name)
		}

		if value, ok := args["value"]; ok {
			tx.Where("value = ?", value)
		}

		ch <- tx.Order("id ASC").Find(&units).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &units, nil
}
