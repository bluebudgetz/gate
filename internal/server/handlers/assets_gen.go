// Code generated for package handlers by go-bindata DO NOT EDIT. (@generated)
// sources:
// delete_account.cyp
// delete_tx.cyp
// get_account.cyp
// get_accounts_tree.cyp
// get_tx.cyp
// get_tx_list.cyp
// patch_account.cyp
// patch_tx.cyp
// post_account.cyp
// post_tx.cyp
// put_account.cyp
// put_tx.cyp
package handlers

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

var _delete_accountCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\x75\x0c\x71\xf6\x50\xd0\xc8\xb3\x72\x4c\x4e\xce\x2f\xcd\x2b\x51\xa8\xce\x4c\xb1\x52\x50\xc9\x4c\xa9\xd5\x54\x70\x71\x0d\x71\x74\xf6\x50\x70\x71\xf5\x71\x0d\x71\x55\xd0\xc8\xd3\x04\x04\x00\x00\xff\xff\x16\x07\x45\xe5\x2d\x00\x00\x00")

func delete_accountCypBytes() ([]byte, error) {
	return bindataRead(
		_delete_accountCyp,
		"delete_account.cyp",
	)
}

func delete_accountCyp() (*asset, error) {
	bytes, err := delete_accountCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "delete_account.cyp", size: 45, mode: os.FileMode(420), modTime: time.Unix(1585334896, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _delete_txCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\x75\x0c\x71\xf6\x50\xd0\xb0\x72\x4c\x4e\xce\x2f\xcd\x2b\xd1\xd4\x8d\x2e\xa9\xb0\x0a\x48\xcc\x4c\x51\xa8\xce\x4c\xb1\x52\x50\xc9\x4c\xa9\x8d\xd5\xb5\x43\xc8\x73\xb9\xb8\xfa\xb8\x86\xb8\x2a\x68\x94\x54\x68\x72\x01\x02\x00\x00\xff\xff\x1a\xe0\x03\xd8\x3d\x00\x00\x00")

func delete_txCypBytes() ([]byte, error) {
	return bindataRead(
		_delete_txCyp,
		"delete_tx.cyp",
	)
}

func delete_txCyp() (*asset, error) {
	bytes, err := delete_txCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "delete_tx.cyp", size: 61, mode: os.FileMode(420), modTime: time.Unix(1585335716, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _get_accountCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\xd2\x51\xab\x9b\x30\x14\x07\xf0\xf7\x7c\x8a\xff\xc3\x28\x76\xa8\xdd\x73\x69\x0b\x6e\x73\xb4\xb0\xd5\x62\x1d\x65\x94\x3e\x9c\xc5\x54\x03\x31\x29\x31\x6e\x74\xec\xc3\x0f\xb5\xda\x6e\xbb\x5c\xb8\x70\xdf\x34\xe7\xe4\xe4\xf7\x27\x99\xcd\xf0\x49\x38\x5e\x82\x38\x37\x8d\x76\x28\x2c\x5d\x4a\x1f\xb5\x23\xeb\xa4\x2e\xf0\x53\xba\x12\x6f\x64\xee\x83\x74\x0e\x52\x0a\xd2\xd5\xc8\xa5\x15\xdc\x61\x02\xa9\x6f\x9f\xbc\x94\x2a\xb7\x42\xe3\x6c\x2c\x04\xf1\x12\xda\xe4\x82\x7d\x89\xb2\x0f\x6b\x78\x34\x8f\xfa\xf9\xd3\x45\x70\x9c\x77\xbd\xc9\xf9\xed\xbb\x30\x3c\x05\x1e\xd7\x63\x91\x1d\xd6\x71\x1a\x83\x42\x99\x63\xd9\x9e\xca\x0e\x9b\x6c\x0d\xf2\x71\xe4\xd8\x6c\xc1\x8d\x52\x82\x3b\x8f\xeb\x29\xfa\x56\xde\xb6\x2e\x56\xfd\x96\xdf\xdd\xef\x09\xd1\x7e\xe4\x30\x36\x26\x9c\xa0\x6e\xaa\x8a\xac\xfc\x25\xba\x1c\xa6\x71\x85\x69\x23\x5e\xe8\x5a\x09\xed\xea\x7f\xe8\xc9\x2e\xdb\x24\xdb\xe8\x33\xfa\x0c\x17\x2c\xe1\xd1\x13\xfe\x51\x1f\x1c\xed\x7c\x47\x32\x3f\x05\x2b\xef\x21\xd2\x2d\xc1\x20\xf2\x5b\x86\x67\x43\xaa\xba\x7a\x8b\x1d\x24\xcf\x60\xa5\xe6\xa6\x7a\x4d\xec\xe2\xae\x7d\x19\x76\x90\xf8\x7f\xb3\x77\x8d\x52\x10\x3f\x84\xbd\xba\xb2\x75\x3a\x53\x08\x57\x0a\xcb\xd2\x38\xfb\x9a\x6e\x41\xa0\x7a\x78\x64\x8f\xf3\xff\x1f\x77\x5f\x43\x70\xbf\xa4\x68\x8f\xef\xa4\x48\x73\xc1\x80\x24\xfd\x18\xa7\x78\xff\x0d\x14\x6a\xaa\x04\xfb\x13\x00\x00\xff\xff\xb9\x09\x31\x28\xc6\x02\x00\x00")

func get_accountCypBytes() ([]byte, error) {
	return bindataRead(
		_get_accountCyp,
		"get_account.cyp",
	)
}

func get_accountCyp() (*asset, error) {
	bytes, err := get_accountCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "get_account.cyp", size: 710, mode: os.FileMode(420), modTime: time.Unix(1585329126, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _get_accounts_treeCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\xd2\x41\x8b\x9b\x40\x14\x07\xf0\xbb\x9f\xe2\x7f\x0a\xa6\xa8\xe9\x59\xd2\x80\x6d\x2c\x09\xb4\x31\x18\x4b\x28\x21\x87\xd7\x71\xa2\x43\xc7\x19\x19\xc7\x96\x2c\xfb\xe1\x17\xcd\x6a\x76\x37\xcb\xc2\xc2\x9e\x74\xe6\x0d\xef\xfd\xfe\xcc\xcc\x66\xf8\xce\x2d\x2b\x41\x52\x82\x18\xd3\xad\xb2\x8d\x07\x92\x5a\x15\xf8\x2f\x6c\x09\x82\x14\x8d\x85\x3e\x21\x17\x86\x33\x8b\x09\x84\x7a\xfc\x65\xa5\x90\xb9\xe1\x0a\x27\x6d\xc0\x89\x95\x50\x3a\xe7\xce\xcf\x28\xfb\xb6\x82\x4b\x61\x74\x69\x38\x9d\xfb\x87\xb0\x3f\x9b\x9c\x3e\x7d\x0e\x82\xa3\xef\x32\x35\x16\x9d\xfd\x3a\x5b\x81\x3c\x1c\x18\xd6\x1b\x30\x2d\x25\x67\xd6\x65\x6a\x8a\xfd\x2a\x4e\x63\xb0\x40\xe4\x98\x2f\x40\xdd\xf7\xbe\x5f\x1e\x11\xed\xc6\xe9\x8e\x33\xa6\x98\xa0\x69\xab\x8a\x8c\xb8\xe3\x7d\x22\xdd\xda\x42\x0b\x55\xa0\xa6\x73\xc5\x95\x6d\x5e\x48\x93\x6d\xb6\x4e\x36\xd1\x0f\x5c\xc8\x35\xbe\xc0\xa5\x57\xb8\x23\xd6\x3f\x98\x70\x4b\x22\x3f\xfa\x0b\xf7\x36\xc1\x20\xf2\x3a\x86\x6b\x02\xaa\xfa\x7a\x87\x1d\x24\x6f\x60\x85\x62\xba\xfa\x48\xec\xfc\xaa\x7d\x1f\x76\x90\x78\xcf\xd9\xdb\x56\x4a\xf0\x7f\xdc\x9c\x6d\xd9\x39\xad\x2e\xb8\x2d\xb9\x71\xd2\x38\xfb\x95\x6e\x40\xa0\x66\x78\x44\x4f\xfb\xdf\xb6\xbb\xee\xc1\xbf\x5e\x52\xb4\xc3\x1f\x92\xa4\x18\x77\x80\x24\x5d\xc6\x29\xbe\xfe\x06\x05\x8a\x2a\xde\xcd\xcf\x92\x65\x12\xa2\x69\xeb\x5a\x1b\x8b\xe6\xaf\xa8\x67\x52\x54\xc2\x3e\x04\x00\x00\xff\xff\xaa\xb2\x19\x8a\xc5\x02\x00\x00")

func get_accounts_treeCypBytes() ([]byte, error) {
	return bindataRead(
		_get_accounts_treeCyp,
		"get_accounts_tree.cyp",
	)
}

func get_accounts_treeCyp() (*asset, error) {
	bytes, err := get_accounts_treeCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "get_accounts_tree.cyp", size: 709, mode: os.FileMode(420), modTime: time.Unix(1585945096, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _get_txCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\x75\x0c\x71\xf6\x50\xd0\x28\x2e\x4a\xb6\x72\x4c\x4e\xce\x2f\xcd\x2b\xd1\xd4\x8d\x2e\xa9\xb0\x0a\x48\xcc\x4c\x51\xa8\xce\x4c\xb1\x52\x50\xc9\x4c\xa9\x8d\xd5\xb5\xd3\x48\x29\x2e\x81\x2b\xe1\x0a\x72\x0d\x09\x0d\xf2\x53\x28\xa9\xd0\x51\x28\x2e\x4a\xd6\x51\x48\x29\x2e\xe1\x02\x04\x00\x00\xff\xff\x29\x13\xda\x9c\x4b\x00\x00\x00")

func get_txCypBytes() ([]byte, error) {
	return bindataRead(
		_get_txCyp,
		"get_tx.cyp",
	)
}

func get_txCyp() (*asset, error) {
	bytes, err := get_txCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "get_tx.cyp", size: 75, mode: os.FileMode(420), modTime: time.Unix(1585944920, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _get_tx_listCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\xc9\xb1\x0e\x82\x30\x10\x06\xe0\xfd\x9e\xe2\x1f\x18\x30\x01\x1f\x80\xc1\x04\xa1\x89\x8d\x22\xa4\xd4\xc1\x18\x07\xd3\x3a\x5c\x54\x34\xdc\x91\xf4\xf1\xdd\x98\xbf\xae\xf6\xcd\x01\xb9\xcc\xa1\xaa\x43\xf8\x2e\x93\x6e\xca\x9b\xa6\x6a\x78\x70\xbc\x97\xbb\x3c\x8a\xae\x40\xce\xf8\x8b\x3b\x43\x53\x01\x99\x43\x81\x28\x4a\x40\xef\x5a\xe3\xb0\xbf\x42\xd3\x96\x45\x96\x67\xec\x27\xb4\x66\x6c\x08\x18\x8f\x76\x40\x26\x2f\xfe\x11\x70\xb2\x9d\xf5\xc8\xde\xfc\x61\xa5\x7f\x00\x00\x00\xff\xff\xf3\xf8\xa0\x5d\x79\x00\x00\x00")

func get_tx_listCypBytes() ([]byte, error) {
	return bindataRead(
		_get_tx_listCyp,
		"get_tx_list.cyp",
	)
}

func get_tx_listCyp() (*asset, error) {
	bytes, err := get_tx_listCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "get_tx_list.cyp", size: 121, mode: os.FileMode(420), modTime: time.Unix(1585083333, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _patch_accountCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8f\xc1\x6a\xb3\x40\x14\x85\xf7\xf3\x14\x67\x21\x64\x06\x8c\x0f\x20\x7f\x02\x62\x86\x3f\x82\x89\x60\xa6\x94\x52\xba\xb8\xcc\x4c\x89\x60\xae\x62\xb5\x5d\x04\xdf\xbd\x68\x4d\x69\xa0\x8b\xae\x86\xb9\xf7\x9c\xef\x9e\x73\xd0\xe5\x7f\x0d\xc9\x71\x62\x6d\x33\x70\x8f\x6b\xe5\x62\x04\x95\x1b\x95\x00\x8a\x23\xd2\x52\x27\x46\xe3\xa4\x0d\x98\x2e\x1e\x1b\xd8\x86\x6a\xff\x66\xbd\x0c\xa6\x41\x08\x8e\xa6\x57\x85\xb0\x9d\xa7\xde\xbb\x82\xb1\x81\xa3\xde\xf7\xd5\xc5\xcb\x05\x73\x48\x4c\xba\xff\x13\x65\x68\xdd\x6f\x94\xc7\xcc\xec\xc1\x22\x4d\xf2\x1c\xd4\x36\x36\x72\x4d\xf4\x71\xf6\x2c\x45\xd0\x52\xe7\xb9\xcf\x1c\xfe\x6d\xc1\x43\x5d\x87\x62\xb5\xd4\x6a\xef\x6a\xdd\x74\xa3\xc2\xb2\x27\xb5\x7e\x8e\xed\xb9\xaa\x5d\xf1\xfa\xb2\xde\xca\x56\xa1\xd4\xe6\xa1\x3c\x82\x56\x13\x65\x0e\x3d\xab\xba\x9f\x32\x85\x9d\xce\xb5\xd1\x90\xdd\x9d\xe1\x4a\x31\x38\xc4\xed\x4e\x8c\xef\x68\xa3\x50\xe2\x29\xd3\xf9\x0e\xef\x54\x0f\x5e\x2c\xa6\xf9\x83\xe4\x04\xfa\x8a\x29\x3e\x03\x00\x00\xff\xff\xde\x5f\xec\x73\x8f\x01\x00\x00")

func patch_accountCypBytes() ([]byte, error) {
	return bindataRead(
		_patch_accountCyp,
		"patch_account.cyp",
	)
}

func patch_accountCyp() (*asset, error) {
	bytes, err := patch_accountCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "patch_account.cyp", size: 399, mode: os.FileMode(420), modTime: time.Unix(1586289152, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _patch_txCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\x75\x0c\x71\xf6\x50\xd0\x28\x2e\x4a\xb6\x72\x4c\x4e\xce\x2f\xcd\x2b\xd1\xd4\x8d\x2e\xa9\xb0\x0a\x48\xcc\x4c\x51\xa8\xce\x4c\xb1\x52\x50\xc9\x4c\xa9\x8d\xd5\xb5\xd3\x48\x29\x2e\x81\x2b\xe1\x0a\x76\x0d\x51\x28\xa9\xd0\xcb\x2c\x2e\x2e\x4d\x4d\xf1\xcf\x53\xb0\x55\x48\xce\x4f\xcc\x49\x2d\x4e\x4e\xd5\x50\x81\x09\xea\x20\xab\x80\x6b\xc9\x2f\xca\x4c\xcf\x44\xd5\x00\x11\xd2\x41\xc8\xc2\x15\x27\xe6\x82\xac\x43\x51\x0c\x11\xd2\x41\xc8\xc2\x15\x27\xe7\xe7\xe6\xa6\xa2\xa9\x86\x8a\xe9\x20\xc9\x6b\x72\x05\xb9\x86\x84\x06\xf9\x29\x14\x17\x25\x83\xc4\x75\x14\x52\x8a\x4b\xb8\x00\x01\x00\x00\xff\xff\xf3\x49\x63\x61\x08\x01\x00\x00")

func patch_txCypBytes() ([]byte, error) {
	return bindataRead(
		_patch_txCyp,
		"patch_tx.cyp",
	)
}

func patch_txCyp() (*asset, error) {
	bytes, err := patch_txCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "patch_tx.cyp", size: 264, mode: os.FileMode(420), modTime: time.Unix(1586288176, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _post_accountCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x8f\x4d\x6a\xc3\x30\x14\x84\xf7\xef\x14\xb3\x08\x58\x02\xc5\x07\x10\x25\x60\x52\xd1\x1a\xdc\x06\x5c\x97\x52\x4a\x17\x42\x7a\x25\x06\x47\x0e\x41\x6a\x17\xc1\x77\x2f\xf5\x1f\x34\xab\x37\xc3\x8c\x3e\x46\xfb\xda\x14\x8d\x81\xb0\xce\xf5\x29\x44\x5d\x4c\x17\xd7\xd6\x6b\xd8\x73\xef\x72\x77\x61\x1b\x39\x4f\xa9\xf5\x42\x2a\x04\x7b\x62\x8d\xcd\xdf\x51\x98\x32\x7f\x08\x1a\xde\x46\x8e\xed\x89\x85\x1c\x24\xbd\x95\xcd\x23\x66\x24\xed\x8b\xaa\x9a\x50\xbe\xcf\x7f\x8e\x1c\x04\x6d\xce\xf6\xc2\x21\x96\x1e\x77\x3b\x84\xd4\x75\x8a\xb2\x27\x53\x3f\x18\x88\x29\xf9\xb7\x63\x29\x0f\x12\x73\xc9\x1d\xdb\xce\xcb\xed\x87\x1e\xc5\xe1\xeb\x73\xbb\x9b\x1f\x4a\xd4\xa6\x79\xad\x9f\x31\x26\x99\xa2\xec\xc6\x5f\x47\xa1\x97\x75\x0a\x0b\x5d\x63\x5d\x35\x90\xa4\xf7\xd2\x54\xf7\xf8\xb6\x5d\x62\x9a\x11\xa3\x41\xf1\xb2\xfe\xec\x37\x00\x00\xff\xff\x82\x04\xbd\xbd\x3d\x01\x00\x00")

func post_accountCypBytes() ([]byte, error) {
	return bindataRead(
		_post_accountCyp,
		"post_account.cyp",
	)
}

func post_accountCyp() (*asset, error) {
	bytes, err := post_accountCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "post_account.cyp", size: 317, mode: os.FileMode(420), modTime: time.Unix(1586290830, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _post_txCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\xd0\xbd\x8a\xc4\x20\x10\xc0\xf1\xde\xa7\x98\xc2\x22\x82\xc9\x03\x58\x1c\x84\x10\xb8\x2b\xee\x83\x90\xab\x8e\x2b\x64\x94\x60\x91\xb8\xe8\x08\x81\x25\xef\xbe\x18\xb3\xb0\x2c\x9b\xa9\xfe\xc8\x0f\x95\xf9\x6c\xc7\xee\x1d\xaa\x18\x50\xb5\x88\x3e\x2d\x04\x57\x67\x14\xf0\xe8\x53\x40\x7b\x9c\x7d\x98\x4d\x48\xa8\x4c\xa4\x27\x45\x3a\x4c\x96\x1e\x14\xeb\x86\xbe\x1d\xfb\x7c\xa1\xa8\xff\x68\x55\x3f\xda\x99\x82\xcb\xe8\x8b\xc7\x06\x83\xd5\x64\x9b\x94\x9c\xa9\x84\x64\xf0\x72\x5c\x8c\xc9\x9a\xef\x45\x01\xbf\xe7\x19\xf5\xc1\x4d\x6e\xc9\x4f\xf0\x92\x67\x50\xcf\xf9\xa3\x3b\x2c\x79\x06\xd1\xcf\xb3\xdd\x25\x3f\x72\xfb\xaf\xdf\xf2\x06\x04\x1b\xfa\xf1\x77\xf8\x82\x18\x50\x02\xad\x12\x4c\x24\x76\x0b\x00\x00\xff\xff\xf7\x96\xfd\x12\x48\x01\x00\x00")

func post_txCypBytes() ([]byte, error) {
	return bindataRead(
		_post_txCyp,
		"post_tx.cyp",
	)
}

func post_txCyp() (*asset, error) {
	bytes, err := post_txCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "post_tx.cyp", size: 328, mode: os.FileMode(420), modTime: time.Unix(1586291732, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _put_accountCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\xd1\x6a\x83\x30\x14\x86\xef\xf3\x14\xff\x45\xa1\x06\x6c\x1f\x40\xd6\x82\xd8\xc3\x2a\xd8\x0a\x2e\x63\x8c\xb1\x8b\x90\x64\x54\xb0\x51\x5c\xdc\x2e\x8a\xef\x3e\x9a\x46\xd9\x58\xaf\xce\x09\xe7\x7c\x3f\x5f\xce\x81\xaa\x47\x42\x24\x95\x6a\x07\xeb\x92\xf4\x56\x71\xa9\x75\x82\x45\xad\x47\xce\x80\xf2\x88\xac\xa2\x54\x10\x9e\x48\xc0\xca\xb3\xc1\x06\x8b\x6b\x8d\xa1\x7a\x23\x9d\xd1\xa5\xc5\x06\x5a\x3a\xe3\xea\xb3\x89\x02\x74\x48\x45\xb6\xbf\xc3\x0c\x9d\xbe\xc7\xbc\xe4\x62\x8f\x20\xc2\xb2\xb4\x28\x20\xbb\x56\xad\x75\xbb\xfe\x3e\x19\x1b\xb1\x45\x27\x7b\x63\x5d\xae\xf1\xb0\x85\x1d\x9a\x26\x66\xcb\xa0\x7f\x9b\xfc\xb1\x9f\x96\x47\x8e\xb0\xa4\x4e\x75\xa3\xf9\xea\x2d\xf1\x4d\xf9\xf1\xbe\xda\x06\x90\xa3\x22\xf1\x5c\x1d\xe1\x27\xcb\x6b\xae\x57\x9f\x91\xfe\x37\xc3\xb1\xa3\x82\x04\x21\xea\xff\x81\x17\xdf\x24\xd3\x37\x62\x4c\x1a\x09\x66\xfd\x91\x71\xf6\x9a\x53\xb1\xc3\x97\x6c\x06\xc3\x42\x84\x7f\x40\x7e\xce\x27\xf8\x09\x00\x00\xff\xff\xcc\x94\x7c\xde\x9b\x01\x00\x00")

func put_accountCypBytes() ([]byte, error) {
	return bindataRead(
		_put_accountCyp,
		"put_account.cyp",
	)
}

func put_accountCyp() (*asset, error) {
	bytes, err := put_accountCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "put_account.cyp", size: 411, mode: os.FileMode(420), modTime: time.Unix(1586293399, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _put_txCyp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\x75\x0c\x71\xf6\x50\xd0\x28\x2e\x4a\xb6\x72\x4c\x4e\xce\x2f\xcd\x2b\xd1\xd4\x8d\x2e\xa9\xb0\x0a\x48\xcc\x4c\x51\xa8\xce\x4c\xb1\x52\x50\xc9\x4c\xa9\x8d\xd5\xb5\xd3\x48\x29\x2e\x81\x2b\xe1\x0a\x76\x0d\x51\x28\xa9\xd0\xcb\x2c\x2e\x2e\x4d\x4d\xf1\xcf\x53\xb0\x55\x50\x81\xb1\x61\x72\xf9\x45\x99\xe9\x99\x60\x19\x08\x0b\x26\x9e\x98\x0b\x32\x02\x24\x0e\x61\xc1\xc4\x93\xf3\x73\x73\x53\x21\x12\x50\x26\x57\x90\x6b\x48\x68\x90\x9f\x42\x71\x51\xb2\x8e\x42\x49\x85\x8e\x42\x4a\x71\x09\x17\x20\x00\x00\xff\xff\xbb\x1d\x4e\x6a\xb1\x00\x00\x00")

func put_txCypBytes() ([]byte, error) {
	return bindataRead(
		_put_txCyp,
		"put_tx.cyp",
	)
}

func put_txCyp() (*asset, error) {
	bytes, err := put_txCypBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "put_tx.cyp", size: 177, mode: os.FileMode(420), modTime: time.Unix(1586293488, 0)}
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
	"delete_account.cyp":    delete_accountCyp,
	"delete_tx.cyp":         delete_txCyp,
	"get_account.cyp":       get_accountCyp,
	"get_accounts_tree.cyp": get_accounts_treeCyp,
	"get_tx.cyp":            get_txCyp,
	"get_tx_list.cyp":       get_tx_listCyp,
	"patch_account.cyp":     patch_accountCyp,
	"patch_tx.cyp":          patch_txCyp,
	"post_account.cyp":      post_accountCyp,
	"post_tx.cyp":           post_txCyp,
	"put_account.cyp":       put_accountCyp,
	"put_tx.cyp":            put_txCyp,
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
	"delete_account.cyp":    &bintree{delete_accountCyp, map[string]*bintree{}},
	"delete_tx.cyp":         &bintree{delete_txCyp, map[string]*bintree{}},
	"get_account.cyp":       &bintree{get_accountCyp, map[string]*bintree{}},
	"get_accounts_tree.cyp": &bintree{get_accounts_treeCyp, map[string]*bintree{}},
	"get_tx.cyp":            &bintree{get_txCyp, map[string]*bintree{}},
	"get_tx_list.cyp":       &bintree{get_tx_listCyp, map[string]*bintree{}},
	"patch_account.cyp":     &bintree{patch_accountCyp, map[string]*bintree{}},
	"patch_tx.cyp":          &bintree{patch_txCyp, map[string]*bintree{}},
	"post_account.cyp":      &bintree{post_accountCyp, map[string]*bintree{}},
	"post_tx.cyp":           &bintree{post_txCyp, map[string]*bintree{}},
	"put_account.cyp":       &bintree{put_accountCyp, map[string]*bintree{}},
	"put_tx.cyp":            &bintree{put_txCyp, map[string]*bintree{}},
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
