package build

import (
	gb "go/build"
	"os/exec"
	"path"
	"strconv"
)

// Tools provides interfaces to build tools.
type Tools interface {
	Assembler
	Compiler
	Linker
}

var (
	// DefaultTools uses tools provided by the current go runtime.
	DefaultTools Tools = &cmdTools{
		Assembler: path.Join(gb.ToolDir, "asm"),
		Compiler:  path.Join(gb.ToolDir, "compile"),
		Linker:    path.Join(gb.ToolDir, "link"),
	}
)

type cmdTools struct {
	Assembler     string
	AssemblerArgs []string
	Compiler      string
	CompilerArgs  []string
	Linker        string
	LinkerArgs    []string
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
	cmd := exec.Command(ct.Assembler, cmdArgs...)
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
	cmd := exec.Command(ct.Compiler, cmdArgs...)
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
	cmd := exec.Command(ct.Linker, cmdArgs...)
	cmd.Dir = args.WorkingDirectory
	cmd.Stdout = args.Stdout
	cmd.Stderr = args.Stderr
	return cmd.Run()
}
