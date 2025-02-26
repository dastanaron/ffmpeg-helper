package ffmpeg

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
)

func ProgressPipe(cmd *exec.Cmd, channel chan int) {
	pipe, err := cmd.StderrPipe()

	if err != nil {
		log.Fatal(err)
	}

	buffer := bufio.NewReader(pipe)

	go asyncParsingProgress(buffer, channel)

	cmd.Start()
}

func asyncParsingProgress(buffer *bufio.Reader, channel chan int) {
	var durationInSeconds int
	for {
		byteLine, err := buffer.Peek(512)
		buffer.Discard(512)
		if err != nil {
			if err == io.EOF {
				channel <- 100
				close(channel)
				break
			}
			fmt.Println(err)
			break
		}

		pattern, err := regexp.Compile("(?i)duration.*")
		if err != nil {
			fmt.Println(err)
		}

		matched := pattern.Match(byteLine)

		if matched {
			line := string(byteLine)
			durationInSeconds = parseTime(`ion:\s([\d]{2}):([\d]{2}):([\d]{2})`, line)
		}

		pattern, err = regexp.Compile("time=.*")
		if err != nil {
			fmt.Println(err)
		}

		matched = pattern.Match(byteLine)

		if matched {
			line := string(byteLine)
			convertedTime := parseTime(`time=([\d]{2}):([\d]{2}):([\d]{2})`, line)
			percent := float64(convertedTime) / float64(durationInSeconds) * float64(100)

			if int(percent) == 100 {
				continue
			}

			channel <- int(math.Round(float64(percent)))
		}
	}
}

func parseTime(eregexp string, time string) int {
	ereg, err := regexp.Compile(eregexp)

	if err != nil {
		fmt.Println(err)
	}

	matches := ereg.FindStringSubmatch(time)
	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])
	return (hours * 3600) + (minutes * 60) + seconds
}
