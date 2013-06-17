package socket

import (
	ws "code.google.com/p/go.net/websocket"
	"darwin/game"
	"log"
	"math/rand"
	"time"
)

type coordinate struct {
	dir game.Direction
	x   int64
	y   int64
}

var startcoords = []coordinate{
	coordinate{
		game.RIGHT,
		8,
		8,
	},
	coordinate{
		game.DOWN,
		640 - 16,
		8,
	},
	coordinate{
		game.LEFT,
		640 - 16,
		480 - 16,
	},
	coordinate{
		game.UP,
		8,
		480 - 16,
	},
}

//command struct, the protocol
type command struct {
	Id    int8
	Value interface{}
}

//the client represents either a desktop machine or a mobile user
type client struct {
	server   *server
	Entity   *game.Entity //the entity in the game
	conn     *ws.Conn     //the websocket connection
	Approved bool         //if the client has approved to start
}

//the server is the game running as a websocket-server
type server struct {
	game      *game.Game       //the game-instance
	tickCount int64            //game ticks
	clients   map[*client]bool //map with clients
	//channels
	approve    chan *client
	register   chan *client //channel for registering clients
	unregister chan *client //channel for unregistering clients
	tick       chan bool    //channel for increment tickCount
	cd         int8         //countdown
	color      []int32
}

//for getting inputs
func (c *client) readConnection() {
	defer c.conn.Close()
	for {
		var cmd command
		err := ws.JSON.Receive(c.conn, &cmd)
		if err != nil {
			log.Println(err)
			break
		}

		if cmd.Id == 1 {
			c.server.approve <- c
		}

		if cmd.Id == 10 {
			c.Entity.SetDir(game.LEFT)
		}

		if cmd.Id == 11 {
			c.Entity.SetDir(game.RIGHT)
		}
	}

	c.server.unregister <- c
}

func (s *server) run() {
	for {
		select {
		case c := <-s.register:
			s.clients[c] = true
			err := ws.JSON.Send(c.conn, command{0, "PING"})
			ws.JSON.Send(c.conn, command{6, c.Entity})

			s.updateLobby()

			if err != nil {
				log.Println(err)
				s.unregister <- c
			}
		case c := <-s.unregister:
			s.game.DeleteEntity(c.Entity)
			c.conn.Close()
			delete(s.clients, c)
			s.updateLobby()
		case <-s.tick:
			if s.game.Running {
				s.tickCount++
				if s.tickCount%20 == 0 {
					s.game.Update()

					ent := make([]*game.Entity, 0)

					for e, _ := range s.game.Entities {
						ent = append(ent, e)
					}

					for c, _ := range s.clients {
						if s.game.End() {
							err := ws.JSON.Send(c.conn,
								command{5, ent})
							if err != nil {
								log.Println(err)
								s.unregister <- c
							}
							continue
						}
						err := ws.JSON.Send(c.conn,
							command{4, ent})
						if err != nil {
							log.Println(err)
							s.unregister <- c
						}
					}
				}

				if s.game.End() {
					s.game.Running = false
					s.game.Reset()
					s.Reset()
				}
			}
		case c := <-s.approve:
			c.Approved = true

			ready := true

			for k, _ := range s.clients {
				if !k.Approved {
					ready = false
				}
			}

			s.updateLobby()

			if ready {
				for c, _ := range s.clients {
					ws.JSON.Send(c.conn, command{6, c.Entity})
				}
				time.AfterFunc(time.Second, s.countdown)
			}

		}
	}
}

func (s *server) updateLobby() {
	clients := make([]*client, 0)

	for c, _ := range s.clients {
		clients = append(clients, c)
	}

	for c, _ := range s.clients {
		ws.JSON.Send(c.conn, command{20, clients})
	}
}

func (s *server) ticker() {
	s.tick <- true
	if s.game.Running {
		time.AfterFunc(time.Second/60, s.ticker)
	}
}

func (s *server) countdown() {
	s.cd--
	if s.cd == 0 {
		s.cd = 10
		for c, _ := range s.clients {
			ws.JSON.Send(c.conn, command{3, ""})
		}
		s.game.Running = true
		time.AfterFunc(time.Second/60, s.ticker)
		return
	}

	for c, _ := range s.clients {
		ws.JSON.Send(c.conn, command{2, s.cd})
	}

	time.AfterFunc(time.Second, s.countdown)
}

func (s *server) Start() {
	go s.run()
}

func (s *server) Reset() {
	for c, _ := range s.clients {
		s.game.DeleteEntity(c.Entity)
		c.conn.Close()
		delete(s.clients, c)
	}
}

//Handler for websocket purpose
func (s *server) Handler(connection *ws.Conn) {
	if len(s.clients) == 4 {
		connection.Close()
		return
	}
	client := new(client)
	client.conn = connection
	client.server = s
	client.Approved = false

	sc := startcoords[len(s.clients)]
	client.Entity = s.game.NewEntity(sc.x, sc.y, sc.dir, s.color[len(s.clients)])
	s.register <- client
	s.cd = 10
	client.readConnection()
}

//create new server
func NewServer() *server {
	rand.Seed(time.Now().UTC().UnixNano())
	s := new(server)
	s.game = game.NewGame(640, 480, 8)
	s.clients = make(map[*client]bool)
	s.register = make(chan *client)
	s.unregister = make(chan *client)
	s.approve = make(chan *client)
	s.tick = make(chan bool)
	s.color = []int32{
		0xFF7800,
		0x128E9B,
		0xD8F800,
		0x9303A7,
	}

	return s
}
