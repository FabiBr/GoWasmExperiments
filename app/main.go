package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
)

type particle struct {
	id         string
	pos_x      float64
	pos_y      float64
	rad        float64
	velocity   float64
	velocity_x float64
	velocity_y float64
	angle      float64
	collided   bool
	activated  bool
}

func newParticle(id string) *particle {
	p := particle{id: id}
	p.rad = 20
	p.pos_x = p.rad * 10
	p.pos_y = p.rad * 10
	p.velocity = 5
	p.velocity_x = 0
	p.velocity_y = 0
	p.angle = 135 * (math.Pi / 180)
	p.collided = false
	p.activated = false

	return &p
}

func main() {

	var document js.Value = js.
		Global().
		Get("document")

	var body js.Value = document.Call("getElementById", "mainBody")
	var canvas js.Value = document.Call("getElementById", "fsCanvas")
	var ctx js.Value = canvas.Call("getContext", "2d")

	width := body.Get("clientWidth").Float()
	height := body.Get("clientHeight").Float()

	canvas.Set("width", width)
	canvas.Set("height", height)
	ctx.Set("fillStyle", "blue")
	ctx.Call("fillRect", "0", "0", width, height)

	done := make(chan struct{})

	var renderFrame js.Func

	var particleArray []*particle
	var ballNr int = 50

	// Creates busted Balls
	particleArray = ballBuster(ballNr, width, height, particleArray)

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		Render(ctx, particleArray, width, height)

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	js.Global().Call("requestAnimationFrame", renderFrame)

	<-done
}

func ballBuster(ballNr int, width float64, height float64, particleArray []*particle) []*particle {
	for i := 0; i < ballNr; i++ {
		p := newParticle(fmt.Sprint(i))
		p.angle = rand.Float64() * 360 * (math.Pi / 180)
		p.pos_x = (rand.Float64() * (width))
		p.pos_y = (rand.Float64() * (height))

		particleArray = append(particleArray, p)
	}
	return particleArray
}

func spawnActivator(width float64, height float64) *particle {
	activator := newParticle("x")
	activator.velocity = 0
	activator.pos_x = (rand.Float64() * (width))
	activator.pos_y = (rand.Float64() * (height))
	activator.collided = true
	activator.activated = true
	return activator
}

func Render(ctx js.Value, particleArray []*particle, width float64, height float64) bool {

	//copy := particleArray
	ctx.Set("fillStyle", "blue")
	ctx.Call("fillRect", "0", "0", width, height)
	updateParticleVelocityAndPosition(particleArray)

	chance := rand.Float64() * 100

	for _, p := range particleArray {
		handleWallCollision(p, width, height)
		handleParticleCollision(particleArray, p)
		drawCircle(ctx, p)
		if chance == 50 {
			//activator := spawnActivator(width, height)
		}
	}

	return true
}

func updateParticleVelocityAndPosition(particleArray []*particle) {
	for _, p := range particleArray {
		p.collided = false
		p.velocity_x = p.velocity * math.Cos(p.angle)
		p.velocity_y = p.velocity * math.Sin(p.angle)

		p.pos_x += p.velocity_x
		p.pos_y += p.velocity_y
	}
}

func handleWallCollision(p *particle, width float64, height float64) {
	if p.pos_x < p.rad {
		p.pos_x = p.rad
		p.angle = math.Atan2(p.velocity_y, p.velocity_x*-1)
	}

	if p.pos_x > width-p.rad {
		p.pos_x = width - p.rad
		p.angle = math.Atan2(p.velocity_y, p.velocity_x*-1)
	}

	if p.pos_y < p.rad {
		p.pos_y = p.rad
		p.angle = math.Atan2(p.velocity_y*-1, p.velocity_x)
	}

	if p.pos_y > height-p.rad {
		p.pos_y = height - p.rad
		p.angle = math.Atan2(p.velocity_y*-1, p.velocity_x)
	}
}

func handleParticleCollision(particleArray []*particle, p *particle) {
	for _, p2 := range particleArray {
		if p.id != p2.id && !p.collided && !p2.collided {
			if math.Pow(p2.pos_x-p.pos_x, 2)+math.Pow(p2.pos_y-p.pos_y, 2) < math.Pow(p2.rad+p.rad, 2) {
				p.collided = true
				p2.collided = true
				distX := p2.pos_x - p.pos_x
				distY := p2.pos_y - p.pos_y
				distCenters := math.Sqrt(distX*distX + distY*distY)

				dist := (p.rad + p2.rad) - distCenters
				if dist > 0 {
					println(int(distX))
					distX /= distCenters
					distY /= distCenters

					p.pos_x -= distX * (dist)
					p.pos_y -= distY * (dist)
					p2.pos_x += distX * (dist)
					p2.pos_y += distY * (dist)
				}

				temp1X := p.velocity_x
				temp1Y := p.velocity_y
				temp2X := p2.velocity_x
				temp2Y := p2.velocity_y

				p.angle = math.Atan2(temp2Y, temp2X)
				p2.angle = math.Atan2(temp1Y, temp1X)

			}
		}
	}
}

func drawCircle(ctx js.Value, p *particle) {
	ctx.Call("beginPath")
	ctx.Call("arc", p.pos_x, p.pos_y, p.rad, 0, 2*math.Pi, false)
	if p.activated {
		ctx.Set("fillStyle", "yellow")
	} else {
		ctx.Set("fillStyle", "green")
	}
	ctx.Call("fill")
	ctx.Set("lineWidth", 5)
	ctx.Set("strokeStyle", "#003300")
	ctx.Call("stroke")
}
