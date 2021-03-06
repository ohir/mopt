# mopt
### a minimalistic Go "flag" substitute that parses cmdline in the getopt style
`import "github.com/ohir/mopt"`

Package mopt provides command arguments parsing in 90 lines of Go.

It is ready to use right after declaring a single `mopt.Usage` type variable optionally containing _usage_ text to be printed with an '-h' flag. For just a single option a single line with Usage literal is all you need:
``` go
niter := mopt.Usage("-n iterations, default: 9").OptN('n', 9)
```
For more options declare Usage string variable, eg:
``` go
var cl mopt.Usage = "\t-n iterations\n\t-v verbose\n ..."
```

then you may use its five option read <a name="pkg-api">methods:</a>
- `cl.OptB('x') bool` tells if '-x' flag is present. Its ok to ask for any type flag presence.
- `cl.OptN('x', def·int) int` returns '-x ±digits' as an int, or the default int.
- `cl.OptF('x', def·float) float64` returns '-x ±digits' as a float64, or the default float.
- `cl.OptS('x', "default string") string` returns text of '-x text', or the default string.
- `cl.OptL() []string` returns slice of arguments following the last option or terminating '--'.

Spaces between flag letter and value are unimportant: ie. `-a bc`, and `-abc` are equivalent.  Same for numbers: `-n-3` and `-n -3` both provide _-3_ number. _For this elasticity a leading dash of string value, if needed, must be given after a backslash: eg. `-s\-dashed` or `-s "\- started with a dash"`. Flag grouping is not supported, too. Ie. `-a -b -c` are three boolean flags, but `-abc` would be an `-a` flag introducing a string value of "bc"_.

Flag `-h` is predefined to print a short "__ProgName__ _purpose, usage & options:_\n" lead, then content of the mopt.Usage variable; then program exits. Lead is kept in a package variable, so it can be changed from the user's code.

Automatic help behaviour can be extended simply by asking about a help topic early on: eg.
``` go
var cl mopt.Usage = "\t-v verbose\n ..."
func main(){
  if htopic := cl.OptS('h',"-"); htopic != "-" {
    switch htopic {
      case "": // bare -h
      case "flip": // -h flip
      // ...
      default:
        println("No help about", htopic, "avaliable!")
    }
    os.Exit(0) // exit after
  }
//...
}
```
----
Mopt package is meant to be used in the PoC code and ad-hoc cli tools. It parses two leading bytes of each os.Args entry anew on every OptX call. Also, there is no user feedback of _"unknown/wrong option"_, nor developer is guarded against opt-letter reuse. _Caveat emptor!_

## <a name="pkg-index">Usage</a>
* [type Usage](#Usage)
  * [func (u Usage) OptB(flag rune) bool](#Usage.OptB)
  * [func (u Usage) OptF(flag rune, def float64) float64](#Usage.OptF)
  * [func (u Usage) OptN(flag rune, def int) int](#Usage.OptN)
  * [func (u Usage) OptS(flag rune, def string) string](#Usage.OptS)
  * [func (u Usage) OptL() (r []string)](#Usage.OptL)


### <a name="Usage">type</a> [Usage](/mopt.go?s=2045:2061#L45)
``` go
type Usage string
```
Usage type string provides help message to be printed if program user will pass the '-h' flag. Mopt package whole api is hooked on an Usage type variable.

### <a name="Usage.OptB">func</a> (Usage) [OptB](/mopt.go?s=3541:3579#L96)
``` go
func (u Usage) OptB(flag rune) (r bool)
```
Method OptB returns true if flag was given, otherwise it returns false.
 It need not to take a default: flag either is present, or not.


### <a name="Usage.OptF">func</a> (Usage) [OptF](/mopt.go?s=3212:3266#L83)
``` go
func (u Usage) OptF(flag rune, def float64) (r float64)
```
Method OptF returns float64 read as f32 from string following the flag.  If flag was not given, or it could not be parsed to the float, OptF returns the def value.

### <a name="Usage.OptL">func</a> (Usage) [OptL](/mopt.go?s=3817:3850#L104)
``` go
func (u Usage) OptL() (r []string)
```
Method OptL returns a slice of strings filled with commandline arguments after the last option, or arguments after the options terminator '--', if given. Or all arguments if no dash-letter was spotted.

### <a name="Usage.OptN">func</a> (Usage) [OptN](/mopt.go?s=2758:2804#L68)
``` go
func (u Usage) OptN(flag rune, def int) (r int)
```
Method OptN returns an int. If flag was not given, or string that followed could not be parsed to an int, OptN returns the def value. Negative values need no special attention: -x-2 and -y -2 both convey -2.

### <a name="Usage.OptS">func</a> (Usage) [OptS](/mopt.go?s=2283:2331#L50)
``` go
func (u Usage) OptS(flag rune, def string) string
```
Method OptS returns following string. If flag was not given, OptS returns the def value. If string after option needs to begin with a dash character, leading dash must be escaped: eg. `-s"\-begins with a dash"`.

### NLS variable
``` go
var HelpLead string = "purpose, usage & options:\n"
```
Allows `-h` to say "propósito, uso y opciones:", or "目的、使用法、オプション"

### Exit variable
``` go
var Exit func(int) = os.Exit
```
os.Exit(0) called by -h support can be hijacked here.

----

* [API](#pkg-api)
* [Index](#pkg-index)
