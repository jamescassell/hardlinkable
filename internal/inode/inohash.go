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

package inode

type Hash uint64

type InoHashes map[Hash]Set

// HashIno produces an equal hash for potentially equal files, based only on
// Inode metadata (size, time, etc.).  Content still has to be verified for
// equality (but unequal hashes indicate files that definitely need not be
// compared)
func HashIno(si StatInfo, ignoreTime, ignorePerm, ignoreOwner bool) Hash {
	h := uint64(si.Size)
	// The main requirement is that files that could be equal have equal
	// hashes.  It's less important if unequal files also have the same
	// hash value, since we will still compare the actual file content
	// later.
	if !ignoreTime {
		h ^= uint64(si.Mtim.UnixNano())
	}
	if !ignorePerm {
		h ^= uint64(si.Mode.Perm())
	}
	if !ignoreOwner {
		h ^= (uint64(si.Uid)<<32 | uint64(si.Gid))
	}
	return Hash(h)
}
