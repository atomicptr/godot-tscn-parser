package convert

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

const (
	resourceTypeExtResource = "ext_resource"
	resourceTypeSubResource = "sub_resource"
	resourceTypeEditable    = "editable"
	resourceTypeConnection  = "connection"
	resourceTypeNode        = "node"
)

const (
	internalNodeUnassignableNodes = "__Internal_UnassignableNodes"
	internalNodeType              = "__Internal"
	internalNodeParentPathField   = "__Internal_ParentPath"
)

const (
	volatileNodeType = "VolatileNode"
)

// ToGodotScene tries to convert a TscnFile structure to an actual Godot Scene with a node tree
func ToGodotScene(tscn *parser.TscnFile) (*godot.Scene, error) {
	if tscn.Key != TscnTypeGodotScene {
		return nil, fmt.Errorf("can't convert %s to gd_scene", tscn.Key)
	}

	scene := &godot.Scene{
		ExtResources: make(map[int64]*godot.ExtResource),
		SubResources: make(map[int64]*godot.SubResource),
		MetaData: godot.MetaData{
			LexerPosition: tscn.Pos,
		},
	}

	// handle everything that isn't a node
	for _, section := range tscn.Sections {
		// External resources
		if section.ResourceType == resourceTypeExtResource {
			res, err := convertSectionToExtResource(section)
			if err != nil {
				return nil, err
			}
			scene.ExtResources[res.ID] = res
			continue
		}

		// Internal resources
		if section.ResourceType == resourceTypeSubResource {
			res, err := convertSectionToSubResource(section)
			if err != nil {
				return nil, err
			}
			scene.SubResources[res.ID] = res
			continue
		}

		// editable scenes
		if section.ResourceType == resourceTypeEditable {
			editable, err := convertSectionToEditable(section)
			if err != nil {
				return nil, err
			}
			scene.Editables = append(scene.Editables, editable)
			continue
		}

		// connections
		if section.ResourceType == resourceTypeConnection {
			connection, err := convertSectionToConnection(section)
			if err != nil {
				return nil, err
			}
			scene.Connections = append(scene.Connections, connection)
			continue
		}

		// ignore nodes for now...
		if section.ResourceType == resourceTypeNode {
			continue
		}

		// something else found? Whoops, throw error
		return nil, fmt.Errorf("invalid resource type found: %s [%s]", section.ResourceType, section.Pos)
	}

	rootNode, err := buildNodeTree(tscn)
	if err != nil {
		return nil, err
	}

	scene.Node = rootNode

	err = postProcessAndCleanSceneFromInternals(scene)
	if err != nil {
		return nil, err
	}

	// TODO: add validator and validate scene
	return scene, nil
}

func postProcessAndCleanSceneFromInternals(scene *godot.Scene) error {
	err := createVolatileNodes(scene)
	if err != nil {
		return err
	}

	unassignableNodes, err := scene.GetNode(internalNodeUnassignableNodes)
	if err != nil {
		return err
	}

	if len(unassignableNodes.Children) > 0 {
		return fmt.Errorf("node tree contains an invalid tree")
	}

	err = scene.RemoveNode(internalNodeUnassignableNodes)
	if err != nil {
		return err
	}

	return nil
}

func createVolatileNodes(scene *godot.Scene) error {
	unassignableNodes, err := scene.GetNode(internalNodeUnassignableNodes)
	if err != nil {
		return err
	}

	var nodesToBeDeleted []string

	// convert map into slice
	children := convertNodeMapIntoSortedSlice(unassignableNodes.Children)

	// create the volatile node structure
	for _, node := range children {
		nodeIdent := node.Name
		parentPath, ok := node.Fields[internalNodeParentPathField]
		if !ok {
			continue
		}

		p := parentPath.(string)
		for _, editable := range scene.Editables {
			// doesn't match the editable, ignore
			if !strings.HasPrefix(p, editable.Path) {
				continue
			}

			// for every path segment check if a child node exists with the given name
			parts := strings.Split(p, "/")
			parentNode := scene.Node
			for _, pathPart := range parts {
				// parent node has the path part? Dig deeper
				if n, ok := parentNode.Children[pathPart]; ok {
					parentNode = n
					continue
				}

				// if not, create a volatile node here
				volatileNode := &godot.Node{
					Name:     pathPart,
					Type:     volatileNodeType,
					Children: make(map[string]*godot.Node),
				}
				parentNode.AddNode(volatileNode)
				parentNode = volatileNode
			}
		}

		// try to add the un-assignable node into the tree
		parentNode, err := scene.GetNode(p)
		if err != nil {
			return err
		}
		parentNode.AddNode(node)
		delete(node.Fields, internalNodeParentPathField)
		nodesToBeDeleted = append(nodesToBeDeleted, nodeIdent)
	}

	for _, nodeIdent := range nodesToBeDeleted {
		delete(unassignableNodes.Children, nodeIdent)
	}

	return nil
}

func convertNodeMapIntoSortedSlice(nodeMap map[string]*godot.Node) []*godot.Node {
	var children []*godot.Node
	for _, node := range nodeMap {
		children = append(children, node)
	}

	parentPathLen := func(n *godot.Node) int {
		path := n.Fields[internalNodeParentPathField]
		p := path.(string)
		parts := strings.Split(p, "/")
		return len(parts)
	}

	sort.SliceStable(children, func(i, j int) bool {
		return parentPathLen(children[i]) < parentPathLen(children[j])
	})

	return children
}

func convertSectionToExtResource(section *parser.GdResource) (*godot.ExtResource, error) {
	if section.ResourceType != resourceTypeExtResource {
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
		ID:   *id.Integer,
		MetaData: godot.MetaData{
			LexerPosition: section.Pos,
		},
	}, nil
}

func convertSectionToSubResource(section *parser.GdResource) (*godot.SubResource, error) {
	if section.ResourceType != resourceTypeSubResource {
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
		Type:   *resType.String,
		ID:     *id.Integer,
		Fields: make(map[string]interface{}),
		MetaData: godot.MetaData{
			LexerPosition: section.Pos,
		},
	}

	insertFieldEntriesFromSection(section, subResource.Fields)

	return &subResource, nil
}

func buildNodeTree(tscn *parser.TscnFile) (*godot.Node, error) {
	root, otherNodes := findNodes(tscn)

	rootNode, err := convertSectionToUnattachedNode(root)
	if err != nil {
		return nil, errors.Wrap(err, "could not determine root node")
	}

	// add a node to allow us to gather un-assignable nodes, they might have volatile ancestors
	unassignableNodeList := &godot.Node{
		Name:     internalNodeUnassignableNodes,
		Type:     internalNodeType,
		Fields:   make(map[string]interface{}),
		Children: make(map[string]*godot.Node),
	}
	rootNode.AddNode(unassignableNodeList)

	// a list of indices of nodes that have been processed
	var processedNodes []int

	// a counter to check how often we couldn't find the parent node
	couldntFindParentNodeCounter := 0

	// if that counter exceeds the threshold, we'll add it to the un-assignable nodes list
	const couldntFindParentNodeCounterThreshold = 1000000

	parseAndAddNode := func(parent *godot.Node, section *parser.GdResource, index int) (*godot.Node, error) {
		node, err := convertSectionToUnattachedNode(section)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse node")
		}
		parent.AddNode(node)
		processedNodes = append(processedNodes, index)
		return node, nil
	}

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
				couldntFindParentNodeCounter++
				if couldntFindParentNodeCounter >= couldntFindParentNodeCounterThreshold {
					node, err := parseAndAddNode(unassignableNodeList, sectionNode, index)
					if err != nil {
						return nil, err
					}
					node.Fields[internalNodeParentPathField] = parentNodePath
					couldntFindParentNodeCounter = 0
					continue
				}

				// couldn't find parent node, continue for now...
				continue
			}

			_, err = parseAndAddNode(parentNode, sectionNode, index)
			if err != nil {
				return nil, err
			}

			couldntFindParentNodeCounter = 0
		}
	}

	return rootNode, nil
}

func findNodes(tscn *parser.TscnFile) (rootNode *parser.GdResource, nodes []*parser.GdResource) {
	for _, section := range tscn.Sections {
		if section.ResourceType != resourceTypeNode {
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

func convertSectionToUnattachedNode(section *parser.GdResource) (*godot.Node, error) {
	if section == nil {
		return nil, fmt.Errorf("section was nil")
	}

	if section.ResourceType != resourceTypeNode {
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

	insertFieldEntriesFromSection(section, node.Fields)

	return &node, nil
}

func attachTypeToNode(node *godot.Node, section *parser.GdResource) error {
	if nodeType, err := section.GetAttribute("type"); err == nil {
		node.Type = *nodeType.String
		return nil
	}

	if instance, err := section.GetAttribute("instance"); err == nil {
		if len(instance.Type.Parameters) != 1 {
			return fmt.Errorf("node instance parameter does not contain a valid reference %v", instance.Raw())
		}

		node.Instance = convertGdValue(instance).(godot.Type)
		return nil
	}

	// nodes don't have to have a type, if it doesn't have one just don't do anything
	return nil
}

func convertSectionToEditable(section *parser.GdResource) (*godot.Editable, error) {
	if section.ResourceType != resourceTypeEditable {
		return nil, fmt.Errorf("you can't convert a %s to editable", section.ResourceType)
	}

	path, err := section.GetAttribute("path")
	if err != nil {
		return nil, errors.Wrap(err, "editable without path")
	}

	editable := godot.Editable{
		Path: *path.String,
		MetaData: godot.MetaData{
			LexerPosition: section.Pos,
		},
	}

	return &editable, nil
}

func convertSectionToConnection(section *parser.GdResource) (*godot.Connection, error) {
	if section.ResourceType != resourceTypeConnection {
		return nil, fmt.Errorf("you can't convert a %s to connection", section.ResourceType)
	}

	from, err := section.GetAttribute("from")
	if err != nil {
		return nil, errors.Wrap(err, "editable without from")
	}

	to, err := section.GetAttribute("to")
	if err != nil {
		return nil, errors.Wrap(err, "editable without to")
	}

	signal, err := section.GetAttribute("signal")
	if err != nil {
		return nil, errors.Wrap(err, "editable without signal")
	}

	method, err := section.GetAttribute("method")
	if err != nil {
		return nil, errors.Wrap(err, "editable without method")
	}

	conn := godot.Connection{
		From:   *from.String,
		To:     *to.String,
		Signal: *signal.String,
		Method: *method.String,
		MetaData: godot.MetaData{
			LexerPosition: section.Pos,
		},
	}

	if flags, err := section.GetAttribute("flags"); err == nil {
		conn.Flags = *flags.Integer
	}

	if binds, err := section.GetAttribute("binds"); err == nil {
		bindsArray := convertGdValue(binds).(godot.Value)
		conn.Binds = bindsArray
	}

	return &conn, nil
}
