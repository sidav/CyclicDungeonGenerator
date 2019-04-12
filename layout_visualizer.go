package main

import (
	"CyclicDungeonGenerator/layout_generation"
	"fmt"
	cw "github.com/sidav/golibrl/console/tcell_console"
	rnd "github.com/sidav/golibrl/random"
)

func doLayoutVisualization() {
	key := "none"
	desiredPatternNum := -1

	for key != "ESCAPE" {
		cw.Clear_console()
		pattNum := rnd.Random(layout_generation.GetTotalPatternsNumber())
		if desiredPatternNum != -1 {
			pattNum = desiredPatternNum
		}
		generatedMap, genRestarts := layout_generation.Generate(pattNum, W, H)

		if generatedMap == nil {
			cw.PutString(":(", 0, 0)
			cw.PutString(fmt.Sprintf("Generation failed even after %d restarts, pattern #%d", genRestarts, pattNum), 0, 1)
			cw.PutString("Press ENTER to generate again or ESCAPE to exit.", 0, 2)
			cw.Flush_console()
			for key != "ESCAPE" && key != "ENTER" {
				key = cw.ReadKey()
			}
			continue
		} else {
			putMap(generatedMap)
			putMiniMapAndPatternNumberAndNumberOfTries(generatedMap, pattNum, desiredPatternNum, genRestarts)
		}
		cw.Flush_console()
	keyread:
		for {
			key = cw.ReadKey()
			switch key {
			case "=":
				if desiredPatternNum < layout_generation.GetTotalPatternsNumber()-1 {
					desiredPatternNum++
				}
				break keyread
			case "-":
				if desiredPatternNum > -1 {
					desiredPatternNum--
				}
				break keyread
			case " ", "ESCAPE":
				break keyread
			}
		}
	}
}

func putCharArray(x, y int, c *[][]rune) {
	for i := 0; i < len(*c); i++ {
		for j := 0; j < len((*c)[0]); j++ {
			setcolorForRune((*c)[i][j])
			cw.PutChar((*c)[i][j], x+i, y+j)
		}
	}
}

func putMap(a *layout_generation.LayoutMap) {
	putCharArray(0, 0, a.WholeMapToCharArray())
}

func putMiniMapAndPatternNumberAndNumberOfTries(a *layout_generation.LayoutMap, pattNum, desiredPNum, restarts int) {
	sx, sy := a.GetSize()
	for y := 0; y < sy; y++ {
		for x := 0; x < sx; x++ {
			chr := a.GetCharOfElementAtCoords(x, y)
			setcolorForRune(chr)
			cw.PutChar(chr, x+sx*5+2, y)
		}
	}
	cw.SetFgColor(cw.BEIGE)
	cw.PutString(fmt.Sprintf("PATTERN SELECTED: #%d  ", desiredPNum), sx*5+2, sy+2)
	cw.PutString(fmt.Sprintf("PATTERN USED: #%d  ", pattNum), sx*5+2, sy+3)
	cw.PutString(fmt.Sprintf("Generation finised with %d restarts", restarts), sx*5+2, sy+4)

}

func setcolorForRune(chr rune) {
	switch chr {
	case '1', '2', '3', '4', '5', '6':
		cw.SetFgColor(cw.DARK_CYAN)
	case '.':
		cw.SetFgColor(cw.BEIGE)
	case '+':
		cw.SetFgColor(cw.DARK_MAGENTA)
	case '#':
		cw.SetFgColor(cw.DARK_GRAY)
	case '%':
		cw.SetFgColor(cw.RED)
	case '=':
		cw.SetFgColor(cw.DARK_BLUE)
	default:
		cw.SetFgColor(cw.DARK_GREEN)
	}
}