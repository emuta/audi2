package repository

import (
	"context"
	"strconv"

	"github.com/lib/pq"
	// log "github.com/sirupsen/logrus"

	"audi/lottery/ssc/cqssc/server/model"
)

func (s *Repository) GetPlay(ctx context.Context, id int32) (*model.Play, error) {
	item := model.Play{Id: id}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&item).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &item, nil
}

func (s *Repository) FindPlay(ctx context.Context, args map[string]interface{}) (*[]model.Play, error) {
	var items []model.Play
	ch := make(chan error)
	go func() {
		defer close(ch)
		tx := s.db.Model(&model.Play{})

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

		if active, ok := args["active"]; ok {
			activeStr, _ := active.(string)
			if v, err := strconv.ParseBool(activeStr); err == nil {
				tx = tx.Where("active = ?", v)
			}
		}

		if pr, ok := args["pr"]; ok {
			tx = tx.Where("pr = ?", pr)
		}

		if catgId, ok := args["catg_id"]; ok {
			tx = tx.Where("catg_id = ?", catgId)
		}

		if groupId, ok := args["group_id"]; ok {
			tx = tx.Where("group_id = ?", groupId)
		}

		if unitId, ok := args["unit_id"]; ok {
			tx = tx.Where("? = any(units)", unitId)
		}

		ch <- tx.Order("id ASC").Find(&items).Error
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ch:
		if err != nil {
			return nil, err
		}
	}

	return &items, nil
}

func (s *Repository) UpdatePlay(ctx context.Context, id int32, pref, active string, units []int64) (*model.Play, error) {
	item := model.Play{Id: id}

	tx := s.db.Begin().Set("gorm:query_option", "for update")
	ch := make(chan error)
	go func() {
		defer close(ch)

		args := map[string]interface{}{}

		if err := tx.Take(&item).Error; err != nil {
			ch <- err
			return
		}

		if pref != "" {
			if v, err := strconv.ParseBool(pref); err == nil {
				args["pref"] = v

				// if pref is true

				if v {
					// set others catg not preferred
					// map[string]bool{"pref": false}
					if err := tx.Model(&model.Catg{}).Where("id != ?", item.CatgId).Update("pref = ?", false).Error; err != nil {
						ch <- err
						return
					}
					// set current parent catg preferred
					if err := tx.Model(&model.Catg{}).Where("id = ?", item.CatgId).Update("pref = ?", true).Error; err != nil {
						ch <- err
						return
					}
					// set others item not preferred
					if err := tx.Model(&model.Play{}).Where("id != ?", id).Update("pref = ?", false).Error; err != nil {
						ch <- err
						return
					}
				}
			}
		}

		if active != "" {
			if v, err := strconv.ParseBool(active); err == nil {
				args["active"] = v
			}
		}

		args["units"] = pq.Int64Array(units)

		ch <- tx.Model(&item).Updates(args).Error
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

	return &item, tx.Commit().Error
}
