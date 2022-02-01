package models

type Position struct {
	X int
	Y int
}

func (p *Position) MoveTop() {
	p.X = p.X - 1
}
func (p *Position) GetNextTop() int {
	return p.X - 1
}

func (p *Position) MoveLeft() {
	p.Y = p.Y - 1
}

func (p *Position) GetNextLeft() int {
	return p.Y - 1
}

func (p *Position) MoveRight() {
	p.Y = p.Y + 1
}

func (p *Position) GetNextRight() int {
	return p.Y + 1
}
