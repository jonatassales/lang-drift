package configuration

import "testing"

func TestImageTag(t *testing.T) {
	tag := ImageTag("book-a-demo", "demo-requests", 3)

	if tag != "langdrift-book-a-demo-demo-requests:r3" {
		t.Fatalf("ImageTag = %q", tag)
	}

	if !IsImageReference(tag) {
		t.Fatalf("ImageTag %q is not an image reference", tag)
	}
}
