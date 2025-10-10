
const BASE = '/api';

function urlOf(path: string) {
  return /^https?:\/\//.test(path) ? path : `${BASE}${path.startsWith('/') ? '' : '/'}${path}`;
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
  // Try common names; adjust to match your backend if needed
  return (
    getCookie('csrf_') ||        // seen in your request
    getCookie('csrftoken') ||    // Django-style
    getCookie('XSRF-TOKEN') ||   // many Node stacks
    getCookie('cors_') ||        // your older naming
    null
  );
}

export async function api<T>(
  path: string,
  init: RequestInit = {}
): Promise<T> {
  const headers = new Headers(init.headers ?? {});
  if (init.body && !headers.has('content-type')) {
    headers.set('content-type', 'application/json');
  }

  // Mirror CSRF token from cookie into typical headers if caller didn’t set one
  const csrf = pickCsrfToken();
  if (csrf) {
    if (!headers.has('X-CSRF-Token')) headers.set('X-CSRF-Token', csrf);
    if (!headers.has('X-CSRFToken')) headers.set('X-CSRFToken', csrf);
    if (!headers.has('X-XSRF-TOKEN')) headers.set('X-XSRF-TOKEN', csrf);
    // keep legacy header if your backend used it
    if (!headers.has('X-CORS')) headers.set('X-CORS', csrf);
  }

  const res = await fetch(urlOf(path), {
    credentials: init.credentials ?? 'include',
    mode: init.mode ?? 'cors',
    ...init,
    headers
  });

  if (!res.ok) {
    const body = await parse<any>(res).catch(() => undefined);
    throw new Error(typeof body === 'string' ? body : JSON.stringify(body ?? res.statusText));
  }

  return parse<T>(res);
}
