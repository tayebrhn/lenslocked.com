package models
import (
	"regexp"
	"strings"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"lenslocked.com/helpers/hash"
	"lenslocked.com/helpers/rand"
)
const (
	pwPepper      = "secret-random-string"
	hmacSecretKey = "secret-hmac-key"
)
const (
	ErrNotFound          modelError = "models: resource not found"
	ErrIDInvalid         modelError = "models: ID provided was invalid"
	ErrEmailRequired     modelError = "models: email address required"
	ErrEmailInvalid      modelError = "models: email address is not valid"
	ErrEmailTaken        modelError = "models: email address already taken"
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	ErrPasswordToShort   modelError = "models: password must " +
		"be atleast 8 characters long"
	ErrPasswordRequired modelError = "models: password is required"
	ErrRemRequired      modelError = "models: remember token is required"
	ErrRemToShort       modelError = "models: remeber token" +
		"must be atleast 32 bytes"
)
func (e modelError) Error() string {
	return string(e)
}
func (e modelError) Public() string {
	str := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(str, " ")
	split[0] = strings.ToUpper(split[0])
	return strings.Join(split, " ")
}
type modelError string
type UserDB interface {
	//Method for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)
	ByAge(age uint) (*User, error)
	//Method for querying for multiple users
	InAgeRange(min, max uint) []User
	//Method for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error	
}
type UserService interface {
	UserDB
	Authenticate(email, password string) (*User, error)
}
func NewUserService(db *gorm.DB) UserService {
	userg := &userGORM{db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(userg, hmac)
	return &userService{
		UserDB: uv,
	}
}
type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null;unique_index"`
	Age          uint   `gorm:"not null"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}
func (us *userService) Authenticate(email, password string) (*User, error) {
	uv := newUserValidator(us.UserDB, hash.NewHMAC(hmacSecretKey))
	user := &User{
		Password: password,
	}
	err := runUserValFns(
		user,
		uv.pwReq,
		uv.pwMinLen)

	if err != nil {
		return nil, err
	}

	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash),
		[]byte(password+pwPepper))

	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrPasswordIncorrect
	default:
		return nil, err
	}
}
type userService struct {
	UserDB
}
var _ UserService = &userService{}
func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}
type userValFn func(*User) error
// Validation and Data Normalization helper Functions
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + pwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes,
		bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}
func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}
func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	//hash the remeber token to use in a cookie
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}
func (uv *userValidator) idGtThan(id uint) userValFn {
	return userValFn(func(user *User) error {
		if user.ID <= id {
			return ErrIDInvalid
		}
		return nil
	})
}
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}
func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}
func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}
func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}
func (uv *userValidator) pwMinLen(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordToShort
	}
	return nil
}
func (uv *userValidator) pwReq(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}
func (uv *userValidator) pwHashReq(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}
func (uv *userValidator) remMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}

	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRemToShort
	}
	return nil
}
// valdation function
func (uv *userValidator) Create(user *User) error {

	err := runUserValFns(
		user,
		uv.pwReq,
		uv.pwMinLen,
		uv.bcryptPassword,
		uv.pwHashReq,
		uv.setRememberIfUnset,
		uv.remMinBytes,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}
func (uv *userValidator) Update(user *User) error {
	err := runUserValFns(
		user,
		uv.pwMinLen,
		uv.bcryptPassword,
		uv.pwHashReq,
		uv.remMinBytes,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	if err := runUserValFns(&user, uv.idGtThan(0)); err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFns(&user,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}
func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB: udb,
		hmac:   hmac,
		emailRegex: regexp.MustCompile(
			`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}
// Database Manipulation function
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err != nil {
		return ErrNotFound
	}
	return err
}
func (us *userGORM) Create(user *User) error {
	return us.db.Create(user).Error
}
func (us *userGORM) Update(user *User) error {
	return us.db.Save(user).Error
}
func (us *userGORM) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}
func (us *userGORM) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (us *userGORM) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (us *userGORM) ByRemember(token string) (*User, error) {
	var user User
	err := first(us.db.Where("remember_hash = ?", token), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (us *userGORM) ByAge(age uint) (*User, error) {
	var user User
	db := us.db.Where("age = ?", age)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (us *userGORM) InAgeRange(min, max uint) []User {
	var users []User
	us.db.Where("age BETWEEN ? AND ?", min, max).Find(&users)
	return users
}
type userGORM struct {
	db *gorm.DB
}
var _ UserDB = &userGORM{}