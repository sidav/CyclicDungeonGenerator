package layout_generation

import (
	"CyclicDungeonGenerator/random"
	"strconv"
)

type LayoutMap struct {
	// room map with connections or smth
	elements       [][]*element
	rnd            *random.FibRandom
	randomizePaths bool
}

func (r *LayoutMap) init(sizex, sizey int, rnd *random.FibRandom, randomizePaths bool) {
	r.elements = make([][]*element, sizex)
	for i := range r.elements {
		r.elements[i] = make([]*element, sizey)
	}
	for x := 0; x < sizex; x++ {
		for y := 0; y < sizey; y++ {
			r.elements[x][y] = &element{connections: map[string]*connection{"north": nil, "east": nil, "south": nil, "west": nil}}
		}
	}
	r.rnd = rnd
	r.randomizePaths = randomizePaths
}

func (r *LayoutMap) placeNodeAtCoords(x, y int, nodeName string) {
	r.elements[x][y].nodeInfo = &nodeCell{nodeName: nodeName}
}

func (r *LayoutMap) placePathAtCoords(x, y int, pathNum int) {
	r.elements[x][y].pathInfo = &pathCell{pathNum}
}

func (r *LayoutMap) placeObstacleAtCoords(x, y int) {
	r.elements[x][y].isObstacle = true
}

func (r *LayoutMap) removeAllObstacles() {
	for x := 0; x < len(r.elements); x++ {
		for y := 0; y < len(r.elements[0]); y++ {
			r.elements[x][y].isObstacle = false
		}
	}
}

func (r *LayoutMap) getRandomPathCoordsAndRandomCellNearPath(pathNum int, allowNearNode bool) (int, int, int, int) {
	const tries = 10
	for try := 0; try < tries; try++ {
		px, py := r.getRandomPathCellCoords(pathNum, allowNearNode)
		if px == -1 && py == -1 {
			continue
		}
		for try2 := 0; try2 < tries; try2++ {
			x, y := r.rnd.RandInRange(px-1, px+1), r.rnd.RandInRange(py-1, py+1)
			if (px-x)*(py-y) != 0 { // diagonal direction is restricted
				continue
			}
			if r.areCoordsValid(x, y) && r.elements[x][y].isEmpty() {
				return px, py, x, y
			}
		}
	}
	return -1, -1, -1, -1
}

func (r *LayoutMap) getRandomNonEmptyCoordsAndRandomCellNearIt() (int, int, int, int) {
	px, py := r.getRandomNonEmptyCellCoords(1)
	if px == -1 && py == -1 {
		return -1, -1, -1, -1
	}
	x, y := r.getRandomEmptyCellNearCoords(px, py)
	if x == -1 && y == -1 {
		return -1, -1, -1, -1
	}
	return px, py, x, y
}

func (r *LayoutMap) getRandomEmptyCellCoords(minEmptyCellsNear int) (int, int) {
	emptiesX := make([]int, 0)
	emptiesY := make([]int, 0)
	for x := 0; x < len(r.elements); x++ {
		for y := 0; y < len(r.elements[0]); y++ {
			if r.elements[x][y].isEmpty() && (r.countEmptyCellsNear(x, y) >= minEmptyCellsNear) {
				emptiesX = append(emptiesX, x)
				emptiesY = append(emptiesY, y)
			}
		}
	}
	if len(emptiesX) == 0 {
		return -1, -1
	}
	index := r.rnd.Rand(len(emptiesX))
	return emptiesX[index], emptiesY[index]
}

func (r *LayoutMap) getRandomEmptyCellCoordsInRange(fx, fy, tx, ty, minEmptyCellsNear int) (int, int) { // range inclusive
	emptiesX := make([]int, 0)
	emptiesY := make([]int, 0)
	if fx > tx {
		t := tx
		tx = fx
		fx = t
	}
	if fy > ty {
		t := ty
		ty = fy
		fy = t
	}
	for x := fx; x <= tx; x++ {
		for y := fy; y <= ty; y++ {
			if r.elements[x][y].isEmpty() && (r.countEmptyCellsNear(x, y) >= minEmptyCellsNear) {
				emptiesX = append(emptiesX, x)
				emptiesY = append(emptiesY, y)
			}
		}
	}
	if len(emptiesX) == 0 {
		return -1, -1
	}
	index := r.rnd.Rand(len(emptiesX))
	return emptiesX[index], emptiesY[index]
}

func (r *LayoutMap) getRandomEmptyCellNearCoords(nx, ny int) (int, int) {
	emptiesX := make([]int, 0)
	emptiesY := make([]int, 0)
	for x := nx - 1; x <= nx+1; x++ {
		for y := ny - 1; y <= ny+1; y++ {
			if (nx-x)*(ny-y) != 0 { // restrict diagonals
				continue
			}
			if (x != nx || y != ny) && r.areCoordsValid(x, y) && r.elements[x][y].isEmpty() {
				emptiesX = append(emptiesX, x)
				emptiesY = append(emptiesY, y)
			}
		}
	}
	if len(emptiesX) == 0 {
		return -1, -1
	}
	index := r.rnd.Rand(len(emptiesX))
	return emptiesX[index], emptiesY[index]
}

func (r *LayoutMap) getRandomNonEmptyCellCoords(minEmptyCellsNear int) (int, int) {
	nonEmptiesX := make([]int, 0)
	nonEmptiesY := make([]int, 0)
	for x := 0; x < len(r.elements); x++ {
		for y := 0; y < len(r.elements[0]); y++ {
			if !r.elements[x][y].isEmpty() && r.countEmptyCellsNear(x, y) >= minEmptyCellsNear {
				nonEmptiesX = append(nonEmptiesX, x)
				nonEmptiesY = append(nonEmptiesY, y)
			}
		}
	}
	if len(nonEmptiesX) == 0 {
		return -1, -1
	}
	index := r.rnd.Rand(len(nonEmptiesX))
	return nonEmptiesX[index], nonEmptiesY[index]
}

func (r *LayoutMap) getRandomPathCellCoords(desiredPathNum int, nodesAllowed bool) (int, int) { // desiredPathNum -1 means any path
	pathsX := make([]int, 0)
	pathsY := make([]int, 0)
	for x := 0; x < len(r.elements); x++ {
		for y := 0; y < len(r.elements[0]); y++ {
			if !r.elements[x][y].isPartOfAPath() {
				continue
			}
			if !nodesAllowed && r.elements[x][y].nodeInfo != nil { // don't take nodes unless allowed
				continue
			}
			if desiredPathNum > -1 && desiredPathNum != r.elements[x][y].pathInfo.pathNumber { // don't take cells of non-desired path numbers
				continue
			}
			pathsX = append(pathsX, x)
			pathsY = append(pathsY, y)
		}
	}
	if len(pathsX) == 0 {
		return -1, -1
	}
	index := r.rnd.Rand(len(pathsX))
	return pathsX[index], pathsY[index]
}

func (r *LayoutMap) areCoordsEmpty(x, y int) bool {
	return r.elements[x][y].isEmpty()
}

func (r *LayoutMap) isPathPresentAtCoords(x, y int) bool {
	return r.elements[x][y].isPartOfAPath()
}

func (r *LayoutMap) areCoordsEmptyOrPathOnly(x, y int) bool {
	return r.elements[x][y].isPathOrEmpty()
}

func (r *LayoutMap) countEmptyCellsNear(x, y int) int {
	count := 0
	w, h := r.GetSize()
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i*j != 0 || i == 0 && j == 0 {
				continue
			}
			if x+i < 0 || x+i >= w || y+j < 0 || y+j >= h {
				continue
			}
			if r.elements[x+i][y+j].isEmpty() {
				count++
			}
		}
	}
	return count
}

func (r *LayoutMap) getCoordsOfNode(nodeName string) (int, int) {
	for x := 0; x < len(r.elements); x++ {
		for y := 0; y < len(r.elements[0]); y++ {
			if r.elements[x][y].IsNode() && r.elements[x][y].nodeInfo.nodeName == nodeName {
				return x, y
			}
		}
	}
	panic("getCoordsOfNode failed with node " + nodeName)
	return -1, -1
}

func (r *LayoutMap) areCoordsValid(x, y int) bool {
	w, h := r.GetSize()
	return x >= 0 && x < w && y >= 0 && y < h
}

// exported

func (r *LayoutMap) GetSize() (int, int) {
	return len(r.elements), len(r.elements[0])
}

func (r *LayoutMap) GetElement(x, y int) *element {
	return r.elements[x][y]
}

func (r *LayoutMap) getPassabilityMapForPathfinder(pathsArePassable bool) *[][]int {
	const (
		minRandomCostIncrease = -100
		maxRandomCostIncrease = 10000
	)
	layoutWidth, layoutHeight := r.GetSize()
	pmap := make([][]int, layoutWidth)
	for i := range pmap {
		pmap[i] = make([]int, layoutHeight)
	}

	for x := 0; x < layoutWidth; x++ {
		for y := 0; y < layoutHeight; y++ {
			if r.areCoordsEmpty(x, y) || pathsArePassable && r.areCoordsEmptyOrPathOnly(x, y) {
				pmap[x][y] = 1
				if r.isPathPresentAtCoords(x, y) {
					pmap[x][y] += maxRandomCostIncrease
				}
				// TODO: think how to better randomize path costs
				if r.randomizePaths {
					// lowering the "from" increases path randomness, but also makes the generator to fail more frequently
					// because it increases the probability for creating a non-existing path
					// "* 10" is to compensate the heuristics in the pathfinder
					pmap[x][y] += r.rnd.RandInRange(minRandomCostIncrease, maxRandomCostIncrease) * 10
				}
			} else {
				pmap[x][y] = -1
			}
		}
	}
	return &pmap
}

// output TODO: remove

func (r *LayoutMap) GetCharOfElementAtCoords(x, y int) rune { // just for rendering, TODO: remove
	elem := r.elements[x][y]
	// rune := '?'
	if elem.isEmpty() {
		return '.'
	}
	if elem.isObstacle {
		return '#'
	}
	if elem.IsNode() {
		if elem.GetName() == "" {
			return 'R'
		}
		return rune(elem.nodeInfo.nodeName[0])
	}
	if elem.isPartOfAPath() {
		number := elem.pathInfo.pathNumber
		return rune(strconv.Itoa(number)[0])
	}
	return '?'
}

func (r *LayoutMap) CellToCharArray(cellx, celly int) [][]rune {
	e := r.elements[cellx][celly]
	ca := make([][]rune, 5)
	for i := range ca {
		ca[i] = make([]rune, 5)
	}

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			ca[x][y] = '#'
		}
	}
	// draw node
	if e.nodeInfo != nil {
		for x := 1; x < 4; x++ {
			for y := 1; y < 4; y++ {
				ca[x][y] = ' '
			}
		}
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 2; y++ {
				conn := e.GetConnectionByCoords(x, y)
				if conn != nil {
					if conn.IsLocked {
						if conn.LockNum == 0 {
							ca[2+x*2][2+y*2] = '%'
						} else {
							ca[2+x*2][2+y*2] = '='
						}
					} else {
						ca[2+x*2][2+y*2] = '+'
					}
				}
			}
		}
		if e.GetName() != "" {
			ca[1][2] = rune(e.nodeInfo.nodeName[0])
			ca[2][2] = rune(e.nodeInfo.nodeName[1])
			ca[3][2] = rune(e.nodeInfo.nodeName[2])
		}
		if e.pathInfo != nil {
			ca[2][1] = rune(strconv.Itoa(e.pathInfo.pathNumber)[0])
		}
		if len(e.nodeInfo.nodeStatus) >= 3 {
			ca[1][3] = rune(e.nodeInfo.nodeStatus[0])
			ca[2][3] = rune(e.nodeInfo.nodeStatus[1])
			ca[3][3] = rune(e.nodeInfo.nodeStatus[2])
		}
		// draw path cell
	} else if e.pathInfo != nil {
		ca[2][2] = rune(strconv.Itoa(e.pathInfo.pathNumber)[0])
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				if e.GetConnectionByCoords(x, y) != nil {
					ca[2+x*2][2+y*2] = ' '
					ca[2+x][2+y] = ' '
				}
			}
		}
	}
	return ca
}

func (r *LayoutMap) WholeMapToCharArray() *[][]rune {
	sx, sy := r.GetSize()
	ca := make([][]rune, 5*sx)
	for i := range ca {
		ca[i] = make([]rune, 5*sy)
	}
	for x := 0; x < len(r.elements); x++ {
		for y := 0; y < len(r.elements[0]); y++ {
			cellArr := r.CellToCharArray(x, y)
			for i := 0; i < 5; i++ {
				for j := 0; j < 5; j++ {
					ca[5*x+i][5*y+j] = cellArr[i][j]
				}
			}
		}
	}
	return &ca
}
