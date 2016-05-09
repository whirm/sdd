

### Super Duper Deduper

SDD is a multithreaded file deduplicator I've written while learning Go.

SDD is roughly 3x faster than fdupes with a cold cache on my i5 laptop. With a warm cache it's just marginally faster.

#### Disclaimer:

This software is still pre-alpha quality and is lacking many tests. Don't use it without backing up first!


```
Super Duper Deduper

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
  -a --auto         Automatically select duplicates for deletion
  -l --link         Hardlink all duplicate files
```
