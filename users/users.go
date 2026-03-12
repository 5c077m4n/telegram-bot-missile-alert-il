// Package users - user store
package users

import (
	"strings"
	"sync"
	"time"
)

var usersMap sync.Map

func GetUser(chatID int64) *User {
	if user, ok := usersMap.Load(chatID); ok {
		return user.(*User)
	}

	user := &User{ChatID: chatID, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	usersMap.Store(chatID, user)
	return user
}

func UpdateUserCity(chatID int64, city string) (*User, error) {
	if user, ok := usersMap.Load(chatID); ok {
		user := user.(*User)

		user.City = city
		user.UpdatedAt = time.Now()
		return user, nil
	}
	return nil, ErrUserNotFound
}

func GetAllUsers() []*User {
	var result []*User
	usersMap.Range(func(_key, value any) bool {
		result = append(result, value.(*User))
		return true
	})
	return result
}

func ContainsCity(alertData, city string) bool {
	return strings.Contains(strings.ToLower(alertData), strings.ToLower(city))
}

func ContainsCityArray(cities []string, city string) bool {
	cityLower := strings.ToLower(city)
	for _, c := range cities {
		if strings.Contains(strings.ToLower(c), cityLower) {
			return true
		}
	}
	return false
}
