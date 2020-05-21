package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
"github.com/satori/go.uuid"
)

type Response struct {
	Error   bool `json:"error"`
	Message string `json:"message"`
	User    UserType `json:"user"`
}
type UserType struct {
	Id       int `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
}
type LoginData struct{
	Email string `json:"email"  binding:"required"`
	Gender string `json:"gender"  binding:"required"`
	Password string `json:"password"  binding:"required"`
	Username string `json:"username"  binding:"required"`
}

func main() {
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

                   var parsedData LoginData
		   err:=c.BindJSON(&parsedData) 
		   if err== nil {
		       log.Println(parsedData)
	       	   }else{ 
		   log.Println("Error is ",err.Error())
		   c.AbortWithStatusJSON(400,gin.H{"error":err.Error()})
		   return 
		   }
		   uid := uuid.NewV4()
		   uidstr := uid.String()
		test:=gin.H{"id":uidstr,"username":parsedData.Username,"email":parsedData.Email,"gender":parsedData.Gender}
		c.JSON(200,gin.H{"error":false,"message":"registration success","user":test})
//		c.JSON(200, resJson)
	})

	router.Run(":" + port)
}
