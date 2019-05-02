package diagram

import (
	"github.com/ajstarks/svgo"
	"github.com/bobappleyard/er"
)

func Generate(s *svg.SVG, m *er.EntityModel) {
	tw := towersFor(m)
	s.Start(tw.width()+5, tw.depth()+5)
	//s.Text(10, 30, m.Name, "font-size:28px")
	x := 0
	for _, tw := range tw.cs {
		generateTower(s, tw, x, 0)
		x += tw.width()
	}
	s.End()
}

func generateTower(s *svg.SVG, tw *tower, x, y int) {
	t := tw.t
	w, h := tw.bounds()
	s.Roundrect(x+5, y+5, w-5, h-5, 5, 5, "fill:lightyellow;stroke:black")
	s.Text(x+8, y+15, t.Name, "font-size:8pt")
	for j, a := range t.Attributes {
		s.Text(x+10, y+15+(j+1)*12, "- "+a.Name, "font-size:8pt; text-decoration:underline")
	}
	y += h
	for _, tw := range tw.cs {
		generateTower(s, tw, x, y+1)
		x += tw.width()
	}
}
