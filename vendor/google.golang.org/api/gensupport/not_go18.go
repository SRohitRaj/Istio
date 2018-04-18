// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.8

package gensupport

import (
	"io"
	"net/http"
)

func SetGetBody(req *http.Request, f func() (io.ReadCloser, error)) {}
