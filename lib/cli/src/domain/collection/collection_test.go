package collection_test

import (
	"slices"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/collection"
)

const ticketText = "# 0070: Run the console with extub as a microruntime\n\n**Spec:** .agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md\n\n**What to build:** `apps/console` and `apps/extub` run together.\n\n**Blocked by:** the Exphere server, not yet ticketed.\n\n**Status:** needs-triage\n\n### Context\n\n- [What the real `ex` reproduces](/.agents/workflow/research/ex-microruntimes-capability.research.md)\n- [The console's manifest](/apps/console/package.json)\n- [The site](https://langdrift.md/platform)\n- [[notes/sizing]]\n- [Section](#context)\n"

func TestParseTicketReadsTheTrackersConvention(t *testing.T) {
	t.Parallel()

	ticket := collection.ParseTicket(".agents/workflow/tickets/0070-console-microruntime.ticket.md", []byte(ticketText))

	if ticket.Number != "0070" || ticket.Slug != "console-microruntime" || ticket.Title != "Run the console with extub as a microruntime" || ticket.Status != "needs-triage" {
		t.Fatalf("ParseTicket() = %+v", ticket)
	}

	if ticket.Name() != "0070-console-microruntime.ticket.md" || ticket.Reference() != "0070-console-microruntime" {
		t.Fatalf("Name() = %q, Reference() = %q", ticket.Name(), ticket.Reference())
	}

	plain := collection.ParseTicket("notes/idea.md", []byte("no heading here\n"))
	if plain.Number != "" || plain.Slug != "idea" || plain.Title != "" || plain.Status != "" {
		t.Fatalf("a plain file parsed as %+v", plain)
	}
}

func TestLinksFollowDocumentOrderAndSkipWhatIsNotAFile(t *testing.T) {
	t.Parallel()

	links := collection.Links(ticketText)

	want := []collection.Link{
		{Target: ".agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md", Syntax: collection.FieldLine},
		{Target: "/.agents/workflow/research/ex-microruntimes-capability.research.md", Syntax: collection.MarkdownLink},
		{Target: "/apps/console/package.json", Syntax: collection.MarkdownLink},
		{Target: "notes/sizing", Syntax: collection.WikiLink},
	}

	if !slices.Equal(links, want) {
		t.Fatalf("Links() = %+v\nwant %+v", links, want)
	}
}

func TestCandidatesReadALinkFromTheRootOrTheLinkingFile(t *testing.T) {
	t.Parallel()

	linking := ".agents/workflow/tickets/0070-console-microruntime.ticket.md"

	tests := []struct {
		link collection.Link
		want []string
	}{
		{collection.Link{Target: "/apps/console/package.json#L6", Syntax: collection.MarkdownLink}, []string{"apps/console/package.json"}},
		{collection.Link{Target: ".agents/workflow/specs/0010.spec.md", Syntax: collection.FieldLine}, []string{".agents/workflow/specs/0010.spec.md"}},
		{collection.Link{Target: "../specs/0010.spec.md", Syntax: collection.MarkdownLink}, []string{".agents/workflow/specs/0010.spec.md", "../specs/0010.spec.md"}},
		{collection.Link{Target: "notes.md", Syntax: collection.WikiLink}, []string{".agents/workflow/tickets/notes.md", "notes.md"}},
	}

	for _, tt := range tests {
		if got := collection.Candidates(tt.link, linking); !slices.Equal(got, tt.want) {
			t.Errorf("Candidates(%+v) = %v, want %v", tt.link, got, tt.want)
		}
	}

	if !collection.Escapes("../outside.md") || collection.Escapes("apps/x") {
		t.Error("Escapes() misjudges a path")
	}
}

func TestRuntimeIDDerivesFromTheSlug(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"console-microruntime":             "console-microruntime",
		"TICKET-early-access-runtime":      "early-access-runtime",
		"0031-TICKET-early-access-runtime": "early-access-runtime",
		"Fix Login_Bug":                    "fix-login-bug",
	}

	for slug, want := range tests {
		got, err := collection.RuntimeID(slug)
		if err != nil || got != want {
			t.Errorf("RuntimeID(%q) = %q, %v; want %q", slug, got, err, want)
		}
	}

	if _, err := collection.RuntimeID("TICKET-"); err == nil {
		t.Error("an empty slug named a runtime")
	}
}

func TestSecretsAndMediaTypes(t *testing.T) {
	t.Parallel()

	for _, secret := range []string{".env", "apps/web/.env.local", "certs/server.pem", "deploy/id_rsa", "ops/aws-credentials.json", ".langdrift/agent/state/inputs/envelope.json"} {
		if !collection.IsSecret(secret) {
			t.Errorf("%s is not seen as a secret", secret)
		}
	}

	for _, plain := range []string{"README.md", "apps/console/package.json", ".langdrift/harness/envelope.json", "environment.md"} {
		if collection.IsSecret(plain) {
			t.Errorf("%s is seen as a secret", plain)
		}
	}

	if collection.MediaType("a.md") != "text/markdown" || collection.MediaType("a.json") != "application/json" || collection.MediaType("go.mod") != "text/plain" {
		t.Error("MediaType() misnames a file")
	}

	if !collection.IsBinary([]byte("PNG\x00data")) || collection.IsBinary([]byte("plain text\n")) {
		t.Error("IsBinary() misjudges content")
	}
}

func TestContextHashFollowsTheRecordedRule(t *testing.T) {
	t.Parallel()

	// The sample's seven sources and its recorded context hash.
	hashes := map[string]string{
		".agents/workflow/tickets/0070-console-microruntime.ticket.md":         "6d7935ac7e4eb3b4036b5a9d6ec8fe1d22fd1d0ec05fc0a000d30305ac8faae7",
		".agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md": "3dde504fc77f063715c765dae8d51fce6535de5631068ad4b215545b8981c98a",
		".agents/workflow/research/ex-microruntimes-capability.research.md":    "7f187a177ccf87f94bae151a50b6692503b61875c98dbbb504047e55058e5e25",
		".agents/workflow/adrs/0002-console-input-is-untrusted.adr.md":         "94cd83fdbd581c5be2d2b6068f498631a074ea0d236995b0b206a9502bae90f4",
		".langdrift/harness/guardrails/0001-finops.gr.md":                      "52dec7dc62975afb1a8b46734a6297e0593a77109aff98c37a224f5c1baea446",
		"apps/console/package.json":                                            "c53acfb6db82dbb7e60a2ca80181137343109de285b96136e9b080cf978437f7",
		"apps/extub/go.mod":                                                    "aafa98e9c16d359d76ab8fce921f64232162e7808132050b02ab0da8e5fc1f3f",
	}

	if got := collection.ContextHash(hashes); got != "0d8ce74c47569bfd35cc820e302898ecde272d93715728f8de68e77ffe8f0dc7" {
		t.Fatalf("ContextHash() = %s", got)
	}
}
