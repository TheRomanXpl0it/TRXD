import { writable, get, readable } from 'svelte/store';
import { getInfo } from '$lib/auth';

export const user = writable<Awaited<ReturnType<typeof getInfo>>>(null);
export const authReady = writable(false);
export const userMode = writable(true);

let loadingPromise: Promise<void> | null = null;

export async function loadUser(force = true) {
  if (!force && get(authReady)) return;

  // If already loading, wait for the existing promise
  if (loadingPromise) return loadingPromise;

  loadingPromise = (async () => {
    try {
      const userfetched:any = await getInfo();

      if (userfetched === "OK"){
        user.set(null);
        authReady.set(true);
        return;
      }

      try{
        userMode.set(userfetched.user_mode);
        user.set(userfetched);
      } catch(e) {
        user.set(null);
      }
      authReady.set(true);
    } finally {
      loadingPromise = null;
    }
  })();

  return loadingPromise;
}

export function clearUser(force = true){
  if(!authReady && !force){
    return
  }
  user.set(null);
  authReady.set(false);
  loadingPromise = null;
}

export function currentUser() { return get(user); }