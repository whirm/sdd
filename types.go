// types.go ---
//
// Filename: types.go
// Description:
// Author: Elric Milon
// Maintainer:
// Created: Sun May  8 22:08:32 2016 (+0200)

// Commentary:
//
//
//
//

// Change Log:
//
//
//
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with GNU Emacs.  If not, see <http://www.gnu.org/licenses/>.
//
//

// Code:

package main

import (
	"fmt"
	"strings"
)

const keep = "Keep"
const delete = "Delete"

type fileSize int64
type filesystemID uint64
type inode uint64

type fileHash string

type fileSlice []string
type inodeMap map[inode]fileSlice
type fileSystemMap map[filesystemID]inodeMap
type sizeMap map[fileSize]fileSystemMap

type hashTree map[fileHash]fileSystemMap

type hashingJob struct {
	hash  fileHash
	size  fileSize
	fsID  filesystemID
	inode inode
	files fileSlice
}

func (files fileSlice) Len() int           { return len(files) }
func (files fileSlice) Swap(i, j int)      { files[i], files[j] = files[j], files[i] }
func (files fileSlice) Less(i, j int) bool { return files[i] < files[j] }

// len counts the number of leaves contained
func (inodes inodeMap) leavesCount() int {
	c := 0
	for _, files := range inodes {
		c += len(files)
	}
	return c
}

func (filesystems fileSystemMap) leavesCount() int {
	c := 0
	for _, inodes := range filesystems {
		if includeHardLinks {
			c += inodes.leavesCount()
		} else {
			c += len(inodes)
		}
	}
	return c
}

func (sizes sizeMap) leavesCount() int {
	c := 0
	for _, fs := range sizes {
		c += fs.leavesCount()
	}
	return c
}

func (hashes hashTree) leavesCount() int {
	c := 0
	for _, sizes := range hashes {
		c += sizes.leavesCount()
	}
	return c
}

func (inodes inodeMap) add(inode inode, name string) {
	inodes[inode] = append(inodes[inode], name)
}

func (filesystems fileSystemMap) add(fs filesystemID, inode inode, name string) {
	if filesystems[fs] == nil {
		filesystems[fs] = make(inodeMap)
	}
	filesystems[fs].add(inode, name)
}

func (sizes sizeMap) add(size fileSize, fs filesystemID, inode inode, name string) {
	if sizes[size] == nil {
		sizes[size] = make(fileSystemMap)
	}
	sizes[size].add(fs, inode, name)
}

func (hashes hashTree) add(hash fileHash, size fileSize, fs filesystemID, inode inode, name string) {
	if hashes[hash] == nil {
		hashes[hash] = make(fileSystemMap)
	}
	hashes[hash].add(fs, inode, name)
}

func (inodes inodeMap) getFileNames() fileSlice {
	var names = make(fileSlice, 0)
	for _, files := range inodes {
		names = append(names, files...)
	}
	return names
}

func (filesystems fileSystemMap) getFileNames() fileSlice {
	var names = make(fileSlice, 0)
	for _, inodes := range filesystems {
		names = append(names, inodes.getFileNames()...)
	}
	return names
}

func (sizes sizeMap) getFileNames() fileSlice {
	var names = make(fileSlice, 0)
	for _, fileSystems := range sizes {
		names = append(names, fileSystems.getFileNames()...)
	}
	return names
}

func (hashes hashTree) getFileNames() fileSlice {
	var names = make(fileSlice, 0)
	for _, sizes := range hashes {
		names = append(names, sizes.getFileNames()...)
	}
	return names
}

func (files fileSlice) prettyPrint(level int) {
	for i, file := range files {
		fmt.Printf("%s%d - %s\n", strings.Repeat(" ", level), i, file)
	}
}

func (inodes inodeMap) prettyPrint(level int) {
	for inode, files := range inodes {
		fmt.Printf("%sInode: %d\n", strings.Repeat(" ", level), inode)
		files.prettyPrint(level + 1)
	}
}

func (filesystems fileSystemMap) prettyPrint(level int) {
	for fs, inodes := range filesystems {
		fmt.Printf("%sFS: %d\n", strings.Repeat(" ", level), fs)
		inodes.prettyPrint(level + 1)
	}
}

func (sizes sizeMap) prettyPrint(level int) {
	for size, fileSystems := range sizes {
		fmt.Printf("%sSize: %d\n", strings.Repeat(" ", level), size)
		fileSystems.prettyPrint(level + 1)
	}
}
func (hashes hashTree) prettyPrint() {
	for hash, sizes := range hashes {
		fmt.Printf("%v:\n", hash)
		sizes.prettyPrint(1)
	}
}

//
// types.go ends here
