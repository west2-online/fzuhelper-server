package environment

import "flag"

func IsTestEnvironment() bool {
	return flag.Lookup("test.v") != nil
}
