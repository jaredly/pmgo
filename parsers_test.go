package main

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

var parseTests = []struct {
	data  string
	value []yamlTask
}{
	{`- name: One Person`, []yamlTask{{"One Person", "", 0, "", nil}}},
	{`- name: one Person
- name: two person`, []yamlTask{{"one Person", "", 0, "", nil},
		{"two person", "", 0, "", nil},
	},
	}, {
		`- name: one
  date: two
  id:   5
  completed: 05-6-1991
`, []yamlTask{{"one", "two", 5, "05-6-1991", nil}},
	}, {
		`- name: one
  items:
  - name: two
    date: three
`, []yamlTask{{"one", "", 0, "", []yamlTask{{"two", "three", 0, "", nil}}}},
	},
}

// Test _parseYamlTask_, We test for a partial task, multiple tasks, and
// nested tasks.
func (s *MySuite) TestParseYamlTask(c *C) {
	for _, item := range parseTests {
		res := parseYamlTask([]byte(item.data))
		c.Check(res, DeepEquals, item.value)
	}
}

/*
func (s *MySuite) TestParsePlanner() {
    res := pm.ParsePlanner(`- name: One Person
  date: today
  id: 5
  completed: ~
  items:
  - name: two
    date: tomorrow
    id: 3
`)
    c.Assert(
    c.Check(42, Equals, "42")
    c.Check(os.Errno(13), Matches, "perm.*accepted")
}
*/
