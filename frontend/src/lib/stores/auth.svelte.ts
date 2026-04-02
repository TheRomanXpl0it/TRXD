import { getInfo } from '$lib/auth';

type UserDataType = Awaited<ReturnType<typeof getInfo>>;

export const authState = $state<{
	user: UserDataType | null;
	ready: boolean;
	userMode: boolean;
	emailVerification: boolean;
	startTime: string | null;
	endTime: string | null;
}>({
	user: null,
	ready: false,
	userMode: true,
	emailVerification: true,
	startTime: null,
	endTime: null
});

let loadingPromise: Promise<void> | null = null;

export async function loadUser(force = true) {
	if (!force && authState.ready) return;

	if (loadingPromise) return loadingPromise;

	loadingPromise = (async () => {
		try {
			const userfetched: any = await getInfo();

			if (userfetched === 'OK') {
				authState.user = null;
				authState.ready = true;
				return;
			}

			if (userfetched) {
				authState.emailVerification = userfetched.email_verification ?? true;
				authState.startTime = userfetched.start_time ?? null;
				authState.endTime = userfetched.end_time ?? null;
			}

			if (userfetched && userfetched.id) {
				authState.userMode = userfetched.user_mode;
				authState.user = userfetched;
			} else {
				authState.user = null;
			}
			authState.ready = true;
		} finally {
			loadingPromise = null;
		}
	})();

	return loadingPromise;
}

export function clearUser(force = true) {
	if (!authState.ready && !force) {
		return;
	}
	authState.user = null;
	authState.ready = false;
	loadingPromise = null;
}

export function currentUser() {
	return authState.user;
}
