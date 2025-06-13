package domain

import (
	"btaskee/config"
	"btaskee/db"
	"btaskee/model/constant"
	dbModel "btaskee/model/db"
	reqModel "btaskee/model/request"
	"btaskee/model/response"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func Login(c *gin.Context, user *reqModel.User) *response.Token {
	if user == nil {
		c.Error(errors.New("User not found"))
		c.Status(http.StatusBadRequest)
		return nil
	}

	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)
	if user.Username == "" || user.Password == "" {
		c.Error(errors.New("missing info"))
		c.Status(http.StatusBadRequest)
		return nil
	}

	// Check if username exists
	existsUser, err := db.GetUserByUsername(c, user.Username)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
	}
	if existsUser == nil {
		c.Error(errors.New("User not found"))
		c.Status(http.StatusNotFound)
		return nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(existsUser.Password), []byte(user.Password)); err != nil {
		c.Error(err)
		c.Status(http.StatusUnauthorized)
		return nil
	}

	token, expireAt, err := config.GenerateToken(existsUser.ID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Status(http.StatusOK)
	return &response.Token{
		Token:    token,
		ExpireAt: expireAt,
	}
}

func Signup(c *gin.Context, user *reqModel.User) *response.User {
	if user == nil {
		c.Error(errors.New("User not found"))
		c.Status(http.StatusBadRequest)
		return nil
	}

	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)
	user.Name = strings.TrimSpace(user.Name)
	if user.Username == "" || user.Password == "" || user.Name == "" {
		c.Error(errors.New("missing info"))
		c.Status(http.StatusBadRequest)
		return nil
	}

	// Check if username exists
	existsUser, err := db.GetUserByUsername(c, user.Username)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
	}
	if existsUser != nil {
		c.Error(errors.New("User already exists"))
		c.Status(http.StatusConflict)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	dbUser := dbModel.User{
		Name:     user.Name,
		Username: user.Username,
		Password: string(hash),
	}
	if err = db.CreateUser(c, &dbUser); err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	c.Status(http.StatusOK)
	return &response.User{
		Username: user.Username,
		Name:     user.Name,
	}
}

func GetMe(c *gin.Context) *response.User {
	userID := c.GetUint(constant.USER_ID_CTX)
	if userID == 0 {
		c.Error(errors.New("missing user_id in context"))
		c.Status(http.StatusUnauthorized)
		return nil
	}

	user, err := db.GetUsersByID(c, userID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	if user == nil {
		c.Error(errors.New("user not found"))
		c.Status(http.StatusNotFound)
		return nil
	}

	// Chuyển đổi từ db model sang response model
	return &response.User{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}
}
