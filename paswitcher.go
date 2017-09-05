package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

func main() {
	exec.Command("sh", "-c", `amixer -c1 sset "Auto-Mute Mode" Disabled`) //Needed to allow headphones and line-in to both be available
	sink, _ := getNextSink()
	inputs, _ := getSinkInputs()

	cmd := fmt.Sprintf("pactl set-sink-port %d %s", sink.index, sink.port)
	cmd = cmd + " & " + fmt.Sprintf("pacmd set-default-sink %d", sink.index)
	for _, i := range inputs {
		cmd = cmd + " & " + fmt.Sprintf("pacmd move-sink-input %d %d", i, sink.index)
	}
	exec.Command("sh", "-c", cmd).Output()
}

// The output type represents a sink/port combination
type output struct {
	index     int
	name      string
	port      string
	available bool
	selected  bool
}

// getSinkInputs returns a slice of input indexes
func getSinkInputs() (ret []int, err error) {
	o, err := exec.Command("sh", "-c", "pacmd list-sink-inputs").Output()
	if err != nil {
		return
	}

	so := string(o[:])

	re := regexp.MustCompile(`index: (\d+)`)
	inp := re.FindAllStringSubmatch(so, -1)

	for _, i := range inp {
		a, err := strconv.Atoi(i[1])
		if err != nil {
			return nil, err
		}
		ret = append(ret, a)
	}
	return
}

// getNextSink returns an output based on the current and available ports
func getNextSink() (output, error) {
	outputs, err := getAvailableOutputs()
	if err != nil {
		return output{}, err
	}

	for i, o := range outputs {
		if o.selected {
			if i+1 < len(outputs) {
				return outputs[i+1], nil
			}
			return outputs[0], nil
		}
	}
	return outputs[0], nil
}

// listSinks returns the output of "pacmd list-sinks" as a string
func listSinks() (string, error) {
	o, err := exec.Command("sh", "-c", "pacmd list-sinks").Output()
	if err != nil {
		return "", err
	}

	so := string(o[:])
	return so, nil
}

// getAvailableOutputs returns a slice of outputs, which have available = true
func getAvailableOutputs() ([]output, error) {
	var ret []output
	o, err := getOutputs()
	if err != nil {
		return nil, err
	}
	for _, a := range o {
		if a.available {
			ret = append(ret, a)
		}
	}
	return ret, nil
}

// getOutputs returns a slice of outputs of all the sinks and ports
func getOutputs() ([]output, error) {
	var ret []output
	re := regexp.MustCompile(`(?:  (\*| ) index: (\d+)[\S\s]+?name: <(.+)>[\S\s]+?ports:([\S\s]+?)active port: <(.+)>)+`)
	rePorts := regexp.MustCompile(`(\S+): .+ \(priority \d+, latency offset \d+ usec, available: (\S+)\)`)
	s, _ := listSinks()
	res := re.FindAllStringSubmatch(s, -1)

	for _, f := range res {
		ports := rePorts.FindAllStringSubmatch(f[4], -1)
		for _, t := range ports {
			index, err := strconv.Atoi(f[2])
			if err != nil {
				return nil, err
			}
			ret = append(ret, output{
				index,
				f[3],
				t[1],
				t[2] != "no",
				t[1] == f[5] && f[1] == "*",
			})
		}
	}
	return ret, nil
}
