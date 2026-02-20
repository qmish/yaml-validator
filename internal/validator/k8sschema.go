package validator

import (
	"fmt"
	"os"

	"yaml-validator/internal/config"
	"yaml-validator/pkg"

	kubeconform "github.com/yannh/kubeconform/pkg/validator"
)

// CheckK8sSchema проверяет манифесты по OpenAPI/JSON Schema Kubernetes (все ресурсы, типы полей).
// Использует kubeconform. Включается через rules.k8s_schema.enabled в конфиге.
func CheckK8sSchema(filename string, opts config.K8sSchemaOptions) []pkg.Error {
	if !opts.Enabled {
		return nil
	}

	version := opts.Version
	if version == "" {
		version = "master"
	}

	kOpts := kubeconform.Opts{
		KubernetesVersion:    version,
		Strict:               opts.Strict,
		IgnoreMissingSchemas: opts.IgnoreMissing,
		Cache:                opts.CacheDir,
	}

	v, err := kubeconform.New(nil, kOpts)
	if err != nil {
		return []pkg.Error{{
			Type:    "K8sSchemaInit",
			Message: fmt.Sprintf("failed to init K8s schema validator: %v", err),
		}}
	}

	f, err := os.Open(filename)
	if err != nil {
		return []pkg.Error{{
			Type:    "K8sSchemaRead",
			Message: err.Error(),
		}}
	}
	defer f.Close()

	results := v.Validate(filename, f)

	var errors []pkg.Error
	for _, res := range results {
		switch res.Status {
		case kubeconform.Invalid:
			for _, ve := range res.ValidationErrors {
				errors = append(errors, pkg.Error{
					Type:    "K8sSchema",
					Message: ve.Msg,
					Path:    ve.Path,
				})
			}
			if res.Err != nil && len(res.ValidationErrors) == 0 {
				errors = append(errors, pkg.Error{
					Type:    "K8sSchema",
					Message: res.Err.Error(),
				})
			}
		case kubeconform.Error:
			if res.Err != nil {
				errors = append(errors, pkg.Error{
					Type:    "K8sSchema",
					Message: res.Err.Error(),
				})
			}
		}
	}

	return errors
}
