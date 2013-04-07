package main

// import "io/ioutil"

import (
   "fmt"
   "log"
   "time"
   "strconv"
   "strings"
   "launchpad.net/goyaml"
)

// 0001-01-01
var dayZero = time.Time

type Task struct {
  id        int64
  name      string
  created   time.Time
  modified  time.Time
  completed time.Time
  priority  uint8
  items      []Task
}

type TaskMap struct {
  tmap  map[int64]*Task
  tasks []Task
  unassigneds []int64 // ? keep track of the unassigned ones
  unassigned int64    // number of unassigned tasks
  max int64           // max id found
}

func processTasks(data) TaskMap {
  taskmap := make(TaskMap)
  var yaml []interface{}
  goyaml.Unmarshall(data, &yaml)
  inflateTasks(yaml, &taskmap.tasks)
  assembleTaskMap(&taskmap)
  reassignTasks(&taskmap)
  return taskmap
}

func inflateTasks(yaml []interface{}, tasks *[]Task) {
  *tasks = make([]Task, len(yaml))
  var task Task
  for item, i := range yaml {
    populateTask(item, &(*tasks)[i])
  }
}

func fillTasks(yaml interface{}, tasks *[]Task) bool {
  switch yaml := yaml.(type) {
    case []interface{}:
      inflateTasks(yaml, tasks)
      return false
    default:
      return true
  }
}

// An individual task is either a dictionary, corresponding to what one would
// expect, or a string [concise]
// All dates follow the form outlined in utils.go.
//   today|yesterday|tomorrow|mm-dd-yyyy
// The title is not allowed to contain "|"
// Condensed form = [id|] title [| created | completed]
func populateTask(yaml interface{}, task *Task, today time.Time) {
  val := reflect.ValueOf(yaml)
  switch yaml.(type) {
    case string:
      if err := expandString(yaml.(string), task, today); err {
        log.Printf("Unable to expand task: %q\n", yaml)
      }
    case map[string]interface{}:
      if len (yaml) == 1 { // condensed string with children
        for k, v := range yaml {
          expandString(k, task, today)
          fillTasks(v, &(*task).items)
        }
      } else { // it's an expanded string
        mapToTask(yaml.(map[string]interface{}), task, today)
      }
    default:
      log.Printf("Unable to parse task: %q\n", yaml)
  }
}

// Condensed form = [id|] title [| created | modified | completed]
func expandString(yaml string, task *Task, today time.Time) bool {
  parts := strings.Split(yaml, "|")
  switch len(parts) {
    case 0:
      return true
    case 1:
      (*task).name = yaml
    default:
      at := 0
      if id, err := strings.Atoi(parts[0]); err == nil {
        (*task).id = id
        at += 1
      }
      (*task).name = parts[at]
      if len(parts) > at + 1 {
        thedate, err := interfaceDate(parts[at + 1], today, today)
        if err != nil {
          log.Printf("Invalid created date: %q\n", parts[at + 1])
        }
        (*task).created = thedate
      }
      if len(parts) > at + 2 {
        thedate, err := interfaceDate(parts[at + 2], today, today)
        if err != nil {
          log.Printf("Invalid modified date: %q\n", parts[at + 2])
        }
        (*task).modified = thedate
      }
      if len(parts) > at + 3 {
        thedate, err := interfaceDate(parts[at + 3], dayZero, today)
        if err != nil {
          log.Printf("Invalid modified date: %q\n", parts[at + 3])
        }
        (*task).completed = thedate
      }
      if len(parts) > at + 4 {
        log.Printf("Extra items in condensed form: %q", parts[at + 4:])
      }
  }
}

// Parse an ID and populate the task
func fillId(val interface{}, task *Task) bool {
  err := false
  switch val := val.(type) {
    case int:
      (*task).id = int64(val)
    case int64:
      (*task).id = val
    default:
      err = true
  }
  return err
}

// Parse a Name and populate the task
func fillName(val interface{}, task *Task) bool {
  err := false
  switch val := val.(type) {
    case string:
      (*task).name = val
    case int:
      (*task).name = strconv.Itoa(val)
    default:
      err = true
  }
  return err
}

// Parse a Date
func interfaceDate(val interface{}, deftime time.Time, today time.Time) (time.Time, error) {
  switch val := val.(type) {
    case string:
      return parseDate(val, deftime, today)
    default:
      var err DateError
      err.Value = fmt.Sprintf("%q", val)
      err.Default = deftime
      return deftime, err
  }
}

// Parse a priority. Can be A, B, C, "", or 0-3
func fillPriority(val interface{}, task *Task) bool {
  switch val := val.(type) {
    case string:
      switch val {
        case "A", "B", "C":
          (*task).priority = strings.Index("ABC", val)
          return false
        case "":
          (*task).priority = 0
          return false
        default:
          return true
      }
    case int, int64:
      switch val {
        case 0, 1, 2, 3:
          (*task).priority = val
          return false
        default:
          return true
      }
    default:
      return true
  }
}

// You have an interface, interpret the values of the map and populate the
// task. today => is the date to use as "today" for interpreting dates
func mapToTask(yaml map[string]interface{}, task *Task, today time.Time) {
  (*task).created = today
  (*task).modified = today
  (*task).completed = dayZero
  for key, val := range yaml {
    switch key {
      case "id":
        // parse the id
        if err := fillId(val, task); err != nil {
          log.Printf("Invalid ID found: %q\n", val)
        }
      case "name":
        // parse the name
        if err := fillName(val, task); err != nil {
          log.Printf("Task name must be text (%q found)\n", val)
        }
      case "created":
        // parse the "created" date
        if thedate, err := interfaceDate(val, today, today); err != nil {
          log.Println(err.Error())
        }
        (*task).created = thedate
      case "modified":
        // parse the modified date
        if thedate, err := interfaceDate(val, today, today); err != nil {
          log.Println(err.Error())
        }
        (*task).modified = thedate
      case "completed":
        // parse the completed date, if any. Default is nil
        thedate, err := interfaceDate(val, dayZero, today)
        if err {
          log.Println(err.Error())
        }
        (*task).completed = thedate
      case "priority":
        // parse the priority string
        if err := fillPriority(val, task); err != nil {
          log.Printf("Invalid priority found %q\n", val)
        }
      case "items":
        // parse the items
        if err := fillTasks(val, &(*task).items); err != nil {
          log.Printf("Unable to parse sub items %q\n", val)
        }
      default:
        log.Printf("Unrecognized key found: %q\n", key)
    }
  }
}

