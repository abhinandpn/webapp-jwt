package controller

import (
	"jwt/pkg/database"
	"jwt/pkg/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

/*
Sign-Up 		function
Log-in 			function
Log-Out 		function
Edit-Profile 	function
*/

func SignUpUser(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// Get name,email and pass from request

	var body struct {
		Name     string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "failed to read new user",
		})
	}

	var user models.User

	// if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
	// }

	findUserQuery := `
					SELECT * FROM users
					WHERE email = $1 LIMIT 1;`

	database.DB.Raw(findUserQuery, body.Email).Scan(&user)

	if user.ID != 0 {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "email id already exists. Try again with another email",
		})

		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 3)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "failed to Hash password",
		})

		return

	}

	// create the user

	user = models.User{Name: body.Name, Email: body.Email, Password: string(hash)}

	result := database.DB.Create(&user)

	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": result.Error.Error(),
		})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{
		"message": "sucess to create user",
	})
}

func LoginUser(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// Get the name email and passs from current user
	var body struct {
		Name     string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "failed to read new user",
		})
	}

	// Lookup requested user

	var user models.User

	database.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "invalid email or password",
		})

		return
	}

	// Compare sent in pass with saved user with pass

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "invalid email or password",
		})

		return
	}

	// generate JWT token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{

		"sub": user.ID,

		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret

	tokenString, err := token.SignedString([]byte(os.Getenv("SCRECTKEY")))

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "invalid to create jwt token",
		})

		return
	}

	// Sent it back

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("UserToken", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{

		"token": tokenString,
	})
}

// validation

func Validation(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{

		"message": user,
	})
}

// Logout

func UserLogout(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("UserToken", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{

		"message": "User Logout succes",
	})

}
