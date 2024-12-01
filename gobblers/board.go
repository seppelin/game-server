package gobblers

import "fmt"

type Pos int
type Size int
type Layer uint16
type Player uint16
type State uint16
type BoardID uint64

const (
	StateNone State = iota
	StateWin
	StateLoss
	StateDraw
)

type Move struct {
	New  bool `json:"new"`
	From Pos  `json:"from"`
	To   Pos  `json:"to"`
	Size Size `json:"size"`
}

type Board struct {
	Layers [2][3]Layer
	Pieces [2][3]int
	Player Player
}

func NewBoard() Board {
	return Board{Pieces: [2][3]int{{2, 2, 2}, {2, 2, 2}}}
}

func BoardFromID(id BoardID) Board {
	var b = NewBoard()
	b.Player = Player(id & 0b11111111)
	for sign := range 2 {
		for size := range 3 {
			id >>= 9
			b.Layers[1-sign][2-size] = Layer(id & 0b111111111)
			for pos := range 9 {
				if b.Layers[1-sign][2-size]&(1<<pos) != 0 {
					b.Pieces[1-sign][2-size] -= 1
				}
			}
		}
	}
	return b
}

func (b Board) ID() BoardID {
	genID := func() BoardID {
		id := BoardID(0)
		for sign := range 2 {
			for size := range 3 {
				id |= BoardID(b.Layers[sign][size])
				id <<= 9
			}
		}
		id |= BoardID(b.Player)
		return id
	}
	reorder := func(order [9]int) {
		layers := b.Layers
		for sign := range 2 {
			for size := range 3 {
				b.Layers[sign][size] = 0
			}
		}
		for sign := range 2 {
			for size := range 3 {
				for pos := range 9 {
					if layers[sign][size]&(1<<pos) == 0 {
						continue
					}
					b.Layers[sign][size] |= 1 << order[pos]
				}
			}
		}
	}

	rotate := [9]int{2, 5, 8, 1, 4, 7, 0, 3, 6}
	mirror := [9]int{6, 7, 8, 3, 4, 5, 0, 1, 2}
	maxId := BoardID(0)
	for range 2 {
		for range 4 {
			maxId = max(maxId, genID())
			reorder(rotate)
		}
		reorder(mirror)
	}
	if maxId == 0 {
		fmt.Printf("%v", b)
	}
	return maxId
}

func (b Board) String() string {
	return fmt.Sprintf(
		"board %v {\n"+
			"	%09b - %09b  |  %d - %d\n"+
			"	%09b - %09b  |  %d - %d\n"+
			"	%09b - %09b  |  %d - %d\n"+
			"}\n",
		b.Player,
		b.Layers[0][2], b.Layers[1][2], b.Pieces[0][2], b.Pieces[1][2],
		b.Layers[0][1], b.Layers[1][1], b.Pieces[0][1], b.Pieces[1][1],
		b.Layers[0][0], b.Layers[1][0], b.Pieces[0][0], b.Pieces[1][0],
	)
}

func (b Board) IsLeft(s Size) bool {
	return b.Pieces[b.Player][s] > 0
}

func (b Board) IsMovable(p Player, s Size, pos Pos) bool {
	for size := Size(2); size > s; size -= 1 {
		if b.Layers[0][size]&(1<<pos) != 0 || b.Layers[1][size]&(1<<pos) != 0 {
			return false
		}
	}
	return b.Layers[p][s]&(1<<pos) != 0
}

func (b Board) IsEmpty(s Size, pos Pos) bool {
	for size := Size(2); size >= s; size -= 1 {
		if b.Layers[0][size]&(1<<pos) != 0 || b.Layers[1][size]&(1<<pos) != 0 {
			return false
		}
	}
	return true
}

func (b Board) IsMoveValid(m Move) bool {
	if m.New {
		if !b.IsLeft(m.Size) {
			return false
		}
	} else {
		if !b.IsMovable(b.Player, m.Size, m.From) {
			return false
		}
	}
	if !b.IsEmpty(m.Size, m.To) {
		return false
	}

	return true
}

func (b Board) GetTop(pos Pos) (p Player, s Size, ok bool) {
	for p = 0; p < 2; p++ {
		for s = 0; s < 3; s++ {
			if b.IsMovable(p, s, pos) {
				ok = true
				return
			}
		}
	}
	return
}

func (b Board) TopView(p Player) Layer {
	np := 1 ^ p
	view := b.Layers[p][0] & ^b.Layers[np][1] |
		b.Layers[p][1] & ^b.Layers[np][2] |
		b.Layers[p][2]
	return view
}

func (b Board) GetState() State {
	isLine := func(l Layer) State {
		// horizontal
		var check = l & (l << 1) & (l << 2) & 0b100100100
		// diag 1
		check |= l & (l << 2) & (l << 4) & 0b001000000
		// vertical
		check |= l & (l << 3) & (l << 6)
		// diag 2
		check |= l & (l << 4) & (l << 8)
		if check != 0 {
			return 1
		}
		return 0
	}

	var win = isLine(b.TopView(b.Player))
	var loss = isLine(b.TopView(b.Player ^ 1))
	return win | (loss << 1)
}

func (b Board) Moves() (moves []Move) {
	for s := Size(0); s < 3; s += 1 {
		for to := Pos(0); to < 9; to += 1 {
			if !b.IsEmpty(s, to) {
				continue
			}
			if b.IsLeft(s) {
				moves = append(moves, Move{
					New:  true,
					To:   to,
					Size: s,
				})
			}
			for from := Pos(0); from < 9; from += 1 {
				if !b.IsMovable(b.Player, s, from) {
					continue
				}
				moves = append(moves, Move{
					New:  false,
					From: from,
					To:   to,
					Size: s,
				})
			}
		}
	}
	return
}

func (b *Board) DoMove(m Move) {
	if m.New {
		b.Pieces[b.Player][m.Size] -= 1
	} else {
		b.Layers[b.Player][m.Size] ^= 1 << m.From
	}

	b.Layers[b.Player][m.Size] |= 1 << m.To

	b.Player = 1 ^ b.Player
}

func (b *Board) UndoMove(m Move) {
	b.Player = 1 ^ b.Player

	b.Layers[b.Player][m.Size] ^= 1 << m.To

	if m.New {
		b.Pieces[b.Player][m.Size] += 1
	} else {
		b.Layers[b.Player][m.Size] |= 1 << m.From
	}
}
