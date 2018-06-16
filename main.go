package main

import (
	"os/exec"
	"fmt"
	"log"
	"bytes"
	"strings"
	"regexp"
	"github.com/hashicorp/go-version"
	"github.com/tivvit/yapm/deb_version"
	"github.com/blang/semver"
)

func ExampleLookPath() {
	path, err := exec.LookPath("apt")
	if err != nil {
		log.Fatal("installing fortune is in your future")
	}
	fmt.Printf("fortune is available at %s\n", path)
}

func main() {
	ExampleLookPath()
	list()
	//update()
}

func update() {
	cmd := exec.Command("apt", "update", "--quiet")
	err := cmd.Run()
	if err != nil {
		log.Print(err)
	}
}

func list() {
	cmd := exec.Command("apt", "list", "--installed", "--quiet")
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%q", out.String())
	lines := strings.Split(out.String(),"\n")
	fmt.Println(len(lines))
	r, _ := regexp.Compile(`(\S+)/\S+ (\S+) .*`)
	for _, s := range lines {
		m := r.FindStringSubmatch(s)
		if m != nil {
			fmt.Println(m[1], m[2])
			//v, err := version.NewVersion(m[2])
			//if err != nil {
			//	fmt.Println(err)
			//} else {
			//	fmt.Println("\t", v.Metadata(), v.Prerelease(), v.Segments())
			//}
			v1, err := deb_version.NewVersion(m[2])
			if err != nil {
				fmt.Println(err)
			} else {
				//fmt.Println("\t", v1.Epoch, v1.DebianRevision, v1.UpstreamVersion)
				v, err := version.NewVersion(v1.UpstreamVersion)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("\t", v1.Epoch, v1.DebianRevision, v1.UpstreamVersion)
					fmt.Println("\t", v.Segments(),v.Metadata(), v.Prerelease())
				}
				v2, err := semver.Make(v1.UpstreamVersion)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("\t", v2.Build, v2.Major, v2.Minor, v2.Patch, v2.Pre)
				}
			}
		}
	}
}
