package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

type Services struct {
	User    UserService
	Gallery GalleryService
	db *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}


func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{},&Gallery{}).Error
}

func (us *Services) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{},&Gallery{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}


func NewServices(connInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User: NewUserService(db),
		Gallery: NewGalleryService(db),
		db: db,
	}, nil
}