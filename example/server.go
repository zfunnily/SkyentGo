package main

import (
	"fmt"
	"pro2d/common/components"
)

func ACbk(ctx *components.Context, cbud interface{}, data []byte, typ int) {
	fmt.Printf("A recv type %d: %s\n", typ, string(data))
}

func BCbk(ctx *components.Context, cbud interface{}, data []byte, typ int) {
	fmt.Printf("B recv type %d: %s\n", typ, string(data))
}

func main() {
	components.MAInst().Start()

	ACtx := components.NewContext()
	BCtx := components.NewContext()
	ACtx.Callback(nil, ACbk)
	BCtx.Callback(nil, BCbk)

	s1 := []byte("hello B, I'm A")
	ret := ACtx.Send(uint32(0), BCtx.Handle(), 1, 1, s1)
	fmt.Println("ACtx ret: ", ret)

	s2 := []byte("hello A, I'm B")
	ret = BCtx.Send(uint32(0), ACtx.Handle(), 1, 1, s2)
	fmt.Println("BCtx ret: ", ret)

	components.TWInst().TimeOut(ACtx.Handle(), 0, 1)
	select {}
}
