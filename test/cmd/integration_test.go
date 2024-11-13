package test

import (
	"testing"

	order_suite "github.com/Na322Pr/route256/test/suite"
	"github.com/stretchr/testify/suite"
)

func TestSuit(t *testing.T) {
	suite.Run(t, &order_suite.OrderSuite{})
}
