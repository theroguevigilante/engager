import { useEffect, useRef, useCallback } from "react";

export interface AggregateTelemetry {
  student_id: string;
  session_duration: number;
  active_duration: number;
  idle_duration: number;
  tab_switches: number;
  interaction_count: number;
}

interface UseActivityTrackerOptions {
  studentId: string;
  attemptId: string | null;
  isActive: boolean;
  apiBase: string;
  flushIntervalMs?: number;
}

export function useActivityTracker({
  studentId,
  attemptId,
  isActive,
  apiBase,
  flushIntervalMs = 15000,
}: UseActivityTrackerOptions) {
  const sessionDurationRef = useRef<number>(0);
  const activeDurationRef = useRef<number>(0);
  const idleDurationRef = useRef<number>(0);
  const tabSwitchesRef = useRef<number>(0);
  const interactionCountRef = useRef<number>(0);

  const lastInteractionTimeRef = useRef<number>(Date.now());
  const isDocumentVisibleRef = useRef<boolean>(!document.hidden);
  const isUserIdleRef = useRef<boolean>(false);

  const flushTelemetry = useCallback(async () => {
    if (!attemptId || !isActive) return;

    const payload: AggregateTelemetry = {
      student_id: studentId,
      session_duration: Math.round(sessionDurationRef.current),
      active_duration: Math.round(activeDurationRef.current),
      idle_duration: Math.round(idleDurationRef.current),
      tab_switches: tabSwitchesRef.current,
      interaction_count: interactionCountRef.current,
    };

    try {
      await fetch(`${apiBase}/attempts/${attemptId}/signals`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
    } catch (err) {
      console.warn("Failed to flush telemetry:", err);
    }
  }, [studentId, attemptId, isActive, apiBase]);

  const lastCountedInteractionRef = useRef<number>(0);
  const recordInteraction = useCallback(() => {
    if (!isActive) return;
    const now = Date.now();
    lastInteractionTimeRef.current = now;

    if (isUserIdleRef.current) {
      isUserIdleRef.current = false;
    }

    if (now - lastCountedInteractionRef.current > 300) {
      interactionCountRef.current += 1;
      lastCountedInteractionRef.current = now;
    }
  }, [isActive]);

  useEffect(() => {
    if (!isActive || !attemptId) {
      sessionDurationRef.current = 0;
      activeDurationRef.current = 0;
      idleDurationRef.current = 0;
      tabSwitchesRef.current = 0;
      interactionCountRef.current = 0;
      return;
    }

    lastInteractionTimeRef.current = Date.now();
    isDocumentVisibleRef.current = !document.hidden;

    const tickInterval = setInterval(() => {
      const now = Date.now();
      const elapsedSinceInteraction = (now - lastInteractionTimeRef.current) / 1000;
      sessionDurationRef.current += 1;

      if (!isDocumentVisibleRef.current || elapsedSinceInteraction >= 10) {
        isUserIdleRef.current = true;
        idleDurationRef.current += 1;
      } else {
        activeDurationRef.current += 1;
      }
    }, 1000);

    const flushTimer = setInterval(() => {
      flushTelemetry();
    }, flushIntervalMs);

    const handleVisibilityChange = () => {
      if (document.hidden) {
        isDocumentVisibleRef.current = false;
        tabSwitchesRef.current += 1;
      } else {
        isDocumentVisibleRef.current = true;
        lastInteractionTimeRef.current = Date.now();
      }
    };

    const handleWindowBlur = () => {
      tabSwitchesRef.current += 1;
    };

    window.addEventListener("visibilitychange", handleVisibilityChange);
    window.addEventListener("blur", handleWindowBlur);
    window.addEventListener("keydown", recordInteraction, { passive: true });
    window.addEventListener("click", recordInteraction, { passive: true });
    window.addEventListener("scroll", recordInteraction, { passive: true });

    return () => {
      clearInterval(tickInterval);
      clearInterval(flushTimer);
      window.removeEventListener("visibilitychange", handleVisibilityChange);
      window.removeEventListener("blur", handleWindowBlur);
      window.removeEventListener("keydown", recordInteraction);
      window.removeEventListener("click", recordInteraction);
      window.removeEventListener("scroll", recordInteraction);
    };
  }, [isActive, attemptId, flushIntervalMs, flushTelemetry, recordInteraction]);

  return {
    recordInteraction,
    flushTelemetry,
  };
}
