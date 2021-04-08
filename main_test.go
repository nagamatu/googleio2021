package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type mainTestSuite struct {
	suite.Suite
}

func Test_mainTestSuite(t *testing.T) {
	suite.Run(t, new(mainTestSuite))
}

func (t *mainTestSuite) Test_getBits() {
	for _, testdata := range []struct {
		c    *cardColumn
		bits []int
	}{
		{
			codeFor('A'),
			[]int{1, 11},
		},
		{
			codeFor('$'),
			[]int{3, 8, 12},
		},
	} {
		t.Assert().Equal(testdata.bits, getBits(testdata.c.code))
	}
}

func (t *mainTestSuite) Test_hasBit() {
	for _, testdata := range []struct {
		c     *cardColumn
		bit   int
		value bool
	}{
		{
			codeFor('A'),
			1,
			true,
		},
		{
			codeFor('A'),
			11,
			true,
		},
		{
			codeFor('A'),
			5,
			false,
		},
	} {
		t.Assert().Equal(testdata.value, hasBit(testdata.c.code, testdata.bit))
	}
}

func (t *mainTestSuite) Test_codeFor() {
	for _, c := range cardColumns {
		t.Assert().Equal(c.c, codeFor(c.c).c, fmt.Sprintf("%c: %d", c.c, c.code))
		t.Assert().Equal(c.code, codeFor(c.c).code, fmt.Sprintf("%c: %d", c.c, c.code))
	}
	t.Assert().Nil(codeFor('a'))
}

func (t *mainTestSuite) Test_string() {
	d := deck{
		codeFor('A'),
		codeFor('2'),
		codeFor('3'),
		codeFor('4'),
		codeFor('5'),
		codeFor('D'),
	}
	t.Assert().Equal("A2345D", d.string())
}

func (t *mainTestSuite) Test_hasCode() {
	d := deck{
		codeFor('A'),
		codeFor('2'),
		codeFor('3'),
		codeFor('4'),
		codeFor('5'),
		codeFor('D'),
	}
	for _, c := range d {
		t.Assert().True(d.hasCode(c))
	}
	t.Assert().False(d.hasCode(codeFor('B')))
}

func (t *mainTestSuite) Test_doCode() {
	d := deck{
		codeFor('A'),
		codeFor('2'),
		codeFor('3'),
		codeFor('4'),
		codeFor('5'),
		codeFor('D'),
	}

	//A2345D#B/!1KT
	availableCode := map[byte]*cardColumn{
		'L': codeFor('L'),
		'B': codeFor('B'),
		'#': codeFor('#'),
		'@': codeFor('@'),
		'N': codeFor('N'),
	}
	actual, _, err := d.doCode(availableCode, 4, up, 1)
	t.Assert().NoError(err, "A2345D#BL@N")
	t.Assert().Equal("A2345D#BL@N", actual.string(), "A2345D#BL@N")

	availableCode = map[byte]*cardColumn{
		'L': codeFor('L'),
		'T': codeFor('T'),
		'@': codeFor('@'),
		'N': codeFor('N'),
		',': codeFor(','),
	}
	_, _, err = d.doCode(availableCode, 4, up, 1)
	t.Assert().Error(err, "No deck found")
}

func (t *mainTestSuite) Test_checkCardCode() {
	t.Assert().Equal(62, len(cardColumns))

	for i := range cardColumns {
		for j := i + 1; j < len(cardColumns); j++ {
			t.Assert().NotEqual(cardColumns[i].c, cardColumns[j].c)
			t.Assert().NotEqual(cardColumns[i].code, cardColumns[j].code)
			t.Assert().NotEqual(cardColumns[i], cardColumns[j])
		}
	}
}

func (t *mainTestSuite) Test_hasNeighbor() {
	for _, testdata := range []struct {
		c0    *cardColumn
		c1    *cardColumn
		value bool
	}{
		{
			codeFor('D'),
			codeFor('S'),
			false,
		},
		{
			codeFor('B'),
			codeFor('S'),
			true,
		},
		{
			codeFor('"'),
			codeFor('|'),
			true,
		},
		{
			codeFor('|'),
			codeFor('*'),
			true,
		},
	} {
		t.Assert().Equal(testdata.value, testdata.c0.hasNeighbor(testdata.c1))
	}
}

func (t *mainTestSuite) Test_hasRoads() {
	for _, testdata := range []struct {
		c     *cardColumn
		n     int
		value []int
	}{
		{
			codeFor('D'),
			3,
			[]int{4},
		},
		{
			codeFor('D'),
			4,
			nil,
		},
		{
			codeFor('D'),
			5,
			[]int{4},
		},
	} {
		t.Assert().Equal(testdata.value, testdata.c.hasRoads(testdata.n))
	}
}
