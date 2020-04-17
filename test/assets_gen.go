// Code generated for package test by go-bindata DO NOT EDIT. (@generated)
// sources:
// neo4j-cleanup.cyp
// neo4j-create-nodes.cyp
// neo4j-create-tx.cyp
package test

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

	info := bindataFileInfo{name: "neo4j-cleanup.cyp", size: 27, mode: os.FileMode(420), modTime: time.Unix(1586970142, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _neo4jCreateNodesCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x93\xcf\x4a\xc3\x40\x10\xc6\xef\x3e\xc5\xde\xda\x82\x7d\x81\x1c\x84\xb6\x14\x14\x94\x40\xf1\x26\x1e\xc6\xcd\x64\x33\x34\xbb\x2b\x9b\xad\x20\xe2\xbb\x4b\xf3\x6f\xb6\x4b\x42\xa3\xf1\x98\xef\x83\xdf\x6f\x32\xcb\xec\x0e\xfb\xcd\xf3\x5e\x2c\xa5\xd5\xef\x60\x3e\x93\x8d\x94\xf6\x64\xbc\xf8\xa2\x2c\x11\x8b\x36\x5d\xdc\x0a\x03\x1a\x13\xb1\xd8\x92\x12\xbb\x3e\x94\x0e\xc1\x63\x96\x9a\x44\x64\xe0\xd1\x93\xc6\xe5\xea\x7b\x75\xd3\x41\xdf\xc0\x1c\x23\xe2\x39\xba\xc4\x6d\x9b\xe4\x1a\x0b\x48\x45\x28\x20\xc5\xa4\x07\x53\x9d\x1c\x18\x89\xbf\x1c\xaf\x45\x0e\x4c\xd9\x26\xac\xd8\x15\x28\x8f\x64\x54\x25\xb8\xba\xa6\x28\x2d\x98\x2a\x82\xd7\x19\x63\x1f\xdb\xcf\x11\xd4\xfa\x25\x91\x05\x95\x59\x9a\xbf\xae\xef\xc2\x91\x59\x22\xc1\x9d\x21\xf1\xdb\x35\x69\x30\x3f\x38\x51\x36\xc9\x24\x57\x3d\x27\x5b\x1c\x1a\xfb\x01\x9e\xac\xa9\x06\x6c\x51\xcb\xd6\x03\x17\x73\xec\xda\x3a\xaf\x40\x61\xbc\xcc\x3e\x67\xe3\x53\x10\xfd\x7d\xa9\x85\xd5\xd8\x91\x22\x67\x58\xb1\xf6\xde\x6a\x14\x9a\xe3\x49\xea\x7e\x7c\x16\xdb\x3c\x27\x39\xa6\xbe\x2c\x59\x9e\xd6\xf9\x7f\xe8\xa9\x3b\xa4\x78\xd3\x5c\x0c\x5c\xdd\xbc\x5d\x97\x94\x63\x8f\x8a\xaf\x25\xec\x82\xab\xa1\x1c\x05\x05\xf9\x24\x3b\xff\x43\xf0\xd0\x08\xa5\x2f\xc6\xf4\x51\x1b\x3c\x77\x5d\xcc\x1b\xe1\x27\x00\x00\xff\xff\xd5\xea\x04\x57\x7d\x05\x00\x00")

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

	info := bindataFileInfo{name: "neo4j-create-nodes.cyp", size: 1405, mode: os.FileMode(420), modTime: time.Unix(1587059640, 0)}
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

	info := bindataFileInfo{name: "neo4j-create-tx.cyp", size: 206, mode: os.FileMode(420), modTime: time.Unix(1587060040, 0)}
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
