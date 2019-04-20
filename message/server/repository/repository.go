package repository

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"audi/message/server/model"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (s *Repository) CreateMessage(ctx context.Context, srv string, event string, data []byte) (*model.Message, error) {
	msg := model.Message{
		Srv:    srv,
		Event:  event,
		Data:   postgres.Jsonb{data},
		Active: true,
	}
	
	ch := make(chan error)
	go func() {
		defer close(ch)

		err := s.db.Create(&msg).Error; 
		if err != nil {
			log.WithError(err).Error("Fail to create message")
		}

		ch <- err
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			return nil, err
		}
	}
	return &msg, nil
}

func (s *Repository) GetMessage(ctx context.Context, id int64) (*model.Message, error) {
	msg := model.Message{Id: id}
	ch := make(chan error)
	go func() {
		defer close(ch)

		if err := s.db.Take(&msg).Error; err != nil {
			log.WithError(err).WithField("id", id).Error("Fail to get message by id")
			ch <- err
			return
		}

		ch <- nil
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			return nil, err
		}
	}

	return &msg, nil
}

func (s *Repository) DeleteMessage(ctx context.Context, id int64) (error) {
	msg := model.Message{Id: id}
	tx := s.db.Begin()
	ch := make(chan error)
	go func() {
		defer close(ch)

		if err := tx.Model(&msg).Update("active", false).Error; err != nil {
			log.WithError(err).WithField("id", id).Error("Fail to SOFT DELETE message by id")
			ch <- err
			return
		}

		ch <- nil
	}()

	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	case err := <- ch:
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func getFindMessageDB(db *gorm.DB, args map[string]interface{}) (*gorm.DB) {
	tx := db.Model(&model.Message{})

	if id, ok := args["id"]; ok {
		tx = tx.Where("id = ?", id)
	}

	if srv, ok := args["srv"]; ok {
		tx = tx.Where("srv = ?", srv)
	}

	if event, ok := args["event"]; ok {
		tx = tx.Where("event = ?", event)
	}

	if createdFrom, ok := args["created_from"]; ok {
		tx = tx.Where("created_at >= ?", createdFrom)
	}

	if createdTo, ok := args["created_to"]; ok {
		tx = tx.Where("created_at >= ?", createdTo)
	}

	if limit, ok := args["limit"]; ok {
		tx = tx.Limit(limit)
	}

	if offset, ok := args["offset"]; ok {
		tx = tx.Offset(offset)
	}

	return tx
}

func (s *Repository) FindMessage(ctx context.Context, args map[string]interface{}) (*[]model.Message, error) {
	var msgs []model.Message
	ch := make(chan error)
	go func() {
		defer close(ch)

		tx := getFindMessageDB(s.db, args)
		// some filter at here
		err := tx.Find(&msgs).Error; 
		if err != nil {
			log.WithError(err).Error("Fail to find message")
		}
		ch <- err
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			return nil, err
		}
	}

	return &msgs, nil
}

func (s *Repository) CountFindMessage(ctx context.Context, args map[string]interface{}) (*int64, error) {
	var total int64
	ch := make(chan error)
	go func() {
		defer close(ch)

		tx := getFindMessageDB(s.db, args)
		// some filter at here
		err := tx.Count(&total).Error; 
		if err != nil {
			log.WithError(err).Error("Fail to count find message")
		}
		ch <- err
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			return nil, err
		}
	}

	return &total, nil
}

func (s *Repository) CreatePublish(ctx context.Context, id int64) (*model.PublishResult, error) {
	st := model.PublishResult{MsgId: id, Ts: time.Now()}

	ch := make(chan error)

	go func() {
		defer close(ch)
		
		ch <- s.db.Create(&st).Error
	}()
	
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		return nil, err
	}

	return &st, nil
}

func (s *Repository) FindPublish(ctx context.Context, id, msgId int64) (*[]model.PublishResult, error) {
	var results []model.PublishResult

	ch := make(chan error)

	go func() {
		defer close(ch)

		tx := s.db.Model(&model.PublishResult{}).Where("id = ? or msg_id = ?", id, msgId).Order("id DESC")
		
		ch <- tx.Find(&results).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			return nil, err
		}
	}

	return &results, nil
}