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

var __000001_init_schemaUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x58\x5d\x6f\xa3\x38\x14\x7d\xef\xaf\xf0\x5b\x12\x69\x1f\xba\xbb\xda\xa7\x3e\xd1\xc4\xad\xd0\x26\xa4\x22\x54\xbb\xd1\x68\x84\xdc\x70\x93\x7a\x0a\x36\x63\x3b\x69\xfb\xef\x47\x50\x3e\x1c\xc0\xe0\x51\x99\x6a\xca\x53\xe4\x7b\xb1\x8f\xcf\x3d\xf7\x83\x5c\xe3\x5b\xd7\xbb\xba\x40\x08\xa1\xb9\x8f\x9d\x00\xa3\xc0\xb9\x5e\x62\xe4\xde\x20\x6f\x1d\x20\xfc\xbf\xbb\x09\x36\xe8\x70\xa4\x71\x24\xd1\x34\xf7\xcb\x1e\x1a\xa1\x13\x11\xbb\x47\x22\xa6\x7f\x5f\xce\xfe\xa8\xd6\x19\x3c\x87\x29\x51\x82\xb3\x30\x01\x29\xc9\x01\x2a\xbf\x3f\x2f\x2f\x75\xcf\x54\xc0\x9e\xbe\x80\xd4\xec\xb3\x2f\x5f\x6b\xfb\x9d\xef\xae\x1c\x7f\x8b\xfe\xc5\xdb\x29\x8d\x66\xf9\xfa\xec\x1c\xe8\xf6\x0e\x23\xc1\x63\x08\xd5\x6b\x0a\x88\x48\x84\xbd\xfb\x95\x06\x72\xc2\xb8\x48\x48\x3c\xa9\x37\x9d\x24\x3c\x02\x41\x14\x17\xfa\x22\x89\x12\xca\xf4\x85\xec\x06\xc0\xd9\xa4\xeb\xd0\x0e\x76\x32\x0c\x36\xe4\xd4\x58\xab\x5f\xb5\x31\xa7\x38\x34\xbd\xda\x60\xa3\x36\xcc\xd7\xde\x26\xf0\x1d\xd7\x0b\xd0\xfe\x29\xcc\x37\xa9\x6c\xd9\x73\xb3\xf6\xb1\x7b\xeb\xe5\x2f\x96\x47\xcc\xce\x3c\xb2\xc7\xc7\x37\xd8\xc7\xde\x1c\x97\xa1\x9e\x76\x79\xad\x3d\xb4\xc0\x4b\x1c\x60\x34\x77\x36\x73\x67\x81\x2d\xf9\xd9\x3d\x12\xc6\x20\xb6\xa1\xa8\x9f\x85\x3c\x50\xe8\x81\xf3\x18\x08\xab\x97\x23\xd8\x93\x63\xac\xc2\xe2\x9c\xb6\x43\xad\xca\xb6\xed\x53\x33\x7b\x94\x20\x6c\x68\x95\x0a\x48\x62\xa6\x75\xa7\xe8\x09\xda\xd4\x48\x45\xf6\xfb\xf6\xb2\x21\x0a\x16\x09\xdb\x71\x83\xa2\x4a\xd8\x5c\x62\x27\x80\x28\x88\x42\xa2\x90\xa2\x09\x48\x45\x92\xb4\xb6\x26\x3c\xa2\x7b\x6a\x34\xef\x38\x53\xc0\x54\xb5\xf3\x5f\xad\x6a\x74\xa2\xfc\x28\xc3\x4e\x3f\xbd\x2e\x15\x22\x33\x93\x79\x54\x8f\x5c\x18\xcd\x90\x3c\x40\x84\xbe\x49\xce\x3a\x0e\xcf\x8d\x32\xb7\xf6\x94\x42\x93\x44\x0b\x64\x46\x91\xd6\xc8\x7b\x65\x5a\x66\xab\x9d\x50\x4d\x60\x32\x65\x1a\x91\x54\x1c\xf5\x02\xc9\xc5\x3d\x80\x62\x83\x03\xe4\xdd\x2f\x97\x96\x6a\x2b\x0a\x7b\x98\x4b\x89\x8b\x73\xd1\x51\xa6\xe0\x00\x02\xdd\x62\x0f\xfb\x4e\x80\x17\xc8\x59\xfe\xe7\x6c\x37\xc8\xd9\x20\x77\x81\xbd\xc0\x0d\xb6\x0d\x35\x72\xa1\xb7\x37\x2d\x32\x31\x65\x4f\xba\x86\xac\x4b\xdc\xa7\x2e\x47\x25\xbd\x8a\x42\x83\xdb\x9f\x24\x97\x91\x04\x0c\xcc\x46\x20\x77\x82\xa6\x8a\x72\x66\x1a\x2d\xca\xd0\x14\x67\x9e\x37\x60\xd3\xdc\xf2\xa2\xde\x60\xb7\x5e\xb2\xcd\xbe\xb7\x43\xcd\xd9\xf7\x66\xef\x0d\x49\x43\x9e\xef\xcb\xc0\xec\xb2\x46\x34\x99\xb1\x17\x4a\x3e\xd3\x58\x26\x9f\x09\x41\xc6\xa7\x11\x41\xc5\xb8\x15\x23\x99\xe3\xb8\xa5\x40\xc0\xf7\x23\x48\x25\xdf\x53\x03\x06\x6a\xfd\x50\xab\x68\xf6\x9a\x86\x8c\x0b\x84\x3d\x2d\x2d\x49\x63\xb0\xb0\xb7\x7b\x75\xfd\xea\xc3\x6b\x37\xb6\xa2\x2f\xf7\x60\x2f\x77\x18\x98\xf3\x2d\xd3\xe7\x83\xfa\xc5\x40\xce\xfc\x56\x2d\x54\x0f\x92\x19\x91\xe6\x34\x5a\x37\x35\x21\x2a\x42\x6d\x04\x53\x8b\xa6\x17\x4a\x39\xf3\x8d\xfb\x91\xc1\x93\x04\xd8\x2f\xcd\xe7\x22\x21\xc3\x7a\xdf\xb3\x7c\xc8\x8e\xb7\x9c\x3e\x3f\x73\x8e\x14\x2c\x98\x5b\x4b\xc5\x52\x7f\x83\x29\xea\xef\xa8\x22\x90\x20\x4e\xad\x2f\xa2\xd1\x46\x0f\x9a\x86\x24\x8a\x04\x48\xed\x4f\x8b\x7f\xf4\x8f\x08\x2e\x54\x5b\x19\x29\x91\xf2\x99\x8b\xc8\x30\x0e\x92\x18\x84\x92\xe1\x50\xaf\xe8\x9f\x1a\x29\xdb\xf3\xc1\x2d\x72\xa7\xa1\xba\x2e\x41\x29\xca\x0e\x72\xd0\xd1\x52\xb3\xf9\xf5\x06\xeb\x6a\x8b\x84\x31\xca\xeb\x50\x35\xfb\xe0\x19\xda\x04\x43\x8f\x9d\x11\x4d\x23\xc0\x1f\x41\x8f\x2e\x97\x7e\x5c\xe3\x96\xfd\x21\x5c\x4d\x85\x1a\xb1\x75\x48\x79\x4c\x7c\x65\x49\x9a\xaf\x57\x2b\x37\xb8\xba\xf8\x11\x00\x00\xff\xff\xb3\xc0\x11\x7f\x4a\x15\x00\x00")

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

	info := bindataFileInfo{name: "000001_init_schema.up.sql", size: 5450, mode: os.FileMode(420), modTime: time.Unix(1611216889, 0)}
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
