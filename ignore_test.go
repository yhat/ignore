package ignore

import "testing"

func TestTrimTrailing(t *testing.T) {
	tests := []struct {
		s, exp string
	}{
		{"", ""},
		{"    ", ""},
		{"\t", ""},
		{"foo ", "foo"},
		{`foo \`, "foo "},
		{"  foo", "  foo"},
	}

	for _, test := range tests {
		s := trimTrailing(test.s)
		if s != test.exp {
			t.Errorf("'%s' was trimmed to '%s' expected '%s'", test.s, s, test.exp)
		}
	}
}

func TestBaseIgnorer(t *testing.T) {

	tests := []struct {
		path    string
		base    string
		pattern string
		exp     bool
	}{
		{"/foo/bar", "/", "bar", true},
		{"/foo/bar", "/", "bar/", false},
		{"/foo/bar", "/", "foo", false},
		{"/foo/bar/", "/", "bar", true},
		{"/foo/bar/", "/foo", "bar", true},
		{"/foo/bar/hi.py", "/", "*.py", true},
		{"/foo/bar/hi.py", "/bar/", "*.py", false},
		{"foo/bar", "", "bar", true},
		{"foo/bar", "", "bar/", false},
		{"foo/bar", "", "foo", false},
		{"foo/bar/", "", "bar", true},
		{"foo/bar/", "foo", "bar", true},
		{"foo/bar/hi.py", "", "*.py", true},
		{"foo/bar/hi.py", "bar/", "*.py", false},
		{"foo/bar/hi.py", "", "*.py", true},
	}

	for _, test := range tests {
		i := &baseIgnorer{test.base, test.pattern}
		result := i.Ignore(file{test.path, false})
		if result == test.exp {
			continue
		}
		msg := "was not expected to match"
		if test.exp {
			msg = "was expected to match"
		}
		t.Errorf("'%s' with basepath '%s' %s '%s'", test.pattern, test.base,
			msg, test.path)
	}
}

func TestPathIgnorer(t *testing.T) {

	tests := []struct {
		path    string
		base    string
		pattern string
		exp     bool
	}{
		{"/foo/bar", "/", "/foo/bar", true},
		{"/foo/bar", "/", "foo/bar", false},
		{"/foo/bar", "/", "foo", false},
		{"/foo/bar/", "/", "/foo/bar/", true},
		{"/foo/bar/hi.py", "/", "/foo/bar/*.py", true},
		{"/foo/bar/hi.py", "/bar/", "/foo/bar/*.py", false},
	}

	for _, test := range tests {
		i := &pathIgnorer{test.base, test.pattern}
		result := i.Ignore(file{test.path, false})
		if result == test.exp {
			continue
		}
		msg := "was not expected to match"
		if test.exp {
			msg = "was expected to match"
		}
		t.Errorf("'%s' with basepath '%s' %s '%s'", test.pattern, test.base,
			msg, test.path)
	}
}
