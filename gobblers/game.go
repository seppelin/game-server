package gobblers

type Selection int

const (
	SelNone Selection = iota
	SelFrom
	SelMove
)

type Game struct {
	history []Move
	board   Board
	sel     Selection
	move    Move
}

func NewGame() *Game {
	return &Game{
		board: NewBoard(),
		sel:   SelNone,
	}
}

func (g Game) IsSelNew(p Player, s Size) bool {
	if g.sel != SelNone && p == g.board.Player {
		if g.move.New && g.move.Size == s {
			return true
		}
	}
	return false
}

func (g Game) IsSelBoard(pos Pos) bool {
	if g.sel != SelNone {
		if !g.move.New && g.move.From == pos {
			return true
		}
		if g.sel == SelMove && g.move.To == pos {
			return true
		}
	}
	return false
}

func (g Game) IsActionNew(p Player, s Size) bool {
	return g.sel == SelNone && p == g.board.Player && g.board.IsLeft(s)
}

func (g Game) IsActionBoard(pos Pos) bool {
	if g.sel == SelNone {
		p, _, ok := g.board.GetTop(pos)
		return ok && p == g.board.Player
	} else {
		return g.board.IsEmpty(g.move.Size, pos)
	}
}

func (g Game) Board() Board {
	return g.board
}

func (g Game) Selection() (Selection, Move) {
	return g.sel, g.move
}

func (g *Game) DoMove(m Move) bool {
	if g.board.IsMoveValid(m) {
		g.board.DoMove(m)
		g.history = append(g.history, m)
		return true
	} else {
		return false
	}
}

func (g *Game) UndoMove(m Move) bool {
	i := len(g.history)
	if i > 0 {
		g.board.UndoMove(g.history[i-1])
		g.history = g.history[:i-1]
		return true
	} else {
		return false
	}
}

func (g *Game) SelectNew(s Size) {
	if g.sel == SelNone {
		if g.board.IsLeft(s) {
			g.sel = SelFrom
			g.move.New = true
			g.move.Size = s
		}
	} else {
		g.sel = SelNone
	}
}

func (g *Game) SelectBoard(pos Pos) {
	if g.sel == SelNone {
		p, s, ok := g.board.GetTop(pos)
		if ok && p == g.board.Player {
			g.sel = SelFrom
			g.move.New = false
			g.move.Size = s
			g.move.From = pos
		}
	} else {
		if g.board.IsEmpty(g.move.Size, g.move.From) {
			g.sel = SelMove
			g.move.To = pos
		} else {
			g.sel = SelNone
		}
	}
}
