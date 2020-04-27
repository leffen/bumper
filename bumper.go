package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/sirupsen/logrus"
)

func main() {
	a := &app{}
	a.Assign("bumper")

	data, err := getVersionData(a.Input, a.FileName)

	if err != nil {
		if err.Error() == "No version info found" {
			a.Usage()
			os.Exit(-1)
		}
		logrus.Fatal(err)
	}

	old, _, _, newcontent, err := BumpInContent(data, a.Part)
	if err != nil {
		log.Fatal(err)
	}

	// Only show current version
	if a.Extract {
		print(a.Format, string(old))
		return
	}

	if a.Input == "" {
		err = ioutil.WriteFile(a.FileName, newcontent, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	print(a.Format, string(newcontent))
	return
}

func getVersionData(input, fileName string) ([]byte, error) {
	if input != "" {
		return []byte(input), nil
	}
	if fileName != "" {
		return ioutil.ReadFile(fileName)
	}
	return nil, fmt.Errorf("No version info found")
}

func print(format, version string) {
	if format == "" {
		fmt.Println(version)
		return
	}

	v := semver.New(version)
	var b bytes.Buffer
	for _, char := range format {
		if char == 'M' {
			b.WriteString(strconv.FormatInt(v.Major, 10))
		} else if char == 'm' {
			b.WriteString(strconv.FormatInt(v.Minor, 10))
		} else if char == 'p' {
			b.WriteString(strconv.FormatInt(v.Patch, 10))
		} else {
			b.WriteRune(char)
		}
	}
	fmt.Println(b.String())
}

// BumpInContent takes finds the first semver string in the content, bumps it, then returns the same content with the new version
func BumpInContent(vbytes []byte, part string) (old, new string, loc []int, newcontents []byte, err error) {
	data := strings.TrimSpace(string(vbytes))
	if len(data) < 1 {
		return "", "", nil, nil, fmt.Errorf("Did not find semantic version")
	}
	fields := strings.Split(data, " ")
	if len(fields) == 0 {
		return "", "", nil, nil, fmt.Errorf("Did not find semantic version")
	}

	re := regexp.MustCompile(`(\d+\.)?(\d+\.)?(\*|\d+)`)
	loc = re.FindIndex([]byte(fields[0]))

	if loc == nil {
		return "", "", nil, nil, fmt.Errorf("Did not find semantic version")
	}
	vs := string(vbytes[loc[0]:loc[1]])

	v := semver.New(vs)
	switch part {
	case "M":
		v.BumpMajor()
	case "m":
		v.BumpMinor()
	default:
		v.BumpPatch()
	}

	len1 := loc[1] - loc[0]
	additionalBytes := len(v.String()) - len1
	// Create and fill an extended buffer
	b := make([]byte, len(vbytes)+additionalBytes)
	copy(b[:loc[0]], vbytes[:loc[0]])
	copy(b[loc[0]:loc[1]+additionalBytes], v.String())
	copy(b[loc[1]+additionalBytes:], vbytes[loc[1]:])
	// fmt.Printf("writing: '%v'", string(b))

	return vs, v.String(), loc, b, nil
}
