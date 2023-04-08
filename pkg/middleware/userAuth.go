package middleware

import (
	"fmt"
	"jwt/pkg/database"
	"jwt/pkg/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func UserAuth(c *gin.Context) {

	fmt.Println("In middleware")

	// Get the cookie of req

	tokenString, err := c.Cookie("UserToken")

	if err != nil {

		c.AbortWithStatus(http.StatusUnauthorized)

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "user not found",
		})
	}

	// decode and validate

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if err != nil {

			c.AbortWithStatus(http.StatusUnauthorized)

			c.JSON(http.StatusBadRequest, gin.H{

				"error": "user not found",
			})
		}

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")

		return []byte(os.Getenv("SCRECTKEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// check the exp

		if float64(time.Now().Unix()) > claims["exp"].(float64) {

			c.AbortWithStatus(http.StatusUnauthorized)

			c.JSON(http.StatusBadRequest, gin.H{

				"error": "user not found",
			})
		}

		// Find the user with token sub

		var user models.User

		database.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)

			c.JSON(http.StatusBadRequest, gin.H{

				"error": "user not found",
			})
		}

		// attach of req

		c.Set("user", user)

		// continue

		c.Next()

		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "user not found",
		})
	}

}
