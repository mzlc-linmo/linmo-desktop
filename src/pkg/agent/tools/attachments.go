package tools

import "strings"

const openOctaAttachmentsMarker = "@@LIMNO_ATTACHMENTS@@"

// StripLinmoAttachmentsMarker removes legacy inline attachment JSON from tool output text.
func StripLinmoAttachmentsMarker(text string) string {
	idx := strings.Index(text, openOctaAttachmentsMarker)
	if idx < 0 {
		return text
	}
	return strings.TrimRight(text[:idx], " \t\n\r")
}
