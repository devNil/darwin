package socket

import (
	ws "code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"log"
	"time"
    "math/rand"
)

var colors = []int32{0xFF00FF, 0xFFFF0F}

type entity struct {
	X     int   `json:"x"`
	Y     int   `json:"y"`
	Dir   int8   `json:"dir"`
	Color string `json:"color"`
	S     int8   `json:"size"`
}

type game struct {
	idc      int8
	entities []*entity
	board    [(640 / 16) * (480 / 8)]int32
}

type server struct {
	clients  map[*client]bool
	register chan *client
	tick     chan bool
	update   chan command
	idc      int8
	board    [(640 / 16) * (480 / 16)]int32
}

var gameserver = &server{
	clients:  make(map[*client]bool),
	register: make(chan *client),
	tick:     make(chan bool),
	update:   make(chan command),
	idc:      0,
}

func (s *server) run() {
	var tick int64
	for {
		select {
		case c := <-s.register:
			s.clients[c] = true
			c.input <- command{-1, []byte("Figg di")}

			val, _ := json.Marshal(c.e)
			c.input <- command{0, val}

		case cmd := <-s.update:
			fmt.Println(cmd)
		case <-s.tick:
			tick += 1
			var result []*entity
			for k, _ := range s.clients {
				if tick%10 == 0 {
					if k.e.Dir == 0 {
						k.e.X += 16
                        fmt.Println(k.e.X)
					}
                    if k.e.Dir == 1 {
                       k.e.Y += 16 
                    }

                    if k.e.Dir == 2{
                        k.e.X -= 16
                    }

                    if k.e.Dir == 3{
                        k.e.Y -= 16
                    }
                fmt.Println(s.board[k.e.X/16+(k.e.Y/16*640/16)])
                if (s.board[k.e.X/16+(k.e.Y/16*640/16)] != 0){
                    k.input <- command{3,[]byte("homo")}
                }
                s.board[k.e.X/16+(k.e.Y/16*640/16)] = 1
                //s.board[k.e.X/16+(k.e.Y/16*640/16)] = 1
				result = append(result, k.e)
				}
			}
			b, _ := json.Marshal(result)
			for k, _ := range s.clients {
				k.input <- command{1, b}
			}
			//fmt.Println("tick")
		}
	}
}

func Run() {
	go gameserver.run()
	tick()
}
func tick() {
	//fmt.Println("tick tack")
	gameserver.tick <- true
	time.AfterFunc(time.Second/60, tick)
}

type command struct {
	Id    int8   `json:"id"`
	Value []byte `json:"v"`
}

type client struct {
	conn  *ws.Conn
	input chan command
	e     *entity
}

func (c *client) send() {
	defer c.conn.Close()
	for {
		var cmd command
		err := ws.JSON.Receive(c.conn, &cmd)
		if err != nil {
			log.Println(err)
			break
		}

        if cmd.Id == 1{
            x := string(cmd.Value)
            if x == "1"{
                if c.e.Dir == 3{
                    c.e.Dir = 0 
                }else{
                    c.e.Dir += 1
                }
            }

            if x == "-1"{
                if c.e.Dir == 0{
                    c.e.Dir = 3
                }else{
                    c.e.Dir -= 1
                }
            }
        }

		gameserver.update <- cmd
	}
}

func (c *client) read() {
	for cmd := range c.input {
		err := ws.JSON.Send(c.conn, cmd)
		if err != nil {
			log.Println(err)
			break
		}
	}
	c.conn.Close()
}

func ConnectionHandler(connection *ws.Conn) {
	cl := &client{
		connection,
		make(chan command),
		&entity{
			X:     rand.Intn(640/16)%16*16,
			Y:     rand.Intn(480/16)%16*16,
			Dir:   0,
			Color: fmt.Sprintf("#%X", colors[gameserver.idc]),
			S:     16,
		},
	}
	gameserver.idc += 1
	log.Println("New Connection")
	gameserver.register <- cl
	go cl.send()
	cl.read()
}
