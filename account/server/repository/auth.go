package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	_ "audi/pkg/logger"

	"audi/account/server/model"
	"audi/account/server/config"
)

func generateToken (token *model.AuthToken, key string) (string, error) {
    claims := &jwt.StandardClaims{
        Id:        strconv.FormatInt(token.AuthId, 10),
        Audience:  strconv.FormatInt(token.UserId, 10),
        Subject:   token.UserName,
        IssuedAt:  token.CreatedAt.Unix(),
        ExpiresAt: token.ExpiredAt.Unix(),
    }

    tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signature, err := tk.SignedString([]byte(key))
    return signature, err
}

func parseToken (val string, key string) (interface{}, error) {
    token, err := jwt.ParseWithClaims(
        val,
        &jwt.StandardClaims{},
        func(token *jwt.Token)(interface{}, error){
            return []byte(key), nil
        })
    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
        return claims.Id, nil
    }

    if ve, ok := err.(*jwt.ValidationError); ok {
        if ve.Errors&jwt.ValidationErrorMalformed != 0 {
            return nil, errors.New("Token is not valid")
        } else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
            return nil, errors.New("Token has been expired")
        } else {
            return nil, errors.New("Token is not valid")
        }
    }

    return nil, errors.New("Token is not valid")
}


func (s *Repository) AuthLogin(ctx context.Context, name, passwd, device, ip, fp string) (*model.AuthToken, error) {
	var token model.AuthToken
	now := time.Now()
	tx := s.db.Begin().Set("gorm:query_option", "for update")
	ch := make(chan error)

	go func() {
		defer close(ch)

		// get user
		var user model.User
		if err := tx.Take(&user, "name = ?", name).Error; err != nil {
			log.WithField("name", name).Error("Failed to get user on login")
			ch <- err
			return
		}

		// compare password
		if err := comparePassword(user.Password, passwd); err != nil {
			log.WithFields(log.Fields{
				"name": name,
				"password": passwd,
			}).Error("Failed to validate user password on login")
			return
		}

		// create auth history
		auth := model.AuthHistory{UserId: user.Id, Device: device, Ip: ip, Fp: fp, CreatedAt: now}
		if err := tx.Create(&auth).Error; err != nil {
			log.Error("Failed to record auth history")
			ch <- err
		}

		// create token
		token.AuthId    = auth.Id
		token.UserId    = user.Id
		token.UserName  = user.Name
		token.CreatedAt = now 
		token.Active    = true
		token.ExpiredAt = token.CreatedAt.Add(time.Duration(config.TOKEN_EXPIRES) * time.Second)
		// generate signature
		signature, err := generateToken(&token, config.TOKEN_SECRET)
		if err != nil {
			log.Error("Failed to generate token")
			ch <- err
			return
		}
		token.Signature = signature

		ch <- tx.Create(&token).Error
	}()

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
	return &token, tx.Commit().Error

}

func (s *Repository) AuthLogout(ctx context.Context, signature string) (error) {
	ch := make(chan error)
	go func() {
		defer close(ch)

		authId, err := parseToken(signature, config.TOKEN_SECRET)
		if err != nil {
			ch <- err
			return
		}

		if err := s.db.Model(&model.AuthToken{}).
				Where("auth_id = ?", authId).
				Update("active = ?", false).Error; err != nil {
			ch <- err
			return
		}
		ch <- nil
	}()
	select {
	case <-ctx.Done():
	    return ctx.Err()
	case err := <- ch:
	    if err != nil {
	        return err
	    }
	}
	return nil
}

func (s *Repository) GetAuthHistory(ctx context.Context, id int64) (*model.AuthHistory, error) {
	result := model.AuthHistory{Id: id}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&result).Error
	}()

	select {
	case <-ctx.Done():
	    return nil, ctx.Err()
	case err := <- ch:
	    if err != nil {
	        return nil, err
	    }
	}
	return &result, nil
}

func getFindAuthHistoryDB(db *gorm.DB, args map[string]interface{}) (tx *gorm.DB) {
	tx = db.Model(&model.AuthHistory{})

	if userId, ok := args["user_id"]; ok {
	    tx = tx.Where("user_id = ?", userId)
	}

	/*
	if userName, ok := args["user_name"]; ok {
	    tx = tx.Where("user_id in (?)", tx.Model(&model.User{}).Select("id").Where("name = ?", userName).QueryExpr())
	}
	*/

	if device, ok := args["device"]; ok {
	    tx = tx.Where("device = ?", device)
	}

	if ip, ok := args["ip"]; ok {
	    tx = tx.Where("ip = ?", ip)
	}

	if fp, ok := args["fp"]; ok {
	    tx = tx.Where("fp = ?", fp)
	}

	if createdFrom, ok := args["created_from"]; ok {
	    tx = tx.Where("created_at >= ?", createdFrom)
	}

	if createdTo, ok := args["created_to"]; ok {
	    tx = tx.Where("created_at <= ?", createdTo)
	}

	if limit, ok := args["limit"]; ok {
	    tx = tx.Limit(limit)
	}

	if offset, ok := args["offset"]; ok {
	    tx = tx.Offset(offset)
	}

	// set order by must id desc
	tx = tx.Order("id DESC")

	return tx
}

func (s *Repository) FindAuthHistory(ctx context.Context, args map[string]interface{}) (*[]model.AuthHistory, error) {
	var result []model.AuthHistory
	ch := make(chan error)
	go func() {
		defer close(ch)
		tx := getFindAuthHistoryDB(s.db, args)
		ch <- tx.Find(&result).Error
	}()

	select {
	case <-ctx.Done():
	    return nil, ctx.Err()
	case err := <- ch:
	    if err != nil {
	        return nil, err
	    }
	}
	return &result, nil
}

func (s *Repository) CountFindAuthHistory(ctx context.Context, args map[string]interface{}) (*int64, error) {
	var result int64
	ch := make(chan error)
	go func() {
		defer close(ch)
		tx := getFindAuthHistoryDB(s.db, args)
		ch <- tx.Count(&result).Error
	}()

	select {
	case <-ctx.Done():
	    return nil, ctx.Err()
	case err := <- ch:
	    if err != nil {
	        return nil, err
	    }
	}
	return &result, nil
}

func (s *Repository) GetAuthToken(ctx context.Context, id int64) (*model.AuthToken, error) {
	result := model.AuthToken{Id: id}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- s.db.Take(&result).Error
	}()

	select {
	case <-ctx.Done():
	    return nil, ctx.Err()
	case err := <- ch:
	    if err != nil {
	        return nil, err
	    }
	}
	return &result, nil
}

func getFindAuthTokenDB(db *gorm.DB, args map[string]interface{}) *gorm.DB {
	tx := db.Model(&model.AuthToken{})

	if authId, ok := args["auth_id"]; ok {
	    tx = tx.Where("auth_id = ?", authId)
	}

	if userId, ok := args["user_id"]; ok {
	    tx = tx.Where("user_id = ?", userId)
	}

	if userName, ok := args["user_name"]; ok {
	    tx = tx.Where("user_name = ?", userName)
	}

	if createdFrom, ok := args["created_from"]; ok {
	    tx = tx.Where("created_at >= ?", createdFrom)
	}

	if createdTo, ok := args["created_to"]; ok {
	    tx = tx.Where("created_at <= ?", createdTo)
	}

	if expiredFrom, ok := args["expired_from"]; ok {
	    tx = tx.Where("expired_at >= ?", expiredFrom)
	}

	if expiredAtTo, ok := args["expired_to"]; ok {
	    tx = tx.Where("expired_at <= ?", expiredAtTo)
	}

	if active, ok := args["active"]; ok {
	    value, _ := active.(string)
	    if v, err := strconv.ParseBool(value); err == nil {
	        tx = tx.Where("active = ?", v)
	    }
	}

	if limit, ok := args["limit"]; ok {
	    tx = tx.Limit(limit)
	}

	if offset, ok := args["offset"]; ok {
	    tx = tx.Offset(offset)
	}

	return tx
}

func (s *Repository) FindAuthToken(ctx context.Context, args map[string]interface{}) (*[]model.AuthToken, error) {
	var result []model.AuthToken
	ch := make(chan error)
	go func() {
		defer close(ch)
		tx := getFindAuthTokenDB(s.db, args)
		ch <- tx.Find(&result).Error
	}()

	select {
	case <-ctx.Done():
	    return nil, ctx.Err()
	case err := <- ch:
	    if err != nil {
	        return nil, err
	    }
	}
	return &result, nil
}

func (s *Repository) CountFindAuthToken(ctx context.Context, args map[string]interface{}) (*int64, error) {
	var result int64
	ch := make(chan error)
	go func() {
		defer close(ch)
		tx := getFindAuthTokenDB(s.db, args)
		ch <- tx.Count(&result).Error
	}()

	select {
	case <-ctx.Done():
	    return nil, ctx.Err()
	case err := <- ch:
	    if err != nil {
	        return nil, err
	    }
	}
	return &result, nil
}

func (s *Repository) InactiveAuthToken(ctx context.Context, id int64) error {
	ch := make(chan error)
	func () {
		defer close(ch)
		ch <- s.db.Model(&model.AuthToken{Id: id}).Update("active = ?", false).Error
	}()

	select {
	case <-ctx.Done():
	    return ctx.Err()
	case err := <- ch:
	    if err != nil {
	        return err
	    }
	}
	return nil
}

func (s *Repository) ParseToken(ctx context.Context, signature string) (*model.AuthToken, error) {
	var token model.AuthToken
	ch := make(chan error)
	go func() {
		defer close(ch)

		authId, err := parseToken(signature, config.TOKEN_SECRET)
		if err != nil {
			ch <- err
			return
		}
		ch <- s.db.Take(&token, "auth_id = ?", authId).Error
	}()
	select {
	case <-ctx.Done():
	    return nil, ctx.Err()
	case err := <- ch:
	    if err != nil {
	        return nil, err
	    }
	}
	return &token, nil
}







































