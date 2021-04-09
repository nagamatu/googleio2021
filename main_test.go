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
			[]int{0, 3},
		},
		{
			codeFor('$'),
			[]int{1, 5, 10},
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
			3,
			true,
		},
		{
			codeFor('A'),
			0,
			true,
		},
		{
			codeFor('A'),
			7,
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
	d := newDeck("A2345D")
	t.Assert().Equal("A2345D", d.string())
}

func (t *mainTestSuite) Test_hasCode() {
	d := newDeck("A2345D")
	for _, c := range d {
		t.Assert().True(d.hasCode(c))
	}
	t.Assert().False(d.hasCode(codeFor('B')))
}

func (t *mainTestSuite) Test_doCode() {
	d := newDeck("A2345D")

	//A2345D$B/KC*V67,:
	availableCode := map[byte]*cardColumn{
		'$': codeFor('$'),
		'B': codeFor('B'),
		'/': codeFor('/'),
		'K': codeFor('K'),
		'C': codeFor('C'),
	}
	actual, _, err := d.doCode(availableCode, 6, up, 1)
	t.Assert().NoError(err, "A2345D$B/KC")
	t.Assert().Equal("A2345D$B/KC", actual.string(), "A2345D$B/KC")

	// A2345D$B/KC*V67,I8
	availableCode = map[byte]*cardColumn{
		'$': codeFor('$'),
		'B': codeFor('B'),
		'/': codeFor('/'),
		'K': codeFor('K'),
		'C': codeFor('C'),
		'*': codeFor('*'),
		'V': codeFor('V'),
		'6': codeFor('6'),
		'7': codeFor('7'),
		',': codeFor(','),
		'I': codeFor('I'),
		'Y': codeFor('Y'),
	}
	actual, _, err = d.doCode(availableCode, 6, up, 1)
	t.Assert().NoError(err)
	t.Assert().True(actual.string() == "A2345D,B/KC*V67YI$" ||
		actual.string() == "A2345D$B/KC*V67YI," ||
		actual.string() == "A2345D$B/KC*V67,IY" ||
		actual.string() == "A2345D,B/KC*V67$IY," ||
		actual.string() == "A2345D,B/KC*V67$IY", actual.string())

	availableCode = map[byte]*cardColumn{
		'L': codeFor('L'),
		'T': codeFor('T'),
		'@': codeFor('@'),
		'N': codeFor('N'),
		',': codeFor(','),
	}
	_, _, err = d.doCode(availableCode, 6, up, 1)
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
		t.Assert().Equal(testdata.value, testdata.c0.hasNeighbor(testdata.c1), fmt.Sprintf("%c %c", testdata.c0.c, testdata.c1.c))
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
			5,
			[]int{6},
		},
		{
			codeFor('D'),
			6,
			nil,
		},
		{
			codeFor('D'),
			7,
			[]int{6},
		},
	} {
		t.Assert().Equal(testdata.value, testdata.c.hasRoads(testdata.n), fmt.Sprintf("%c", testdata.c.c))
	}
}

func (t *mainTestSuite) Test_switchCost() {
	for _, testdata := range []struct {
		cost    int
		pos     int
		dir     int
		nextDir int
		value   int
	}{
		{
			1, 4, up, down, 11,
		},
		{
			1, 3, up, down, 2,
		},
		{
			1, 11, down, up, 2,
		},
		{
			1, 10, down, up, 11,
		},
	} {
		t.Assert().Equal(testdata.value, switchCost(testdata.cost, testdata.pos, testdata.dir, testdata.nextDir), fmt.Sprintf("pos: %d", testdata.pos))
	}

}

func (t *mainTestSuite) Test_codeStatistics() {
	bitsMap := make(map[int]int)
	for _, c := range cardColumns {
		bits := getBits(c.code)
		for _, bit := range bits {
			bitsMap[bit]++
		}
	}
	for i := 0; i <= 11; i++ {
		fmt.Printf("[%d] %d\n", i, bitsMap[i])
	}
}

func (t *mainTestSuite) Test_log() {
	d := newDeck("A2345D")
	d.log()
}

func (t *mainTestSuite) Test_showDeck() {
	d := newDeck("A2345D")
	d.showDeck()
}
