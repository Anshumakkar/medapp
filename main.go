package main

import (
	"log"
	"fmt"
	"net/http"
	"os"
"encoding/json"
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
//		values := c.Request.URL.Query()
//		for key, value := range values {
//
//			log.Printf("Key = %v value = %v\n", key, value)
//		}
		resp := &Response{Error: false,
			Message: "This is Registered",
			User: UserType{
				Id:       1,
				Username: "ANshu",
				Email:    "hello@test.com",
				Gender:   "Male",
			},
		}
		log.Println(resp)
		resJson,err :=json.Marshal(resp)
		if err!=nil{
		log.Println("Error is " + err.Error())
		}
		fmt.Fprintf(os.Stdout, "json resp : %s",resJson)
		c.JSON(200, string(resJson))
	})

	router.Run(":" + port)
}
