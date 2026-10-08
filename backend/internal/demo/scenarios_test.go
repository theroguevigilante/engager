package demo

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/adaptive-learning/backend/internal/assessment"
	"github.com/adaptive-learning/backend/internal/intervention"
	"github.com/adaptive-learning/backend/internal/loop"
	"github.com/adaptive-learning/backend/internal/models"
)

func assess(tel Telemetry) (*models.EngagementAssessment, *models.ActivityAttempt) {
	att := &models.ActivityAttempt{Status: tel.Status, AttemptCount: tel.Attempts, Duration: int64(tel.Session)}
	sig := func(k string, v int) models.EngagementSignal {
		return models.EngagementSignal{SignalType: k, SignalValue: strconv.Itoa(v)}
	}
	m := assessment.ComputeMetrics([]models.EngagementSignal{
		sig("session_duration", tel.Session), sig("active_duration", tel.Active), sig("idle_duration", tel.Idle),
		sig("tab_switches", tel.TabSwitches), sig("interaction_count", tel.Interactions),
	}, *att)
	st, conf, reasons := assessment.Classify(m)
	return &models.EngagementAssessment{State: st, Confidence: conf, Metrics: m, Reasons: reasons}, att
}

func TestScenariosProduceDistinctStoriesAndOutcomes(t *testing.T) {
	want := map[string][2]string{
		"struggling": {"STRUGGLING", "SCAFFOLD"},
		"bored":      {"BORED", "INCREASE_CHALLENGE"},
		"drifting":   {"DISENGAGED", "RE_ENGAGE"},
	}
	for _, sc := range Scenarios {
		pre, preAtt := assess(sc.Pre)
		iv, _, _ := intervention.SelectIntervention(pre, preAtt)
		post, postAtt := assess(sc.Post)
		eff, _, _, _, _ := loop.EvaluateLoop(pre, post, postAtt, iv)
		fmt.Println(sc.Key, pre.State, iv, "->", post.State, eff)
		if string(pre.State) != want[sc.Key][0] || string(iv) != want[sc.Key][1] {
			t.Errorf("%s: got %s/%s want %v", sc.Key, pre.State, iv, want[sc.Key])
		}
		if post.State != models.StateEngaged || eff != models.EffectivenessEffective {
			t.Errorf("%s: post %s eff %s", sc.Key, post.State, eff)
		}
	}
}
