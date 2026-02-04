export function formatDate(iso?: string): string {
	if (!iso) return '-';
	const date = new Date(iso);
	return Number.isNaN(+date) ? '-' : date.toLocaleString();
}

export function formatTimeSince(iso?: string): string {
	if (!iso) return '-';
	const seconds = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / 1000));

	const minutes = Math.floor(seconds / 60);
	const hours = Math.floor(minutes / 60);
	const days = Math.floor(hours / 24);
	const months = Math.floor(days / 30);

	if (months > 0) {
		const remDays = days % 30;
		return `${months}m ${remDays}d`;
	}
	if (days > 0) {
		const remHours = hours % 24;
		return `${days}d ${remHours}hrs`;
	}
	if (hours > 0) {
		const remMinutes = minutes % 60;
		return `${hours}hrs ${remMinutes}min`;
	}
	if (minutes > 0) {
		const remSeconds = seconds % 60;
		return `${minutes}min ${remSeconds}s`;
	}
	return `${seconds}second${seconds !== 1 ? 's' : ''}`;
}

export function formatNumber(value: any): number {
	return Number(value ?? 0);
}

export function truncateText(text: string, maxLength: number = 32): string {
	if (!text || text.length <= maxLength) return text;
	return text.slice(0, maxLength) + '...';
}
