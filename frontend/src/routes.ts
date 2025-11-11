import SignIn from './routes/signIn/+page.svelte';
import SignUp from './routes/signUp/+page.svelte';
import SignOut from './routes/signOut/+page.svelte';
import Home from './routes/+page.svelte';
import Challenges from './routes/challenges/+page.svelte';
import Team from './routes/team/+page.svelte';
import Teams from './routes/teams/+page.svelte';
import Account from './routes/account/+page.svelte';
import Scoreboard from './routes/scoreboard/+page.svelte';

import NotFound from './routes/404/+page.svelte';

import { requireGuest, requireAuth, requireAuthLazy, requireAdminLazy } from './guards';

export default {
  // user lands on these pages so i want to load them statically
  '/signIn': requireGuest(SignIn),
  '/signUp': requireGuest(SignUp),
  '/signOut': requireAuth(SignOut),

  // Main pages eagerly loaded to prevent flickering
  '/':        requireAuth(Home),
  '/challenges': requireAuth(Challenges),
  '/scoreboard': requireAuth(Scoreboard),
  '/teams':      requireAuth(Teams),
  '/team':       requireAuth(Team),
  '/team/:id':   requireAuth(Team),
  '/account':    requireAuth(Account),
  '/account/:id':requireAuth(Account),

  // Less frequently used pages lazy loaded
  '/writeups': requireAuthLazy(() => import('./routes/writeups/+page.svelte')),

  // Admin only lazy route
  '/configs': requireAdminLazy(() => import('./routes/configs/+page.svelte')),

  // 404
  '*': NotFound
};
