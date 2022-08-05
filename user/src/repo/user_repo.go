package repo

import (
	"errors"
	"github.com/cemayan/event-scraper/user/src/database"
	"github.com/cemayan/event-scraper/user/src/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userrepo struct {
	db  *gorm.DB
	log *log.Logger
}

func NewUserRepo(db *gorm.DB, log *log.Logger) model.UserRepository {
	return &userrepo{
		db:  db,
		log: log,
	}
}

func (r userrepo) UpdateUser(user *model.User) {
	r.db.Save(user)
}

func (r userrepo) DeleteUser(id uint) {
	user, err := r.GetUserById(id)
	if user == nil || err != nil {
		return
	}
	r.db.Delete(user)
}

func (r userrepo) CreateUser(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetUserById returns user based on given id
func (r userrepo) GetUserById(id uint) (*model.User, error) {
	var user model.User
	if err := database.DB.Find(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername returns user based on given username
func (r userrepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where(&model.User{Username: username}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail returns user based on given id
func (r userrepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where(&model.User{Email: email}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
