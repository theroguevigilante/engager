package planner

import (
	"strings"
	"testing"

	"github.com/adaptive-learning/backend/internal/models"
)

func TestActivityPlanner(t *testing.T) {
	current := models.Activity{
		ID:         "act-recursion-1",
		Title:      "Recursion Basics",
		Topic:      "Computer Science",
		Difficulty: "beginner",
		Content:    "Understand base cases.",
	}

	catalog := []models.Activity{
		{
			ID:         "act-recursion-2",
			Title:      "Binary Tree Recursion",
			Topic:      "Computer Science",
			Difficulty: "intermediate",
			Content:    "Traverse trees recursively.",
		},
		{
			ID:         "act-math-1",
			Title:      "Matrix Multiplication",
			Topic:      "Linear Algebra",
			Difficulty: "beginner",
			Content:    "Multiply 2x2 matrices.",
		},
	}

	planner := NewPlanner(NewDeterministicAdapter())

	t.Run("BORED + INCREASE_CHALLENGE -> Selects harder activity from catalog", func(t *testing.T) {
		ass := &models.EngagementAssessment{State: models.StateBored}
		inter := &models.AttemptIntervention{ID: "int-1", InterventionType: models.InterventionIncreaseChallenge}
		result := planner.PlanActivity(ass, inter, current, catalog)

		if result.Difficulty != "intermediate" && result.Difficulty != "advanced" {
			t.Errorf("expected higher difficulty, got %s", result.Difficulty)
		}
		if result.SourceActivityID != "act-recursion-2" {
			t.Errorf("expected 'act-recursion-2', got %s", result.SourceActivityID)
		}
		if result.IsGenerated {
			t.Errorf("expected is_generated to be false when matched from catalog")
		}
	})

	t.Run("BORED + INCREASE_CHALLENGE (No harder catalog) -> Adapts current activity", func(t *testing.T) {
		ass := &models.EngagementAssessment{State: models.StateBored}
		inter := &models.AttemptIntervention{ID: "int-2", InterventionType: models.InterventionIncreaseChallenge}
		hardestCurrent := models.Activity{
			ID:         "hard-1",
			Title:      "Hard Activity",
			Topic:      "Computer Science",
			Difficulty: "advanced",
			Content:    "Maxed content.",
		}

		result := planner.PlanActivity(ass, inter, hardestCurrent, []models.Activity{})
		if !result.IsGenerated {
			t.Errorf("expected is_generated to be true when adapting fallback")
		}
		if result.Difficulty != "advanced" {
			t.Errorf("expected difficulty 'advanced', got %s", result.Difficulty)
		}
		if result.SourceActivityID != "hard-1" {
			t.Errorf("expected source activity ID 'hard-1', got %s", result.SourceActivityID)
		}
	})

	t.Run("STRUGGLING + SCAFFOLD (No lower catalog) -> Generates guided scaffold", func(t *testing.T) {
		ass := &models.EngagementAssessment{State: models.StateStruggling}
		inter := &models.AttemptIntervention{ID: "int-3", InterventionType: models.InterventionScaffold}

		// Intermediate activity with empty catalog to ensure scaffold drops difficulty to beginner
		intermedCurrent := models.Activity{
			ID:         "act-recursion-inter",
			Title:      "Tree Traversal",
			Topic:      "Computer Science",
			Difficulty: "intermediate",
			Content:    "Traverse binary trees.",
		}

		result := planner.PlanActivity(ass, inter, intermedCurrent, []models.Activity{})
		if !result.IsGenerated {
			t.Errorf("expected is_generated to be true for scaffold adaptation")
		}
		if result.Difficulty != "beginner" {
			t.Errorf("expected difficulty 'beginner', got %s", result.Difficulty)
		}
	})

	t.Run("DISENGAGED + CHANGE_ACTIVITY -> Selects different topic from catalog", func(t *testing.T) {
		ass := &models.EngagementAssessment{State: models.StateDisengaged}
		inter := &models.AttemptIntervention{ID: "int-4", InterventionType: models.InterventionChangeActivity}

		result := planner.PlanActivity(ass, inter, current, catalog)
		if result.SourceActivityID != "act-math-1" {
			t.Errorf("expected 'act-math-1' for topic change, got %s", result.SourceActivityID)
		}
		if result.IsGenerated {
			t.Errorf("expected is_generated to be false when topic switched via catalog")
		}
	})

	t.Run("Regression: Non-programming activity avoids CS jargon in challenge adaptation", func(t *testing.T) {
		moviePuzzle := models.Activity{
			ID:         "act-puzzle-1",
			Title:      "Emoji Movie Decoder",
			Topic:      "Cinema & Pop Culture",
			Content:    "Decode this movie: 🦁👑🌅",
			Difficulty: "beginner",
		}

		_, content := GenerateChallengeExtension(moviePuzzle)
		lower := strings.ToLower(content)

		forbidden := []string{
			"auxiliary memory",
			"time complexity",
			"space complexity",
			"unbounded",
			"collection is empty",
			"pseudocode",
			"recursive sub-problem",
		}
		for _, f := range forbidden {
			if strings.Contains(lower, f) {
				t.Errorf("found forbidden CS jargon %q in non-programming challenge adaptation:\n%s", f, content)
			}
		}
	})

	t.Run("Regression: Non-programming activity avoids CS jargon in scaffold adaptation", func(t *testing.T) {
		moviePuzzle := models.Activity{
			ID:         "act-puzzle-1",
			Title:      "Emoji Movie Decoder",
			Topic:      "Cinema & Pop Culture",
			Content:    "Decode this movie: 🦁👑🌅",
			Difficulty: "beginner",
		}

		_, content := GenerateScaffoldContent(moviePuzzle)
		lower := strings.ToLower(content)

		forbidden := []string{
			"base case",
			"recursive",
			"recursion",
			"pseudocode",
			"auxiliary memory",
			"time complexity",
		}
		for _, f := range forbidden {
			if strings.Contains(lower, f) {
				t.Errorf("found forbidden CS jargon %q in non-programming scaffold adaptation:\n%s", f, content)
			}
		}
	})
}
