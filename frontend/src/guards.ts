// TODO: Refactor this file

import { wrap } from 'svelte-spa-router/wrap';
import { currentUser, loadUser, user as authUser } from '$lib/stores/auth';
import { get } from 'svelte/store';
import { push } from 'svelte-spa-router';
import { toast } from 'svelte-sonner';

// Eager guards (for already-imported components)
export const requireAuth = (component: any) =>
  wrap({
    component,
    conditions: [async () => {
      await loadUser(false);
      if (currentUser()) return true;
      push('/signIn');
      return false;
    }]
  });

export const requireGuest = (component: any) =>
  wrap({
    component,
    conditions: [async () => {
      await loadUser(false);
      if (currentUser()) {
        push('/challenges');
        toast.error('already signed in!');
        return false;
      }
      return true;
    }]
  });

// only import the component AFTER auth passes
export const requireAuthLazy = (importer: () => Promise<any>) =>
  wrap({
    // condition runs FIRST – blocks the import if not authed
    conditions: [async () => {
      await loadUser(false);
      if (currentUser()) return true;
      push('/signIn');
      return false;
    }],
    asyncComponent: () => importer().then((m) => m.default ?? m)
  });

// auth + role check before import
export const requireAdminLazy = (importer: () => Promise<any>) =>
  wrap({
    conditions: [async () => {
      await loadUser(false);
      const u = get(authUser);
      if (!u) { push('/signIn'); return false; }
      if (u.role !== 'Admin') { push('/'); toast.error('Forbidden'); return false; }
      return true;
    }],
    asyncComponent: () => importer().then((m) => m.default ?? m)
  });
