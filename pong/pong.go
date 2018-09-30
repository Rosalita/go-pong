package main

import (
	"fmt"
	"time"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

type colour struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius    float32
	xVelocity float32
	yVelocity float32
	colour    colour
}

type paddle struct {
	pos
	w  float32
	h float32
	speed float32
	colour colour
}

func (p *paddle) draw(pixels []byte) {
	startX := int(p.x) - int(p.w)/2
	startY := int(p.y) - int(p.h)/2

	for y := 0; y < int(p.h); y++ {
		for x := 0; x < int(p.w); x++ {
			setPixel(startX+x, startY+y, p.colour, pixels)
		}
	}
}

func (p *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		p.y -= p.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		p.y += p.speed * elapsedTime
	}
}

func (p *paddle) aiUpdate(ball *ball, elapsedTime float32){
 p.y = ball.y
}

func (b *ball) draw(pixels []byte) {
	for y := int(-b.radius); y < int(b.radius); y++ {
		for x := int(-b.radius); x < int(b.radius); x++ {
			if (x*x)+(y*y) < (int(b.radius) * int(b.radius)) {
				setPixel(int(b.x)+x, int(b.y)+y, b.colour, pixels)
			}
		}
	}
}

func (b *ball) update( leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32) {
	b.x += b.xVelocity * elapsedTime
	b.y += b.yVelocity * elapsedTime

	if b.y - b.radius < 0 || b.y + b.radius > windowHeight {
		b.yVelocity = -b.yVelocity
	}

	if b.x < 0 || b.x > windowWidth {
		b.pos = getCentre()
	}

	if b.x - b.radius < leftPaddle.x + leftPaddle.w/2{
		if b.y > leftPaddle.y - leftPaddle.h/2 && b.y < leftPaddle.y + leftPaddle.h/2{
			b.xVelocity = -b.xVelocity
		}
	}

	if b.x + b.radius > rightPaddle.x - rightPaddle.w/2{
		if b.y > rightPaddle.y - rightPaddle.h/2 && b.y < rightPaddle.y + rightPaddle.h/2{
			b.xVelocity = -b.xVelocity
		}
	}
}

func clear(pixels []byte){
  for i := range pixels{
	  pixels[i] = 0
  }
}

func getCentre() pos {
	return pos{float32(windowWidth) /2, float32(windowHeight) /2}
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

	player1 := paddle{pos{100, 100}, 20, 100, 300, colour{255, 255, 255}}
	player2 := paddle{pos{700, 100}, 20, 100, 300, colour{255, 255, 255}}

	ball := ball{pos{300, 300}, 20, 400, 400, colour{255, 255, 255}}
	

	keyState := sdl.GetKeyboardState()
	var frameStart time.Time
	var elapsedTime float32

	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clear(pixels)

		player1.update(keyState, elapsedTime)
		player2.aiUpdate(&ball, elapsedTime)
		ball.update(&player1, &player2, elapsedTime)


		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		texture.Update(nil, pixels, windowWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
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
