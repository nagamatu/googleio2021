package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type mainTestSuite struct {
	suite.Suite
}

func Test_mainTestSuite(t *testing.T) {
	suite.Run(t, new(mainTestSuite))
}

func (t *mainTestSuite) Test_checkCardCode() {
	t.Assert().Equal(62, len(cardColumns))

	for i := range cardColumns {
		for j := i + 1; j < len(cardColumns); j++ {
			t.Assert().NotEqual(cardColumns[i], cardColumns[j])
		}
	}
}

func (t *mainTestSuite) Test_checkNoNeighbor() {
	for _, testdata := range []struct {
		deck  deck
		value bool
	}{
		{
			[]*cardColumn{
				codeFor('A'),
				codeFor('B'),
			},
			false,
		},
		{
			[]*cardColumn{
				codeFor('A'),
				codeFor('2'),
				codeFor('3'),
				codeFor('4'),
				codeFor('5'),
				codeFor('D'),
			},
			true,
		},
		{
			[]*cardColumn{
				codeFor('A'),
			},
			true,
		}} {
		t.Assert().Equal(testdata.value, checkNoNeighbor(testdata.deck), testdata.deck.string())
	}
}

func (t *mainTestSuite) Test_checkRoad() {
	for _, testdata := range []struct {
		deck  deck
		value bool
	}{
		{
			[]*cardColumn{
				codeFor('A'),
				codeFor('3'),
				codeFor('4'),
				codeFor('5'),
				codeFor('D'),
			},
			false,
		},
		{
			[]*cardColumn{
				codeFor('A'),
				codeFor('2'),
				codeFor('3'),
				codeFor('4'),
				codeFor('5'),
				codeFor('D'),
			},
			true,
		},
		{
			[]*cardColumn{
				codeFor('A'),
				codeFor('2'),
				codeFor('3'),
				codeFor('4'),
				codeFor('5'),
				codeFor('D'),
				codeFor('B'),
			},
			false,
		},
		{
			[]*cardColumn{
				codeFor('A'),
			},
			true,
		}} {
		t.Assert().Equal(testdata.value, checkRoad(testdata.deck, 1), testdata.deck.string())
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
