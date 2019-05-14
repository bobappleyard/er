package diagram

import (
	"github.com/ajstarks/svgo"
	"github.com/bobappleyard/er"
)

// Draw a diagram for the model in question into the provided SVG.
func Draw(s *svg.SVG, m *er.EntityModel) {
	tw := buildTowers(m)
	tw.layoutDiagram()
	s.Start(tw.body.w+5, tw.body.h+5)
	for _, t := range tw.down {
		t.draw(s)
	}
	s.End()
}

func (t *tower) draw(s *svg.SVG) {
	// s.Rect(t.body.x, t.body.y, t.body.w, t.body.h, "fill:gray;opacity:0.5")
	t.drawHead(s)
	for _, t := range t.down {
		t.draw(s)
	}
	for _, l := range t.lines {
		s.Line(l.x1, l.y1, l.x2, l.y2, "stroke:black")
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
