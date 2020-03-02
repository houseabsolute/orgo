// Code generated for package templatebin by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/rs.tpl
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

var _templatesRsTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x56\x4b\x6f\xa3\x3c\x14\x5d\xc3\xaf\xb8\x42\xfd\xa4\xa4\x4a\xc9\x1e\x29\xab\x0a\xe5\x8b\x14\x31\x4d\x89\x34\x8b\xaa\x1a\xd1\xd4\x21\xa8\xc4\x76\x6c\xd3\x87\x18\xff\xf7\x91\x6d\x28\x8f\x50\xf2\x98\x89\x46\x93\x0d\xc1\xbe\xdc\x73\x8e\x7d\x38\x86\x46\xab\x97\x28\x46\x90\xe7\xe0\x4e\xc9\xdd\x4b\x1c\x6e\x08\x13\x41\xb4\x45\x20\xa5\x6d\x27\x5b\x4a\x98\x00\x27\x4e\xc4\x26\x7b\x72\x57\x64\x3b\xde\x90\x8c\xa3\xe8\x89\x93\x34\x13\x68\x4c\x58\x4c\xc6\xf4\x25\x1e\x3f\x45\x1c\x39\xb6\xfd\x1a\x31\x10\x30\xd1\xfd\x96\xe4\x96\x3c\x9b\x3e\x79\x7e\x03\x57\x8c\x83\x37\x51\x30\xba\xfd\x4f\xa0\x2c\xc1\x62\x0d\xce\x7f\xfc\x3e\x74\x54\x99\xf8\xa0\x9a\x89\xaa\x94\x12\xb8\x60\xd9\x4a\x40\x6e\x03\x00\xf0\xd5\x06\x6d\x23\x80\xeb\x50\xff\xd1\x63\x8c\x83\xfe\x5d\x2b\x70\xf7\x3e\xb4\xa5\x6d\xaf\x33\xbc\x82\x00\xbd\x0d\x86\x70\x5d\xb5\x32\x3d\x18\x12\x19\xc3\x15\x82\x19\xad\xba\x7b\xc0\x47\x9f\x43\x8c\x7b\xfa\xaa\x7b\x07\xe8\xed\x3e\x1c\x88\xa1\x99\x96\x76\x21\x89\x45\x38\x46\xe0\xde\x92\x34\xdb\x62\x5e\x97\xe0\xde\xb1\xe4\x35\x12\xa8\x10\x2b\xe5\xf7\x0d\x62\xe8\x7f\x94\x52\xc4\x78\x53\x99\x6a\x94\xac\xc1\x0d\x17\xf3\xe5\x07\x45\xee\x6d\x84\xfd\x85\xea\xa5\x26\xfd\x05\x28\x45\x03\xdd\xb3\xac\x98\x12\x75\x51\xba\x86\x80\xde\xa9\xeb\xbf\x53\x86\x38\x4f\x08\xd6\xcf\x04\xfe\xe9\xcf\xcc\x70\xf1\xcc\xc3\xe3\x29\x48\x44\x9c\xf7\xe0\x17\xba\x43\xc1\x12\x1c\x57\xda\x43\xc1\x4c\x77\xae\x27\xbe\x50\x7b\x4c\xd5\x0c\x7f\x56\x3d\x3c\xf6\x75\x53\x8a\x8e\x2b\x2d\x35\x20\xfc\x5c\x32\xee\xb8\x6d\x29\x9c\x2f\xa7\xcb\x72\x7a\xbe\x3c\x7d\x9f\x8e\x5a\xb8\xf9\xf2\x98\x25\xe9\x20\x3c\xbd\x14\xa3\xe9\xb9\x8c\x0e\xaf\xe8\x8c\x07\x59\x9a\x96\x05\xc5\x9d\x46\xea\x36\x02\x0f\x88\x38\x50\x53\x03\x95\x76\xed\xc6\xbc\xdd\xdd\x6f\x72\x77\x1c\x94\x9a\x8a\xd0\x53\x49\x94\x37\xa3\xa1\x99\x0c\x76\x03\x5a\xa7\xa9\x9e\x86\x49\x03\xb6\x13\xee\x46\x1a\xb2\x57\x3b\xac\x22\x47\x45\x6d\x3b\x68\x77\x3a\x66\x75\x11\x65\xc9\x2b\x4d\x75\x55\x33\xaa\xea\xb9\xcc\x75\xbd\x65\x0e\x87\x22\xc8\xbc\x03\x12\xf2\x3e\x57\xf8\x0b\x4d\x13\x6a\xbf\x5e\xef\x58\x96\xbf\xd0\x31\x5c\xd9\xb2\x1d\xac\x47\x7a\x15\xf2\x06\x6a\x71\x14\xec\xc8\x2e\x73\xfd\xf7\x5c\x9f\x09\x66\xdd\xa4\x1c\x75\xe2\xb8\x86\xd5\x60\x58\xd1\x97\x23\xdb\xb2\x02\xff\x2f\x10\x8c\xd5\xf0\x37\x9a\x3b\x18\xed\x1c\xaf\x9f\x2f\x48\x68\x51\x9e\x61\xaf\xe2\x5b\x3a\xe1\x14\xa2\x96\x65\x29\x6b\x26\x18\xca\x84\x54\x43\x6b\xc2\xe0\xc7\x08\x12\x65\x2a\xe3\xcd\x56\x7b\x55\x64\x25\x18\x26\x10\x51\x8a\xf0\xf3\x20\xc1\x23\x48\x2a\xa2\x43\x55\xa0\x36\xdd\xea\x53\xef\x29\x5c\x55\x65\x96\x5f\x05\xf6\xbf\x2d\x67\x57\x6c\x26\x38\x58\x89\x71\x8c\xc0\x4a\xa2\x3e\x07\xbd\x3e\x87\x75\x67\xea\x41\x47\x19\x0e\xe6\x95\xee\xe8\xba\x67\xf4\xe3\x68\x5c\xdc\xd9\x5d\x86\x6e\x52\xab\xed\xd2\xd7\x9c\x0e\xee\x4b\xb3\x51\xd3\x71\x97\xc0\xdb\xf7\x41\x8b\x41\xc9\xa1\x19\x9f\x80\x52\x8e\x3a\xc3\xb2\x1d\xd3\x6a\xa8\xf9\x42\xfc\x7e\x50\xee\x61\xf4\xc7\xe3\x61\x4a\x7f\x3e\x1b\xf7\x29\x5e\x22\x10\xcf\xf7\xd2\xa5\xb1\xcf\xf4\x95\xf9\x10\xa9\x1f\xd7\xb5\x8f\x13\x55\xde\xf8\x56\xf9\x15\x00\x00\xff\xff\x59\xeb\x45\xd6\x49\x0e\x00\x00")

func templatesRsTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesRsTpl,
		"templates/rs.tpl",
	)
}

func templatesRsTpl() (*asset, error) {
	bytes, err := templatesRsTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/rs.tpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesSchemaTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\xc1\x6a\xc3\x30\x10\x44\xef\xfa\x8a\xc1\xb4\x60\x07\x22\xdf\x0b\x3d\xe7\x16\x4a\xdc\x1f\x90\x95\x8d\x6c\x6c\x4b\x8e\xb4\xa2\x14\xd7\xff\x5e\x24\x9b\x06\xda\x1e\x77\xf6\x49\x33\xb3\xb3\xd2\x83\x32\x84\x65\x81\x6c\x74\x47\x93\x3a\xab\x89\xb0\xae\x42\xf4\xd3\xec\x3c\xa3\x14\x00\x50\x98\x9e\xbb\xd8\x4a\xed\xa6\xba\x73\x31\x90\x6a\x83\x1b\x23\x53\x3d\x0f\xa6\x6e\x55\xa0\x22\x73\xcb\x72\x84\x57\xd6\x10\xe4\xbb\x6a\x47\x0a\xe9\xab\x6d\x01\x79\x72\x6f\x83\xc1\x17\x66\xdf\x5b\xbe\xa1\x78\xbe\x17\x8f\xf5\x11\x64\xaf\x69\xac\x84\xe0\xcf\x99\xb0\xc5\x41\x60\x1f\x35\x63\xc9\xd8\xb5\x57\x23\x69\x4e\x62\x6f\x4d\x96\x0e\xc9\x7c\xcf\x2e\x56\x21\xfe\x4d\x90\xc4\xa7\xbb\x4d\xd5\x5e\x5e\x21\x73\xc7\x3f\x39\x32\xe4\x43\x26\x4e\xee\x37\x13\x2e\x4d\xa6\xc4\x2d\x5a\x8d\x32\xe0\xb0\x79\x56\xa9\x5a\x7a\xb6\xae\x65\x85\xc3\xcf\xb0\x27\xf6\xc4\xd1\xdb\x47\xfd\xa6\x73\x9e\xf7\x1b\xcb\x33\x7d\x5c\x9a\x32\x54\x62\x33\xdf\x0f\xf0\x1d\x00\x00\xff\xff\xd8\x84\xe1\x5f\x94\x01\x00\x00")

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
	"templates/rs.tpl":     templatesRsTpl,
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
		"rs.tpl":     &bintree{templatesRsTpl, map[string]*bintree{}},
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
