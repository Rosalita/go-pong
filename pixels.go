package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

type colour struct {
	r, g, b byte
}

func main() {
	window, err := sdl.CreateWindow("Hello I'm a window",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(windowWidth),
		int32(windowHeight),
		sdl.WINDOW_SHOWN)
	if err != nil {
		log.Println(err)
		return
	}
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Println(err)
		return
	}
	defer renderer.Destroy()
	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(windowWidth),
		int32(windowHeight))
	if err != nil {
		log.Println(err)
		return
	}
	defer texture.Destroy()
	pixels := make([]byte, windowWidth*windowHeight*4)
	for y := 0; y < windowHeight; y++ {
		for x := 0; x < windowWidth; x++ {
			if y%2 == 0 {
				setPixel(x, y, colour{byte(x % 255), 0, 0}, pixels)
			} else {
				setPixel(x, y, colour{0, 0, byte(y % 255)}, pixels)
			}
		}
	}
	texture.Update(nil, pixels, windowWidth*4)
	renderer.Copy(texture, nil, nil)
	renderer.Present()
	sdl.Delay(5000)
}
func setPixel(x, y int, colour colour, pixels []byte) {
	index := (y*windowWidth + x) * 4
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = colour.r
		pixels[index+1] = colour.g
		pixels[index+2] = colour.b
	}
}
