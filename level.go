package gogame

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
)

// TODO this is stupid to store it as an interface
var staticLines []*chipmunk.Shape

func LevelInit() {
	// bounding box
	space = chipmunk.NewSpace()

	staticLines = []*chipmunk.Shape{
		chipmunk.NewSegment(vect.Vect{0,0}, vect.Vect{size, 0}, 0),
		chipmunk.NewSegment(vect.Vect{size,0}, vect.Vect{size, size}, 0),
		chipmunk.NewSegment(vect.Vect{size,size}, vect.Vect{0, size}, 0),
		chipmunk.NewSegment(vect.Vect{0,size}, vect.Vect{0, 0}, 0),
	}
	for i, segment := range staticLines {
		segment.SetElasticity(0.1)
		staticBody := chipmunk.NewBodyStatic()
		staticBody.AddShape(segment)
		staticBody.CallbackHandler = &HandleCollisions{i}
		space.AddBody(staticBody)
	}

}

func DrawLevel(screen *ebiten.Image) {
	seg1 := staticLines[0].GetAsSegment()
	img1, _ := ebiten.NewImage(size, 10, ebiten.FilterNearest)
	img1.Fill(color.RGBA{0xFF, 0x00, 0x00, 0xFF})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(seg1.A.X), 0)
	screen.DrawImage(img1, opts)

	//for _, segment := range staticLines {
	//	seg := segment.GetAsSegment()
	//	img, _ := ebiten.NewImage(size, 1, ebiten.FilterNearest)
	//	img.Fill(color.White)
	//	opts = &ebiten.DrawImageOptions{}
	//	opts.GeoM.Translate(float64(-size/2), float64(-size/2))
	//
	//	opts.GeoM.Rotate(float64(p.Shape.Body.Angle() * chipmunk.DegreeConst))
		//opts.GeoM.Translate(float64(p.Shape.Body.Position().X), float64(p.Shape.Body.Position().Y))

}

type HandleCollisions struct {
i int
}

func (h *HandleCollisions) CollisionEnter(arbiter *chipmunk.Arbiter) bool {
	log.Println("CollisionEnter", h.i)
	return false
}

func (h *HandleCollisions) CollisionPreSolve(arbiter *chipmunk.Arbiter) bool {
	log.Println("CollisionPreSolve", h.i)
	return false
}

func (h *HandleCollisions) CollisionPostSolve(arbiter *chipmunk.Arbiter) {
	log.Println("CollisionPostSolve", h.i)
}

func (h *HandleCollisions) CollisionExit(arbiter *chipmunk.Arbiter) {
	log.Println("CollisionExit", h.i)
}
