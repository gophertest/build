package build_test

import (
	"bytes"
	"fmt"
	gb "go/build"
	"os"
	"strings"
	"testing"

	"github.com/gophertest/build"
	"github.com/stretchr/testify/assert"
)

func TestAssembler(t *testing.T) {
	testCases := []struct {
		Args     build.AssembleArgs
		Expected string
	}{
		{
			build.AssembleArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout:      &bytes.Buffer{},
				TrimPath:    "tp",
				OutputFile:  "of",
				IncludeDirs: []string{"DirA", "DirB"},
				Defines:     []string{"A", "B"},
				GenSymABIs:  true,
				Shared:      true,
				DynamicLink: true,
				Files:       []string{"a", "b", "c"},
			},
			"-trimpath tp -o of -I DirA -I DirB -D A -D B -gensymabis -shared -dynlink a b c goos goarch go/path go/root 1",
		},
		{
			build.AssembleArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout: &bytes.Buffer{},
			},
			"goos goarch go/path go/root 1",
		},
	}
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		args := []string(nil)
		for i, v := range os.Args {
			if v == "--" {
				args = os.Args[i+1:]
			}
		}
		args = append(args, os.Getenv("GOOS"))
		args = append(args, os.Getenv("GOARCH"))
		args = append(args, os.Getenv("GOPATH"))
		args = append(args, os.Getenv("GOROOT"))
		args = append(args, os.Getenv("CGO_ENABLED"))
		fmt.Fprint(os.Stdout, strings.Join(args, " "))
		os.Exit(0)
	} else {
		os.Setenv("TEST_SUBPROCESS", "1")
		defer os.Setenv("TEST_SUBPROCESS", "")
		for c, tc := range testCases {
			tools := build.NewCmdTools()
			tools.Assembler = os.Args[0]
			tools.AssemblerArgs = []string{"-test.run=TestAssembler", "--"}
			err := tools.Assemble(tc.Args)
			assert.NoError(t, err)
			out := tc.Args.Stdout.(*bytes.Buffer)
			assert.Equalf(t, tc.Expected, out.String(), "failed with case %d", c)
		}
	}
}

func TestCompiler(t *testing.T) {
	testCases := []struct {
		Args     build.CompileArgs
		Expected string
	}{
		{
			build.CompileArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout:                   &bytes.Buffer{},
				TrimPath:                 "tp",
				OutputFile:               "of",
				DisableBoundsChecking:    true,
				CompilingRuntimeLibrary:  true,
				DisableOptimizations:     true,
				RelativeImportPath:       "rip",
				IncludeDirs:              []string{"includeDirA", "includeDirB"},
				Concurrency:              5,
				AsmHeaderFile:            "aho",
				Complete:                 true,
				DynamicLink:              true,
				GoVersion:                "",
				HaltOnError:              true,
				ImportConfigFile:         "icf",
				ImportMap:                []string{"importMapA", "importMapB"},
				InstallSuffix:            "",
				DisableInlining:          true,
				LinkObjectOutputFile:     "loof",
				MSan:                     true,
				NoLocalImports:           true,
				PackageImportPath:        "pip",
				Pack:                     true,
				Race:                     true,
				Shared:                   true,
				SmallFrames:              true,
				CompilingStandardLibrary: true,
				SymABIsFile:              "saf",
				Files:                    []string{"a", "b", "c"},
			},
			"-trimpath tp -o of -B -+ -N -D rip -I includeDirA -I includeDirB -D 5 -asmhdr aho -complete -dynlink -h -importcfg icf -importmap importMapA -importmap importMapB -l -linkobj loof -msan -nolocalimports -p pip -pack -race -shared -smallframes -std -symabis saf a b c goos goarch go/path go/root 1",
		},
		{
			build.CompileArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout: &bytes.Buffer{},
			},
			"goos goarch go/path go/root 1",
		},
	}
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		args := []string(nil)
		for i, v := range os.Args {
			if v == "--" {
				args = os.Args[i+1:]
			}
		}
		args = append(args, os.Getenv("GOOS"))
		args = append(args, os.Getenv("GOARCH"))
		args = append(args, os.Getenv("GOPATH"))
		args = append(args, os.Getenv("GOROOT"))
		args = append(args, os.Getenv("CGO_ENABLED"))
		fmt.Fprint(os.Stdout, strings.Join(args, " "))
		os.Exit(0)
	} else {
		os.Setenv("TEST_SUBPROCESS", "1")
		defer os.Setenv("TEST_SUBPROCESS", "")
		for c, tc := range testCases {
			tools := build.NewCmdTools()
			tools.Compiler = os.Args[0]
			tools.CompilerArgs = []string{"-test.run=TestCompiler", "--"}
			err := tools.Compile(tc.Args)
			assert.NoError(t, err)
			out := tc.Args.Stdout.(*bytes.Buffer)
			assert.Equalf(t, tc.Expected, out.String(), "failed with case %d", c)
		}
	}
}

func TestLinker(t *testing.T) {
	testCases := []struct {
		Args     build.LinkArgs
		Expected string
	}{
		{
			build.LinkArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout:                     &bytes.Buffer{},
				EntrySymbolName:            "esn",
				HeaderType:                 "ht",
				ELFDynamicLinker:           "edl",
				LibraryPaths:               []string{"lpa", "lpb"},
				StringDefines:              []string{"sda", "sdb"},
				BuildID:                    "bi",
				BuildMode:                  "bm",
				ExternalTar:                "et",
				ExternalLinker:             "el",
				ExternalLinkerFlags:        "elf",
				IgnoreVersionMismatch:      true,
				DisableGoPackageDataChecks: true,
				HaltOnError:                true,
				ImportConfigFile:           "icf",
				InstallSuffix:              "is",
				FieldTrackingSymbol:        "fts",
				LibGCC:                     "lgcc",
				LinkMode:                   "lm",
				LinkShared:                 true,
				MSan:                       true,
				OutputFile:                 "of",
				PluginPath:                 "pp",
				Race:                       true,
				TempDir:                    "td",
				RejectUnsafePackages:       true,
				Files:                      []string{"a", "b", "c"},
			},
			"-E esn -H ht -I edl -L lpa -L lpb -X sda -X sdb -buildid bi -buildmode bm -extar et -extld el -extldflags elf -f -g -h -importcfg icf -installsuffix is -k fts -libgcc lgcc -linkmode lm -linkshared -msan -o of -pluginpath pp -race -tmpdir td -u a b c goos goarch go/path go/root 1",
		},
		{
			build.LinkArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout: &bytes.Buffer{},
			},
			"goos goarch go/path go/root 1",
		},
	}
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		args := []string(nil)
		for i, v := range os.Args {
			if v == "--" {
				args = os.Args[i+1:]
			}
		}
		args = append(args, os.Getenv("GOOS"))
		args = append(args, os.Getenv("GOARCH"))
		args = append(args, os.Getenv("GOPATH"))
		args = append(args, os.Getenv("GOROOT"))
		args = append(args, os.Getenv("CGO_ENABLED"))
		fmt.Fprint(os.Stdout, strings.Join(args, " "))
		os.Exit(0)
	} else {
		os.Setenv("TEST_SUBPROCESS", "1")
		defer os.Setenv("TEST_SUBPROCESS", "")
		for c, tc := range testCases {
			tools := build.NewCmdTools()
			tools.Linker = os.Args[0]
			tools.LinkerArgs = []string{"-test.run=TestLinker", "--"}
			err := tools.Link(tc.Args)
			assert.NoError(t, err)
			out := tc.Args.Stdout.(*bytes.Buffer)
			assert.Equalf(t, tc.Expected, out.String(), "failed with case %d", c)
		}
	}
}

func TestPacker(t *testing.T) {
	testCases := []struct {
		Args     build.PackArgs
		Expected string
		Error    string
	}{
		{
			build.PackArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout:     &bytes.Buffer{},
				ObjectFile: "obj",
				Op:         build.Append,
				Names:      []string{"a", "b", "c"},
			},
			"r obj a b c goos goarch go/path go/root 1",
			"",
		},
		{
			build.PackArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout:     &bytes.Buffer{},
				ObjectFile: "obj",
				Op:         build.AppendNew,
				Names:      []string{"a", "b", "c"},
			},
			"c obj a b c goos goarch go/path go/root 1",
			"",
		},
		{
			build.PackArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout:     &bytes.Buffer{},
				ObjectFile: "obj",
				Op:         build.Extract,
				Names:      []string{"a", "b", "c"},
			},
			"x obj a b c goos goarch go/path go/root 1",
			"",
		},
		{
			build.PackArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout:     &bytes.Buffer{},
				ObjectFile: "obj",
				Op:         build.List,
				Names:      []string{"a", "b", "c"},
			},
			"t obj a b c goos goarch go/path go/root 1",
			"",
		},
		{
			build.PackArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout:     &bytes.Buffer{},
				ObjectFile: "obj",
				Op:         build.Print,
				Names:      []string{"a", "b", "c"},
			},
			"p obj a b c goos goarch go/path go/root 1",
			"",
		},
		{
			build.PackArgs{
				Context: gb.Context{
					GOOS:       "goos",
					GOARCH:     "goarch",
					GOPATH:     "go/path",
					GOROOT:     "go/root",
					CgoEnabled: true,
				},
				Stdout: &bytes.Buffer{},
			},
			"",
			"unknown pack operation 0",
		},
	}
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		args := []string(nil)
		for i, v := range os.Args {
			if v == "--" {
				args = os.Args[i+1:]
			}
		}
		args = append(args, os.Getenv("GOOS"))
		args = append(args, os.Getenv("GOARCH"))
		args = append(args, os.Getenv("GOPATH"))
		args = append(args, os.Getenv("GOROOT"))
		args = append(args, os.Getenv("CGO_ENABLED"))
		fmt.Fprint(os.Stdout, strings.Join(args, " "))
		os.Exit(0)
	} else {
		os.Setenv("TEST_SUBPROCESS", "1")
		defer os.Setenv("TEST_SUBPROCESS", "")
		for c, tc := range testCases {
			tools := build.NewCmdTools()
			tools.Packer = os.Args[0]
			tools.PackerArgs = []string{"-test.run=TestPacker", "--"}
			err := tools.Pack(tc.Args)
			if tc.Error == "" {
				assert.NoErrorf(t, err, "failed with case %d", c)
			} else {
				assert.EqualError(t, err, tc.Error, "failed with case %d", c)
			}
			out := tc.Args.Stdout.(*bytes.Buffer)
			assert.Equalf(t, tc.Expected, out.String(), "failed with case %d", c)
		}
	}
}

func TestVersion(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		fmt.Fprint(os.Stdout, "go version go99.99.99 linux/amd64")
		os.Exit(0)
		return
	}
	os.Setenv("TEST_SUBPROCESS", "1")
	defer os.Setenv("TEST_SUBPROCESS", "")
	tools := build.NewCmdTools()
	tools.Go = os.Args[0]
	tools.GoArgs = []string{"-test.run=TestVersion"}
	version, err := tools.Version()
	assert.NoError(t, err)
	assert.Equal(t, "go version go99.99.99 linux/amd64", version)
}
