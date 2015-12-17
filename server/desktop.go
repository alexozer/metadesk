package server

import "sort"

type DesktopTree struct {
	*Desktop             // root
	wm                   *Bspwm
	focused, lastFocused *Desktop
}

func NewDesktopTree(wm *Bspwm) *Desktop {
	tree := new(DesktopTree)

	root := tree.newDesktop()
	root.id = wm.RootDesktop()
	tree.Desktop = root
	tree.focused, tree.lastFocused = root, root

	tree.wm = wm

	return root
}

func (this *DesktopTree) newDesktop() *Desktop {
	return &Desktop{
		DesktopTree: this,
		attrs:       make(map[string]string),
		children:    make([]*Desktop, 0),
	}
}

type Desktop struct {
	id string
	*DesktopTree

	parent                    *Desktop
	children                  []*Desktop
	focusedChild, parentIndex int

	attrs map[string]string
}

func (this *Desktop) ChildAt(i int) *Desktop {
	return this.children[i]
}

func (this *Desktop) NumChildren() int {
	return len(this.children)
}

func (this *Desktop) IsValidIndex(i int) bool {
	return i >= 0 && i < len(this.children)
}

func (this *Desktop) Parent() *Desktop {
	return this.parent
}

func (this *Desktop) AddChild() {
	d := this.newDesktop()
	if len(this.children) == 0 {
		d.id = this.id
		if this.focused == this {
			this.focused = d
		}
	} else {
		d.id = this.wm.AddDesktop()
	}

	d.parent = this
	d.parentIndex = len(this.children)
	this.children = append(this.children, d)
}

/* Requires: this != root, len(this.children) == 0, !this.IsOccupied() */
func (this *Desktop) Remove() {
	parent := this.parent

	left := parent.children[:this.parentIndex]
	right := parent.children[this.parentIndex+1:]
	parent.children = append(left, right...)
	if len(parent.children) == 0 {
		parent.id = this.id
		this.focused, this.lastFocused = parent, parent
		return
	}
	this.wm.RemoveDesktop(this.id)

	for _, child := range parent.children[this.parentIndex:] {
		child.parentIndex--
	}

	if parent.focusedChild > this.parentIndex {
		parent.focusedChild--
	} else if parent.focusedChild == this.parentIndex {
		if parent.focusedChild >= len(parent.children) {
			parent.focusedChild--
		}
	}

	if this.isLeafFocused() {
		parent.focus()
	}

	if this.lastFocused == this {
		this.lastFocused = this.focused
	}
}

func (this *Desktop) IsOccupied() bool {
	return this.wm.IsDesktopOccupied(this.id)
}

func (this *Desktop) Focus() {
	this.focusIndices()
	if !this.focus() {
		return
	}

	if !this.lastFocused.IsOccupied() &&
		this.lastFocused.parent != nil &&
		len(this.lastFocused.parent.children) > 1 {

		this.lastFocused.Remove()
		this.lastFocused = this.focused
	}
}

func (this *Desktop) FocusNext() {
	if this.focusedChild == len(this.children)-1 {
		if this.children[this.focusedChild].IsOccupied() {
			this.AddChild()
		} else {
			return
		}
	}
	this.focusedChild++
	this.Focus()
}

func (this *Desktop) FocusPrev() {
	if this.focusedChild > 0 {
		this.focusedChild--
		this.Focus()
	}
}

func (this *Desktop) FocusedChild() int {
	return this.focusedChild
}

func (this *Desktop) Focused() *Desktop {
	return this.focused
}

func (this *Desktop) LastFocused() *Desktop {
	return this.lastFocused
}

func (this *Desktop) ClaimFocusedWindow() {
	this.wm.ClaimFocusedWindow(this.id)
}

func (this *Desktop) SwapNext() {
	index := this.parentIndex + 1
	if this.parent.IsValidIndex(index) {
		this.SwapWith(index)
	}
}

func (this *Desktop) SwapPrev() {
	index := this.parentIndex - 1
	if this.parent.IsValidIndex(index) {
		this.SwapWith(index)
	}
}

// Assumes: this.IsValidIndex(index) && this.parent != nil
func (this *Desktop) SwapWith(index int) {
	parent := this.parent

	// fix the parent's focused pointer
	if parent.focusedChild == this.parentIndex {
		parent.focusedChild = index
	} else if this.parentIndex == index {
		parent.focusedChild = this.parentIndex
	}

	// swap the parent's children
	swapee := parent.children[index]
	parent.children[this.parentIndex], parent.children[index] = swapee, this
	this.parentIndex, swapee.parentIndex = swapee.parentIndex, this.parentIndex
}

func (this *Desktop) focusIndices() {
	if this.parent != nil {
		this.parent.focusedChild = this.parentIndex
		this.parent.focusIndices()
	}
}

func (this *Desktop) focus() bool {
	if len(this.children) == 0 {
		if this == this.focused {
			return false
		}

		this.wm.FocusDesktop(this.id)
		this.lastFocused, this.focused = this.focused, this
	} else {
		this.children[this.focusedChild].focus()
	}

	return true
}

func (this *Desktop) isLeafFocused() bool {
	if len(this.children) > 0 {
		return this.children[this.focusedChild].isLeafFocused()
	} else {
		return this == this.focused
	}
}

func (this *Desktop) Attr(name string) string {
	return this.attrs[name]
}

func (this *Desktop) SetAttr(name, value string) {
	this.attrs[name] = value
}

func (this *Desktop) UnsetAttr(name string) {
	delete(this.attrs, name)
}

type AttrList [][2]string

func (this AttrList) Len() int {
	return len(this)
}

func (this AttrList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this AttrList) Less(i, j int) bool {
	return this[i][0] < this[j][0]
}

func (this *Desktop) AttrList() AttrList {
	list := make(AttrList, 0, len(this.attrs))
	for key, val := range this.attrs {
		list = append(list, [2]string{key, val})
	}
	sort.Sort(list)

	return list
}
