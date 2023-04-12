package controller

import (
	"fmt"
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

var UserLogStatus = false

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
		return
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
	UserLogStatus = true
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
		return
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

		"message": "login sucess",
	})

	UserLogStatus = true
}

// user profile

func UserProfile(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	ok := UserLogStatus

	if ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read user",
		})
		return
	}

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{

		"username": user.(models.User).Name,
		"email":    user.(models.User).Email,
	})
}

// Edit User
func UserEdit(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	ok := UserLogStatus

	if ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read user",
		})
		return
	}

	user, _ := c.Get("user")

	var body struct {
		Name  string
		Email string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read new user",
		})
		return
	}

	userid := user.(models.User).ID
	// username := user.(models.User).Name
	// usermail := user.(models.User).Email
	var editQuery string

	if body.Email == "" && body.Name == "" {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "update with current detail",
		})

	} else if body.Email == "" {
		// for updating name only
		editQuery = fmt.Sprintf("UPDATE users SET name = '%s' WHERE id = %v", body.Name, userid)

		c.JSON(http.StatusBadRequest, gin.H{

			"message": "updating name......!!",
		})

	} else if body.Name == "" {
		// for updating email only

		editQuery = fmt.Sprintf("UPDATE users SET email = '%s' WHERE id = %v", body.Email, userid)

		c.JSON(http.StatusBadRequest, gin.H{

			"message": "updating Email......!!",
		})

	} else {
		// for updating both
		editQuery = fmt.Sprintf("UPDATE users SET email = '%s' , name = '%s' WHERE id = %v", body.Email, body.Name, userid)

		c.JSON(http.StatusBadRequest, gin.H{

			"message": "updating Name and Email......!!",
		})
	}

	fmt.Println(editQuery)

	// var user models.User

	database.DB.Raw(editQuery).Scan(&user)

	if body.Email == "" {
		// for updating name only

		c.JSON(http.StatusBadRequest, gin.H{

			"message":  "Name Updated",
			"new name": body.Name,
		})

	} else if body.Name == "" {
		// for updating email only

		c.JSON(http.StatusBadRequest, gin.H{

			"message":   "Email Updated",
			"new email": body.Email,
		})

	} else {
		// for updating both

		c.JSON(http.StatusBadRequest, gin.H{

			"message":   "Name and Email Updated",
			"New name":  body.Name,
			"New Email": body.Email,
		})
	}

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
