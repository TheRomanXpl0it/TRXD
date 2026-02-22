/**
 * Format a countdown in seconds to a human-readable string.
 *   3661 → "1:01:01"
 *    125 → "2:05"
 *     45 → "45"
 */
export function fmtTimeLeft(total: number): string {
    if (!total || total < 0) total = 0;
    const h = Math.floor(total / 3600);
    const m = Math.floor((total % 3600) / 60);
    const s = Math.floor(total % 60);
    if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
    if (m > 0) return `${m}:${String(s).padStart(2, '0')}`;
    return `${s}`;
}
