// Code generated for package migrations by go-bindata DO NOT EDIT. (@generated)
// sources:
// 000001_init_schema.down.sql
// 000001_init_schema.up.sql
package migrations

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

var __000001_init_schemaDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\xcf\x3d\x0e\xc2\x30\x0c\x86\xe1\x3d\xa7\xf0\x3d\x32\x51\x08\x28\x12\xa5\x15\xcd\x00\x53\x15\x15\xab\x54\xca\x4f\xb1\x13\x24\x6e\xcf\x0c\x83\x61\xf6\x23\x7d\x7e\x1b\x73\xb0\x27\xad\x00\x00\x76\xe7\xae\x07\xb7\x69\x8e\x06\xec\x1e\xcc\xc5\x0e\x6e\x00\x46\x7a\x22\xb1\x20\xa6\x1c\x23\xa6\x22\x11\xc2\x47\x45\x16\xc9\xea\x0b\x61\x4e\x63\x59\x90\xfe\x60\x13\xa1\x2f\x59\x92\x11\x99\xfd\x8c\xd2\x66\xe5\x1f\x65\x77\x9f\x12\x06\xb1\x2c\x87\xcf\x89\x6b\xff\x7d\x1e\xcb\x6b\x15\xdf\x98\xeb\x12\x6e\xac\xd5\xb6\x6b\x5b\xeb\xb4\x7a\x07\x00\x00\xff\xff\x57\xde\x8f\x03\x93\x01\x00\x00")

func _000001_init_schemaDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_init_schemaDownSql,
		"000001_init_schema.down.sql",
	)
}

func _000001_init_schemaDownSql() (*asset, error) {
	bytes, err := _000001_init_schemaDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_init_schema.down.sql", size: 403, mode: os.FileMode(420), modTime: time.Unix(1611206279, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000001_init_schemaUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x58\x4d\x6f\xab\x38\x14\xdd\xf7\x57\x78\x97\x44\xea\xa2\x33\xa3\x59\x75\x45\x13\xb7\x42\x93\x90\x8a\x50\xcd\x44\xa3\x11\x72\xc3\x4d\xea\x29\xd8\x3c\xdb\x49\xdb\x7f\xff\x04\xe5\xc3\x01\x0c\x7e\x2a\xaf\x7a\x65\x15\xf9\x5e\xec\xe3\x73\xcf\xfd\x08\x37\xf8\xce\xf5\xae\x2f\x10\x42\x68\xee\x63\x27\xc0\x28\x70\x6e\x96\x18\xb9\xb7\xc8\x5b\x07\x08\xff\xe3\x6e\x82\x0d\x3a\x1c\x69\x1c\x49\x34\xcd\xfd\xb2\x87\x46\xe8\x44\xc4\xee\x89\x88\xe9\x1f\x57\xb3\xcb\x6a\x9d\xc1\x4b\x98\x12\x25\x38\x0b\x13\x90\x92\x1c\xa0\xf2\xfb\xed\xea\x4a\xf7\x4c\x05\xec\xe9\x2b\x48\xcd\x3e\xfb\xf7\xbf\xda\x7e\xef\xbb\x2b\xc7\xdf\xa2\xbf\xf0\x76\x4a\xa3\x59\xbe\x3e\x3b\x07\xba\xbd\x6f\xe2\x14\x3c\x86\x50\xbd\xa5\x80\x88\x44\xd8\x7b\x58\x69\x90\x27\x8c\x8b\x84\xc4\x93\xfa\x88\x49\xc2\x23\x10\x44\x71\xa1\x2f\x92\x28\xa1\x4c\x5f\xc8\xee\x03\x9c\x4d\xba\x20\x74\x70\x95\x61\xb0\xa1\xaa\xc6\x5a\xfd\xaa\x8d\x39\xe1\xa1\xe9\xd5\x06\x37\xb5\x61\xbe\xf6\x36\x81\xef\xb8\x5e\x80\xf6\xcf\x61\xbe\x49\x65\xcb\x9e\xdb\xb5\x8f\xdd\x3b\x2f\x7f\xb1\x3c\x62\x76\xe6\x91\x3d\x3e\xbe\xc5\x3e\xf6\xe6\xb8\x0c\xfc\xb4\xcb\x6b\xed\xa1\x05\x5e\xe2\x00\xa3\xb9\xb3\x99\x3b\x0b\x6c\xc9\xcf\xee\x89\x30\x06\xb1\x0d\x45\xfd\x2c\xe4\x81\x42\x8f\x9c\xc7\x40\x58\xbd\x1c\xc1\x9e\x1c\x63\xd5\x36\xd4\xda\x6c\xdb\xbe\x34\xa3\x47\x09\xc2\x86\x4e\xa9\x80\x24\x66\x3a\x77\x8a\x9e\xa0\x4d\x8d\x54\x64\xbf\x6f\x2f\x1b\xd8\xb7\x48\xdb\x8e\x1b\x14\xb5\xc2\xe6\x12\x3b\x01\x44\x41\x14\x12\x85\x14\x4d\x40\x2a\x92\xa4\xb5\x35\xe1\x11\xdd\x53\xa3\x79\xc7\x99\x02\xa6\xaa\x9d\x7f\xcf\x6a\x92\x5e\x92\x4e\x94\x1f\x65\xd8\xe9\xa6\x17\xa7\x42\xc3\x66\x2e\x8f\xea\x89\x0b\xa3\x19\x92\x47\x88\xd0\xff\x92\xb3\xcb\xf6\xe1\xb9\x51\xe6\xd6\x9e\x7a\x68\x52\x68\x81\xcc\xa8\xd1\x1a\x79\xaf\x4a\xcb\x24\xb5\xd3\xa9\x09\x4c\x26\x4c\x23\x92\x8a\xa3\x5e\x20\xb9\xb6\x07\x50\x6c\x70\x80\xbc\x87\xe5\xd2\x52\x6c\x45\x3d\x0f\x73\x25\x71\x71\xae\x39\xca\x14\x1c\x40\xa0\x3b\xec\x61\xdf\x09\xf0\x02\x39\xcb\xbf\x9d\xed\x06\x39\x1b\xe4\x2e\xb0\x17\xb8\xc1\xb6\x21\x46\x2e\xf4\x1e\xa7\x45\x26\xa6\xec\x59\xd7\x90\x75\x65\xfb\xd2\xd5\xa8\xa4\x57\x51\x68\x70\xfb\x83\xe4\x32\x92\x80\x81\xd9\x08\xe4\x4e\xd0\x54\x51\xce\x4c\xf3\x45\x19\x9a\xe2\xcc\xf3\xbe\x6b\x1a\x5e\x5e\xd5\x3b\xec\xd6\x4b\xb6\xd9\xf7\x7e\xa8\x39\xfb\xde\xed\xbd\x21\x69\xc8\xf3\x63\x19\x98\x5d\xd6\x88\x26\x33\xf6\x42\xc9\x47\x19\xcb\xe4\x33\x21\xc8\xf8\xec\x42\x10\x66\x08\x2a\xc6\xad\x18\xc9\x1c\xc7\x2d\x05\x02\xbe\x1d\x41\x2a\xf9\x91\x1a\x30\x50\xeb\x87\x5a\x45\xb3\xd7\x34\x64\x5c\x20\xec\xe9\x68\x49\x1a\x83\x85\xbd\xdd\xaa\xeb\x57\x1f\xdf\xba\xb1\x15\x6d\xb9\x07\x7b\xb9\xc3\xc0\xb0\x6f\x99\x3e\x9f\xd4\x2f\x06\x72\xe6\x97\x6a\xa1\x7a\x90\xcc\x88\x34\xa7\xd1\xba\xa9\x09\x51\x11\x6a\x23\x98\x5a\x34\xbd\x50\xca\x91\x6f\xdc\xff\x16\x3c\x49\x80\xfd\xd4\x7c\x2e\x12\x32\xac\xf7\x3d\xcb\x87\xec\x78\xcb\xe1\xf3\x2b\xe7\x48\xc1\x82\xb9\xb5\x54\x2c\xf5\x37\x98\xa2\xfe\x8e\x2a\x02\x09\xe2\xd4\xfa\x43\x34\xda\xe8\x41\xd3\x90\x44\x91\x00\xa9\x7d\xb9\xf8\x53\xff\xae\xc1\x85\x6a\x2b\x23\x25\x52\xbe\x70\x11\x19\xc6\x41\x12\x83\x50\x32\x1c\xea\x15\xfd\x53\x23\x65\x7b\x3e\xb8\x45\xee\x34\x54\xd7\x25\x28\x45\xd9\x41\x0e\x3a\x5a\x6a\x36\xbf\xde\x60\x5d\x6d\x91\x30\x46\x79\x1d\xaa\x66\x9f\x3c\x43\x9b\x60\xe8\xb1\x33\xa2\x69\x04\xf8\x33\xe8\xd1\xe5\xd2\x8f\x6b\xdc\xb2\x3f\x84\xab\xa9\x50\x23\xb6\x0e\x29\x8f\x89\xaf\x2c\x49\xf3\xf5\x6a\xe5\x06\xd7\x17\xdf\x03\x00\x00\xff\xff\x77\x2d\xda\x08\x4f\x15\x00\x00")

func _000001_init_schemaUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_init_schemaUpSql,
		"000001_init_schema.up.sql",
	)
}

func _000001_init_schemaUpSql() (*asset, error) {
	bytes, err := _000001_init_schemaUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_init_schema.up.sql", size: 5455, mode: os.FileMode(420), modTime: time.Unix(1611206743, 0)}
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
	"000001_init_schema.down.sql": _000001_init_schemaDownSql,
	"000001_init_schema.up.sql":   _000001_init_schemaUpSql,
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
	"000001_init_schema.down.sql": &bintree{_000001_init_schemaDownSql, map[string]*bintree{}},
	"000001_init_schema.up.sql":   &bintree{_000001_init_schemaUpSql, map[string]*bintree{}},
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
