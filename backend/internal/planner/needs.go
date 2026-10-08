package planner

// NeedTwist returns a short, domain-agnostic prompt that responds to the
// student's self-reported need (connection, play, curiosity, ...). It never
// adds subject-specific constraints, so it is safe for any activity type.
func NeedTwist(need string) (twist string, why string) {
	switch need {
	case "lonely":
		return "🤝 Team up: explain your answer to a friend, classmate or family member, then ask how they would approach it. Note one idea you hadn't thought of.",
			"You said you wanted some company, so we added a way to learn with someone else."
	case "playful":
		return "🎮 Play mode: turn this into a 60-second challenge, then beat your own time, or turn it into a riddle for someone else.",
			"You felt like playing, so we added a game twist."
	case "curious":
		return "🔭 Go deeper: write one \"what if…?\" question about this topic and make a guess at the answer.",
			"You were feeling curious, so we added a question to explore."
	case "stuck":
		return "🧭 Being stuck is normal: say which part is confusing, then work on just that one part first.",
			"You said you felt stuck, so we added a way to break it down."
	case "bored":
		return "🚀 Make it yours: add your own twist or a harder rule, and see if you can still solve it.",
			"You said you were bored, so we added a way to make it your own."
	}
	return "", ""
}
