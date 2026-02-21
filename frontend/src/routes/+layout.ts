import { loadUser, authState } from '$lib/stores/auth';
import { redirect, type RequestEvent } from '@sveltejs/kit';

export const ssr = false; // Disable SSR since TRXd relies heavily on client-side state/API 
export const prerender = false;

export async function load({ url }: { url: URL }) {
    // Let the user state initialize
    await loadUser(false);

    // Public routes
    const publicRoutes = ['/signIn', '/signUp'];
    const isPublicRoute = publicRoutes.some(r => url.pathname.startsWith(r));

    if (!authState.user && !isPublicRoute) {
        throw redirect(302, '/signIn');
    }

    if (authState.user && isPublicRoute) {
        throw redirect(302, '/challenges');
    }

    return {};
}
