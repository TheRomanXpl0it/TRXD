// useCountdown.ts
import { useEffect, useState } from "react";

/**
 * Count down to an absolute deadline in ms (e.g., Date.now() + 30_000).
 * Returns whole seconds remaining (never negative).
 */
export function useCountdown(deadlineMs: number | null, tickMs = 1000): number {
    const [now, setNow] = useState(() => Date.now());
    useEffect(() => {
        if (!deadlineMs) return;
    setNow(Date.now()); // resync immediately when deadline changes
    const id = setInterval(() => setNow(Date.now()), tickMs);
    return () => clearInterval(id);
  }, [deadlineMs, tickMs]);
  if (!deadlineMs) return 0;
  return Math.max(0, Math.floor((deadlineMs - now) / 1000));
}