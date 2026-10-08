import React from "react";
import { Activity, AdaptationLoopData, DemoResult } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { NEEDS, SCENARIOS } from "@/lib/needs";

interface DashboardProps {
  activitiesCount: number;
  latestLoop: AdaptationLoopData | null;
  recommendedActivity: Activity | null;
  onStartRecommended: (activity: Activity) => void;
  onNavigateToActivities: () => void;
  onRunDemo: (scenario: string, need: string) => void;
  isDemoRunning: boolean;
  demoResult: DemoResult | null;
  demoError: string | null;
}

export const Dashboard: React.FC<DashboardProps> = ({
  activitiesCount,
  latestLoop,
  recommendedActivity,
  onStartRecommended,
  onNavigateToActivities,
  onRunDemo,
  isDemoRunning,
  demoResult,
  demoError,
}) => {
  const [scenario, setScenario] = React.useState<string>("struggling");
  const [need, setNeed] = React.useState<string>("");
  const story = SCENARIOS.find((x) => x.key === scenario)?.story;
  return (
    <div className="space-y-6">
      {/* One-click live demo of the closed adaptation loop */}
      <Card className="border border-indigo-200 dark:border-indigo-800 bg-white dark:bg-slate-900">
        <CardHeader className="pb-2">
          <CardTitle className="text-base">See the adaptation loop in action</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-slate-600 dark:text-slate-400">
            Simulates a student who struggles with an activity, then shows how Engage adapts the
            next one and measures whether it helped.
          </p>
          <div className="grid gap-3 md:grid-cols-2">
            <label className="text-xs font-medium text-slate-600 dark:text-slate-400 space-y-1 block">
              Student story
              <select
                className="w-full rounded-md border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 px-3 py-2 text-sm text-slate-900 dark:text-slate-100"
                value={scenario}
                onChange={(e) => setScenario(e.target.value)}
              >
                {SCENARIOS.map((x) => (
                  <option key={x.key} value={x.key}>{x.label}</option>
                ))}
              </select>
              <span className="block text-slate-400 dark:text-slate-500 font-normal">{story}</span>
            </label>
            <label className="text-xs font-medium text-slate-600 dark:text-slate-400 space-y-1 block">
              Student says they feel (optional)
              <select
                className="w-full rounded-md border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 px-3 py-2 text-sm text-slate-900 dark:text-slate-100"
                value={need}
                onChange={(e) => setNeed(e.target.value)}
              >
                <option value="">No check-in</option>
                {NEEDS.map((n) => (
                  <option key={n.key} value={n.key}>{n.emoji} {n.label}</option>
                ))}
              </select>
            </label>
          </div>
          <Button onClick={() => onRunDemo(scenario, need)} disabled={isDemoRunning}>
            {isDemoRunning ? "Running demo…" : "Run live demo"}
          </Button>
          {demoError && <p className="text-sm text-rose-600 dark:text-rose-400">{demoError}</p>}
          {demoResult && (
            <div className="rounded-lg border border-emerald-200 dark:border-emerald-800 bg-emerald-50 dark:bg-emerald-950/50 p-3 text-sm text-emerald-800 dark:text-emerald-200 space-y-1">
              <p className="font-semibold">Demo complete</p>
              <p>Initial state: {demoResult.pre_state} → Intervention: {demoResult.intervention_type} → Post state: {demoResult.post_state}</p>
              <p>Effectiveness: {demoResult.effectiveness}</p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Latest closed loop summary */}
      {latestLoop && (
        <Card className="border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900">
          <CardHeader className="pb-2">
            <CardTitle className="text-base">Latest adaptation cycle</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm text-slate-600 dark:text-slate-400">
            <p>Pre: {latestLoop.pre_state} ({(latestLoop.pre_engagement_score * 100).toFixed(0)}%) → Post: {latestLoop.post_state} ({(latestLoop.post_engagement_score * 100).toFixed(0)}%)</p>
            <p>Intervention: {latestLoop.intervention_type} · Effectiveness: {latestLoop.effectiveness}</p>
          </CardContent>
        </Card>
      )}

      {/* Recommended activity */}
      {recommendedActivity && (
        <Card className="border border-indigo-200 dark:border-indigo-800 bg-white dark:bg-slate-900">
          <CardHeader className="pb-2">
            <CardTitle className="text-base">Recommended for you</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <p className="font-semibold text-slate-900 dark:text-white">{recommendedActivity.title}</p>
            <p className="text-sm text-slate-600 dark:text-slate-400">{recommendedActivity.description}</p>
            <div className="flex gap-2">
              <Button onClick={() => onStartRecommended(recommendedActivity)}>Start activity</Button>
              <Button
                onClick={onNavigateToActivities}
                className="bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-600 shadow-none hover:bg-slate-50 dark:hover:bg-slate-700"
              >
                Browse all
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Quick stats */}
      <Card className="border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900">
        <CardHeader className="pb-2">
          <CardTitle className="text-base">Platform snapshot</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm text-slate-600 dark:text-slate-400">
          <p>Total activities in library: {activitiesCount}</p>
          <p>Adaptation cycles closed: {latestLoop ? "1+" : "0"}</p>
        </CardContent>
      </Card>
    </div>
  );
};
