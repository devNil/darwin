package web

import(
	"log"
	ws "code.google.com/p/go.net/websocket"
	"time"
	"strconv"
)

type darwinServer struct{
	browser map[string]*browser
}
var server = darwinServer{
	make(map[string]*browser),
}

type browser struct{
	id string
	input chan string
	conn *ws.Conn
}

func(b *browser)send(){
	defer b.conn.Close()
	for{
		var message string
		err := ws.Message.Receive(b.conn, &message)
		if err != nil{
			log.Println(err)
			break
		}
	}
	delete(server.browser, b.id)
	log.Println(len(server.browser))
}

func (b *browser)read(){
	for message := range b.input{
		err:= ws.Message.Send(b.conn, message)
		if err != nil{
			log.Println("Close of client")
			break
		}
	}
	b.conn.Close()
	delete(server.browser, b.id)
	log.Println(len(server.browser))
}

type mobile struct{
	b *browser
	conn *ws.Conn
}

func (m *mobile)send(){
	defer m.conn.Close()
	for{
		var message string
		err := ws.Message.Receive(m.conn, &message)
		if err != nil{
			log.Println(err)
			m.b.input<-"Disconnect of phone"
			break
		}
		m.b.input<-message
	}
}

//Desktopbrowser registers here
func BrowserSocketHandler(connection *ws.Conn){
	id := time.Now().Unix()
	b := browser{
		strconv.FormatInt(id,10),
		make(chan string),
		connection,
	}
	err := ws.Message.Send(connection, strconv.FormatInt(id, 10))
	if err != nil{
		log.Println("Registered and closed immediatly")
		return
	}
	server.browser[strconv.FormatInt(id, 10)] = &b
	go b.send()
	b.read()

}

func MobileSocketHandler(connection *ws.Conn){
	var m string

	ws.Message.Receive(connection, &m)

	if b := server.browser[m]; b != nil{

		mobile := mobile{
			b,
			connection,
		}
		mobile.send()
	}else{
		return
	}
}
