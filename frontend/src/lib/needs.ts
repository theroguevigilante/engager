export const NEEDS = [
  { key: "stuck", label: "Stuck", emoji: "🧭" },
  { key: "bored", label: "Bored", emoji: "🥱" },
  { key: "curious", label: "Curious", emoji: "🔭" },
  { key: "lonely", label: "Want company", emoji: "🤝" },
  { key: "playful", label: "Want to play", emoji: "🎮" },
  { key: "fine", label: "I'm good", emoji: "🙂" },
] as const;

export const SCENARIOS = [
  { key: "struggling", label: "Struggling", story: "Works hard, retries, but can't crack it." },
  { key: "bored", label: "Bored", story: "Rushes through with no challenge." },
  { key: "drifting", label: "Drifting away", story: "Idle, switching tabs, barely interacting." },
] as const;
