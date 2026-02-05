// Package resources manages static assets and build outputs.
package resources

const (
	// LibsDirectoryPath is the directory for component libraries.
	LibsDirectoryPath = "web/ui/components"
	// StylesDirectoryPath is the directory for stylesheets.
	StylesDirectoryPath = "web/ui/styles"
	// TokensDirectoryPath is the directory for design tokens.
	TokensDirectoryPath = "web/ui/styles/layers/tokens" //nolint:gosec // G101: False positive - this is a directory path
	// StaticDirectoryPath is the directory for static assets.
	StaticDirectoryPath = "web/resources/static"
)
