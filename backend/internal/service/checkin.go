package service

import "errors"

var ErrInvalidNeed = errors.New("unknown need")

var validNeeds = map[string]bool{"stuck": true, "bored": true, "curious": true, "lonely": true, "playful": true, "fine": true}

// RecordCheckin stores a student's self-reported need.
func (s *Service) RecordCheckin(studentID, need string) error {
	if !validNeeds[need] {
		return ErrInvalidNeed
	}
	_ = s.repo.EnsureStudent(studentID, "Demo Student")
	return s.repo.SaveCheckin(studentID, need)
}

func (s *Service) LatestNeed(studentID string) string {
	need, _ := s.repo.LatestCheckin(studentID)
	return need
}

func (s *Service) CommunityPulse() (map[string]int, error) {
	return s.repo.PulseCounts()
}
