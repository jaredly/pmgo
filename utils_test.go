package main

import (
	. "launchpad.net/gocheck"
	"testing"
	"time"
)

// Hook up gocheck into the "go test" runner.
func TestUtils(t *testing.T) { TestingT(t) }

type MyUtilsSuite struct{}

var _ = Suite(&MyUtilsSuite{})

var deft = time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)
var today = time.Date(2000, 1, 2, 1, 2, 3, 0, time.UTC)
var yesterday = time.Date(2000, 1, 1, 1, 2, 3, 0, time.UTC)

var pdData = []struct {
	data  string
	value time.Time
}{
	{"", deft},
	{"today", today},
	{"yesterday", yesterday},
	{"02-30-2001", time.Date(2001, 2, 30, 0, 0, 0, 0, time.UTC)},
	{"05-20-1991", time.Date(1991, 5, 20, 0, 0, 0, 0, time.UTC)},
	{"5-20-1991", time.Date(1991, 5, 20, 0, 0, 0, 0, time.UTC)},
	{"5-2-1991", time.Date(1991, 5, 2, 0, 0, 0, 0, time.UTC)},
}

// Test _parseYamlTask_, We test for a partial task, multiple tasks, and
// nested tasks.
func (s *MyUtilsSuite) TestParseDate(c *C) {
	for _, item := range pdData {
		res := parseDate(item.data, deft, today)
		c.Check(res, DeepEquals, item.value)
	}
}
