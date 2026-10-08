package intervention

import "github.com/adaptive-learning/backend/internal/models"

var Catalog = map[models.InterventionType]models.InterventionMetadata{
	models.InterventionContinue: {
		Type:        models.InterventionContinue,
		Title:       "Keep Momentum",
		Description: "Maintain current progression and difficulty.",
		Category:    "REINFORCEMENT",
	},
	models.InterventionIncreaseChallenge: {
		Type:        models.InterventionIncreaseChallenge,
		Title:       "Level Up Challenge",
		Description: "Introduce higher difficulty or advanced problem constraints.",
		Category:    "ACCELERATION",
	},
	models.InterventionScaffold: {
		Type:        models.InterventionScaffold,
		Title:       "Guided Scaffolding",
		Description: "Provide hints, breakdown steps, or prerequisite review.",
		Category:    "ASSISTANCE",
	},
	models.InterventionReEngage: {
		Type:        models.InterventionReEngage,
		Title:       "Interactive Prompt",
		Description: "Offer an interactive check-in or quick knowledge checkpoint.",
		Category:    "ATTENTION",
	},
	models.InterventionChangeActivity: {
		Type:        models.InterventionChangeActivity,
		Title:       "Alternative Activity",
		Description: "Pivot to an alternative learning modality or topic format.",
		Category:    "ADAPTATION",
	},
}
