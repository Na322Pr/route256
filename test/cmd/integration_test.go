package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	order_suite "gitlab.ozon.dev/marchenkosasha2/homework/test/suite"
)

func TestSuit(t *testing.T) {
	suite.Run(t, &order_suite.OrderSuite{})
}
