package subcommand

import (
	"fmt"
	"strings"
)

//FlagType defines the different flag types. Options have values associated to the flag, Switches have no value associated.
type FlagType int

const (
	Option FlagType = iota
	Switch
)

//Convinience type for funcions passed flags
type FlagFunction func(string, string) error

//Flag structure
type Flag struct {
	//long definition (--option OPTION)
	Long string
	//Short definition (-o )
	Short string
	//Description
	Description string
	//FlagType, option or switch
	Type FlagType
	//Function to call when the flag is found during the parsing process
	fn func(string, string) error
	//Says if the flag is optional or mandatory
	Mandatory bool
}

//Must sets the flag as mandatory. The parser will raise an error in case it isn't present in the arguments
//TODO make sure that switches are not allowed to get mandatory
func (f *Flag) Must(isIt bool) {
	f.Mandatory = isIt
}

//Gets a help friendly flag representation:
//-o,--option  OPTION           This option does this and that
//-s,--switch                   This is a switch
//-i,--ignoreme [IGNOREME]      Optional option
func (f Flag) String() string {
	return fmt.Sprintf("%s\t%s", f.FlagStringPrefix(), f.Description)
}

func (f Flag) FlagStringPrefix() string {
	var format string
	var prefix string
	shortFormat := "%v"
	if f.Short != "" {
		shortFormat = "-%v,"
	}
	if f.Type == Option {
		if f.Mandatory {
			format = "--%v %v"
		} else {
			format = "--%v [%v]"
		}
		prefix = fmt.Sprintf(shortFormat+format, f.Short, f.Long, strings.ToUpper(f.Long))
	} else {
		format = "--%v"
		prefix = fmt.Sprintf(shortFormat+format, f.Short, f.Long)
	}
	return prefix
}

//Checks that the definition is just one word
func checkDefinition(flag string) bool {
	parts := strings.Split(flag, " ")
	return len(parts) == 1
}

//builds the flag struct panicking if errors are encountered
func buildFlag(long string, short string, desc string, fn FlagFunction, kind FlagType) *Flag {
	long = strings.Trim(long, " ")
	short = strings.Trim(short, " ")
	if len(long) == 0 {
		panic("Long definition is empty")
	}

	if !checkDefinition(long) {
		panic(fmt.Sprintf("Long definition %v has two words. Only one is accepted", long))
	}

	if !checkDefinition(short) {
		panic(fmt.Sprintf("Short definition %v has two words. Only one is accepted", long))
	}
	return &Flag{
		Type:        kind,
		Long:        long,
		Short:       short,
		fn:          fn,
		Description: desc,
		Mandatory:   false,
	}
}
