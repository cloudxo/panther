package mage

/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const swaggerGlob = "api/gateway/*/api.yml"

// Build contains targets for compiling source code.
type Build mg.Namespace

// Other mage commands can use this to invoke build commands directly.
var build = Build{}

// API Generate API source files from GraphQL + Swagger
func (b Build) API() {
	mg.Deps(b.generateSwaggerClients, b.generateWebTypescript)
}

func (b Build) generateSwaggerClients() error {
	specs, err := filepath.Glob(swaggerGlob)
	if err != nil {
		return fmt.Errorf("failed to glob %s: %v", swaggerGlob, err)
	}

	logger.Infof("build:api: generating Go SDK for %d APIs (%s)", len(specs), swaggerGlob)

	cmd := filepath.Join(setupDirectory, "swagger")
	if _, err = os.Stat(cmd); err != nil {
		return fmt.Errorf("%s not found (%v): run 'mage setup'", cmd, err)
	}

	// This import has to be fixed, see below
	clientImport := regexp.MustCompile(
		`"github.com/panther-labs/panther/api/gateway/[a-z]+/client/operations"`)

	for _, spec := range specs {
		dir := filepath.Dir(spec)
		client, models := filepath.Join(dir, "client"), filepath.Join(dir, "models")

		args := []string{"generate", "client", "-q", "-f", spec, "-c", client, "-m", models}
		if err := sh.Run(cmd, args...); err != nil {
			return fmt.Errorf("%s %s failed: %v", cmd, strings.Join(args, " "), err)
		}

		// TODO - delete unused models
		// If an API model is removed, "swagger generate" will leave the Go file in place.
		// We tried to remove generated files based on timestamp, but that had issues in Docker.
		// We tried removing the client/ and models/ every time, but mage itself depends on some of these.
		// For now, developers just need to manually remove unused swagger models.

		// There is a bug in "swagger generate" which can lead to incorrect import paths.
		// To reproduce: comment out this section, clone to /tmp and "mage build:api" - note the diffs.
		// The most reliable workaround has been to just rewrite the import ourselves.
		//
		// For example, in api/gateway/remediation/client/panther_remediation_client.go:
		//     import "github.com/panther-labs/panther/api/gateway/analysis/client/operations"
		// should be
		//     import "github.com/panther-labs/panther/api/gateway/remediation/client/operations"
		walk(client, func(path string, info os.FileInfo) {
			if info.IsDir() || filepath.Ext(path) != ".go" {
				return
			}

			body, err := ioutil.ReadFile(path)
			if err != nil {
				logger.Fatalf("failed to open %s: %v", path, err)
			}

			correctImport := fmt.Sprintf(
				`"github.com/panther-labs/panther/api/gateway/%s/client/operations"`,
				filepath.Base(filepath.Dir(filepath.Dir(path))))

			newBody := clientImport.ReplaceAll(body, []byte(correctImport))
			if err := ioutil.WriteFile(path, newBody, info.Mode()); err != nil {
				logger.Fatalf("failed to rewrite %s: %v", path, err)
			}
		})

		// Format generated files with our license header and import ordering.
		// "swagger generate client" can embed the header, but it's simpler to keep the whole repo
		// formatted the exact same way.
		fmtLicense(client, models)
		if err := gofmt(client, models); err != nil {
			logger.Warnf("gofmt %s %s failed: %v", client, models, err)
		}
	}
	return nil
}

func (b Build) generateWebTypescript() error {
	logger.Info("build:api: generating web typescript from graphql")
	if err := sh.Run("npm", "run", "graphql-codegen"); err != nil {
		return fmt.Errorf("graphql generation failed: %v", err)
	}
	fmtLicense("web/__generated__")
	if err := prettier("web/__generated__/*"); err != nil {
		logger.Warnf("prettier web/__generated__/ failed: %v", err)
	}
	return nil
}

// Lambda Compile Go Lambda function source
func (b Build) Lambda() {
	if err := b.lambda(); err != nil {
		logger.Fatal(err)
	}
}

// "go build" in parallel for each Lambda function.
//
// If you don't already have all go modules downloaded, this may fail because each goroutine will
// automatically modify the go.mod/go.sum files which will cause conflicts with itself.
//
// Run "go mod download" or "mage setup" before building to download the go modules.
// If you're adding a new module, run "go get ./..." before building to fetch the new module.
func (b Build) lambda() error {
	var packages []string
	walk("internal", func(path string, info os.FileInfo) {
		if info.IsDir() && strings.HasSuffix(path, "main") {
			packages = append(packages, path)
		}
	})

	logger.Infof("build:lambda: compiling %d Go Lambda functions (internal/.../main) using %s",
		len(packages), runtime.Version())

	for _, pkg := range packages {
		if err := buildLambdaPackage(pkg); err != nil {
			return err
		}
	}

	return nil
}

func buildLambdaPackage(pkg string) error {
	targetDir := filepath.Join("out", "bin", pkg)
	binary := filepath.Join(targetDir, "main")
	var buildEnv = map[string]string{"GOARCH": "amd64", "GOOS": "linux"}

	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return fmt.Errorf("failed to create %s directory: %v", targetDir, err)
	}
	if err := sh.RunWith(buildEnv, "go", "build", "-p", "1", "-ldflags", "-s -w", "-o", targetDir, "./"+pkg); err != nil {
		return fmt.Errorf("go build %s failed: %v", binary, err)
	}

	return nil
}

// Tools Compile devtools and opstools
func (b Build) Tools() {
	if err := b.tools(); err != nil {
		logger.Fatal(err)
	}
}

func (b Build) tools() error {
	// cross compile so tools can be copied to other machines easily
	buildEnvs := []map[string]string{
		// darwin:arm is not compatible
		{"GOOS": "darwin", "GOARCH": "amd64"},
		{"GOOS": "linux", "GOARCH": "amd64"},
		{"GOOS": "linux", "GOARCH": "arm"},
		{"GOOS": "windows", "GOARCH": "amd64"},
		{"GOOS": "windows", "GOARCH": "arm"},
	}

	var paths []string
	walk("cmd", func(path string, info os.FileInfo) {
		if !info.IsDir() && filepath.Base(path) == "main.go" {
			paths = append(paths, path)
		}
	})

	// Set global gitVersion, warn if not deploying a tagged release
	getGitVersion(false)

	for _, path := range paths {
		parts := strings.SplitN(path, `/`, 3)
		// E.g. "out/bin/cmd/devtools/" or "out/bin/cmd/opstools"
		outDir := filepath.Join("out", "bin", parts[0], parts[1])

		// used in tools to check/display which Panther version was compiled
		setVersionVar := fmt.Sprintf("-X 'main.version=%s'", gitVersion)

		logger.Infof("build:tools: compiling %s to %s with %d os/arch combinations",
			path, outDir, len(buildEnvs))
		for _, env := range buildEnvs {
			// E.g. "requeue-darwin-amd64"
			binaryName := filepath.Base(filepath.Dir(path)) + "-" + env["GOOS"] + "-" + env["GOARCH"]
			if env["GOOS"] == "windows" {
				binaryName += ".exe"
			}

			err := sh.RunWith(env, "go", "build",
				"-ldflags", "-s -w "+setVersionVar,
				"-o", filepath.Join(outDir, binaryName), "./"+path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Generate CloudFormation: deployments/dashboards.yml and out/deployments/
func (b Build) Cfn() {
	if err := b.cfn(); err != nil {
		logger.Fatal(err)
	}
}

func (b Build) cfn() error {
	if err := embedAPISpec(); err != nil {
		return err
	}

	return generateDashboards()
}
