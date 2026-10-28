package application

import "context"

// ValidateRuntime checks a runtime's descriptor and native files without
// provisioning anything or contacting a cluster.
type ValidateRuntime struct {
	Projects      ProjectOpener
	Configuration RuntimeConfiguration
}

// ValidationResult summarizes a valid configuration.
type ValidationResult struct {
	Runtime string   `json:"runtime"`
	Units   []string `json:"units"`
	Objects int      `json:"objects"`
	Status  string   `json:"status"`
}

// Execute loads and validates the configuration.
func (useCase ValidateRuntime) Execute(_ context.Context, request RuntimeRequest) (ValidationResult, error) {
	runtime, err := parseRuntime(request.Runtime)
	if err != nil {
		return ValidationResult{}, err
	}

	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return ValidationResult{}, err
	}

	descriptor, snapshot, err := loadConfiguration(useCase.Configuration, owner, runtime)
	if err != nil {
		return ValidationResult{}, err
	}

	units := make([]string, 0, len(descriptor.Units))
	for _, name := range descriptor.UnitNames() {
		units = append(units, name.String())
	}

	return ValidationResult{
		Runtime: runtime.String(),
		Units:   units,
		Objects: len(snapshot.Manifests),
		Status:  "valid configuration; execution unverified",
	}, nil
}
