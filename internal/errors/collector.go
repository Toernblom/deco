package errors

import (
	"fmt"
	"sort"

	"github.com/Toernblom/deco/internal/domain"
)

// Collector aggregates multiple errors during validation.
// It handles deduplication, sorting, and limiting the number of errors returned.
type Collector struct {
	errors     []domain.DecoError
	maxErrors  int             // 0 means no limit
	totalCount int             // Total errors added (including duplicates and truncated)
	seen       map[string]bool // For deduplication
}

// NewCollector creates a new error collector with no limit.
func NewCollector() *Collector {
	return &Collector{
		errors:    make([]domain.DecoError, 0),
		maxErrors: 0, // No limit
		seen:      make(map[string]bool),
	}
}

// NewCollectorWithLimit creates a new error collector with a maximum error limit.
// Once the limit is reached, additional errors are counted but not stored.
func NewCollectorWithLimit(maxErrors int) *Collector {
	return &Collector{
		errors:    make([]domain.DecoError, 0),
		maxErrors: maxErrors,
		seen:      make(map[string]bool),
	}
}

// Add adds an error to the collector.
// Duplicate errors (same code + location) are automatically deduplicated.
func (c *Collector) Add(err domain.DecoError) {
	c.totalCount++

	// Check for duplicate (same code + location)
	key := c.deduplicationKey(err)
	if c.seen[key] {
		return // Skip duplicate
	}
	c.seen[key] = true

	// Add error only if we haven't hit the limit
	if c.maxErrors == 0 || len(c.errors) < c.maxErrors {
		c.errors = append(c.errors, err)
	}
}

// AddBatch adds multiple errors at once.
func (c *Collector) AddBatch(errs []domain.DecoError) {
	for _, err := range errs {
		c.Add(err)
	}
}

// Errors returns all collected errors, sorted by location.
// Sorting order: file (alphabetically), line (numerically), column (numerically).
// Errors without location come last.
func (c *Collector) Errors() []domain.DecoError {
	// Make a copy to avoid mutating internal state
	result := make([]domain.DecoError, len(c.errors))
	copy(result, c.errors)

	// Sort by file, then line, then column
	sort.Slice(result, func(i, j int) bool {
		// Errors with location come before errors without location
		if result[i].Location == nil && result[j].Location != nil {
			return false
		}
		if result[i].Location != nil && result[j].Location == nil {
			return true
		}

		// Both have no location - maintain stable order
		if result[i].Location == nil && result[j].Location == nil {
			return false
		}

		// Both have locations - compare them
		locI := result[i].Location
		locJ := result[j].Location

		// Compare files
		if locI.File != locJ.File {
			return locI.File < locJ.File
		}

		// Same file, compare lines
		if locI.Line != locJ.Line {
			return locI.Line < locJ.Line
		}

		// Same line, compare columns
		return locI.Column < locJ.Column
	})

	return result
}

// HasErrors returns true if any errors have been collected.
func (c *Collector) HasErrors() bool {
	return len(c.errors) > 0
}

// Count returns the total number of unique errors added to the collector.
// This may be higher than len(Errors()) if deduplication or truncation occurred.
func (c *Collector) Count() int {
	return len(c.seen)
}

// Truncated returns true if the collector hit its error limit.
// This means some errors were counted but not stored.
func (c *Collector) Truncated() bool {
	if c.maxErrors == 0 {
		return false
	}
	return len(c.seen) > c.maxErrors
}

// Reset clears all collected errors.
func (c *Collector) Reset() {
	c.errors = make([]domain.DecoError, 0)
	c.totalCount = 0
	c.seen = make(map[string]bool)
}

// deduplicationKey generates a unique key for an error based on code and location.
// Errors with the same code and location are considered duplicates.
func (c *Collector) deduplicationKey(err domain.DecoError) string {
	if err.Location == nil {
		// No location - use code + summary as key
		return fmt.Sprintf("%s:%s", err.Code, err.Summary)
	}

	// Use code + file + line + column as key
	// Different columns on the same line are different errors
	loc := err.Location
	if loc.Column > 0 {
		return fmt.Sprintf("%s:%s:%d:%d", err.Code, loc.File, loc.Line, loc.Column)
	}
	return fmt.Sprintf("%s:%s:%d", err.Code, loc.File, loc.Line)
}
