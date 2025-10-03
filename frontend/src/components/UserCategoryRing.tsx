import { useState } from "react";
/**
 * UserCategoryRing
 * A circular border (donut-ring) segmented by category proportions.
 * Center shows total solves.
 *
 * Props:
 * - size: diameter in px
 * - strokeWidth: ring thickness in px
 * - totalSolves: number in the middle
 * - categories: array of { key, label?, count, color }
 * - showLegend: render a legend under the ring
 * - className: extra wrapper classes
 */
export type CategorySlice = {
  key: string;
  color: string; // any valid CSS color
  label?: string;
  // total challenges in this category
  total?: number;
  // solved challenges (defaults to 0). If only `count` is provided, it is treated as solved for backward compatibility.
  solved?: number;
  // legacy: when provided, treated as solved count
  count?: number;
};

export type UserCategoryRingProps = {
  size?: number;
  strokeWidth?: number;
  totalSolves: number;
  teamSolves?: unknown;
  categories: CategorySlice[];
  showLegend?: boolean;
  className?: string;
};

export function UserCategoryRing({
  size = 200,
  strokeWidth = 14,
  totalSolves,
  teamSolves: _teamSolves = undefined,
  categories,
  showLegend = true,
  className = "",
}: UserCategoryRingProps) {
  // Normalize categories to have total and solved numbers
  const normalized = (categories || [])
    .map((c) => {
      const total = Math.max(0, c.total ?? c.count ?? 0);
      const solved = Math.min(Math.max(0, c.solved ?? c.count ?? 0), total);
      return { ...c, total, solved };
    })
    .filter((c) => c.total > 0);

  const totalAll = normalized.reduce((sum, c) => sum + c.total, 0);

  const radius = (size - strokeWidth) / 2;
  const circumference = 2 * Math.PI * radius;
  const innerRadius = Math.max(0, radius - strokeWidth / 2);

  // Build solved segments, anchored at the start of each category slice.
  // Each category gets an angular share equal to total/totalAll, and within that share we draw a colored arc
  // proportional to solved/total. The remaining area stays as the gray background ring (unsolved).
  let acc = 0; // accumulated category fraction
  const _solvedSegments = normalized.map((c) => {
    const catFrac = totalAll > 0 ? c.total / totalAll : 0;
    const solvedFracWithinCat = c.total > 0 ? c.solved / c.total : 0;

    const segLen = circumference * (catFrac * solvedFracWithinCat);
    const offset = circumference * (1 - acc); // start at the beginning of the category slice

    const startAcc = acc;
    acc += catFrac;

    return {
      key: c.key,
      label: c.label,
      color: c.color,
      segLen,
      offset,
      catFrac,
      solvedFracWithinCat,
      startAcc,
      total: c.total,
      solved: c.solved,
    };
  });

  void _solvedSegments;
  let acc2 = 0;
  const categorySegments = normalized.map((c) => {
    const catFrac = totalAll > 0 ? c.total / totalAll : 0;
    const catLen = circumference * catFrac;
    const offset = circumference * (1 - acc2);
    acc2 += catFrac;
    return {
      key: c.key,
      label: c.label,
      color: c.color,
      catLen,
      offset,
      total: c.total,
      solved: c.solved,
      pct: c.total > 0 ? Math.round((c.solved / c.total) * 100) : 0,
    };
  });

  const hasData = totalAll > 0;
  void _teamSolves;
  const [hover, setHover] = useState<null | {
    x: number;
    y: number;
    label: string;
    color: string;
    solved: number;
    unsolved: number;
    pct: number;
    summary?: boolean;
  }>(null);

  return (
    <div className={`inline-flex flex-col items-center ${className}`}>
      <div
        className="relative"
        style={{ width: size, height: size }}
        aria-label="User category distribution"
        role="img"
      >
        <svg
          width={size}
          height={size}
          viewBox={`0 0 ${size} ${size}`}
          className="block"
        >
          <g transform={`rotate(-90 ${size / 2} ${size / 2})`}>
            {/* Unsolved per-category arcs */}
            {hasData &&
              categorySegments.map((seg) => {
                const solvedLen =
                  seg.catLen * (seg.total > 0 ? seg.solved / seg.total : 0);
                const unsolvedLen = Math.max(0, seg.catLen - solvedLen);
                if (unsolvedLen <= 0) return null;
                const unsolvedOffset =
                  (seg.offset - solvedLen + circumference) % circumference;
                return (
                  <circle
                    key={`unsolved-${seg.key}`}
                    cx={size / 2}
                    cy={size / 2}
                    r={radius}
                    fill="none"
                    stroke="#e5e7eb"
                    strokeWidth={strokeWidth}
                    strokeLinecap="butt"
                    strokeDasharray={`${unsolvedLen} ${circumference - unsolvedLen}`}
                    strokeDashoffset={unsolvedOffset}
                  />
                );
              })}

            {hasData &&
              categorySegments
                .map((seg) => {
                  const solvedLen =
                    seg.catLen * (seg.total > 0 ? seg.solved / seg.total : 0);
                  return { ...seg, solvedLen };
                })
                .filter((seg) => seg.solvedLen > 0)
                .map((seg) => (
                  <circle
                    key={`solved-${seg.key}`}
                    cx={size / 2}
                    cy={size / 2}
                    r={radius}
                    fill="none"
                    stroke={seg.color}
                    strokeWidth={strokeWidth}
                    strokeLinecap="butt"
                    strokeDasharray={`${seg.solvedLen} ${circumference - seg.solvedLen}`}
                    strokeDashoffset={seg.offset}
                  />
                ))}

            {hasData &&
              categorySegments.map((seg) => (
                <circle
                  key={`hit-${seg.key}`}
                  cx={size / 2}
                  cy={size / 2}
                  r={radius}
                  fill="none"
                  stroke={seg.color}
                  strokeOpacity={0.08}
                  style={{ pointerEvents: "stroke" }}
                  strokeWidth={strokeWidth}
                  strokeLinecap="butt"
                  strokeDasharray={`${seg.catLen} ${circumference - seg.catLen}`}
                  strokeDashoffset={seg.offset}
                  onMouseEnter={(e) => {
                    const rect = (
                      e.currentTarget as SVGCircleElement
                    ).ownerSVGElement?.getBoundingClientRect();
                    setHover({
                      x: e.clientX - (rect?.left ?? 0),
                      y: e.clientY - (rect?.top ?? 0),
                      label: seg.label ?? seg.key,
                      color: seg.color,
                      solved: seg.solved,
                      unsolved: seg.total - seg.solved,
                      pct: seg.pct,
                      summary: false,
                    });
                  }}
                  onMouseMove={(e) => {
                    const rect = (
                      e.currentTarget as SVGCircleElement
                    ).ownerSVGElement?.getBoundingClientRect();
                    setHover((prev) => ({
                      ...(prev ?? {
                        label: seg.label ?? seg.key,
                        color: seg.color,
                        solved: seg.solved,
                        unsolved: seg.total - seg.solved,
                        pct: seg.pct,
                        summary: false,
                      }),
                      x: e.clientX - (rect?.left ?? 0),
                      y: e.clientY - (rect?.top ?? 0),
                    }));
                  }}
                  onMouseLeave={() => setHover(null)}
                />
              ))}
            {/* Center hover target */}
            {hasData && (
              <circle
                cx={size / 2}
                cy={size / 2}
                r={innerRadius}
                fill="rgba(0,0,0,0.001)"
                style={{ pointerEvents: "all" }}
                onMouseEnter={(e) => {
                  const rect = (
                    e.currentTarget as SVGCircleElement
                  ).ownerSVGElement?.getBoundingClientRect();
                  setHover({
                    x: e.clientX - (rect?.left ?? 0),
                    y: e.clientY - (rect?.top ?? 0),
                    label: "All categories",
                    color: "#999",
                    solved: 0,
                    unsolved: 0,
                    pct: 0,
                    summary: true,
                  });
                }}
                onMouseMove={(e) => {
                  const rect = (
                    e.currentTarget as SVGCircleElement
                  ).ownerSVGElement?.getBoundingClientRect();
                  setHover((prev) => ({
                    ...(prev ?? {
                      label: "All categories",
                      color: "#999",
                      solved: 0,
                      unsolved: 0,
                      pct: 0,
                      summary: true,
                    }),
                    x: e.clientX - (rect?.left ?? 0),
                    y: e.clientY - (rect?.top ?? 0),
                  }));
                }}
                onMouseLeave={() => setHover(null)}
              />
            )}
          </g>
        </svg>

        {/* Center label */}
        <div className="absolute inset-0 flex flex-col items-center justify-center select-none pointer-events-none">
          <div className="text-3xl font-bold leading-none">{totalSolves}</div>
          <div className="text-xs text-gray-500">solves</div>
        </div>
        {hover && (
          <div
            className="absolute z-10 pointer-events-none px-2 py-1 rounded bg-black/80 text-white text-xs shadow-md"
            style={{
              left: Math.min(Math.max(hover.x + 8, 0), size - 160),
              top: Math.min(Math.max(hover.y + 8, 0), size - 60),
              maxWidth: 160,
            }}
          >
            {hover.summary ? (
              <div>
                <div className="flex flex-col gap-1">
                  {categorySegments.map((seg) => {
                    const total = seg.total;
                    const solved = seg.solved;
                    const pct = Math.round((solved / (total || 1)) * 100);
                    return (
                      <div key={seg.key} className="flex items-center gap-2">
                        <span
                          className="inline-block w-3 h-3 rounded"
                          style={{ background: seg.color }}
                        />
                        <span className="truncate">{seg.label ?? seg.key}</span>
                        <span className="ml-auto tabular-nums">
                          {solved}/{total} {pct}%
                        </span>
                      </div>
                    );
                  })}
                </div>
              </div>
            ) : (
              <>
                <div className="flex items-center gap-2">
                  <span
                    className="inline-block w-3 h-3 rounded"
                    style={{ background: hover.color }}
                  />
                  <span className="font-medium">{hover.label}</span>
                </div>
                <div className="mt-1 tabular-nums">
                  {hover.solved}/{hover.solved + hover.unsolved}{"  "}
                  {Math.round(
                    (hover.solved / (hover.solved + hover.unsolved || 1)) * 100,
                  )}
                  %
                </div>
              </>
            )}
          </div>
        )}
      </div>

      {showLegend && (
        <div className="mt-4 grid grid-cols-2 gap-2 text-sm w-full max-w-sm">
          {(categories || []).map((c) => {
            const total = Math.max(0, c.total ?? c.count ?? 0);
            const solved = Math.min(
              Math.max(0, c.solved ?? c.count ?? 0),
              total,
            );
            const pct = total > 0 ? Math.round((solved / total) * 100) : 0;
            return (
              <div key={c.key} className="flex items-center gap-2">
                <span
                  className="inline-block w-3 h-3 rounded"
                  style={{ background: c.color }}
                  aria-hidden
                />
                <span className="truncate">{c.label ?? c.key}</span>
                <span className="ml-auto tabular-nums text-gray-500">
                  {pct}%
                </span>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

// ---------------- Demo ----------------
export function UserCategoryRingDemo() {
  // Example: totals per category and how many are solved
  const data: CategorySlice[] = [
    { key: "web", label: "Web", total: 100, solved: 18, color: "#10b981" },
    { key: "pwn", label: "Pwn", total: 100, solved: 9, color: "#3b82f6" },
    { key: "rev", label: "Reversing", total: 100, solved: 6, color: "#a855f7" },
    { key: "crypto", label: "Crypto", total: 100, solved: 3, color: "#f59e0b" },
  ];
  const totalSolves = data.reduce((s, d) => s + (d.solved ?? 0), 0);
  return (
    <div className="p-6">
      <UserCategoryRing
        totalSolves={totalSolves}
        categories={data}
        size={220}
        strokeWidth={16}
      />
    </div>
  );
}
