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
  count: number;
  color: string; // any valid CSS color
  label?: string;
};

export type UserCategoryRingProps = {
  size?: number;
  strokeWidth?: number;
  totalSolves: number;
  categories: CategorySlice[];
  showLegend?: boolean;
  className?: string;
};

export function UserCategoryRing({
  size = 200,
  strokeWidth = 14,
  totalSolves,
  categories,
  showLegend = true,
  className = "",
}: UserCategoryRingProps) {
  const totalCount = Math.max(
    0,
    categories?.reduce((sum, c) => sum + Math.max(0, c.count), 0) ?? 0
  );

  const radius = (size - strokeWidth) / 2;
  const circumference = 2 * Math.PI * radius;

  // Create normalized slices [0..1]
  const slices = (categories || [])
    .filter((c) => c && c.count > 0)
    .map((c) => ({ ...c, frac: totalCount > 0 ? c.count / totalCount : 0 }));

  // For segment placement along the ring
  let acc = 0; // accumulated fraction
  const segments = slices.map((s) => {
    const segLen = circumference * s.frac;
    const offset = circumference * (1 - acc); // dashoffset from end
    acc += s.frac;
    return { ...s, segLen, offset };
  });

  const hasData = totalCount > 0 && segments.length > 0;

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
            {/* Background ring */}
            <circle
              cx={size / 2}
              cy={size / 2}
              r={radius}
              fill="none"
              stroke={hasData ? "#e5e7eb" : "#e5e7eb"}
              strokeWidth={strokeWidth}
            />

            {hasData &&
              segments.map((seg) => (
                <circle
                  key={seg.key}
                  cx={size / 2}
                  cy={size / 2}
                  r={radius}
                  fill="none"
                  stroke={seg.color}
                  strokeWidth={strokeWidth}
                  strokeLinecap="butt"
                  strokeDasharray={`${seg.segLen} ${circumference - seg.segLen}`}
                  strokeDashoffset={seg.offset}
                >
                  <title>
                    {`${seg.label ?? seg.key}: ${Math.round(seg.frac * 100)}% (${seg.count})`}
                  </title>
                </circle>
              ))}
          </g>
        </svg>

        {/* Center label */}
        <div className="absolute inset-0 flex flex-col items-center justify-center select-none">
          <div className="text-3xl font-bold leading-none">{totalSolves}</div>
          <div className="text-xs text-gray-500">solves</div>
        </div>
      </div>

      {showLegend && (
        <div className="mt-4 grid grid-cols-2 gap-2 text-sm w-full max-w-sm">
          {categories.map((c) => {
            const pct = totalCount > 0 ? Math.round((c.count / totalCount) * 100) : 0;
            return (
              <div key={c.key} className="flex items-center gap-2">
                <span
                  className="inline-block w-3 h-3 rounded"
                  style={{ background: c.color }}
                  aria-hidden
                />
                <span className="truncate">{c.label ?? c.key}</span>
                <span className="ml-auto tabular-nums text-gray-500">{pct}%</span>
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
  const data: CategorySlice[] = [
    { key: "web", label: "Web", count: 18, color: "#10b981" },
    { key: "pwn", label: "Pwn", count: 9, color: "#3b82f6" },
    { key: "rev", label: "Reversing", count: 6, color: "#a855f7" },
    { key: "crypto", label: "Crypto", count: 3, color: "#f59e0b" },
  ];
  const totalSolves = data.reduce((s, d) => s + d.count, 0);
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
