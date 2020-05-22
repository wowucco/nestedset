
package nestedset

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
)

// SortedNodes represent nodes array sorted by left value.
type SortedNodes []NodeInterface

func (sn SortedNodes) Len() int           { return len(sn) }
func (sn SortedNodes) Swap(i, j int)      { sn[i], sn[j] = sn[j], sn[i] }
func (sn SortedNodes) Less(i, j int) bool { return sn[i].Left() < sn[j].Left() }

// NestedSet represents a nested set management type.
type NestedSet struct {
	nodes    []NodeInterface
	rootNode NodeInterface
	maxId    int64
	mutex    sync.Mutex
}

// NewNestedSet creates and returns a new instance of NestedSet with root node.
func NewNestedSet(rootNode NodeInterface) *NestedSet {

	rootNode.SetRight(1)

	s := NestedSet{
		nodes:    make([]NodeInterface, 0),
		rootNode: rootNode,
	}

	s.nodes = append(s.nodes, s.rootNode)

	return &s
}

// Overrides json.Marshaller.MarshalJSON().
func (s NestedSet) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(s.nodes, "", "  ")
}

// Adds new node to nested set. If `parent` nil, add node to root node.
func (s *NestedSet) Add(newNode, parent NodeInterface) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if parent != nil {
		if !s.exists(parent) {
			return errors.New("Parent node not found in structure")
		}

	} else {
		parent = s.rootNode
	}

	level := parent.Level() + 1
	right := parent.Right()

	newNode.SetLevel(level)
	s.maxId++

	newNode.SetLeft(right)
	newNode.SetRight(right + 1)

	for _, n := range s.nodes {

		if n.Right() >= right {
			n.SetRight(n.Right() + 2)
			if n.Left() > right {
				n.SetLeft(n.Left() + 2)
			}
		}
	}

	s.nodes = append(s.nodes, newNode)

	return nil
}

// Deletes node from nested set.
func (s *NestedSet) Delete(node NodeInterface) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if node == nil || node == s.rootNode {
		return errors.New("Can't delete root node")
	}

	if !s.exists(node) {
		return errors.New("Node not found in structure")
	}

	newNodes := make([]NodeInterface, 0)

	for _, n := range s.nodes {

		if n.Left() < node.Left() || n.Right() > node.Right() {

			if n.Right() > node.Right() {
				n.SetRight(n.Right() - (node.Right() - node.Left() + 1))
			}

			if n.Left() > node.Left() {
				n.SetLeft(n.Left() - (node.Right() - node.Left() + 1))
			}

			newNodes = append(newNodes, n)

		}
	}

	s.nodes = newNodes

	return nil
}

// Moves node and its branch to another parent node.
func (s *NestedSet) Move(node, parent NodeInterface) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if node.Level() == 0 {
		return errors.New("Can't move root node")
	}

	if parent == nil {
		parent = s.rootNode
	}

	if parent.Left() >= node.Left() && parent.Right() <= node.Right() {
		return errors.New("Can't move branch to node within itself")
	}

	currentParent := s.parent(node)
	if currentParent == nil {
		return errors.New("Parent node not found, the structure broken")
	}
	if currentParent == parent {
		return errors.New("Moving in same branch not implemented")
	}

	level := node.Level()
	left := node.Left()
	right := node.Right()
	levelUp := parent.Level()
	rightNear := parent.Right() - 1
	skewLevel := levelUp - level + 1
	skewTree := right - left + 1
	skewEdit := rightNear - left + 1
	isUp := rightNear < right

	toUpdate := s.branch(node)

	if isUp {
		for _, n := range s.nodes {

			if n.Right() < left && n.Right() > rightNear {
				n.SetRight(n.Right() + skewTree)
			}
			if n.Left() < left && n.Left() > rightNear {
				n.SetLeft(n.Left() + skewTree)
			}
		}
	} else {
		skewEdit = rightNear - left + 1 - skewTree

		for _, n := range s.nodes {

			if n.Right() > right && n.Right() <= rightNear {
				n.SetRight(n.Right() - skewTree)
			}

			if n.Left() > right && n.Left() <= rightNear {
				n.SetLeft(n.Left() - skewTree)
			}
		}
	}

	for _, n := range toUpdate {
		n.SetLeft(n.Left() + skewEdit)
		n.SetRight(n.Right() + skewEdit)
		n.SetLevel(n.Level() + skewLevel)
	}

	return nil
}

// Returns parent for node.
func (s *NestedSet) Parent(node NodeInterface) NodeInterface {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.parent(node)
}

func (s *NestedSet) parent(node NodeInterface) NodeInterface {

	for _, n := range s.nodes {
		if n.Left() <= node.Left() && n.Right() >= node.Right() && n.Level() == (node.Level()-1) {
			return n
		}
	}

	return nil
}

// Finds and returns node by id.
func (s *NestedSet) FindById(id int64) NodeInterface {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, n := range s.nodes {
		if n.Id() == id {
			return n
		}
	}

	return nil
}

// Returns branch for node, including itself.
func (s *NestedSet) Branch(node NodeInterface) []NodeInterface {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.branch(node)
}

func (s *NestedSet) branch(node NodeInterface) []NodeInterface {

	sort.Sort(SortedNodes(s.nodes))

	var res []NodeInterface

	// Return full tree
	if node == nil {
		res = make([]NodeInterface, len(s.nodes))
		copy(res, s.nodes)
		return res
	}

	if !s.exists(node) {
		return nil
	}

	res = make([]NodeInterface, 0)

	for _, n := range s.nodes {
		if n.Left() >= node.Left() && n.Right() <= node.Right() {
			res = append(res, n)
		}
	}

	return res
}

func (s *NestedSet) exists(node NodeInterface) bool {

	bFound := false
	for _, n := range s.nodes {
		if n == node {
			bFound = true
			break
		}
	}

	return bFound
}
