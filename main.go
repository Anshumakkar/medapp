package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"medapp/db"
)

type Response struct {
	Error   bool     `json:"error"`
	Message string   `json:"message"`
	User    UserType `json:"user"`
}

type UserType struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
}

type RegistrationData struct {
	Email       string `json:"email"  binding:"required"`
	Gender      string `json:"gender"  binding:"required"`
	Password    string `json:"password"  binding:"required"`
	Username    string `json:"username"  binding:"required"`
	PhoneNumber string `json:"phonenumber" binding:"required"`
}

type LoginData struct{
	PhoneNumber string `json:"phonenumber" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	err:=db.CreateClient()
	if err!=nil{
		log.Fatal("DB Client should be started")
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	//router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.tmpl.html", nil)
	})

	router.POST("/registerUser", func(c *gin.Context) {
		log.Println("url is", c.Request.URL.Query())
		values := c.Request.URL.Query()
		for key, value := range values {

			log.Printf("Key = %v value = %v\n", key, value)
		}

		var parsedData RegistrationData
		err := c.BindJSON(&parsedData)
		if err == nil {
			log.Println(parsedData)
		} else {
			log.Println("Error is ", err.Error())
			c.AbortWithStatusJSON(400, gin.H{"error": "Json body ill formatted"})
			return
		}
		uid := uuid.NewV4()
		uidstr := uid.String()
		phoneNumber:= parsedData.PhoneNumber
		username:=parsedData.Username
		email:=parsedData.Email
		gender:=parsedData.Gender
		password:=parsedData.Password
		err= db.StoreItem(phoneNumber,email,username,password,gender,uidstr)
		if err!=nil{
			log.Println("Error while inserting into DB " ,err.Error())
			c.JSON(500,gin.H{"error":true,"message":"Internal Service Error"})
			return
		}
		test := gin.H{"id": uidstr, "username": parsedData.Username, "email": parsedData.Email,
			"gender": parsedData.Gender, "phonenumber": parsedData.PhoneNumber}
		c.JSON(200, gin.H{"error": false, "message": "registration success", "user": test})
	})

	router.POST("/loginUser",func(c *gin.Context){
   		var parsedData LoginData
                err := c.BindJSON(&parsedData)
                if err == nil {
                        log.Println(parsedData)
                } else {
                        log.Println("Error is ", err.Error())
                        c.AbortWithStatusJSON(400, gin.H{"error": "Json body ill formatted"})
                        return
                }
		var userData *db.LoginData
		userData,err= db.GetItem(parsedData.PhoneNumber)
		if err!=nil{
		  log.Println("Error while retrieving data for ",parsedData.PhoneNumber)
		  log.Println(err)
		  c.JSON(500,gin.H{"error":true,"message":"Internal Service Error"})
		  return
		}
		if userData.Password != parsedData.Password{
		log.Println("Password required is ",userData.Password," and password given is ",parsedData.Password)
		c.JSON(403,gin.H{"error":true,"message":"Incorrect Credentials"})
		return
		}
		test := gin.H{"id": userData.ID, "username": userData.Username, "email": userData.Email,
                        "gender": userData.Gender, "phonenumber": userData.PhoneNumber}
                c.JSON(200, gin.H{"error": false, "message": "registration success", "user": test})

	})

	router.Run(":" + port)
}
