import { wrap } from 'svelte-spa-router/wrap';
import { currentUser, loadUser } from '$lib/stores/auth';
import { push } from 'svelte-spa-router';
import { toast } from 'svelte-sonner';
import { get } from 'svelte/store';
import { user as authUser } from '$lib/stores/auth';

export const requireAuth = (component: any) =>
  wrap({
    component,
    conditions: [async (detail) => {
      await loadUser();                    // no-op if already ready
      if (currentUser()) return true;
      push("/signIn")
      return false;
    }]
  });

export const requireGuest = (component: any) =>
  wrap({
    component,
    conditions: [async (detail) => {
      await loadUser();                    // no-op if already ready
      if (currentUser()){
        push("/challenges")
        toast.error("already signed in!")
        return false;
      }
      return true
    }]
  });

export const AdminConfigs = wrap({
  asyncComponent: () =>
    import('./routes/configs/+page.svelte').then((m) => m.default as any),
  conditions: [
    async () => {
      await loadUser();
      const u = get(authUser);
      if (!u) {
        push('/signIn');
        return false;
      }
      if (u.role !== 'Admin') {
        push('/');
        toast.error('Forbidden');
        return false;
      }
      return true;
    }
  ]
});