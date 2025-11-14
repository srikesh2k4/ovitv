package utils

import "github.com/microcosm-cc/bluemonday"

var sanitizer = bluemonday.UGCPolicy()

func Sanitize(input string) string {
    return sanitizer.Sanitize(input)
}