import React, { useCallback, useEffect, useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { NEEDS } from "@/lib/needs";

interface Props {
  apiBase: string;
  studentId: string;
}

export const CheckIn: React.FC<Props> = ({ apiBase, studentId }) => {
  const [chosen, setChosen] = useState<string | null>(null);
  const [counts, setCounts] = useState<Record<string, number>>({});
  const [error, setError] = useState<string | null>(null);

  const loadPulse = useCallback(async () => {
    try {
      const res = await fetch(`${apiBase}/community/pulse`);
      if (res.ok) setCounts((await res.json()).counts ?? {});
    } catch {
      /* offline: pulse is optional */
    }
  }, [apiBase]);

  useEffect(() => {
    loadPulse();
  }, [loadPulse]);

  const choose = async (need: string) => {
    setError(null);
    try {
      const res = await fetch(`${apiBase}/students/${studentId}/checkin`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ need }),
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      setChosen(need);
      loadPulse();
    } catch (e) {
      setError(e instanceof Error ? `Could not save check-in (${e.message})` : "Could not save check-in");
    }
  };

  const pulse = NEEDS.filter((n) => (counts[n.key] ?? 0) > 0);

  return (
    <Card className="border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900">
      <CardHeader className="pb-2">
        <CardTitle className="text-base">How are you feeling about learning today?</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex flex-wrap gap-2" role="group" aria-label="Check-in">
          {NEEDS.map((n) => (
            <button
              key={n.key}
              onClick={() => choose(n.key)}
              aria-pressed={chosen === n.key}
              className={`rounded-full border px-3 py-1.5 text-sm transition-colors ${
                chosen === n.key
                  ? "bg-indigo-600 text-white border-indigo-600"
                  : "bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 border-slate-300 dark:border-slate-600 hover:bg-slate-50 dark:hover:bg-slate-700"
              }`}
            >
              <span aria-hidden="true">{n.emoji}</span> {n.label}
            </button>
          ))}
        </div>
        {chosen && (
          <p className="text-sm text-emerald-700 dark:text-emerald-400">
            Thanks. Your next adapted activity will take this into account.
          </p>
        )}
        {error && <p className="text-sm text-rose-600 dark:text-rose-400">{error}</p>}
        <div className="border-t border-slate-100 dark:border-slate-800 pt-3">
          <p className="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide">Community pulse (last 24h, anonymous)</p>
          {pulse.length === 0 ? (
            <p className="text-sm text-slate-400 dark:text-slate-500 mt-1">No check-ins yet. Be the first to share how you feel.</p>
          ) : (
            <div className="mt-2 flex flex-wrap gap-2 text-sm text-slate-700 dark:text-slate-300">
              {pulse.map((n) => (
                <span key={n.key} className="rounded bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 px-2 py-0.5">
                  {n.emoji} {counts[n.key]} {n.label.toLowerCase()}
                </span>
              ))}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
};
