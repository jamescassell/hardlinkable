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

package main

import (
	"hash/fnv"
	"os"
)

type Digest uint32

// Return a short digest of the first part of the given pathname, to help
// determine if two files are definitely not equivalent, without doing a full
// comparison.  Typically this will be used when a full file comparison will be
// performed anyway (incurring the IO overhead), and saving the digest to help
// quickly reduce the set of possibly equal inodes later (ie. reducing the
// length of the repeated linear searches).
func contentDigest(pathname string) (Digest, error) {
	const bufSize = 8192

	f, err := os.Open(pathname)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	buf := make([]byte, bufSize)
	_, err = f.Read(buf)
	if err != nil {
		return 0, err
	}

	Stats.computedDigest()

	hash := fnv.New32a()
	hash.Write(buf)
	return Digest(hash.Sum32()), nil
}