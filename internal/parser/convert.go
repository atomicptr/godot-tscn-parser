package parser

import (
	"fmt"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
	"github.com/pkg/errors"
)

func (tscn *TscnFile) ConvertToGodotScene() (*godot.Scene, error) {
	if tscn.Key != "gd_scene" {
		return nil, fmt.Errorf("can't convert %s to gd_scene", tscn.Key)
	}

	scene := godot.Scene{
		ExtResources: make(map[int64]*godot.ExtResource),
		SubResources: make(map[int64]*godot.SubResource),
		MetaData: godot.MetaData{
			LexerPosition: tscn.Pos,
		},
	}

	// handle everything that isn't a node
	for _, section := range tscn.Sections {
		// External resources
		if section.ResourceType == "ext_resource" {
			res, err := convertSectionToExtResource(section)
			if err != nil {
				return nil, err
			}
			scene.ExtResources[res.Id] = res
			continue
		}

		// Internal resources
		if section.ResourceType == "sub_resource" {
			res, err := convertSectionToSubResource(section)
			if err != nil {
				return nil, err
			}
			scene.SubResources[res.Id] = res
			continue
		}

		// TODO: Add support for "editable" (attributes path)
		// TODO: Add support for "connection" (attributes from, to, method, signal)
	}

	rootNode, err := buildNodeTree(tscn)
	if err != nil {
		return nil, err
	}

	scene.Node = rootNode

	return &scene, nil
}

func convertSectionToExtResource(section *GdResource) (*godot.ExtResource, error) {
	if section.ResourceType != "ext_resource" {
		return nil, fmt.Errorf("you can't convert a %s to ext_resource", section.ResourceType)
	}

	path, err := section.GetAttribute("path")
	if err != nil {
		return nil, errors.Wrap(
			err,
			"could not convert, because ext_resource doesn't have required field path",
		)
	}

	resType, err := section.GetAttribute("type")
	if err != nil {
		return nil, errors.Wrap(
			err,
			"could not convert, because ext_resource doesn't have required field type",
		)
	}

	id, err := section.GetAttribute("id")
	if err != nil {
		return nil, errors.Wrap(
			err,
			"could not convert, because ext_resource doesn't have required field id",
		)
	}

	return &godot.ExtResource{
		Path: *path.String,
		Type: *resType.String,
		Id:   *id.Integer,
		MetaData: godot.MetaData{
			LexerPosition: section.Pos,
		},
	}, nil
}

func convertSectionToSubResource(section *GdResource) (*godot.SubResource, error) {
	if section.ResourceType != "sub_resource" {
		return nil, fmt.Errorf("you can't convert a %s to sub_resource", section.ResourceType)
	}

	resType, err := section.GetAttribute("type")
	if err != nil {
		return nil, errors.Wrap(
			err,
			"could not convert, because sub_resource doesn't have required field type",
		)
	}

	id, err := section.GetAttribute("id")
	if err != nil {
		return nil, errors.Wrap(
			err,
			"could not convert, because sub_resource doesn't have required field id",
		)
	}

	subResource := godot.SubResource{
		Type: *resType.String,
		Id:   *id.Integer,
		MetaData: godot.MetaData{
			LexerPosition: section.Pos,
		},
	}

	for _, field := range section.Fields {
		// TODO: properly parse structures like bones/0/name = "Bone", bones/0/parent = -1, etc.
		subResource.Fields[field.Key] = convertGdValue(field.Value)
	}

	return &subResource, nil
}

func buildNodeTree(tscn *TscnFile) (*godot.Node, error) {
	root, otherNodes := findNodes(tscn)

	rootNode, err := convertSectionToUnattachedNode(root)
	if err != nil {
		return nil, errors.Wrap(err, "could not determine root node")
	}

	// a list of indices of nodes that have been processed
	var processedNodes []int

	// a counter to check how often we couldn't find the parent node
	couldntFindParentNodeCounter := 0

	// if that counter exceeds the threshold, we'll throw an error to stop execution
	const couldntFindParentNodeCounterThreshold = 1000000

	// while not all nodes have been processed
	for len(otherNodes) != len(processedNodes) {
		for index, sectionNode := range otherNodes {
			if intArrayContains(index, processedNodes) {
				continue
			}

			// since we've removed the root sectionNode there shouldn't be another one without parent
			parentAttribute, _ := sectionNode.GetAttribute("parent")
			parentNodePath, ok := parentAttribute.Raw().(string)
			if !ok {
				return nil, fmt.Errorf("section attribute parent is not a string: %v", parentAttribute.Raw())
			}

			parentNode, err := rootNode.GetNode(parentNodePath)
			if err != nil {
				couldntFindParentNodeCounter += 1
				if couldntFindParentNodeCounter >= couldntFindParentNodeCounterThreshold {
					return nil, errors.New("can't build node tree, either its invalid or way too big (over a million nodes)")
				}

				// couldn't find parent node, continue for now...
				continue
			}

			node, err := convertSectionToUnattachedNode(sectionNode)
			if err != nil {
				return nil, errors.Wrap(err, "could not parse node")
			}
			parentNode.AddNode(node)
			processedNodes = append(processedNodes, index)

			couldntFindParentNodeCounter = 0
		}
	}

	return rootNode, nil
}

func findNodes(tscn *TscnFile) (*GdResource, []*GdResource) {
	var rootNode *GdResource
	var nodes []*GdResource

	for _, section := range tscn.Sections {
		if section.ResourceType != "node" {
			continue
		}

		parent, _ := section.GetAttribute("parent")
		// node without parent field is the root node
		if parent == nil {
			rootNode = section
			continue
		}

		nodes = append(nodes, section)
	}

	return rootNode, nodes
}

func convertSectionToUnattachedNode(section *GdResource) (*godot.Node, error) {
	if section.ResourceType != "node" {
		return nil, fmt.Errorf("you can't convert a %s to node", section.ResourceType)
	}

	name, err := section.GetAttribute("name")
	if err != nil {
		return nil, errors.Wrap(err, "node without name")
	}

	node := godot.Node{
		Name:     *name.String,
		Fields:   make(map[string]interface{}),
		Children: make(map[string]*godot.Node),
		MetaData: godot.MetaData{
			LexerPosition: section.Pos,
		},
	}

	err = attachTypeToNode(&node, section)
	if err != nil {
		return nil, err
	}

	for _, field := range section.Fields {
		node.Fields[field.Key] = convertGdValue(field.Value)
	}

	return &node, nil
}

func attachTypeToNode(node *godot.Node, section *GdResource) error {
	if nodeType, err := section.GetAttribute("type"); err == nil {
		node.Type = *nodeType.String
		return nil
	}

	if instance, err := section.GetAttribute("instance"); err == nil {
		if len(instance.Type.Parameters) != 1 {
			return fmt.Errorf("node instance parameter does not contain a valid reference %v", instance.Raw())
		}

		node.Instance = fmt.Sprintf(
			"%s:%d",
			instance.Type.Key,
			*instance.Type.Parameters[0].Integer,
		)
		return nil
	}

	// nodes don't have to have a type, if it doesn't have one just don't do anything
	return nil
}

func convertGdValue(val *GdValue) godot.Value {
	// TODO: handle special values like maps, array and parse them properly if necessary (Vectors, Transforms etc)
	return val
}
