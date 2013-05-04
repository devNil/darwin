package socket

import(
	"fmt"
	"log"
	ws "code.google.com/p/go.net/websocket"
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

func Run(){
    var (
        lastTime, lastTimer1 int64
        unprocessed, nsPerTick float64
        ticks int32
    )
    lastTime = time.Now().UnixNano()
    unprocessed = 0
    nsPerTick = 1000000000.0 / 60
    ticks = 0
    lastTimer1 = time.Now().UnixNano()

    for {
        now := time.Now().UnixNano()
        unprocessed += (float64((now - lastTime)) / nsPerTick)
        lastTime = now
        //fmt.Println(unprocessed)
        for (unprocessed >= 1) {
            ticks++
            //tick()
            unprocessed--
        }
        if ((time.Now().UnixNano() - lastTimer1) > 1000000000.0) {
            lastTimer1 += 1000000000.0
            fmt.Printf("ticks: %d \n",ticks)
            ticks = 0
        }

    }

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
