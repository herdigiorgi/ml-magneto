package mutant

import (
	"errors"
	"unicode"
)

const maxDirections int = 8
const maxDepth int = 4

// Node used to represent each DNA letter
type Node struct {
	value     rune
	nodes     []*Node
	discarted []bool
}

func (n *Node) initValue(value rune) {
	n.value = value
	n.nodes = make([]*Node, maxDirections)
	n.discarted = make([]bool, maxDirections)
}

func (n *Node) getSubNode(direction int) *Node {
	return n.nodes[direction]
}

func (n *Node) setSubNode(direction int, subNode *Node) {
	n.nodes[direction] = subNode
}

func (n *Node) isDiscarted(direction int) bool {
	return n.discarted[direction]
}

func (n *Node) isEmpty() bool {
	return len(n.nodes) <= 0
}

func (n *Node) discart(direction int) {
	if !n.isEmpty() {
		n.discarted[direction] = true
	}
}

func (n *Node) inverseDirection(direction int) int {
	return (direction + maxDirections/2) % maxDirections
}

func (n *Node) discartInverse(direction int) {
	n.discart(n.inverseDirection(direction))
}

func (n *Node) checkDirected(value rune, direction int, depth int) bool {
	if n.value == value && depth >= maxDepth {
		return true
	} else if n.value == value {
		neighbour := n.getSubNode(direction)
		if neighbour != nil &&
			neighbour.checkDirected(value, direction, depth+1) {
			return true
		}
	}
	n.discartInverse(direction)
	return false
}

func (n *Node) check() bool {
	for dir := 0; dir < maxDirections; dir++ {
		if n.isDiscarted(dir) {
			continue
		}
		neighbour := n.getSubNode(dir)
		if neighbour != nil &&
			neighbour.checkDirected(n.value, dir, 2) {
			return true
		}
		n.discart(dir)
	}
	return false
}

func checkNodes(nxnNodes [][]Node) bool {
	for _, nodes := range nxnNodes {
		for _, node := range nodes {
			if node.check() {
				return true
			}
		}
	}
	return false
}

func isValidAdnLetter(x rune) bool {
	return x == 'A' || x == 'T' || x == 'C' || x == 'G'
}

func sanitizeMutantInput(input []string) ([][]rune, error) {
	var result [][]rune
	if len(input) < 0 {
		return nil, errors.New("empty input")
	}
	length := len(input[0])
	for _, e := range input {
		if len(e) != length {
			return nil, errors.New("invalid length")
		}
		var chrs []rune
		for _, chr := range e {
			chr = unicode.ToUpper(chr)
			if !isValidAdnLetter(chr) {
				return nil, errors.New("invalid adn sequence")
			}
			chrs = append(chrs, chr)
		}
		result = append(result, chrs)
	}
	return result, nil
}

// Transforms an NxM array of runes into is Node array equivalent. The goal is
// to make the input array navigable. The "direction" indexes have this shape:
// 0  1  2
// 7  X  3
// 6  5  4
func mutantInputToNodes(input [][]rune) [][]Node {
	result := make([][]Node, len(input))
	for i := 0; i < len(input); i++ {
		result[i] = make([]Node, len(input[0]))
	}

	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input[i]); j++ {
			node := &result[i][j]
			node.initValue(input[i][j])
			if i > 0 {
				if j > 0 {
					result[i-1][j-1].setSubNode(4, node)
					node.setSubNode(0, &result[i-1][j-1])
				}
				if j < len(input[i])-1 {
					result[i-1][j+1].setSubNode(6, node)
					node.setSubNode(2, &result[i-1][j+1])
				}
				result[i-1][j].setSubNode(5, node)
				node.setSubNode(1, &result[i-1][j])
			}
			if j > 0 {
				result[i][j-1].setSubNode(3, node)
				node.setSubNode(7, &result[i][j-1])
			}
		}
	}
	return result
}

// IsMutant magneto's detection algorithm
func IsMutant(input []string) (bool, error) {
	sanitized, err := sanitizeMutantInput(input)
	if err != nil {
		return false, err
	}
	nodes := mutantInputToNodes(sanitized)
	return checkNodes(nodes), nil
}
