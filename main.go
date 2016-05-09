// main.go ---
//
// Filename: main.go
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
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
)

func main() {
	parseOptions()
	// First, walk through all dirs and group files by syze, file system and inode number
	// Stage 1: group files by size, filesystem and inode number
	sizes := processDirs(dirs)

	// Stage 2: find duplicates
	fmt.Println("Looking for duplicates")
	hashes := findDupes(sizes)

	// fmt.Println("STAGE 3: deleting")
	deletionHandler(hashes)
}

func findDupes(sizes sizeMap) hashTree {
	hashes := make(hashTree)
	workers := runtime.NumCPU()
	var mergerLock = &sync.Mutex{}
	var wg sync.WaitGroup

	wg.Add(workers)

	jobs := make(chan hashingJob, workers*10)
	results := make(chan hashingJob, workers*10)

	// Start a bunch of workers
	for w := 1; w <= workers; w++ {
		go hashWorker(jobs, results, &wg)
	}

	go hashPossibleDupes(sizes, jobs, &wg)

	go mergeResults(hashes, results, mergerLock)

	// Wait for all the workers to be done and close results channel
	wg.Wait()
	close(results)
	mergerLock.Lock()

	return hashes
}

// mergeResults Merges all hashing job results into a single map
func mergeResults(hashes hashTree, results <-chan hashingJob, lock *sync.Mutex) {
	lock.Lock()
	defer lock.Unlock()

	for job := range results {
		if includeHardLinks {
			for _, file := range job.files {
				hashes.add(job.hash, job.size, job.fsID, job.inode, file)
			}
		} else {
			hashes.add(job.hash, job.size, job.fsID, job.inode, job.files[0])
		}
	}
}
func processDirs(dirs []string) sizeMap {
	sizes := make(sizeMap)
	for c, dir := range dirs {
		fmt.Printf("Scanning %s (%d/%d)\n", dir, c+1, len(dirs))
		groupBySizes(sizes, dir)

	}
	return sizes
}

func groupBySizes(sizes sizeMap, dir string) sizeMap {
	err := filepath.Walk(dir, makeSizeWalker(sizes))
	if err != nil {
		panic(err)
	}
	return sizes
}

func hashWorker(jobs <-chan hashingJob, results chan<- hashingJob, wg *sync.WaitGroup) {
	for job := range jobs {
		job.hash = hashFile(job.files[0])
		results <- job
	}
	wg.Done()
}

// hashFile Hashes a file and returns it as a string
func hashFile(path string) fileHash {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	h := newHasher()
	_, err = io.Copy(h, f)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
	return fileHash(fmt.Sprintf("%x", h.Sum(nil)))

}

func makeSizeWalker(sizemap sizeMap) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {

		if !info.Mode().IsRegular() || err != nil {
			return nil
		}
		var stat syscall.Stat_t
		if err := syscall.Stat(path, &stat); err != nil {
			fmt.Println("Oopsie while scanning ", path)
			panic(err)
		}
		sizemap.add(fileSize(info.Size()), filesystemID(stat.Dev), inode(stat.Ino), path)
		return nil
	}
}

// hashPossibleDupes creates hashing jobs for all files that could be duplicates
func hashPossibleDupes(sizes sizeMap, jobs chan<- hashingJob, wg *sync.WaitGroup) {
	nullHash := fileHash("")
	for size, fsMap := range sizes {
		if fsMap.leavesCount() > 1 {
			for fsID, inodeMap := range fsMap {
				for inode, files := range inodeMap {
					jobs <- hashingJob{hash: nullHash, size: size, fsID: fsID, inode: inode, files: files}
				}
			}
		}
	}
	close(jobs)
}

//
// main.go ends here
