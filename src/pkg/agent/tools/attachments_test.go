package tools

import "testing"

func TestStripLinmoAttachmentsMarker(t *testing.T) {
	input := "summary\n@@LIMNO_ATTACHMENTS@@\n[{\"type\":\"file\"}]"
	got := StripLinmoAttachmentsMarker(input)
	if got != "summary" {
		t.Fatalf("got %q want summary", got)
	}
	if StripLinmoAttachmentsMarker("plain text") != "plain text" {
		t.Fatal("expected plain text unchanged")
	}
}
