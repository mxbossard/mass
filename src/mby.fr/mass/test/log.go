package main

import(
	"io"
	"math"
	"math/rand"
	"time"
	"strings"
	"strconv"
	"sync"
	"fmt"
	"log"

	"mby.fr/mass/internal/display"
	//"mby.fr/mass/internal/logger"
	//"mby.fr/mass/internal/output"
)

const loremString = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin facilisis mi sapien, vitae accumsan libero malesuada in. Suspendisse sodales finibus sagittis. Proin et augue vitae dui scelerisque imperdiet. Suspendisse et pulvinar libero. Vestibulum id porttitor augue. Vivamus lobortis lacus et libero ultricies accumsan. Donec non feugiat enim, nec tempus nunc. Mauris rutrum, diam euismod elementum ultricies, purus tellus faucibus augue, sit amet tristique diam purus eu arcu. Integer elementum urna non justo fringilla fermentum. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Quisque sollicitudin elit in metus imperdiet, et gravida tortor hendrerit. In volutpat tellus quis sapien rutrum, sit amet cursus augue ultricies. Morbi tincidunt arcu id commodo mollis. Aliquam laoreet purus sed justo pulvinar, quis porta risus lobortis. In commodo leo id porta mattis.`

var random = rand.New(rand.NewSource(99))

func writeLorem(size int, writer io.Writer) error {
	// Write full loremString
	for i := 0; i < size / len(loremString); i++ {
		writer.Write([]byte(loremString))
	}
	// Writer last part of loremString
	k := size % len(loremString)
	writer.Write([]byte(loremString[:k]))

	return nil
}

func lorem(size int) string {
	sb := strings.Builder{}
	err := writeLorem(size, &sb)
	if err != nil {
		log.Fatal(err)
	}
	return sb.String()
}

func produceText() string {
	size := random.Intn(80)
	return lorem(size)
}

const logPeriodInMs = 200

func infiniteActionLogger(d display.Displayer, action, subject string) {
	actionLogger(d, action, subject, math.MaxInt32)
}

func actionLogger(d display.Displayer, action, subject string, logCount int) {
	al := d.ActionLogger(action, subject, true)
	//outs := output.NewStandardOutputs()
	//al := logger.NewAction(outs, action, subject)
	//funcs := []func(string, ...interface{}){al.Trace, al.Debug, al.Warn, al.Error, al.Fatal}
	outputs := []io.Writer{al.Out(), al.Err()}
	loggingFuncs := []func(string, ...interface{}){al.Trace, al.Debug, al.Warn, al.Error}
	al.Info("Will log %d line(s).", logCount)
	var i = 0
	Loop:
	for {
		// Use logging functions
		for _, lf := range loggingFuncs {
			// Use outputs
			for _, o := range outputs {
				for k := 0; k < 3; k++ {
					text := strconv.Itoa(i) + " " + produceText()
					fmt.Fprintf(o, "%s\n", text)
					time.Sleep(logPeriodInMs * time.Millisecond)
				}
			}

			if i == logCount {
				al.Info("End logging. Logged %d line(s).", logCount)
				break Loop
			}
			text := strconv.Itoa(i) + " " + produceText()
			lf(text)
			time.Sleep(logPeriodInMs * time.Millisecond)
			i ++
		}
	}
	//al.Flush()
	//outs.Flush()
}

func main() {
	d := display.Service()
	actionLoggerCount := 20
	maxActiveActionLogger := 5
	maxLogCount := 10

	ch := make(chan int, maxActiveActionLogger)
	defer close(ch)
	var wg sync.WaitGroup

	for i := 0; i < actionLoggerCount; i++ {
		ch <- i
		wg.Add(1)
		action := "action_" + strconv.Itoa(i)
		subject := "subject_" + strconv.Itoa(i)
		logCount := random.Intn(maxLogCount)
		go func() {
			defer wg.Done()
			//fmt.Printf("Launching ActionLoger: %s %s for %d lines ...\n", action, subject, logCount)
			actionLogger(d, action, subject, logCount)
			//fmt.Printf("Finished ActionLoger: %s %s printed %d lines.\n", action, subject, logCount)
			<- ch
		}()
	}

	wg.Wait()
	d.Flush()
	fmt.Println("finished")
}
