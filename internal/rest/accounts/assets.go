// Code generated for package accounts by go-bindata DO NOT EDIT. (@generated)
// sources:
// get_accounts_query.json
package accounts

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

var _get_accounts_queryJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x55\xc1\x6a\xdc\x30\x10\xbd\xfb\x2b\x06\xe1\xa3\x1b\x92\xeb\xde\x4a\x4b\x20\x50\xe8\x25\xd0\xc3\xb2\x94\x89\xa4\xee\xaa\xb5\x25\x57\x1a\x43\xa0\xec\xbf\x17\x49\x96\x6c\xd9\xde\xb0\x81\x40\x9b\xdb\xea\xcd\xec\x9b\x37\xa3\x27\xcf\xbe\x02\xf8\x03\xac\x3e\x5a\xec\x4f\x5f\x8c\xf9\x35\xf4\x6c\xe7\x91\x1f\xd6\x74\x6c\x07\x0c\x39\x37\x83\x26\xc7\x1a\x60\x8e\xd0\xd2\x37\x45\x27\x1f\xa8\xbf\x2b\xe1\x41\x6e\xb4\x96\x9c\xee\xad\xe9\xee\x95\x6c\x85\x8f\x95\xa1\x47\x93\x03\x3d\x5a\xa9\xe9\x21\x44\xd1\x79\x84\x9f\x54\x2b\xac\xd4\x0c\xce\x70\x6e\xbc\x9a\x0a\x00\x80\xd5\xbd\x35\x3f\x25\x27\xaf\x26\x20\x10\x69\x77\x70\xd7\xa4\x33\xb7\x12\x49\x8a\xaf\xba\x40\x87\x5e\x6c\xa0\x1a\x3b\x59\x00\x59\x4a\xc1\x38\xaa\x79\x10\x41\x5c\x9d\xce\x37\xbe\x74\x48\x3a\x57\x50\xea\x6c\xf3\xd0\x12\x49\x1a\x1d\x59\xd4\x0e\x39\x29\xa3\x1d\xcb\x25\x5a\x49\x71\xc2\xe3\x64\x3f\x5d\x28\xe9\xcf\x4d\xce\x0a\x3a\xe3\xcc\x63\xf9\xd8\x83\xea\x65\xab\xb4\x6f\x6c\x3f\x62\x90\x75\x44\x79\x1d\x12\x3f\xcd\xd4\x25\x5c\x3e\xf7\x76\x05\xfb\x80\xb1\x05\xdb\x16\xeb\x8c\xe5\xf7\x66\x72\x0c\x3a\x33\x58\x2e\x3f\xe6\x06\x9a\x0b\x79\xf5\xd4\xe3\x46\xc6\x61\x85\x9d\xd7\x44\xdb\xea\x94\x7e\x43\x75\xf3\x7b\xba\x4e\x66\xf5\x52\xc6\x3c\x3a\xfd\x4e\xbf\x0e\xf9\x8a\xe3\x2b\x31\x03\x1d\x8d\xd2\xc7\xc7\xe7\xb7\x75\x61\x2e\x7c\x95\x1b\x57\xd9\x33\x57\x8e\xb1\xf7\xe2\x4d\x42\x7b\x94\xf4\xbf\x7a\xf3\xb5\xea\xfe\xa9\x37\x95\xe6\xa6\x7b\xc9\x9b\x28\x44\x58\x00\x6e\x6e\xcf\xf4\xaf\xd2\x85\x35\x99\xcf\x92\xab\x0e\xdb\xf8\x8d\xac\xdd\x10\x4c\x5c\x4f\x45\x6e\xb0\xf3\x2d\xb3\xac\x67\xb2\x5c\x7a\x25\x57\x72\x4e\x8f\x6a\xcd\xf9\x8a\x4e\x9e\xb0\x45\xcd\xe5\xa2\xa8\x1b\x9e\xc8\x62\xd8\x5f\xfb\xc2\xf4\xb9\xf1\xa6\x80\xb3\xf6\x6a\x79\x27\x17\xb4\x6c\xec\xc7\xd9\x4d\xec\xe0\x76\x35\x96\x05\x5c\xee\xba\xdb\x45\x95\xf0\x85\xb4\xe3\xaa\x1a\x97\xe7\x87\x3b\xbf\xa4\xab\xc3\xdf\x00\x00\x00\xff\xff\xf9\xf3\x8a\x1b\x35\x08\x00\x00")

func get_accounts_queryJsonBytes() ([]byte, error) {
	return bindataRead(
		_get_accounts_queryJson,
		"get_accounts_query.json",
	)
}

func get_accounts_queryJson() (*asset, error) {
	bytes, err := get_accounts_queryJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "get_accounts_query.json", size: 2101, mode: os.FileMode(420), modTime: time.Unix(1571332004, 0)}
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
	"get_accounts_query.json": get_accounts_queryJson,
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
	"get_accounts_query.json": &bintree{get_accounts_queryJson, map[string]*bintree{}},
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
