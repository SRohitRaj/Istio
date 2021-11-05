// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"

	"istio.io/istio/pkg/test"
)

// AsBytes is a simple wrapper around os.ReadFile provided for completeness.
func AsBytes(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// AsBytesOrFail calls AsBytes and fails the test if any errors occurred.
func AsBytesOrFail(t test.Failer, filename string) []byte {
	t.Helper()
	content, err := AsBytes(filename)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

// AsString is a convenience wrapper around os.ReadFile that converts the content to a string.
func AsString(filename string) (string, error) {
	bytes, err := AsBytes(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// AsStringOrFail calls AsBytesOrFail and then converts to string.
func AsStringOrFail(t test.Failer, filename string) string {
	t.Helper()
	return string(AsBytesOrFail(t, filename))
}

// NormalizePath expands the homedir (~) and returns an error if the file doesn't exist.
func NormalizePath(originalPath string) (string, error) {
	if originalPath == "" {
		return "", nil
	}
	// trim leading/trailing spaces from the path and if it uses the homedir ~, expand it.
	var err error
	out := strings.TrimSpace(originalPath)
	out, err = homedir.Expand(out)
	if err != nil {
		return "", err
	}

	// Verify that the file exists.
	if _, err := os.Stat(out); os.IsNotExist(err) {
		return "", fmt.Errorf("failed normalizing file %s: %v", originalPath, err)
	}

	return out, nil
}

// ReadTarFile reads a tar compress file from the embedded
func ReadTarFile(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	tr := tar.NewReader(bytes.NewBuffer(b))
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return "", err
		}
		if hdr.Name != strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)) {
			continue
		}
		contents, err := io.ReadAll(tr)
		if err != nil {
			return "", err
		}
		return string(contents), nil
	}
	return "", fmt.Errorf("file not found %v", filePath)
}
