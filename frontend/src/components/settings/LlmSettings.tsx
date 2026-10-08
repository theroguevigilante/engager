import React, { useCallback, useEffect, useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface LlmStatus {
  configured: boolean;
  provider: string;
  model: string;
  base_url: string;
  last_source: "ai" | "rules";
  last_error?: string;
}

interface Props {
  apiBase: string;
  refreshKey?: unknown;
}

const inputCls =
  "w-full rounded-md border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 px-3 py-2 text-sm text-slate-900 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-500";

export const LlmSettings: React.FC<Props> = ({ apiBase, refreshKey }) => {
  const [status, setStatus] = useState<LlmStatus | null>(null);
  const [provider, setProvider] = useState("anthropic");
  const [apiKey, setApiKey] = useState("");
  const [model, setModel] = useState("");
  const [baseUrl, setBaseUrl] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    try {
      const res = await fetch(`${apiBase}/settings/llm`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data: LlmStatus = await res.json();
      setStatus(data);
      setProvider(data.provider);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Could not load AI settings");
    }
  }, [apiBase]);

  useEffect(() => {
    load();
  }, [load, refreshKey]);

  const save = async () => {
    setBusy(true);
    setError(null);
    try {
      const res = await fetch(`${apiBase}/settings/llm`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          provider,
          api_key: apiKey,
          model,
          base_url: provider === "openai" || provider === "gemini" ? baseUrl : "",
        }),
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      setStatus(await res.json());
      setApiKey("");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Save failed");
    } finally {
      setBusy(false);
    }
  };

  const disable = async () => {
    setBusy(true);
    try {
      const res = await fetch(`${apiBase}/settings/llm`, { method: "DELETE" });
      if (res.ok) setStatus(await res.json());
    } finally {
      setBusy(false);
    }
  };

  return (
    <Card className="border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="text-base">AI content adaptation</CardTitle>
          <span
            className={`text-[10px] uppercase font-bold px-2 py-0.5 rounded border ${
              status?.configured
                ? "bg-emerald-50 dark:bg-emerald-950/50 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-800"
                : "bg-slate-50 dark:bg-slate-800 text-slate-500 dark:text-slate-400 border-slate-200 dark:border-slate-700"
            }`}
          >
            {status?.configured ? `On · ${status.provider} · ${status.model}` : "Rules-only mode"}
          </span>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="text-sm text-slate-600 dark:text-slate-400">
          Optional. An LLM only rewrites activity content. Engagement states and interventions always come from the
          explainable rule engine, and any AI failure falls back to rules. The key is kept in server memory only.
          You can also set <code className="text-xs bg-slate-100 dark:bg-slate-800 px-1 rounded">ANTHROPIC_API_KEY</code>, <code className="text-xs bg-slate-100 dark:bg-slate-800 px-1 rounded">OPENAI_API_KEY</code> or <code className="text-xs bg-slate-100 dark:bg-slate-800 px-1 rounded">GEMINI_API_KEY</code> before starting the backend.
        </p>
        <div className="grid gap-3 md:grid-cols-2">
          <label className="text-xs font-medium text-slate-600 dark:text-slate-400 space-y-1">
            Provider
            <select className={inputCls} value={provider} onChange={(e) => setProvider(e.target.value)}>
              <option value="anthropic">Anthropic (Claude)</option>
              <option value="openai">OpenAI-compatible</option>
              <option value="gemini">Google Gemini</option>
            </select>
          </label>
          <label className="text-xs font-medium text-slate-600 dark:text-slate-400 space-y-1">
            API key
            <input
              type="password"
              className={inputCls}
              value={apiKey}
              placeholder={status?.configured ? "•••• saved (leave blank to keep)" : "Paste key"}
              onChange={(e) => setApiKey(e.target.value)}
              autoComplete="off"
            />
          </label>
          <label className="text-xs font-medium text-slate-600 dark:text-slate-400 space-y-1">
            Model (optional)
            <input className={inputCls} value={model} placeholder="default" onChange={(e) => setModel(e.target.value)} />
          </label>
          {provider === "openai" && (
            <label className="text-xs font-medium text-slate-600 dark:text-slate-400 space-y-1">
              Base URL (optional)
              <input
                className={inputCls}
                value={baseUrl}
                placeholder="https://api.openai.com/v1"
                onChange={(e) => setBaseUrl(e.target.value)}
              />
            </label>
          )}
          {provider === "gemini" && (
            <label className="text-xs font-medium text-slate-600 dark:text-slate-400 space-y-1">
              Base URL (optional)
              <input
                className={inputCls}
                value={baseUrl}
                placeholder="https://generativelanguage.googleapis.com"
                onChange={(e) => setBaseUrl(e.target.value)}
              />
            </label>
          )}
        </div>
        <div className="flex gap-2">
          <Button onClick={save} disabled={busy || (!apiKey && !status?.configured)}>
            Save
          </Button>
          <Button
            onClick={disable}
            disabled={busy || !status?.configured}
            className="bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-600 shadow-none hover:bg-slate-50 dark:hover:bg-slate-700"
          >
            Disable
          </Button>
        </div>
        {status?.last_error && (
          <p className="text-xs text-amber-700 dark:text-amber-400">Last AI call failed, so rules were used: {status.last_error}</p>
        )}
        {error && <p className="text-xs text-rose-600 dark:text-rose-400">{error}</p>}
      </CardContent>
    </Card>
  );
};
