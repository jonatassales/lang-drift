package configuration

func validConfig() Configuration {
	return Configuration{
		Version:    2,
		RuntimeUID: "runtime-uid",
		RuntimeID:  "runtime",
		ProjectID:  "project",
		Revision:   1,
		Name:       "Example",
		CPUMillis:  1000,
		MemoryMiB:  512,
		Model:      "test-model",
		Effort:     "medium",
		Components: []Component{
			{
				Name: "api",
				Path: ".",
			},
		},
	}
}
