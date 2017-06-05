package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const usage = "usage: manifest <file> <component> <version>"

type tuple struct {
	Version string    `yaml:"version"`
	Release time.Time `yaml:"release"`
}

type manifest struct {
	VCLI       tuple   `yaml:"vcli"`
	Kernels    []tuple `yaml:"kernels"`
	Bootloader tuple   `yaml:"bootloader"`
	Trampoline tuple   `yaml:"trampoline"`
}

func main() {

	if len(os.Args) != 4 {
		fmt.Println(usage)
		os.Exit(1)
	}

	filename := os.Args[1]
	component := os.Args[2]
	version := os.Args[3]

	in, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error opening manifest: %s\n", err.Error())
		os.Exit(2)
	}

	manifest := new(manifest)
	err = yaml.Unmarshal(in, manifest)
	if err != nil {
		fmt.Printf("error unmarshalling manifest: %s\n", err.Error())
		os.Exit(3)
	}

	switch component {
	case "vcli":
		manifest.VCLI.Version = version
		manifest.VCLI.Release = time.Now()
	case "kernel":
		manifest.Kernels = append(manifest.Kernels, tuple{
			Version: version,
			Release: time.Now(),
		})
	case "bootloader":
		manifest.Bootloader.Version = version
		manifest.Bootloader.Release = time.Now()
	case "trampoline":
		manifest.Trampoline.Version = version
		manifest.Trampoline.Release = time.Now()
	default:
		fmt.Println("error: <component> must be one of 'vcli', 'kernels', 'bootloader', 'trampoline'")
		os.Exit(4)
	}

	out, err := yaml.Marshal(manifest)
	if err != nil {
		fmt.Printf("something went wrong: %s\n", err.Error())
		os.Exit(5)
	}

	err = ioutil.WriteFile(filename, out, 0777)
	if err != nil {
		fmt.Printf("something went wrong: %s\n", err.Error())
		os.Exit(6)
	}

}
