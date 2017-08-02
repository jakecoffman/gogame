package gogame

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

// TODO this is stupid to store it as an interface
var staticLines []*chipmunk.Shape

func LevelInit() {
	// bounding box
	space = chipmunk.NewSpace()

	staticLines = []*chipmunk.Shape{
		chipmunk.NewSegment(vect.Vect{0,0}, vect.Vect{Size, 0}, 0),
		chipmunk.NewSegment(vect.Vect{Size,0}, vect.Vect{Size, Size}, 0),
		chipmunk.NewSegment(vect.Vect{Size,Size}, vect.Vect{0, Size}, 0),
		chipmunk.NewSegment(vect.Vect{0,Size}, vect.Vect{0, 0}, 0),
	}
	for _, segment := range staticLines {
		segment.SetElasticity(0.1)
		staticBody := chipmunk.NewBodyStatic()
		staticBody.AddShape(segment)
		staticBody.CallbackHandler = &HandleCollisions{}
		space.AddBody(staticBody)
	}
}

func DrawLevel(screen *ebiten.Image) {
	seg1 := staticLines[0].GetAsSegment()
	img1, _ := ebiten.NewImage(Size, 10, ebiten.FilterNearest)
	img1.Fill(color.RGBA{0xFF, 0x00, 0x00, 0xFF})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(seg1.A.X), 0)
	screen.DrawImage(img1, opts)

	//for _, segment := range staticLines {
	//	seg := segment.GetAsSegment()
	//	img, _ := ebiten.NewImage(Size, 1, ebiten.FilterNearest)
	//	img.Fill(color.White)
	//	opts = &ebiten.DrawImageOptions{}
	//	opts.GeoM.Translate(float64(-Size/2), float64(-Size/2))
	//
	//	opts.GeoM.Rotate(float64(p.Shape.Body.Angle() * chipmunk.DegreeConst))
		//opts.GeoM.Translate(float64(p.Shape.Body.Position().X), float64(p.Shape.Body.Position().Y))

}

type HandleCollisions struct {

}

func (h *HandleCollisions) CollisionEnter(arbiter *chipmunk.Arbiter) bool {

	return true
}

func (h *HandleCollisions) CollisionPreSolve(arbiter *chipmunk.Arbiter) bool {

	return true
}

func (h *HandleCollisions) CollisionPostSolve(arbiter *chipmunk.Arbiter) {

}

func (h *HandleCollisions) CollisionExit(arbiter *chipmunk.Arbiter) {

}
