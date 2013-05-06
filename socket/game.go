package socket

import (
	ws "code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"log"
	"time"
    "strconv"
    "math/rand"
)

var colors = []int32{
    0xFF00FF,
    0xFFFF0F,
    0xAFAFAF,
    0x00AAEE
    }
var id = strconv.FormatInt(time.Now().Unix(),10) //this is the id a remote uses to connect
const (
    BoardX = 640 //Board width
    BoardY = 480 //Board heigth
    BoardFS = 16 //Boardfield size
    MaxX = BoardX/BoardFS -1
    MaxY = BoardY/BoardFS -1
    LenX = BoardX/BoardFS
    LenY = BoardY/BoardFS
    NumTicks = 10
    NumPlayer = 2 //number of players per game
)
type entity struct {
	X     int   `json:"x"`
	Y     int   `json:"y"`
	Dir   int8   `json:"dir"`
	Color string `json:"color"`
	S     int8   `json:"size"`
}


type server struct {
	clients  map[*client]bool
	register chan *client
	tick     chan bool
	update   chan command
	idc      int8
	board    [LenX * LenY]int32
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
			c.input <- command{-1, []byte("Welcome on this Server")}
            //send the remote id to a client
            c.input <- command{2, []byte(id)}
			val, _ := json.Marshal(c.e)
			c.input <- command{0, val}

		case cmd := <-s.update:
			fmt.Println(cmd)
		case <-s.tick:
			tick += 1
			var result []*entity
			for k, _ := range s.clients {
				if tick%NumTicks == 0 {
					if k.e.Dir == 0 {
						k.e.X += BoardFS
					}
                    if k.e.Dir == 1 {
                       k.e.Y += BoardFS
                    }

                    if k.e.Dir == 2{
                        k.e.X -= BoardFS
                    }

                    if k.e.Dir == 3{
                        k.e.Y -= BoardFS
                    }
                fmt.Println(s.board[k.e.X/BoardFS+(k.e.Y/BoardFS*LenX)])
                if (s.board[k.e.X/BoardFS+(k.e.Y/BoardFS*LenX)] != 0){
                    k.input <- command{3,[]byte("died")}
                    //kill player and all connections
                } else {
                    s.board[k.e.X/BoardFS+(k.e.Y/BoardFS*LenX)] = 1
                }
				result = append(result, k.e)
				}
			}
            if tick%NumTicks == 0 {
			    b, _ := json.Marshal(result)
			    for k, _ := range s.clients {
		            k.input <- command{1, b}
			    }
            }
		}
	}
}

func Run() {
    //add the walls to the map
    for y:=0; y < 480/16; y++{
        gameserver.board[y*LenX] = 1
        gameserver.board[MaxX+y*LenX] = 1
    }
    for x:=0; x < 640/16; x++{
        gameserver.board[x] = 1
        gameserver.board[x+MaxY*LenX] = 1
    }
    rand.Seed(time.Now().UnixNano())
	go gameserver.run()
	startUp()
}
func startUp() {
    if gameserver.idc == NumPlayer {
        for k, _ := range gameserver.clients {
            k.input <- command{4,nil}
        }
        time.Sleep(10*time.Second)
        tick()
    } else {
	    time.AfterFunc(time.Second/60, startUp)
    }
}
func tick() {
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
//each new connection will first 
func ConnectionHandler(connection *ws.Conn) {
    var x,y int
    //search for a empty field
    for {
        x = rand.Intn(LenX)%BoardFS*BoardFS
        y = rand.Intn(LenY)%BoardFS*BoardFS
        if (gameserver.board[x+y*LenX] == 0){
            break;
        }
    }
	cl := &client{
		connection,
		make(chan command),
		&entity{
			X:     x,
			Y:     y,
			Dir:   rand.Intn(4),
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
func RemoteConnectionHandler(connection *ws.Conn){
    var m string
    var x,y int
    ws.Message.Receive(connection, &m)
    //check if the sended id is correct
    if (m == id){
        //search for a empty field
        for {
            x = rand.Intn(LenX)%BoardFS*BoardFS
            y = rand.Intn(LenY)%BoardFS*BoardFS
            if (gameserver.board[x+y*LenX] == 0){
                break;
            }
        }
        cl := &client{
            connection,
            make(chan command),
            &entity{
                X:     x,
                Y:     y,
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
}
