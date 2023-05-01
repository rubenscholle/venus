package authbundle

import (
	"errors"
	"log"
	"net/http"

	"crypto/md5"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func (con *authController) LoginHandler(c *gin.Context) {
	var user AuthUser

	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	hashedBytes := md5.Sum([]byte(user.Password))
	hashedPassword := hex.EncodeToString(hashedBytes[:])
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	log.Println(err)
	//	c.Error(err)
	//	return
	//}

	var userDb AuthUser
	con.OrmDb.Where("username=? AND password=?", user.Username, hashedPassword).First(&userDb)
	if userDb.ID == 0 {
		err := errors.New("invalid user or password")
		log.Println(err)
		c.Error(err)
		return
	}

	jwt, err := GenerateJWT(userDb)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": jwt})
}
