package main

import "io/ioutil"
import "fmt"

// import "strconv"
// import "launchpad.net/goyaml"

// import "strings"

type Goal struct {
	Id        int
	Completed bool
}

func (self *Goal) SetYAML(tag string, value interface{}) bool {
	fmt.Printf("Set %s\n", tag)
	return true
}

type Day struct {
	Day   string
	Goals []Goal
}

func (self *Day) SetYAML(tag string, value []Goal) bool {
	fmt.Printf("sat: %s %s", tag, value)
	self.Day = tag
	self.Goals = value
	return true
}

func (self *Day) String() string {
	return self.Day
}

/*
func ParseTasks(data []byte) []Task {
    var res []string
    err := goyaml.Unmarshal(data, &res)
    fmt.Printf("Things %s\n", res)
    if err != nil {
        fmt.Println("Failed to unmarshall : %s", err)
        panic(err)
    }
    var tasks = make([]Task, len(res))
    var parts = make([]string, 2)
    for i := 0; i < len(res); i++ {
        parts = strings.Split(res[i], " ")
        num, err := strconv.Atoi(parts[0])
        if err != nil {
            fmt.Printf("Failed to conv %i yeah %q\n", num, err)
            continue
        }
        tasks[i].Id = num
        tasks[i].Title = parts[1]
    }
    return tasks
}
*/

func main() {
	text, err := ioutil.ReadFile("tasks.yaml")
	if err != nil {
		fmt.Println("No a.yaml found")
	}
	// tasks := ParseTasks(text)
	fmt.Printf("Got Text: %s\n", text)
	// fmt.Printf("Got Text: %s\n", tasks)
	fmt.Println("Hello, 世界")
}

