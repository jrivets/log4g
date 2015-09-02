package gorivets

import (
	. "github.com/jrivets/log4g/Godeps/_workspace/src/gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type sortedSliceSuite struct {
}

var _ = Suite(&sortedSliceSuite{})

type paramType struct {
	v int
}

func (p1 *paramType) Compare(val Comparator) int {
	p2v := val.(*paramType).v
	switch {
	case p1.v < p2v:
		return -1
	case p1.v > p2v:
		return 1
	default:
		return 0
	}
}

func (s *sortedSliceSuite) TestNewSortedSlice(c *C) {
	ss, err := NewSortedSlice(-1)
	c.Assert(ss, IsNil)
	c.Assert(err, NotNil)

	ss, err = NewSortedSlice(1)
	c.Assert(ss, NotNil)
	c.Assert(err, IsNil)
}

func (s *sortedSliceSuite) TestNewSortedSliceByParams(c *C) {
	ss, err := NewSortedSliceByParams()
	c.Assert(ss, IsNil)
	c.Assert(err, NotNil)

	ss, err = NewSortedSliceByParams(make([]Comparator, 1)...)
	c.Assert(ss, NotNil)
	c.Assert(err, IsNil)

	ss, err = NewSortedSliceByParams([]Comparator{&paramType{2}, &paramType{1}}...)
	c.Assert(ss.Len(), Equals, 2)
	c.Assert(ss.At(0).(*paramType).v, Equals, 1)
	c.Assert(ss.At(1).(*paramType).v, Equals, 2)
}

func (s *sortedSliceSuite) TestLen(c *C) {
	ss, _ := NewSortedSliceByParams(&paramType{2}, &paramType{1})
	c.Assert(ss.Len(), Equals, 2)

	ss, _ = NewSortedSlice(24)
	c.Assert(ss.Len(), Equals, 0)
}

func (s *sortedSliceSuite) TestAdd(c *C) {
	ss, _ := NewSortedSlice(1)
	ss.Add(nil)
	c.Assert(ss.Len(), Equals, 0)

	ss.Add(&paramType{3})
	c.Assert(ss.Len(), Equals, 1)

	ss.Add(&paramType{1})
	c.Assert(ss.Len(), Equals, 2)
	c.Assert(ss.At(0).(*paramType).v, Equals, 1)
	c.Assert(ss.At(1).(*paramType).v, Equals, 3)
}

func (s *sortedSliceSuite) TestFind(c *C) {
	ss, _ := NewSortedSlice(1)
	ss.Add(&paramType{3})
	c.Check(ss.At(0).(*paramType).v, Equals, 3)

	idx, ok := ss.Find(&paramType{2})
	c.Assert(ok, Equals, false)
	c.Assert(idx < 0, Equals, true)

	idx, ok = ss.Find(&paramType{3})
	c.Check(ok, Equals, true)
	c.Assert(idx, Equals, 0)
}

func (s *sortedSliceSuite) TestDelete(c *C) {
	ss, _ := NewSortedSlice(1)
	ss.Add(&paramType{3})
	ss.Add(&paramType{4})

	c.Assert(ss.Delete(&paramType{2}), Equals, false)
	c.Assert(ss.Len(), Equals, 2)

	c.Assert(ss.Delete(&paramType{3}), Equals, true)
	c.Assert(ss.Len(), Equals, 1)
	c.Check(ss.At(0).(*paramType).v, Equals, 4)

	ss.Delete(&paramType{4})
	c.Assert(ss.Len(), Equals, 0)
}

func (s *sortedSliceSuite) TestDeleteAt(c *C) {
	ss, _ := NewSortedSlice(1)
	ss.Add(&paramType{3})
	ss.Add(&paramType{4})

	e := ss.DeleteAt(1).(*paramType)
	c.Assert(e.v, Equals, 4)
	c.Assert(ss.Len(), Equals, 1)
	c.Check(ss.At(0).(*paramType).v, Equals, 3)

	ss.DeleteAt(0)
	c.Assert(ss.Len(), Equals, 0)
}

func (s *sortedSliceSuite) TestCopy(c *C) {
	ss, _ := NewSortedSlice(1)
	ss.Add(&paramType{3})
	ss.Add(&paramType{4})

	slice := ss.Copy()
	c.Check(slice[0].(*paramType).v, Equals, 3)
	c.Check(slice[1].(*paramType).v, Equals, 4)
	slice[0] = &paramType{25}
	c.Check(slice[0].(*paramType).v, Equals, 25)
	c.Check(ss.At(0).(*paramType).v, Equals, 3)
}
