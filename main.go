package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
)

const (
	volumesFlag    = "skip-volumes"
	networksFlag   = "skip-networks"
	imagesFlag     = "skip-images"
	containersFlag = "skip-containers"
	pruneFlag      = "skip-prune"
	forceFlag      = "force"
	helpFlag       = "help"
)

var (
	volumes    bool
	networks   bool
	images     bool
	containers bool
	prune      bool
	force      bool
	help       bool
)

func init() {
	flag.BoolVar(&volumes, volumesFlag, false, "Skip the deletion of volumes")
	flag.BoolVar(&networks, networksFlag, false, "Skip the deletion of networks")
	flag.BoolVar(&images, imagesFlag, false, "Skip the remove of images")
	flag.BoolVar(&containers, containersFlag, false, "Skip the remove of containers")
	flag.BoolVar(&prune, pruneFlag, false, "Skip the prune")
	flag.BoolVar(&force, forceFlag, false, "Force remove of containers and images")
	flag.BoolVar(&help, helpFlag, false, "Help")

	flag.Parse()
}

func main() {
	if help {
		flag.PrintDefaults()
		return
	}

	deleteVolumes()
	deleteNetworks()
	removeImages()
	removeContainers()
	prunes()
}

func deleteVolumes() {
	if volumes {
		fmt.Println(">> Volumes skipped")
	}
	command := "docker volume rm $(docker volume ls -qf dangling=true)"
	fmt.Printf(">> %s\n", command)
	executor(command)
}

func deleteNetworks() {
	if networks {
		fmt.Println(">> Networks skipped")
	}
	command := `docker network rm $(docker network ls | grep "bridge" | awk '/ / { print $1 }')`
	fmt.Printf(">> %s\n", command)
	executor(command)
}

func removeImages() {
	if images {
		fmt.Println(">> Images skipped")
	}
	var command string
	if force {
		command = "docker rmi -f $(docker images -aq)"
	} else {
		command = `docker rmi $(docker images --filter "dangling=true" -q --no-trunc)`
	}
	fmt.Printf(">> %s\n", command)
	executor(command)
}

func removeContainers() {
	if containers {
		fmt.Println(">> Containers skipped")
	}
	var command string
	if force {
		command = "docker rm -f $(docker ps -aq)"
	} else {
		command = `docker rm $(docker ps -qa --no-trunc --filter "status=exited")`
	}
	fmt.Printf(">> %s\n", command)
	executor(command)
}

func prunes() {
	if prune {
		fmt.Println(">> Prune skipped")
	}
	commandNetwork := `docker network prune -f`
	fmt.Printf(">> %s\n", commandNetwork)
	executor(commandNetwork)

	commandSystem := `docker system prune -f`
	fmt.Printf(">> %s\n", commandSystem)
	executor(commandSystem)
}

func executor(command string) {
	cmd := exec.Command("bash", "-c", command)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()
	if errb.Len() > 0 {
		fmt.Printf("error: \n")
		fmt.Printf("%s \n", errb.String())
	} else {
		fmt.Printf("out: \n")
		fmt.Printf("%s \n", outb.String())
	}
	fmt.Println("----------------------------")
}
