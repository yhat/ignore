package ignore

import (
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Ignorer struct {
	ignorers []ignorer
}

// Should the file be ignored? If the fullpath is a directory, it must contain
// a trailing slash. Use of a OS specific file separator is okay.
func (i Ignorer) Ignore(fullpath string) bool {

	fullpath = filepath.ToSlash(fullpath) // for windows

	f := file{fullpath, false}
	for _, ig := range i.ignorers {
		f.ignored = ig.Ignore(f)
	}
	return f.ignored
}

// Append an Ignorer to another.
func (i Ignorer) Append(ignorers ...Ignorer) Ignorer {
	islice := i.ignorers
	for _, ignorer := range ignorers {
		islice = append(islice, ignorer.ignorers...)
	}
	return Ignorer{islice}
}

type file struct {
	fullpath string
	ignored  bool // has this file already been ignored? (for '!')
}

type ignorer interface {
	Ignore(file) bool
}

// negator implements the '!' functionality, unignoring files if they match
// the underlying pattern.
type negator struct {
	p ignorer
}

func (n *negator) Ignore(f file) bool {
	if !f.ignored {
		return true
	}
	f.ignored = false
	return !n.p.Ignore(f)
}

// baseIgnorer ignores by the basename of the provided path
type baseIgnorer struct {
	basepath string
	pattern  string
}

func (i *baseIgnorer) Ignore(f file) bool {
	if f.ignored {
		return true
	}
	if !strings.HasPrefix(f.fullpath, i.basepath) {
		return false
	}
	b := path.Base(f.fullpath)
	ok, err := path.Match(i.pattern, b)
	return ok && (err == nil)
}

// pathIgnorer ignores by the fullpath of the provided file
type pathIgnorer struct {
	basepath string
	pattern  string
}

func (pi *pathIgnorer) Ignore(f file) bool {
	if f.ignored {
		return true
	}
	if !strings.HasPrefix(f.fullpath, pi.basepath) {
		return false
	}
	ok, err := path.Match(pi.pattern, f.fullpath)
	return ok && (err == nil)
}

// Parse the contents of an ignore file.
func Parse(ignore string) Ignorer {
	return ParseRel(ignore, "")
}

// Parse an ignore file relative to it's basepath. For instance the file
// 'foo/.ignore' would have a basepath of 'foo/'. Use of a OS specific
// separator is okay.
func ParseRel(ignore, basepath string) Ignorer {
	basepath = filepath.ToSlash(basepath)
	if basepath == "." {
		basepath = ""
	}
	lines := strings.Split(ignore, "\n")
	ignorers := []ignorer{}
	for _, line := range lines {
		line = trimTrailing(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		isNegator := line[0] == '!'
		if line[0] == '\\' && len(line) > 1 {
			if line[1] == '#' || line[1] == '!' {
				line = line[1:]
			}
		}

		ig := parseLine(line, basepath)
		if isNegator {
			ig = &negator{ig}
		}

		ignorers = append(ignorers, ig)
	}
	return Ignorer{ignorers}
}

// TODO: add '**' support
func parseLine(line, basepath string) ignorer {
	if 0 > strings.IndexRune(line, '/') {
		return &baseIgnorer{basepath, line}
	}
	return &pathIgnorer{basepath, line}
}

// 'Trailing spaces are ignored unless they are quoted with backslash ("\").'
// NOTE: Was only implemented in git 2.0.0
func trimTrailing(s string) string {
	for len(s) > 0 {
		r, size := utf8.DecodeLastRuneInString(s)
		switch r {
		case '\t', ' ':
		case '\\':
			return s[:len(s)-size]
		default:
			return s
		}
		s = s[:len(s)-size]
	}
	return s
}
