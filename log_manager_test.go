package log4g

import (
	"fmt"

	"github.com/jrivets/gorivets"
	"gopkg.in/check.v1"
)

type logManagerSuite struct {
}

var _ = check.Suite(&logManagerSuite{})

func (s *logManagerSuite) TestInit(c *check.C) {
	lm()
	lm()
	err := gorivets.CheckPanic(func() {
		lm().registerInGMap()
	})
	c.Assert(err, check.NotNil)
	fmt.Printf("Ok, with the err=%s", err)
}
