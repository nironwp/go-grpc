package database

import (
	"database/sql"
	"net/mail"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/nironwp/grpc/configs"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	db       *sql.DB
	ID       string
	Name     string
	Email    string
	password string
}

func NewUser(db *sql.DB) *User {
	return &User{db: db}
}

func (u *User) validEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Create(name string, email string, password string) (*User, error) {
	id := uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = u.validEmail(email)
	if err != nil {
		return nil, err
	}
	_, err = u.db.Exec("INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)", id, name, email, string(hash))
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		password: string(hash),
	}, nil
}

func (u *User) Find(id string) (*User, error) {
	var user User
	err := u.db.QueryRow("SELECT id, name, email, password FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) FindByEmail(email string) (*User, error) {
	var user User
	err := u.db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) Login(email string, password string) (string, error) {

	user, err := u.FindByEmail(email)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password))

	if err != nil {
		return "", err
	}
	configs := configs.LoadConfig(".")
	secretKey := []byte(configs.APPSecret)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["name"] = user.Name
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *User) Update(id string, name string, email string, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	err = u.validEmail(email)
	if err != nil {
		return nil, err
	}
	_, err = u.db.Exec("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?", name, email, string(hash), id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		password: string(hash),
	}, nil
}

func (u *User) Delete(id string) error {
	_, err := u.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
