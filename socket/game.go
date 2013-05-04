package socket

import(
	"fmt"
	ws "code.google.com/p/go.net/websocket"
)

var colors = {0xFF00FF, 0xFFFF00}

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

func Run(){
}

type command struct{
	Id int8
	Value []byte
}

type client struct{
	conn *ws.Conn
	e *entity
}

