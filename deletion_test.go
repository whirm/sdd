// deletion_test.go ---

//
// Filename: deletion_test.go
// Description:
// Author: Elric Milon
// Maintainer:
// Created: Mon May 30 20:53:23 2016 (+0200)

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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func genericDeleterTestBase(t *testing.T, confirmer func() bool) map[string]bool {
	dir, err := ioutil.TempDir("", "sddtest")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up

	content := []byte("banana")

	var files []string

	for _, file := range []string{"a", "b", "c"} {
		tmpfn := filepath.Join(dir, file)
		if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
			log.Fatal(err)
		}
		files = append(files, tmpfn)
	}

	tree := make(hashTree)

	for i, file := range []string{"a", "b", "c"} {
		tree.add(fileHash("yay"), fileSize(42), filesystemID(1), inode(i), file)
	}

	genericDeleter(tree, func(f fileSlice, d []string) []string { return files[1:] }, confirmer)

	deletedFiles := map[string]bool{}

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Println("DELETED", file, err)
			deletedFiles[file[len(file)-1:]] = true
		}
	}
	return deletedFiles
}

func TestGenericDeleterOk(t *testing.T) {
	deletedFiles := genericDeleterTestBase(t, func() bool { return true })

	// "a" has been kept
	if deletedFiles["a"] {
		t.Error("'a' Should not have been deleted")
	}

	// The rest haven't
	for _, file := range []string{"b", "c"} {
		if !deletedFiles[file] {
			t.Error(file, "Should have been deleted")
		}
	}
}

func TestGenericDeleterCancel(t *testing.T) {
	deletedFiles := genericDeleterTestBase(t, func() bool { return false })
	// No file has been deleted
	if len(deletedFiles) != 0 {
		t.Error("Some files where unexpectedly deleted", deletedFiles)
	}
}

func TestAutomaticSelector(t *testing.T) {
	tree := make(hashTree)

	for i, file := range []string{"a", "b", "c"} {
		tree.add(fileHash("yay"), fileSize(42), filesystemID(1), inode(i), file)
	}
	tree.add(fileHash("nope"), fileSize(42), filesystemID(1), inode(42), "X")

	// files := fileSlice{"a", "b", "c"}
	var toDelete []string
	toDelete = automaticFileSelector(fileSlice{"a", "b", "c"}, toDelete)
	// The first element "a" should be saved and "b" and "c" should be selected for deletion
	if len(toDelete) != 2 {
		t.Error("Expecting toDelete to contain 2 elements, but instead got ", toDelete)
	}
	toDelete = automaticFileSelector(fileSlice{"X", "Y", "Z"}, toDelete)
	// The elements passed to toDelete should be kept
	if len(toDelete) != 4 {
		t.Error("Expecting toDelete to contain 4 elements, but instead got ", toDelete)
	}
	expectedToDelete := fileSlice{"b", "c", "Y", "Z"}
	for i, file := range toDelete {
		if expectedToDelete[i] != file {
			t.Error("Expecting toDelete to be ", expectedToDelete, " but got ", toDelete, " instead.")
		}
	}
}

//
// deletion_test.go ends here
