// Code generated for package services by go-bindata DO NOT EDIT. (@generated)
// sources:
// neo4j-cleanup.cyp
// neo4j-create-nodes.cyp
// neo4j-create-tx.cyp
package services

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

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

var _neo4jCleanupCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\x75\x0c\x71\xf6\x50\xd0\xc8\xd3\x54\x70\x71\x0d\x71\x74\xf6\x50\x70\x71\xf5\x71\x0d\x71\x05\x89\x00\x02\x00\x00\xff\xff\x3d\xed\xc8\xcd\x1b\x00\x00\x00")

func neo4jCleanupCypBytes() ([]byte, error) {
	return bindataRead(
		_neo4jCleanupCyp,
		"neo4j-cleanup.cyp",
	)
}

func neo4jCleanupCyp() (*asset, error) {
	bytes, err := neo4jCleanupCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "neo4j-cleanup.cyp", size: 27, mode: os.FileMode(420), modTime: time.Unix(1587165408, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _neo4jCreateNodesCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x95\x4d\x8b\xdb\x3c\x14\x85\xf7\xef\xaf\x30\xb3\x99\x19\x78\x25\xce\xf5\x57\x12\x2d\x0a\x33\xc3\x40\x0b\x2d\x81\x41\xab\x96\x2e\x54\x47\x71\xc4\xc4\x72\x71\x3c\x85\x52\xfa\xdf\x4b\x6e\x3e\xe4\xb8\x76\xe3\x65\xce\x81\x3c\xcf\xe5\x4a\xf2\xd3\xcb\xf3\x83\x7e\x8e\xee\x8a\xba\xfa\x6e\xfc\x4f\xf5\x50\x14\xf5\x9b\x6f\xa3\x5f\x6e\xa5\xa2\xdb\x63\x7a\xfb\x7f\xe4\x4d\x65\x55\x74\xfb\xe8\xca\xe8\xe9\x1c\x16\x8d\x35\xad\x5d\x2d\xbd\x8a\x56\xa6\xb5\xad\xab\xec\xdd\x4d\x0c\xcc\x04\x72\x81\x4c\x23\x55\x48\x14\x62\x09\xe0\xf3\xcd\xfd\xef\xfb\xff\x4e\xb8\x6f\xc6\xbf\xf6\x58\xfb\xe8\x12\xf4\x78\x48\xc6\x28\x73\xc1\x20\x8d\x4c\x31\x48\x02\xf1\x25\xc5\xb8\xb2\x07\x31\xae\x0c\x8c\x0f\x7e\xf7\xd6\x18\x5f\xd8\x09\x23\x2d\x04\xf3\x34\x72\xc5\x3c\x09\x24\x7f\x8f\x74\x84\x0d\x4c\x76\x4c\x02\xfc\x69\x63\x8b\x57\xe7\xcb\x5d\x14\xaa\x11\x38\x41\x30\x5f\x63\xa6\x98\x2f\x81\xf4\x12\xbe\xad\x8d\xdf\xf5\xb0\x9c\x05\xe0\xc7\xe3\xcf\x31\x08\x09\xe6\x68\xcc\x15\x73\x24\x90\x31\x44\x7c\x51\xc5\xc6\x6d\x57\xcb\xf5\x57\xf1\xae\x3b\x66\xc0\x17\xa6\xd9\xff\x7d\xff\xf4\x1c\xd2\xce\xcc\xa6\x89\xb6\x87\x64\xcc\x22\x16\x2c\xa2\xb1\x50\x2c\x22\x81\x7c\xc0\x82\x67\x0b\xfc\xc6\xfa\xfa\x87\x69\x5d\xed\x77\x03\x1e\xbd\x36\xf8\xbc\x84\xe2\x9a\x57\x22\x58\x4d\x13\x14\xab\x49\x60\x76\xdd\xab\xaa\x9b\xb6\x34\xa5\xed\xaf\xe6\x9c\x07\x97\x4f\x9d\x68\x4c\x22\x15\x20\x41\xb1\x26\x52\xec\x21\x81\xf9\xd4\x15\x6d\xea\xca\x9e\x18\x3d\x9b\x6e\x15\x84\xde\xd7\x95\x8d\xaa\x10\x8f\x49\x65\x02\xb1\xa0\x44\x53\xac\xd8\x4b\x02\x8b\x01\xa9\xf3\xc8\x41\xa9\x5e\xaf\x5d\x31\x26\x75\x59\x06\xad\x25\xe7\x53\xc4\x72\x81\x44\x50\xaa\x29\x51\xec\x26\x41\x98\x26\xe6\x4e\x0f\x43\x7f\x6f\xa1\x18\x78\x45\xfe\xb5\xb9\x99\x40\x2a\x28\xd3\x94\x2a\xd6\x91\x20\x9a\xba\xb9\xad\x5b\xdb\x33\xa4\x7f\xc7\xbb\x5d\xe7\xae\xbb\xb5\x8d\x5c\x27\x1f\xf3\x9a\x0b\x64\x82\x72\x4d\x99\x62\x35\x09\x8a\x07\xbc\xc2\xdc\x9d\x03\x65\xcd\xb6\xdd\x8c\x89\xf5\xda\xce\xb1\xe2\x62\x92\xdc\x62\xff\x19\xa1\x99\xa6\x5c\xb1\x9f\x04\x25\x57\xe4\xfe\x04\x00\x00\xff\xff\x82\x99\x9c\xa4\xcf\x06\x00\x00")

func neo4jCreateNodesCypBytes() ([]byte, error) {
	return bindataRead(
		_neo4jCreateNodesCyp,
		"neo4j-create-nodes.cyp",
	)
}

func neo4jCreateNodesCyp() (*asset, error) {
	bytes, err := neo4jCreateNodesCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "neo4j-create-nodes.cyp", size: 1743, mode: os.FileMode(420), modTime: time.Unix(1587492025, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _neo4jCreateTxCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\x8e\x31\xaa\xc3\x30\x0c\x40\x77\x9f\x42\x83\x21\x31\x24\x17\xc8\xf0\xe1\x53\x42\xdb\x21\xb4\x94\x6e\xa5\x83\x89\x4c\xd1\xe0\x84\x5a\xca\x14\x72\xf7\x22\xdb\x9a\x1e\x4f\x0f\xa1\xe9\xff\x79\xba\x40\xcb\x69\x86\x9d\x70\x00\xcb\x69\xbe\xe2\xe1\x3a\x68\x91\xa5\x3a\x64\x51\x67\xa6\xf1\x71\x1e\x73\xec\xfa\xd7\x70\xf7\x84\xb0\x1b\x00\x6d\xea\x58\xc2\xce\x00\xf8\xb8\x6e\x8b\x64\x6b\x0b\xab\x25\xe6\x2d\xe0\x6d\x19\xb4\xab\xac\x7e\x5e\x63\x0c\x25\xb7\x95\xb3\x4e\xc1\x4b\xe9\xd1\x4b\x10\x8a\xa1\x75\xba\x58\x13\x7d\x68\xc9\xd7\x9b\xaf\x6f\x0e\xf3\xee\xff\xf4\x5b\x67\x7e\x01\x00\x00\xff\xff\x27\xd9\x89\x6a\xce\x00\x00\x00")

func neo4jCreateTxCypBytes() ([]byte, error) {
	return bindataRead(
		_neo4jCreateTxCyp,
		"neo4j-create-tx.cyp",
	)
}

func neo4jCreateTxCyp() (*asset, error) {
	bytes, err := neo4jCreateTxCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "neo4j-create-tx.cyp", size: 206, mode: os.FileMode(420), modTime: time.Unix(1587165408, 0)}
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
	"neo4j-cleanup.cyp":      neo4jCleanupCyp,
	"neo4j-create-nodes.cyp": neo4jCreateNodesCyp,
	"neo4j-create-tx.cyp":    neo4jCreateTxCyp,
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
	"neo4j-cleanup.cyp":      &bintree{neo4jCleanupCyp, map[string]*bintree{}},
	"neo4j-create-nodes.cyp": &bintree{neo4jCreateNodesCyp, map[string]*bintree{}},
	"neo4j-create-tx.cyp":    &bintree{neo4jCreateTxCyp, map[string]*bintree{}},
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
