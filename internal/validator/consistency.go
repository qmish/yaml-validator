package validator

import (
	"fmt"

	"yaml-validator/internal/parser"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// CheckConsistency проверяет консистентность значений по заданным путям между файлами (8.5).
// files — пути к файлам, paths — dot-пути для сравнения (напр. "version", "metadata.annotations.app-version").
func CheckConsistency(files []string, paths []string) []pkg.Error {
	var errors []pkg.Error
	if len(files) < 2 || len(paths) == 0 {
		return errors
	}

	// Парсим первый документ каждого файла
	type fileVal struct {
		file  string
		root  *yaml.Node
		valid bool
	}
	parsed := make([]fileVal, len(files))
	for i, f := range files {
		nodes, err := parser.ParseFileMulti(f)
		if err != nil {
			errors = append(errors, pkg.Error{
				Type:    "ConsistencyParseError",
				Message: fmt.Sprintf("cannot parse %s for consistency check: %v", f, err),
			})
			continue
		}
		root := parser.GetRootMapping(nodes[0])
		parsed[i] = fileVal{file: f, root: root, valid: root != nil}
	}

	for _, path := range paths {
		fileToVal := make(map[string]string)
		for _, pv := range parsed {
			if !pv.valid || pv.root == nil {
				continue
			}
			val, ok := parser.GetValueAtPath(pv.root, path)
			if ok {
				fileToVal[pv.file] = val
			}
		}
		if len(fileToVal) < 2 {
			continue
		}
		// Сравниваем значения
		var firstFile, firstVal string
		for f, v := range fileToVal {
			if firstFile == "" {
				firstFile, firstVal = f, v
				continue
			}
			if v != firstVal {
				errors = append(errors, pkg.Error{
					Type:    "InconsistentValue",
					Message: fmt.Sprintf("inconsistent value for '%s': %s has %q, %s has %q", path, firstFile, firstVal, f, v),
					Path:    path,
				})
				break
			}
		}
	}
	return errors
}
