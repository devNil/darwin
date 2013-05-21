package socket

import (
	ws "code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var colors = []int32{
							0xFF00FF,
							0xFFFF0F,
							0xAFAFAF,
							0x00AAEE}
var id = strconv.FormatInt(time.Now().Unix(), 10)[6:] //this is the id a remote uses to connect
const (
	BoardX    = 640 //Board width
	BoardY    = 480 //Board heigth
	BoardFS   = 16  //Boardfield size
	MaxX      = BoardX/BoardFS - 1
	MaxY      = BoardY/BoardFS - 1
	LenX      = BoardX / BoardFS
	LenY      = BoardY / BoardFS
	NumTicks  = 10
	NumBots   = 1
	NumPlayer = 1 + NumBots //number of players per game
)

type entity struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Dir   int8   `json:"dir"`
	Color string `json:"color"`
	S     int8   `json:"size"`
	died  bool   //flag if player died, so he won't be updated
}

type server struct {
	clients  map[*client]bool
	bots     map[*client]bool
	register chan *client
	tick     chan bool
	update   chan command
	idc      int8
	board    [LenX * LenY]int32
}

var gameserver = &server{
	clients:  make(map[*client]bool),
	bots:     make(map[*client]bool),
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
			if c.conn != nil{
			    s.clients[c] = true
				c.input <- command{-1, []byte("Welcome on this Server")}
				//send the remote id to a client
				c.input <- command{2, []byte(id)}
				val, _ := json.Marshal(c.e)
				c.input <- command{0, val}
			} else {
				s.bots[c] = true
			}
		case cmd := <-s.update:
			fmt.Println(cmd)
		case <-s.tick:
			tick += 1
			var result []*entity
			for b, _ := range s.bots {
				if tick%NumTicks == 0 && !b.e.died {
					tempX := b.e.X
					tempY := b.e.Y
					if b.e.Dir == 0 {
						tempX += BoardFS
					}
					if b.e.Dir == 1 {
						tempY += BoardFS
					}

					if b.e.Dir == 2 {
						tempX -= BoardFS
					}

					if b.e.Dir == 3 {
						tempY -= BoardFS
					}
					if s.board[tempX/BoardFS+(tempY/BoardFS*LenX)] != 0 {
						b.e.findNewDirection()
                        b.e.died = true
                        log.Println("bot has problems")
					} else {
                        b.e.X = tempX
                        b.e.Y = tempY
                        result = append(result, b.e)
                    }
				}
			}
			for k, _ := range s.clients {
				if tick%NumTicks == 0 && !k.e.died {
					if k.e.Dir == 0 {
						k.e.X += BoardFS
					}
					if k.e.Dir == 1 {
						k.e.Y += BoardFS
					}

					if k.e.Dir == 2 {
						k.e.X -= BoardFS
					}

					if k.e.Dir == 3 {
						k.e.Y -= BoardFS
					}
					fmt.Println(s.board[k.e.X/BoardFS+(k.e.Y/BoardFS*LenX)])
					if s.board[k.e.X/BoardFS+(k.e.Y/BoardFS*LenX)] != 0 {
						k.input <- command{3, []byte("died")}
						k.e.died = true
						//kill player and all connections
					} else {
						s.board[k.e.X/BoardFS+(k.e.Y/BoardFS*LenX)] = 1
						result = append(result, k.e)
					}
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
	for y := 0; y < LenY; y++ {
		gameserver.board[y*LenX] = 1
		gameserver.board[MaxX+y*LenX] = 1
	}
	for x := 0; x < LenX; x++ {
		gameserver.board[x] = 1
		gameserver.board[x+MaxY*LenX] = 1
	}
	rand.Seed(time.Now().UnixNano())
	go gameserver.run()
	startUp()
}
func startUp() {
	if len(gameserver.bots) != NumBots {
		log.Println(len(gameserver.bots))
		addBot()
	}
	if gameserver.idc == NumPlayer {
        log.Println("game should start")
		for k, _ := range gameserver.clients {
            if k.conn != nil {
			    k.input <- command{4, nil}
            }
		}
		time.Sleep(10 * time.Second)
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

		if cmd.Id == 1 {
			x := string(cmd.Value)
            log.Println("test 1 " + string(cmd.Value))
			if x == "1" {
                log.Println("1")
				if c.e.Dir == 3 {
					c.e.Dir = 0
				} else {
					c.e.Dir += 1
				}
			}

			if x == "-1" {
                log.Println("-1")
				if c.e.Dir == 0 {
					c.e.Dir = 3
				} else {
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
func addBot(){
    var x, y int
	for {
		x = rand.Intn(LenX) % BoardFS
		y = rand.Intn(LenY) % BoardFS
		if gameserver.board[x+y*LenX] == 0 {
			break
		}
	}
	cl := &client{
		nil,
		nil,
		&entity{
			X:     x * 16,
			Y:     y * 16,
			Dir:   int8(rand.Intn(4)),
			Color: fmt.Sprintf("#%X", colors[gameserver.idc]),
			S:     16,
		},
	}
    gameserver.idc += 1
    log.Println("new Bot")
    gameserver.register <- cl

}
//each new connection will first 
func ConnectionHandler(connection *ws.Conn) {
	var x, y int
	//search for a empty field
	for {
		x = rand.Intn(LenX) % BoardFS
		y = rand.Intn(LenY) % BoardFS
		if gameserver.board[x+y*LenX] == 0 {
			break
		}
	}
	cl := &client{
		connection,
		make(chan command),
		&entity{
			X:     x * 16,
			Y:     y * 16,
			Dir:   int8(rand.Intn(4)),
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
func RemoteConnectionHandler(connection *ws.Conn) {
	var m string
	var x, y int
	ws.Message.Receive(connection, &m)
	//check if the sended id is correct
	if m == id {
		//search for a empty field
		for {
			x = rand.Intn(LenX) % BoardFS
			y = rand.Intn(LenY) % BoardFS
			if gameserver.board[x+y*LenX] == 0 {
				break
			}
		}
		cl := &client{
			connection,
			make(chan command),
			&entity{
				X:     x * 16,
				Y:     y * 16,
				Dir:   int8(rand.Intn(4)),
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
func (e *entity) findNewDirection() {
	for {
		x := rand.Intn(2)
		if x == 1 {
			if e.Dir == 3 {
				e.Dir = 0
			} else {
				e.Dir += 1
			}
		}
		if x == 0 {
			if e.Dir == 0 {
				e.Dir = 3
			} else {
				e.Dir -= 1
			}
		}
	}
}
