
const (

	// Process States

	// CREATED : process created
	CREATED = iota

	// RUNNING : process running
	RUNNING

	// WAITING : process waiting
	WAITING

	// BLOCKED : process blocked
	BLOCKED

	// TERMINATED : process terminated
	TERMINATED
)

var (

	// TimeQuantum : time quantum for process
	TimeQuantum = 50

	// ProcNum : PID for the highest process
	ProcNum int = 0
)

// Process : Running set of code
type Process struct {
	PID     int
	Name    string
	state   int
	runtime int
	memory  int
}

// Scheduler : Controller to schedule process to run
type Scheduler struct {
	inMsg     chan *Process
	processes []*Process
}

// Run : Start the schedule and process execution
func (s *Scheduler) Run() {
	for {

		// Check for new processes to schedule
		select {
		case x, ok := <-s.inMsg:
			if ok {
				// New process ready to be executed
				s.processes = append(s.processes, x)

			} else {
				// Channel is closed to execution must exit
				return
			}
		default:
			// No new processes
			break
		}

		for i, curProc := range s.processes {
			curProc.state = RUNNING

			// I'm assuming this will get much more complex beyond just subtracting runtime
			// Fortunately, as of now it is basic round robin execution
			curProc.runtime -= TimeQuantum
			time.Sleep(200 * time.Millisecond)

			if curProc.runtime <= 0 {
				s.processes = remove(s.processes, i)
			} else {
				curProc.state = WAITING
			}

		}

	}
}

// CreateProc : create a new process correctly
func CreateProc(name string, runtime int, mem int) *Process {

	ProcNum++

	return &Process{
		PID:     ProcNum,
		Name:    name,
		state:   CREATED,
		runtime: runtime,
		memory:  mem,
	}
}

func createRandomProcessFromTemplate(templateName string, instructions [][]string, ch chan *Process) {

	totalRuntime := 0
	for _, instruction := range instructions {
		if len(instruction) < 2 {
			continue
		}

		templateRuntime, err := strconv.Atoi(instruction[1])
		if err != nil {
			fmt.Println("error converting runtime to int", err)
		}

		// Jitter values by +-20
		templateRuntime += rand.Intn(20) - 10

		if templateRuntime < 0 {
			templateRuntime = 0
		}

		if instruction[0] == "CALCULATE" {
			totalRuntime += templateRuntime
		}

		instruction[1] = strconv.Itoa(templateRuntime)
	}

	p := CreateProc("From template: "+templateName, totalRuntime, rand.Intn(100)+1)
	ch <- p
}
