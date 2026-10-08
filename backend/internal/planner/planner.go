package planner

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/adaptive-learning/backend/internal/models"
)

// Repository defines minimal data access methods if passed into planner.
type Repository interface {
	GetActivity(id string) (*models.Activity, error)
	ListActivities() ([]models.Activity, error)
}

// ActivityPlanner plans the next activity using catalog matching and LLM/deterministic fallbacks.
type ActivityPlanner struct {
	adapter ActivityAdapter
	repo    Repository
	llm     *LLMClient
}

// NewPlanner supports both legacy adapter-based initialization and repository-based initialization.
func NewPlanner(args ...interface{}) *ActivityPlanner {
	p := &ActivityPlanner{
		adapter: NewDeterministicAdapter(),
		llm:     NewLLMClient(),
	}

	for _, arg := range args {
		switch v := arg.(type) {
		case ActivityAdapter:
			p.adapter = v
		case Repository:
			p.repo = v
		}
	}
	return p
}

func NewActivityPlanner(repo Repository) *ActivityPlanner {
	return NewPlanner(repo)
}

// InterventionSpec defines intervention parameters for service-layer planning.
type InterventionSpec struct {
	InterventionType models.InterventionType
	Reasons          []string
	StudentState     string
	Need             string // self-reported need from the student check-in
}

// ActivityDomain categorizes the domain of an activity to prevent cross-domain contamination.
type ActivityDomain string

const (
	DomainProgramming ActivityDomain = "programming"
	DomainMathematics ActivityDomain = "mathematics"
	DomainPuzzle      ActivityDomain = "puzzle"
	DomainCreative    ActivityDomain = "creative"
	DomainGeneral     ActivityDomain = "general"
)

func containsWord(s string, words []string) bool {
	for _, w := range words {
		// Match exact word boundaries so "decode" does not trigger "code"
		pattern := `\b` + regexp.QuoteMeta(w) + `\b`
		matched, _ := regexp.MatchString(pattern, s)
		if matched {
			return true
		}
	}
	return false
}

func DetectDomain(act models.Activity) ActivityDomain {
	combined := strings.ToLower(fmt.Sprintf("%s %s %s", act.Topic, act.Title, act.Content))

	// 1. Puzzle / Riddle / Emoji check first to avoid false matches on words like "decode"
	puzzleKeywords := []string{"puzzle", "emoji", "movie", "decode", "riddle", "guess", "trivia", "teaser", "cinema", "pop culture"}
	if containsWord(combined, puzzleKeywords) {
		return DomainPuzzle
	}

	// 2. Mathematics check
	mathKeywords := []string{"math", "algebra", "calculus", "geometry", "equation", "solve for x", "matrix", "derivative", "integral", "fraction", "linear algebra"}
	if containsWord(combined, mathKeywords) {
		return DomainMathematics
	}

	// 3. Creative / Humanities check
	creativeKeywords := []string{"essay", "story", "write", "creative", "poem", "reflect", "discussion", "literature"}
	if containsWord(combined, creativeKeywords) {
		return DomainCreative
	}

	// 4. Programming check (exact word boundaries prevent "decode" -> "code")
	progKeywords := []string{"code", "programming", "python", "javascript", "golang", "function", "array", "binary tree", "recursion", "class", "syntax", "compile"}
	if containsWord(combined, progKeywords) {
		return DomainProgramming
	}

	return DomainGeneral
}

func GenerateChallengeExtension(orig models.Activity) (string, string) {
	domain := DetectDomain(orig)
	title := fmt.Sprintf("%s (Advanced Challenge)", orig.Title)
	var content string

	switch domain {
	case DomainProgramming:
		content = fmt.Sprintf(
			"⚡ CHALLENGE EXTENSION:\n%s\n\nPerformance & Engineering Constraints:\n1. Optimize your solution to run within optimal time complexity.\n2. Handle edge cases such as empty, single-element, or unbounded inputs without auxiliary memory overhead.",
			orig.Content,
		)
	case DomainMathematics:
		content = fmt.Sprintf(
			"⚡ CHALLENGE EXTENSION:\n%s\n\nAdvanced Investigation:\n1. Solve the problem using an alternative method (e.g., graphical or algebraic substitution).\n2. Formulate a general formula that holds true for any arbitrary constant.",
			orig.Content,
		)
	case DomainPuzzle:
		if strings.Contains(strings.ToLower(orig.Content), "movie") || strings.Contains(strings.ToLower(orig.Topic), "movie") || strings.Contains(strings.ToLower(orig.Topic), "cinema") {
			content = fmt.Sprintf(
				"⚡ CHALLENGE EXTENSION:\n%s\n\nBonus Puzzles:\n1. Decode these additional movie titles:\n   🧊🚢💔\n   🕷️👨\n2. Create your own movie puzzle using exactly three emojis.",
				orig.Content,
			)
		} else {
			content = fmt.Sprintf(
				"⚡ CHALLENGE EXTENSION:\n%s\n\nBonus Puzzles:\n1. Solve the puzzle in under 30 seconds.\n2. Construct a multi-step clue based on the same pattern for a peer.",
				orig.Content,
			)
		}
	case DomainCreative:
		content = fmt.Sprintf(
			"⚡ CHALLENGE EXTENSION:\n%s\n\nCreative Extension:\n1. Re-tell or re-frame your response from an opposing character's or stakeholder's point of view.\n2. Integrate vivid sensory details or rhetorical devices into your reflection.",
			orig.Content,
		)
	default:
		content = fmt.Sprintf(
			"⚡ CHALLENGE EXTENSION:\n%s\n\nDeep Dive Questions:\n1. Connect your answer to a real-world scenario outside this immediate context.\n2. Identify two potential counter-examples or limitations to your conclusion.",
			orig.Content,
		)
	}

	return title, content
}

func GenerateScaffoldContent(orig models.Activity) (string, string) {
	domain := DetectDomain(orig)
	title := fmt.Sprintf("%s (Guided Walkthrough)", orig.Title)
	var content string

	switch domain {
	case DomainProgramming:
		content = fmt.Sprintf(
			"🪜 GUIDED SCAFFOLD:\n%s\n\nStep-by-Step Breakdown:\nStep 1: Write down the expected input types and return values in pseudocode.\nStep 2: Solve for a small sample input with paper and pencil.\nStep 3: Implement the basic loop or condition without worrying about optimization.",
			orig.Content,
		)
	case DomainMathematics:
		content = fmt.Sprintf(
			"🪜 GUIDED SCAFFOLD:\n%s\n\nStep-by-Step Breakdown:\nStep 1: Identify and isolate the known quantities and variables.\nStep 2: Review the relevant formula or property needed for the first simplification step.\nStep 3: Check each arithmetic step before computing the final value.",
			orig.Content,
		)
	case DomainPuzzle:
		content = fmt.Sprintf(
			"🪜 GUIDED SCAFFOLD:\n%s\n\nHints:\n• Clue 1: Look at each symbol/emoji individually and identify what literal object or setting it depicts.\n• Clue 2: The setting and characters are key clues to a well-known classic.",
			orig.Content,
		)
	default:
		content = fmt.Sprintf(
			"🪜 GUIDED SCAFFOLD:\n%s\n\nGuided Walkthrough:\n• Core Focus: Start by summarizing the core question in your own words.\n• Step 1: Brainstorm 2-3 immediate observations before forming your conclusion.",
			orig.Content,
		)
	}

	return title, content
}

func GenerateReEngageContent(orig models.Activity) (string, string) {
	title := fmt.Sprintf("%s (Quick Check)", orig.Title)
	content := fmt.Sprintf(
		"🎯 FAST TRACK ACTIVITY:\nLet's switch things up with a rapid-fire check:\n%s\n\nSubmit your quickest initial hypothesis or single key takeaway in 1-2 sentences.",
		orig.Content,
	)
	return title, content
}

// PlanActivity matches catalog items first, and falls back to adaptation generation.
func (p *ActivityPlanner) PlanActivity(
	assessment *models.EngagementAssessment,
	interv *models.AttemptIntervention,
	current models.Activity,
	catalog []models.Activity,
) models.PersonalizedActivity {

	switch interv.InterventionType {
	case models.InterventionIncreaseChallenge:
		for _, act := range catalog {
			if act.ID != current.ID && act.Topic == current.Topic && isHigherDifficulty(act.Difficulty, current.Difficulty) {
				return models.PersonalizedActivity{
					ID:                    fmt.Sprintf("plan-%d", time.Now().UnixNano()),
					SourceActivityID:      act.ID,
					Title:                 act.Title,
					Content:               act.Content,
					Difficulty:            act.Difficulty,
					InterventionType:      interv.InterventionType,
					PersonalizationReason: "Selected a more challenging activity from the catalog.",
					IsGenerated:           false,
				}
			}
		}

	case models.InterventionChangeActivity:
		for _, act := range catalog {
			if act.ID != current.ID && act.Topic != current.Topic {
				return models.PersonalizedActivity{
					ID:                    fmt.Sprintf("plan-%d", time.Now().UnixNano()),
					SourceActivityID:      act.ID,
					Title:                 act.Title,
					Content:               act.Content,
					Difficulty:            act.Difficulty,
					InterventionType:      interv.InterventionType,
					PersonalizationReason: "Selected an alternative topic to re-engage learning.",
					IsGenerated:           false,
				}
			}
		}

	case models.InterventionScaffold:
		for _, act := range catalog {
			if act.ID != current.ID && act.Topic == current.Topic && isLowerDifficulty(act.Difficulty, current.Difficulty) {
				return models.PersonalizedActivity{
					ID:                    fmt.Sprintf("plan-%d", time.Now().UnixNano()),
					SourceActivityID:      act.ID,
					Title:                 act.Title,
					Content:               act.Content,
					Difficulty:            act.Difficulty,
					InterventionType:      interv.InterventionType,
					PersonalizationReason: "Selected a foundational activity from the catalog.",
					IsGenerated:           false,
				}
			}
		}
	}

	reason := string(interv.InterventionType)
	if interv.Reason != "" {
		reason = interv.Reason
	}

	targetDifficulty := current.Difficulty
	switch interv.InterventionType {
	case models.InterventionIncreaseChallenge:
		targetDifficulty = "advanced"
	case models.InterventionScaffold:
		targetDifficulty = "beginner"
	}

	// Call adapter with (source, intervention, reason, targetDifficulty)
	adapted := p.adapter.Adapt(current, interv.InterventionType, reason, targetDifficulty)

	// Preserve adapter's computed difficulty rather than overwriting with current.Difficulty
	if adapted.Difficulty == "" {
		adapted.Difficulty = targetDifficulty
	}
	adapted.ID = fmt.Sprintf("gen-%d", time.Now().UnixNano())
	adapted.SourceActivityID = current.ID
	adapted.InterventionType = interv.InterventionType
	adapted.IsGenerated = true

	if adapted.PersonalizationReason == "" {
		adapted.PersonalizationReason = fmt.Sprintf("Adapted based on %s state for targeted learning.", interv.InterventionType)
	}

	return adapted
}

// PlanPersonalizedActivity supports the service layer with LLM generation and catalog fallbacks.
func (p *ActivityPlanner) PlanPersonalizedActivity(
	ctx context.Context,
	attempt *models.ActivityAttempt,
	spec InterventionSpec,
	orig models.Activity,
) (*models.PersonalizedActivity, error) {

	if p.llm != nil && p.llm.IsConfigured() {
		state := spec.StudentState
		if state == "" {
			state = attempt.Status
		}
		llmReasons := append([]string{}, spec.Reasons...)
		if spec.Need != "" && spec.Need != "fine" {
			llmReasons = append(llmReasons, "the student says they feel: "+spec.Need+" (respond to this need in a caring, non-clinical way)")
		}
		llmRes, err := p.llm.GenerateAdaptedActivity(ctx, orig, string(spec.InterventionType), llmReasons, state)
		if err == nil && llmRes != nil {
			RecordLLMOutcome("ai", "")
			difficulty := orig.Difficulty
			switch spec.InterventionType {
			case models.InterventionScaffold:
				difficulty = "beginner"
			case models.InterventionIncreaseChallenge:
				difficulty = "advanced"
			}
			reason := llmRes.PersonalizationReason
			if reason == "" {
				reason = fmt.Sprintf("Adapted based on %s state for targeted learning.", spec.InterventionType)
			}
			return &models.PersonalizedActivity{
				ID:                    fmt.Sprintf("gen-%d", time.Now().UnixNano()),
				SourceActivityID:      orig.ID,
				Title:                 llmRes.Title,
				Content:               llmRes.Content,
				Difficulty:            difficulty,
				InterventionType:      spec.InterventionType,
				PersonalizationReason: reason,
				IsGenerated:           true,
			}, nil
		}
		// Any LLM failure falls back to the deterministic planner below.
		msg := "invalid LLM response"
		if err != nil {
			msg = err.Error()
		}
		RecordLLMOutcome("rules", msg)
	} else {
		RecordLLMOutcome("rules", "")
	}

	interv := &models.AttemptIntervention{
		InterventionType: spec.InterventionType,
	}
	if len(spec.Reasons) > 0 {
		interv.Reason = spec.Reasons[0]
	}

	var catalog []models.Activity
	if p.repo != nil {
		if all, err := p.repo.ListActivities(); err == nil {
			catalog = all
		}
	}

	planned := p.PlanActivity(&models.EngagementAssessment{State: models.EngagementState(attempt.Status)}, interv, orig, catalog)
	if twist, why := NeedTwist(spec.Need); twist != "" {
		planned.Content += "\n\n" + twist
		planned.PersonalizationReason = strings.TrimSpace(planned.PersonalizationReason + " " + why)
	}
	return &planned, nil
}

func isHigherDifficulty(target, current string) bool {
	weights := map[string]int{"beginner": 1, "intermediate": 2, "advanced": 3}
	return weights[strings.ToLower(target)] > weights[strings.ToLower(current)]
}

func isLowerDifficulty(target, current string) bool {
	weights := map[string]int{"beginner": 1, "intermediate": 2, "advanced": 3}
	return weights[strings.ToLower(target)] < weights[strings.ToLower(current)]
}
