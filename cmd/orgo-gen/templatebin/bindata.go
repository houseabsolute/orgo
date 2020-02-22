// Code generated for package templatebin by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/schema.tpl
package templatebin

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

var _templatesSchemaTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x91\xb1\x8e\xa3\x30\x10\x86\x7b\x3f\xc5\x08\xe5\x24\x88\x72\x49\x8f\x94\xea\x8a\xeb\x52\x84\x74\x51\x0a\x43\x26\x04\x25\xd8\xc4\x63\xee\xb4\x9a\x9d\x77\x5f\xd9\x86\x44\xbb\xd2\x36\x2b\x0a\xc4\xcf\x37\xff\xc0\xe7\x41\x37\x37\xdd\x22\x30\xaf\xab\xe6\x8a\xbd\xde\xe9\x1e\x45\x94\xea\xfa\xc1\x3a\x0f\x59\xdb\xf9\xeb\x58\xaf\x1b\xdb\x6f\xae\x76\x24\xd4\x35\xd9\xfb\xe8\x71\x33\xdc\xda\x4d\xad\x09\x33\xa5\xfc\xdb\x80\x90\xc6\x81\xbc\x1b\x1b\x0f\xac\x00\x00\xce\x9d\xbe\x63\xe3\x43\xd8\x99\x36\x46\xcb\x30\x33\xed\x52\xa2\xd4\x3f\xed\xc0\xeb\xfa\x8e\x04\xbd\x1e\x8e\x89\x3c\x25\xea\x10\x72\xa5\x2e\xa3\x69\xa0\x33\x9d\xcf\x8b\xa9\x77\x1a\xd8\x7e\x33\xc2\x8a\xd9\x69\xd3\x22\xa4\x67\x12\x51\xcc\xbf\x61\xf1\x30\xba\x47\x28\xb7\xb0\x0e\xbf\x09\xef\x30\xb8\xce\xf8\x0b\x64\xbf\x1e\x19\x88\xc4\x6e\xe6\x84\x89\x94\xc1\xca\xc1\xfe\xb1\x67\x14\x59\xc5\x06\x34\xe7\x09\x13\x15\xae\x9f\x2e\x8a\x90\xa3\x48\xfc\xb5\x5f\x19\xda\x57\x91\x8a\x5e\x99\x17\x8e\x44\x3e\x8b\xa5\x24\x1b\x96\x93\xc8\x90\x39\x02\x78\x19\xde\x57\xe1\xf3\xa2\xba\x9c\x66\xae\x98\xdb\xf2\x02\x96\x73\x71\x6a\x74\xe8\x47\x67\xe6\xf7\x29\x7b\x6d\x2a\x81\x56\xcf\xc8\x51\x19\xef\x71\xcf\x0e\xff\xef\xab\x3c\x1d\xc8\xf1\xe5\xee\x54\xac\x9e\x9a\x98\x93\xb6\x8f\x00\x00\x00\xff\xff\x8f\xfc\xd8\xcd\x6c\x02\x00\x00")

func templatesSchemaTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesSchemaTpl,
		"templates/schema.tpl",
	)
}

func templatesSchemaTpl() (*asset, error) {
	bytes, err := templatesSchemaTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/schema.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"templates/schema.tpl": templatesSchemaTpl,
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
	"templates": &bintree{nil, map[string]*bintree{
		"schema.tpl": &bintree{templatesSchemaTpl, map[string]*bintree{}},
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
