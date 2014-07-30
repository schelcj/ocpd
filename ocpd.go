package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

func send_result(send_nsca string, config string, host string, result string) {
	cmd := exec.Command(send_nsca, "-H", host, "-c", config)
	cmd.Stdin = strings.NewReader(result)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	send_nsca := flag.String("send_nsca", "/usr/sbin/send_nsca", "Path to send_nsca binary")
	config := flag.String("config", "/etc/send_nsca.cfg", "Path to the send_nsca configuration file")
	host := flag.String("host", "nagios2", "Host to send nsca events")
	fifo := flag.String("fifo", "service-perfdata.fifo", "perfdata fifo")

	flag.Parse()

	fd, err := os.OpenFile(*fifo, os.O_RDWR, 0444)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		go send_result(*send_nsca, *config, *host, scanner.Text())
	}
}
