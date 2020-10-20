package utils

import (
	"fmt"

	"github.com/jhoonb/archivex"
)

func TarDirectory(folderName string) error {
	tar := new(archivex.TarFile)
	tarName := fmt.Sprintf("./%s.tar", folderName)
	if err := tar.Create(tarName); err != nil {
		return err
	}
	if err := tar.AddAll(folderName, true); err != nil {
		return err
	}
	if err := tar.Close(); err != nil {
		return err
	}
	return nil
}
