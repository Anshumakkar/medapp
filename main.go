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
type DoctorInfo struct{
	Name string `json:"name"`
	Slots []string `json:"slots"`
}

type DoctorsInfo struct{
Info []DoctorInfo `json:"doctorsInfo"`
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


	router.GET("/doctorsInfo",func(c *gin.Context){
		var doctorsArray []DoctorInfo
		slots := []string{"10:00 AM","10:15 AM","10:30 AM"}
		slots1 := []string{"11:00 AM","11:15 AM","11:30 AM"}
		docInfo1 := DoctorInfo{
			Name:"Doctor A",
			Slots:slots,
		}
		docInfo2 := DoctorInfo{Name:"Doctor B",Slots:slots}
		docInfo3 := DoctorInfo{Name:"Doctor C",Slots:slots1}
		docInfo4 := DoctorInfo{Name:"Doctor D",Slots: slots}
		doctorsArray = append(doctorsArray,docInfo1,docInfo2,docInfo3,docInfo4)
		info := DoctorsInfo{Info:doctorsArray}
		c.JSON(200,info)

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

		isExisting,err := db.CheckExistence(phoneNumber)
		if err!=nil{
			c.JSON(503,gin.H{"error":true,"message":"Internal Service Error"})
			return
		}else if isExisting == true{
			c.JSON(409,gin.H{"error":true,"message":"User already exists"})
			return
		}

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
		if err == nil && userData.PhoneNumber == "" {
		   log.Println("No User Data with this ",parsedData.PhoneNumber)
		   c.JSON(404,gin.H{"error":true,"message":"Not Registered"})
		   return
		}
		if userData.Password != parsedData.Password{
		  log.Println("Password required is ",userData.Password," and password given is ",parsedData.Password)
		  c.JSON(403,gin.H{"error":true,"message":"PhoneNumber or password is incorrect"})
		  return
		}
		test := gin.H{"id": userData.ID, "username": userData.Username, "email": userData.Email,
                        "gender": userData.Gender, "phonenumber": userData.PhoneNumber}
                c.JSON(200, gin.H{"error": false, "message": "login success", "user": test})

	})

	router.Run(":" + port)
}
