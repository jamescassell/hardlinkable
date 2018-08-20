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
	"fmt"
	"os"
	"syscall"
)

type Linkable struct {
	FSDevs map[int64]FSDev
}

var MyLinkable *Linkable

func NewLinkable() *Linkable {
	var l Linkable
	l.FSDevs = make(map[int64]FSDev)
	return &l
}

func init() {
	MyLinkable = NewLinkable()
}

func (ln *Linkable) Dev(dev int64) FSDev {
	if fsdev, ok := ln.FSDevs[dev]; ok {
		return fsdev
	} else {
		fsdev = NewFSDev(dev)
		ln.FSDevs[dev] = fsdev
		return fsdev
	}
}

func (ln *Linkable) FindIdenticalFiles(pathname string) {
	fi, err := os.Lstat(pathname)
	if err != nil {
		os.Exit(2)
	}
	sysStat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		os.Exit(2)
	}
	//fmt.Printf("%+v\n", sysStat)
	fsdev := ln.Dev(int64(sysStat.Dev))
	fsdev.findIdenticalFiles(pathname, fi)
}

func Run(dirs []string) {
	var options *Options = &MyOptions
	c := MatchedPathnames(dirs)
	for pathname := range c {
		fi, err := os.Lstat(pathname)
		if err != nil {
			continue
		}
		if fi.Size() < options.MinFileSize {
			Stats.FoundFileTooSmall()
			continue
		}
		if options.MaxFileSize > 0 &&
			fi.Size() > options.MaxFileSize {
			Stats.FoundFileTooLarge()
			continue
		}
		// If the file hasn't been rejected by this
		// point, add it to the found count
		Stats.FoundFile()

		//fmt.Printf("%+v %s\n", stat, pathname)
		//fmt.Println(pathname)
		MyLinkable.FindIdenticalFiles(pathname)
	}
	//fmt.Printf("\n%+v\n", MyLinkable)
	fmt.Printf("\n%+v\n", Stats)
}