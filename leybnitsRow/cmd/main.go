package main

import "leybnitsRow/internal"

func main() {

	pc := internal.NewPiCounter(3, 3)
	pc.Start()
	pc.Print()

}
