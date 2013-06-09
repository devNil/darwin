package game

import (
	"testing"
)

func TestCreateGame(t *testing.T) {
	game := NewGame(640, 480, 8)

	if game.W != 640/8 || game.H != 480/8 || game.Ps != 8 {
		t.Log("game init is false")
		t.Fail()
	}

	if len(game.Map) != (640/8)*(480/8) {
		t.Log("map is not big enough")
		t.Fail()
	}

	if game.Running {
		t.Log("new games must be paused")
		t.Fail()
	}
}

func TestCreateEntity(t *testing.T) {
	game := NewGame(640, 480, 8)
	game.NewEntity(0, 0, RIGHT, 0xFF00FF)

	if len(game.Entities) != 1 {
		t.Log("entity is not added")
		t.Fail()
	}
}

func TestCollide(t *testing.T) {
	game := NewGame(640, 480, 8)

	e1 := game.NewEntity(0, 0, RIGHT, 0xFF00FF)
	e2 := game.NewEntity(16, 0, LEFT, 0xF0F0F0)

	game.Update()

	if !e2.Died {
		t.Log(e2)
		t.Fail()
	}

	//e1 should win because of first place
	if e1.Died {
		t.Log(e1)
		t.Fail()
	}
}

func TestMove(t *testing.T) {
	game := NewGame(640, 480, 8)

	e1 := game.NewEntity(0, 0, RIGHT, 0xFF00FF)

	game.Update()

	if e1.X != 8 && e1.Y != 0 {
		t.Log("entity doesn't move in right direction")
		t.Fail()
	}

	e1.Dir = DOWN

	game.Update()

	// 8 | 8
	if e1.X != 8 && e1.Y != 8 {
		t.Log("entity doesn't move down")
		t.Fail()
	}

	e1.Dir = LEFT

	game.Update()

	// 0 | 8
	if e1.X != 0 && e1.Y != 8 {
		t.Log("entity doesn't move left")
		t.Fail()
	}

	e1.Dir = UP

	game.Update()

	if e1.X != 0 && e1.Y != 0 && !e1.Died {
		t.Log("entity doesn't move up and dies")
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	game := NewGame(640, 480, 8)

	game.NewEntity(0, 0, RIGHT, 0xFF)
	e2 := game.NewEntity(0, 0, LEFT, 0x0FF)

	game.DeleteEntity(e2)

	if len(game.Entities) != 1 {
		t.Log("entites is not deleted")
		t.Fail()
	}

}

func TestEntitySetDirCW(t *testing.T) {
	game := NewGame(640, 480, 8)

	e := game.NewEntity(0, 0, RIGHT, 0xFF00FF)

	e.SetDir(RIGHT)

	if e.Dir != DOWN {
		t.Log("direction should be down")
		t.Fail()
	}

	e.SetDir(RIGHT)

	if e.Dir != LEFT {
		t.Log("direction should be left")
		t.Fail()
	}

	e.SetDir(RIGHT)

	if e.Dir != UP {
		t.Log("direction should be up")
		t.Fail()
	}

	e.SetDir(RIGHT)

	if e.Dir != RIGHT {
		t.Log("direction should be right")
		t.Fail()
	}
}

func TestEntitySetDirCCW(t *testing.T) {
	game := NewGame(640, 480, 8)

	e := game.NewEntity(0, 0, RIGHT, 0xFF00FF)

	e.SetDir(LEFT)

	if e.Dir != UP {
		t.Log("direction should be up")
		t.Fail()
	}

	e.SetDir(LEFT)

	if e.Dir != LEFT {
		t.Log("direction should be left")
		t.Fail()
	}

	e.SetDir(LEFT)

	if e.Dir != DOWN {
		t.Log("direction should be down")
		t.Fail()
	}

	e.SetDir(LEFT)

	if e.Dir != RIGHT {
		t.Log("direction should be left")
		t.Fail()
	}
}

func TestGameEnd(t *testing.T) {
	game := NewGame(640, 480, 8)

	game.NewEntity(0, 0, LEFT, 0xFF)
	game.Update()

	if !game.End() {
		t.Log("Game should be finished")
		t.Fail()
	}
}
