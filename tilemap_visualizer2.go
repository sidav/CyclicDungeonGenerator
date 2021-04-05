package main

import (
	cw "CyclicDungeonGenerator/console_wrapper"
	"CyclicDungeonGenerator/layout_generation"
	"CyclicDungeonGenerator/layout_to_tilemap"
	"CyclicDungeonGenerator/layout_to_tiles2"
	"CyclicDungeonGenerator/random"
	"fmt"
)

type vis struct {}

func (g *vis) doTilemapVisualization() {
	roomSize := 3
	key := "none"
	desiredPatternNum := -1
	rnd := random.FibRandom{}
	rnd.InitDefault()
	layout_to_tilemap.Random = &rnd
	parser := layout_generation.PatternParser{}
	filenames := parser.ListPatternFilenamesInPath("patterns/")

	for key != "ESCAPE" {
		cw.Clear_console()
		pattNum := rnd.Rand(len(filenames))
		if desiredPatternNum != -1 {
			pattNum = desiredPatternNum
		}
		gen := layout_generation.InitCyclicGenerator(true, W, H, -1)
		gen.TriesForPattern = 100
		generatedMap, genRestarts := gen.GenerateLayout(parser.ParsePatternFile(filenames[pattNum]))

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
			g.putTileMap(generatedMap, roomSize)
			// putMiniMapAndPatternNumberAndNumberOfTries(generatedMap, pattNum, desiredPatternNum, genRestarts)
		}
		cw.Flush_console()
	keyread:
		for {
			key = cw.ReadKey()
			switch key {
			case "=", "+":
				if desiredPatternNum < len(filenames)-1 {
					desiredPatternNum++
				}
				break keyread
			case "-":
				if desiredPatternNum > -1 {
					desiredPatternNum--
				}
				break keyread
			case "b":
				roomSize++
				break keyread
			case " ", "ESCAPE":
				break keyread
			}
		}
	}
}

func (g *vis) putTileMap(a *layout_generation.LayoutMap, roomSize int) {
	cw.Clear_console()
	g.putTileArray(layout_to_tiles2.MakeCharmap(roomSize, a), 0, 0)
	roomSize+=1
	rw, rh := a.GetSize()
	for rx := 0; rx < rw; rx++ {
		for ry := 0; ry < rh; ry++ {
			node := a.GetElement(rx, ry)
			conns := node.GetAllConnectionsCoords()
			if len(conns) > 0 {
				cw.SetFgColor(cw.GREEN)
				if node.IsNode() {
					name := node.GetName()
					namelen := len(name)
					offset := roomSize / 2 - namelen / 2
					cw.PutString(name, rx*roomSize + offset, ry*roomSize+roomSize/2)
				}
			}
		}
	}
}

func (g *vis) putTileArray(arr [][]rune, sx, sy int) {
	for x :=0; x <len(arr); x++{
		for y :=0; y <len((arr)[x]); y++ {
			chr := (arr)[x][y]
			setcolorForRune(chr)
			cw.PutChar(chr, sx+x, sy+y)
			cw.SetColor(cw.WHITE, cw.BLACK)
		}
	}
}
