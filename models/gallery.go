package models

import "github.com/jinzhu/gorm"

const(
	ErrUserIDRequired modelError = "model: user ID is requered"
	ErrTitleRequierd modelError = "model: title is requered"
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

type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
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

func (gv *galleryValidator) userIDReq(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleReq(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequierd
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

type galleryValidator struct {
	GalleryDB
}

func (gg *galleryGORM) ByID(id uint) (*Gallery ,error) {
	var gallery Gallery
	db := gg.db.Where("id = ?",id)
	err := first(db,&gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (gg *galleryGORM) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGORM) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

type galleryGORM struct {
	db *gorm.DB
}

var _ GalleryDB = &galleryGORM{}
