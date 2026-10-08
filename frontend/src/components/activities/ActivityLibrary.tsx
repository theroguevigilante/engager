import React, { useState, useMemo } from "react";
import { Activity, UserRole } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface ActivityLibraryProps {
  activities: Activity[];
  isLoading: boolean;
  role: UserRole;
  onStartActivity: (activity: Activity) => void;
  onEditActivity: (activity: Activity) => void;
  onCreateNew: () => void;
  onRefresh: () => void;
}

export const ActivityLibrary: React.FC<ActivityLibraryProps> = ({
  activities,
  isLoading,
  role,
  onStartActivity,
  onEditActivity,
  onCreateNew,
  onRefresh,
}) => {
  const [searchTerm, setSearchTerm] = useState("");
  const [difficultyFilter, setDifficultyFilter] = useState("all");

  const filtered = useMemo(() => {
    return activities.filter((act) => {
      const matchSearch =
        act.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        act.description?.toLowerCase().includes(searchTerm.toLowerCase());
      const matchDiff =
        difficultyFilter === "all" ||
        act.difficulty?.toLowerCase() === difficultyFilter.toLowerCase();
      return matchSearch && matchDiff;
    });
  }, [activities, searchTerm, difficultyFilter]);

  const getDifficultyBadge = (difficulty?: string) => {
    switch (difficulty?.toLowerCase()) {
      case "beginner":
        return "bg-emerald-50 dark:bg-emerald-950/50 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-800";
      case "intermediate":
        return "bg-amber-50 dark:bg-amber-950/50 text-amber-700 dark:text-amber-300 border-amber-200 dark:border-amber-800";
      case "advanced":
        return "bg-rose-50 dark:bg-rose-950/50 text-rose-700 dark:text-rose-300 border-rose-200 dark:border-rose-800";
      default:
        return "bg-slate-50 dark:bg-slate-800 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-700";
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">
            {role === "student" ? "Available Learning Activities" : "Curriculum Activity Library"}
          </h1>
          <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
            {role === "student"
              ? "Select an activity to begin your personalized learning session."
              : "Manage, create, and adapt modular activities for the recommendation engine."}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            onClick={onRefresh}
            className="bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 text-xs"
          >
            Refresh
          </Button>
          {role === "author" && (
            <Button
              onClick={onCreateNew}
              className="bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-semibold"
            >
              + Create Activity
            </Button>
          )}
        </div>
      </div>

      {/* Filter / Search Bar */}
      <div className="flex flex-col sm:flex-row gap-3">
        <input
          type="text"
          placeholder="Search activities by title or concept..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="flex-1 p-2 text-xs border border-slate-300 dark:border-slate-700 rounded bg-white dark:bg-slate-900 text-slate-900 dark:text-slate-100 placeholder-slate-400 dark:placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-600"
        />
        <select
          value={difficultyFilter}
          onChange={(e) => setDifficultyFilter(e.target.value)}
          className="p-2 text-xs border border-slate-300 dark:border-slate-700 rounded bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-600"
        >
          <option value="all">All Difficulties</option>
          <option value="beginner">Beginner</option>
          <option value="intermediate">Intermediate</option>
          <option value="advanced">Advanced</option>
        </select>
      </div>

      {/* Activities Grid */}
      {isLoading ? (
        <div className="p-16 text-center text-slate-400 dark:text-slate-500 text-sm">
          Loading modular activities...
        </div>
      ) : filtered.length === 0 ? (
        <Card className="border-dashed border-2 border-slate-200 dark:border-slate-800 bg-transparent">
          <CardContent className="p-12 text-center text-slate-500 dark:text-slate-400 space-y-3">
            <p className="text-sm font-medium">No learning activities found matching your criteria.</p>
            {role === "author" && (
              <Button
                onClick={onCreateNew}
                className="bg-indigo-600 hover:bg-indigo-700 text-white text-xs mt-2"
              >
                Create your first activity
              </Button>
            )}
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {filtered.map((activity) => (
            <Card
              key={activity.id}
              className="border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 hover:border-slate-300 dark:hover:border-slate-700 hover:shadow-xs transition-all flex flex-col justify-between"
            >
              <CardHeader className="pb-3">
                <div className="flex justify-between items-start gap-2">
                  <CardTitle className="text-base text-slate-900 dark:text-white font-semibold leading-snug">
                    {activity.title}
                  </CardTitle>
                  {activity.difficulty && (
                    <span
                      className={`text-[10px] uppercase font-bold px-2 py-0.5 rounded border shrink-0 ${getDifficultyBadge(
                        activity.difficulty
                      )}`}
                    >
                      {activity.difficulty}
                    </span>
                  )}
                </div>
                <div className="flex items-center gap-3 text-xs text-slate-400 dark:text-slate-500 mt-1">
                  <span>Topic: {activity.topic || "Core Curriculum"}</span>
                  <span>⏱ {activity.estimated_duration_minutes || 5} min</span>
                </div>
              </CardHeader>
              <CardContent className="space-y-4 pt-0">
                <p className="text-xs text-slate-600 dark:text-slate-400 line-clamp-2">
                  {activity.description || "Interactive exercise designed for active problem solving."}
                </p>
                <div className="flex items-center justify-end gap-2 pt-2 border-t border-slate-100 dark:border-slate-800">
                  {role === "author" && (
                    <Button
                      onClick={() => onEditActivity(activity)}
                      className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700 text-xs px-3"
                    >
                      Edit
                    </Button>
                  )}
                  <Button
                    onClick={() => onStartActivity(activity)}
                    className="bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-semibold px-4"
                  >
                    Start Session →
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
};
