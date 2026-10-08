import React, { useState, useEffect } from "react";
import { Activity } from "@/types";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

interface ActivityModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSaved: () => void;
  activityToEdit?: Activity | null;
  apiBase: string;
}

const inputCls =
  "w-full rounded-md border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 px-3 py-2 text-sm text-slate-900 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-500";

export const ActivityModal: React.FC<ActivityModalProps> = ({
  isOpen,
  onClose,
  onSaved,
  activityToEdit,
  apiBase,
}) => {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [content, setContent] = useState("");
  const [correctAnswer, setCorrectAnswer] = useState("");
  const [topic, setTopic] = useState("Mathematics");
  const [difficulty, setDifficulty] = useState("intermediate");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (activityToEdit) {
      setTitle(activityToEdit.title || "");
      setDescription(activityToEdit.description || "");
      setContent(activityToEdit.content || "");
      setCorrectAnswer(activityToEdit.correct_answer || "");
      setTopic(activityToEdit.topic || "Mathematics");
      setDifficulty(activityToEdit.difficulty || "intermediate");
    } else {
      setTitle("");
      setDescription("");
      setContent("");
      setCorrectAnswer("");
      setTopic("Mathematics");
      setDifficulty("intermediate");
    }
    setError(null);
  }, [activityToEdit, isOpen]);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) {
      setError("Title and content are required.");
      return;
    }

    setIsSubmitting(true);
    setError(null);

    const isEdit = Boolean(activityToEdit?.id);
    const url = isEdit
      ? `${apiBase}/activities/${activityToEdit!.id}`
      : `${apiBase}/activities`;
    const method = isEdit ? "PUT" : "POST";

    const payload: Partial<Activity> = {
      title: title.trim(),
      description: description.trim(),
      content: content.trim(),
      correct_answer: correctAnswer.trim(),
      topic,
      difficulty,
    };

    if (!isEdit) {
      payload.id = `act-${Date.now().toString().slice(-4)}`;
    }

    try {
      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        throw new Error(`Failed to save activity: HTTP ${res.status}`);
      }

      onSaved();
      onClose();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Error saving activity");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-slate-900/40 dark:bg-slate-950/60 backdrop-blur-xs flex items-center justify-center p-4">
      <Card className="max-w-xl w-full bg-white dark:bg-slate-900 shadow-xl max-h-[90vh] flex flex-col border border-slate-200 dark:border-slate-700">
        <CardHeader className="pb-3 border-b border-slate-100 dark:border-slate-800">
          <div className="flex justify-between items-center">
            <CardTitle className="text-lg">
              {activityToEdit ? "Edit Activity" : "Create New Activity"}
            </CardTitle>
            <button
              onClick={onClose}
              type="button"
              className="text-slate-400 dark:text-slate-500 hover:text-slate-600 dark:hover:text-slate-300 text-xl leading-none"
            >
              ×
            </button>
          </div>
        </CardHeader>
        <CardContent className="pt-4 overflow-y-auto">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-xs font-medium text-slate-600 dark:text-slate-400 mb-1">
                Title *
              </label>
              <input
                className={inputCls}
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="e.g. Fractions Basics"
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 dark:text-slate-400 mb-1">
                Description
              </label>
              <input
                className={inputCls}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Short summary"
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-xs font-medium text-slate-600 dark:text-slate-400 mb-1">
                  Topic
                </label>
                <select
                  className={inputCls}
                  value={topic}
                  onChange={(e) => setTopic(e.target.value)}
                >
                  <option>Mathematics</option>
                  <option>Science</option>
                  <option>History</option>
                  <option>Language Arts</option>
                  <option>Computer Science</option>
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium text-slate-600 dark:text-slate-400 mb-1">
                  Difficulty
                </label>
                <select
                  className={inputCls}
                  value={difficulty}
                  onChange={(e) => setDifficulty(e.target.value)}
                >
                  <option value="beginner">Beginner</option>
                  <option value="intermediate">Intermediate</option>
                  <option value="advanced">Advanced</option>
                </select>
              </div>
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 dark:text-slate-400 mb-1">
                Content *
              </label>
              <textarea
                className={`${inputCls} min-h-[100px]`}
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="Activity instructions or problem set..."
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 dark:text-slate-400 mb-1">
                Correct Answer
              </label>
              <input
                className={inputCls}
                value={correctAnswer}
                onChange={(e) => setCorrectAnswer(e.target.value)}
                placeholder="Expected answer"
              />
            </div>
            {error && <p className="text-sm text-rose-600 dark:text-rose-400">{error}</p>}
            <div className="flex gap-2 pt-2">
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Saving…" : activityToEdit ? "Update" : "Create"}
              </Button>
              <Button
                type="button"
                onClick={onClose}
                className="bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-600 shadow-none hover:bg-slate-50 dark:hover:bg-slate-700"
              >
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};
