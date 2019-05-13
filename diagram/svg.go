package diagram

import (
	"github.com/ajstarks/svgo"
	"github.com/bobappleyard/er"
)

func Generate(s *svg.SVG, m *er.EntityModel) {
	tw := towersFor(m)
	tw.calcLayout(0, 0)
	s.Start(tw.body.w+5, tw.body.h+5)
	for _, t := range tw.down {
		t.draw(s)
	}
	s.End()
}

func (t *tower) draw(s *svg.SVG) {
	t.drawHead(s)
	for _, t := range t.down {
		t.draw(s)
	}
}

func (t *tower) drawHead(s *svg.SVG) {
	x, y, w, h := t.head.x, t.head.y, t.head.w, t.head.h
	s.Roundrect(x+5, y+5, w-5, h-5, 5, 5, "fill:lightyellow;stroke:black")
	s.Text(x+8, y+15, t.t.Name, "font-size:8pt")
	for j, a := range t.t.Attributes {
		style := "font-size:8pt"
		if a.Identifying {
			style += "; text-decoration:underline"
		}
		s.Text(x+10, y+15+(j+1)*12, "- "+a.Name, style)
	}
}
