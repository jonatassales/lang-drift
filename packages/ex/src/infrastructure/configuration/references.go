package configuration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
)

// CheckReferences checks the local paths used by diagnostics and hooks. It
// deliberately does not parse native files: broken native configuration must
// not prevent observing a previously recorded workload or running a declared
// repair hook. Up, down, and validate retain their full snapshot validation.
func (Files) CheckReferences(owner project.Project, descriptor microruntime.Descriptor) error {
	directory, err := runtimeDirectory(owner, descriptor.ID)
	if err != nil {
		return err
	}
	if err := requireComponents(owner.Root, descriptor); err != nil {
		return err
	}
	if err := requireHooks(directory, descriptor); err != nil {
		return err
	}
	if err := requireDirectory(directory, descriptor.TerraformDirectory); err != nil {
		return fmt.Errorf("terraform: %w", err)
	}
	preparation, err := files.Existing(directory, descriptor.TerraformDirectory.String())
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(preparation)
	if err != nil {
		return err
	}
	found := false
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".tf") || strings.HasSuffix(entry.Name(), ".tf.json")) {
			found = true
		}
	}
	if !found {
		return errors.New("terraform directory contains no .tf or .tf.json configuration")
	}
	for _, unit := range descriptor.UnitNames() {
		path, err := files.Existing(directory, descriptor.Units[unit].Manifests.String())
		if err != nil {
			return fmt.Errorf("unit %s: %w", unit, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !info.IsDir() && (!info.Mode().IsRegular() || !manifestExtensions[filepath.Ext(path)]) {
			return fmt.Errorf("unit %s manifests must be a directory or a JSON or YAML file", unit)
		}
	}

	return nil
}

// ResolveHook returns one trusted extension file reached without symlinks.
// It performs no process execution and does not grant a sandbox boundary.
func (Files) ResolveHook(owner project.Project, runtime microruntime.Identifier, script microruntime.RelativePath) (application.HookScript, error) {
	directory, err := runtimeDirectory(owner, runtime)
	if err != nil {
		return application.HookScript{}, err
	}
	path, err := files.Existing(directory, script.String())
	if err != nil {
		return application.HookScript{}, err
	}
	if filepath.Ext(path) != ".sh" {
		return application.HookScript{}, errors.New("hooks must be .sh files")
	}
	info, err := os.Stat(path)
	if err != nil {
		return application.HookScript{}, err
	}
	if !info.Mode().IsRegular() {
		return application.HookScript{}, errors.New("hook must be a regular file")
	}

	return application.HookScript{Path: path, Directory: directory}, nil
}
