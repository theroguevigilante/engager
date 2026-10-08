package planner

import (
	"github.com/adaptive-learning/backend/internal/models"
	"github.com/google/uuid"
)

// ActivityAdapter defines the contract for adapting or generating activity variations.
// This interface allows swapping the deterministic generator with an LLM provider later.
type ActivityAdapter interface {
	Adapt(source models.Activity, intervention models.InterventionType, reason string, targetDifficulty string) models.PersonalizedActivity
}

// DeterministicAdapter produces rule-based scaffolding and challenge extensions.
type DeterministicAdapter struct{}

func NewDeterministicAdapter() *DeterministicAdapter {
	return &DeterministicAdapter{}
}

func (a *DeterministicAdapter) Adapt(source models.Activity, intervention models.InterventionType, reason string, targetDifficulty string) models.PersonalizedActivity {
	diff := targetDifficulty
	if diff == "" {
		diff = source.Difficulty
	}

	switch intervention {
	case models.InterventionIncreaseChallenge:
		title, content := GenerateChallengeExtension(source)
		if targetDifficulty == "" {
			diff = "advanced"
		}

		return models.PersonalizedActivity{
			ID:                    uuid.New().String(),
			SourceActivityID:      source.ID,
			Title:                 title,
			Content:               content,
			Difficulty:            diff,
			InterventionType:      intervention,
			PersonalizationReason: reason,
			IsGenerated:           true,
		}

	case models.InterventionScaffold:
		title, content := GenerateScaffoldContent(source)
		if targetDifficulty == "" {
			diff = "beginner"
		}

		return models.PersonalizedActivity{
			ID:                    uuid.New().String(),
			SourceActivityID:      source.ID,
			Title:                 title,
			Content:               content,
			Difficulty:            diff,
			InterventionType:      intervention,
			PersonalizationReason: reason,
			IsGenerated:           true,
		}

	case models.InterventionReEngage:
		title, content := GenerateReEngageContent(source)

		return models.PersonalizedActivity{
			ID:                    uuid.New().String(),
			SourceActivityID:      source.ID,
			Title:                 title,
			Content:               content,
			Difficulty:            diff,
			InterventionType:      intervention,
			PersonalizationReason: reason,
			IsGenerated:           true,
		}

	default:
		return models.PersonalizedActivity{
			ID:                    uuid.New().String(),
			SourceActivityID:      source.ID,
			Title:                 source.Title,
			Content:               source.Content,
			Difficulty:            diff,
			InterventionType:      intervention,
			PersonalizationReason: reason,
			IsGenerated:           false,
		}
	}
}
