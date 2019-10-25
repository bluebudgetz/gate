// Code generated for package accounts by go-bindata DO NOT EDIT. (@generated)
// sources:
// get_accounts_list_query.json
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

var _get_accounts_list_queryJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x55\xc1\x8a\xdb\x30\x10\xbd\xfb\x2b\x06\xe1\xa3\xbb\xec\x5e\x73\x2b\x2d\x0b\x0b\x85\x5e\x16\x7a\x08\xa1\x4c\x64\x35\x51\x6b\x6b\x5c\x69\x0c\x0b\x25\xff\x5e\x24\xd9\xb2\x15\x3b\x4b\x16\x02\xed\xde\xa2\x37\x93\x37\x6f\x46\x4f\x9e\x6d\x01\xf0\x07\x44\x79\xb0\xd8\x1d\xbf\x10\xfd\xea\x3b\xb1\xf1\xc8\x0f\x4b\xad\xd8\x80\x40\x29\xa9\x37\xec\x44\x05\xc2\x31\x5a\xfe\xa6\xf9\xe8\x03\xe5\x77\x5d\x7b\x50\x92\x31\x4a\xf2\xa3\xa5\xf6\x51\xab\xa6\xf6\xb1\x3c\xf4\x4c\x29\xd0\xa1\x55\x86\x9f\x42\x14\x9d\x47\xe4\x51\x37\xb5\x55\x46\xc0\x09\x4e\x95\x57\x53\x00\x00\x88\xb2\xb3\xf4\x53\x49\xf6\x6a\x02\x02\x91\x76\x03\x0f\xd5\x78\x96\x56\x21\xab\xfa\xab\xc9\xd0\xbe\xab\x57\x50\x83\xad\xca\x80\x24\x25\x63\x1c\xd4\x3c\xd5\x41\x5c\x39\x9e\xef\x7c\xe9\x90\x74\x2a\x20\xd7\xd9\xa4\xa1\x8d\x24\xe3\xe8\xd8\xa2\x71\x28\x59\x93\x71\x22\x95\x68\x14\xc7\x09\x0f\x93\xfd\x74\xa1\xa4\x3f\x57\x29\x2b\xe8\x8c\x33\x8f\xe5\x63\x0f\xba\x53\x8d\x36\xbe\xb1\xed\x80\x41\xd2\x11\xe5\xb5\xc8\xf2\x38\x53\x37\xe2\xea\xa5\xb3\x0b\xd8\x07\xc8\x66\x6c\x6b\xac\x33\x96\xdf\xab\xc9\x31\xe8\xa8\xb7\x52\x7d\x4c\x0d\x54\x17\xf2\xca\xa9\xc7\x95\x8c\xdd\x02\x3b\x2d\x89\xd6\xd5\x69\x73\x43\x75\xf3\x7b\xba\x4e\x66\xf1\x5a\xc6\x3c\x3a\xfd\x1e\x7f\xed\xd2\x15\xc7\x57\x42\x3d\x1f\x48\x9b\xc3\xf3\xcb\x6d\x5d\x98\x0a\x5f\xe5\xc6\x45\xf6\xcc\x95\x43\xec\xbd\x78\x93\xd1\x1e\x14\xff\xaf\xde\x7c\xab\xba\x7f\xea\x4d\x6d\x24\xb5\xaf\x79\x13\xeb\x3a\x2c\x00\x37\xb7\xe7\xf8\xaf\xdc\x85\x25\xd3\x67\xea\xf7\x8d\x8a\x9f\xc8\xd2\xf5\xc1\xc3\xe5\x54\xe3\x0e\x5b\xdf\xb1\x48\x72\x26\xc7\x8d\x8f\xe4\x3a\xca\xe9\x49\x2d\x29\xdf\xd0\xc7\x1e\x1b\x34\x52\x9d\xd5\x74\xfd\x9e\x2d\x86\xed\xb5\xcd\x2c\x9f\xda\xae\x32\x38\x49\x2f\xce\x6f\xe4\x82\x96\x95\xed\x38\xbb\x87\x0d\xdc\x2f\xa6\x72\x06\xe7\x9b\xee\xfe\xac\x4a\xf8\x3e\xda\x61\x51\x0d\xab\xf3\xc3\x83\x5f\xd1\xc5\xae\xf8\x1b\x00\x00\xff\xff\x7d\x1e\xb0\x58\x34\x08\x00\x00")

func get_accounts_list_queryJsonBytes() ([]byte, error) {
	return bindataRead(
		_get_accounts_list_queryJson,
		"get_accounts_list_query.json",
	)
}

func get_accounts_list_queryJson() (*asset, error) {
	bytes, err := get_accounts_list_queryJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "get_accounts_list_query.json", size: 2100, mode: os.FileMode(420), modTime: time.Unix(1571579105, 0)}
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
	"get_accounts_list_query.json": get_accounts_list_queryJson,
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
	"get_accounts_list_query.json": &bintree{get_accounts_list_queryJson, map[string]*bintree{}},
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
