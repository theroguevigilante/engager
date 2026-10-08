import React from "react";
import { AdaptationLoopData } from "@/types";
import { Card, CardHeader, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface AdaptationHistoryProps {
  loops: AdaptationLoopData[];
  onRefresh: () => void;
}

export const AdaptationHistory: React.FC<AdaptationHistoryProps> = ({ loops, onRefresh }) => {
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">
            Closed-Loop Adaptation History
          </h1>
          <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
            Inspection log showing pre-intervention states, adaptive interventions, and post-outcome deltas.
          </p>
        </div>
        <Button
          onClick={onRefresh}
          className="h-8 px-3 bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-600 shadow-none hover:bg-slate-50 dark:hover:bg-slate-700"
        >
          Refresh
        </Button>
      </div>

      {loops.length === 0 ? (
        <Card className="border-dashed border-2 border-slate-200 dark:border-slate-700">
          <CardContent className="p-12 text-center text-slate-500 dark:text-slate-400 space-y-2">
            <p className="text-sm font-medium">No adaptation cycles closed yet.</p>
            <p className="text-xs text-slate-400 dark:text-slate-500">
              Complete an activity, trigger an adaptive challenge or scaffold, and submit the adapted task to see the closed-loop comparison here.
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {loops.map((item) => (
            <Card key={item.id} className="border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900">
              <CardHeader className="pb-3 border-b border-slate-100 dark:border-slate-800">
                <div className="flex justify-between items-center">
                  <span className="text-xs font-medium text-slate-400 dark:text-slate-500">
                    {item.created_at ? new Date(item.created_at).toLocaleString() : "Adaptation cycle"}
                  </span>
                  <span
                    className={`text-[10px] uppercase font-bold px-2 py-0.5 rounded border ${
                      item.effectiveness === "EFFECTIVE"
                        ? "bg-emerald-50 dark:bg-emerald-950/50 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-800"
                        : item.effectiveness === "INEFFECTIVE"
                        ? "bg-rose-50 dark:bg-rose-950/50 text-rose-700 dark:text-rose-300 border-rose-200 dark:border-rose-800"
                        : "bg-amber-50 dark:bg-amber-950/50 text-amber-700 dark:text-amber-300 border-amber-200 dark:border-amber-800"
                    }`}
                  >
                    Outcome: {item.effectiveness}
                  </span>
                </div>
              </CardHeader>
              <CardContent className="pt-4">
                {/* Horizontal Flow Diagram */}
                <div className="flex flex-col md:flex-row items-center justify-between gap-3 text-center text-xs">
                  <div className="p-3 bg-slate-50 dark:bg-slate-800 rounded border border-slate-200 dark:border-slate-700 w-full md:w-1/4">
                    <span className="text-[10px] text-slate-400 dark:text-slate-500 font-bold uppercase block">
                      Initial State
                    </span>
                    <span className="font-bold text-slate-800 dark:text-slate-200 text-sm mt-0.5 block">
                      {item.pre_state}
                    </span>
                    <span className="text-[10px] text-slate-500 dark:text-slate-400">
                      Score: {(item.pre_engagement_score * 100).toFixed(0)}%
                    </span>
                  </div>

                  <span className="text-slate-400 dark:text-slate-500 font-bold hidden md:inline">→</span>

                  <div className="p-3 bg-indigo-50 dark:bg-indigo-950/50 rounded border border-indigo-200 dark:border-indigo-800 w-full md:w-1/4">
                    <span className="text-[10px] text-indigo-500 dark:text-indigo-400 font-bold uppercase block">
                      Engine Action
                    </span>
                    <span className="font-bold text-indigo-900 dark:text-indigo-200 text-sm mt-0.5 block">
                      ADAPTED
                    </span>
                    <span className="text-[10px] text-indigo-600 dark:text-indigo-400">
                      {item.intervention_type}
                    </span>
                  </div>

                  <span className="text-slate-400 dark:text-slate-500 font-bold hidden md:inline">→</span>

                  <div className="p-3 bg-slate-50 dark:bg-slate-800 rounded border border-slate-200 dark:border-slate-700 w-full md:w-1/4">
                    <span className="text-[10px] text-slate-400 dark:text-slate-500 font-bold uppercase block">
                      Post State
                    </span>
                    <span className="font-bold text-slate-800 dark:text-slate-200 text-sm mt-0.5 block">
                      {item.post_state}
                    </span>
                    <span className="text-[10px] text-slate-500 dark:text-slate-400">
                      Score: {(item.post_engagement_score * 100).toFixed(0)}%
                    </span>
                  </div>

                  <span className="text-slate-400 dark:text-slate-500 font-bold hidden md:inline">→</span>

                  <div
                    className={`p-3 rounded border w-full md:w-1/4 ${
                      item.effectiveness === "EFFECTIVE"
                        ? "bg-emerald-50 dark:bg-emerald-950/50 border-emerald-200 dark:border-emerald-800"
                        : item.effectiveness === "INEFFECTIVE"
                        ? "bg-rose-50 dark:bg-rose-950/50 border-rose-200 dark:border-rose-800"
                        : "bg-amber-50 dark:bg-amber-950/50 border-amber-200 dark:border-amber-800"
                    }`}
                  >
                    <span
                      className={`text-[10px] font-bold uppercase block ${
                        item.effectiveness === "EFFECTIVE"
                          ? "text-emerald-600 dark:text-emerald-400"
                          : item.effectiveness === "INEFFECTIVE"
                          ? "text-rose-600 dark:text-rose-400"
                          : "text-amber-600 dark:text-amber-400"
                      }`}
                    >
                      Outcome
                    </span>
                    <span
                      className={`font-bold text-sm mt-0.5 block ${
                        item.effectiveness === "EFFECTIVE"
                          ? "text-emerald-900 dark:text-emerald-200"
                          : item.effectiveness === "INEFFECTIVE"
                          ? "text-rose-900 dark:text-rose-200"
                          : "text-amber-900 dark:text-amber-200"
                      }`}
                    >
                      {item.effectiveness}
                    </span>
                  </div>
                </div>

                {/* Reasoning */}
                <div className="mt-4 pt-3 border-t border-slate-100 dark:border-slate-800">
                  <p className="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide mb-1">
                    Reasoning
                  </p>
                  <p className="text-sm text-slate-700 dark:text-slate-300">
                    {item.reasoning || "No reasoning recorded."}
                  </p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
};
