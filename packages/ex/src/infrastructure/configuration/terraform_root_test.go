package configuration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/infrastructure/terraform"
)

// generatedRootDocument is the root generation writes, as a JSON document
// that each case can change.
func generatedRootDocument(t *testing.T) map[string]any {
	t.Helper()

	data, err := json.Marshal(generateTerraformRoot(identifier(t, "demo"), microruntime.ResourceQuotas{
		microruntime.QuotaPods: quantity(t, "10"),
	}))
	if err != nil {
		t.Fatal(err)
	}

	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}

	return document
}

func writeRoot(t *testing.T, files map[string]string) string {
	t.Helper()

	directory := tempDir(t)
	for name, content := range files {
		write(t, filepath.Join(directory, name), content)
	}

	return directory
}

func encode(t *testing.T, document map[string]any) string {
	t.Helper()

	data, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func TestGenerateTerraformRoot(t *testing.T) {
	t.Parallel()

	document := generatedRootDocument(t)

	digests, err := validateTerraformRoot(writeRoot(t, map[string]string{"main.tf.json": encode(t, document)}))
	if err != nil {
		t.Fatalf("the generated root violates the policy: %v", err)
	}

	if len(digests) != 1 || len(digests["main.tf.json"]) != 64 {
		t.Errorf("digests = %v, want one digest per source", digests)
	}

	resources := document["resource"].(map[string]any)
	namespace := resources[terraform.NamespaceResource].(map[string]any)["runtime"].(map[string]any)["metadata"].([]any)[0].(map[string]any)
	labels := namespace["labels"].(map[string]any)

	if namespace["name"] != "${var.namespace}" || labels[microruntime.LabelOwnershipToken] != "${var.ownership_token}" ||
		labels[microruntime.LabelRuntime] != "demo" {
		t.Errorf("namespace metadata = %v, want the namespace and token variables", namespace)
	}

	quota := resources[terraform.QuotaResource].(map[string]any)["runtime"].(map[string]any)
	hard := quota["spec"].([]any)[0].(map[string]any)["hard"].(map[string]any)
	if hard["pods"] != "10" {
		t.Errorf("quota hard = %v, want the blueprint's quotas", hard)
	}

	variables := document["variable"].(map[string]any)
	for _, name := range terraform.Variables {
		if variables[name] == nil {
			t.Errorf("variable %s is not declared", name)
		}
	}
}

func TestValidateTerraformRoot(t *testing.T) {
	t.Parallel()

	edit := func(change func(map[string]any)) func(*testing.T) map[string]string {
		return func(t *testing.T) map[string]string {
			document := generatedRootDocument(t)
			change(document)

			return map[string]string{"main.tf.json": encode(t, document)}
		}
	}
	provider := func(document map[string]any) map[string]any {
		return document["provider"].(map[string]any)["kubernetes"].(map[string]any)
	}
	namespaceMetadata := func(document map[string]any) map[string]any {
		return document["resource"].(map[string]any)[terraform.NamespaceResource].(map[string]any)["runtime"].(map[string]any)["metadata"].([]any)[0].(map[string]any)
	}

	tests := []struct {
		name    string
		files   func(*testing.T) map[string]string
		wantErr bool
	}{
		{name: "generated root", files: edit(func(map[string]any) {})},
		{
			name: "outputs and locals",
			files: edit(func(d map[string]any) {
				d["output"] = map[string]any{"namespace": map[string]any{"value": "${var.namespace}"}}
				d["locals"] = map[string]any{"owner": "${var.owner}"}
			}),
		},
		{
			name: "provider lock file and Terraform data beside the sources",
			files: func(t *testing.T) map[string]string {
				files := edit(func(map[string]any) {})(t)
				files[".terraform.lock.hcl"] = "# lock"
				files[".terraform/providers/README"] = "cache"

				return files
			},
		},
		{name: "no sources", files: func(*testing.T) map[string]string { return map[string]string{"README.md": "# notes"} }, wantErr: true},
		{
			name: "HCL source",
			files: func(t *testing.T) map[string]string {
				files := edit(func(map[string]any) {})(t)
				files["extra.tf"] = `resource "null_resource" "x" {}`

				return files
			},
			wantErr: true,
		},
		{
			name: "variable file",
			files: func(t *testing.T) map[string]string {
				files := edit(func(map[string]any) {})(t)
				files["terraform.tfvars.json"] = `{"namespace": "elsewhere"}`

				return files
			},
			wantErr: true,
		},
		{
			name: "module directory",
			files: func(t *testing.T) map[string]string {
				files := edit(func(map[string]any) {})(t)
				files["modules/main.tf.json"] = `{}`

				return files
			},
			wantErr: true,
		},
		{
			name: "second provider configuration",
			files: func(t *testing.T) map[string]string {
				files := edit(func(map[string]any) {})(t)
				files["second.tf.json"] = `{"provider": {"kubernetes": {"config_path": "${var.kubeconfig}", "config_context": "${var.context}"}}}`

				return files
			},
			wantErr: true,
		},
		{name: "no provider configuration", files: edit(func(d map[string]any) { delete(d, "provider") }), wantErr: true},
		{name: "module block", files: edit(func(d map[string]any) { d["module"] = map[string]any{} }), wantErr: true},
		{name: "data block", files: edit(func(d map[string]any) { d["data"] = map[string]any{} }), wantErr: true},
		{name: "import block", files: edit(func(d map[string]any) { d["import"] = []any{} }), wantErr: true},
		{name: "remote execution", files: edit(func(d map[string]any) { d["terraform"].(map[string]any)["cloud"] = map[string]any{} }), wantErr: true},
		{
			name: "remote state",
			files: edit(func(d map[string]any) {
				d["terraform"].(map[string]any)["backend"] = map[string]any{"s3": map[string]any{}}
			}),
			wantErr: true,
		},
		{
			name: "another provider requirement",
			files: edit(func(d map[string]any) {
				d["terraform"].(map[string]any)["required_providers"] = map[string]any{"null": map[string]any{"source": "hashicorp/null"}}
			}),
			wantErr: true,
		},
		{
			name: "provider from another source",
			files: edit(func(d map[string]any) {
				d["terraform"].(map[string]any)["required_providers"] = map[string]any{"kubernetes": map[string]any{"source": "example/kubernetes"}}
			}),
			wantErr: true,
		},
		{name: "ambient kubeconfig", files: edit(func(d map[string]any) { provider(d)["config_path"] = "~/.kube/config" }), wantErr: true},
		{name: "ambient context", files: edit(func(d map[string]any) { delete(provider(d), "config_context") }), wantErr: true},
		{name: "extra provider option", files: edit(func(d map[string]any) { provider(d)["insecure"] = true }), wantErr: true},
		{
			name: "aliased provider list",
			files: edit(func(d map[string]any) {
				d["provider"] = map[string]any{"kubernetes": []any{map[string]any{"alias": "other"}}}
			}),
			wantErr: true,
		},
		{name: "other provider", files: edit(func(d map[string]any) { d["provider"].(map[string]any)["null"] = map[string]any{} }), wantErr: true},
		{
			name: "resource outside the namespace and quota",
			files: edit(func(d map[string]any) {
				d["resource"].(map[string]any)["kubernetes_secret_v1"] = map[string]any{"x": map[string]any{}}
			}),
			wantErr: true,
		},
		{
			name: "resource bound to an aliased provider",
			files: edit(func(d map[string]any) {
				d["resource"].(map[string]any)[terraform.NamespaceResource].(map[string]any)["runtime"].(map[string]any)["provider"] = "kubernetes.other"
			}),
			wantErr: true,
		},
		{
			name: "provisioner",
			files: edit(func(d map[string]any) {
				d["resource"].(map[string]any)[terraform.NamespaceResource].(map[string]any)["runtime"].(map[string]any)["provisioner"] = []any{}
			}),
			wantErr: true,
		},
		{name: "file function", files: edit(func(d map[string]any) { namespaceMetadata(d)["name"] = `${file("/etc/hostname")}` }), wantErr: true},
		{name: "data reference", files: edit(func(d map[string]any) { namespaceMetadata(d)["name"] = "${data.external.x.result}" }), wantErr: true},
		{name: "provider function", files: edit(func(d map[string]any) { namespaceMetadata(d)["name"] = "${provider::kubernetes::manifest_encode(x)}" }), wantErr: true},
		{
			name: "duplicate keys",
			files: func(*testing.T) map[string]string {
				return map[string]string{"main.tf.json": `{"variable": {}, "variable": {}}`}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := validateTerraformRoot(writeRoot(t, tt.files(t)))

			if tt.wantErr != (err != nil) {
				t.Fatalf("validateTerraformRoot() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTerraformRoot_RefusesSymlinks(t *testing.T) {
	t.Parallel()

	directory := writeRoot(t, map[string]string{"main.tf.json": encode(t, generatedRootDocument(t))})
	elsewhere := filepath.Join(tempDir(t), "extra.tf.json")
	write(t, elsewhere, `{}`)

	if err := os.Symlink(elsewhere, filepath.Join(directory, "extra.tf.json")); err != nil {
		t.Fatal(err)
	}

	if _, err := validateTerraformRoot(directory); err == nil {
		t.Fatal("validateTerraformRoot() followed a symlinked source")
	}
}
