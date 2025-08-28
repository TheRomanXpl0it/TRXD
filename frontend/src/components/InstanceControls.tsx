"use client";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useEffect, useState } from "react";
import { useInstance } from "@/context/InstanceProvider";
import { Container, OctagonX } from "lucide-react";
import { formatHHMMSS } from "@/components/TimeoutBadge"



export function InstanceControls({
  remote,
}: {
  remote?: string;
  className?: string;
}) {
  const { isInstanced, isRunning, remaining, starting, stopping, start, stop } = useInstance();
  const [dots, setDots] = useState("");

  // animate ". .. ..." while starting
  useEffect(() => {
    if (!starting) { setDots(""); return; }
    const seq = [".", "..", "..."]; let i = 0;
    const id = setInterval(() => { setDots(seq[i]); i = (i + 1) % seq.length; }, 400);
    return () => clearInterval(id);
  }, [starting]);

  if (!isInstanced) return null;

  if (!isRunning) {
    // Show Start button (same look as before)
    return (
      <div className={"w-full"}>
        <Button
          className={cn("bg-blue-500 w-full text-white", starting && "cursor-wait", !starting && "cursor-pointer")}
          aria-busy={starting ? "true" : "false"}
          aria-disabled={starting ? "true" : "false"}
          onClick={() => { if (!starting) start(); }}
        >
          <Container className="text-white" size={24} />
          <span className="ml-2">{starting ? `Starting${dots}` : "Start Instance"}</span>
        </Button>
      </div>
    );
  }

    const href =
    remote
        ? (/^[a-zA-Z][a-zA-Z\d+\-.]*:\/\//.test(remote) ? remote : `//${remote}`)
        : undefined;

  // Running: show blue bar with timer + Stop
  return (
    <div className={"flex items-center justify-between gap-3 w-100"}>
        <div className="rounded-lg bg-blue-500 text-white p-2">
            <div className="flex items-center justify-between gap-3">
                <div className="flex items-center gap-2">
                    <Container size={24} className="text-blue-200" />
                    {remote ? (
                        <a
                            href={href}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-white hover:text-blue-200 transition-colors"
                        >
                            {remote}
                        </a>
                    ) : "Instance Running"}
                </div>

                <span className="font-mono tabular-nums px-2">
                    {formatHHMMSS(remaining)}
                </span>
            </div>
        </div>
        <Button
        onClick={() => { if (!stopping) stop(); }}
        className={cn("bg-red-500 hover:cursor-pointer", stopping && "cursor-wait")}
        aria-busy={stopping ? "true" : "false"}
        aria-disabled={stopping ? "true" : "false"}
        >
            <span className="hidden sm:inline">{stopping ? "Stopping…" : <OctagonX className="text-white" size={20} />}</span>
        </Button>
    </div>
  );
}
