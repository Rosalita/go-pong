package main

import (
	"fmt"
	"time"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

type gameState int

const (
	start gameState = iota
	play
)

var state = start

var nums = [][]byte{
	{1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
}

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
	w      float32
	h      float32
	speed  float32
	score  int
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

	numX := lerp(p.x, getCentre().x, 0.2)
	drawNumber(pos{numX, 35}, p.colour, 10, p.score, pixels)
}

func (p *paddle) update(keyState []uint8, controllerAxis int16, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		p.y -= p.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		p.y += p.speed * elapsedTime
	}

	if math.Abs(float64(controllerAxis)) > 1500 {
		pct := float32(controllerAxis) / 	32767.0
		p.y += p.speed * pct * elapsedTime
	}
}

func (p *paddle) aiUpdate(ball *ball, elapsedTime float32) {
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

func (b *ball) update(leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32) {
	b.x += b.xVelocity * elapsedTime
	b.y += b.yVelocity * elapsedTime

	if b.y-b.radius < 0 {
		b.yVelocity = -b.yVelocity
		b.y = b.radius
	}

	if b.y+b.radius > windowHeight {
		b.yVelocity = -b.yVelocity
		b.y = windowHeight - b.radius
	}

	if b.x < 0 {
		rightPaddle.score++
		b.pos = getCentre()
		state = start
	}

	if b.x > windowWidth {
		leftPaddle.score++
		b.pos = getCentre()
		state = start
	}

	if b.x-b.radius < leftPaddle.x+leftPaddle.w/2 {
		if b.y > leftPaddle.y-leftPaddle.h/2 && b.y < leftPaddle.y+leftPaddle.h/2 {
			b.xVelocity = -b.xVelocity
			b.x = leftPaddle.x + leftPaddle.w/2.0 + b.radius
		}
	}

	if b.x+b.radius > rightPaddle.x-rightPaddle.w/2 {
		if b.y > rightPaddle.y-rightPaddle.h/2 && b.y < rightPaddle.y+rightPaddle.h/2 {
			b.xVelocity = -b.xVelocity
			b.x = rightPaddle.x - rightPaddle.w/2.0 - b.radius
		}
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func drawNumber(pos pos, colour colour, size int, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2 // numbers are minimum 3 pixels wide
	startY := int(pos.y) - (size*5)/2 // numbers are minimum 5 pixels high

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, colour, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}

}

func getCentre() pos {
	return pos{float32(windowWidth) / 2, float32(windowHeight) / 2}
}

// linear interpolation
func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
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

	var controllerHandlers []*sdl.GameController
	for i := 0;  i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()
	}	

	pixels := make([]byte, windowWidth*windowHeight*4)

	player1 := paddle{pos{100, 100}, 20, 100, 300, 0, colour{255, 255, 255}}
	player2 := paddle{pos{700, 100}, 20, 100, 300, 0, colour{255, 255, 255}}

	ball := ball{pos{300, 300}, 20, 400, 400, colour{255, 255, 255}}

	keyState := sdl.GetKeyboardState()
	var frameStart time.Time
	var elapsedTime float32
	var controllerAxis int16

	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		for _, controller := range controllerHandlers {
			if controller != nil {
				controllerAxis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
			
			}
		}

		if state == play {

			player1.update(keyState, controllerAxis, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)

		}

		if state == start{
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3{
					player1.score = 0
					player2.score = 0
				}
			
				state = play
			}
		}

		clear(pixels)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		texture.Update(nil, pixels, windowWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())

		if elapsedTime < 0.005 { // less than 5ms limit to 200fps
			sdl.Delay(5 - uint32(elapsedTime/1000))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}

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
