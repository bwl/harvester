package rendering

type Rect struct{ X,Y,W,H int }

func (r Rect) union(o Rect) Rect {
	if r.W<=0 || r.H<=0 { return o }
	if o.W<=0 || o.H<=0 { return r }
	x1 := min(r.X, o.X)
	y1 := min(r.Y, o.Y)
	x2 := max(r.X+r.W, o.X+o.W)
	y2 := max(r.Y+r.H, o.Y+o.H)
	return Rect{X:x1, Y:y1, W:x2-x1, H:y2-y1}
}

func min(a,b int) int { if a<b { return a }; return b }
func max(a,b int) int { if a>b { return a }; return b }
