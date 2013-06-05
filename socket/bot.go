package socket

import (
    "fmt"
)
type bot struct {
    X int `json:"x"`
    Y int `json:"y"`
    Dir int8 `json:"dir"`
    Color int `json:"color"`
    S int `json:"size"`
    died bool
    Score int
    CurrentState int
    Gene[16][4]Genom
}
type Genom struct {
    NewState int
    WhatToDo int
}

NewBot(x,y,color int) *bot {
    b := &bot{
        X: x*16,
        Y: y*16,
        Dir: int8(rand.Intn(4))
        Color: color,
        S: 16,
        died: true,
        Score: 1,
        CurrentState: 0}
    for i:= 0; i<16; i++ {
        for j:= 0; j < 4; j++ {
            b.Gene[i][j].NewState = rand.Intn(16)
            b.Gene[i][j].WhatToDo = rand.Intn(4)
        }
    }
    return b
}

func (b *bot) Mutate() {
    for i:= 0; i<16; i++ {
        for j:= 0; j < 4; j++ {
            if rand.Float64() < settings.MutateProp {
                b.Gene[i][j].NewState = rand.Intn(16)
            }
            if rand.Float64() < settings.MutateProp {
                b.Gene[i][j].WhatToDo = rand.Intn(4)
            }
        }
    }
}

func (b *bot) CrossWith(o *bot){

}

func (b *bot) Clone() *bot {


}
