package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	timeout "github.com/vearne/gin-timeout"
)

type errResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func MyRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				info := p.(timeout.PanicInfo)
				log.Printf("unknow error:%v\n", info.Value)
				log.Printf("stack:%v\n", info.Stack)
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					errResponse{Code: -1, Msg: fmt.Sprintf("unknow error:%v", p)})
				return
			}
		}()
		c.Next()
	}
}

func main() {
	router := gin.Default()
	// In order to maintain flexibility,
	// you should define your own recovery middleware
	router.Use(MyRecovery())
	defaultMsg := `{"code": -1, "msg":"http: Handler timeout"}`
	router.Use(timeout.Timeout(timeout.WithTimeout(3*time.Second),
		timeout.WithDefaultMsg(defaultMsg)))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, time.Now().String())
	})
	router.GET("/panic", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		dive()
	})
	log.Fatal(router.Run(":8080"))
}

func dive() {
	x := 0
	fmt.Println(100 / x)
}
