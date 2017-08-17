package gomagtek

import "github.com/google/gousb"

// ============================================================================
// Context Type.
// ============================================================================

/*
 * Context manages all resources related to USB device handling.
 */
type Context struct {
	*gousb.Context
}

/*
 * NewContext returns a new Context instance.
 */
func NewContext(d *gousb.Context) (*Context) {
	return &Context{}
}
