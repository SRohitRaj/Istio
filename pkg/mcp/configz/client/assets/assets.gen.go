// Code generated for package assets by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/config.html
package assets

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesConfigHtml = []byte(`{{ define "content" }}

    <p>
        The Mesh Configuration Protocol (MCP) client state for this process.
    </p>

    <style>
    .metadata-table-cell {
        margin: 0;
        padding: 0;
    }

    .metadata-table {
        margin: 0;
        padding: 0;
    }
    </style>

    <div>
        <table>
            <thead>
            <tr>
                <th colspan="2">Client Info</th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td>ID</td>
                <td>{{.ID}}</td>
            </tr>
            <tr>
                <td>Metadata</td>
                <td class="metadata-table-cell">
                    <table class="metadata-table">
                    {{ range $key, $value := .Metadata }}
                    <tr>
                        <td>{{$key}}</td>
                        <td>{{$value}}</td>
                    {{end}}
                    </table>
                </td>
            </tr>
            </tbody>
        </table>
    </div>

    <div>
        <table>
            <thead>
            <tr>
                <th>Supported Collections</th>
            </tr>
            </thead>
            <tbody>
            {{ range $value := .Collections }}
            <tr>
                <td>{{$value}}</td>
            {{end}}
            </tbody>
        </table>
    </div>

    <div>
        <table id="recent-requests-table">
            <thead>
            <tr>
                <th colspan="4">Recent requestsChan</th>
            </tr>
            <tr>
                <th>Time</th>
                <th>Collection</th>
                <th>Acked</th>
                <th>Nonce</th>
            </tr>
            </thead>

            <tbody>

            </tbody>
        </table>
    </div>
{{ template "last-refresh" .}}

<script>
    "use strict";

    function refreshRecentRequests() {
        var url = window.location.protocol + "//" + window.location.host + "/configj/";

        var ajax = new XMLHttpRequest();
        ajax.onload = onload;
        ajax.onerror = onerror;
        ajax.open("GET", url, true);
        ajax.send();

        function onload() {
            if (this.status == 200) { // request succeeded
                var data = JSON.parse(this.responseText);

                var table = document.getElementById("recent-requests-table");

                var tbody = document.createElement("tbody");
                for (var i = 0; i < data.LatestRequests.length; i++) {
                    var row = document.createElement("tr");

                    var c1 = document.createElement("td");
                    c1.innerText = data.LatestRequests[i].Time;
                    row.appendChild(c1);

                    var c2 = document.createElement("td");
                    c2.innerText = data.LatestRequests[i].Request.Collection;
                    row.appendChild(c2);


                    var c3 = document.createElement("td");
                    if (data.LatestRequests[i].Request.ErrorDetail === null) {
                        c3.innerText = "true"
                    } else {
                        c3.innerText = "false"
                    }
                    row.appendChild(c3);

                    var c4 = document.createElement("td");
                    c4.innerText = data.LatestRequests[i].Request.ResponseNonce;
                    row.appendChild(c4);

                    tbody.appendChild(row)
                }
                table.removeChild(table.tBodies[0]);
                table.appendChild(tbody);

                updateRefreshTime();
            }
        }

        function onerror(e) {
            console.error(e);
        }
    }

    refreshRecentRequests();
    window.setInterval(refreshRecentRequests, 1000);

</script>

{{ end }}
`)

func templatesConfigHtmlBytes() ([]byte, error) {
	return _templatesConfigHtml, nil
}

func templatesConfigHtml() (*asset, error) {
	bytes, err := templatesConfigHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/config.html", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/config.html": templatesConfigHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"config.html": &bintree{templatesConfigHtml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
