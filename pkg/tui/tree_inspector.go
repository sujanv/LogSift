package tui

type TreeNode struct {
	Key      string
	Value    string
	Expanded bool
	Children []*TreeNode
}

type TreeInspector struct {
	Visible bool
	Root    *TreeNode
}

func (ti *TreeInspector) Toggle() bool {
	ti.Visible = !ti.Visible
	return ti.Visible
}
