package repository

import (
	"context"
	"strconv"
    "time"

    "github.com/lib/pq"
    "github.com/jinzhu/gorm"

	"audi/lottery/ssc/cqssc/server/model"
)

func (s *Repository) CreateBet(
		ctx context.Context,
		userId int64,
		odds float64,
		playId int32,
		unitId int64,
		comm float64,
		chaseStop bool,
		codes []string,
		plans []map[string]interface{}) (*model.Bet, error) {

	bet := model.Bet{
		UserId:    userId,
		Odds:      odds,
		PlayId:    playId,
		UnitId:    unitId,
		Comm:      comm,
		ChaseStop: chaseStop,
		Codes:     codes,
	}
	tx := s.db.Begin()
	ch := make(chan error)

	go func() {
		defer close(ch)

		if err := tx.Create(&bet).Error; err != nil {
			ch <- err
			return
		}

		for _, p := range plans {
			plan := model.BetPlan{BetId: bet.Id}
			if v, ok := p["term_id"].(int64); ok {
				plan.TermId = v
			}
			if v, ok := p["times"].(int64); ok {
				plan.Times = v
			}
			if v, ok := p["qty"].(int64); ok {
				plan.Qty = v
			}
			if v, ok := p["payment"].(float64); ok {
				plan.Payment = v
			}
			if v, ok := p["bonus"].(float64); ok {
				plan.Bonus = v
			}
			if v, ok := p["rebate"].(float64); ok {
				plan.Rebate = v
			}

			if err := tx.Create(&plan).Error; err != nil {
				ch <- err
				return
			}
		}

		ch <- nil

	} ()

	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return &bet, tx.Commit().Error
}


func (s *Repository) GetBet(ctx context.Context, betId int64) (*model.Bet, error) {
	result := model.Bet{Id: betId}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&result).Error
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

func getFindBetDB(db *gorm.DB, args map[string]interface{}) *gorm.DB {
	tx := db.Model(&model.Bet{})

	if userId, ok := args["user_id"]; ok {
		tx = tx.Where("user_id = ?", userId)
	}

	if playId, ok := args["play_id"]; ok {
		tx = tx.Where("play_id = ?", playId)
	}

	if code, ok := args["code"]; ok {
		tx = tx.Where("? = any(codes)", code)
	}

	if createdFrom, ok := args["created_from"]; ok {
		tx = tx.Where("created_from >= ?", createdFrom)
	}

	if createdTo, ok := args["created_to"]; ok {
		tx = tx.Where("created_to <= ?", createdTo)
	}

	if termId, ok := args["term_id"]; ok {
		tx = tx.Where("term_id = ?", termId)
	}

	if win, ok := args["win"]; ok {
		winStr, _ := win.(string)
		if v, err := strconv.ParseBool(winStr); err == nil {
			if v {
				tx = tx.Where(" id in (?)", db.Model(&model.BetPlanStats{}).Select("bet_id").Where("win > 0").QueryExpr())
			} else {
				tx = tx.Where(" id in (?)", db.Model(&model.BetPlanStats{}).Select("bet_id").Where("win = 0").QueryExpr())
			}
			
		}
	}

	if settled, ok := args["settled"]; ok {
		settledStr, _ := settled.(string)
		if v, err := strconv.ParseBool(settledStr); err == nil {
			tx = tx.Where(" id in (?)", db.Model(&model.BetPlanStats{}).Select("bet_id").Where("settled = ?", v).QueryExpr())
		}
	}

	if revoked, ok := args["revoked"]; ok {
		revokedStr, _ := revoked.(string)
		if v, err := strconv.ParseBool(revokedStr); err == nil {
			tx = tx.Where(" id in (?)", db.Model(&model.BetPlanStats{}).Select("bet_id").Where("revoked = ?", v).QueryExpr())
		}
	}

	

	return tx
}

func (s *Repository) FindBet(ctx context.Context, args map[string]interface{}) (*[]model.Bet, error) {
	var result []model.Bet
	ch := make(chan error)
	go func() {
		defer close(ch)
		tx := getFindBetDB(s.db, args)

		if limit, ok := args["limit"]; ok {
			tx = tx.Limit(limit)
		}

		if offset, ok := args["offset"]; ok {
			tx = tx.Offset(offset)
		}

		ch <- tx.Order("id DESC").Find(&result).Error
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

func (s *Repository) CountFindBet(ctx context.Context, args map[string]interface{}) (*int32, error) {
	var result int32
	ch := make(chan error)
	go func() {
		defer close(ch)
		tx := getFindBetDB(s.db, args)
		ch <- tx.Count(&result).Error
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

func (s *Repository) GetBetPlan(ctx context.Context, planId int64) (*model.BetPlan, error) {
	result := model.BetPlan{Id: planId}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&result).Error
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

func (s *Repository) GetBetStats(ctx context.Context, betId int64) (*[]model.BetPlanStats, error) {
	var result []model.BetPlanStats
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Where("bet_id = ?", betId).Order("id ASC").Find(&result).Error
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

func (s *Repository) GetBetPlanStats(ctx context.Context, planId int64) (*model.BetPlanStats, error) {
	result := model.BetPlanStats{Id: planId}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&result).Error
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

func (s *Repository) SettleBetPlan(ctx context.Context, planId int64, win int64) (*model.BetPlanStats, error) {
	var result model.BetPlanStats
	tx := s.db.Begin().Set("gorm:query_option", "for update")
	ch := make(chan error)
	go func() {
		defer close(ch)
		if err := tx.Take(&result, "id = ? and settled = ?", planId, false).Error; err != nil {
			ch <- err
			return
		}

		plan := model.BetPlan{Id: result.Id}
		if err := tx.Take(&plan).Error; err != nil {
			ch <- err
			return
		}

		var settledAt pq.NullTime
		settledAt.Scan(time.Now())

		result.Win = win
		result.Settled = true
		result.SettledAt = settledAt
		result.Bonus  = float64(win) * plan.Bonus
		result.Rebate = plan.Rebate

		ch <- tx.Save(&result).Error


	}()

	select {
	case <- ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return &result, tx.Commit().Error
}


func (s *Repository) RevokeBetPlan(ctx context.Context, planId int64) (*model.BetPlanStats, error) {
	var result model.BetPlanStats
	tx := s.db.Begin().Set("gorm:query_option", "for update")
	ch := make(chan error)
	go func() {
		defer close(ch)
		if err := tx.Take(&result, "id = ? and revoked = ?", planId, false).Error; err != nil {
			ch <- err
			return
		}

		plan := model.BetPlan{Id: result.Id}
		if err := tx.Take(&plan).Error; err != nil {
			ch <- err
			return
		}

		var revokedAt pq.NullTime
		revokedAt.Scan(time.Now())

		args := map[string]interface{}{
			"revoked": true,
			"revoked_at": revokedAt,
			"payment":0,
			"bonus": 0,
			"rebate": 0,
		}

		ch <- tx.Model(&model.BetPlanStats{Id: planId}).Updates(args).Error;

		/*

		result.Revoked = true
		result.RevokedAt = revokedAt
		result.Payment = 0
		result.Bonus   = 0
		result.Rebate  = 0

		ch <- tx.Save(&result).Error
		*/
	}()

	select {
	case <- ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return &result, tx.Commit().Error
}