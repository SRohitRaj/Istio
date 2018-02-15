// Copyright 2015 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package xstrings

import (
	"bytes"
)

const _BUFFER_INIT_GROW_SIZE_MAX = 2048

// Lazy initialize a buffer.
func allocBuffer(orig, cur string) *bytes.Buffer {
	output := &bytes.Buffer{}
	maxSize := len(orig) * 4

	// Avoid to reserve too much memory at once.
	if maxSize > _BUFFER_INIT_GROW_SIZE_MAX {
		maxSize = _BUFFER_INIT_GROW_SIZE_MAX
	}

	output.Grow(maxSize)
	output.WriteString(orig[:len(orig)-len(cur)])
	return output
}
