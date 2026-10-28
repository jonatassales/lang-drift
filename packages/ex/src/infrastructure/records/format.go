// Package records writes the records and mirrors of the .langdrift folder in
// the layout the consumer's own formatter keeps, so a formatted file never
// changes its bytes, and the JSON fences of the Markdown records.
package records

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// Profile is one layout: how to indent, how wide a line may be, and whether a
// small object may sit on one line.
type Profile struct {
	// Indent is written once per level.
	Indent string
	// IndentColumns is how many columns one level takes when a line's width
	// is measured; a tab counts as the consumer's indent width.
	IndentColumns int
	// Width is the widest line, in columns.
	Width int
	// CollapseObjects lets an object whose values are all scalars sit on one
	// line when it fits, as the fences do.
	CollapseObjects bool
}

// FileProfile is Biome's layout for the JSON files: tabs, lines up to 80
// columns, every object expanded, a scalar array on one line when it fits.
var FileProfile = Profile{Indent: "\t", IndentColumns: 2, Width: 80}

// FenceProfile is the layout of a JSON fence inside Markdown: two spaces,
// and a small scalar-only object or array on one line when it fits.
var FenceProfile = Profile{Indent: "  ", IndentColumns: 2, Width: 80, CollapseObjects: true}

// Format renders a record in the profile's layout, keys in the order the
// type declares them, nothing escaped that JSON does not require, and a
// trailing newline.
func Format(value any, profile Profile) ([]byte, error) {
	compact, err := encodeCompact(value)
	if err != nil {
		return nil, err
	}

	tree, err := parse(compact)
	if err != nil {
		return nil, err
	}

	var out strings.Builder

	render(&out, tree, profile, 0, "", false)
	out.WriteByte('\n')

	return []byte(out.String()), nil
}

// encodeCompact serializes a value without HTML escaping, so a record may
// hold "<" or "&" as the consumer wrote it.
func encodeCompact(value any) ([]byte, error) {
	var buffer bytes.Buffer

	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(value); err != nil {
		return nil, err
	}

	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}

type nodeKind int

const (
	scalarNode nodeKind = iota
	objectNode
	arrayNode
)

// node is one JSON value with its members in document order.
type node struct {
	kind   nodeKind
	raw    string
	keys   []string
	values []*node
}

func parse(compact []byte) (*node, error) {
	decoder := json.NewDecoder(bytes.NewReader(compact))
	decoder.UseNumber()

	tree, err := parseValue(decoder)
	if err != nil {
		return nil, err
	}

	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return nil, errors.New("a record holds one JSON value")
	}

	return tree, nil
}

func parseValue(decoder *json.Decoder) (*node, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	delimiter, isDelimiter := token.(json.Delim)
	if !isDelimiter {
		return scalar(token)
	}

	switch delimiter {
	case '{':
		return parseObject(decoder)

	case '[':
		return parseArray(decoder)
	}

	return nil, fmt.Errorf("unexpected %v", delimiter)
}

func parseObject(decoder *json.Decoder) (*node, error) {
	object := &node{kind: objectNode}

	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		key, isString := token.(string)
		if !isString {
			return nil, fmt.Errorf("object key %v is not a string", token)
		}

		value, err := parseValue(decoder)
		if err != nil {
			return nil, err
		}

		object.keys = append(object.keys, key)
		object.values = append(object.values, value)
	}

	// The closing brace.
	_, err := decoder.Token()

	return object, err
}

func parseArray(decoder *json.Decoder) (*node, error) {
	array := &node{kind: arrayNode}

	for decoder.More() {
		item, err := parseValue(decoder)
		if err != nil {
			return nil, err
		}

		array.values = append(array.values, item)
	}

	// The closing bracket.
	_, err := decoder.Token()

	return array, err
}

func scalar(token json.Token) (*node, error) {
	switch value := token.(type) {
	case nil:
		return &node{kind: scalarNode, raw: "null"}, nil

	case bool:
		if value {
			return &node{kind: scalarNode, raw: "true"}, nil
		}

		return &node{kind: scalarNode, raw: "false"}, nil

	case json.Number:
		return &node{kind: scalarNode, raw: value.String()}, nil

	case string:
		quoted, err := encodeCompact(value)
		if err != nil {
			return nil, err
		}

		return &node{kind: scalarNode, raw: string(quoted)}, nil
	}

	return nil, fmt.Errorf("unexpected token %v", token)
}

// render writes one value at depth. prefix is what the line already holds
// before the value, the key, and trailing says whether a comma follows, so
// the width check measures the whole line.
func render(out *strings.Builder, value *node, profile Profile, depth int, prefix string, trailing bool) {
	switch value.kind {
	case scalarNode:
		out.WriteString(value.raw)

	case arrayNode:
		renderArray(out, value, profile, depth, prefix, trailing)

	case objectNode:
		renderObject(out, value, profile, depth, prefix, trailing)
	}
}

func renderArray(out *strings.Builder, array *node, profile Profile, depth int, prefix string, trailing bool) {
	if len(array.values) == 0 {
		out.WriteString("[]")

		return
	}

	if allScalars(array) {
		flat := "[" + joinScalars(array, ", ") + "]"
		if fits(profile, depth, prefix, flat, trailing) {
			out.WriteString(flat)

			return
		}
	}

	out.WriteString("[\n")

	for index, item := range array.values {
		last := index == len(array.values)-1

		writeIndent(out, profile, depth+1)
		render(out, item, profile, depth+1, "", !last)
		writeSeparator(out, last)
	}

	writeIndent(out, profile, depth)
	out.WriteString("]")
}

func renderObject(out *strings.Builder, object *node, profile Profile, depth int, prefix string, trailing bool) {
	if len(object.keys) == 0 {
		out.WriteString("{}")

		return
	}

	if profile.CollapseObjects && allScalars(object) {
		flat := "{ " + joinMembers(object) + " }"
		if fits(profile, depth, prefix, flat, trailing) {
			out.WriteString(flat)

			return
		}
	}

	out.WriteString("{\n")

	for index, key := range object.keys {
		last := index == len(object.keys)-1
		member := quote(key) + ": "

		writeIndent(out, profile, depth+1)
		out.WriteString(member)
		render(out, object.values[index], profile, depth+1, member, !last)
		writeSeparator(out, last)
	}

	writeIndent(out, profile, depth)
	out.WriteString("}")
}

// fits reports whether a line of indent, prefix, flat value and a possible
// trailing comma stays inside the profile's width.
func fits(profile Profile, depth int, prefix, flat string, trailing bool) bool {
	width := depth*profile.IndentColumns + utf8.RuneCountInString(prefix) + utf8.RuneCountInString(flat)
	if trailing {
		width++
	}

	return width <= profile.Width
}

func allScalars(value *node) bool {
	for _, member := range value.values {
		if member.kind != scalarNode {
			return false
		}
	}

	return true
}

func joinScalars(array *node, separator string) string {
	parts := make([]string, len(array.values))
	for index, item := range array.values {
		parts[index] = item.raw
	}

	return strings.Join(parts, separator)
}

func joinMembers(object *node) string {
	parts := make([]string, len(object.keys))
	for index, key := range object.keys {
		parts[index] = quote(key) + ": " + object.values[index].raw
	}

	return strings.Join(parts, ", ")
}

func quote(key string) string {
	quoted, err := encodeCompact(key)
	if err != nil {
		return `"` + key + `"`
	}

	return string(quoted)
}

func writeIndent(out *strings.Builder, profile Profile, depth int) {
	for range depth {
		out.WriteString(profile.Indent)
	}
}

func writeSeparator(out *strings.Builder, last bool) {
	if !last {
		out.WriteString(",")
	}

	out.WriteString("\n")
}
