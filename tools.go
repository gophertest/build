package build

import (
	"bufio"
	"bytes"
	"fmt"
	gb "go/build"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Tools provides interfaces to build tools.
type Tools interface {
	Assembler
	Compiler
	Linker
	Packer
	BuildIDer
	Version() (string, error)
	BuildCtx() (gb.Context, error)
}

var (
	DebugLog bool = false
)

var (
	// DefaultTools uses tools provided by the current go runtime.
	DefaultTools Tools = &cmdTools{
		Go:        "go",
		Assembler: path.Join(gb.ToolDir, "asm"),
		Compiler:  path.Join(gb.ToolDir, "compile"),
		Linker:    path.Join(gb.ToolDir, "link"),
		Packer:    path.Join(gb.ToolDir, "pack"),
		BuildIDer: path.Join(gb.ToolDir, "buildid"),
	}
)

type cmdTools struct {
	mutex sync.Mutex

	Go            string
	GoArgs        []string
	Assembler     string
	AssemblerArgs []string
	Compiler      string
	CompilerArgs  []string
	Linker        string
	LinkerArgs    []string
	Packer        string
	PackerArgs    []string
	BuildIDer     string
	BuildIDerArgs []string

	version string
}

var envRegex = regexp.MustCompile(`([a-zA-Z0-9_]+)="(.*)"`)

func (ct *cmdTools) BuildCtx() (gb.Context, error) {
	ctx := gb.Default

	cmdArgs := append([]string(nil), ct.GoArgs...)
	cmdArgs = append(cmdArgs, "env")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := exec.Command(ct.Go, cmdArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		io.Copy(os.Stderr, stderr)
		return ctx, err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		values := envRegex.FindStringSubmatch(scanner.Text())
		if values == nil {
			continue
		}
		key := values[1]
		value := values[2]
		switch key {
		case "GOOS":
			ctx.GOOS = value
		case "GOARCH":
			ctx.GOARCH = value
		case "GOPATH":
			ctx.GOPATH = value
		case "GOROOT":
			ctx.GOROOT = value
		}
	}

	if scanner.Err() != nil {
		return ctx, scanner.Err()
	}

	return ctx, nil
}

func (ct *cmdTools) Version() (string, error) {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()
	if ct.version != "" {
		return ct.version, nil
	}
	cmdArgs := append([]string(nil), ct.GoArgs...)
	cmdArgs = append(cmdArgs, "version")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := exec.Command(ct.Go, cmdArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		io.Copy(os.Stderr, stderr)
		return "", err
	}
	ct.version = strings.TrimSpace(stdout.String())
	return ct.version, nil
}

func (ct *cmdTools) env(buildCtx gb.Context) []string {
	env := os.Environ()
	newEnv := make([]string, 0, len(env))
	for _, v := range env {
		switch {
		case strings.HasPrefix(v, "GOARCH"):
		case strings.HasPrefix(v, "GOOS"):
		case strings.HasPrefix(v, "GOROOT"):
		case strings.HasPrefix(v, "GOPATH"):
		case strings.HasPrefix(v, "CGO_ENABLED"):
		default:
			newEnv = append(newEnv, v)
		}
	}
	newEnv = append(newEnv,
		"GOARCH="+buildCtx.GOARCH,
		"GOOS="+buildCtx.GOOS,
		"GOROOT="+buildCtx.GOROOT,
		"GOPATH="+buildCtx.GOPATH,
	)
	if buildCtx.CgoEnabled {
		newEnv = append(newEnv, "CGO_ENABLED=1")
	} else {
		newEnv = append(newEnv, "CGO_ENABLED=0")
	}
	return newEnv
}

func (ct *cmdTools) Assemble(args AssembleArgs) error {
	cmdArgs := append([]string(nil), ct.AssemblerArgs...)
	if args.TrimPath != "" {
		cmdArgs = append(cmdArgs, "-trimpath", args.TrimPath)
	}
	if args.OutputFile != "" {
		cmdArgs = append(cmdArgs, "-o", args.OutputFile)
	}
	for _, v := range args.IncludeDirs {
		cmdArgs = append(cmdArgs, "-I", v)
	}
	for _, v := range args.Defines {
		cmdArgs = append(cmdArgs, "-D", v)
	}
	if args.GenSymABIs {
		cmdArgs = append(cmdArgs, "-gensymabis")
	}
	if args.Shared {
		cmdArgs = append(cmdArgs, "-shared")
	}
	if args.DynamicLink {
		cmdArgs = append(cmdArgs, "-dynlink")
	}
	for _, v := range args.Files {
		cmdArgs = append(cmdArgs, v)
	}
	if DebugLog {
		fmt.Printf("cd %s\n", args.WorkingDirectory)
		fmt.Printf("%s %s\n", ct.Assembler, strings.Join(cmdArgs, " "))
	}
	cmd := exec.Command(ct.Assembler, cmdArgs...)
	cmd.Env = ct.env(args.Context)
	cmd.Dir = args.WorkingDirectory
	cmd.Stdout = args.Stdout
	cmd.Stderr = args.Stderr
	return cmd.Run()
}

func (ct *cmdTools) Compile(args CompileArgs) error {
	cmdArgs := append([]string(nil), ct.CompilerArgs...)
	if args.TrimPath != "" {
		cmdArgs = append(cmdArgs, "-trimpath", args.TrimPath)
	}
	if args.OutputFile != "" {
		cmdArgs = append(cmdArgs, "-o", args.OutputFile)
	}
	if args.BuildID != "" {
		cmdArgs = append(cmdArgs, "-buildid", args.BuildID)
	}
	if args.DisableBoundsChecking {
		cmdArgs = append(cmdArgs, "-B")
	}
	if args.CompilingRuntimeLibrary {
		cmdArgs = append(cmdArgs, "-+")
	}
	if args.DisableOptimizations {
		cmdArgs = append(cmdArgs, "-N")
	}
	if args.RelativeImportPath != "" {
		cmdArgs = append(cmdArgs, "-D", args.RelativeImportPath)
	}
	for _, v := range args.IncludeDirs {
		cmdArgs = append(cmdArgs, "-I", v)
	}
	if args.Concurrency != 0 {
		cmdArgs = append(cmdArgs, "-D", strconv.Itoa(args.Concurrency))
	}
	if args.AsmHeaderFile != "" {
		cmdArgs = append(cmdArgs, "-asmhdr", args.AsmHeaderFile)
	}
	if args.Complete {
		cmdArgs = append(cmdArgs, "-complete")
	}
	if args.DynamicLink {
		cmdArgs = append(cmdArgs, "-dynlink")
	}
	if args.GoVersion != "" {
		cmdArgs = append(cmdArgs, "-goversion", args.GoVersion)
	}
	if args.HaltOnError {
		cmdArgs = append(cmdArgs, "-h")
	}
	if args.ImportConfigFile != "" {
		cmdArgs = append(cmdArgs, "-importcfg", args.ImportConfigFile)
	}
	for _, v := range args.ImportMap {
		cmdArgs = append(cmdArgs, "-importmap", v)
	}
	if args.InstallSuffix != "" {
		cmdArgs = append(cmdArgs, "-installsuffix", args.InstallSuffix)
	}
	if args.DisableInlining {
		cmdArgs = append(cmdArgs, "-l")
	}
	if args.LinkObjectOutputFile != "" {
		cmdArgs = append(cmdArgs, "-linkobj", args.LinkObjectOutputFile)
	}
	if args.MSan {
		cmdArgs = append(cmdArgs, "-msan")
	}
	if args.NoLocalImports {
		cmdArgs = append(cmdArgs, "-nolocalimports")
	}
	if args.PackageImportPath != "" {
		cmdArgs = append(cmdArgs, "-p", args.PackageImportPath)
	}
	if args.Pack {
		cmdArgs = append(cmdArgs, "-pack")
	}
	if args.Race {
		cmdArgs = append(cmdArgs, "-race")
	}
	if args.Shared {
		cmdArgs = append(cmdArgs, "-shared")
	}
	if args.SmallFrames {
		cmdArgs = append(cmdArgs, "-smallframes")
	}
	if args.CompilingStandardLibrary {
		cmdArgs = append(cmdArgs, "-std")
	}
	if args.SymABIsFile != "" {
		cmdArgs = append(cmdArgs, "-symabis", args.SymABIsFile)
	}
	for _, v := range args.Files {
		cmdArgs = append(cmdArgs, v)
	}
	if DebugLog {
		fmt.Printf("cd %s\n", args.WorkingDirectory)
		fmt.Printf("%s %s\n", ct.Compiler, strings.Join(cmdArgs, " "))
	}
	cmd := exec.Command(ct.Compiler, cmdArgs...)
	cmd.Env = ct.env(args.Context)
	cmd.Dir = args.WorkingDirectory
	cmd.Stdout = args.Stdout
	cmd.Stderr = args.Stderr
	return cmd.Run()
}

func (ct *cmdTools) Link(args LinkArgs) error {
	cmdArgs := append([]string(nil), ct.LinkerArgs...)
	if args.EntrySymbolName != "" {
		cmdArgs = append(cmdArgs, "-E", args.EntrySymbolName)
	}
	if args.HeaderType != "" {
		cmdArgs = append(cmdArgs, "-H", args.HeaderType)
	}
	if args.ELFDynamicLinker != "" {
		cmdArgs = append(cmdArgs, "-I", args.ELFDynamicLinker)
	}
	for _, v := range args.LibraryPaths {
		cmdArgs = append(cmdArgs, "-L", v)
	}
	for _, v := range args.StringDefines {
		cmdArgs = append(cmdArgs, "-X", v)
	}
	if args.BuildID != "" {
		cmdArgs = append(cmdArgs, "-buildid", args.BuildID)
	}
	if args.BuildMode != "" {
		cmdArgs = append(cmdArgs, "-buildmode", args.BuildMode)
	}
	if args.ExternalTar != "" {
		cmdArgs = append(cmdArgs, "-extar", args.ExternalTar)
	}
	if args.ExternalLinker != "" {
		cmdArgs = append(cmdArgs, "-extld", args.ExternalLinker)
	}
	if args.ExternalLinkerFlags != "" {
		cmdArgs = append(cmdArgs, "-extldflags", args.ExternalLinkerFlags)
	}
	if args.IgnoreVersionMismatch {
		cmdArgs = append(cmdArgs, "-f")
	}
	if args.DisableGoPackageDataChecks {
		cmdArgs = append(cmdArgs, "-g")
	}
	if args.HaltOnError {
		cmdArgs = append(cmdArgs, "-h")
	}
	if args.ImportConfigFile != "" {
		cmdArgs = append(cmdArgs, "-importcfg", args.ImportConfigFile)
	}
	if args.InstallSuffix != "" {
		cmdArgs = append(cmdArgs, "-installsuffix", args.InstallSuffix)
	}
	if args.FieldTrackingSymbol != "" {
		cmdArgs = append(cmdArgs, "-k", args.FieldTrackingSymbol)
	}
	if args.LibGCC != "" {
		cmdArgs = append(cmdArgs, "-libgcc", args.LibGCC)
	}
	if args.LinkMode != "" {
		cmdArgs = append(cmdArgs, "-linkmode", args.LinkMode)
	}
	if args.LinkShared {
		cmdArgs = append(cmdArgs, "-linkshared")
	}
	if args.MSan {
		cmdArgs = append(cmdArgs, "-msan")
	}
	if args.OutputFile != "" {
		cmdArgs = append(cmdArgs, "-o", args.OutputFile)
	}
	if args.PluginPath != "" {
		cmdArgs = append(cmdArgs, "-pluginpath", args.PluginPath)
	}
	if args.Race {
		cmdArgs = append(cmdArgs, "-race")
	}
	if args.TempDir != "" {
		cmdArgs = append(cmdArgs, "-tmpdir", args.TempDir)
	}
	if args.RejectUnsafePackages {
		cmdArgs = append(cmdArgs, "-u")
	}
	for _, v := range args.Files {
		cmdArgs = append(cmdArgs, v)
	}
	if DebugLog {
		fmt.Printf("cd %s\n", args.WorkingDirectory)
		fmt.Printf("%s %s\n", ct.Linker, strings.Join(cmdArgs, " "))
	}
	cmd := exec.Command(ct.Linker, cmdArgs...)
	cmd.Env = ct.env(args.Context)
	cmd.Dir = args.WorkingDirectory
	cmd.Stdout = args.Stdout
	cmd.Stderr = args.Stderr
	return cmd.Run()
}

func (ct *cmdTools) Pack(args PackArgs) error {
	cmdArgs := append([]string(nil), ct.PackerArgs...)
	op := ""
	switch args.Op {
	case AppendNew:
		op = "c"
	case Print:
		op = "p"
	case Append:
		op = "r"
	case List:
		op = "t"
	case Extract:
		op = "x"
	default:
		return fmt.Errorf("unknown pack operation %#v", args.Op)
	}
	cmdArgs = append(cmdArgs, op)
	cmdArgs = append(cmdArgs, args.ObjectFile)
	for _, v := range args.Names {
		cmdArgs = append(cmdArgs, v)
	}
	if DebugLog {
		fmt.Printf("cd %s\n", args.WorkingDirectory)
		fmt.Printf("%s %s\n", ct.Packer, strings.Join(cmdArgs, " "))
	}
	cmd := exec.Command(ct.Packer, cmdArgs...)
	cmd.Env = ct.env(args.Context)
	cmd.Dir = args.WorkingDirectory
	cmd.Stdout = args.Stdout
	cmd.Stderr = args.Stderr
	return cmd.Run()
}

func (ct *cmdTools) BuildID(args BuildIDArgs) (string, error) {
	cmdArgs := append([]string(nil), ct.BuildIDerArgs...)
	if args.Write {
		cmdArgs = append(cmdArgs, "-w")
	}
	cmdArgs = append(cmdArgs, args.ObjectFile)
	if DebugLog {
		fmt.Printf("cd %s\n", args.WorkingDirectory)
		fmt.Printf("%s %s\n", ct.BuildIDer, strings.Join(cmdArgs, " "))
	}
	stdout := &bytes.Buffer{}
	cmd := exec.Command(ct.BuildIDer, cmdArgs...)
	cmd.Env = ct.env(args.Context)
	cmd.Dir = args.WorkingDirectory
	cmd.Stdout = stdout
	cmd.Stderr = args.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}
