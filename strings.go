package tUtils

import (
	"bytes"
	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
)

func Env(t Tester, key string) string {
	environment.RequireEnvVar(t.T(), key)
	return os.Getenv(key)
}

// Path returns a representation to an absolute path of the relativePath
func Path(t Tester, relativePath string) string {
	path, err := filepath.Abs(relativePath)
	require.NoError(t.T(), err)
	return path
}

// Substr replaces by the `sub` string all occurrence of the matched `pattern`
//  path 		string		 	path of the file to edit
//	pattern		string 			pattern to replace
//  sub			string 			substitution to the pattern
func Substr(t Tester, path, pattern, sub string) {
	path = Path(t, path)

	b, err := os.ReadFile(path)
	require.NoError(t.T(), err)

	b = bytes.Replace(
		b,
		[]byte(pattern),
		[]byte(sub),
		-1,
	)

	err = os.WriteFile(path, b, 0666)
	require.NoError(t.T(), err)
}

func Uuid() string {
	return strings.ToLower(random.UniqueId())
}
