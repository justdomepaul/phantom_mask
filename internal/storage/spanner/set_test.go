package spanner

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type StorageSetSuite struct {
	suite.Suite
}

func TestStorageExtendOptionsSuite(t *testing.T) {
	suite.Run(t, new(StorageSetSuite))
}
