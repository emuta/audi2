package repository

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"audi/lottery/ssc/cqssc/server/model"
)

func (s *Repository) CreateTerm(ctx context.Context, id int64, f, t time.Time) (*model.Term, error) {
	p := model.Term{Id: id, StartFrom: f, EndTo: t}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Create(&p).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &p, nil
}

func (s *Repository) GetTerm(ctx context.Context, id int64) (*model.Term, error) {
	p := model.Term{Id: id}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&p).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &p, nil
}

func (s *Repository) SetTermResult(ctx context.Context, id int64, codes []string) (*model.Term, error) {
	term := model.Term{Id: id}

	tx := s.db.Begin().Set("gorm:query_option", "for update")
	ch := make(chan error)
	go func() {
		defer close(ch)

		if err := tx.Where("codes is null").Take(&term).Error; err != nil {
			ch <- err
			return
		}

		ch <- tx.Model(&term).Updates(map[string]interface{}{"codes": codes, "opened_at": time.Now()}).Error
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

	return &term, tx.Commit().Error
}

func (s *Repository) UpdateResult(ctx context.Context, id int64, codes []string) (*model.Term, error) {
	term := model.Term{Id: id}

	tx := s.db.Begin().Set("gorm:query_option", "for update")
	ch := make(chan error)
	go func() {
		defer close(ch)

		if err := tx.Take(&term).Error; err != nil {
			ch <- err
			return
		}

		if err := tx.Model(&term).Updates(map[string]interface{}{"codes": codes}).Error; err != nil {
			log.WithError(err).Error("Failed to updated term codes")
			ch <- err
			return
		}

		log.WithFields(log.Fields{
			"id":    id,
			"codes": codes,
		}).Info("Term codes updated success")

		ch <- nil
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

	return &term, tx.Commit().Error
}

func getFindTermDB(db *gorm.DB, args map[string]interface{}) *gorm.DB {
	tx := db.Model(&model.Term{})

	if id, ok := args["id"]; ok {
		tx = tx.Where("id = ?", id)
	}

	if code, ok := args["code"]; ok {
		tx = tx.Where("? = any(codes)", code)
	}

	if startFrom, ok := args["start_from"]; ok {
		tx = tx.Where("start_from >= ?", startFrom)
	}

	if endTo, ok := args["end_to"]; ok {
		tx = tx.Where("end_to <= ?", endTo)
	}

	return tx
}

func (s *Repository) FindTerm(ctx context.Context, args map[string]interface{}) (*[]model.Term, error) {
	var results []model.Term
	ch := make(chan error)
	go func() {
		defer close(ch)

		tx := getFindTermDB(s.db, args)

		if limit, ok := args["limit"]; ok {
			tx = tx.Limit(limit)
		}

		if offset, ok := args["offset"]; ok {
			tx = tx.Offset(offset)
		}

		if orderBy, ok := args["order_by"]; ok {
			tx = tx.Order(orderBy)
		} else {
			tx = tx.Order("id DESC")
		}

		ch <- tx.Find(&results).Error

	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &results, nil
}

func (s *Repository) CountFindTerm(ctx context.Context, args map[string]interface{}) (*int32, error) {
	var total int32
	ch := make(chan error)
	go func() {
		defer close(ch)

		tx := getFindTermDB(s.db, args)

		ch <- tx.Count(&total).Error

	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &total, nil
}

func (s *Repository) SettleTerm(ctx context.Context, id int64) (*model.Term, error) {
	term := model.Term{Id: id}

	tx := s.db.Begin().Set("gorm:query_option", "for update")
	ch := make(chan error)
	go func() {
		defer close(ch)

		if err := tx.Take(&term).Error; err != nil {
			ch <- err
			return
		}

		ch <- tx.Model(&term).Updates(map[string]interface{}{"settled_at": time.Now()}).Error
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

	return &term, tx.Commit().Error
}

func (s *Repository) RevokeTerm(ctx context.Context, id int64) (*model.Term, error) {
	term := model.Term{Id: id}

	tx := s.db.Begin().Set("gorm:query_option", "for update")
	ch := make(chan error)
	go func() {
		defer close(ch)

		if err := tx.Take(&term).Error; err != nil {
			ch <- err
			return
		}

		ch <- tx.Model(&term).Updates(map[string]interface{}{"revoked_at": time.Now()}).Error
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

	return &term, tx.Commit().Error
}
