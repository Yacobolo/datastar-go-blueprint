// Package ds provides small helpers for composing Datastar request expressions.
package ds

import (
	"fmt"
	"strings"
)

// Opt formats a Datastar request option as a JavaScript object pair.
func Opt(key, value string) string {
	return fmt.Sprintf("%s: %q", key, value)
}

// Get returns a Datastar @get expression.
func Get(url string, opts ...string) string {
	return request("get", url, opts...)
}

// Post returns a Datastar @post expression.
func Post(url string, opts ...string) string {
	return request("post", url, opts...)
}

// Put returns a Datastar @put expression.
func Put(url string, opts ...string) string {
	return request("put", url, opts...)
}

// Delete returns a Datastar @delete expression.
func Delete(url string, opts ...string) string {
	return request("delete", url, opts...)
}

// Getf returns a formatted Datastar @get expression.
func Getf(urlFormat string, args ...any) string {
	return request("get", fmt.Sprintf(urlFormat, args...))
}

// Postf returns a formatted Datastar @post expression.
func Postf(urlFormat string, args ...any) string {
	return request("post", fmt.Sprintf(urlFormat, args...))
}

// Putf returns a formatted Datastar @put expression.
func Putf(urlFormat string, args ...any) string {
	return request("put", fmt.Sprintf(urlFormat, args...))
}

// Deletef returns a formatted Datastar @delete expression.
func Deletef(urlFormat string, args ...any) string {
	return request("delete", fmt.Sprintf(urlFormat, args...))
}

func request(method, url string, opts ...string) string {
	if len(opts) == 0 {
		return fmt.Sprintf("@%s(%q)", method, url)
	}

	return fmt.Sprintf("@%s(%q,{%s})", method, url, strings.Join(opts, ", "))
}
