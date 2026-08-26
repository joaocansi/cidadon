package entity

import (
	"fmt"
	"strconv"
	"strings"
)

type AuthSubject struct {
	UserID uint
	Role   UserRole
}

func (s *AuthSubject) ToString() string {
	return fmt.Sprintf("%d:%s", s.UserID, s.Role)
}

func (s *AuthSubject) FromString(subject string) error {
	splitSubject := strings.Split(subject, ":")
	if len(splitSubject) != 2 {
		return fmt.Errorf("invalid subject")
	}
	userIDStrToInt, err := strconv.Atoi(splitSubject[0])
	if err != nil {
		return fmt.Errorf("invalid subject")
	}
	s.UserID = uint(userIDStrToInt)
	s.Role = UserRole(splitSubject[1])
	return nil
}
