package service

import "github.com/adaptive-learning/backend/internal/models"

var defaultActivities = []models.Activity{
	{Title: "Decode the movie", Description: "Guess the film from emoji clues.", Topic: "Movie Puzzles", Difficulty: "beginner", EstimatedDurationMinutes: 3,
		Content: "Decode this movie: 🦁👑🌅", CorrectAnswer: "The Lion King"},
	{Title: "What comes next?", Description: "Spot the pattern in the sequence.", Topic: "Mathematics", Difficulty: "intermediate", EstimatedDurationMinutes: 4,
		Content: "2, 6, 12, 20, 30, ?  What is the next number?", CorrectAnswer: "42"},
	{Title: "The missing variable", Description: "Find the value that balances the equation.", Topic: "Mathematics", Difficulty: "beginner", EstimatedDurationMinutes: 3,
		Content: "Given 3x + 9 = 24. What is the value of x?", CorrectAnswer: "5"},
	{Title: "Crack the Caesar cipher", Description: "Decode a word shifted by one letter.", Topic: "Cybersecurity", Difficulty: "beginner", EstimatedDurationMinutes: 4,
		Content: "Decrypt the word 'IFMMP' (each letter was shifted forward by 1, so shift backwards by 1).", CorrectAnswer: "HELLO"},
	{Title: "Spot the phishing email", Description: "Choose the safest response to a suspicious message.", Topic: "Cybersecurity", Difficulty: "beginner", EstimatedDurationMinutes: 3,
		Content: "An email from 'IT-Support' asks you to confirm your password using a link. Which is best?\nA) Click the link\nB) Reply with your password\nC) Report it to the security team\nD) Forward it to friends\nAnswer with a letter.", CorrectAnswer: "C"},
	{Title: "Binary search by hand", Description: "Trace a classic algorithm.", Topic: "Programming", Difficulty: "intermediate", EstimatedDurationMinutes: 6,
		Content: "Sorted array: [2, 5, 8, 12, 16, 23, 38]. Using binary search for 23, at which 0-based index is it found?", CorrectAnswer: "5"},
	{Title: "Who is the shortest?", Description: "A small logic puzzle.", Topic: "Logic", Difficulty: "beginner", EstimatedDurationMinutes: 3,
		Content: "Alice is taller than Bob. Bob is taller than Carol. Dana is shorter than Carol. Who is the shortest?", CorrectAnswer: "Dana"},
	{Title: "Design an alien", Description: "An invention challenge.", Topic: "Creative Writing", Difficulty: "intermediate", EstimatedDurationMinutes: 8,
		Content: "Invent an alien that can survive on a planet with no sunlight. Describe it in three sentences."},
	{Title: "Why do traffic jams persist?", Description: "Think in feedback loops.", Topic: "Systems Thinking", Difficulty: "advanced", EstimatedDurationMinutes: 8,
		Content: "A town adds a new bus line, but the roads are still jammed. Describe one feedback effect that could explain why."},
	{Title: "A time I got stuck", Description: "A short reflection.", Topic: "Reflection", Difficulty: "beginner", EstimatedDurationMinutes: 5,
		Content: "Describe a time you got stuck learning something. What helped you get unstuck?"},
}

// SeedDefaults inserts a varied starter library. Without force it only runs on an empty library.
func (s *Service) SeedDefaults(force bool) error {
	if !force {
		existing, err := s.repo.ListActivities()
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return nil
		}
	}
	for _, a := range defaultActivities {
		if _, err := s.CreateActivity(a); err != nil {
			return err
		}
	}
	return nil
}
