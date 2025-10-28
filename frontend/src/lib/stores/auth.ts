import { writable, get, readable } from 'svelte/store';
import { getInfo } from '$lib/auth';

export const user = writable<Awaited<ReturnType<typeof getInfo>>>(null);
export const authReady = writable(false);
export const userMode = writable(true);

export async function loadUser(force = false) {
  if (!force && get(authReady)) return;
  const userfetched:any = await getInfo();
  //console.log(userfetched);
  
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
}

export function currentUser() { return get(user); }