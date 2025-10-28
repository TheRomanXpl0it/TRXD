import { config } from '$lib/env';

const BASE = '/api';

function urlOf(path: string) {
	// If it's already a full URL, return as-is
	if (/^https?:\/\//.test(path)) return path;

	// For relative paths, use the BASE (which gets proxied in dev)
	return `${BASE}${path.startsWith('/') ? '' : '/'}${path}`;
}

async function parse<T>(res: Response): Promise<T> {
	if (res.status === 204) return undefined as unknown as T;
	const ct = res.headers.get('content-type') || '';
	if (ct.includes('application/json')) {
		try {
			return (await res.json()) as T;
		} catch (e) {
			throw new Error(`Failed to parse JSON response: ${(e as Error).message}`);
		}
	}
	return (await res.text()) as unknown as T;
}

// --- Safe cookie helpers (no regex) ---
function getCookie(name: string): string | null {
	if (typeof document === 'undefined') return null;
	const prefix = `${encodeURIComponent(name)}=`;
	const parts = document.cookie ? document.cookie.split('; ') : [];
	for (const part of parts) {
		if (part.startsWith(prefix)) {
			return decodeURIComponent(part.slice(prefix.length));
		}
	}
	return null;
}

function pickCsrfToken(): string | null {
	return getCookie('csrf_') || null;
}

// Type guards for BodyInit detection
const isFormData = (v: unknown): v is FormData =>
	typeof FormData !== 'undefined' && v instanceof FormData;

const isURLSearchParams = (v: unknown): v is URLSearchParams =>
	typeof URLSearchParams !== 'undefined' && v instanceof URLSearchParams;

const isBlob = (v: unknown): v is Blob => typeof Blob !== 'undefined' && v instanceof Blob;

const isReadableStream = (v: unknown): v is ReadableStream =>
	typeof ReadableStream !== 'undefined' && v instanceof ReadableStream;

// Plain object = not null, typeof 'object', and not any BodyInit
function isPlainObject(v: unknown): v is Record<string, unknown> {
	if (v === null || typeof v !== 'object') return false;
	return !(
		isFormData(v) ||
		isURLSearchParams(v) ||
		isBlob(v) ||
		isReadableStream(v) ||
		ArrayBuffer.isView(v) || // BufferSource (TypedArrays/DataView)
		v instanceof ArrayBuffer
	);
}

export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
	const headers = new Headers(init.headers ?? {});
	let body = init.body;

	// Only set Content-Type for real JSON payloads.
	if (body !== undefined && body !== null && !headers.has('content-type')) {
		if (isFormData(body)) {
			// Let the browser set: multipart/form-data; boundary=...
			// Ensure we DO NOT set a Content-Type.
			// (No-op: just don't touch the header)
		} else if (isURLSearchParams(body)) {
			// Let the browser set: application/x-www-form-urlencoded;charset=UTF-8
		} else if (isBlob(body)) {
			// If the Blob has a type, the browser will use it. No need to set.
			// If you *must* force it, you could read body.type, but best to leave it.
		} else if (typeof body === 'string') {
			// Caller provided a string; assume they know what it is. Don’t force JSON.
			// If you *want* to treat strings as JSON, set the header explicitly at call site.
		} else if (isPlainObject(body)) {
			// Caller passed a plain object → serialize to JSON and set header.
			headers.set('content-type', 'application/json');
			body = JSON.stringify(body);
		} else {
			// Other BodyInit (ArrayBuffer, TypedArray, ReadableStream, etc.) → do nothing.
		}
	}

	// CSRF: mirror cookie into header if caller didn’t set it
	const csrf = pickCsrfToken();
	if (csrf && !headers.has('X-Csrf-Token')) {
		headers.set('X-Csrf-Token', csrf);
	}

	const res = await fetch(urlOf(path), {
		credentials: init.credentials ?? 'include',
		mode: init.mode ?? 'cors',
		...init,
		body, // use possibly stringified JSON
		headers
	});

	if (!res.ok) {
		const bodyTextOrJson = await parse<any>(res).catch(() => {
			throw new Error(res.statusText);
		});
		throw new Error(
			typeof bodyTextOrJson === 'string'
				? bodyTextOrJson
				: JSON.stringify(bodyTextOrJson ?? res.statusText)
		);
	}

	return parse<T>(res);
}
