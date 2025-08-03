package rendering

func CalculatePosition(pos Position, bounds Bounds, containerW, containerH int) (int, int) {
	var x int
	switch pos.Horizontal {
	case Left:
		x = 0
	case CenterH:
		x = (containerW - bounds.Width) / 2
	case Right:
		x = containerW - bounds.Width
	}
	var y int
	switch pos.Vertical {
	case Top:
		y = 0
	case CenterV:
		y = (containerH - bounds.Height) / 2
	case Bottom:
		y = containerH - bounds.Height
	}
	x += pos.OffsetX
	y += pos.OffsetY
	return x, y
}
