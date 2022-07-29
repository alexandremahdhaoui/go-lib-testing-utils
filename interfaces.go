package tUtils

import (
	"testing"
)

type Tester interface {
	T() *testing.T
}

type Identifier interface {
	Id() string
}
