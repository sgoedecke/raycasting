package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	fullscreen   = true
	width        = 420
	height       = 200
	scale        = 3.0
	wallDistance = 8.0

	pos, dir, plane pixel.Vec
)

func setup() {
	pos = pixel.V(12.0, 14.5)
	dir = pixel.V(-1.0, 0.0)
	plane = pixel.V(0.0, 0.66)
}

var world = [24][24]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 7, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 7, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 2, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 6, 0, 4, 0, 0, 0, 4, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 5, 0, 0, 0, 1},
	{1, 0, 6, 0, 4, 0, 7, 0, 4, 0, 0, 0, 0, 0, 5, 0, 0, 0, 5, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 4, 0, 0, 0, 4, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 0, 4, 0, 0, 0, 5, 5, 0, 5, 5, 5, 0, 5, 5, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 4, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 5, 0, 5, 5, 5, 5, 5, 5, 5, 0, 5, 0, 1},
	{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 5, 0, 5, 0, 0, 0, 0, 0, 5, 0, 5, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 5, 5, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}


// from the global variables that hold the game state, draw a frame of the game into
// the image.RGBA buffer
func frame() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, width, height))

    // starting at x, we draw each column of pixels like so: trace a ray in that direction,
    // so we know how far away the closest wall is. then draw a vertically centered column of
    // pixels with a height that is directly proportional to the distance from the player.
	for x := 0; x < width; x++ {
		var (
			step         image.Point // the next step of the ray we're tracing (really this is a vector)
			sideDist     pixel.Vec
			perpWallDist float64
			hit, side    bool

            // rayPos is the position of the player (i.e. the source of the ray)
            // worldX, worldY start at rayPos and get advanced as we step & represent the tip of the extending ray
			rayPos, worldX, worldY = pos, int(pos.X), int(pos.Y)

			cameraX = 2*float64(x)/float64(width) - 1

			rayDir = pixel.V(
				dir.X+plane.X*cameraX,
				dir.Y+plane.Y*cameraX,
			)

			deltaDist = pixel.V(
				math.Sqrt(1.0+(rayDir.Y*rayDir.Y)/(rayDir.X*rayDir.X)),
				math.Sqrt(1.0+(rayDir.X*rayDir.X)/(rayDir.Y*rayDir.Y)),
			)
		)

		if rayDir.X < 0 {
			step.X = -1
			sideDist.X = (rayPos.X - float64(worldX)) * deltaDist.X
		} else {
			step.X = 1
			sideDist.X = (float64(worldX) + 1.0 - rayPos.X) * deltaDist.X
		}

		if rayDir.Y < 0 {
			step.Y = -1
			sideDist.Y = (rayPos.Y - float64(worldY)) * deltaDist.Y
		} else {
			step.Y = 1
			sideDist.Y = (float64(worldY) + 1.0 - rayPos.Y) * deltaDist.Y
		}

        // advance the ray (sideDist) in tiny zigzags until we hit a wall.
        // if we hit the wall going x-wards, we know it's a left-right wall
        // if we hit the wall going y-wards, we know it's a forward-back wall
		for !hit {
			if sideDist.X < sideDist.Y { // zig
				sideDist.X += deltaDist.X // why do we advance both sideDist and worldX/Y?
				worldX += step.X
				side = false
			} else {
				sideDist.Y += deltaDist.Y // zag
				worldY += step.Y
				side = true
			}

			if world[worldX][worldY] > 0 {
				hit = true
			}
		}

		var wallX float64

		if side {
			perpWallDist = (float64(worldY) - rayPos.Y + (1-float64(step.Y))/2) / rayDir.Y
			wallX = rayPos.X + perpWallDist*rayDir.X
		} else {
			perpWallDist = (float64(worldX) - rayPos.X + (1-float64(step.X))/2) / rayDir.X
			wallX = rayPos.Y + perpWallDist*rayDir.Y
		}

        // if the column we are drawing is in the very middle of the screen, we know that the distance of
        // the current ray is the distance between the player and the wall. so we record it to use it for 
        // collision detection later.
		if x == width/2 {
			wallDistance = perpWallDist
		}

		wallX -= math.Floor(wallX)

		lineHeight := int(float64(height) / perpWallDist) // how high to draw the wall. height is screen height

		if lineHeight < 1 { // since we draw in blocks, we have to round this off to at least 1.
			lineHeight = 1
		}

        // we have to draw our wall centered in the screen (so that perspective looks right), so we pick our start and end point so lineHeight has
        // the same distance both above and below it
		drawStart := -lineHeight/2 + height/2
		if drawStart < 0 {
			drawStart = 0
		}

		drawEnd := lineHeight/2 + height/2
		if drawEnd >= height {
			drawEnd = height - 1
		}

        for y := 0; y < height; y++ {
          var c color.RGBA

          // based on how high up we are, set the color
          if y >= drawStart && y < drawEnd + 1 {
            // we're drawing a wall
			c = color.RGBA{3, 30, 3, 30} // green
           // draw side walls darker to highlight corners
			if side {
				c.R = c.R / 2
				c.G = c.G / 2
				c.B = c.B / 2
			}
          }

          if y < drawStart {
            // we're drawing the floor
			c = color.RGBA{10,10,10, 1} // black 
          }

          if y >= drawEnd {
            // we're drawing the ceiling 
			c = color.RGBA{40,40,40, 1} // black 
          }

          m.Set(x, y, c)

        }
	}

	return m
}

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width)*scale, float64(height)*scale),
		VSync:       true,
		Undecorated: false,
	}

	if fullscreen {
		cfg.Monitor = pixelgl.PrimaryMonitor()
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	c := win.Bounds().Center()

	last := time.Now()

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		win.Clear(color.Black)

		dt := time.Since(last).Seconds()
		last = time.Now()

		if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			moveForward(3.5 * dt)
		}

		if win.Pressed(pixelgl.KeyA) {
			moveLeft(3.5 * dt)
		}

		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			moveBackwards(3.5 * dt)
		}

		if win.Pressed(pixelgl.KeyD) {
			moveRight(3.5 * dt)
		}

		if win.Pressed(pixelgl.KeyRight) {
			turnRight(1.2 * dt)
		}

		if win.Pressed(pixelgl.KeyLeft) {
			turnLeft(1.2 * dt)
		}

		p := pixel.PictureDataFromImage(frame())

		pixel.NewSprite(p, p.Bounds()).
			Draw(win, pixel.IM.Moved(c).Scaled(c, scale))

		win.Update()
	}
}

func moveForward(s float64) {
	if wallDistance > 0.3 {
       // only move forward if we're not too close to a wall.
       // why can't we just rely on the '0' check? because if we do,
       // it feels very slidy when you run into a wall.
		if world[int(pos.X+dir.X*s)][int(pos.Y)] == 0 {
			pos.X += dir.X * s
		}

		if world[int(pos.X)][int(pos.Y+dir.Y*s)] == 0 {
			pos.Y += dir.Y * s
		}
	}
}

func moveLeft(s float64) {
	if world[int(pos.X-plane.X*s)][int(pos.Y)] == 0 {
		pos.X -= plane.X * s
	}

	if world[int(pos.X)][int(pos.Y-plane.Y*s)] == 0 {
		pos.Y -= plane.Y * s
	}
}

func moveBackwards(s float64) {
	if world[int(pos.X-dir.X*s)][int(pos.Y)] == 0 {
		pos.X -= dir.X * s
	}

	if world[int(pos.X)][int(pos.Y-dir.Y*s)] == 0 {
		pos.Y -= dir.Y * s
	}
}

func moveRight(s float64) {
	if world[int(pos.X+plane.X*s)][int(pos.Y)] == 0 {
		pos.X += plane.X * s
	}

	if world[int(pos.X)][int(pos.Y+plane.Y*s)] == 0 {
		pos.Y += plane.Y * s
	}
}

func turnRight(s float64) {
	oldDirX := dir.X

	dir.X = dir.X*math.Cos(-s) - dir.Y*math.Sin(-s)
	dir.Y = oldDirX*math.Sin(-s) + dir.Y*math.Cos(-s)

	oldPlaneX := plane.X

	plane.X = plane.X*math.Cos(-s) - plane.Y*math.Sin(-s)
	plane.Y = oldPlaneX*math.Sin(-s) + plane.Y*math.Cos(-s)
}

func turnLeft(s float64) {
	oldDirX := dir.X

	dir.X = dir.X*math.Cos(s) - dir.Y*math.Sin(s)
	dir.Y = oldDirX*math.Sin(s) + dir.Y*math.Cos(s)

	oldPlaneX := plane.X

	plane.X = plane.X*math.Cos(s) - plane.Y*math.Sin(s)
	plane.Y = oldPlaneX*math.Sin(s) + plane.Y*math.Cos(s)
}

func main() {
	setup()
	pixelgl.Run(run)
}
