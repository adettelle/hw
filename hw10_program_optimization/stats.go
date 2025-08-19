package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tidwall/gjson" //nolint:depguard
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]string

func getUsers(r io.Reader) (result users, err error) { //nolint:unparam
	scanner := bufio.NewScanner(r)

	i := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		val := gjson.GetBytes(line, "Email")
		result[i] = val.String()
		i++
	}
	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	suffix := "." + domain

	for _, email := range u {
		matched := strings.HasSuffix(email, suffix)

		if matched {
			idx := strings.Index(email, "@")
			domains := strings.ToLower(email[idx+1:])
			result[domains]++
		}
	}
	return result, nil
}
