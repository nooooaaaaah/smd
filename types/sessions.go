package types

import (
	"errors"
	"time"
)

func (s *ActiveSessions) AddActiveSession(session Session) error {
	switch {
	case !session.Token.ExpiresAt.After(time.Now()):
		return errors.New("Session is expired")
	case s.CheckForDuplicateSession(session):
		return errors.New("Session already exists")
	default:
		*s = append(*s, session)
		return nil
	}
}

func (s *ActiveSessions) RemoveActiveSession(session Session) {
	for i, v := range *s {
		if v.ID == session.ID {
			*s = append((*s)[:i], (*s)[i+1:]...)
		}
	}
}

func (s *ActiveSessions) RemoveActiveSessionByToken(token string) {
	for i, v := range *s {
		if v.Token.Token == token {
			*s = append((*s)[:i], (*s)[i+1:]...)
		}
	}
}

func (s *ActiveSessions) RemoveActiveSessionByUserID(userID string) {
	for i, v := range *s {
		if v.UserID == userID {
			*s = append((*s)[:i], (*s)[i+1:]...)
		}
	}
}

func (s *ActiveSessions) RemoveAllActiveSessions() {
	*s = ActiveSessions{}
}

func (s *ActiveSessions) GetAllActiveSessions() ActiveSessions {
	return *s
}

func (s *ActiveSessions) GetActiveSession(session Session) Session {
	for _, v := range *s {
		if v.ID == session.ID {
			return v
		}
	}
	return Session{}
}

func (s *ActiveSessions) GetActiveSessionByToken(token string) Session {
	for _, v := range *s {
		if v.Token.Token == token {
			return v
		}
	}
	return Session{}
}

func (s *ActiveSessions) GetActiveSessionByUserID(userID string) Session {
	for _, v := range *s {
		if v.UserID == userID {
			return v
		}
	}
	return Session{}
}

func (s *ActiveSessions) CountActiveSessions() int {
	return len(*s)
}

func (s *ActiveSessions) UpdateActiveSession(session Session) error {
	for i, v := range *s {
		if v.ID == session.ID {
			(*s)[i] = session
			return nil
		}
	}
	return errors.New("Session not found")
}

func (s *ActiveSessions) CheckForDuplicateSession(session Session) bool {
	for _, v := range *s {
		if v.ID == session.ID {
			return true
		}
	}
	return false
}

func (s *ActiveSessions) RemoveExpiredSessions() {
	for _, v := range *s {
		if !v.Token.ExpiresAt.After(time.Now()) {
			s.RemoveActiveSession(v)
		}
	}
}
