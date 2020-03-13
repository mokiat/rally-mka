package graphics

import "github.com/mokiat/go-whiskey/math"

type DepthFunc int

const (
	DepthFuncLess DepthFunc = iota
	DepthFuncLessOrEqual
)

func createSequence(items []Item) Sequence {
	return Sequence{
		items: items,
	}
}

type Sequence struct {
	items          []Item
	itemStartIndex int
	itemEndIndex   int

	BackgroundColor  math.Vec4
	TestDepth        bool
	ClearColor       bool
	ClearDepth       bool
	WriteDepth       bool
	DepthFunc        DepthFunc
	ProjectionMatrix math.Mat4x4
	ViewMatrix       math.Mat4x4
}

func (s *Sequence) BeginItem() *Item {
	if s.itemEndIndex == len(s.items) {
		panic("max number of render items reached")
	}
	item := &s.items[s.itemEndIndex]
	s.itemEndIndex++
	item.reset()
	return item
}

func (s *Sequence) EndItem(item *Item) {
}

func (s *Sequence) reset(index int) {
	s.itemStartIndex = index
	s.itemEndIndex = index
	s.WriteDepth = true
	s.DepthFunc = DepthFuncLess
	s.TestDepth = true
}

func (s *Sequence) itemsView() []Item {
	return s.items[s.itemStartIndex:s.itemEndIndex]
}
