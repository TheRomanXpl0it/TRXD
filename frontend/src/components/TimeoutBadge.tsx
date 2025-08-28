import { Container } from "lucide-react";
import { useEffect, useState } from "react";
import { useCountdown } from "@/hooks/use-countdown";

function formatHHMMSS(total: number) {
  let t = Math.max(0, Math.floor(total)); // guard negatives/decimals
  const h = Math.floor(t / 3600);
  const m = Math.floor((t % 3600) / 60);
  const s = t % 60;
  const pad = (n: number) => String(n).padStart(2, "0");

  if (h > 0) return `${h}:${pad(m)}:${pad(s)}`; // e.g., 1:05:09
  if (m > 0) return `${m}:${pad(s)}`;          // e.g., 20:20
  return `${s}s`;                               // e.g., 45s
}

function TimeoutBadge({ seconds }: { seconds?: number }) {
    // Use STATE, not ref, so we re-render when setting it
    const [deadlineMs, setDeadlineMs] = useState<number | null>(null);

  useEffect(() => {
      if (typeof seconds === "number" && seconds > 0) {
      setDeadlineMs(Date.now() + seconds * 1000);
    } else {
        setDeadlineMs(null);
    }
  }, [seconds]);

  const remaining = useCountdown(deadlineMs);
  
  if (!deadlineMs) return null;
  return (
      <span className="text-sm flex items-center gap-1">
      <Container size={18} className="text-blue-500" />
      {formatHHMMSS(remaining)}
    </span>
  );
}

export { TimeoutBadge, formatHHMMSS };