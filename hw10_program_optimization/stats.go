package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
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

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	buf := bufio.NewReader(r)

	var i int
	for {
		var line []byte
		line, _, err = buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
			}

			return
		}

		content := strings.TrimSpace(string(line))
		var user User
		if err = easyjson.Unmarshal([]byte(content), &user); err != nil {
			continue
		}

		result[i] = user
		i++
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.HasSuffix(user.Email, "."+domain) {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}

	return result, nil
}
