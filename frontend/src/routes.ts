import SignIn from './routes/signIn/+page.svelte';
import SignUp from './routes/signUp/+page.svelte';
import SignOut from './routes/signOut/+page.svelte';

import NotFound from './routes/404/+page.svelte'; // can be lazy too if you want

import { requireGuest, requireAuth, requireAuthLazy, requireAdminLazy } from './guards';

export default {
  // Auth pages: user lands on these pages so i want to load them statically
  '/signIn': requireGuest(SignIn),
  '/signUp': requireGuest(SignUp),
  '/signOut': requireAuth(SignOut),

  // Everything else: LAZY + protected as needed
  '/':        requireAuthLazy(() => import('./routes/+page.svelte')),
  '/writeups': requireAuthLazy(() => import('./routes/writeups/+page.svelte')),
  '/scoreboard': requireAuthLazy(() => import('./routes/scoreboard/+page.svelte')),
  '/challenges': requireAuthLazy(() => import('./routes/challenges/+page.svelte')),
  '/team':       requireAuthLazy(() => import('./routes/team/+page.svelte')),
  '/team/:id':   requireAuthLazy(() => import('./routes/team/+page.svelte')),
  '/account':    requireAuthLazy(() => import('./routes/account/+page.svelte')),
  '/account/:id':requireAuthLazy(() => import('./routes/account/+page.svelte')),

  // Admin-only lazy route
  '/configs': requireAdminLazy(() => import('./routes/configs/+page.svelte')),

  // 404 (do not lazy load this)
  '*': NotFound
};
