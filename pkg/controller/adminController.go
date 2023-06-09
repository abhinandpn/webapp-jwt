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
Login 	function
Logout 	function

User 	View
Delete 	User
Edit 	user
*/

var AdminLogStatus = false

// Admin Login   Data collect from database
func AdminLogin(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// Get the name email and passs from current admin
	var body struct {
		Name     string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "failed to Admin",
		})
	}

	// Lookup requested admin

	var admin models.Admin

	database.DB.First(&admin, "email = ?", body.Email)

	if admin.ID == 0 {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "invalid email or password 1",
		})

		return
	}

	// Compare sent in pass with saved admin with pass

	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.Password))

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "invalid email or password 2",
		})

		return
	}

	// generate JWT token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{

		"sub": admin.ID,

		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret

	tokenString, err := token.SignedString([]byte(os.Getenv("ADMINKEY")))

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "invalid to create jwt token",
		})

		return
	}

	// Sent it back

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("AdminToken", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{

		// "token":   tokenString,
		"message": "admin login sucess",
	})

	AdminLogStatus = true
}

/*
// Admin Login  Data collect from .env
func AdminLogIn(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	var body struct {
		// Name     string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "failed to read Admin",
		})
	}

	AdminEmail := os.Getenv("ADMINMAIL")        // get Admin mail
	AdminPassword := os.Getenv("ADMINPASSWORD") // get admin pass
	AdminId := os.Getenv("ADMINID")

	if body.Email == AdminEmail {

		if body.Password == AdminPassword {

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{

				"sub": AdminId,

				"exp": time.Now().Add(time.Hour * 24).Unix(),
			})

			// Sign and get the complete encoded token as a string using the secret

			tokenString, err := token.SignedString([]byte(os.Getenv("ADMINKEY")))

			if err != nil {

				c.JSON(http.StatusBadRequest, gin.H{

					"error": "invalid to create jwt token",
				})

				return
			}
			c.SetSameSite(http.SameSiteLaxMode)

			c.SetCookie("AdminToken", tokenString, 3600*24*30, "", "", false, true)

			c.JSON(http.StatusOK, gin.H{

				// "token":   tokenString,
				"message": "admin login sucess",
			})

			AdminLogStatus = true
		} else {

			c.JSON(http.StatusBadRequest, gin.H{

				"error": "invalid email or password",
			})

			return
		}
	} else {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "invalid email or password",
		})

		return
	}
}
*/

// Admin Logout
func AdminLogout(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("AdminToken", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{

		"message": "Admin Logout succes",
	})
}

func AddAdmin(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	ok := AdminLogStatus

	if ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read Admin",
		})
		return
	}

	// Get name,email and pass from request

	var body struct {
		Name     string
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "failed to read new admin",
		})
	}

	var admin models.Admin

	findAdminQuery := `
					SELECT * FROM admins
					WHERE email = $1 LIMIT 1;`

	database.DB.Raw(findAdminQuery, body.Email).Scan(&admin)

	if admin.ID != 0 {

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

	// create the Admin

	admin = models.Admin{Name: body.Name, Email: body.Email, Password: string(hash)}

	result := database.DB.Create(&admin)

	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error":   result.Error.Error(),
			"message": "error !!!",
		})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{

		"message": "sucess new admin is created",
	})

}

func UserView(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	ok := AdminLogStatus

	if ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read Admin",
		})
		return
	}

	UserViewQuery := `
	select id,name,email from users order by id asc;
	`

	var user []models.User

	database.DB.Raw(UserViewQuery).Scan(&user)

	c.JSON(http.StatusOK, user)

}

func UserDelete(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "failed to read request",
		})

		return
	}

	// Check if the user exists
	var user models.User

	if err := database.DB.Where("id = ?", req.ID).First(&user).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{

			"error": "user not found ",
		})

		return
	}

	// Delete the user
	if err := database.DB.Delete(&user).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{

			"error": "failed to delete user",
		})

		return

	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{

		"message": "user deleted successfully",
	})

}

func EditUser(c *gin.Context) {

	ok := AdminLogStatus

	if ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read Admin",
		})
		return
	}

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// Parse the JSON request body
	var updateUser models.User
	if err := c.BindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	// Check if the user exists
	var user models.User
	if err := database.DB.Where("id = ?", updateUser.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	// Update the user in the database
	if err := database.DB.Model(&user).Updates(updateUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update user",
		})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
	})

}

func BlockUser(c *gin.Context) {

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	var UserId struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&UserId); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{

			"error": "failed to read request",
		})

		return
	}

	// Check if the user exists

	var user models.User

	if err := database.DB.Where("id = ?", UserId.ID).First(&user).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{

			"error": "user not found ",
		})

		return
	}

	if user.BlockStatus {
		unblock := `update users set block_status=$1 where id=$2 `
		if err := database.DB.Exec(unblock, false, user.ID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"allert": "user is un-blocked",
		})
	} else {
		block := `update users set block_status=$1 where id=$2 `
		if err := database.DB.Exec(block, true, user.ID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"allert": "user is blocked",
		})
	}

	c.JSON(http.StatusOK, gin.H{

		"user id ":          user.ID,
		"user Name":         user.Name,
		"user Block Status": user.BlockStatus,
	})

}
