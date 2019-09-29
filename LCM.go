package LCA

import "fmt"

type Tree struct {
	root Node
}
type Node struct {
	data   int
	childs []*Node
}

// LowestCommonAncestor find the lowest common ancestor of two nodes
// if found, the node with it is returned, if two nodes dont have a common ancestor
// empty node is returned
func (t *Tree) LowestCommonAncestor(n1, n2 *Node) (Node, error) {
	var p1 []Node
	var p2 []Node
	if n1 == nil || n2 == nil {
		return Node{}, fmt.Errorf("One of the nodes is nil. n1: %v, n2: %v", n1, n2)
	}
	p1 = n1.findPath(&(t.root), []Node{})
	p1 = p1[:len(p1)-1] // last elements is itself, not needed
	p2 = n2.findPath(&(t.root), []Node{})
	p2 = p2[:len(p2)-1] // last elements is itself, not needed
	var lcm Node
	for i := 0; i < len(p1) && i < len(p2); i++ {
		if p1[i].data == p2[i].data {
			lcm = p1[i]
		}
	}
	return lcm, nil
}
func (n *Node) findPath(other *Node, path []Node) []Node {
	if n.data == other.data {
		path = append(path, *other)
		return path
	}
	path = append(path, *other)
	for _, node := range other.childs {
		path = n.findPath(node, path)
		if len(path) > 0 && path[len(path)-1].data == n.data {
			return path
		}
		if len(path) > 0 {
			path = path[:len(path)-1]
		}
	}
	return path
}
