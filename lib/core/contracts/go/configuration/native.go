package configuration

// Native preserves an existing v2 descriptor and its owned resources on migration.
type Native struct {
	RuntimeID string `json:"runtimeId"`
}
