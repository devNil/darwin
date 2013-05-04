package socket

import(
	"fmt"
	"log"
    "time"
	ws "code.google.com/p/go.net/websocket"
	"time"
)

var colors = []int32{0xFF00FF, 0xFFFF00}

type entity struct{
	id,X,Y,Dir int8
	Color string
	S int8
}

type game struct{
	idc int8
	entities []*entity
	board [(640/16)*(480/8)]int32
}

type server struct{
	clients map[*client]bool
	register chan *client
	tick chan bool
	update chan command
}

var gameserver = &server{
	make(map[*client]bool),
	make(chan *client),
	make(chan bool),
	make(chan command),
}

func (s *server)run(){
    for{
        select{
            case c := <-s.register:
                s.clients[c]=true
            case cmd := <-s.update:
                fmt.Println(cmd)
            case <-tick:
                fmt.Println("tick")
        }
    }
}

func Run(){
    tick()
}
func tick() {
    fmt.Println("tick tack")
    time.AfterFunc(time.Second/60, tick)
}

type command struct{
	Id int8
	Value []byte
}

type client struct{
	conn *ws.Conn
	input chan command
}

func(c *client)send(){
	defer c.conn.Close()
	for{
		var cmd command
		err := ws.JSON.Receive(c.conn, &cmd)
		if err != nil{
			log.Println(err)
			break
		}
	}
}

func (c *client)read(){
	for cmd := range c.input{
		err := ws.JSON.Send(c.conn, cmd)
		if err != nil {
			log.Println("Close of client")
			break
		}
	}
	c.conn.Close()
}

func ConnectionHandler(connection *ws.Conn){
	cl := &client{
		connection,
		make(chan command),
	}

	go cl.send()
	cl.read()
}
