package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
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

                   var parsedData map[string]interface{}
		   err:=c.BindJSON(&parsedData) 
		   if err== nil {
		       log.Println(parsedData)
	       	   }else{ 
		   log.Println("Error is ",err.Error())
		   }

		test:=gin.H{"id":1,"username":"anshu","email":"testing@test.com","gender":"male"}
		c.JSON(200,gin.H{"error":false,"message":"testing","user":test})
//		c.JSON(200, resJson)
	})

	router.Run(":" + port)
}
