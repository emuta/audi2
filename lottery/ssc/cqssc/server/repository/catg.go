package repository

import (
	"context"
	"strconv"

	"audi/lottery/ssc/cqssc/server/model"
)

func (s *Repository) GetCatg(ctx context.Context, id int32) (*model.Catg, error) {
	catg := model.Catg{Id: id}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&catg).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &catg, nil
}

func (s *Repository) FindCatg(ctx context.Context, args map[string]interface{}) (*[]model.Catg, error) {
	var catgs []model.Catg
	ch := make(chan error)
	go func() {
		defer close(ch)
		tx := s.db.Model(&model.Catg{})

		if id, ok := args["id"]; ok {
			tx = tx.Where("id = ?", id)
		}

		if name, ok := args["name"]; ok {
			tx = tx.Where("name = ?", name)
		}

		if tag, ok := args["tag"]; ok {
			tx = tx.Where("tag = ?", tag)
		}

		if pref, ok := args["pref"]; ok {
			prefStr, _ := pref.(string)
			if v, err := strconv.ParseBool(prefStr); err == nil {
				tx = tx.Where("pref = ?", v)
			}
		}
		ch <- tx.Order("id ASC").Find(&catgs).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &catgs, nil
}
