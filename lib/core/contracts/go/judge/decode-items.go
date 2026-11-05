package judge

import (
	"encoding/json"
	"strconv"
)

// decodeItems reads a JSON array of objects, each through decodeObject,
// checks every item with its path, such as sources[1], and returns them.
func decodeItems[Item any](
	raw json.RawMessage,
	path string,
	fields []string,
	check func(path string, item Item) error,
) ([]Item, error) {
	var rawItems []json.RawMessage
	if err := json.Unmarshal(raw, &rawItems); err != nil {
		return nil, refuse(path, "must be an array")
	}

	items := make([]Item, 0, len(rawItems))
	for index, rawItem := range rawItems {
		itemPath := path + "[" + strconv.Itoa(index) + "]"

		var item Item
		if err := decodeObject(rawItem, itemPath, fields, &item); err != nil {
			return nil, err
		}

		if err := check(itemPath, item); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
