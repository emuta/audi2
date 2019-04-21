package repository

import (
	"context"

	"audi/account/server/model"
)

func (s *Repository) GetRole(ctx context.Context, id int32, name string) (*model.Role, error) {
	var result model.Role

	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&result, "id = ? or name = ?", id, name).Error
	}()

	select {
	case <- ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *Repository) FindRole(ctx context.Context, args map[string]interface{}) (*[]model.Role, error) {
    var roles []model.Role
    ch := make(chan error)
    go func(){
        defer close(ch)
        tx := s.db.Model(&model.Role{})

        if id, ok := args["id"]; ok {
            tx = tx.Where("id = ?", id)
        }

        if name, ok := args["name"]; ok {
            tx = tx.Where("name = ?", name)
        }

        if err := tx.Find(&roles).Error; err != nil {
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
    return &roles, nil
}