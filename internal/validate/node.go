package validate

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
)

func validatorFirstNodeHasNoParent(tscnFile *parser.TscnFile) error {
	for _, section := range tscnFile.Sections {
		if section.ResourceType != parser.ResourceTypeNode {
			continue
		}

		_, err := section.GetAttribute("parent")
		if err == nil {
			return fmt.Errorf(
				"the first node in the file, which is also the scene root, must not have a 'parent' attribute. %s",
				section.Pos,
			)
		}
		// we only care about the first node
		return nil
	}

	return nil
}

func validatorOnlyOneRootNode(tscnFile *parser.TscnFile) error {
	foundRoot := false
	for _, section := range tscnFile.Sections {
		if section.ResourceType != parser.ResourceTypeNode {
			continue
		}

		_, err := section.GetAttribute("parent")
		if err != nil {
			if foundRoot {
				return fmt.Errorf("found a second root node (a node without parent): %s", section.Pos)
			}

			foundRoot = true
		}
	}

	return nil
}

func validatorTestIfAllResourceReferencesActuallyExist(tscnFile *parser.TscnFile) error {
	types := findTypeReferencesInTscnFile(tscnFile)

	for _, t := range types {
		err := doesTypeReferenceExistInTscnFile(tscnFile, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func findTypeReferencesInTscnFile(tscnFile *parser.TscnFile) (types []*parser.GdType) {
	t := findTypeReferencesInGdFieldSlice(tscnFile.Attributes)
	t = append(t, findTypeReferencesInGdFieldSlice(tscnFile.Fields)...)

	for _, section := range tscnFile.Sections {
		t = append(t, findTypeReferencesInGdFieldSlice(section.Attributes)...)
		t = append(t, findTypeReferencesInGdFieldSlice(section.Fields)...)
	}

	for _, typeElem := range t {
		types = append(types, typeElem)
		types = append(types, findTypeReferencesInGdValueSlice(typeElem.Parameters)...)
	}

	return types
}

func findTypeReferencesInGdFieldSlice(slice []*parser.GdField) (types []*parser.GdType) {
	for _, elem := range slice {
		types = append(types, findTypeReferenceInGdValue(elem.Value)...)
	}
	return types
}

func findTypeReferencesInGdValueSlice(slice []*parser.GdValue) (types []*parser.GdType) {
	for _, v := range slice {
		types = append(types, findTypeReferenceInGdValue(v)...)
	}
	return types
}

func findTypeReferencesInGdKeyValueSlice(slice []*parser.GdMapField) (types []*parser.GdType) {
	for _, elem := range slice {
		types = append(types, findTypeReferenceInGdValue(elem.Value)...)
	}
	return types
}

func findTypeReferenceInGdValue(v *parser.GdValue) (types []*parser.GdType) {
	if v.Type != nil {
		types = append(types, v.Type)
	}

	if (v.IsEmptyMap != nil && *v.IsEmptyMap) || len(v.Map) != 0 {
		types = append(types, findTypeReferencesInGdKeyValueSlice(v.Map)...)
	}

	if (v.IsEmptyArray != nil && *v.IsEmptyArray) || len(v.Array) != 0 {
		types = append(types, findTypeReferencesInGdValueSlice(v.Array)...)
	}

	return types
}

func doesTypeReferenceExistInTscnFile(tscnFile *parser.TscnFile, typeRef *parser.GdType) error {
	// if type reference isn't ExtResource/SubResource we don't need to check
	if typeRef.Key != "ExtResource" && typeRef.Key != "SubResource" {
		return nil
	}

	// type references have to be ExtResource(ID) or SubResource(ID)
	if len(typeRef.Parameters) != 1 {
		return fmt.Errorf(
			"type reference %s(%v) is not a valid type reference",
			typeRef.Key,
			convertGdValuesIntoString(typeRef.Parameters),
		)
	}

	for _, section := range tscnFile.Sections {
		if section.ResourceType != parser.ResourceTypeExtResource && section.ResourceType != parser.ResourceTypeSubResource {
			continue
		}

		idValue, err := section.GetAttribute("id")
		if err != nil {
			return errors.Wrapf(err, "could not retrieve attribute idValue %s", typeRef.Pos)
		}
		if idValue.Integer == nil {
			return fmt.Errorf("id field must be integer %s", typeRef.Pos)
		}

		typeRefID := *typeRef.Parameters[0].Integer
		extResourceID := *idValue.Integer

		// found referenced element, stop...
		if typeRefID == extResourceID {
			return nil
		}
	}

	return fmt.Errorf(
		"could not find type reference %s(%v) %s",
		typeRef.Key,
		convertGdValuesIntoString(typeRef.Parameters),
		typeRef.Pos,
	)
}

func convertGdValuesIntoString(values []*parser.GdValue) string {
	var parts []string
	for _, v := range values {
		parts = append(parts, v.ToString())
	}
	return strings.Join(parts, ", ")
}
