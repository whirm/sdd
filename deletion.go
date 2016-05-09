// deletion.go ---

//
// Filename: deletion.go
// Description:
// Author: Elric Milon
// Maintainer:
// Created: Sun May 22 10:30:56 2016 (+0200)

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
	"math"
	"os"
	"sort"
	"strings"
)

type fileSelector func(fileSlice, []string) []string

// Generates a list of all duplicate files and ask user for confirmation.
func deleteDupes(hashes hashTree) {
	genericDeleter(hashes, automaticFileSelector, interactiveConfirmRequester)
}

func askForDeletion(hashes hashTree) {
	genericDeleter(hashes, interactiveFileSelector, interactiveConfirmRequester)
}

func automaticFileSelector(files fileSlice, toDelete []string) []string {
	action := keep
	for _, file := range files {
		fmt.Printf(" %v:   %v\n", action, file)
		if action == delete {
			toDelete = append(toDelete, file)
		}
		action = delete
	}
	return toDelete
}

func interactiveFileSelector(files fileSlice, toDelete []string) []string {
	var kept = math.MaxUint32
	for i, file := range files {
		fmt.Printf("%6d: %v\n", i+1, file)
	}
	for kept > len(files) {
		fmt.Printf("Preserve file: (1 - %d, 0 for all): ", len(files))
		_, err := fmt.Scan(&kept)
		if err != nil {
			panic(err)
		}
	}
	if kept == 0 {
		fmt.Println("Keeping all")
	} else {
		fmt.Println("pressed", kept)
		kept--
		fileKept := files[kept]
		fmt.Printf("Keeping %s\n", fileKept)
		for i, file := range files {
			if i != kept {
				toDelete = append(toDelete, file)
			}
		}
	}
	return toDelete
}

func genericDeleter(hashes hashTree, fileSelector fileSelector, confirmRequester func() bool) {
	var toDelete []string
	for _, sizeMap := range hashes {
		files := sizeMap.getFileNames()
		sort.Sort(files)
		if len(files) > 1 {
			fmt.Println("")
			toDelete = fileSelector(files, toDelete)
		}
	}
	if len(toDelete) == 0 {
		fmt.Println("Nothing to be deleted.")
		return
	}

	fmt.Println("\nThe following files will be deleted:", strings.Join(toDelete, ", "))
	ok := confirmRequester()
	if ok {
		for _, file := range toDelete {
			fmt.Println("Deleting", file)
			err := os.Remove(file)
			if err != nil {
				fmt.Printf("Failed to remove %s, error was: %s\n", file, err)
			}
		}
	} else {
		fmt.Println("OK then! Aborting.")
	}
}

func simulateDeletion(hashes hashTree) {
	for _, sizeMap := range hashes {
		files := sizeMap.getFileNames()
		sort.Sort(files)
		if len(files) > 1 {
			fmt.Println("")
			action := keep
			for _, file := range files {
				fmt.Printf("%6v: %v\n", action, file)
				action = delete
			}
		}
	}
	fmt.Println("\nNothing has been deleted, this was a dry run")
}

func hardLinkDupes(hashes hashTree) {
	for _, sizeMap := range hashes {
		files := sizeMap.getFileNames()
		sort.Sort(files)
		if len(files) > 1 {
			target := files[0]
			fmt.Printf("\nHard linking:\n %v\n To:\n", target)
			for _, file := range files[1:] {
				fmt.Printf(" %v\n", file)
				err := os.Rename(file, file+".bak")
				if err != nil {
					panic(err)
				}
				err = os.Link(target, file)
				if err != nil {
					panic(err)
				}
				err = os.Remove(file + ".bak")
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func interactiveConfirmRequester() bool {
	var ok string
	for ok != "y" && ok != "n" {
		fmt.Print("Is that OK? y/n: ")
		_, err := fmt.Scan(&ok)
		if err != nil {
			panic(err)
		}
	}
	return ok == "y"
}

//
// deletion.go ends here
