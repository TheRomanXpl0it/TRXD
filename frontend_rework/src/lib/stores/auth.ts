import { writable, get } from 'svelte/store';
import { getInfo } from '$lib/auth';

export const user = writable<Awaited<ReturnType<typeof getInfo>>>(null);
export const authReady = writable(false);

export async function loadUser(force = false) {
  if (!force && get(authReady)) return;
  const userfetched:any = await getInfo();
  
  if (userfetched === "OK"){
    user.set(null);
    authReady.set(true);
    return;
  }
  
  try{
    user.set(userfetched);
  } catch(e) {
    user.set(null);
  }
  authReady.set(true);
}

export function currentUser() { return get(user); }