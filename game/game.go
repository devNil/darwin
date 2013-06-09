//The game package provides basic game entites
package game

import "sync"

type Direction int8

const (
	UP Direction = iota
	LEFT
	DOWN
	RIGHT
)

//Basic Entity struct
type Entity struct {
	X     int64
	Y     int64
	Dir   Direction //Direction
	Color int32     //Color of the player
	Died  bool      //flag if dead
	lock  *sync.Mutex
}

//The game struct is a container for all things in the game
type Game struct {
	Entities map[*Entity]bool //Map with every entity in game
	Running  bool             //Flag if game is running
	Map      []bool
	W        int64 //Width of the map
	H        int64 //Height of the map
	Ps       int64 //Player size
}

func (d Direction) clockwise() Direction {
	if d == RIGHT {
		return DOWN
	}

	if d == DOWN {
		return LEFT
	}

	if d == LEFT {
		return UP
	}

	return RIGHT
}

func (d Direction) counterclockwise() Direction {
	if d == RIGHT {
		return UP
	}

	if d == DOWN {
		return RIGHT
	}

	if d == LEFT {
		return DOWN
	}

	return LEFT
}

//only accepts RIGTH or LEFT
func (e *Entity) SetDir(d Direction) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if d == RIGHT {
		e.Dir = e.Dir.clockwise()
	}

	if d == LEFT {
		e.Dir = e.Dir.counterclockwise()
	}

}

func (g *Game) DeleteEntity(e *Entity) {
	delete(g.Entities, e)
}

//If game is finished
func (g *Game) End() bool {
	end := true

	count := 0

	for e := range g.Entities {
		if !e.Died {
			end = false
			count++
		}
	}

	return end || count == (len(g.Entities)-(len(g.Entities)-1))
}

//Update game
func (g *Game) Update() {
	for e, _ := range g.Entities {
		if e.Died {
			continue
		}

		if e.Dir == UP {
			e.Y -= g.Ps
		}

		if e.Dir == RIGHT {
			e.X += g.Ps
		}

		if e.Dir == DOWN {
			e.Y += g.Ps
		}

		if e.Dir == LEFT {
			e.X -= g.Ps
		}

		if g.Collide(e) {
			e.Died = true
			continue
		}

		g.SetPosition(e)
	}
}

//Set the position on the map, not safe!!!
func (g *Game) SetPosition(e *Entity) {
	g.Map[(e.X/g.Ps)+(e.Y/g.Ps)*g.W] = true
}

//Checks if the point of the map is already setted
func (g *Game) Collide(e *Entity) bool {
	if e.X < 0 || e.X >= (g.W*g.Ps) || e.Y < 0 || e.Y >= (g.H*g.Ps) {
		return true
	}

	return g.Map[(e.X/g.Ps)+(e.Y/g.Ps)*g.W]
}

//Reset map and entities, entities will be marked as alive
func (g *Game) Reset() {
	for i := 0; i < len(g.Map); i++ {
		g.Map[i] = false
	}

	for e, _ := range g.Entities {
		e.Died = false
	}

}

//Create a new Entity
func (g *Game) NewEntity(x, y int64, dir Direction, color int32) *Entity {
	e := new(Entity)
	e.X = x
	e.Y = y
	e.Dir = dir
	e.Color = color
	e.Died = false
	e.lock = new(sync.Mutex)

	g.Entities[e] = true
	g.SetPosition(e)
	return e
}

//Create new game with desired map width, height and size of players
func NewGame(width, height int64, psize int64) *Game {
	g := new(Game)
	g.Entities = make(map[*Entity]bool)
	g.Running = false
	g.Map = make([]bool, (width/psize)*(height/psize))
	g.W = width / psize
	g.H = height / psize
	g.Ps = psize
	return g
}
