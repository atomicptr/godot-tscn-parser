package godot

import (
	"fmt"
	"strings"
)

// Node is Godots central building block. Godot stores scenes as a tree of nodes.
type Node struct {
	Name     string
	Type     string
	Instance Type
	Fields   map[string]interface{}
	Children map[string]*Node
	Parent   *Node
	MetaData
}

// AddNode adds a node as the child of the current node
func (n *Node) AddNode(node *Node) {
	node.Parent = n
	n.Children[node.Name] = node
}

// GetNode retrieves a node for a given path
func (n *Node) GetNode(path string) (*Node, error) {
	if path == "." {
		return n, nil
	}

	parts := strings.Split(path, "/")
	root := n
	for _, p := range parts {
		node, hasChild := root.Children[p]

		if !hasChild {
			return nil, fmt.Errorf("could not get node path: %s", path)
		}

		root = node
	}

	return root, nil
}

// RemoveNode deletes a child node
func (n *Node) RemoveNode(path string) error {
	node, err := n.GetNode(path)
	if err != nil {
		return err
	}

	parent := node.Parent

	if parent == nil {
		return fmt.Errorf("can't remove root node, or this node is not attached to anything")
	}

	delete(parent.Children, node.Name)
	return nil
}
