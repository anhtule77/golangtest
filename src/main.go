package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool) // khai bao  bien toan cuc dang map, key dang con tro, value kieu bool,
// su dung map vi de dang ket noi va xoa
var broadcast = make(chan Message) // khai bao bien la mot chan: hang doi cho cac tin nhan,
var upgrader = websocket.Upgrader{ //la mot doi tuong voi cac phuong thuc dung de lay ket noi http thong thuong va upgrade len websocket
	CheckOrigin: func(r *http.Request) bool { // goi ham CheckOrigin de kiem tra nguon goc
		return true
	},
}

type Message struct { // dinh nghia mot doi tuong voi cac thuoc tinh don gian
	Email    string `json:"email"` //chuyen doi tuong Message thanh json va nguoc lai
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	// tao mot static fileserver lien ket voi route"/" de nguoi dung co the
	// truy cap xem trang index
	http.Handle("/", fs)
	// thiet lap thuc thi fs khi co request den duong dan root /

	// thiet lap thuc thi ham handleConnection khi co request den duong dan root /ws
	http.HandleFunc("/ws", handleConnection)

	//sd goroutine goi "handleMessages", qua trinh nay se chay song song voi cac phan con lai cua ung dung
	//va chi nhan tin nhan tu channel tu truoc va chuyen cho khach hang qua ket noi Websocket
	go handleMessage()

	log.Println("http server start on: 8000")
	// dua ra thong bao va
	// khoi dong may chu tren cong localhost 8000 va ghi lai loi
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

// khi co request gui den, ham handleConnection se duoc goi va 2 parameter
//se duoc set tuong ung la 2 gia tri dai dien cho response va request
func handleConnection(w http.ResponseWriter, r *http.Request) { //ham xu ly cac ket noi den Websocket

	ws, err := upgrader.Upgrade(w, r, nil) // chuyen ket noi hien co sang giao thuc Websoket
	if err != nil {
		log.Fatal(err) // ghi lai loi neu co nhung khong dung ctr
	}
	defer ws.Close()
	//khi gap lenh defer ctr se bo qua, chua thuc thi cau lenh nay, sd de dong ket noi Websocket khi ham tra ve
	clients[ws] = true // tao mot may khach moi bang cach them no vao map"client" toan cuc
	for {              // tao vong for de cho doi message moi duoc ghi vao Websocket
		var msg Message
		//doc mot tin nhan moi duoi dang JSON va anh xa no toi mot doi tuong Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			// neu client bi ngat ket noi, ta se ghi lai va xoa khach hang ra khoi bien toan cuc
			// de chung ta khong doc hoac gui message moi cho client do
			log.Printf("error: %v", err)
			delete(clients, ws) // xoa phan tu trong map clients
			break
		}
		//gui tin nhan moi nhan duoc den channel
		broadcast <- msg
	}
}
func handleMessage() {
	for {
		//vong lap lien tuc doc tu channel broadcast
		//chuyen  tiep thong tin den tat ca cac may khach qua ket noi Websocket
		msg := <-broadcast
		// gui message den moi client duoc ket noi
		for client := range clients { //lap qua tat ca cac clients
			err := client.WriteJSON(msg)
			if err != nil {
				// neu client bi ngat ket noi, ta se ghi lai va xoa khach hang ra khoi bien toan cuc
				// de chung ta khong doc hoac gui message moi cho clien do
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client) // xoa mot phan tu trong map clients
			}
		}
	}
}
