package repo

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"

	"audi/message/server/model"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (s *Repo) CreateMessage(ctx context.Context, srv string, event string, data postgres.Jsonb) (*model.Message, error) {
	msg := model.Message{
		Srv:    srv,
		Event:  event,
		Data:   data,
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

func (s *Repo) GetMessage(ctx context.Context, id int64) (*model.Message, error) {
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

func (s *Repo) DeleteMessage(ctx context.Context, id int64) (error) {
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

func (s *Repo) FindMessage(ctx context.Context, args map[string]interface{}) (*[]model.Message, error) {
	var msgs []model.Message
	ch := make(chan error)
	go func() {
		defer close(ch)

		tx := s.db.Model(&model.Message{})

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

func (s *Repo) CreatePublish(ctx context.Context, msgId int64) (*model.Publish, error) {
	st := model.Publish{MessageId: msgId, Ts: time.Now()}

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

func (s *Repo) GetPublishes(ctx context.Context, id int64, limit interface{}, offset interface{}) (*[]model.Publish, error) {
	var sts []model.Publish

	ch := make(chan error)

	go func() {
		defer close(ch)
		
		ch <- s.db.Where("message_id = ?", id).Order("id DESC").Limit(limit).Offset(offset).Find(&sts).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			return nil, err
		}
	}

	return &sts, nil
}















