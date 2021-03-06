// Copyright © 2018 Chad Netzer <chad.netzer@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package hardlinkable

import "fmt"

const DefaultSearchThresh = 1
const DefaultMinFileSize = 1
const DefaultUseNewestLink = true
const DefaultStoreExistingLinkResults = true // Non-cli default
const DefaultStoreNewLinkResults = true      // Non-cli default
const DefaultShowExtendedRunStats = false    // Non-cli default
const DefaultShowRunStats = true             // Non-cli default

// Options is passed to the Run() func, and controls the operation of the
// hardlinkable algorithm, including what inode parameters much match for files
// to be compared for equality, what files and directories are included or
// excluded, and whether linking is actually enabled or not.
type Options struct {
	// SameName enabled ensures only files with matching filenames can be
	// linked
	SameName bool

	// IgnoreTime enabled allows files with different mtime values can be
	// linked
	IgnoreTime bool

	// IgnorePerm enabled allows files with different inode mode values
	// can be linked
	IgnorePerm bool

	// IgnoreOwner enabled allows files with different uid or gid can be
	// linked
	IgnoreOwner bool

	// IgnoreXAttr enabled allows files with different xattrs can be linked
	IgnoreXAttr bool

	// LinkingEnabled causes the Run to perform the linking step
	LinkingEnabled bool

	// MinFileSize controls the minimum size of files that are eligible to
	// be considered for linking.
	MinFileSize uint64

	// MaxFileSize controls the maximum size of files that are eligible to
	// be considered for linking.
	MaxFileSize uint64

	// DebugLevel controls the amount of debug information reported in the
	// results output, as well as debug logging.
	DebugLevel uint

	// UseNewestLink requests setting the inode to the mtime and uid/gid of
	// the more recent inode when files are linked.
	UseNewestLink bool

	// FileIncludes is a slice of regex expressions that control what
	// filenames will be considered for linking.  If given without any
	// FileExcludes, the walked files must match one of the includes.  If
	// FileExcludes are provided, the FileIncludes can override them.
	FileIncludes []string

	// FileExcludes is a slice of regex expressions that control what
	// filenames will be excluded from consideration for linking.
	FileExcludes []string

	// DirExcludes is a slice of regex expressions that control what
	// directories will be excluded from the file discovery walk.
	DirExcludes []string

	// StoreExistingLinkResults allows controlling whether to store
	// discovered existing links in Results. Command line option Verbosity
	// > 2 can override.
	StoreExistingLinkResults bool

	// StoreNewLinkResults allows controlling whether to store discovered
	// new hardlinkable pathnames in Results. Command line option Verbosity
	// > 1 can override.
	StoreNewLinkResults bool

	// ShowExtendedRunStats enabled displays additional Result stats
	// output.  Command line option Verbosity > 0 can override.
	ShowExtendedRunStats bool

	// ShowRunStats enabled displays Result stats output.
	ShowRunStats bool

	// IgnoreWalkErrors allows Run to continue when errors occur during the
	// walk phase, such as not having permission to walk a directory, or
	// being unable to read a file for comparision.
	IgnoreWalkErrors bool

	// IgnoreLinkErrors allows Run to continue when linking fails (or any
	// errors during the Link phase)
	IgnoreLinkErrors bool

	// CheckQuiescence enabled looks for signs of the filesystems changing
	// during walk.  Always enabled when LinkingEnabled is true.
	CheckQuiescence bool

	// SearchThresh determines the length that the lists of files with
	// equivalent inode hashes can grow to, before also enabling content
	// digests (which can drastically reduce the number of compared files
	// when there are many with the same hash, but differing content at the
	// start of the file).  Can be disabled with -1.  May save a small
	// amount of memory, but potentially at greatly increased runtime in
	// worst case scenarios with many, many files.
	SearchThresh int
}

// SetupOptions returns a Options struct with the defaults initialized and the
// given setup functions also applied.
func SetupOptions(args ...func(*Options)) Options {
	o := Options{
		SearchThresh:             DefaultSearchThresh,
		MinFileSize:              DefaultMinFileSize,
		UseNewestLink:            DefaultUseNewestLink,
		StoreExistingLinkResults: DefaultStoreExistingLinkResults,
		StoreNewLinkResults:      DefaultStoreNewLinkResults,
		ShowExtendedRunStats:     DefaultShowExtendedRunStats,
		ShowRunStats:             DefaultShowRunStats,
	}
	for _, fn := range args {
		fn(&o)
	}
	return o
}

// Not all the settable Options are represented as SetupOptions() funcs, just
// some common ones that user's might typically want to change from the
// default.  Option members can still be set directly, after the call to
// SetupOptions().

// SameName requires linked files to have equal filenames
func SameName(o *Options) {
	o.SameName = true
}

// IgnoreTime allows linked files to have unequal modification times
func IgnoreTime(o *Options) {
	o.IgnoreTime = true
}

// IgnorePerm allows linked files to have unequal mode bits
func IgnorePerm(o *Options) {
	o.IgnorePerm = true
}

// IgnoreOwner allows linked files to have unequal uid or gid
func IgnoreOwner(o *Options) {
	o.IgnoreOwner = true
}

// IgnoreXAttr allows linked files to have unequal xattrs
func IgnoreXAttr(o *Options) {
	o.IgnoreXAttr = true
}

// ContentOnly uses only file content to determine equality (not inode
// parameters like time, permission, ownership, etc.)
func ContentOnly(o *Options) {
	o.IgnoreTime = true
	o.IgnorePerm = true
	o.IgnoreOwner = true
	o.IgnoreXAttr = true
}

// LinkingEnabled allows Run() to actually perform linking of files
func LinkingEnabled(o *Options) {
	o.LinkingEnabled = true
}

// LinkingDisabled forbids Run() from actually linking the files
func LinkingDisabled(o *Options) {
	o.LinkingEnabled = false
}

// MinFileSize sets the minimum size of files that can be linked
func MinFileSize(size uint64) func(*Options) {
	return func(o *Options) {
		o.MinFileSize = size
	}
}

// MaxFileSize sets the maximum size of files that can be linked
func MaxFileSize(size uint64) func(*Options) {
	return func(o *Options) {
		o.MaxFileSize = size
	}
}

// DebugLevel sets the debugging level (1,2,or 3)
func DebugLevel(debugLevel uint) func(*Options) {
	return func(o *Options) {
		o.DebugLevel = debugLevel
	}
}

// ShowExtendedRunStats enabled prints more in OutputRunStats()
func ShowExtendedRunStats(o *Options) {
	o.ShowExtendedRunStats = true
}

// IgnoreWalkErrors allows the Run to continue during Walk phase errors (such
// as permission errors reading dirs or files)
func IgnoreWalkErrors(o *Options) {
	o.IgnoreWalkErrors = true
}

// IgnoreLinkErrors allows the Run to continue during Link phase errors
// (typically the actual linking itself)
func IgnoreLinkErrors(o *Options) {
	o.IgnoreLinkErrors = true
}

// CheckQuiescence enables quiescence checking which can detect changes to the
// filesystem during the file/directory walk.
func CheckQuiescence(o *Options) {
	o.CheckQuiescence = true
}

// Validate will ensure that contradictory Options aren't set, and that
// dependent Options are set.  An error will be returned if Options is invalid.
func (o *Options) Validate() error {
	if o.MaxFileSize > 0 && o.MaxFileSize < o.MinFileSize {
		return fmt.Errorf("MinFileSize (%v) cannot be larger than MaxFileSize (%v)",
			o.MinFileSize, o.MaxFileSize)
	}

	if o.ShowExtendedRunStats {
		o.ShowRunStats = true
	}

	if o.LinkingEnabled {
		o.CheckQuiescence = true
	}

	return nil
}
