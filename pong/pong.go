package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

type colour struct {
	r, g, b byte
}

type position struct {
	x, y float32
}

type ball struct {
	position
	radius int
	xVelocity float32
	yVelocity float32
	colour colour
}

type paddle struct {
	position 
	width int
	height int
	colour colour
}

func(p *paddle) draw(pixels []byte){
	startX := int(p.x) - p.width/2
	startY := int(p.y) - p.height/2

	for y :=0; y < p.height; y++{
		for x := 0;  x < p.width; x++ {
			setPixel(startX+x, startY+y, colour{255, 255, 255}, pixels)
		}
	}

}

func main() {

err := sdl.Init(sdl.INIT_EVERYTHING)

if err != nil {
	fmt.Println(err)
	return
}
defer sdl.Quit()

	window, err := sdl.CreateWindow("Hello I'm a window",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(windowWidth),
		int32(windowHeight),
		sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()
	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(windowWidth),
		int32(windowHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer texture.Destroy()
	pixels := make([]byte, windowWidth*windowHeight*4)
	for y := 0; y < windowHeight; y++ {
		for x := 0; x < windowWidth; x++ {
			if y%2 == 0 {
				setPixel(x, y, colour{byte(x % 255), 0, 0}, pixels)
			} else if y%3 == 0{
				setPixel(x, y, colour{0, byte(y % 255), 0}, pixels)
			} else {
				setPixel(x, y, colour{0, 0, byte(y % 255)}, pixels)
			}
		}
	}
	texture.Update(nil, pixels, windowWidth*4)
	renderer.Copy(texture, nil, nil)
	renderer.Present()


	//sdl.Delay(5000)

	for{
		for event := sdl.PollEvent(); event != nil; sdl.PollEvent(){
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		sdl.Delay(16)
	}



}
func setPixel(x, y int, colour colour, pixels []byte) {
	index := (y*windowWidth + x) * 4
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = colour.r
		pixels[index+1] = colour.g
		pixels[index+2] = colour.b
	}
}
