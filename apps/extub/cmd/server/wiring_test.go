package main

import "testing"

const longToken = "a-shared-preview-token-of-forty-characters"

func TestListenAddressFollowsTheMode(t *testing.T) {
	for _, check := range []struct {
		settings config
		address  string
	}{
		{config{}, "127.0.0.1:3401"},
		{config{extubPort: "3402", port: "8080"}, "127.0.0.1:3402"},
		{config{token: longToken}, ":3000"},
		{config{token: longToken, port: "8080", extubPort: "3402"}, ":8080"},
	} {
		address, err := listenAddress(check.settings)
		if err != nil || address != check.address {
			t.Fatalf("%+v listens on %q (%v), want %q", check.settings, address, err, check.address)
		}
	}
	for _, settings := range []config{{extubPort: "0"}, {token: longToken, port: "70000"}, {extubPort: "x"}, {token: "too-short"}} {
		if _, err := listenAddress(settings); err == nil {
			t.Fatalf("%+v accepted", settings)
		}
	}
}

func TestConfigReadsTheEnvironment(t *testing.T) {
	environment := map[string]string{"EXTUB_TOKEN": " " + longToken + "\n", "PORT": "8080", "EXTUB_PORT": "3402"}
	settings := configFromEnv(func(name string) string { return environment[name] })
	if settings != (config{token: longToken, port: "8080", extubPort: "3402"}) {
		t.Fatalf("read %+v", settings)
	}
}
