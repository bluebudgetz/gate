// Package assets Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// assets/rdbms/init.sql
// deployments/rdbms/migrations/01_create_tables.down.sql
// deployments/rdbms/migrations/01_create_tables.up.sql
// deployments/rdbms/schema.sql
package assets

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

var _assetsRdbmsInitSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x93\x5f\x6b\xdb\x30\x14\xc5\xdf\xfd\x29\x0e\x7d\xb1\x5d\x44\x5a\xe7\x5f\x19\x7d\xf2\x96\x8c\x1a\x8c\x3d\xf2\xa7\x63\x4f\xe1\xc6\xba\x4d\x44\x1d\x29\xc8\x72\x20\xdf\x7e\x28\x59\x3b\xaf\x7b\x1b\x9e\x1f\xcc\xd1\x85\xfb\x3b\xe8\x88\x73\x77\x1b\xe0\x16\xab\xbd\x6a\xf0\xa2\x6a\x86\x6a\x50\x51\x5d\xb3\xc4\xf6\x8c\x9b\x1d\x39\xbe\x81\x33\x20\x29\xa1\xb4\x72\x8a\x6a\x48\x72\xe4\x67\x6e\xcf\x17\xbd\xa5\x86\x07\x01\x2e\xa0\xcc\x21\x5b\xa2\x2c\xf2\x1f\x50\xfa\x64\x5e\xd9\xaf\x41\xf2\x89\x6b\x73\x3c\xb0\x76\x60\x7d\x52\xd6\x68\xaf\xfd\xd2\x5d\x90\x15\xcb\xf9\x62\x85\xac\x58\x95\xd8\x6e\x07\x54\x55\xa6\xd5\xae\x41\xa4\xa4\x40\x65\x99\x1c\xcb\x8d\xd1\x02\xed\x51\xbe\x6b\xc9\x35\xbf\x69\x4d\x07\x16\x38\x92\x65\xed\x36\x4a\xc6\xc1\x73\x9a\xaf\xe7\x4b\x44\x89\xc0\x2c\x5d\xcd\x37\xe9\x6c\x16\x79\x11\x15\xe5\xf7\x28\x8e\x85\xf7\x9a\x2f\x9e\xd3\x1c\x5f\xf3\xb2\x5c\x44\x8b\xb4\x98\x45\x31\x6e\x31\x1c\xc7\x78\x2a\xd7\x8b\x58\xa0\x58\xe7\xf9\xdb\x3f\xfc\x4c\xfa\x35\x09\xaf\xc7\x58\x04\xb8\x7e\xd1\xb0\x4f\xfe\xf0\x6f\xfe\xa8\x2f\x7e\x6e\x48\x37\xa1\x40\xd2\x81\x8f\xfb\x82\x2f\xa9\x26\x7b\xfe\x40\x9f\xf4\x45\xcf\x74\xd3\x5a\xd2\x15\x87\x02\xc3\x8e\xc1\xb4\x2f\x83\x2f\x7b\xa6\x63\xaa\xe5\x13\x53\xed\xf6\x67\x64\xba\x1a\x84\x02\x93\x8e\xd7\x43\x7f\x51\xbd\xf0\xd2\xd9\xb6\x72\xc8\x9d\xbc\xda\x3c\x06\x1f\x1b\xe0\x2c\xe9\x86\x2a\xa7\x8c\x6e\x10\x75\x1b\x60\xac\xda\x29\x2d\xd0\x98\xd6\x56\xbc\xf9\x55\x95\x8d\x2f\x8a\x23\xbb\x63\xf7\xc7\x88\x0e\x5e\x0a\x54\xe6\xe0\xdb\xd6\xfc\x2e\xc6\x3f\xdf\x26\xcc\xb4\x72\xa1\xc0\x58\x60\x2a\x90\xdc\xdf\x0b\x84\xd7\xe0\xa0\xde\x1f\xaa\x93\x5c\x1f\x46\x0f\x02\xa3\x89\x0f\xef\x12\x5c\x6b\xf9\xff\x59\x8d\x04\x92\xe1\xf4\x93\x40\xf8\x8d\xce\xa8\x0d\x69\x24\x61\xfc\x18\xfc\x0c\x00\x00\xff\xff\xc4\x0b\xce\xc6\x25\x05\x00\x00")

func assetsRdbmsInitSqlBytes() ([]byte, error) {
	return bindataRead(
		_assetsRdbmsInitSql,
		"assets/rdbms/init.sql",
	)
}

func assetsRdbmsInitSql() (*asset, error) {
	bytes, err := assetsRdbmsInitSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/rdbms/init.sql", size: 1317, mode: os.FileMode(420), modTime: time.Unix(1557076824, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _deploymentsRdbmsMigrations01_create_tablesDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x28\x29\x4a\xcc\x2b\x4e\x4c\x2e\xc9\xcc\xcf\x2b\xb6\xe6\xc2\xaa\x24\x31\x39\x39\xbf\x34\xaf\xa4\xd8\x9a\x0b\x10\x00\x00\xff\xff\xa4\x3c\x65\xad\x42\x00\x00\x00")

func deploymentsRdbmsMigrations01_create_tablesDownSqlBytes() ([]byte, error) {
	return bindataRead(
		_deploymentsRdbmsMigrations01_create_tablesDownSql,
		"deployments/rdbms/migrations/01_create_tables.down.sql",
	)
}

func deploymentsRdbmsMigrations01_create_tablesDownSql() (*asset, error) {
	bytes, err := deploymentsRdbmsMigrations01_create_tablesDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "deployments/rdbms/migrations/01_create_tables.down.sql", size: 66, mode: os.FileMode(420), modTime: time.Unix(1555667022, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _deploymentsRdbmsMigrations01_create_tablesUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x93\x31\x6f\xb3\x40\x0c\x86\x77\x7e\x85\x47\x90\xbe\xe9\x93\x32\x65\xba\x82\x49\x51\xe1\x88\x0e\x53\x35\x13\xba\xc0\x35\x42\x2d\x10\x11\xf2\xff\xab\x2b\x81\x40\x01\xb5\x89\xda\x9b\x90\xfd\xda\xb2\x9f\xd7\xc4\x11\xc2\x7e\xbf\x36\x0c\x5b\x20\x23\x04\x62\x0f\x3e\x82\x4c\xd3\xea\x5c\x36\x27\xc3\x34\x00\x00\xf2\x0c\xba\xe7\x71\x82\x98\x47\xde\x86\xa3\x03\x3c\x24\xe0\xb1\xef\x03\x8b\x29\x4c\x3c\x6e\x0b\x0c\x90\x13\x6c\x85\x17\x30\xb1\x83\x27\xdc\xfd\xfb\x6c\x90\xd6\x4a\x36\x2a\x4b\xaa\x12\xc8\x0b\x30\x22\x16\x6c\x75\xbc\x6f\xe0\xa0\xcb\x62\x9f\xc0\x8e\x85\x40\x4e\x49\xaf\x6a\xeb\xcf\xc7\x6c\xbe\x5e\xd7\x86\x1c\xe2\xad\xa3\x87\x5f\xa8\xce\xd4\xbb\x5a\xaa\x6e\x15\xa5\x2c\xd4\x65\xc1\x67\x26\xec\x47\x26\xcc\xff\xab\x95\xd5\xcf\xd7\xaa\x8e\xb2\x56\x65\x93\x68\x1a\x63\x0c\xbd\xc2\x0e\x79\x44\x82\xe9\xac\x1b\x0a\xf4\x36\x5c\x33\x80\xd7\xb7\xa4\x23\x9a\x5c\x9b\x98\xfd\xa7\x05\x02\x5d\x14\xc8\x6d\x8c\x7a\xf6\x60\xe6\x99\x65\x74\xdc\xaf\x4b\x0a\x8c\x48\x78\x36\x0d\x53\x0e\xfa\x38\x4c\x59\xeb\xb1\x9f\x4d\x2d\xcb\x93\x4c\x9b\xbc\x2a\x67\x3c\x9d\xb3\x16\xee\x32\xf7\xf2\x46\x94\xef\x72\x79\xa1\xd1\xad\x76\x2f\xb7\x69\xa5\x55\x9d\x1f\xf2\x72\x08\x61\x64\x3f\x7c\x39\x80\x53\x75\xae\x53\xd5\x79\xa9\x3d\x5c\x60\xd6\xca\x1b\x59\x1f\x54\xf3\x63\xb9\x2c\xb4\x6c\x38\x8c\xeb\x87\x6c\xe6\x77\xbb\x60\xaf\x8a\x42\xe9\x3b\xe9\x1f\xe1\x0b\x7d\x77\x87\xc3\x4b\x48\xa6\xfb\x98\x93\xd0\x6f\xdf\xe6\x4d\x03\x4e\x09\x9a\x93\xd0\x5f\xfc\x3c\x1f\x01\x00\x00\xff\xff\xca\x1a\x6a\x3a\x16\x05\x00\x00")

func deploymentsRdbmsMigrations01_create_tablesUpSqlBytes() ([]byte, error) {
	return bindataRead(
		_deploymentsRdbmsMigrations01_create_tablesUpSql,
		"deployments/rdbms/migrations/01_create_tables.up.sql",
	)
}

func deploymentsRdbmsMigrations01_create_tablesUpSql() (*asset, error) {
	bytes, err := deploymentsRdbmsMigrations01_create_tablesUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "deployments/rdbms/migrations/01_create_tables.up.sql", size: 1302, mode: os.FileMode(420), modTime: time.Unix(1557088053, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _deploymentsRdbmsSchemaSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\x8c\xc1\x6a\x83\x40\x14\x45\xf7\x7e\xc5\x25\x50\x04\xa1\x66\xd3\x45\x21\x9b\x4e\xf4\xa5\x15\x06\x2d\xe3\xeb\xa2\xab\x32\x4f\x5f\x88\x54\x6b\x71\xc6\x42\xfe\xbe\x98\x76\x79\x2f\xe7\x9c\x7d\x96\x20\x03\x5f\x86\x80\xf3\x30\x2a\x86\x80\xce\x8f\xa3\xf6\x90\x2b\x76\xfd\xdc\x7d\xea\x72\xdf\xcd\xd3\xf7\x1c\x34\xbf\x4e\xe3\x0e\x71\x46\xb7\xa8\x8f\x8a\x78\x51\xf4\x3e\x7a\xf1\x41\xf3\x04\xb7\x52\x15\x51\xb5\xa8\x1b\xc6\x1a\xfe\x22\x1b\xe5\xbb\xb8\xfa\x11\x41\x97\x1f\x5d\xc2\xc6\xee\x93\xc2\x91\x61\x42\xe3\xe0\xe8\xd5\x9a\x82\x50\x1a\x36\x47\xd3\x12\x44\x50\xd2\xc9\xbc\x59\x46\xf1\x62\x9c\x29\x98\x1c\x5a\x62\xac\xf1\xfc\x38\xc9\x03\x8a\xc6\xda\x4d\xfe\xdf\x1f\x32\x7c\x1d\x92\x67\x67\x6a\x86\xb1\x16\x4d\x0d\x91\x3c\x03\x37\x48\x45\xd2\xa7\xf4\x2e\x45\x55\x52\xcd\xd5\xa9\xa2\x12\xc7\xf7\xdb\x7d\x48\x7e\x03\x00\x00\xff\xff\xf3\x42\x1b\x54\xfd\x00\x00\x00")

func deploymentsRdbmsSchemaSqlBytes() ([]byte, error) {
	return bindataRead(
		_deploymentsRdbmsSchemaSql,
		"deployments/rdbms/schema.sql",
	)
}

func deploymentsRdbmsSchemaSql() (*asset, error) {
	bytes, err := deploymentsRdbmsSchemaSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "deployments/rdbms/schema.sql", size: 253, mode: os.FileMode(420), modTime: time.Unix(1555666955, 0)}
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
	"assets/rdbms/init.sql":                                  assetsRdbmsInitSql,
	"deployments/rdbms/migrations/01_create_tables.down.sql": deploymentsRdbmsMigrations01_create_tablesDownSql,
	"deployments/rdbms/migrations/01_create_tables.up.sql":   deploymentsRdbmsMigrations01_create_tablesUpSql,
	"deployments/rdbms/schema.sql":                           deploymentsRdbmsSchemaSql,
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
	"assets": &bintree{nil, map[string]*bintree{
		"rdbms": &bintree{nil, map[string]*bintree{
			"init.sql": &bintree{assetsRdbmsInitSql, map[string]*bintree{}},
		}},
	}},
	"deployments": &bintree{nil, map[string]*bintree{
		"rdbms": &bintree{nil, map[string]*bintree{
			"migrations": &bintree{nil, map[string]*bintree{
				"01_create_tables.down.sql": &bintree{deploymentsRdbmsMigrations01_create_tablesDownSql, map[string]*bintree{}},
				"01_create_tables.up.sql":   &bintree{deploymentsRdbmsMigrations01_create_tablesUpSql, map[string]*bintree{}},
			}},
			"schema.sql": &bintree{deploymentsRdbmsSchemaSql, map[string]*bintree{}},
		}},
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
