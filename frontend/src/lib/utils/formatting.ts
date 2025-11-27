export function formatDate(iso?: string): string {
	if (!iso) return '-';
	const date = new Date(iso);
	return Number.isNaN(+date) ? '-' : date.toLocaleString();
}

export function formatTimeSince(iso?: string): string {
	if (!iso) return '-';
	const seconds = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / 1000));
	const hours = Math.floor(seconds / 3600);
	const minutes = Math.floor((seconds % 3600) / 60);
	const secs = seconds % 60;

	if (hours > 0) return `${hours}h ${minutes}m`;
	if (minutes > 0) return `${minutes}m ${secs}s`;
	return `${secs}s`;
}

export function formatNumber(value: any): number {
	return Number(value ?? 0);
}

export function truncateText(text: string, maxLength: number = 32): string {
	if (!text || text.length <= maxLength) return text;
	return text.slice(0, maxLength) + '...';
}
