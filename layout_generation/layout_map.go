package layout_generation

import (
	rnd "github.com/sidav/golibrl/random"
	"strconv"
)

type LayoutMap struct {
	// room map with connections or smth
	elements [][] *element
}

func (r *LayoutMap) init(sizex, sizey int) {
	r.elements = make([][]*element, size)
	for i := range r.elements {
		r.elements[i] = make([]*element, size)
	}
	for x := 0; x < sizex; x++ {
		for y := 0; y < sizey; y++ {
			r.elements[x][y] = &element{connections: map[string]*connection {"north": nil, "east": nil, "south": nil, "west": nil}}
		}
	}
}

func (r *LayoutMap) placeNodeAtCoords(x, y int, nodeName string) {
	r.elements[x][y].nodeInfo = &node_cell{nodeName: nodeName}
}

func (r *LayoutMap) placePathAtCoords(x, y int, pathNum int) {
	r.elements[x][y].pathInfo = &path_cell{pathNum}
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

func (r *LayoutMap) getRandomPathCell(desiredPathNum int, nodesAllowed bool) (int, int) { // desiredPathNum -1 means any path
	const tries = 50
	for try := 0; try < tries; try++ {
		x, y := rnd.Random(size), rnd.Random(size)
		if !r.elements[x][y].isPartOfAPath() {
			continue
		}
		if !nodesAllowed && r.elements[x][y].nodeInfo != nil { // don't take nodes unless allowed
			continue
		}
		if desiredPathNum > -1 && desiredPathNum != r.elements[x][y].pathInfo.pathNumber { // don't take cells of non-desired path numbers
			continue
		}
		return x, y
	}
	return -1, -1
}

func (r *LayoutMap) getRandomPathCoordsAndRandomCellNearPath(pathNum int, allowNearNode bool) (int, int, int, int) {
	const tries = 10
	for try := 0; try < tries; try++ {
		px, py := r.getRandomPathCell(pathNum, allowNearNode)
		if px == -1 && py == -1 {
			continue
		}
		for try2 := 0; try2 < tries; try2++ {
			x, y := rnd.RandInRange(px-1, px+1), rnd.RandInRange(py-1, py+1)
			if x >= 0 && y >= 0 && x < len(r.elements) && y < len(r.elements[0]) && r.elements[x][y].isEmpty() {
				return px, py, x, y
			}
		}
	}
	return -1, -1, -1, -1
}

func (r *LayoutMap) getRandomEmptyCellCoords() (int, int) { // desiredPathNum -1 means any path
	const tries = 25
	for i := 0; i < tries; i++ {
		x, y := rnd.Random(size), rnd.Random(size)
		if r.elements[x][y].isEmpty() {
			return x, y
		}
	}
	return -1, -1
}

func (r *LayoutMap) areCoordsEmpty(x, y int) bool {
	return r.elements[x][y].isEmpty()
}

func (r *LayoutMap) getCoordsOfNode(nodeName string) (int, int) {
	for x := 0; x < len(r.elements); x++ {
		for y := 0; y < len(r.elements[0]); y++ {
			if r.elements[x][y].isNode() && r.elements[x][y].nodeInfo.nodeName == nodeName {
				return x, y
			}
		}
	}
	panic("getCoordsOfNode failed with node "+nodeName)
	return -1, -1
}

// exported

func (r *LayoutMap) GetSize() (int, int) {
	return len(r.elements), len(r.elements[0])
}

func (r *LayoutMap) getPassabilityMapForPathfinder() *[][]int {
	pmap := make([][]int, size)
	for i := range pmap {
		pmap[i] = make([]int, size)
	}

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if layout.areCoordsEmpty(x, y) {
				pmap[x][y] = 1
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
	if elem.isNode() {
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
	for i := range (ca) {
		ca[i] = make([]rune, 5)
	}

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			ca[x][y] = '#'
		}
	}
	if e.nodeInfo != nil {
		for x := 1; x < 4; x++ {
			for y := 1; y < 4; y++ {
				ca[x][y] = ' '
			}
		}
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 2; y++ {
				if e.getConnectionByCoords(x, y) != nil {
					ca[2+x*2][2+y*2] = '+'
				}
			}
		}
		ca[1][2] = rune(e.nodeInfo.nodeName[0])
		ca[2][2] = rune(e.nodeInfo.nodeName[1])
		ca[3][2] = rune(e.nodeInfo.nodeName[2])
		if e.pathInfo != nil {
			ca[2][1] = rune(strconv.Itoa(e.pathInfo.pathNumber)[0])
		}
	} else if e.pathInfo != nil {
		ca[2][2] = rune(strconv.Itoa(e.pathInfo.pathNumber)[0])
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				if e.getConnectionByCoords(x, y) != nil {
					ca[2+x*2][2+y*2] = ' '
					ca[2+x][2+y] = ' '
				}
			}
		}
	}
	return ca
}