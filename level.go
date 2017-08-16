package gogame

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/jakecoffman/physics"
)

var LevelLines []*physics.Shape

func LevelInit() {
	// bounding box
	space = physics.NewSpace()
	//body := space.StaticBody()

	LevelLines = []*physics.Shape{
	//space.AddShape(physics.NewSegment(body, &physics.Vector{0, 0}, &physics.Vector{Size, 0}, 0)),
	//space.AddShape(physics.NewSegment(body, &physics.Vector{Size, 0}, &physics.Vector{Size, Size}, 0)),
	//space.AddShape(physics.NewSegment(body, &physics.Vector{Size, Size}, &physics.Vector{0, Size}, 0)),
	//space.AddShape(physics.NewSegment(body, &physics.Vector{0, Size}, &physics.Vector{0, 0}, 0)),
	}
	for _, segment := range LevelLines {
		segment.E = 1.0
		segment.U = 1.0
	}
}

func DrawLevel(screen *ebiten.Image) {
	//seg1 := LevelLines[0].Class.(*physics.Segment)
	//img1, _ := ebiten.NewImage(Size, 10, ebiten.FilterNearest)
	//img1.Fill(color.RGBA{0xFF, 0x00, 0x00, 0xFF})
	//opts := &ebiten.DrawImageOptions{}
	//opts.GeoM.Translate(seg1.A.X, seg1.A.Y)
	//screen.DrawImage(img1, opts)

	//for _, segment := range LevelLines {
	//	seg := segment.GetAsSegment()
	//	img, _ := ebiten.NewImage(Size, 1, ebiten.FilterNearest)
	//	img.Fill(color.White)
	//	opts = &ebiten.DrawImageOptions{}
	//	opts.GeoM.Translate(float64(-Size/2), float64(-Size/2))
	//
	//	opts.GeoM.Rotate(float64(p.Shape.Body.Angle() * physics.DegreeConst))
	//opts.GeoM.Translate(float64(p.Shape.Body.Position().X), float64(p.Shape.Body.Position().Y))

}
