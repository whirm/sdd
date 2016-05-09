// args.go ---
//
// Filename: args.go
// Description:
// Author: Elric Milon
// Maintainer:
// Created: Thu May 19 22:06:47 2016 (+0200)

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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"strings"

	"github.com/docopt/docopt-go"
)

var includeHardLinks = false

var newHasher func() hash.Hash

var deletionHandler func(hashTree)

var dirs []string

func parseOptions() {

	usage := `Super Duper Deduper

Usage:
  sdd [--md5 | --sha1| --sha256] [--interactive | --dry-run | --auto | --link] [-H] DIR ...
  sdd -h | --help
  sdd -V | --version

Options:
  -h --help         Show this screen.
  -V --version      Show version
  -H --hardlinks    Consider hardlinks as different files
  -5 --md5          Hash using MD5
  -1 --sha1         Hash using SHA1 (default)
  -6 --sha256       Hash using SHA256
  -i --interactive  Ask for each duplicate group (default)
  -n --dry-run      Don't actually delete anything
  -a --auto         Automatically mark duplicates for deletion
  -l --link         Hardlink all duplicate files
`

	arguments, err := docopt.Parse(usage, nil, true, `0.0.0.1 "Anxaneta"`, false)

	if err != nil {
		panic(err)
	}
	switch {
	case arguments["--md5"]:
		newHasher = md5.New
	case arguments["--sha256"]:
		newHasher = sha256.New
	default:
		newHasher = sha1.New
	}

	switch {
	case arguments["--link"]:
		deletionHandler = hardLinkDupes
	case arguments["--auto"]:
		deletionHandler = deleteDupes
	case arguments["--dry-run"]:
		deletionHandler = simulateDeletion
	default:
		deletionHandler = askForDeletion
	}
	includeHardLinks = arguments["--hardlinks"].(bool)

	dirs = cleanDirs(arguments["DIR"].([]string))
}

func cleanDirs(dirs []string) []string {
	var absDirs []string
	for _, dir := range dirs {
		if isDir(dir) {
			abs, err := filepath.Abs(dir)
			if err != nil {
				panic(err)
			}
			absDirs = append(absDirs, abs)
		} else {
			fmt.Printf("%s is not a directory!\n", dir)
			os.Exit(2)
		}
	}
	var cleaned []string
	for i, dir := range absDirs {
		kick := false
		for j, dir2 := range dirs {
			if i == j {
				continue
			} else {
				if strings.HasPrefix(dir, dir2) {
					fmt.Printf("%s is a subdir of %s, ignoring it.\n", dir, dir2)
					kick = true
					break
				}
			}
		}
		if !kick {
			cleaned = append(cleaned, dir)
		}
	}
	return cleaned
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fi.Mode().IsDir()
}

//
// args.go ends here
