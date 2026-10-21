package microruntime

// Object is a Kubernetes object as decoded from JSON or YAML.
//
// Accessors read a missing or mistyped field as its zero value, so a malformed
// object fails the checks that inspect it instead of panicking.
type Object map[string]any

// Text returns a string field.
func (object Object) Text(key string) string {
	text, _ := object[key].(string)

	return text
}

// Number returns a JSON number field.
func (object Object) Number(key string) float64 {
	number, _ := object[key].(float64)

	return number
}

// Flag returns a boolean field.
func (object Object) Flag(key string) bool {
	flag, _ := object[key].(bool)

	return flag
}

// Child returns an object field.
func (object Object) Child(key string) Object { return asObject(object[key]) }

// Children returns an array field, one entry per element. Elements that are
// not objects become nil, so lookups on them read zero values.
func (object Object) Children(key string) []Object {
	elements := asSlice(object[key])

	children := make([]Object, len(elements))
	for index, element := range elements {
		children[index] = asObject(element)
	}

	return children
}

func (object Object) Kind() string       { return object.Text("kind") }
func (object Object) APIVersion() string { return object.Text("apiVersion") }
func (object Object) Metadata() Object   { return object.Child("metadata") }
func (object Object) Name() string       { return object.Metadata().Text("name") }
func (object Object) Namespace() string  { return object.Metadata().Text("namespace") }
func (object Object) UID() string        { return object.Metadata().Text("uid") }
func (object Object) Labels() Object     { return object.Metadata().Child("labels") }
func (object Object) Spec() Object       { return object.Child("spec") }
func (object Object) Status() Object     { return object.Child("status") }

// PodSpec returns the pod spec of a workload's pod template.
func (object Object) PodSpec() Object {
	return object.Spec().Child("template").Child("spec")
}

// Containers returns a pod spec's containers followed by its init containers.
func (object Object) Containers() []Object {
	return append(object.Children("containers"), object.Children("initContainers")...)
}

// OwnerUIDs returns the UIDs of every owner reference.
func (object Object) OwnerUIDs() []string { return object.ownerUIDs(false) }

// ControllerUIDs returns the UIDs of owner references marked as controller.
func (object Object) ControllerUIDs() []string { return object.ownerUIDs(true) }

func (object Object) ownerUIDs(controllersOnly bool) []string {
	var uids []string

	for _, owner := range object.Metadata().Children("ownerReferences") {
		if controllersOnly && !owner.Flag("controller") {
			continue
		}

		uids = append(uids, owner.Text("uid"))
	}

	return uids
}

// IsOwnedBy reports whether object is the recorded one, identified by uid,
// and still carries the execution's ownership token.
func (object Object) IsOwnedBy(uid, token string) bool {
	return uid != "" &&
		object.UID() == uid &&
		object.Labels().Text(LabelOwnershipToken) == token
}

// ensureChild returns an object field, creating it when absent.
func (object Object) ensureChild(key string) Object {
	if child := object.Child(key); child != nil {
		return child
	}

	child := map[string]any{}
	object[key] = child

	return child
}

func asObject(value any) Object {
	switch typed := value.(type) {
	case Object:
		return typed
	case map[string]any:
		return typed
	}

	return nil
}

func asSlice(value any) []any {
	slice, _ := value.([]any)

	return slice
}

// strictSlice tells an absent array, which is valid, apart from a mistyped one.
func strictSlice(value any) ([]any, bool) {
	if value == nil {
		return nil, true
	}

	slice, ok := value.([]any)

	return slice, ok
}
