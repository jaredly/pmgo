package main

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
type MySuite struct{}
var _ = Suite(&MySuite{})

// Setup of things:
//   TaskMap {
//      map map[id int64]*Task  // positive = real ids, negative = placeholder
//      tasks []Task
//      unassigneds []int64
//      max int64
//      unassigned int64
//   }
//   Task {
//      parent *Task
//      id int64
//      name string
//      created date
//      modified date
//      completed date [string or null]
//      priority [A|B|C|D]
//   }
//   tasks.yml
//      processTasks(bytes) -> TaskMap
//          goyaml.Unmarshall(bytes, &[]interface{})
//              take the raw yaml, turn parse it into a raw list of interfaces
//              [might be maps, might be strings].
//          inflateTasks([]interface{}, *[]Task)
//              take the interfaces, inflate them into objects. Should be
//              injected into a taskmap
//          assembleTaskMap(*TaskMap) // do I want this to be copied?
//              Go through the task tree and populate a task map.
//          reassignTasks(TaskMap)
//              Go through the unassigned tasks, assigning them new IDs









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
  completed: 05-6-1991`, []yamlTask{{"one", "two", 5, "05-6-1991", nil}},
	}, {
		`- name: one
  items:
  - name: two
    date: three`, []yamlTask{{"one", "", 0, "", []yamlTask{{"two", "three", 0, "", nil}}}},
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

var countYamlTests = []struct {
  data []yamlTask
  count int64
} {
  {
    []yamlTask {
    },
    0,
  }, {
    []yamlTask {
      {"", "", 0, "", nil},
    },
    1,
  }, {
    []yamlTask {
      {"", "", 0, "", nil},
      {"", "", 0, "", nil},
      {"", "", 0, "", nil},
    },
    3,
  }, {
    []yamlTask {
      {"", "", 0, "", []yamlTask {
        {"", "", 0, "", nil},
      }},
    },
    2,
  }, {
    []yamlTask {
      {"", "", 0, "", []yamlTask {
        {"", "", 0, "", nil},
      }},
      {"", "", 0, "", []yamlTask {
        {"", "", 0, "", nil},
      }},
    },
    4,
  }, {
    []yamlTask {
      {"", "", 0, "", []yamlTask {
          {"", "", 0, "", []yamlTask {
            {"", "", 0, "", nil},
          }},
          {"", "", 0, "", nil},
        },
      },
      {"", "", 0, "", []yamlTask {
        {"", "", 0, "", nil},
      }},
    },
    6,
  },
}

func (s *MySuite) TestCountYamlTests(c *C) {
  for _, item := range countYamlTests {
    var count int64 = 0
    countYamlTasks(&item.data, &count)
    c.Check(count, Equals, item.count)
  }
}

/*
var unfoldTasks = []struct {
  input []yamlTask,
  output map[int64]Task
  unassigneds []int64
  unassigned int64
} {

}

// Test _unfoldYamlTasks_
func (s *MySuite) TestUnfoldYamlTasks(c *C) {
  for _, item := range unfoldTests {
    unassigneds := make([]int64, item.count)
    firstParent := 0
    unfoldYamlTasks(item.input, tasks, 
    res := 
    c.Check(res, DeepEquals, item.output
  }
}
*/


/**
var processTests = []struct {
  input []yamlTask
  output map[int64]Task
} {

}

// Test _processYamlTasks_
func (s *MySuite) TestProcessYamlTasks(c *C) {
  for _, item := range processTests {
    c.Check(res, DeepEquals, item.output)
  }
}
*/

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

