package godot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createPlayerNodeTree() Node {
	player := Node{
		Name:     "Player",
		Type:     "Spatial",
		Children: make(map[string]*Node),
	}

	arm := Node{
		Name:     "Arm",
		Type:     "Spatial",
		Children: make(map[string]*Node),
	}

	hand := Node{
		Name:     "Hand",
		Type:     "Spatial",
		Children: make(map[string]*Node),
	}

	thumb := Node{
		Name: "Thumb",
		Type: "Spatial",
	}

	indexFinger := Node{
		Name: "Index Finger",
		Type: "Spatial",
	}

	hand.AddNode(&thumb)
	hand.AddNode(&indexFinger)
	arm.AddNode(&hand)
	player.AddNode(&arm)

	return player
}

func TestGetNodeWithSelf(t *testing.T) {
	player := createPlayerNodeTree()
	node, err := player.GetNode(".")
	assert.NoError(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, "Player", node.Name)
}

func TestGetNodeWithDirectChild(t *testing.T) {
	player := createPlayerNodeTree()
	node, err := player.GetNode("Arm")
	assert.NoError(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, "Arm", node.Name)
	assert.Equal(t, "Player", node.Parent.Name)
}

func TestGetNodeWithDeepPath(t *testing.T) {
	player := createPlayerNodeTree()
	node, err := player.GetNode("Arm/Hand")
	assert.NoError(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, "Hand", node.Name)
	assert.Equal(t, "Arm", node.Parent.Name)
}

func TestGetNodeOnChildNode(t *testing.T) {
	player := createPlayerNodeTree()
	node, err := player.GetNode("Arm")
	assert.NoError(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, "Arm", node.Name)
	hand, err := node.GetNode("Hand")
	assert.NoError(t, err)
	assert.Equal(t, "Hand", hand.Name)
}

func TestGetNodeWithInvalidPath(t *testing.T) {
	player := createPlayerNodeTree()
	_, err := player.GetNode("Arm/Leg")
	assert.Error(t, err)
}

func TestAddNode(t *testing.T) {
	player := createPlayerNodeTree()
	hand, err := player.GetNode("Arm/Hand")
	assert.NoError(t, err)
	assert.Len(t, hand.Children, 2)

	hand.AddNode(&Node{
		Name: "Middle Finger",
		Type: "Spatial",
	})
	assert.Len(t, hand.Children, 3)

	middleFinger, err := hand.GetNode("Middle Finger")
	assert.NoError(t, err)
	assert.Equal(t, "Middle Finger", middleFinger.Name)
}

func TestRemoveNodeWithDeepPath(t *testing.T) {
	player := createPlayerNodeTree()
	hand, err := player.GetNode("Arm/Hand")
	assert.NoError(t, err)
	assert.Len(t, hand.Children, 2)
	err = player.RemoveNode("Arm/Hand/Thumb")
	assert.NoError(t, err)
	assert.Len(t, hand.Children, 1)
}

func TestRemoveNodeWithDirectChild(t *testing.T) {
	player := createPlayerNodeTree()
	hand, err := player.GetNode("Arm/Hand")
	assert.NoError(t, err)
	assert.Len(t, hand.Children, 2)
	err = hand.RemoveNode("Thumb")
	assert.NoError(t, err)
	assert.Len(t, hand.Children, 1)
}

func TestRemoveNodeWithChildren(t *testing.T) {
	player := createPlayerNodeTree()
	err := player.RemoveNode("Arm/Hand")
	assert.NoError(t, err)
	arm, err := player.GetNode("Arm")
	assert.NoError(t, err)
	assert.Len(t, arm.Children, 0)
}

func TestRemoveNodeWithSelfReturnsError(t *testing.T) {
	player := createPlayerNodeTree()
	err := player.RemoveNode(".")
	assert.Error(t, err)
}

func TestRemoveInvalidNode(t *testing.T) {
	player := createPlayerNodeTree()
	err := player.RemoveNode("Arm/Leg")
	assert.Error(t, err)
}
