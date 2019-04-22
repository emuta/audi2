package repository

import (
	"context"
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	_ "audi/pkg/logger"

	"audi/account/server/model"
)

func encryptPassword(val string) (string, error) {
    hashStr, err := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost)
    return string(hashStr), err
}

func comparePassword(hashStr string, val string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashStr), []byte(val))
}

func (s *Repository) GetUser(ctx context.Context, id int64) (*model.User, error) {
    user := model.User{Id: id}
    ch := make(chan error) 
    go func() {
        defer close(ch)
        if err := s.db.Take(&user).Error; err != nil {
            log.WithError(err).WithField("id", id).Error("Fail to get user")
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
    return &user, nil
}

func (s *Repository) ValidatePassword(ctx context.Context, name, password string) (*model.User, error) {
    var user model.User
    ch := make(chan error) 
    go func() {
        defer close(ch)
        if err := s.db.Take(&user, "name = ?", name).Error; err != nil {
            log.WithError(err).WithField("name", name).Error("Fail to get user")
            ch <- err
            return
        }

        if comparePassword(user.Password, password) != nil {
            err := errors.New("Incorrect old password")
            log.WithField("name", user.Name).Error(err)
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
    return &user, nil
} 

func (s *Repository) CreateUser(ctx context.Context, name, passwd string, roleId int32) (*model.User, error) {
    user := model.User{Name: name, Password: passwd, RoleId: roleId}
    tx := s.db.Begin()
    if tx.Error != nil {
        return nil, tx.Error
    }

    ch := make(chan error) 
    go func() {
        defer close(ch)

        // check user name
        if err := tx.Take(&model.User{}, "name = ?", name).Error; err == nil {
            err := errors.New("duplicate conflict of name exists")
            log.WithError(err).WithField("name", name).Error("Fail to create user")
            ch <- err
            return
        }

        encryptedPass, err := encryptPassword(passwd)
        if err != nil {
            ch <- err
            return 
        }
        user.Password = encryptedPass
        
        if err := tx.Create(&user).Error; err != nil {
            log.WithError(err).Error("Fail to create user")
            ch <- err
            return 
        }

        ch <- nil
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

        if err := tx.Commit().Error; err != nil {
            return nil, err
        }
    }
    return &user, nil
}

func getFindUserDB(db *gorm.DB, args map[string]interface{}) *gorm.DB {
	tx := db.Model(&model.User{})

	if id, ok := args["id"]; ok {
	    tx = tx.Where("id = ?", id)
	}

	if name, ok := args["name"]; ok {
	    tx = tx.Where("name = ?", name)
	}

	if joinFrom, ok := args["join_from"]; ok {
	    tx = tx.Where("join_time >= ?", joinFrom)
	}

	if joinTo, ok := args["join_to"]; ok {
	    tx = tx.Where("join_time <= ?", joinTo)
	}

	if roleId, ok := args["role_id"]; ok {
	    tx = tx.Where("role_id = ?", roleId)
	}

	if active, ok := args["active"]; ok {
	    value, _ := active.(string)
	    if v, err := strconv.ParseBool(value); err == nil {
	        tx = tx.Where("active = ?", v)
	    }
	}

	if offset, ok := args["offset"]; ok  {
	    tx = tx.Offset(offset)
	}

	if limit, ok := args["limit"]; ok  {
	    tx = tx.Limit(limit)
	}

	return tx
}

func (s *Repository) FindUser(ctx context.Context, args map[string]interface{}) (*[]model.User, error){
    var users []model.User

    ch := make(chan error)
    go func(){
        defer close(ch)
        
        tx := getFindUserDB(s.db, args)
        err := tx.Order("id DESC").Find(&users).Error
        if err != nil {
            log.WithError(err).Error("Fail to find users")
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
    return &users, nil
}

func (s *Repository) CountFindUser(ctx context.Context, args map[string]interface{}) (*int64, error){
    var total int64

    ch := make(chan error)
    go func(){
        defer close(ch)
        
        tx := getFindUserDB(s.db, args)
        err := tx.Count(&total).Error
        if err != nil {
            log.WithError(err).Error("Fail to count find users")
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
    return &total, nil
}


func (s *Repository) ChangeUserPassword(ctx context.Context, id int64, oldPass string, newPass string) (bool, error) {
    ch := make(chan error)
    tx := s.db.Begin().Set("gorm:query_option", "for update")
    go func(){
        defer close(ch)
        var user model.User
        if err := tx.Take(&user, "id = ?", id).Error; err != nil {
            log.WithError(err).Error("Fail to fetch user for update password")
            ch <- err
            return
        }
        if comparePassword(user.Password, oldPass) != nil {
            err := errors.New("Incorrect old password")
            log.WithField("name", user.Name).WithField("value", oldPass).Error(err)
            ch <- err
            return
        }

        encryptedPass, err := encryptPassword(newPass)
        if err != nil {
            ch <- err
            return 
        }

        if err := tx.Model(&user).Update("password", encryptedPass).Error; err != nil {
            log.WithError(err).Error("Fail to change user password")
            ch <- err
            return 
        }

        log.WithField("name", user.Name).Info("password changed")

        ch <- nil

        }()

    select {
    case <-ctx.Done():
        tx.Rollback()
        return false, ctx.Err()
    case err := <- ch:
        if err != nil {
            tx.Rollback()
            return false, err
        }

        if err := tx.Commit().Error; err != nil {
            return false, err
        }
    }
    return true, nil
}

func (s *Repository) UpdateUserRole(ctx context.Context, id int64, roleId int32) (bool, error) {
    ch := make(chan error)
    tx := s.db.Begin().Set("gorm:query_option", "for update")
    go func(){
        defer close(ch)
        var user model.User
        if err := tx.Take(&user, "id = ?", id).Error; err != nil {
            log.WithError(err).Error("Fail to fetch user for update role")
            ch <- err
            return
        }

        if err := tx.Model(&user).Update("role_id", roleId).Error; err != nil {
            log.WithError(err).Error("Fail to update user role")
            ch <- err
            return 
        }

        log.WithField("name", user.Name).Info("role updated")

        ch <- nil

        }()

    select {
    case <-ctx.Done():
        tx.Rollback()
        return false, ctx.Err()
    case err := <- ch:
        if err != nil {
            tx.Rollback()
            return false, err
        }

        if err := tx.Commit().Error; err != nil {
            return false, err
        }
    }
    return true, nil
}

func (s *Repository) UpdateUserActive(ctx context.Context, id int64, active bool) (bool, error) {
    ch := make(chan error)
    tx := s.db.Begin().Set("gorm:query_option", "for update")
    go func(){
        defer close(ch)
        var user model.User
        if err := tx.Take(&user, "id = ?", id).Error; err != nil {
            log.WithError(err).Error("Fail to fetch user for Ban")
            ch <- err
            return
        }

        if err := tx.Model(&user).Updates(map[string]interface{}{"active": active}).Error; err != nil {
            log.WithError(err).Error("Fail to Ban user")
            ch <- err
            return 
        }

        log.WithField("name", user.Name).Info("user banned")

        ch <- nil

        }()

    select {
    case <-ctx.Done():
        tx.Rollback()
        return false, ctx.Err()
    case err := <- ch:
        if err != nil {
            tx.Rollback()
            return false, err
        }

        if err := tx.Commit().Error; err != nil {
            return false, err
        }
    }
    return true, nil
}
