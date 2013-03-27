package main

// import "io/ioutil"
import "fmt"
import "log"
import "time"

// import "strconv"
// import "strings"
import "launchpad.net/goyaml"

type A struct {
	string
}

type m string

func (a m) SetYAML(tag string, value string) bool {
	fmt.Printf("Out: %q %q\n", tag, value)
	return true
}

/*
The planner looks like a list of tasks:

  - name:      string
    date:      string
    id:        num
    completed: ~ or string (date)
    items:
    - [task] ...
*/
type Task struct {
	Name      string
	Date      time.Time
	Id        int64
	Completed time.Time
	Parent    *int64 // the ID of the parent
}

type yamlTask struct {
	Name      string
	Date      string
	Id        int64
	Completed string
	Items     []yamlTask
}

func parseYamlTask(data []byte) []yamlTask {
	var res []yamlTask
	goyaml.Unmarshal(data, &res)
	return res
}

func countYamlTasks(yaml *[]yamlTask, count *int64) {
	for _, item := range *yaml {
		*count += 1
		countYamlTasks(&item.Items, count)
	}
}

// Here we take a list of embedded yamlTasks, unfold them, and fill in IDs
// where needed.
func processYamlTasks(yaml *[]yamlTask) map[int64]Task {
	var count int64
	countYamlTasks(yaml, &count)
	// TODO make a few extra boxes?
	tasks := make(map[int64]Task)
	var unassigned int64 = -1
	// prepare a list to map our unassigned task ids
	unassigneds := make([]int64, count)
	var max int64 = 1
	var firstParent int64 = 0
	unfoldYamlTasks(yaml, &tasks, &unassigned, &max, &firstParent, &unassigneds)
	return tasks
}

// func processYamlTask(yaml *[]

func unfoldYamlTasks(yaml *[]yamlTask, tasks *map[int64]Task,
	unassigned *int64, max *int64, parent *int64,
	unassigneds *[]int64) {
	var no_time time.Time
	today := time.Now()
	var good bool
	var cparent []int64
	for _, item := range *yaml {
		good = true
		if item.Id == 0 {
			// if there was no ID given (comes out as zero), mark it as
			// unassigned
			item.Id = *unassigned
			*unassigned -= 1
			good = false
		} else if item.Id < 0 {
			// IDs should never be negative. Reassign.
			log.Printf("WARNING: negative ID found %d (for task %q). Reassigning\n",
				item.Id, item.Name)
			item.Id = *unassigned
			*unassigned -= 1
			good = false
		} else if _, present := (*tasks)[item.Id]; present {
			// If there's a duplicate ID while processing, we just reassign
			// the second one. This should never happen anyway unless you mess
			// up the file yourself.
			log.Printf("WARNING:  ID already used; %d (for task %q). Reassigning\n",
				item.Id, item.Name)
			item.Id = *unassigned
			*unassigned -= 1
			good = false
		}
		if item.Id > *max {
			// keep track of the maximum used id
			*max = item.Id
		}
		(*tasks)[item.Id] = Task{item.Name, parseDate(item.Date, today, today),
			item.Id, parseDate(item.Completed, no_time, today), parent}
		if len(item.Items) != 0 {
			if good {
				// The item has a good task ID, we just need a dummy in to pass
				// the pointer of
				cparent = make([]int64, 1)
				cparent[0] = item.Id
				unfoldYamlTasks(&item.Items, tasks, unassigned, max, &cparent[0], unassigneds)
			} else {
				// pass a reference to the unassigneds list
				unfoldYamlTasks(&item.Items, tasks, unassigned, max,
					&(*unassigneds)[-*unassigned-1], unassigneds)
			}
		}
	}
}

/*
func ParsePlanner(data []byte) map[int64]Task, error {
    var result []yamlTask
    err := goyaml.Unmarshal(data, &result)
    if err != nil {
        fmt.Println("Unable to parse planner.yaml")
        return nil, err
    }

    tasks := make([]Task, len(res))
}
*/

/***
The timesheet looks like:
"01-02-2006"

weekname:
    dayname:
        - id-id start:stop OR
        - id:    num
          tid:   num
          start: time
          end:   time
          completed: done
**/
// func ParseTimesheet(data []byte) map[Week]D
