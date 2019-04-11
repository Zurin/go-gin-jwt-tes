package main

import (
	"fmt"
	"gin-jwt/config"
	"gin-jwt/controllers"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	db := config.DBInit()
	inDB := &controllers.InDB{DB: db}

	router := gin.Default()

	router.POST("/login", loginHandler)
	router.GET("/person/:id", auth, inDB.GetPerson)
	router.GET("/persons", auth, inDB.GetPersons)
	router.POST("/person", auth, inDB.CreatePerson)
	router.PUT("/person", auth, inDB.UpdatePerson)
	router.DELETE("/person/:id", auth, inDB.DeletePerson)
	router.Run(":3000")
}

func loginHandler(c *gin.Context) {
	var user Credential
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "can't bind struct",
		})
	}
	if user.Username != "zurin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "wrong username or password",
		})
		c.Abort()
	} else {
		if user.Password != "123456" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "wrong username or password",
			})
			c.Abort()
		} else {
			claims := &jwt.StandardClaims{
				ExpiresAt: 15000,
				Issuer:    "zurin",
			}

			// sign := jwt.New(jwt.GetSigningMethod("HS256"))
			sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token, err := sign.SignedString([]byte("secret"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				c.Abort()
			}
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		}
	}
}

func auth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if token != nil && err == nil {
		fmt.Println("token verified")
	} else {
		result := gin.H{
			"message": "not authorized",
			"error":   err.Error(),
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
	}
}
