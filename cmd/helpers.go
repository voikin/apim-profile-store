package main

import (
	"fmt"
	"runtime/debug"
)

func getSwaggerVersion() string {
	const moduleName = "github.com/voikin/apim-proto/gen/go"

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "latest"
	}

	for _, dep := range info.Deps {
		if dep.Path == moduleName {
			return dep.Version
		}
	}

	return "unknown"
}

func getSwaggerURL() string {
	version := getSwaggerVersion()
	return fmt.Sprintf(
		"https://raw.githubusercontent.com/voikin/apim-proto/gen/go/%s/gen/openapi/apim_profile_store/v1/apim_profile_store.swagger.json", //nolint:lll // single string
		version,
	)
}
