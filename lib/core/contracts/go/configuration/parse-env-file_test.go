package configuration

import "testing"

func TestParseEnvFile(t *testing.T) {
	variables, err := ParseEnvFile("# the log\nLOG_URL=http://log:3500\n\nPORT=3402\nEMPTY=\nQUERY=a=b\n")
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]string{
		"LOG_URL": "http://log:3500",
		"PORT":    "3402",
		"EMPTY":   "",
		"QUERY":   "a=b",
	}
	if len(variables) != len(want) {
		t.Fatalf("variables = %v, want %v", variables, want)
	}

	for name, value := range want {
		if variables[name] != value {
			t.Fatalf("%s = %q, want %q", name, variables[name], value)
		}
	}
}
