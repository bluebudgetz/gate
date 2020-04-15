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

var _neo4jCreateNodesCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x93\xc1\x4a\xc3\x40\x10\x86\xef\x3e\xc5\xdc\xda\x82\x7d\x81\x1c\x84\x36\x14\x14\x2a\x81\xe2\x4d\xa4\x4c\x37\x93\xcd\xd0\xec\xac\xec\x6e\x05\x11\xdf\x5d\x8c\x69\x76\x0b\x06\x52\xcd\x71\xfe\x81\xef\xfb\x33\x49\xf2\xdd\x66\xf5\xb4\x81\x79\xb6\x52\xca\x9e\x24\xc0\x07\x97\x19\xcc\x94\x35\xaf\x28\xef\xb3\x5b\x10\x34\x94\xc1\x6c\xcd\x1a\xf2\x3e\x54\x8e\x30\x50\x59\x48\x06\x25\x06\x0a\x6c\x68\xbe\xf8\x5c\xdc\x0c\xd0\x0e\x28\xc7\x4b\xd4\xfa\x27\xb9\x92\x83\xac\x23\xe6\x41\xfc\xc9\xa1\x28\xfa\x5f\xaf\x2e\x89\xdc\xbc\x26\x75\x64\xd1\x1e\xe2\xea\x4a\x6e\x63\x51\x7c\x24\x6e\xbb\x71\x80\xb2\x7c\xce\x54\xcd\x4d\x59\x54\x2f\xcb\xbb\xf9\xa1\xad\x34\x88\x56\xe8\xbe\xe9\x49\x5d\x74\xd0\x25\xa3\xf8\x6d\xb7\x41\xbc\x23\xb1\x6f\x18\xd8\x8a\xbf\xd4\xec\xe2\x62\x42\x9d\xb1\x2e\x68\xd4\x94\x1c\xeb\x31\x89\xa6\x38\x58\xaf\xd8\xd7\xd6\x50\xf4\xdc\x5b\x43\x70\x5e\x8e\x75\x19\x17\xf4\x08\x93\xad\x2a\x56\x89\xab\x68\xe7\x69\x6d\x7c\xfe\xfc\xfd\x2f\xbf\xc4\x54\xc7\x8b\x92\x7d\xc3\x55\xf2\x44\x5b\xae\x08\xfa\xed\x58\x1b\x8b\x77\x63\x5c\x35\x61\x13\xea\xe4\x5d\xb5\xf3\x5f\x7d\x5f\x01\x00\x00\xff\xff\xb7\xf7\xcc\x10\xe3\x04\x00\x00")

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

	info := bindataFileInfo{name: "neo4j-create-nodes.cyp", size: 1251, mode: os.FileMode(420), modTime: time.Unix(1586971608, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _neo4jCreateTxCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\xd2\x5d\x8b\xe2\x30\x14\x06\xe0\x7b\x7f\xc5\xb9\x10\xea\x47\xc5\xd6\x5d\x61\xed\xc5\x82\x2c\xb2\x2b\x58\x76\xd9\xdd\xbb\x61\x90\x33\x49\xda\x06\x9b\x64\x48\x52\x41\xc4\xff\x3e\xa4\xd6\xaf\x3a\xd8\x41\x06\xaf\x1a\x4e\xdf\xf0\x3e\x69\x13\x4f\xff\xff\xf8\x05\x1d\xa3\x09\x6c\x39\x8d\xa0\x6d\x34\x99\xd3\x5d\xd7\x87\x0e\x35\xb6\x9a\x51\x63\xdd\xac\x15\xcf\xfe\xfe\x9c\x95\xe1\xee\xe0\x29\xfa\x83\x9c\x56\x01\x4e\x7d\x40\xa1\x0a\x69\x23\x68\xef\x17\x3e\x70\x63\x0a\x46\x7f\x4b\xf7\xbe\x5a\xfa\x40\x94\x10\xac\x8c\x55\xab\xdd\xf3\xe0\xbb\xab\xea\xb6\x5a\xc3\xde\xa1\x81\x88\xd7\xcb\x06\xcf\x60\x8e\x7a\xb3\x1c\x05\xe1\xc4\x83\x3e\x08\x25\x6d\x06\x7d\xf0\x82\x89\x77\xaa\x0e\x83\x20\x08\xce\x8b\x29\x5a\x66\xb9\x60\x9d\xed\x86\xa1\x8e\xc0\x6d\xf7\xf7\x9b\xa3\xfd\xc3\x07\x8a\x9b\x08\x26\xee\xc4\x47\x9b\xf7\xaf\x6c\xf3\x4a\xdb\x8b\x5c\x4d\x09\x39\x1e\x9e\xa0\x5e\xc8\x1a\x8e\xa0\xce\x15\xca\x77\x74\xe1\x99\x6e\x34\xbe\xcf\x16\x5e\xda\x16\x0a\x25\x70\x69\x2c\xe6\xb9\x1b\x55\x4a\x94\xab\xa3\x51\xcb\xf5\x95\x51\xb3\x35\x93\x05\xfb\x80\xf3\xcb\xc3\x9c\x99\x12\x2c\xd6\x36\xad\x51\x85\xd2\x36\xc5\x94\x99\xa5\x0b\x34\x7d\xd5\x7b\x7f\x79\x8d\x1b\x57\xa5\x0d\x64\x95\x24\xe4\x26\x59\x25\x09\x27\x4d\xe8\xaf\x8f\x45\xe7\x3c\x61\x73\x69\x74\x0d\xcd\xa5\x29\x34\x4a\xc2\xcc\xd2\x25\x1a\xcc\xe1\xb7\xf1\x67\x90\xe7\x87\xd2\x6b\x33\xf2\xf4\x74\x35\x72\x9b\xdd\x26\x67\x0c\x73\x9b\x35\xdf\x8e\x47\xa0\x7b\xc3\xb7\x00\x00\x00\xff\xff\xd8\x10\xf2\x6d\x48\x05\x00\x00")

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

	info := bindataFileInfo{name: "neo4j-create-tx.cyp", size: 1352, mode: os.FileMode(420), modTime: time.Unix(1586971774, 0)}
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
