// (c) 2015-2021 Ohir Ripe. MIT license.
package mopt

import (
	"os"
	"testing"
)

var (
	get    Usage     = "Help message of the mopt package (testing)."
	osA              = os.Args
	exited int       = 12345
	texit  func(int) = func(code int) { exited = code }
)

func init() {
	Exit = texit
}

func TestBool(t *testing.T) {
	os.Args = []string{"test", "a", "-a", "-b"}
	if !get.OptB('b') {
		t.Logf("OptB did not found bool flag 'b' while it should!")
		t.Fail()
	}
	if !get.OptB('a') {
		t.Logf("OptB did not found bool flag 'a' while it should!")
		t.Fail()
	}
	if get.OptB('x') {
		t.Logf("OptB found bool flag 'x' but it should not!")
		t.Fail()
	}
	os.Args = osA
}

func TestString(t *testing.T) {
	os.Args = []string{
		"test", "-aStringA", "-b", "String2",
		`-c\-dashed`, "-d", "-dashed", "-e", `\-dashed`,
	}
	if s := get.OptS('b', "x"); s != "String2" {
		t.Logf("OptS flag 'b' errs: expected String2, got >%s<", s)
		t.Fail()
	}
	if s := get.OptS('a', "x"); s != "StringA" {
		t.Logf("OptS flag 'a' errs: expected StringA, got >%s<", s)
		t.Fail()
	}
	if s := get.OptS('x', "xxx"); s != "xxx" {
		t.Logf("OptS default subst errs: expected xxx, got >%s<", s)
		t.Fail()
	}
	if s := get.OptS('c', "ups!"); s != "-dashed" {
		t.Logf("OptS default subst errs: expected >-dashed<, got >%s<", s)
		t.Fail()
	}
	if s := get.OptS('d', "ups!"); s != "" {
		t.Logf("OptS value with unescaped dash should return empty, but gave >%s<", s)
		t.Fail()
	}
	if s := get.OptS('e', "ups!"); s != "-dashed" {
		t.Logf("OptS default subst errs: expected >-dashed<, got >%s<", s)
		t.Fail()
	}
	os.Args = osA
}

func TestNum(t *testing.T) {
	os.Args = []string{"test", "-a-222", "-b", "11", "-c", "-55", "-d0xFFEG"}
	if n := get.OptN('d', 77); n != 77 {
		t.Logf("OptN flag 'd' wrong hex: expected default, got %d", n)
		t.Fail()
	}
	if n := get.OptN('c', 1); n != -55 {
		t.Logf("OptN flag 'b' errs: expected -55, got %d", n)
		t.Fail()
	}
	if n := get.OptN('b', 2); n != 11 {
		t.Logf("OptN flag 'b' errs: expected 11, got %d", n)
		t.Fail()
	}
	if n := get.OptN('a', 3); n != -222 {
		t.Logf("OptN flag 'a' errs: expected -222, got %d", n)
		t.Fail()
	}
	if n := get.OptN('x', 333); n != 333 {
		t.Logf("OptN default subst errs: expected 333, got %d", n)
		t.Fail()
	}
	os.Args = osA
}

func TestFloat(t *testing.T) {
	feq := func(a, b float64) bool {
		return (a-b) < float64(0.000001) && (b-a) < float64(0.000001)
	}
	os.Args = []string{"test", "-a-222.22", "-b", "11.1", "-c", "-5.5", "-d-0,77e-4"}
	if n := get.OptF('d', 1.0); !feq(n, 1.0) {
		t.Logf("OptF flag 'd' errs: expected default 1.0, got %f", n)
		t.Fail()
	}
	if n := get.OptF('c', 1.0); !feq(n, -5.5) {
		t.Logf("OptF flag 'b' errs: expected -5.5, got %f", n)
		t.Fail()
	}
	if n := get.OptF('b', 2.0); !feq(n, 11.1) {
		t.Logf("OptF flag 'b' errs: expected 11.1, got %f", n)
		t.Fail()
	}
	if n := get.OptF('a', 3.0); !feq(n, -222.220001) {
		t.Logf("OptF flag 'a' errs: expected -222.22, got %f", n)
		t.Fail()
	}
	if n := get.OptF('x', 33.3); !feq(n, 33.3) {
		t.Logf("OptF default subst errs: expected 333, got %f", n)
		t.Fail()
	}
	os.Args = osA
}

func TestHelp(t *testing.T) {
	os.Args = []string{"test", "-h", "subtopic"}
	if s := get.OptS('h', "bad"); s != "subtopic" {
		t.Logf("No -h subtopic retrieved. Bad! >%s<", s)
		t.Fail()
	}
	get.OptB('x') // none -x but -h is there
	if exited != 0 {
		t.Logf("The -h option did not called to Exit. Bad!")
		t.Fail()
	}
}

func TestList(t *testing.T) {
	os.Args = []string{"test", "-a", "string", "-n12", "one", "two", "three"}
	r := get.OptL()
	if len(r) != 3 {
		t.Logf("No list, or bad list retrieved. Len should be 3, is %d => %v", len(r), r)
		t.Fail()
	}
	os.Args = []string{"test", "--", "one", "-x", "two", "three", "-y"}
	r = get.OptL()
	if len(r) != 5 {
		t.Logf("No list, or bad list retrieved after --. Len:%d => %v", len(r), r)
		t.Fail()
	}
	os.Args = []string{"test", "one", "two", "three"}
	r = get.OptL()
	if len(r) != 3 {
		t.Logf("No list, or bad list retrieved. Len should be 3, is %d => %v", len(r), r)
		t.Fail()
	}
	os.Args = []string{"test", "one", "two", "--"}
	r = get.OptL()
	if len(r) != 0 {
		t.Logf("Tail list out of thin air. Len should be 0, is %d => %v", len(r), r)
		t.Fail()
	}
	os.Args = []string{"test", "une", "--", "-duo", "tre"}
	r = get.OptL()
	if len(r) != 2 || r[0] != "-duo" || r[1] != "tre" {
		t.Logf("No list, or bad list retrieved. Len should be 3, is %d => %v", len(r), r)
		t.Fail()
	}
}

func TestTerminator(t *testing.T) {
	os.Args = []string{"test", "-astr", "-n12", "-d", "--", "-e", "end"}
	if !get.OptB('d') {
		t.Logf("Flag -d not found at the -- terminator! Bad!")
		t.Fail()
	}
	if get.OptB('e') {
		t.Logf("Flag -e found after -- terminator! Bad!")
		t.Fail()
	}
}
