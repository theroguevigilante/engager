import React from "react";
import { UserRole } from "@/types";

interface AppShellProps {
  currentTab: string;
  setCurrentTab: (tab: string) => void;
  role: UserRole;
  setRole: (role: UserRole) => void;
  studentId: string;
  isSessionActive: boolean;
  isDark: boolean;
  onToggleTheme: () => void;
  children: React.ReactNode;
}

export const AppShell: React.FC<AppShellProps> = ({
  currentTab,
  setCurrentTab,
  role,
  setRole,
  studentId,
  isSessionActive,
  isDark,
  onToggleTheme,
  children,
}) => {
  const studentNav = [
    { id: "dashboard", label: "Dashboard", icon: "📊" },
    { id: "activities", label: "Activities", icon: "📚" },
    { id: "history", label: "Adaptation History", icon: "🔄" },
  ];

  const authorNav = [
    { id: "dashboard", label: "Overview", icon: "📈" },
    { id: "activities", label: "Activity Library", icon: "🗂️" },
    { id: "create", label: "+ Create Activity", icon: "✏️" },
    { id: "history", label: "Adaptation Insights", icon: "🔍" },
  ];

  const navItems = role === "student" ? studentNav : authorNav;

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-950 flex flex-col md:flex-row font-sans text-slate-800 dark:text-slate-100 transition-colors duration-200">
      {/* Sidebar Navigation */}
      <aside className="w-full md:w-64 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col shrink-0">
        <div className="p-5 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <div className="w-9 h-9 rounded-xl bg-gradient-to-tr from-indigo-600 to-violet-500 flex items-center justify-center text-white font-black text-lg shadow-md shadow-indigo-500/20">
              E
            </div>
            <div>
              <span className="font-bold text-slate-900 dark:text-white tracking-tight text-lg leading-tight block">
                Engage
              </span>
              <span className="text-[10px] text-slate-400 dark:text-slate-500 font-semibold uppercase tracking-wider block">
                Adaptive Learning
              </span>
            </div>
          </div>

          {/* Theme Toggle Button */}
          <button
            onClick={onToggleTheme}
            type="button"
            title={isDark ? "Switch to light mode" : "Switch to dark mode"}
            className="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 text-xs transition-colors"
          >
            {isDark ? "☀️" : "🌙"}
          </button>
        </div>

        {/* Role Switcher */}
        <div className="p-3 mx-4 my-3 bg-slate-100 dark:bg-slate-800/60 rounded-xl flex items-center justify-between text-xs border border-slate-200/50 dark:border-slate-700/50">
          <span className="font-semibold text-slate-500 dark:text-slate-400">Mode:</span>
          <div className="inline-flex rounded-lg p-0.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 shadow-xs">
            <button
              onClick={() => setRole("student")}
              type="button"
              className={`px-3 py-1 rounded-md text-xs font-semibold transition-all ${
                role === "student"
                  ? "bg-indigo-600 text-white shadow-sm"
                  : "text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white"
              }`}
            >
              Student
            </button>
            <button
              onClick={() => setRole("author")}
              type="button"
              className={`px-3 py-1 rounded-md text-xs font-semibold transition-all ${
                role === "author"
                  ? "bg-indigo-600 text-white shadow-sm"
                  : "text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white"
              }`}
            >
              Author
            </button>
          </div>
        </div>

        {/* Navigation Links */}
        <nav className="flex-1 p-3 space-y-1">
          {navItems.map((item) => (
            <button
              key={item.id}
              type="button"
              disabled={isSessionActive && item.id !== "activities"}
              onClick={() => setCurrentTab(item.id)}
              className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all ${
                currentTab === item.id
                  ? "bg-indigo-50 dark:bg-indigo-950/60 text-indigo-700 dark:text-indigo-300 font-semibold shadow-xs"
                  : "text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-900 dark:hover:text-white"
              } ${isSessionActive && item.id !== "activities" ? "opacity-35 cursor-not-allowed" : ""}`}
            >
              <span className="text-base">{item.icon}</span>
              <span>{item.label}</span>
            </button>
          ))}
        </nav>

        {/* Profile Info */}
        <div className="p-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50">
          <div className="text-[10px] uppercase font-bold tracking-wider text-slate-400 dark:text-slate-500">
            {role === "student" ? "Active Learner" : "Course Creator"}
          </div>
          <div className="text-xs font-mono font-semibold text-slate-700 dark:text-slate-300 mt-0.5 truncate">
            {role === "student" ? studentId : "admin@engage.internal"}
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col min-h-screen overflow-y-auto">
        <main className="p-6 md:p-8 max-w-6xl w-full mx-auto">{children}</main>
      </div>
    </div>
  );
};
