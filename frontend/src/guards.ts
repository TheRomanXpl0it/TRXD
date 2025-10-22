import { wrap } from 'svelte-spa-router/wrap';
import { currentUser, loadUser } from '$lib/stores/auth';

export const requireAuth = (component: any) =>
  wrap({
    component,
    conditions: [async (detail) => {
      await loadUser();                    // no-op if already ready
      if (currentUser()) return true;
      const dest = `/signIn?redirect=${encodeURIComponent(detail.location)}`;
      window.location.hash = `#${dest}`;   // hash router
      return false;
    }]
  });
