package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type exampleMelodyHandler struct {
	melo *melody.Melody
}

func createExampleMelodyHandler() exampleMelodyHandler {
	mel := exampleMelodyHandler{}
	m := melody.New()

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Printf("websocket connection open. [session: %#v]\n", s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Printf("websocket connection close. [session: %#v]\n", s)
	})

	mel.melo = m
	return mel
}

func chatFunc(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "html/chat.html")
}

func (e *exampleMelodyHandler) wsHandler(c *gin.Context) {
	e.melo.HandleRequest(c.Writer, c.Request)
}

func main() {
	r := gin.Default() //ginは基本的にgin.Default()の返す構造体のメソッド経由でないと操作できない．
	r.LoadHTMLGlob("html/*.html")

	emelody := createExampleMelodyHandler()

	v1 := r.Group("/")
	{
		v1.GET("chat", chatFunc)
		v1.GET("ws", emelody.wsHandler)
	}
	r.Run(":8080")
}
