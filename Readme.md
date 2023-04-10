
WEB-APP
 
Simple web application written in golang using postgres database


Package used for this WebApp

`GORM`        | `https://gorm.io/`                                 | For Database opration

`GIN`         | `https://gin-gonic.com/`                           | For HTTP web framework

`CRYPTO`      | `https://pkg.go.dev/golang.org/x/crypto@v0.8.0`    | Password encryption and decryption

`JWT`         | `https://pkg.go.dev/github.com/golang-jwt/jwt`     | For implementation of JSON Web Tokens.

`GODOTENV`    | `https://github.com/joho/godotenv`                 | For loads env vars from a .env file

## Token for JWT Authentication and also using if the user or admin is alredy logged or not through this

* User token is saved as USERTOKEN

* Admin token is seved as ADMINTOKEN

## USER SPECIFICATION

LogIn

SignUp

Edit Profile

User Profile Viewer

## ADMIN SPECIFICATION

LogIn

SignUp

Add Admi

List All User

Edit User

Delete User

# Run this Code

    1. Run your server eg: local server with postgres database

    2. Give the database details to .env file

    3. go run main.go

    * Use this route for working eg: localhost://8080/login

# ------------Routes And HTTP Method-----------

> Admin route

    POST     /adminlogin
	POST     /addadmin
	GET      /userview
	GET      /deleteuser
	PUT      /edituser

> User route

	POST    /signup
	POST    /login
	GET     /userprofile
    PUT     /useredit



# .env Details

* customize with your detail
    > user      = "your postgres username"
    > password  = "password of your postgres user"
    > database  = "your server name"
    > localhost = "your host"

    > PORT      = 8080
    > DATABASE  = "host=localhost user="user" password="password" dbname="database" port="port"             sslmode=disable"

    // User key
    // You can edit the key

    > SCRECTKEY = AHDBHSIHSNNSGHSHSJHJ

    > ADMINKEY  = KAHNDKJDHBMKJSNND


# Simple CURD Oprations and the paths

# Signup 

1. Receive Name, email , password from request
2. Check if user is already registered. (Check with email)
3. If user if not present, hash the password
4. Create new user in the db with email and hashed password
5. Send success response back

# Login

1. Recieve email and password from request
2. Find user info from database
3. If user is not found (user.ID == 0), return invalid username or password response
4. Compare password from request and database
5. If password is incorrect, return invalid username or password response
6. If password is correct, create token, set in in cookie

# Delete

1. Recieve the id from request
2. check the id if exist or not
3. find the user from database
4. delete and update the database
5. send sucess responce

# Update

1. Collect the user information
2. Check the use id exist or not
3. Find the user from databse
4. Update the user information
5. Update the database 
6. Sent responce


---------------------------------------------------------------------------------
Upcomming Function

 * Admin can block user
 * When Admin Delete user only update database Delete_status value
 * User can delete Profile
 * When user delete profile only update database Delete_Status value
 * Google Signup
 * Twitter signup
 * Facebook signup

Upcomming Functionality

 * Template Loading and Front End Design