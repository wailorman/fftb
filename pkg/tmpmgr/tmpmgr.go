package tmpmgr

import (
	"os"
	"path/filepath"
)

type TmpManagerClient struct {
	location string
}

func New(location string) *TmpManagerClient {
	return &TmpManagerClient{location}
}

func (i *TmpManagerClient) Create(name string, options ...func(*directory)) (string, error) {
	dir := &directory{parentPath: i.location, name: name}

	for _, option := range options {
		option(dir)
	}

	err := dir.create()

	if err != nil {
		return "", err
	}

	return filepath.Join(dir.parentPath, dir.name), nil
}

func (i *TmpManagerClient) Destroy(name string) error {
	return os.RemoveAll(filepath.Join(i.location, name))
}

type directory struct {
	parentPath string
	name       string
}

func (d *directory) create() error {
	return os.MkdirAll(filepath.Join(d.parentPath, d.name), os.ModePerm)
}
