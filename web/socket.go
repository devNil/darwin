package web

import(
	"log"
	ws "code.google.com/p/go.net/websocket"
	"time"
	"strconv"
)

type client struct{
	browser *ws.Conn
	phone *ws.Conn
	send chan string
}

func (c *client)close(){
	c.browser.Close()
	c.phone.Close()
}



//read phone inputs
func (c *client)read(){
	for{
		var message string
		err := ws.Message.Receive(c.phone, &message)
		if err != nil{
			log.Println("Error read: ")
			break
		}

	}
	c.close()
}

//write to browser
func (c *client)write(){
	for message := range c.send{
		err := ws.Message.Send(c.browser, &message)
		if err != nil{
			log.Println("Error write: ")
			break
		}
	}
	c.close()
}

type darwinServer struct{
	browser map[string]*browser
}
var server = darwinServer{
	make(map[string]*browser),
}

type browser struct{
	input chan string
	conn *ws.Conn
}

func (b *browser)read(){
	for message := range b.input{
		err:= ws.Message.Send(b.conn, message)
		if err != nil{
			log.Println(err)
			break
		}
	}
	b.conn.Close()
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
		}
		m.b.input<-message
	}
}

//Desktopbrowser registers here
func BrowserSocketHandler(connection *ws.Conn){
	id := time.Now().Unix()
	b := browser{
		make(chan string),
		connection,
	}
	server.browser[strconv.FormatInt(id, 10)] = &b
	ws.Message.Send(connection, strconv.FormatInt(id, 10))
	b.read()

}

func MobileSocketHandler(connection *ws.Conn){
	var m string

	ws.Message.Receive(connection, &m)
	
	mobile := mobile{
		server.browser[m],
		connection,
	}
	mobile.send()
}
