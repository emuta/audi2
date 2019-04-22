package repository

import (
	"context"

	"audi/lottery/ssc/cqssc/server/model"
)

func (s *Repository) GetGroup(ctx context.Context, id int32) (*model.Group, error) {
	group := model.Group{Id: id}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&group).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &group, nil
}

func (s *Repository) FindGroup(ctx context.Context, args map[string]interface{}) (*[]model.Group, error) {
	var groups []model.Group
	ch := make(chan error)
	go func() {
		defer close(ch)
		tx := s.db.Model(&model.Group{})

		if id, ok := args["id"]; ok {
			tx = tx.Where("id = ?", id)
		}

		if name, ok := args["name"]; ok {
			tx = tx.Where("name = ?", name)
		}

		if tag, ok := args["tag"]; ok {
			tx = tx.Where("tag = ?", tag)
		}
		ch <- tx.Order("id ASC").Find(&groups).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &groups, nil
}
