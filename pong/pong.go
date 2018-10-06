package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

type gameState int

const (
	titleScreen gameState = iota
	restart
	play
)

var state = titleScreen

var chars = map[string][]byte{
	"0": {1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	"1": {1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	"2": {1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	"3": {1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	"a": {1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
	},
	"b": {1, 1, 1,
		1, 0, 1,
		1, 1, 0,
		1, 0, 1,
		1, 1, 1,
	},
	"e": {1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	"g": {1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	"i": {0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	},
	"n": {1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
	},
	"p": {1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
	},
	"r": {1, 1, 1,
		1, 0, 1,
		1, 1, 0,
		1, 0, 1,
		1, 0, 1,
	},
	"s": {1, 1, 1,
		1, 0, 0,
		1, 1, 0,
		0, 0, 1,
		1, 1, 1,
	},
	"t": {1, 1, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	},
	"w": {1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
	},
	" ": {0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
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

type text struct {
	pos
	size       int
	characters string
	colour     colour
}

func (t *text) draw(pixels []byte) {
	x := int(t.x) - (t.size*3)/2
	y := int(t.y) - (t.size*5)/2

	for i, c := range t.characters {
		if i == 0 {
			drawChar(pos{float32(x), float32(y)}, t.colour, t.size, string(c), pixels)
			x += (t.size * 3) + t.size
			continue
		}
		drawChar(pos{float32(x), float32(y)}, t.colour, t.size, string(c), pixels)
		x += (t.size * 3) + t.size
	}
}

func (t *text) rainbowUpdate() {

	rainbow := []colour{
		colour{255, 0, 0},
		colour{255, 128, 0},
		colour{255, 255, 0},
		colour{0, 255, 0},
		colour{0, 255, 255},
		colour{0, 0, 255},
		colour{128, 0, 255},
		colour{255, 0, 255},
		colour{255, 0, 128},
	}

	rand.Seed(time.Now().Unix())
	randnum := rand.Intn(9)

	t.colour = rainbow[randnum]

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
	drawChar(pos{numX, 35}, p.colour, 10, strconv.Itoa(p.score), pixels)
}

func (p *paddle) update(keyState []uint8, controllerAxis int16, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		p.y -= p.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		p.y += p.speed * elapsedTime
	}

	if math.Abs(float64(controllerAxis)) > 1500 {
		pct := float32(controllerAxis) / 32767.0
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
		state = restart
	}

	if b.x > windowWidth {
		leftPaddle.score++
		b.pos = getCentre()
		state = restart
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

func drawChar(pos pos, colour colour, size int, char string, pixels []byte) {
	startX := int(pos.x) - (size*3)/2 // numbers are minimum 3 pixels wide
	startY := int(pos.y) - (size*5)/2 // numbers are minimum 5 pixels high

	for i, v := range chars[char] {
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
	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()
	}

	pixels := make([]byte, windowWidth*windowHeight*4)

	titleText := text{pos{100, 100}, 10, "rainb0w p0ng", colour{255, 255, 255}}
	pressStartText := text{pos{300, 300}, 10, "press start", colour{255, 255, 255}}
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

		if state == titleScreen {
			clear(pixels)
			titleText.rainbowUpdate()
			titleText.draw(pixels)
			pressStartText.draw(pixels)

			if keyState[sdl.SCANCODE_SPACE] != 0 {
				state = play
			}

		}

		if state == play {
			clear(pixels)
			player1.update(keyState, controllerAxis, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)

			player1.draw(pixels)
			player2.draw(pixels)
			ball.draw(pixels)

		}

		if state == restart {

			player1.draw(pixels)
			player2.draw(pixels)
			ball.draw(pixels)

			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}

				state = play
			}
		}

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
