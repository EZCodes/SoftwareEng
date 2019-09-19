package LCA

struct Tree{
	root Node
}
//TODO make compatible with multiple childs
struct Node{
	data int
	childs []*Node
}

(t *Tree) lowestCommonAncestor(n1, n2 *Node) Node {
	var p1 []Node
	var p2 []Node
	p1 = n1.findPath(t.root)
	p2 = n2.findPath(t.root) 
	var lcm Node
	for i:=0; i<len(p1) && i<len(p2); i++ {
		if p1[i].data == p2[i].data{
			lcm = p1[i]
		}
	}
	return lcm
}
(n *Node) findPath(other *Node, path []Node) []Node{
	if n.data == other.data {
		path = append(path, &other)
		return path
	}
	path = append(path, &other)
	for _, node :=  range other.childs{
		n.findPath(&node, path)
		if len(path)>0 && path[len(path)-1].data == n.data{
			return path
		}
		path = path[:len(path)-1]
	}
	return path
}