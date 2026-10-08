import React, { useState, useEffect, useCallback } from "react";
import {
  Activity,
  Attempt,
  PersonalizedActivityData,
  AdaptationLoopData,
  DemoResult,
  ActivityRunnerState,
  UserRole,
} from "@/types";
import { AppShell } from "@/components/layout/AppShell";
import { ActivityLibrary } from "@/components/activities/ActivityLibrary";
import { ActivityModal } from "@/components/activities/ActivityModal";
import { Dashboard } from "@/components/dashboard/Dashboard";
import { CheckIn } from "@/components/checkin/CheckIn";
import { LlmSettings } from "@/components/settings/LlmSettings";
import { AdaptationHistory } from "@/components/adaptation/AdaptationHistory";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useActivityTracker } from "@/hooks/useActivityTracker";

const API_BASE = "http://localhost:8080/api";
const STUDENT_ID = "student-demo-101";

export default function App() {
  const [role, setRole] = useState<UserRole>("student");
  const [currentTab, setCurrentTab] = useState<string>("dashboard");

  // Dark mode state management
  const [isDark, setIsDark] = useState<boolean>(() => {
    return (
      localStorage.getItem("engage-theme") === "dark" ||
      (!localStorage.getItem("engage-theme") &&
        window.matchMedia("(prefers-color-scheme: dark)").matches)
    );
  });

  useEffect(() => {
    if (isDark) {
      document.documentElement.classList.add("dark");
      localStorage.setItem("engage-theme", "dark");
    } else {
      document.documentElement.classList.remove("dark");
      localStorage.setItem("engage-theme", "light");
    }
  }, [isDark]);

  const toggleTheme = () => setIsDark((prev) => !prev);

  // State machine & Runner
  const [runnerState, setRunnerState] = useState<ActivityRunnerState>("AVAILABLE");
  const [activities, setActivities] = useState<Activity[]>([]);
  const [selectedActivity, setSelectedActivity] = useState<Activity | null>(null);
  const [currentAttempt, setCurrentAttempt] = useState<Attempt | null>(null);
  const [activeParentAttemptId, setActiveParentAttemptId] = useState<string | null>(null);
  const [studentInput, setStudentInput] = useState<string>("");

  // Adaptation & loop data
  const [personalizedNext, setPersonalizedNext] = useState<PersonalizedActivityData | null>(null);
  const [latestLoop, setLatestLoop] = useState<AdaptationLoopData | null>(null);
  const [allLoops, setAllLoops] = useState<AdaptationLoopData[]>([]);

  // Author modal
  const [isActivityModalOpen, setIsActivityModalOpen] = useState(false);
  const [activityToEdit, setActivityToEdit] = useState<Activity | null>(null);

  // Status alerts & telemetry
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [usingCache, setUsingCache] = useState<boolean>(false);
  const [feedbackNotice, setFeedbackNotice] = useState<string | null>(null);

  const { flushTelemetry, recordInteraction } = useActivityTracker({
    studentId: STUDENT_ID,
    attemptId: currentAttempt?.id ?? null,
    isActive: runnerState === "ACTIVE",
    apiBase: API_BASE,
    flushIntervalMs: 15000,
  });

  const loadActivities = useCallback(async () => {
    setIsLoading(true);
    setErrorMessage(null);
    try {
      const response = await fetch(`${API_BASE}/activities`, {
        headers: { Accept: "application/json" },
      });
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      const data = await response.json();
      const list: Activity[] = Array.isArray(data) ? data : [];
      setActivities(list);
      setUsingCache(false);
      try {
        localStorage.setItem("engage-activities", JSON.stringify(list));
      } catch {
        /* storage unavailable */
      }
    } catch (err: unknown) {
      let cached: Activity[] = [];
      try {
        cached = JSON.parse(localStorage.getItem("engage-activities") || "[]");
      } catch {
        cached = [];
      }
      if (cached.length > 0) {
        setActivities(cached);
        setUsingCache(true);
      } else {
        setErrorMessage(err instanceof Error ? err.message : "Failed to load activities");
      }
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    loadActivities();
  }, [loadActivities]);

  const loadLoops = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE}/loops?student_id=${STUDENT_ID}`, {
        headers: { Accept: "application/json" },
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      setAllLoops(Array.isArray(data) ? data : []);
    } catch (err) {
      console.warn("Failed to load adaptation history:", err);
    }
  }, []);

  useEffect(() => {
    loadLoops();
  }, [loadLoops]);

  const triggerAdaptationPipeline = async (attemptId: string) => {
    try {
      await fetch(`${API_BASE}/attempts/${attemptId}/evaluate`, { method: "POST" });
      await fetch(`${API_BASE}/attempts/${attemptId}/intervene`, { method: "POST" });
      const pRes = await fetch(`${API_BASE}/attempts/${attemptId}/personalize`, { method: "POST" });
      if (pRes.ok) {
        const pData: PersonalizedActivityData = await pRes.json();
        setPersonalizedNext({ ...pData, parent_attempt_id: attemptId });
      }
    } catch (err) {
      console.warn("Pipeline error:", err);
    }
  };

  const [demoResult, setDemoResult] = useState<DemoResult | null>(null);
  const [isDemoRunning, setIsDemoRunning] = useState(false);
  const [demoError, setDemoError] = useState<string | null>(null);

  const runDemo = async (scenario: string, need: string) => {
    if (!activities[0]) return;
    setIsDemoRunning(true);
    setDemoError(null);
    try {
      const res = await fetch(`${API_BASE}/demo/run`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ activity_id: activities[0].id, student_id: STUDENT_ID, scenario, need }),
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data: DemoResult = await res.json();
      setDemoResult(data);
      setLatestLoop(data.loop);
      loadLoops();
    } catch (err: unknown) {
      setDemoError(err instanceof Error ? `Demo failed: ${err.message}` : "Demo failed");
    } finally {
      setIsDemoRunning(false);
    }
  };

  const triggerCloseLoop = async (origAttemptId: string) => {
    try {
      const res = await fetch(`${API_BASE}/attempts/${origAttemptId}/close-loop`, {
        method: "POST",
      });
      if (res.ok) {
        const data: AdaptationLoopData = await res.json();
        setLatestLoop(data);
        setAllLoops((prev) => [data, ...prev.filter((l) => l.id !== data.id)]);
      }
    } catch (err) {
      console.warn("Close loop error:", err);
    }
  };

  const handleStartActivity = async (activity: Activity, parentAttemptId?: string | null) => {
    setRunnerState("STARTING");
    setErrorMessage(null);
    setFeedbackNotice(null);
    setCurrentTab("activities");

    try {
      const payload: { student_id: string; parent_attempt_id?: string } = {
        student_id: STUDENT_ID,
      };
      if (parentAttemptId) payload.parent_attempt_id = parentAttemptId;

      const response = await fetch(`${API_BASE}/activities/${activity.id}/start`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify(payload),
      });

      if (!response.ok) throw new Error(`HTTP ${response.status}`);

      const attemptData: Attempt = await response.json();
      setSelectedActivity(activity);
      setCurrentAttempt(attemptData);
      setActiveParentAttemptId(parentAttemptId || null);
      setStudentInput("");
      setRunnerState("ACTIVE");
    } catch (err: unknown) {
      setErrorMessage(err instanceof Error ? err.message : "Failed to start activity");
      setRunnerState("AVAILABLE");
    }
  };

  const handleStartPersonalizedActivity = async () => {
    if (!personalizedNext) return;
    const parentId = personalizedNext.parent_attempt_id;
    const syntheticActivity: Activity = {
      id: personalizedNext.source_activity_id,
      title: personalizedNext.title,
      description: personalizedNext.personalization_reason,
      content: personalizedNext.content,
      difficulty: personalizedNext.difficulty,
    };
    setPersonalizedNext(null);
    await handleStartActivity(syntheticActivity, parentId);
  };

  const handleCompleteActivity = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedActivity || !currentAttempt) return;

    if (!studentInput.trim()) {
      setErrorMessage("Please enter your answer or work before submitting.");
      return;
    }

    setRunnerState("COMPLETING");
    setErrorMessage(null);
    await flushTelemetry();

    const completedAttemptId = currentAttempt.id;
    const completedTitle = selectedActivity.title;
    const parentId = activeParentAttemptId;

    try {
      const response = await fetch(`${API_BASE}/activities/${selectedActivity.id}/complete`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify({
          student_id: STUDENT_ID,
          attempt_id: completedAttemptId,
          response: studentInput.trim(),
        }),
      });

      if (!response.ok) throw new Error(`HTTP ${response.status}`);

      setFeedbackNotice(`Submitted "${completedTitle}".`);
      setSelectedActivity(null);
      setCurrentAttempt(null);
      setActiveParentAttemptId(null);
      setStudentInput("");
      setRunnerState("AVAILABLE");

      if (parentId) {
        await triggerCloseLoop(parentId);
      } else {
        await triggerAdaptationPipeline(completedAttemptId);
      }
    } catch (err: unknown) {
      setErrorMessage(err instanceof Error ? err.message : "Submission failed");
      setRunnerState("ACTIVE");
    }
  };

  const handleAbandonActivity = async () => {
    if (!selectedActivity || !currentAttempt) return;

    setRunnerState("ABANDONING");
    setErrorMessage(null);
    await flushTelemetry();

    const abandonedAttemptId = currentAttempt.id;
    const abandonedTitle = selectedActivity.title;
    const parentId = activeParentAttemptId;

    try {
      await fetch(`${API_BASE}/activities/${selectedActivity.id}/abandon`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify({ student_id: STUDENT_ID, attempt_id: abandonedAttemptId }),
      });
      setFeedbackNotice(`Session ended for "${abandonedTitle}".`);
    } catch (err: unknown) {
      console.warn("Abandon error:", err);
    } finally {
      setSelectedActivity(null);
      setCurrentAttempt(null);
      setActiveParentAttemptId(null);
      setStudentInput("");
      setRunnerState("AVAILABLE");

      if (parentId) {
        await triggerCloseLoop(parentId);
      } else {
        await triggerAdaptationPipeline(abandonedAttemptId);
      }
    }
  };

  const handleOpenCreateModal = () => {
    setActivityToEdit(null);
    setIsActivityModalOpen(true);
  };

  const handleOpenEditModal = (activity: Activity) => {
    setActivityToEdit(activity);
    setIsActivityModalOpen(true);
  };

  return (
    <AppShell
      currentTab={currentTab}
      setCurrentTab={setCurrentTab}
      role={role}
      setRole={setRole}
      studentId={STUDENT_ID}
      isSessionActive={runnerState === "ACTIVE"}
      isDark={isDark}
      onToggleTheme={toggleTheme}
    >
      {/* Alert Banners */}
      {errorMessage && (
        <div className="mb-4 p-3 bg-rose-50 dark:bg-rose-950/60 border border-rose-200 dark:border-rose-900/60 text-rose-800 dark:text-rose-200 rounded-md text-xs flex justify-between items-center">
          <span>{errorMessage}</span>
          <button onClick={() => setErrorMessage(null)} className="font-bold ml-2 text-rose-600 dark:text-rose-400">✕</button>
        </div>
      )}

      {feedbackNotice && (
        <div className="mb-4 p-3 bg-emerald-50 dark:bg-emerald-950/60 border border-emerald-200 dark:border-emerald-900/60 text-emerald-800 dark:text-emerald-200 rounded-md text-xs flex justify-between items-center">
          <span>{feedbackNotice}</span>
          <button onClick={() => setFeedbackNotice(null)} className="font-bold ml-2 text-emerald-600 dark:text-emerald-400">✕</button>
        </div>
      )}

      {/* Step 8 Closed-Loop Verification Result Banner */}
      {latestLoop && runnerState === "AVAILABLE" && (
        <Card className="mb-6 border-emerald-300 dark:border-emerald-800 bg-emerald-50/80 dark:bg-emerald-950/40 shadow-sm border-2">
          <CardHeader className="pb-2 flex flex-row items-center justify-between">
            <div>
              <span className="text-[10px] tracking-wider uppercase font-bold px-2 py-0.5 rounded bg-emerald-200 dark:bg-emerald-900 text-emerald-800 dark:text-emerald-200 inline-block mb-1">
                ✓ Adaptation Feedback Cycle Closed
              </span>
              <CardTitle className="text-base text-slate-900 dark:text-slate-100 font-bold">
                Intervention Result: {latestLoop.effectiveness}
              </CardTitle>
            </div>
            <button
              onClick={() => setLatestLoop(null)}
              className="text-xs text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 font-bold"
            >
              ✕
            </button>
          </CardHeader>
          <CardContent className="space-y-3 text-xs">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-2 bg-white/95 dark:bg-slate-900 p-3 rounded border border-emerald-100 dark:border-emerald-900/60 text-center">
              <div>
                <div className="text-slate-400 dark:text-slate-500 text-[10px] uppercase font-bold">Initial State</div>
                <div className="font-semibold text-slate-800 dark:text-slate-200 mt-0.5">{latestLoop.pre_state}</div>
              </div>
              <div>
                <div className="text-slate-400 dark:text-slate-500 text-[10px] uppercase font-bold">Post-Adaptation</div>
                <div className="font-semibold text-slate-800 dark:text-slate-200 mt-0.5">{latestLoop.post_state}</div>
              </div>
              <div>
                <div className="text-slate-400 dark:text-slate-500 text-[10px] uppercase font-bold">Delta</div>
                <div
                  className={`font-semibold mt-0.5 ${
                    latestLoop.engagement_delta >= 0 ? "text-emerald-600 dark:text-emerald-400" : "text-rose-600 dark:text-rose-400"
                  }`}
                >
                  {latestLoop.engagement_delta >= 0 ? "+" : ""}
                  {(latestLoop.engagement_delta * 100).toFixed(0)}%
                </div>
              </div>
              <div>
                <div className="text-slate-400 dark:text-slate-500 text-[10px] uppercase font-bold">Outcome</div>
                <div className="font-bold text-slate-800 dark:text-slate-100 mt-0.5">{latestLoop.effectiveness}</div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Step 6 Personalized Activity Banner */}
      {personalizedNext && runnerState === "AVAILABLE" && (
        <Card className="mb-6 border-indigo-300 dark:border-indigo-800 bg-indigo-50/80 dark:bg-indigo-950/40 shadow-sm border-2">
          <CardHeader className="pb-2 flex flex-row items-center justify-between">
            <div>
              <span className="text-[10px] tracking-wider uppercase font-bold px-2 py-0.5 rounded bg-indigo-200 dark:bg-indigo-900 text-indigo-800 dark:text-indigo-200 inline-block mb-1">
                ✨ Learning Path Adapted ({personalizedNext.intervention_type})
              </span>
              <CardTitle className="text-base text-slate-900 dark:text-slate-100 font-bold">
                {personalizedNext.title}
              </CardTitle>
            </div>
            <button
              onClick={() => setPersonalizedNext(null)}
              className="text-xs text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 font-bold"
            >
              ✕
            </button>
          </CardHeader>
          <CardContent className="space-y-3 text-xs">
            <div className="bg-white/95 dark:bg-slate-900 p-3 rounded border border-indigo-100 dark:border-indigo-900/60">
              <span className="font-semibold text-slate-800 dark:text-slate-200 block">Why this was chosen:</span>
              <p className="text-slate-600 dark:text-slate-400 mt-0.5">{personalizedNext.personalization_reason}</p>
            </div>
            <Button
              onClick={handleStartPersonalizedActivity}
              className="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs px-4 py-2"
            >
              Start Adapted Activity →
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Active Runner Session */}
      {runnerState === "ACTIVE" && selectedActivity ? (
        <Card className="border border-slate-300 dark:border-slate-800 shadow-sm bg-white dark:bg-slate-900">
          <CardHeader className="border-b border-slate-100 dark:border-slate-800 pb-4">
            <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-2">
              <div>
                <div className="flex items-center gap-2">
                  <CardTitle className="text-xl font-bold text-slate-900 dark:text-slate-100">
                    {selectedActivity.title}
                  </CardTitle>
                  {activeParentAttemptId && (
                    <span className="text-[10px] font-bold px-2 py-0.5 rounded bg-purple-100 dark:bg-purple-950 text-purple-800 dark:text-purple-300 border border-purple-200 dark:border-purple-800 uppercase">
                      Adapted Loop Active
                    </span>
                  )}
                </div>
                <div className="flex gap-4 text-xs text-slate-500 dark:text-slate-400 mt-1">
                  <span>Topic: {selectedActivity.topic || "Core Curriculum"}</span>
                  <span>Attempt: <b className="font-mono">{currentAttempt?.id.slice(0, 8)}</b></span>
                </div>
              </div>
              <span className="text-xs font-semibold px-2.5 py-1 rounded-full bg-blue-50 dark:bg-blue-950 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800">
                ● Live Telemetry Active
              </span>
            </div>
          </CardHeader>
          <CardContent className="pt-6 space-y-6">
            <div>
              <label className="text-xs font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500 block mb-1">
                Activity Prompt
              </label>
              <div className="bg-slate-50 dark:bg-slate-800/70 p-4 rounded-md border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-200 text-sm whitespace-pre-wrap">
                {selectedActivity.content}
              </div>
            </div>

            <form onSubmit={handleCompleteActivity} className="space-y-4">
              <div>
                <label className="text-xs font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500 block mb-1">
                  Your Response / Solution
                </label>
                <textarea
                  value={studentInput}
                  onChange={(e) => {
                    setStudentInput(e.target.value);
                    recordInteraction();
                  }}
                  rows={4}
                  placeholder="Type your response here..."
                  className="w-full text-sm p-3 border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-950 text-slate-900 dark:text-slate-100 placeholder-slate-400 dark:placeholder-slate-500 rounded focus:ring-2 focus:ring-indigo-600 focus:outline-none"
                />
              </div>

              <div className="flex justify-between items-center pt-2">
                <Button
                  type="button"
                  onClick={handleAbandonActivity}
                  className="bg-transparent text-rose-600 dark:text-rose-400 border border-rose-300 dark:border-rose-800 hover:bg-rose-50 dark:hover:bg-rose-950/40 text-xs px-4"
                >
                  Abandon / Exit
                </Button>
                <Button
                  type="submit"
                  className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs px-5 font-semibold"
                >
                  Submit & Complete
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      ) : (
        <>
          {usingCache && (
            <div role="status" className="rounded-md border border-amber-200 bg-amber-50 px-4 py-2 text-sm text-amber-800">
              You're offline or the server can't be reached. Showing your saved activities; progress can't be
              recorded until you reconnect.{" "}
              <button className="underline font-medium" onClick={loadActivities}>Retry</button>
            </div>
          )}

          {/* Main Views */}
          {currentTab === "dashboard" && (
            <div className="space-y-6">
            <Dashboard
              activitiesCount={activities.length}
              latestLoop={latestLoop}
              recommendedActivity={activities[0] || null}
              onStartRecommended={(act) => handleStartActivity(act, null)}
              onNavigateToActivities={() => setCurrentTab("activities")}
              onRunDemo={runDemo}
              isDemoRunning={isDemoRunning}
              demoResult={demoResult}
              demoError={demoError}
            />
            <CheckIn apiBase={API_BASE} studentId={STUDENT_ID} />
            <LlmSettings apiBase={API_BASE} refreshKey={demoResult} />
            </div>
          )}

          {currentTab === "activities" && (
            <ActivityLibrary
              activities={activities}
              isLoading={isLoading}
              role={role}
              onStartActivity={(act) => handleStartActivity(act, null)}
              onEditActivity={handleOpenEditModal}
              onCreateNew={handleOpenCreateModal}
              onRefresh={loadActivities}
            />
          )}

          {currentTab === "create" && (
            <ActivityLibrary
              activities={activities}
              isLoading={isLoading}
              role={role}
              onStartActivity={(act) => handleStartActivity(act, null)}
              onEditActivity={handleOpenEditModal}
              onCreateNew={handleOpenCreateModal}
              onRefresh={loadActivities}
            />
          )}

          {currentTab === "history" && (
            <AdaptationHistory loops={allLoops} onRefresh={loadLoops} />
          )}
        </>
      )}

      {/* Activity Create/Edit Modal */}
      <ActivityModal
        isOpen={isActivityModalOpen}
        onClose={() => setIsActivityModalOpen(false)}
        onSaved={loadActivities}
        activityToEdit={activityToEdit}
        apiBase={API_BASE}
      />
    </AppShell>
  );
}
