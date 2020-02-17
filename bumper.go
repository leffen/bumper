package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"

	"github.com/coreos/go-semver/semver"
)

func main() {
	a := &app{}
	a.Assign("bumper")

	var err error
	data := []byte{}

	if a.FileName != "" {
		data, err = ioutil.ReadFile(a.FileName)
		if err != nil {
			log.Fatal(err)
		}
	} else if a.Input != "" {
		data = []byte(a.Input)
	} else {
		fmt.Printf("Usage: ......")
		return
	}

	old, new, _, newcontent, err := BumpInContent(data, a.Part)
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
	print(a.Format, string(new))
	return
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
	re := regexp.MustCompile(`(\d+\.)?(\d+\.)?(\*|\d+)`)
	loc = re.FindIndex(vbytes)

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
