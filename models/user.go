package models
import (
  "errors"
  "golang.org/x/crypto/bcrypt"
  _ "github.com/lib/pq"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "lenslocked.com/helpers/rand"
  "lenslocked.com/helpers/hash"
)

const (
  pwPepper = "secret-random-string"
  hmacSecretKey = "secret-hmac-key"
)
var (
  ErrNotFound = errors.New("models: resource not found")
  ErrInvalidID = errors.New("models: ID provided was invalid")
  ErrInvalidPassword = errors.New("models: incorrect password provided")
)



type User struct {
  gorm.Model
  Name string `gorm:"not null"`
  Email string `gorm:"not null;unique_index"`
  Age uint `gorm:"not null"`
  Password string `gorm:"-"`
  PasswordHash string `gorm:"not null"`
  Remember string `gorm:"-"`
  RememberHash string `gorm:not null;unique_index`
}
type UserDB interface {
  ByID(id uint) (*User,error)
  ByEmail(email string) (*User,error)
  ByRemember(token string) (*User,error)
  ByAge(age uint) (*User,error)
  InAgeRange(min, max uint) []User
  Create(userPtr *User) error
  Update(user *User) error
  Delete(id uint) error
  DestructiveReset() error
  AutoMigrate() error
  Close() error
}

type userGORM struct {
  db *gorm.DB
  hmac hash.HMAC
}

cont _ UserDB = &userGORM{}

func newUserGORM(connInfo string) (*userGORM,error) {
  db, err := gorm.Open("postgres",connInfo)
  if err != nil {
    return nil, err
  }
  db.LogMode(true)
  hmac := hash.NewHMAC(hmacSecretKey)
  return &userGORM {
    db : db,
    hmac: hmac,
  }, nil
}




func (usPtr *userGORM) Create(userPtr *User) error {
  pwBytes := []byte(userPtr.Password + pwPepper)
  hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes,
  bcrypt.DefaultCost)
  if err != nil {
    return err
  }
  userPtr.PasswordHash = string(hashedBytes)
  userPtr.Password = ""

  if userPtr.Remember == "" {
    token, err := rand.RememberToken()
    if err != nil {
      return err
    }
    userPtr.Remember = token
  }

//hash the remeber token to use in a cookie
  userPtr.RememberHash = usPtr.hmac.Hash(userPtr.Remember)
  return usPtr.db.Create(userPtr).Error
}

func (us *userGORM) Update(user *User) error {
  if user.Remember != "" {
    user.RememberHash = us.hmac.Hash(user.Remember)
  }
  return us.db.Save(user).Error
}

func (us *userGORM) Delete(id uint) error {
  if id == 0 {
    return ErrInvalidID
  }
  user := User{Model: gorm.Model{ID: id}}
  return us.db.Delete(&user).Error
}

func first(db *gorm.DB, dst interface{}) error {
  err := db.First(dst).Error
  if err != nil {
    return ErrNotFound
  }
  return err
}

func (us *userGORM) ByID(id uint) (*User,error) {
  var user User
  db := us.db.Where("id = ?",id)
  err := first(db,&user)
  if err != nil {
    return nil,err
  }
  return &user, nil
}

func (us *userGORM) ByEmail(email string) (*User,error) {
  var user User
  db := us.db.Where("email = ?",email)
  err := first(db,&user)
  if err != nil {
    return nil, err
  }
  return &user, nil
}

func (us *userGORM) ByRemember(token string) (*User,error) {
  var user User
  rememberHash := us.hmac.Hash(token)
  err := first(us.db.Where("remember_hash = ?", rememberHash), &user)

  if err != nil {
    return nil, err
  }
  return &user, nil
}

func (us *userGORM) ByAge(age uint) (*User,error) {
  var user User
  db := us.db.Where("age = ?",age)
  err := first(db, &user)
  if err != nil {
    return nil, err
  }
  return &user, nil
}

func (us *userGORM) InAgeRange(min, max uint) []User {
  var users []User
  us.db.Where("age BETWEEN ? AND ?",min,max).Find(&users)
  return users
}

func (us *userGORM) AutoMigrate() error {
  if err := us.db.AutoMigrate(&User{}).Error; err != nil {
    return err
  }
  return nil
}

func (us *userGORM) DestructiveReset() error {
  err := us.db.DropTableIfExists(&User{}).Error
  if err != nil {
    return err
  }
  return us.AutoMigrate()
}

func (us *userGORM) Close() error {
  return us.db.Close()
}

type UserService struct {
  db *gorm.DB
  hmac hash.HMAC
}

func (us *UserService) Authenticate(email, password string) (*User,error) {
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
    return nil, ErrInvalidPassword
  default:
    return nil, err
  }
}
func NewUserService(connInfo string) (*UserService,error) {
  db, err := gorm.Open("postgres",connInfo)
  if err != nil {
    return nil, err
  }
  db.LogMode(true)
  hmac := hash.NewHMAC(hmacSecretKey)
  return &UserService {
    db : db,
    hmac: hmac,
  }, nil
}
