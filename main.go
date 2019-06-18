package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// Listen for kill commandos
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill, syscall.SIGKILL)
		<-s
		os.Exit(0)
	}()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/checkHaProxy", func(c *gin.Context) {
		postFile, err := c.FormFile("file")
		if err != nil {
			c.String(400, err.Error())
			return
		}
		file, err := postFile.Open()
		if err != nil {
			c.String(400, err.Error())
			return
		}

		data := make([]byte, postFile.Size)
		_, err = file.Read(data)
		if err != nil {
			c.String(400, err.Error())
			return
		}

		err = Validate(string(data))
		if err != nil {
			c.String(400, err.Error())
			return
		}
		c.String(200, "OK")
	})
	r.Run(":8223")
}

// Validate validates the config using haproxy
func Validate(config string) error {
	parsedConf, err := ParseConf(config)
	if err != nil {
		return err
	}
	data := []byte(parsedConf)

	tmpfile, err := ioutil.TempFile(os.TempDir(), "haproxyConf")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		return err
	}

	if err := tmpfile.Close(); err != nil {
		return err
	}

	cmd := exec.Command("haproxy", "-c", "-V", "-f", tmpfile.Name())
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()

	outErr := errb.String()
	if outErr != "" {
		return errors.New(outErr)
	}

	return nil
}

// ParseConf replaces groups and missing keys
func ParseConf(config string) (string, error) {
	lines := strings.Split(config, "\n")
	out := ""

	currentSection := ""
	for _, line := range lines {
		if Match(`^\s*#`, line) {
			continue
		}

		sectionMated := FullMatch("^(global|defaults|frontend|backend)", line)
		for _, match := range sectionMated {
			if len(match) >= 2 {
				currentSection = match[1]
			}
		}

		if currentSection == "global" {
			userOrGroupMatched := FullMatch(`^\s*(user|group)\s+(\w+)`, line)
			for _, match := range userOrGroupMatched {
				if len(match) >= 3 {
					line = strings.Replace(line, match[2], "root", 1)
				}
			}
		}

		if currentSection == "frontend" {
			certMatched := FullMatch(`^\s*bind\s+.{0,50}ssl\s+crt\s+((\/|\w|\d|\.|\\ )+)`, line)
			for _, match := range certMatched {
				if len(match) >= 2 {
					line = strings.Replace(line, match[1], "/etc/certs/fullKey.pem", 1)
				}
			}
		}

		if currentSection == "backend" {
			certMatched := FullMatch(`^\s*server\s+.{0,50}ssl.+ca-file\s+(.+\.(pem|cer))`, line)
			for _, match := range certMatched {
				if len(match) >= 2 {
					line = strings.Replace(line, match[1], "/etc/certs/fullKey.pem", 1)
				}
			}
		}

		extra := "\n"
		if out == "" {
			extra = ""
		}
		out = out + extra + line
	}

	return out, nil
}

// FullMatch matches a string with a regex
func FullMatch(regx string, toMatch string) [][]string {
	return regexp.MustCompile(regx).FindAllStringSubmatch(string(toMatch), -1)
}

// Match returns true if a match is found
func Match(regx string, toMatch string) bool {
	matched, err := regexp.MatchString(regx, toMatch)
	if err != nil {
		return false
	}
	return matched
}
