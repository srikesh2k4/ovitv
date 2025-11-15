package utils

import "github.com/microcosm-cc/bluemonday"

var sanitizer = bluemonday.UGCPolicy()

// Sanitize removes dangerous HTML/script content from input.
func Sanitize(input string) string {
    return sanitizer.Sanitize(input)
}
