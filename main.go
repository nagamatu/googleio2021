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
	{'0', valueOfBits(2)}, // 0
	{'1', valueOfBits(3)},
	{'2', valueOfBits(4)},
	{'3', valueOfBits(5)},
	{'4', valueOfBits(6)},
	{'5', valueOfBits(7)},
	{'6', valueOfBits(8)},
	{'7', valueOfBits(9)},
	{'8', valueOfBits(10)},
	{'9', valueOfBits(11)},

	{'A', valueOfBits(3, 0)}, // 10
	{'B', valueOfBits(4, 0)},
	{'C', valueOfBits(5, 0)},
	{'D', valueOfBits(6, 0)},
	{'E', valueOfBits(7, 0)},
	{'F', valueOfBits(8, 0)},
	{'G', valueOfBits(9, 0)},
	{'H', valueOfBits(10, 0)},
	{'I', valueOfBits(11, 0)},
	{'J', valueOfBits(3, 1)},
	{'K', valueOfBits(4, 1)},
	{'L', valueOfBits(5, 1)},
	{'M', valueOfBits(6, 1)},
	{'N', valueOfBits(7, 1)},
	{'O', valueOfBits(8, 1)},
	{'P', valueOfBits(9, 1)},
	{'Q', valueOfBits(10, 1)},
	{'R', valueOfBits(11, 1)},
	{'S', valueOfBits(2, 4)},
	{'T', valueOfBits(2, 5)},
	{'U', valueOfBits(2, 6)},
	{'V', valueOfBits(2, 7)},
	{'W', valueOfBits(2, 8)},
	{'X', valueOfBits(2, 9)},
	{'Y', valueOfBits(2, 10)},
	{'Z', valueOfBits(2, 11)},

	{'&', valueOfBits(0)},
	{'-', valueOfBits(1)},
	{':', valueOfBits(4, 10)},
	{'#', valueOfBits(5, 10)},
	{'@', valueOfBits(6, 10)},
	{'\'', valueOfBits(7, 10)},
	{'=', valueOfBits(8, 10)},
	{'"', valueOfBits(9, 10)},
	{'c', valueOfBits(4, 10, 0)},
	{'.', valueOfBits(5, 10, 0)},
	{'<', valueOfBits(6, 10, 0)},
	{'(', valueOfBits(7, 10, 0)},
	{'+', valueOfBits(8, 10, 0)},
	{'|', valueOfBits(9, 10, 0)},
	{'!', valueOfBits(4, 10, 1)},
	{'$', valueOfBits(5, 10, 1)},
	{'*', valueOfBits(6, 10, 1)},
	{')', valueOfBits(7, 10, 1)},
	{';', valueOfBits(8, 10, 1)},
	{'^', valueOfBits(9, 10, 1)},
	{',', valueOfBits(2, 5, 10)},
	{'%', valueOfBits(2, 6, 10)},
	{'_', valueOfBits(2, 7, 10)},
	{'>', valueOfBits(2, 8, 10)},
	{'?', valueOfBits(2, 9, 10)},
	{'/', valueOfBits(2, 3)}, // 61
}

func getBits(v int) []int {
	bits := []int{}
	for i := 0; i < 12; i++ {
		if hasBit(v, i) {
			bits = append(bits, i)
		}
	}
	return bits
}

func (c *cardColumn) hasNeighbor(t *cardColumn) bool {
	return (c.code & t.code) != 0
}

func (c *cardColumn) hasRoads(start int) []int {
	var roads []int = nil
	if start < 11 && hasBit(c.code, start+1) {
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

var lastLogTime time.Time = time.Now().Add(time.Duration(-logRateLimitSeconds))

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
					if nextDir == down {
						if start >= 4 {
							count += 10
						} else {
							count++
						}
					} else {
						if start == 11 {
							count++
						} else {
							count += 10
						}
					}
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
		for i := 0; i < 12; i++ {
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

	d, _, err := d.doCode(availableCode, 6, up, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %+v\n", err)
		return
	}
	d.showDeck()
}
