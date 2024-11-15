package gobblers

// #cgo CFLAGS: -Ofast
// #include "search.h"
import "C"
import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"
	"time"
)

const (
	maxScore    = 10000
	minScore    = -maxScore
	winScore    = 1000      // +depth
	lossScore   = -winScore // -depth
	tooFarScore = 100       // == depth
	drawScore   = 0         // +depth
)

type EvalKind int

const (
	EvalTooFar EvalKind = iota
	EvalWin
	EvalLoss
	EvalDrawMe
	EvalDrawOther
)

type Evaluation struct {
	kind  EvalKind
	depth int
	time  time.Duration
	nodes uint64
}

func boardToCBoard(b Board) C.Board {
	cboard := C.bInit()
	for sign := range 2 {
		for size := range 3 {
			cboard.layers[sign][size] = C.int(b.Layers[sign][size])
			cboard.pieces[sign][size] = C.int(b.Pieces[sign][size])
		}
	}
	cboard.player = C.int(b.Player)
	return cboard
}

func Evaluate(b Board, startDepth int, maxDepth int) (eval Evaluation) {
	cboard := boardToCBoard(b)
	cnodes := C.longlong(0)
	depth := startDepth

	start := time.Now()
	for {
		score := int(C.negamax(&cboard, C.int(minScore), C.int(maxScore), C.int(depth), &cnodes))
		if score != tooFarScore || depth >= maxDepth {
			eval.time = time.Since(start)
			eval.nodes = uint64(cnodes)

			if score == tooFarScore {
				eval.kind = EvalTooFar
				eval.depth = depth
			} else if score >= winScore {
				eval.kind = EvalWin
				eval.depth = (depth - (score - winScore))
			} else if score <= lossScore {
				eval.kind = EvalLoss
				eval.depth = (depth + (score - lossScore))
			} else if score >= drawScore {
				eval.kind = EvalDrawMe
				eval.depth = (depth - (score - drawScore))
			} else {
				eval.kind = EvalDrawOther
				eval.depth = (depth + (score - drawScore))
			}
			return
		}
		depth += 1
	}
}

type Search struct {
	mutex sync.Mutex
	cache map[BoardID]Evaluation
	evals map[BoardID]struct{}
	done  *sync.Cond
}

func NewSearch() *Search {
	var s Search
	s.done = sync.NewCond(&s.mutex)
	s.evals = make(map[BoardID]struct{})
	b, err := os.ReadFile("scorebook")
	if err != nil {
		println("Error loading scorebook: ", err.Error())
		s.cache = make(map[BoardID]Evaluation)
		return &s
	}
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	dec.Decode(&s.cache)
	return &s
}

func (s *Search) Evaluate(b Board, maxDepth int) Evaluation {
	startDepth := min(10, maxDepth)
	id := b.ID()
	s.mutex.Lock()
	if _, ok := s.evals[id]; ok {
		for {
			s.done.Wait()
			if _, ok := s.evals[id]; !ok {
				break
			}
		}
	}
	if eval, ok := s.cache[id]; ok {
		if eval.kind != EvalTooFar || eval.depth >= maxDepth {
			eval.nodes = 0
			eval.time = 0
			return eval
		}
		startDepth = max(startDepth, eval.depth)
	}
	s.evals[id] = struct{}{}
	s.mutex.Unlock()

	eval := Evaluate(b, startDepth, maxDepth)

	s.mutex.Lock()
	s.cache[id] = eval
	delete(s.evals, id)
	s.done.Broadcast()
	s.mutex.Unlock()

	return eval
}

func (s *Search) Flush() error {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s.cache)
	if err != nil {
		return err
	}
	os.WriteFile("scorebook", buf.Bytes(), 0644)
	return nil
}
