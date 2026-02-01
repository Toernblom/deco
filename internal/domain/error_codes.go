package domain

// ErrorCode represents a registered error code with metadata
type ErrorCode struct {
	Code     string
	Category string
	Message  string
}

// ErrorCodeRegistry maintains a registry of all error codes
type ErrorCodeRegistry struct {
	codes map[string]ErrorCode
}

// NewErrorCodeRegistry creates a new error code registry with all defined codes
func NewErrorCodeRegistry() *ErrorCodeRegistry {
	registry := &ErrorCodeRegistry{
		codes: make(map[string]ErrorCode),
	}

	// Schema errors: E001-E019
	registry.register("E001", "schema", "Node not found")
	registry.register("E002", "schema", "Invalid reference")
	registry.register("E003", "schema", "Invalid schema")
	registry.register("E004", "schema", "Circular dependency")
	registry.register("E005", "schema", "Simple error")
	registry.register("E006", "schema", "Duplicate node ID")
	registry.register("E007", "schema", "Invalid node structure")
	registry.register("E008", "schema", "Missing required field")
	registry.register("E009", "schema", "Invalid field type")
	registry.register("E010", "schema", "Unknown field")
	registry.register("E011", "schema", "Unsupported schema version")
	registry.register("E012", "schema", "Invalid metadata")
	registry.register("E013", "schema", "Reserved for future use")
	registry.register("E014", "schema", "Reserved for future use")
	registry.register("E015", "schema", "Reserved for future use")
	registry.register("E016", "schema", "Reserved for future use")
	registry.register("E017", "schema", "Reserved for future use")
	registry.register("E018", "schema", "Reserved for future use")
	registry.register("E019", "schema", "Reserved for future use")

	// Reference errors: E020-E039
	registry.register("E020", "refs", "Reference not found")
	registry.register("E021", "refs", "Dangling reference")
	registry.register("E022", "refs", "Invalid reference format")
	registry.register("E023", "refs", "Circular reference")
	registry.register("E024", "refs", "Ambiguous reference")
	registry.register("E025", "refs", "Reference type mismatch")
	registry.register("E026", "refs", "Reserved for future use")
	registry.register("E027", "refs", "Reserved for future use")
	registry.register("E028", "refs", "Reserved for future use")
	registry.register("E029", "refs", "Reserved for future use")
	registry.register("E030", "refs", "Reserved for future use")
	registry.register("E031", "refs", "Reserved for future use")
	registry.register("E032", "refs", "Reserved for future use")
	registry.register("E033", "refs", "Reserved for future use")
	registry.register("E034", "refs", "Reserved for future use")
	registry.register("E035", "refs", "Reserved for future use")
	registry.register("E036", "refs", "Reserved for future use")
	registry.register("E037", "refs", "Reserved for future use")
	registry.register("E038", "refs", "Reserved for future use")
	registry.register("E039", "refs", "Reserved for future use")

	// Validation errors: E040-E059
	registry.register("E040", "validation", "Validation failed")
	registry.register("E041", "validation", "Constraint violation")
	registry.register("E042", "validation", "CEL expression error")
	registry.register("E043", "validation", "Invalid value")
	registry.register("E044", "validation", "Value out of range")
	registry.register("E045", "validation", "Type constraint violation")
	registry.register("E046", "validation", "Reserved for future use")
	registry.register("E047", "validation", "Reserved for future use")
	registry.register("E048", "validation", "Reserved for future use")
	registry.register("E049", "validation", "Reserved for future use")
	registry.register("E050", "validation", "Reserved for future use")
	registry.register("E051", "validation", "Reserved for future use")
	registry.register("E052", "validation", "Reserved for future use")
	registry.register("E053", "validation", "Reserved for future use")
	registry.register("E054", "validation", "Reserved for future use")
	registry.register("E055", "validation", "Reserved for future use")
	registry.register("E056", "validation", "Reserved for future use")
	registry.register("E057", "validation", "Reserved for future use")
	registry.register("E058", "validation", "Reserved for future use")
	registry.register("E059", "validation", "Reserved for future use")

	// I/O errors: E060-E079
	registry.register("E060", "io", "File not found")
	registry.register("E061", "io", "Cannot read file")
	registry.register("E062", "io", "Cannot write file")
	registry.register("E063", "io", "Permission denied")
	registry.register("E064", "io", "Directory not found")
	registry.register("E065", "io", "Cannot create directory")
	registry.register("E066", "io", "YAML parse error")
	registry.register("E067", "io", "YAML format error")
	registry.register("E068", "io", "File already exists")
	registry.register("E069", "io", "Disk full")
	registry.register("E070", "io", "Reserved for future use")
	registry.register("E071", "io", "Reserved for future use")
	registry.register("E072", "io", "Reserved for future use")
	registry.register("E073", "io", "Reserved for future use")
	registry.register("E074", "io", "Reserved for future use")
	registry.register("E075", "io", "Reserved for future use")
	registry.register("E076", "io", "Reserved for future use")
	registry.register("E077", "io", "Reserved for future use")
	registry.register("E078", "io", "Reserved for future use")
	registry.register("E079", "io", "Reserved for future use")

	// Graph errors: E080-E099
	registry.register("E080", "graph", "Graph cycle detected")
	registry.register("E081", "graph", "Disconnected graph")
	registry.register("E082", "graph", "Invalid graph structure")
	registry.register("E083", "graph", "Node already in graph")
	registry.register("E084", "graph", "Edge already exists")
	registry.register("E085", "graph", "Node not in graph")
	registry.register("E086", "graph", "Reserved for future use")
	registry.register("E087", "graph", "Reserved for future use")
	registry.register("E088", "graph", "Reserved for future use")
	registry.register("E089", "graph", "Reserved for future use")
	registry.register("E090", "graph", "Reserved for future use")
	registry.register("E091", "graph", "Reserved for future use")
	registry.register("E092", "graph", "Reserved for future use")
	registry.register("E093", "graph", "Reserved for future use")
	registry.register("E094", "graph", "Reserved for future use")
	registry.register("E095", "graph", "Reserved for future use")
	registry.register("E096", "graph", "Reserved for future use")
	registry.register("E097", "graph", "Reserved for future use")
	registry.register("E098", "graph", "Reserved for future use")
	registry.register("E099", "graph", "Reserved for future use")

	return registry
}

// register adds an error code to the registry
func (r *ErrorCodeRegistry) register(code, category, message string) {
	r.codes[code] = ErrorCode{
		Code:     code,
		Category: category,
		Message:  message,
	}
}

// Lookup retrieves an error code by its code string
func (r *ErrorCodeRegistry) Lookup(code string) (ErrorCode, bool) {
	ec, exists := r.codes[code]
	return ec, exists
}

// AllCodes returns all registered error codes
func (r *ErrorCodeRegistry) AllCodes() []ErrorCode {
	codes := make([]ErrorCode, 0, len(r.codes))
	for _, code := range r.codes {
		codes = append(codes, code)
	}
	return codes
}

// ByCategory returns all error codes in a specific category
func (r *ErrorCodeRegistry) ByCategory(category string) []ErrorCode {
	var codes []ErrorCode
	for _, code := range r.codes {
		if code.Category == category {
			codes = append(codes, code)
		}
	}
	return codes
}

// Categories returns all unique categories
func (r *ErrorCodeRegistry) Categories() []string {
	seen := make(map[string]bool)
	var categories []string
	for _, code := range r.codes {
		if !seen[code.Category] {
			seen[code.Category] = true
			categories = append(categories, code.Category)
		}
	}
	return categories
}
