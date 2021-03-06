package build

import (
	gb "go/build"
	"io"
)

// Assembler provides access to the `go tool asm` tool.
type Assembler interface {
	// Assemble runs the asm tool.
	Assemble(args AssembleArgs) error
}

// AssembleArgs passed to Assemble.
type AssembleArgs struct {
	Context          gb.Context
	WorkingDirectory string
	Stdout           io.Writer
	Stderr           io.Writer
	// Files to assemble.
	Files []string
	// TrimPath is "-trimpath string"
	TrimPath string
	// OutputFile is "-o string"
	OutputFile string
	// IncludeDirs is "-I string [-I string ...]"
	IncludeDirs []string
	// Defines is "-D string [-D string ...]"
	Defines []string
	// GenSymABIs is "-gensymabis"
	GenSymABIs bool
	// Shared is "-shared"
	Shared bool
	// DynamicLink is "-dynlink"
	DynamicLink bool
}

// Compiler provides access to the `go tool compile` tool.
type Compiler interface {
	// Compile runs the compile tool.
	Compile(args CompileArgs) error
}

// CompileArgs passed to Compile.
type CompileArgs struct {
	Context          gb.Context
	WorkingDirectory string
	Stdout           io.Writer
	Stderr           io.Writer
	// Files to compile.
	Files []string
	// TrimPath is "-trimpath string"
	TrimPath string
	// OutputFile is "-o string"
	OutputFile string
	// BuildID is "-buildid string"
	BuildID string
	// DisableBoundsChecking is "-B"
	DisableBoundsChecking bool
	// CompilingRuntimeLibrary is "-+"
	CompilingRuntimeLibrary bool
	// DisableOptimizations is "-N"
	DisableOptimizations bool
	// RelativeImportPath is "-D string"
	RelativeImportPath string
	// IncludeDirs is "-I string [-I string ...]"
	IncludeDirs []string
	// Concurrency is "-c=int"
	Concurrency int
	// AsmHeaderFile is "-asmhdr string"
	AsmHeaderFile string
	// Complete is "-complete"
	Complete bool
	// DynamicLink is "-dynlink"
	DynamicLink bool
	// GoVersion is "-goversion string"
	GoVersion string
	// HaltOnError is "-h"
	HaltOnError bool
	// ImportConfigFile is "-importcfg string"
	ImportConfigFile string
	// ImportMap is "-importmap string [-importmap string ...]"
	ImportMap []string
	// InstallSuffix is "-installsuffix string"
	InstallSuffix string
	// DisableInlining is "-l"
	DisableInlining bool
	// LinkObjectOutputFile is "-linkobj path"
	LinkObjectOutputFile string
	// MSan is "-msan"
	MSan bool
	// NoLocalImports is "-nolocalimports"
	NoLocalImports bool
	// PackageImportPath is "-p string"
	PackageImportPath string
	// Pack is "-pack"
	Pack bool
	// Race is "-race"
	Race bool
	// Shared is "-shared"
	Shared bool
	// SmallFrames is "-smallframes"
	SmallFrames bool
	// CompilingStandardLibrary is "-std"
	CompilingStandardLibrary bool
	// SymABIsFile is "-symabis string"
	SymABIsFile string
}

// Linker provides access to the `go tool link` tool.
type Linker interface {
	// Link runs the link tool.
	Link(args LinkArgs) error
}

// LinkArgs passed to Link.
type LinkArgs struct {
	Context          gb.Context
	WorkingDirectory string
	Stdout           io.Writer
	Stderr           io.Writer
	// Files to link.
	Files []string
	// EntrySymbolName is "-E string"
	EntrySymbolName string
	// HeaderType is "-H string"
	HeaderType string
	// ELFDynamicLinker is "-I string"
	ELFDynamicLinker string
	// LibraryPaths is "-L string [-L string ...]"
	LibraryPaths []string
	// StringDefines is "-X string [-X string ...]"
	StringDefines []string
	// BuildID is "-buildid string"
	BuildID string
	// BuildMode is "-buildmode string"
	BuildMode string
	// ExternalTar is "-extar string"
	ExternalTar string
	// ExternalLinker is "-extld string"
	ExternalLinker string
	// ExternalLinkerFlags is "-extldflags string"
	ExternalLinkerFlags string
	// IgnoreVersionMismatch is "-f"
	IgnoreVersionMismatch bool
	// DisableGoPackageDataChecks is "-g"
	DisableGoPackageDataChecks bool
	// HaltOnError is "-h"
	HaltOnError bool
	// ImportConfigFile is "-importcfg string"
	ImportConfigFile string
	// InstallSuffix is "-installsuffix string"
	InstallSuffix string
	// FieldTrackingSymbol is "-k string"
	FieldTrackingSymbol string
	// LibGCC is "-libgcc string"
	LibGCC string
	// LinkMode is "-linkmode string"
	LinkMode string
	// LinkShared is "-linkshared"
	LinkShared bool
	// MSan is "-msan"
	MSan bool
	// OutputFile is "-o string"
	OutputFile string
	// PluginPath is "-pluginpath string"
	PluginPath string
	// Race is "-race"
	Race bool
	// TempDir is "-tmpdir string"
	TempDir string
	// RejectUnsafePackages is "-u"
	RejectUnsafePackages bool
}

// Packer provides access to the `go tool pack` tool.
type Packer interface {
	// Pack runs the link pack.
	Pack(args PackArgs) error
}

// PackOp is the operation to perform on the object file.
type PackOp int

const (
	_ = iota
	// AppendNew passes the operation "c"
	AppendNew = iota
	// Print passes the operation "p"
	Print = iota
	// Append passes the operation "r"
	Append = iota
	// List passes the operation "t"
	List = iota
	// Extract passes the operation "x"
	Extract = iota
)

// PackArgs passed to Pack.
type PackArgs struct {
	Context          gb.Context
	WorkingDirectory string
	Stdout           io.Writer
	Stderr           io.Writer
	// Op the operation to perform on the object file.
	Op PackOp
	// ObjectFile to operate on
	ObjectFile string
	// Names to pass to the operation.
	Names []string
}

// BuildIDer can read and write BuildID
type BuildIDer interface {
	// BuildID either reads or write the BuildID
	BuildID(args BuildIDArgs) (string, error)
}

// BuildIDArgs passed to BuildID
type BuildIDArgs struct {
	Context          gb.Context
	WorkingDirectory string
	Stderr           io.Writer

	// ObjectFile to read or write BuildID
	ObjectFile string
	// Write is "-w"
	Write bool
}
