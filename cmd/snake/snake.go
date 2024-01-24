package main

type point [2]uint8 // {x,y}
const (
	pointX = iota
	pointY
)

func (p point) norm(p2 point) point {
	return point{
		pointX: p[pointX] % p2[pointX],
		pointY: p[pointY] % p2[pointY],
	}
}

type direction uint8

const (
	dirUp = direction(iota) + 1
	dirRight
	dirDown
	dirLeft
)

type rotation uint8

const (
	clockwise = rotation(iota)
	anticlockwise
)

func (d direction) Rotate(r rotation) direction {
	switch r {
	case clockwise:
		d += 1
		if d > 4 {
			d = 1
		}
	case anticlockwise:
		d -= 1
		if d < 1 {
			d = 4
		}
	}
	return d
}

type snake struct {
	head snakeCell
	len  uint8
	dir  direction
}

// Move moves the snake one step in the chosen direction, constrained by the
// game board.
func (s *snake) Move() {
	p := s.head.p
	s.head.Push(p) // shunt current pos into next

	switch s.dir {
	case dirUp:
		if p[pointY] == 0 {
			p[pointY] = 4
		} else {
			p[pointY] = p[pointY] - 1
		}

	case dirDown:
		if p[pointY] == 4 {
			p[pointY] = 0
		} else {
			p[pointY] = p[pointY] + 1
		}

	case dirLeft:
		if p[pointX] == 0 {
			p[pointX] = 4
		} else {
			p[pointX] = p[pointX] - 1
		}

	case dirRight:
		if p[pointX] == 4 {
			p[pointX] = 0
		} else {
			p[pointX] = p[pointX] + 1
		}
	}

	s.head.p = p
}

// Grow adds a new cell to the snake. It will remain stationary for one tick.
func (s *snake) Grow() {
	s.len++
	cell := &s.head
	for cell.next != nil {
		cell = cell.next
	}
	cell.next = &snakeCell{p: cell.p}
}

// CollidesSelf checks if the snake's head is colliding with itself.
func (s *snake) CollidesSelf() bool {
	cur := s.head.next
	for cur != nil {
		if s.head.p == cur.p {
			return true
		}
		cur = cur.next
	}
	return false
}

// Collides checks if any part of the snake is colliding with the given point.
func (s *snake) Collides(p point) bool {
	cur := &s.head
	for cur != nil {
		if cur.p == p {
			return true
		}
		cur = cur.next
	}
	return false
}

type snakeCell struct {
	p point

	next *snakeCell
}

// Push pushes the current point of this snakeCell down into the next, and
// overwrites it with the received point.
func (sc *snakeCell) Push(p point) {
	if sc.next != nil {
		sc.next.Push(sc.p)
	}
	sc.p = p
}
