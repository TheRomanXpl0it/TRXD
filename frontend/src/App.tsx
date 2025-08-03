import { lazy, Suspense } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Layout } from "@/Layout";
import { AuthProvider } from "@/context/AuthProvider";
import Loading from "@/components/Loading";
import PrivateRoute from "@/context/PrivateRoute";
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
      <Router>
        <AuthProvider>
            <Routes>
              <Route element={<Layout />}>
                <Route path="/" element={<LazyHome />} />
                <Route path="/leaderboard" element={<LazyLeaderboard />} />
                <Route path="/writeups" element={<LazyWriteups />} />
                <Route path="/login" element={<LazyLogin />} /> 

                { /* Authenticated routes */ }
                <Route path="/challenges" element={<PrivateRoute><LazyChallenges /></PrivateRoute>} />
                <Route path="/settings" element={<PrivateRoute><LazySettings /></PrivateRoute>} />
                <Route path="/account" element={<PrivateRoute><LazyAccount /></PrivateRoute>} />
                <Route path="/team" element={<PrivateRoute><LazyTeam /></PrivateRoute>} />
                <Route path="/createteam" element={<PrivateRoute><LazyCreateTeam /></PrivateRoute>} />
                <Route path="/jointeam" element={<PrivateRoute><LazyJoinTeam /></PrivateRoute>} />
              </Route>
            </Routes>
          </AuthProvider>
        </Router>
    </Suspense>
  );
}

export default App;
