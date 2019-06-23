// Package schema Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// v1_accounts_POST.json5
// v1_accounts_id_PATCH.json5
// v1_accounts_id_PUT.json5
// v1_transactions_id_PATCH.json5
// v1_transactions_id_POST.json5
// v1_transactions_id_PUT.json5
package schema

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

var _v1_accounts_postJson5 = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8f\xb1\x4e\x04\x21\x10\x86\x7b\x9e\x62\x32\x5a\xde\x89\x9d\x09\x9d\xa5\x95\x0f\x60\x2c\x90\x9d\xbb\x63\xb3\x0b\x38\xcc\x16\xc6\xec\xbb\x1b\x16\xdc\x90\x4b\xb6\x63\x3e\xfe\x6f\x7e\xf8\x55\x00\xf8\x98\xdd\x8d\x66\x8b\x06\xf0\x26\x92\x8c\xd6\x63\x8e\xe1\x5c\xe9\x53\xe4\xab\x1e\xd8\x5e\xe4\xfc\xfc\xa2\x2b\x7b\xc0\x53\xf1\xc4\xcb\x44\xc5\x7a\x75\x2e\x2e\x41\x1a\xfd\x49\x1b\x8c\x5f\x23\xb9\xc6\x12\xc7\x44\x2c\x9e\x32\x1a\x28\x9d\x00\x18\xec\x4c\xfb\xd4\x79\x59\xd8\x87\x2b\x6e\x78\x3d\xd5\x6c\xb2\x4c\x41\xde\x86\x3e\x1f\x03\xbd\x5f\xd0\xc0\x47\x03\xb0\x5f\x1d\xae\xeb\x56\x1e\xe6\xc3\x32\x4d\x5d\xba\x9d\x3e\xd5\xff\xb4\xf9\xc8\xf4\xbd\x78\xa6\x61\xef\xaf\xdf\xb9\x7f\xae\x2a\xe6\xaa\xfe\x02\x00\x00\xff\xff\x8f\x36\xcd\xfb\x66\x01\x00\x00")

func v1_accounts_postJson5Bytes() ([]byte, error) {
	return bindataRead(
		_v1_accounts_postJson5,
		"v1_accounts_POST.json5",
	)
}

func v1_accounts_postJson5() (*asset, error) {
	bytes, err := v1_accounts_postJson5Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v1_accounts_POST.json5", size: 358, mode: os.FileMode(420), modTime: time.Unix(1560173156, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v1_accounts_id_patchJson5 = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8d\x3d\xae\xc3\x20\x10\x84\x7b\x4e\xb1\xda\xf7\x4a\x3b\xa4\x8b\x44\x97\x32\x55\x0e\x10\xa5\x20\x78\xfd\x27\x1b\x10\xac\x8b\x28\xf2\xdd\x23\xdb\xc4\xa2\x71\xc7\x7c\x7c\x33\xfb\x11\x00\xf8\x1f\x4d\x4b\xa3\x46\x05\xd8\x32\x7b\x25\x65\x1f\x9d\x2d\x37\x7a\x72\xa1\x91\x55\xd0\x35\x97\xe7\x8b\xdc\xd8\x1f\x16\x4b\x8f\x3b\x1e\x68\x69\x5d\x8d\x71\x93\xe5\x44\xdf\x7e\x85\xee\xd5\x93\x49\xcc\x07\xe7\x29\x70\x47\x11\x15\x2c\x37\x01\xd0\xea\x91\xf6\x94\xf5\x22\x87\xce\x36\xb8\xe2\xb9\xd8\x5c\xaf\x03\x59\xbe\x55\xb9\xef\x2c\xdd\x6b\x54\xf0\x48\x00\xf6\xaf\xc3\xb9\x6c\xf2\xd0\xb7\xd3\x30\x64\x76\x7a\x3d\xc5\x2f\xcd\x62\x16\xdf\x00\x00\x00\xff\xff\xf8\xf6\x6b\xa9\x36\x01\x00\x00")

func v1_accounts_id_patchJson5Bytes() ([]byte, error) {
	return bindataRead(
		_v1_accounts_id_patchJson5,
		"v1_accounts_id_PATCH.json5",
	)
}

func v1_accounts_id_patchJson5() (*asset, error) {
	bytes, err := v1_accounts_id_patchJson5Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v1_accounts_id_PATCH.json5", size: 310, mode: os.FileMode(420), modTime: time.Unix(1561236900, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v1_accounts_id_putJson5 = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8f\xb1\x4e\x04\x21\x10\x86\x7b\x9e\x62\x32\x5a\xde\x89\x9d\x09\x9d\xa5\x95\x0f\x60\x2c\x90\x9d\xbb\x63\xb3\x0b\x38\xcc\x16\xc6\xec\xbb\x1b\x16\xdc\x90\x4b\xb6\x63\x3e\xfe\x6f\x7e\xf8\x55\x00\xf8\x98\xdd\x8d\x66\x8b\x06\xf0\x26\x92\x8c\xd6\x63\x8e\xe1\x5c\xe9\x53\xe4\xab\x1e\xd8\x5e\xe4\xfc\xfc\xa2\x2b\x7b\xc0\x53\xf1\xc4\xcb\x44\xc5\x7a\x75\x2e\x2e\x41\x1a\xfd\x49\x1b\x8c\x5f\x23\xb9\xc6\x12\xc7\x44\x2c\x9e\x32\x1a\x28\x9d\x00\x18\xec\x4c\xfb\xd4\x79\x59\xd8\x87\x2b\x6e\x78\x3d\xd5\x6c\xb2\x4c\x41\xde\x86\x3e\x1f\x03\xbd\x5f\xd0\xc0\x47\x03\xb0\x5f\x1d\xae\xeb\x56\x1e\xe6\xc3\x32\x4d\x5d\xba\x9d\x3e\xd5\xff\xb4\xf9\xc8\xf4\xbd\x78\xa6\x61\xef\xaf\xdf\xb9\x7f\xae\x2a\xe6\xaa\xfe\x02\x00\x00\xff\xff\x8f\x36\xcd\xfb\x66\x01\x00\x00")

func v1_accounts_id_putJson5Bytes() ([]byte, error) {
	return bindataRead(
		_v1_accounts_id_putJson5,
		"v1_accounts_id_PUT.json5",
	)
}

func v1_accounts_id_putJson5() (*asset, error) {
	bytes, err := v1_accounts_id_putJson5Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v1_accounts_id_PUT.json5", size: 358, mode: os.FileMode(420), modTime: time.Unix(1560173124, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v1_transactions_id_patchJson5 = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x8d\x31\x4f\xc3\x40\x0c\x46\xf7\xfc\x0a\xcb\x30\xb6\x84\x0d\x29\x1b\x23\x53\x17\x36\xc4\x70\xbd\xb8\xe9\x55\x39\x3b\xb2\x7d\x03\x42\xf9\xef\x28\x4d\xa8\x82\x44\xa5\x76\x4b\xde\xf7\xde\xf9\xbb\x02\xc0\x47\x8b\x47\xca\x01\x1b\xc0\xa3\xfb\xd0\xd4\xf5\xc9\x84\xb7\x33\x7d\x12\xed\xea\x56\xc3\xc1\xb7\xcf\x2f\xf5\xcc\x1e\x70\x33\x75\x9e\xbc\xa7\xa9\x7a\xd7\xc0\x16\xa2\x27\xe1\x65\xf9\x1a\xce\x83\xec\x4f\x14\x7d\x66\x83\xca\x40\xea\x89\x0c\x1b\x98\xee\x02\x60\x32\x2b\xd4\xee\xf8\x42\x56\xad\xb9\x26\xee\xf0\x8c\xc7\xcd\xec\x8b\xa6\x2e\xdd\x6c\x9b\x14\x8d\xf4\x1a\xa3\x14\xf6\xb7\xf6\xd6\xcc\x83\x76\xe4\x77\x67\x21\x4f\xfe\x7f\x36\x97\xbc\x27\xfd\x6b\x47\xc9\x99\xd8\x6d\xed\x0b\xd3\xee\x80\x0d\x7c\x2c\x00\x2e\xd3\xd5\xe3\xab\x27\xaf\xfa\x5c\xfa\x7e\x65\x2f\x5f\x9f\xd5\xef\xdf\x58\x8d\xd5\x4f\x00\x00\x00\xff\xff\x8d\x0a\x75\x3c\x08\x02\x00\x00")

func v1_transactions_id_patchJson5Bytes() ([]byte, error) {
	return bindataRead(
		_v1_transactions_id_patchJson5,
		"v1_transactions_id_PATCH.json5",
	)
}

func v1_transactions_id_patchJson5() (*asset, error) {
	bytes, err := v1_transactions_id_patchJson5Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v1_transactions_id_PATCH.json5", size: 520, mode: os.FileMode(420), modTime: time.Unix(1561317397, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v1_transactions_id_postJson5 = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x91\x31\x4f\xc3\x30\x10\x85\xf7\xfc\x0a\xeb\x60\x6c\x09\x1b\x52\x36\x46\xa6\x2e\x6c\x88\xc1\x75\xae\xa9\xab\xe4\x2e\xdc\x9d\x07\x84\xfa\xdf\x91\x93\x10\x12\x89\x48\x65\xb3\xdf\xfb\x9e\x9f\xf4\xfc\x55\x38\x07\xf7\x1a\xce\xd8\x79\xa8\x1c\x9c\xcd\xfa\xaa\x2c\x2f\xca\xb4\x1f\xd5\x07\x96\xa6\xac\xc5\x9f\x6c\xff\xf8\x54\x8e\xda\x1d\xec\x72\xce\xa2\xb5\x98\x53\xaf\xe2\x49\x7d\xb0\xc8\x34\x39\x9f\xfd\x60\xf0\xf1\x82\xc1\x46\xad\x17\xee\x51\x2c\xa2\x42\xe5\x72\xaf\x73\x10\x55\x13\xd6\x07\x9a\x95\x45\x56\x4d\x22\x35\x30\xc8\xd7\xdd\xc8\xb3\xc4\x26\xde\x4c\x2b\x27\x09\xf8\x1c\x02\x27\xb2\x97\xfa\xd6\x98\x79\x69\xd0\xfe\x1d\xf3\x5d\xe6\xff\xa2\x29\x75\x47\x94\x35\x1d\xb8\xeb\x90\x4c\x97\x3c\x13\x1e\x4e\x50\xb9\xb7\x49\x70\xb3\xb5\x59\xbe\x78\x72\x93\xa7\xd4\xb6\x0b\x7a\x3a\xbd\x17\x3f\xb7\x21\x0f\x82\x1f\x29\x0a\xd6\x73\xff\xef\xe7\xac\xc7\xdf\x18\x77\x63\xbc\xf5\x38\x45\xee\xbd\x16\xdf\x01\x00\x00\xff\xff\xcd\xb3\xdc\x00\x76\x02\x00\x00")

func v1_transactions_id_postJson5Bytes() ([]byte, error) {
	return bindataRead(
		_v1_transactions_id_postJson5,
		"v1_transactions_id_POST.json5",
	)
}

func v1_transactions_id_postJson5() (*asset, error) {
	bytes, err := v1_transactions_id_postJson5Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v1_transactions_id_POST.json5", size: 630, mode: os.FileMode(420), modTime: time.Unix(1561318157, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v1_transactions_id_putJson5 = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x91\x31\x4f\xc3\x30\x10\x85\xf7\xfc\x0a\xcb\x30\xb6\x84\x0d\x29\x1b\x23\x53\x17\x36\xc4\xe0\x3a\xd7\xd4\x55\x7c\x17\xee\xce\x03\x42\xfd\xef\xc8\x49\x08\x09\x22\x52\xba\xd9\xef\xbe\xe7\x67\x3f\x7f\x15\xc6\xd8\x7b\xf1\x67\x88\xce\x56\xc6\x9e\x55\xbb\xaa\x2c\x2f\x42\xb8\x1f\xd4\x07\xe2\xa6\xac\xd9\x9d\x74\xff\xf8\x54\x0e\xda\x9d\xdd\x65\x9f\x06\x6d\x21\xbb\x5e\xd9\xa1\x38\xaf\x81\x70\x9c\x7c\x76\xfd\x80\x8e\x17\xf0\x3a\x68\x1d\x53\x07\xac\x01\xc4\x56\x26\xe7\x1a\x63\x83\x48\x82\xfa\x80\x93\x32\xf3\x8a\x72\xc0\xc6\xf6\xf2\x75\x37\xf0\xc4\xa1\x09\x9b\x69\xa1\xc4\x1e\x9e\xbd\xa7\x84\xfa\x52\x6f\xb5\xa9\xe3\x06\xf4\x66\x9b\x8b\x99\xff\x8f\xc6\x14\x8f\xc0\x4b\xda\x53\x8c\x80\x2a\x73\x9e\x10\x0e\x27\x5b\x99\xb7\x51\x30\xd3\x68\x35\x7c\x76\xe4\x2a\x8f\xa9\x6d\x67\xf4\xb8\x7a\x2f\x7e\x76\xbd\xdf\x32\x7c\xa4\xc0\x50\x4f\xf9\xbf\x9f\xb3\x2c\x7f\xa5\xdc\x95\xf2\x96\xe5\xfc\x7d\x7c\x91\xef\x71\x2d\xbe\x03\x00\x00\xff\xff\x36\x1b\xd7\x90\x86\x02\x00\x00")

func v1_transactions_id_putJson5Bytes() ([]byte, error) {
	return bindataRead(
		_v1_transactions_id_putJson5,
		"v1_transactions_id_PUT.json5",
	)
}

func v1_transactions_id_putJson5() (*asset, error) {
	bytes, err := v1_transactions_id_putJson5Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v1_transactions_id_PUT.json5", size: 646, mode: os.FileMode(420), modTime: time.Unix(1561318327, 0)}
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
	"v1_accounts_POST.json5":         v1_accounts_postJson5,
	"v1_accounts_id_PATCH.json5":     v1_accounts_id_patchJson5,
	"v1_accounts_id_PUT.json5":       v1_accounts_id_putJson5,
	"v1_transactions_id_PATCH.json5": v1_transactions_id_patchJson5,
	"v1_transactions_id_POST.json5":  v1_transactions_id_postJson5,
	"v1_transactions_id_PUT.json5":   v1_transactions_id_putJson5,
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
	"v1_accounts_POST.json5":         &bintree{v1_accounts_postJson5, map[string]*bintree{}},
	"v1_accounts_id_PATCH.json5":     &bintree{v1_accounts_id_patchJson5, map[string]*bintree{}},
	"v1_accounts_id_PUT.json5":       &bintree{v1_accounts_id_putJson5, map[string]*bintree{}},
	"v1_transactions_id_PATCH.json5": &bintree{v1_transactions_id_patchJson5, map[string]*bintree{}},
	"v1_transactions_id_POST.json5":  &bintree{v1_transactions_id_postJson5, map[string]*bintree{}},
	"v1_transactions_id_PUT.json5":   &bintree{v1_transactions_id_putJson5, map[string]*bintree{}},
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
