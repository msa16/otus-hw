package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
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

type users *[]*User

func getUsers(r io.Reader) (users, error) {
	var json = jsoniter.ConfigFastest
	result := make([]*User, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		user := User{}
		if err := json.Unmarshal(line, &user); err != nil {
			return nil, err
		}
		result = append(result, &user)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainRegexp, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}
	for _, user := range *u {
		if domainRegexp.Match([]byte(user.Email)) {
			idx := strings.IndexByte(user.Email, '@')
			if idx > 0 {
				domain := strings.ToLower(user.Email[idx+1:])
				num := result[domain]
				num++
				result[domain] = num
			}
		}
	}
	return result, nil
}
