package hw10programoptimization

import (
	"bufio"
	"io"
	"log"
	"strings"

	"github.com/tidwall/gjson" //nolint:depguard
)

type DomainStat map[string]int

const (
	emailField      = "Email"
	domainSeparator = "@"
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	suffix := "." + domain

	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatalf("error %v", scanner.Err())
		}
		line := scanner.Bytes()
		email := gjson.GetBytes(line, emailField).String()

		matched := strings.HasSuffix(email, suffix)

		if matched {
			idx := strings.LastIndex(email, domainSeparator)
			domains := strings.ToLower(email[idx+1:])
			result[domains]++
		}
	}
	return result, nil
}
