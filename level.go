package gogame

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

// TODO this is stupid to store it as an interface
var LevelLines []*chipmunk.Shape

func LevelInit() {
	// bounding box
	space = chipmunk.NewSpace()

	LevelLines = []*chipmunk.Shape{
		chipmunk.NewSegment(vect.Vect{0, 0}, vect.Vect{Size, 0}, 0),
		chipmunk.NewSegment(vect.Vect{Size, 0}, vect.Vect{Size, Size}, 0),
		chipmunk.NewSegment(vect.Vect{Size, Size}, vect.Vect{0, Size}, 0),
		chipmunk.NewSegment(vect.Vect{0, Size}, vect.Vect{0, 0}, 0),
	}
	for _, segment := range LevelLines {
		segment.SetElasticity(1.0)
		segment.SetFriction(1.0)
		staticBody := chipmunk.NewBodyStatic()
		staticBody.AddShape(segment)
		staticBody.CallbackHandler = &HandleCollisions{}
		space.AddBody(staticBody)
	}
}

func DrawLevel(screen *ebiten.Image) {
	seg1 := LevelLines[0].GetAsSegment()
	img1, _ := ebiten.NewImage(Size, 10, ebiten.FilterNearest)
	img1.Fill(color.RGBA{0xFF, 0x00, 0x00, 0xFF})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(seg1.A.X), 0)
	screen.DrawImage(img1, opts)

	//for _, segment := range LevelLines {
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
