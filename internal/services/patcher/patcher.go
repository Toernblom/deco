package patcher

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Toernblom/deco/internal/domain"
)

// PatchOperation represents a single patch operation in a batch.
type PatchOperation struct {
	Op    string      `json:"op"`    // Operation type: "set", "append", "unset"
	Path  string      `json:"path"`  // Field path (supports dot notation and array indices)
	Value interface{} `json:"value"` // Value to set or append (not used for unset)
}

// Patcher handles patch operations on nodes.
type Patcher struct {
	// List of required fields that cannot be unset
	requiredFields map[string]bool
}

// New creates a new Patcher instance.
func New() *Patcher {
	return &Patcher{
		requiredFields: map[string]bool{
			"id":      true,
			"kind":    true,
			"version": true,
			"status":  true,
			"title":   true,
		},
	}
}

// Set sets a field value at the given path.
// Path can be simple ("title") or use dot notation ("meta.title").
// Array indices are supported using bracket notation ("tags[0]").
func (p *Patcher) Set(node *domain.Node, path string, value interface{}) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	parts := parsePath(path)
	if len(parts) == 0 {
		return fmt.Errorf("empty path")
	}

	return p.setValue(reflect.ValueOf(node).Elem(), parts, value)
}

// Append appends one or more values to an array field.
// Returns an error if the field is not an array.
func (p *Patcher) Append(node *domain.Node, path string, values ...interface{}) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	parts := parsePath(path)
	if len(parts) == 0 {
		return fmt.Errorf("empty path")
	}

	// Get the field
	field, err := p.getField(reflect.ValueOf(node).Elem(), parts)
	if err != nil {
		return err
	}

	// Check if it's a slice
	if field.Kind() != reflect.Slice {
		return fmt.Errorf("cannot append to non-array field at path %q", path)
	}

	// Append values
	for _, value := range values {
		val := reflect.ValueOf(value)

		// Convert value to match slice element type
		if val.Type() != field.Type().Elem() {
			val = val.Convert(field.Type().Elem())
		}

		field.Set(reflect.Append(field, val))
	}

	return nil
}

// Unset removes a field or array element at the given path.
// For arrays, use bracket notation ("tags[1]") to remove a specific element.
func (p *Patcher) Unset(node *domain.Node, path string) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	parts := parsePath(path)
	if len(parts) == 0 {
		return fmt.Errorf("empty path")
	}

	// Check if trying to unset a required field
	if len(parts) == 1 && p.requiredFields[strings.ToLower(parts[0])] {
		return fmt.Errorf("cannot unset required field %q", parts[0])
	}

	return p.unsetValue(reflect.ValueOf(node).Elem(), parts)
}

// Apply applies a batch of patch operations to a node.
// Operations are applied in order. If any operation fails, all changes are rolled back.
// This provides transactional semantics for batch operations.
func (p *Patcher) Apply(node *domain.Node, operations []PatchOperation) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	// If no operations, nothing to do
	if len(operations) == 0 {
		return nil
	}

	// Create a snapshot for rollback
	snapshot, err := p.snapshot(node)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Apply each operation in order
	for i, op := range operations {
		var err error

		switch op.Op {
		case "set":
			err = p.Set(node, op.Path, op.Value)
		case "append":
			// For append, we need to handle single value vs array of values
			err = p.Append(node, op.Path, op.Value)
		case "unset":
			err = p.Unset(node, op.Path)
		default:
			err = fmt.Errorf("unknown operation type %q", op.Op)
		}

		if err != nil {
			// Rollback on error
			p.restore(node, snapshot)
			return fmt.Errorf("operation %d (%s %s) failed: %w", i, op.Op, op.Path, err)
		}
	}

	return nil
}

// snapshot creates a deep copy of a node for rollback purposes
func (p *Patcher) snapshot(node *domain.Node) (*domain.Node, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	// Encode the node
	if err := enc.Encode(node); err != nil {
		return nil, err
	}

	// Decode into a new node
	var snapshot domain.Node
	if err := dec.Decode(&snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

// restore restores a node from a snapshot
func (p *Patcher) restore(node *domain.Node, snapshot *domain.Node) {
	*node = *snapshot
}

// setValue sets a value at the given path parts
func (p *Patcher) setValue(v reflect.Value, parts []string, value interface{}) error {
	if len(parts) == 0 {
		return fmt.Errorf("empty path")
	}

	// Handle the current part
	part := parts[0]

	// Check for array index notation (e.g., "tags[0]")
	if fieldName, idx, isArray := parseArrayIndex(part); isArray {
		// First get the field
		field := v.FieldByName(capitalizeFirst(fieldName))
		if !field.IsValid() {
			return fmt.Errorf("field %q not found", fieldName)
		}

		if field.Kind() != reflect.Slice {
			return fmt.Errorf("expected slice for field %q, got %v", fieldName, field.Kind())
		}

		if idx < 0 || idx >= field.Len() {
			return fmt.Errorf("index %d out of bounds for slice of length %d", idx, field.Len())
		}

		elem := field.Index(idx)

		if len(parts) == 1 {
			// Last part, set the value
			val := reflect.ValueOf(value)
			if val.Type() != elem.Type() && val.Type().ConvertibleTo(elem.Type()) {
				val = val.Convert(elem.Type())
			}
			elem.Set(val)
			return nil
		}

		// More parts to traverse
		return p.setValue(elem, parts[1:], value)
	}

	// Regular field access
	field := v.FieldByName(capitalizeFirst(part))
	if !field.IsValid() {
		return fmt.Errorf("field %q not found", part)
	}

	if !field.CanSet() {
		return fmt.Errorf("field %q cannot be set", part)
	}

	if len(parts) == 1 {
		// Last part, set the value
		val := reflect.ValueOf(value)

		// Handle type conversion
		if val.Type() != field.Type() {
			if val.Type().ConvertibleTo(field.Type()) {
				val = val.Convert(field.Type())
			} else {
				return fmt.Errorf("cannot convert %v to %v", val.Type(), field.Type())
			}
		}

		field.Set(val)
		return nil
	}

	// More parts to traverse
	// Handle maps
	if field.Kind() == reflect.Map {
		if field.IsNil() {
			field.Set(reflect.MakeMap(field.Type()))
		}

		key := reflect.ValueOf(parts[1])
		if len(parts) == 2 {
			// Set map value
			val := reflect.ValueOf(value)
			field.SetMapIndex(key, val)
			return nil
		}

		return fmt.Errorf("nested path in map not supported yet")
	}

	// Traverse to next level
	return p.setValue(field, parts[1:], value)
}

// getField retrieves a field value at the given path
func (p *Patcher) getField(v reflect.Value, parts []string) (reflect.Value, error) {
	if len(parts) == 0 {
		return v, nil
	}

	part := parts[0]

	// Check for array index notation
	if fieldName, idx, isArray := parseArrayIndex(part); isArray {
		field := v.FieldByName(capitalizeFirst(fieldName))
		if !field.IsValid() {
			return reflect.Value{}, fmt.Errorf("field %q not found", fieldName)
		}
		if field.Kind() != reflect.Slice {
			return reflect.Value{}, fmt.Errorf("expected slice for field %q, got %v", fieldName, field.Kind())
		}
		if idx < 0 || idx >= field.Len() {
			return reflect.Value{}, fmt.Errorf("index out of bounds")
		}
		return p.getField(field.Index(idx), parts[1:])
	}

	// Regular field access
	field := v.FieldByName(capitalizeFirst(part))
	if !field.IsValid() {
		return reflect.Value{}, fmt.Errorf("field %q not found", part)
	}

	if len(parts) == 1 {
		return field, nil
	}

	// Handle maps
	if field.Kind() == reflect.Map {
		key := reflect.ValueOf(parts[1])
		return field.MapIndex(key), nil
	}

	return p.getField(field, parts[1:])
}

// unsetValue unsets a value at the given path
func (p *Patcher) unsetValue(v reflect.Value, parts []string) error {
	if len(parts) == 0 {
		return fmt.Errorf("empty path")
	}

	part := parts[0]

	// Check for array index notation
	if fieldName, idx, isArray := parseArrayIndex(part); isArray {
		field := v.FieldByName(capitalizeFirst(fieldName))
		if !field.IsValid() {
			return fmt.Errorf("field %q not found", fieldName)
		}

		if field.Kind() != reflect.Slice {
			return fmt.Errorf("expected slice for field %q, got %v", fieldName, field.Kind())
		}

		if idx < 0 || idx >= field.Len() {
			return nil // Index doesn't exist, nothing to unset
		}

		if len(parts) == 1 {
			// Remove element from slice
			newSlice := reflect.MakeSlice(field.Type(), 0, field.Len()-1)
			newSlice = reflect.AppendSlice(newSlice, field.Slice(0, idx))
			if idx+1 < field.Len() {
				newSlice = reflect.AppendSlice(newSlice, field.Slice(idx+1, field.Len()))
			}
			field.Set(newSlice)
			return nil
		}

		return p.unsetValue(field.Index(idx), parts[1:])
	}

	// Regular field access
	field := v.FieldByName(capitalizeFirst(part))
	if !field.IsValid() {
		return fmt.Errorf("field %q not found", part)
	}

	if !field.CanSet() {
		return fmt.Errorf("field %q cannot be set", part)
	}

	if len(parts) == 1 {
		// Last part, unset the value
		field.Set(reflect.Zero(field.Type()))
		return nil
	}

	// More parts to traverse
	if field.Kind() == reflect.Map {
		if field.IsNil() {
			return nil // Map doesn't exist, nothing to unset
		}

		key := reflect.ValueOf(parts[1])
		if len(parts) == 2 {
			// Remove from map
			field.SetMapIndex(key, reflect.Value{})
			return nil
		}

		return fmt.Errorf("nested path in map not fully supported")
	}

	return p.unsetValue(field, parts[1:])
}

// parsePath splits a path into parts
// Examples:
//
//	"title" -> ["title"]
//	"meta.title" -> ["meta", "title"]
//	"tags[0]" -> ["tags[0]"]
func parsePath(path string) []string {
	return strings.Split(path, ".")
}

// parseArrayIndex checks if a part contains array index notation
// Returns the field name, index, and true if it's an array access
// Example: "tags[0]" -> "tags", 0, true
func parseArrayIndex(part string) (string, int, bool) {
	if !strings.Contains(part, "[") {
		return "", -1, false
	}

	start := strings.Index(part, "[")
	end := strings.Index(part, "]")

	if start == -1 || end == -1 || end <= start {
		return "", -1, false
	}

	fieldName := part[:start]
	indexStr := part[start+1 : end]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", -1, false
	}

	return fieldName, index, true
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
