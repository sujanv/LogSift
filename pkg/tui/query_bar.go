package tui

type QueryBar struct {
	Active  bool
	Query   string
	Cursor  int
	History []string
}

func NewQueryBar() QueryBar {
	return QueryBar{
		History: make([]string, 0, 16),
	}
}

func (qb *QueryBar) Toggle() bool {
	qb.Active = !qb.Active
	return qb.Active
}
