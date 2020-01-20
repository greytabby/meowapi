package handler

import (
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/greytabby/meowapi/lib/model"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type jwtCustomClaims struct {
	UID  int64  `json:"uid"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var signingKey = []byte(os.Getenv("JWT_SIGNING_KEY"))

var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

// UserDbAccessor Userテーブルへのアクセスを行う
type UserDbAccessor interface {
	FindUser(name string) (model.User, error)
	AddUser(user model.User) error
	DeleteUser(user model.User) error
}

// AuthHandler 認証に関するapihandler
type AuthHandler struct {
	Db UserDbAccessor
}

// Signup ユーザ登録を行う
func (ah *AuthHandler) Signup(c echo.Context) error {
	var user model.User
	var err error
	if err = c.Bind(&user); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Invalid field")
	}

	// check the user already exist.
	u, err := ah.Db.FindUser(user.Name)
	if u.Id != 0 {
		c.Logger().Errorf("User already exist", u.Name)
		return c.String(http.StatusConflict, "User already exist")
	}

	// bcrypt password verify ignores 73 characters and more
	if len(user.Password) > 72 {
		return c.String(http.StatusBadRequest, "password length is 72 or less")
	}

	// create hash password
	// save hashed password not plain password.
	user.Password, err = passwordHash(user.Password)
	if err != nil {
		c.Logger().Error("passwordhash failed.", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = ah.Db.AddUser(user)
	if err != nil {
		c.Logger().Error("Create user failed.", err)
		return c.String(http.StatusInternalServerError, "")
	}

	// delete password info from response
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// Login ユーザ情報を照合しjwttokenを発行する
func (ah *AuthHandler) Login(c echo.Context) error {
	var requser model.User
	if err := c.Bind(&requser); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Invalid field")
	}

	// check username and password
	loginUser, err := ah.Db.FindUser(requser.Name)
	if err != nil {
		c.Logger().Errorf("Login: Invalid name", requser.Name, err)
		return c.String(http.StatusBadRequest, "Invalid name or password")
	}
	if err := passwordVerify(loginUser.Password, requser.Password); err != nil {
		c.Logger().Errorf("Login: Invalid password", requser.Name, err)
		return c.String(http.StatusBadRequest, "Invalid name or password")
	}

	// jwt authentication
	claims := &jwtCustomClaims{
		loginUser.Id,
		loginUser.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	// create jwttoken
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	t, err := token.SignedString(signingKey)
	if err != nil {
		c.Logger().Errorf("Login: create token failed", err)
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func passwordHash(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}

func passwordVerify(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

// UserIdFromToken tokenからuseridを得る
func UserIdFromToken(c echo.Context) int64 {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	uid := claims.UID
	return uid
}
