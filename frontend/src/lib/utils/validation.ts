export function isValidUrl(url: string): boolean {
	if (!url) return false;
	try {
		const parsed = new URL(url);
		return parsed.protocol === 'http:' || parsed.protocol === 'https:';
	} catch {
		return false;
	}
}

export function validateRequired(value: string, fieldName: string = 'Field'): string | null {
	if (!value.trim()) {
		return `${fieldName} is required.`;
	}
	return null;
}

export function validateUrl(value: string, fieldName: string = 'URL'): string | null {
	if (value && !isValidUrl(value)) {
		return `${fieldName} must be a valid URL.`;
	}
	return null;
}

export function validateMinLength(value: string, minLength: number, fieldName: string = 'Field'): string | null {
	if (value.trim().length < minLength) {
		return `${fieldName} must be at least ${minLength} characters.`;
	}
	return null;
}
