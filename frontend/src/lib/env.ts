
export const config = {
	/**
	 * Backend URL for API requests and file downloads
	 */
	get backendUrl(): string {
		return import.meta.env.VITE_BACKEND_URL || 'http://localhost:1337';
	},

	/**
	 * Get the full URL for a backend endpoint
	 * @param path - The path to append to the backend URL
	 * @returns Full URL to the backend endpoint
	 */
	getBackendUrl(path: string = ''): string {
		const baseUrl = this.backendUrl;
		if (!path) return baseUrl;

		// Ensure path starts with /
		const normalizedPath = path.startsWith('/') ? path : `/${path}`;
		return `${baseUrl}${normalizedPath}`;
	}
};
