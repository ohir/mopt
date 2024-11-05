// (c) 2015-2024 Ohir Ripe. MIT license.

/*
	Package mopt provides getopt style options parsing.

Its API consist of five OptX methods, with 'X' being of 'B'ool, 'S'tring,
'F'loat, 'N'umber (int), and finally 'L'ist - that returns list (slice) of
arguments after the last option (or after terminating --).

Declaration `var cl mopt.Usage = "usage/help"` is the only chore. Then you
just call one of cl.OptX(flag, default) methods where needed. If flag was
given you will get its value. If it was not - you have the default.

Mopt parses oldschool single letter options, and option "-h" is predefined
to print var Usage content (ie. "usage/help" string). Spaces between flag
letter and value are unimportant: ie. -a bc, and -abc are equivalent.
Same for numbers: -n-3 and -n -3 both provide -3 number.
For this elasticity a leading dash of string value must be given escaped
with a backslash: eg. -s\-dashed or -s "\- started with a dash"; and flag
grouping is not supported, too. Ie. -a -b -c are three boolean flags, but
-abc would be an -a option introducing a string value of "bc".

Mopt is meant to be used in the PoC code and ad-hoc cli tools. It parses
whole os.Args anew on each OptX call. There is no user feedback of "unknown
option", nor developer is guarded against opt-letter reuse. Caveat Emptor!
*/
package mopt

import (
	"os"
	"strconv"
	"strings"
)

/*
	Usage is the API holder type - to be filled with the message to be

printed as "help/usage" after the "prog - purpose, usage & options:\n"
predefined lead, if user gave the '-h' option.  After printing Usage
content program terminates returning zero.

Expanding help:
The -h flag and subsequent string can be retrieved early in a program without
stepping into the default `print Usage then Exit` path. If first call to the
Usage methods will be for the 'h', eg. `subhelp := cl.OptS("h","-")`, returned
"subhelp" string (!= "-") will tell that -h option is present, amd that user
possibly wants more help on topic. After servicing user needs, program should
terminate.
*/
type Usage string

// Method OptS returns string that followed the flag. If flag was not given,
// OptS returns the def value. If string after option needs to begin with
// dash character, it must be escaped: eg. -s"\-begins with a dash".
func (u Usage) OptS(flag rune, def string) string {
	if s, ok := u.optss(flag); ok {
		switch {
		case len(s) == 0:
			return def
		case s[0] == '-': // next option, allow string to be optional
			return def
		case len(s) < 2:
		case s[0] == '\\' && s[1] == '-':
			return s[1:]
		}
		return s
	}
	return def
}

// Method OptN returns int read from string that follows the flag.
// If flag was not given, or it could not be parsed to an int, OptN
// returns the def value. Negative values need no special attention:
// -a-2 and -a -2 both resolve to -2.
func (u Usage) OptN(flag rune, def int) (r int) {
	if s, ok := u.optss(flag); ok {
		r, err := strconv.Atoi(s)
		if err != nil {
			return def
		}
		return r
	}
	return def
}

// Method OptF returns float64 read as f32 from string following the flag.
// If flag was not given, or it could not be parsed to the float, OptF
// returns the def value. String is parsed to the float64 but of value
// that is convertible to the float32 without value changing.
func (u Usage) OptF(flag rune, def float64) (r float64) {
	if s, ok := u.optss(flag); ok {
		r, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return def
		}
		return r
	}
	return def
}

// Method OptB returns true if flag was given, otherwise it returns false.
//  It need not to take a default: flag either is present, or not.
func (u Usage) OptB(flag rune) (r bool) {
	_, r = u.optss(flag)
	return
}

// Method OptCSF returns bitflags given as a comma separated string past the
// flag.  It returns the copy of _current_ parameter with bit set or zeroed
// at position where the bit-name is present in the _all_ parameter.
//   - Returned value bit is set only if a "bitname" entry was given.
//   - Returned value bit is zeroed only if a "no-bitname" entry was given.
//
// Ex: -Fflag1,no-flag3
func (u Usage) OptCSF(flag rune, current uint32, all string) (r uint32) {
	fs, no := u.optss(flag)
	if r = current; !no {
		return
	}
	allf := strings.Split(all, ",")
	for _, fl := range strings.Split(fs, ",") {
		sb, sm := uint32(1), ^uint32(1)
		if no = len(fl) > 2 && fl[:3] == "no-"; no {
			fl = fl[3:]
		}
		for _, tf := range allf {
			if fl == tf && no {
				r &= sm
			} else if fl == tf {
				r |= sb
			}
			sb <<= 1
			sm = ^sb
		}
	}
	return
}

// Method OptL returns a slice of commandline arguments after the last
// option, or arguments after the terminating -- dashes, if given;
// or all arguments, if no dash-letter was spotted.
func (u Usage) OptL() (r []string) {
	r = os.Args[1:]
	i, lo := 0, 0
	for ; i < len(r); i++ {
		if len(r[i]) < 2 || r[i][0] != '-' {
			continue
		} else if r[i][1] == '-' && i < len(r)-1 {
			return r[i+1:]
		}
		lo = i + 1
	}
	if lo < len(r) {
		return r[lo:]
	} else {
		return []string{}
	}
}

func (u Usage) optss(flag rune) (s string, ok bool) {
	var i int
args:
	for i = 1; i < len(os.Args); i++ {
		s = os.Args[i]
		switch {
		case len(s) < 2:
			continue args
		case s[0] == '-' && s[1] == byte(flag):
			i++
			ok = true
			break args
		case s[0] == '-' && s[1] == 'h':
			println(os.Args[0], HelpLead, u)
			Exit(0)
		case s[0] == '-' && s[1] == '-':
			break args
		}
	}
	switch {
	case !ok:
	case len(s) > 2:
		return s[2:], ok
	case i < len(os.Args):
		return os.Args[i], ok
	}
	return "", ok
}

// Exit(0) can be hijacked here
var Exit func(int) = os.Exit

// Mopt can also say "propósito, uso y opciones:\n"
var HelpLead string = "purpose, usage & options:\n"
