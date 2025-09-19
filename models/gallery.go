package models

import (
	"github.com/jinzhu/gorm"
)

const (
	ErrUserIDRequired modelError = "model: user ID is required"
	ErrTitleRequired  modelError = "model: title is required"
)

type galleryFunc func(*Gallery) error

func runGalleryValFns(gallery *Gallery, fns ...galleryFunc) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gallery) ImageSplitN(n int) [][]Image {
	ret := make([][]Image, n)

	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}
	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)

	}
	return ret
}

type Gallery struct {
	gorm.Model
	UserID uint    `gorm:"not_null;index"`
	Title  string  `gorm:"not_null"`
	Images []Image `gorm:"-"`
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	ByUserID(userID uint) ([]Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGORM{
				db: db,
			},
		},
	}
}

type GalleryService interface {
	GalleryDB
}
type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) nonZero(gallery *Gallery) error {
	if gallery.ID <= 0 {
		return ErrIDInvalid
	}
	return nil
}

func (gv *galleryValidator) userIDReq(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleReq(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFns(
		gallery,
		gv.userIDReq,
		gv.titleReq,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

func (gv *galleryValidator) Delete(id uint) error {
	var gallery Gallery
	gallery.ID = id

	err := runGalleryValFns(
		&gallery,
		gv.nonZero,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Delete(gallery.ID)
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValFns(
		gallery,
		gv.userIDReq,
		gv.titleReq,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

type galleryGORM struct {
	db *gorm.DB
}

var _ GalleryDB = &galleryGORM{}

func (gg *galleryGORM) ByUserID(userID uint) ([]Gallery, error) {
	var galleries []Gallery
	db := gg.db.Where("user_id = ?", userID)
	err := db.Find(&galleries).Error
	if err != nil {
		return nil, err
	}
	return galleries, nil
}

func (gg *galleryGORM) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (gg *galleryGORM) Delete(id uint) error {
	gallery := Gallery{
		Model: gorm.Model{ID: id}}
	return gg.db.Delete(&gallery).Error
}

func (gg *galleryGORM) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGORM) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
