package microruntime

// secretKind is the one kind the LangDrift CLI creates that no manifest declares: a
// unit's environment, written at start from values that never reach a file.
const secretKind = "Secret"

// Environment is what one unit's container receives as its environment:
// the component's env file and the variables it passes through from the
// agent. the LangDrift CLI writes it into one Secret at start and records the Secret, never
// its values.
type Environment struct {
	Unit      Identifier
	Secret    Identifier
	Variables map[string]string
}

// Manifest is the Secret as a manifest of the unit, for the startup to
// claim, create and record with the unit's other objects.
func (environment Environment) Manifest() Manifest {
	data := make(map[string]any, len(environment.Variables))
	for key, value := range environment.Variables {
		data[key] = value
	}

	return Manifest{
		Unit: environment.Unit,
		Object: Object{
			"apiVersion": "v1",
			"kind":       secretKind,
			"metadata":   map[string]any{"name": environment.Secret.String()},
			"type":       "Opaque",
			"stringData": data,
		},
	}
}
