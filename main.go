package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

func root(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "index.html")
}

func pingPong(w http.ResponseWriter, r *http.Request) {
	u, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Println(err)
		return
	}
	defer u.Close()

	go func() {
		for {
			err = u.WriteMessage(websocket.TextMessage, []byte("ping from server"))
			if err != nil{
				log.Println("send ping: ", err)
				break
			}

			<-time.After(time.Second * time.Duration(5))
		}
	}()

	for{
		mt, msg, err := u.ReadMessage()
		if err != nil{
			log.Println(err)
			break
		}

		err = u.WriteMessage(mt, msg)
		if err != nil{
			log.Println(err)
			break
		}
	}
}

func main(){

	http.HandleFunc("/", root)
	http.HandleFunc("/ping", pingPong)

	http.ListenAndServe(":8080", nil)
}
