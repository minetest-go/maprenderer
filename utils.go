package maprenderer

func SortPos(p1, p2 *Pos) (*Pos, *Pos) {
	return &Pos{
			min(p1[0], p2[0]),
			min(p1[1], p2[1]),
			min(p1[2], p2[2]),
		}, &Pos{
			max(p1[0], p2[0]),
			max(p1[1], p2[1]),
			max(p1[2], p2[2]),
		}
}
