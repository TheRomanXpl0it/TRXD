import Home from './routes/+page.svelte';
import SignIn from './routes/signIn/+page.svelte';
import Challenges from './routes/challenges/+page.svelte';
import Team from './routes/team/+page.svelte';
import Profile from './routes/account/+page.svelte';
import SignOut from './routes/signOut/+page.svelte';
import SignUp from './routes/signUp/+page.svelte';
import { requireAuth } from './guards';

export default {
  '/': Home,
  '/signIn': SignIn,
  '/signOut': SignOut,
  '/signUp': SignUp,
  '/challenges': requireAuth(Challenges),
  '/team': requireAuth(Team),
  '/team/:id': requireAuth(Team),
  '/account': requireAuth(Profile),
  '/account/:id': requireAuth(Profile),
  
};
