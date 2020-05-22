// Copyright 2018 Ara Israelyan. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


package nestedset

// NodeInterface is the interface implemented by types that can be used by nodes in nested set
type NodeInterface interface {
	GetId() int64    // Returns id of node
	GetName() string // Returns name of node

	GetLevel() int64 // Returns level of node
	GetLeft() int64  // Returns left of node
	GetRight() int64 // Returns right of node

	SetLevel(int64)   // Sets node level
	SetLeft(int64)    // Sets node left
	SetRight(int64)   // Sets node right
}

// Node represents generic node type with NodeInterface implementation
type Node struct {
	NodeId    int64    `json:"id"`
	NodeName  string   `json:"node_name"`
	NodeLevel int64    `json:"level"`
	NodeLeft  int64    `json:"left"`
	NodeRight int64    `json:"right"`
}

func (n Node) GetId() int64 {
	return n.NodeId
}

func (n Node) GetName() string {
	return n.NodeName
}

func (n Node) GetLevel() int64 {

	return n.NodeLevel
}

func (n Node) GetLeft() int64 {
	return n.NodeLeft
}

func (n Node) GetRight() int64 {
	return n.NodeRight
}

func (n *Node) SetLevel(level int64) {
	n.NodeLevel = level
}

func (n *Node) SetLeft(left int64) {
	n.NodeLeft = left
}

func (n *Node) SetRight(right int64) {
	n.NodeRight = right
}
