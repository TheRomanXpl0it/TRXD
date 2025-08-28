"use client";
import React, { createContext, useContext, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { startInstance, stopInstance } from "@/lib/backend-interaction";
import { Challenge as ChallengeType, ChallengeContext } from "@/context/ChallengeProvider";

type InstanceState = {
  isInstanced: boolean;
  isRunning: boolean;
  remaining: number;         // seconds remaining
  starting: boolean;
  stopping: boolean;
  start: () => Promise<void>;
  stop: () => Promise<void>;
};

const NOOP = async () => {};
const InstanceContext = createContext<InstanceState | undefined>(undefined);

// Accepts several shapes: number, { timeout }, { remaining }, { lifetime }, { expiresAt }
function computeRemainingFromResponse(data: any): number | null {
  if (data == null) return null;
  if (typeof data === "number") return Math.max(0, Math.floor(data));
  if (typeof data.timeout === "number") return Math.max(0, Math.floor(data.timeout));
  if (typeof data.remaining === "number") return Math.max(0, Math.floor(data.remaining));
  if (typeof data.lifetime === "number") return Math.max(0, Math.floor(data.lifetime));
  if (typeof data.expiresAt === "string") {
    const ms = new Date(data.expiresAt).getTime() - Date.now();
    if (!Number.isNaN(ms)) return Math.max(0, Math.floor(ms / 1000));
  }
  return null;
}

function computeHostFromResponse(data: any): string | null {
  let host: string | null = data?.host ?? data?.remote ?? null;
  if (!host) return null;
  if (data?.port) host += `:${data.port}`;
  return host;
}

export function InstanceProvider({
  challenge,
  children,
}: {
  challenge: ChallengeType;
  children: React.ReactNode;
}) {
  const challengeCtx = useContext(ChallengeContext);
  const setChallenges =
    challengeCtx?.setChallenges ?? ((_: React.SetStateAction<ChallengeType[]>) => {});

  const isInstanced = Boolean(challenge.instanced);
  const [starting, setStarting] = useState(false);
  const [stopping, setStopping] = useState(false);

  // initialize from challenge.timeout (number of seconds left)
  const initialRemaining = useMemo(
    () => (typeof challenge.timeout === "number" ? Math.max(0, challenge.timeout) : 0),
    [challenge.timeout]
  );
  const [remaining, setRemaining] = useState<number>(initialRemaining);

  // keep local remaining in sync when parent updates challenge
  useEffect(() => {
    setRemaining(initialRemaining);
  }, [initialRemaining]);

  const isRunning = isInstanced && remaining > 0;

  // tick countdown
  useEffect(() => {
    if (!isRunning) return;
    const id = setInterval(() => setRemaining((r) => Math.max(0, r - 1)), 1000);
    return () => clearInterval(id);
  }, [isRunning]);

  const start = async () => {
    if (starting || !isInstanced) return;
    setStarting(true);
    try {
      // Create promise once, show toasts, then await it for data.
      const p = startInstance(challenge.id);

      toast.promise(p, {
        loading: "Starting instance…",
        success: (r: any) => {
          const secs = computeRemainingFromResponse(r?.data) ?? 0; // reads r.data.timeout correctly
          const host = computeHostFromResponse(r?.data);
          return r?.data?.status === "AlreadyRunning"
            ? `Instance already running (${new Date(secs * 1000).toISOString().slice(11, 19)} left).`
            : `Instance started (${new Date(secs * 1000).toISOString().slice(11, 19)}).` +
              (host ? ` ${host}` : "");
        },
        error: (err: any) => err?.response?.data?.message || "Failed to start the challenge.",
      });

      const res = await p; // <-- actual response here
      const secs = computeRemainingFromResponse(res?.data) ?? 0;  // handles { timeout: 1799 }
      const host = computeHostFromResponse(res?.data);            // handles { host, port? }

      setRemaining(secs);
      setChallenges((prev) =>
        prev.map((c) =>
          c.id === challenge.id ? { ...c, timeout: secs, remote: host ?? c.remote } : c
        )
      );
    } catch (e) {
      console.error(e);
    } finally {
      setStarting(false);
    }
  };

  const stop = async () => {
    if (stopping || !isInstanced) return;
    setStopping(true);
    try {
      await toast.promise(stopInstance(challenge.id), {
        loading: "Stopping instance…",
        success: "Instance stopped.",
        error: (err: any) => err?.response?.data?.message || "Failed to stop instance.",
      });
      setRemaining(0);
      setChallenges((prev) =>
        prev.map((c) => (c.id === challenge.id ? { ...c, timeout: 0, remote: "" } : c))
      );
    } catch (e) {
      console.error(e);
    } finally {
      setStopping(false);
    }
  };

  const value: InstanceState = {
    isInstanced,
    isRunning,
    remaining,
    starting,
    stopping,
    start,
    stop,
  };

  return <InstanceContext.Provider value={value}>{children}</InstanceContext.Provider>;
}

export function useInstance() {
  return useContext(InstanceContext) ?? {
    isInstanced: false,
    isRunning: false,
    remaining: 0,
    starting: false,
    stopping: false,
    start: NOOP,
    stop: NOOP,
  };
}
