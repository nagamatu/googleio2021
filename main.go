package main

import (
	"fmt"
	"os"
	"time"
)

type cardColumn struct {
	c    byte
	code int
}

type deck []*cardColumn

func valueOfBits(bits ...int) int {
	v := 0
	for _, bit := range bits {
		v |= 1 << bit
	}
	return v
}

func hasBit(v, bit int) bool {
	return (v & (1 << bit)) != 0
}

func codeFor(c byte) *cardColumn {
	for _, code := range cardColumns {
		if code.c == c {
			return code
		}
	}
	return nil
}

var cardColumns = []*cardColumn{
	{'0', valueOfBits(0)}, // 0
	{'1', valueOfBits(1)},
	{'2', valueOfBits(2)},
	{'3', valueOfBits(3)},
	{'4', valueOfBits(4)},
	{'5', valueOfBits(5)},
	{'6', valueOfBits(6)},
	{'7', valueOfBits(7)},
	{'8', valueOfBits(8)},
	{'9', valueOfBits(9)},

	{'A', valueOfBits(1, 11)}, // 10
	{'B', valueOfBits(2, 11)},
	{'C', valueOfBits(3, 11)},
	{'D', valueOfBits(4, 11)},
	{'E', valueOfBits(5, 11)},
	{'F', valueOfBits(6, 11)},
	{'G', valueOfBits(7, 11)},
	{'H', valueOfBits(8, 11)},
	{'I', valueOfBits(9, 11)},
	{'J', valueOfBits(1, 12)},
	{'K', valueOfBits(2, 12)},
	{'L', valueOfBits(3, 12)},
	{'M', valueOfBits(4, 12)},
	{'N', valueOfBits(5, 12)},
	{'O', valueOfBits(6, 12)},
	{'P', valueOfBits(7, 12)},
	{'Q', valueOfBits(8, 12)},
	{'R', valueOfBits(9, 12)},
	{'S', valueOfBits(0, 2)},
	{'T', valueOfBits(0, 3)},
	{'U', valueOfBits(0, 4)},
	{'V', valueOfBits(0, 5)},
	{'W', valueOfBits(0, 6)},
	{'X', valueOfBits(0, 7)},
	{'Y', valueOfBits(0, 8)},
	{'Z', valueOfBits(0, 9)},

	{'&', valueOfBits(11)},
	{'-', valueOfBits(12)},
	{':', valueOfBits(2, 8)},
	{'#', valueOfBits(3, 8)},
	{'@', valueOfBits(4, 8)},
	{'\'', valueOfBits(5, 8)},
	{'=', valueOfBits(6, 8)},
	{'"', valueOfBits(7, 8)},
	{'c', valueOfBits(2, 8, 11)},
	{'.', valueOfBits(3, 8, 11)},
	{'<', valueOfBits(4, 8, 11)},
	{'(', valueOfBits(5, 8, 11)},
	{'+', valueOfBits(6, 8, 11)},
	{'|', valueOfBits(7, 8, 11)},
	{'!', valueOfBits(2, 8, 12)},
	{'$', valueOfBits(3, 8, 12)},
	{'*', valueOfBits(4, 8, 12)},
	{')', valueOfBits(5, 8, 12)},
	{';', valueOfBits(6, 8, 12)},
	{'^', valueOfBits(7, 8, 12)},
	{',', valueOfBits(0, 3, 8)},
	{'%', valueOfBits(0, 4, 8)},
	{'_', valueOfBits(0, 5, 8)},
	{'>', valueOfBits(0, 6, 8)},
	{'?', valueOfBits(0, 7, 8)},
	{'/', valueOfBits(0, 1)}, // 61
}

func getBits(v int) []int {
	bits := []int{}
	for i := 0; i < 13; i++ {
		if hasBit(v, i) {
			bits = append(bits, i)
		}
	}
	return bits
}

func (c *cardColumn) hasNeighbor(t *cardColumn) bool {
	bits := getBits(c.code)
	for _, bit := range bits {
		if hasBit(t.code, bit) {
			return true
		}
	}
	return false
}

func (c *cardColumn) hasRoads(start int) []int {
	var roads []int = nil
	if start < 9 && hasBit(c.code, start+1) {
		roads = append(roads, start+1)
	}
	if start > 0 && hasBit(c.code, start-1) {
		roads = append(roads, start-1)
	}
	return roads
}

func (d deck) string() string {
	var s []byte = nil
	for _, c := range d {
		s = append(s, c.c)
	}
	return string(s)
}

func (d deck) hasCode(c *cardColumn) bool {
	for _, s := range d {
		if s.code == c.code {
			return true
		}
	}
	return false
}

func copyMapExcept(m map[byte]*cardColumn, c byte) map[byte]*cardColumn {
	ret := make(map[byte]*cardColumn)
	for k, v := range m {
		if k != c {
			ret[k] = v
		}
	}
	return ret
}

const (
	up             = -1
	down           = 1
	maxSwitchCount = 9
)

func getDirection(cur, next int) int {
	if cur < next {
		return down
	}
	return up
}

const logRateLimitSeconds = 1000 * 1000 * 1000

var lastLogTime time.Time = time.Now()

func (d deck) log() {
	now := time.Now()
	if now.Sub(lastLogTime) > time.Duration(logRateLimitSeconds) {
		d.showDeck()
		lastLogTime = now
	}
}

func (d deck) doCode(availables map[byte]*cardColumn, start, dir, switchCount int) (deck, map[byte]*cardColumn, error) {
	d.log()
	if len(availables) == 0 {
		return d, availables, nil
	}
	if switchCount >= maxSwitchCount {
		return nil, nil, fmt.Errorf("too many direction switch")
	}
	if len(d) >= 8 && d[7].code != codeFor('B').code {
		return nil, nil, fmt.Errorf("7th code must be 'B'")
	}
	for _, c := range availables {
		if len(d) > 0 && d[len(d)-1].hasNeighbor(c) {
			continue
		}
		if roads := c.hasRoads(start); len(roads) > 0 {
			newDeck := append(d, c)
			for _, r := range roads {
				nextDir := getDirection(start, r)
				count := switchCount
				if dir != nextDir {
					count++
				}
				dd, aa, err := newDeck.doCode(copyMapExcept(availables, c.c), r, nextDir, count)
				if err == nil {
					return dd, aa, nil
				}
			}
		}
	}
	return nil, nil, fmt.Errorf("no deck available")
}

func (d deck) showDeck() {
	fmt.Printf("%s\n", d.string())
	for _, c := range d {
		fmt.Printf("%c: ", c.c)
		bits := getBits(c.code)
		for i := 0; i < 13; i++ {
			found := false
			for _, bit := range bits {
				if i == bit {
					found = true
					break
				}
			}
			if found {
				fmt.Printf("* ")
			} else {
				fmt.Printf(". ")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	d := deck{
		codeFor('A'),
		codeFor('2'),
		codeFor('3'),
		codeFor('4'),
		codeFor('5'),
		codeFor('D'),
	}
	availableCode := make(map[byte]*cardColumn)
	for _, c := range cardColumns {
		if !d.hasCode(c) {
			availableCode[c.c] = c
		}
	}

	d, _, err := d.doCode(availableCode, 4, up, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %+v\n", err)
		return
	}
	d.showDeck()
}
