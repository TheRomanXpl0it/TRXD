import { lazy, Suspense } from "react";
import { Routes, Route } from "react-router-dom";
import { Layout } from "@/Layout";
import Loading from "@/components/loading";
import "@/App.css";

// Lazy loaded components
const LazyHome = lazy(() => import('@/pages/home').then(module => ({ default: module.Home })));
const LazyChallenges = lazy(() => import('@/pages/challenges').then(module => ({ default: module.Challenges })));
const LazyLeaderboard = lazy(() => import('@/pages/leaderboard').then(module => ({ default: module.Leaderboard })));
const LazyWriteups = lazy(() => import('@/pages/writeups').then(module => ({ default: module.Writeups })));
const LazyLogin = lazy(() => import('@/pages/login').then(module => ({ default: module.Login })));
const LazySettings = lazy(() => import('@/pages/settings').then(module => ({ default: module.Settings })));
const LazyAccount = lazy(() => import('@/pages/account').then(module => ({ default: module.Account })));
const LazyTeam = lazy(() => import('@/pages/team').then(module => ({ default: module.Team })));
const LazyCreateTeam = lazy(() => import('@/pages/createteam').then(module => ({ default: module.CreateTeam })));
const LazyJoinTeam = lazy(() => import('@/pages/jointeam').then(module => ({ default: module.JoinTeam })));

function App() {
  return (
    <Suspense fallback={<Loading />}>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<LazyHome />} />
          <Route path="/challenges" element={<LazyChallenges />} />
          <Route path="/leaderboard" element={<LazyLeaderboard />} />
          <Route path="/writeups" element={<LazyWriteups />} />
          <Route path="/login" element={<LazyLogin />} />
          <Route path="/settings" element={<LazySettings />} />
          <Route path="/account" element={<LazyAccount />} />
          <Route path="/team" element={<LazyTeam />} />
          <Route path="/createteam" element={<LazyCreateTeam />} />
          <Route path="/jointeam" element={<LazyJoinTeam />} />
        </Route>
      </Routes>
    </Suspense>
  );
}

export default App;
