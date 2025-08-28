import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Layout } from "@/Layout";
import { AuthProvider } from "@/context/AuthProvider";
import PrivateRoute from "@/context/PrivateRoute";
import { Home } from "@/pages/home";
import { Challenges } from "@/pages/challenges";
import { Login } from "@/pages/login";
import { Leaderboard } from "@/pages/leaderboard";
import { Writeups } from "@/pages/writeups";
import { CreateTeam } from "@/pages/createteam";
import { JoinTeam } from "@/pages/jointeam";
import { Account } from "@/pages/account";
import { Settings } from "@/pages/settings";
import { Team } from "@/pages/team";
import { ErrorPage } from "@/pages/error";
import { ErrorBoundary } from "@/components/ErrorBoundary";


import "@/App.css";


function App() {
  return (
    <Router>
      <ErrorBoundary fallback={<ErrorPage />}>
        <AuthProvider>
          <Routes>
            <Route element={<Layout />}>
              <Route path="/" element={<Home />} />
              <Route path="/leaderboard" element={<Leaderboard />} />
              <Route path="/writeups" element={<Writeups />} />
              <Route path="/login" element={<Login />} />
              <Route path="/account/:username" element={<Account />} />

              { /* Authenticated routes */ }
              <Route path="/challenges" element={<PrivateRoute><Challenges /></PrivateRoute>} />
              <Route path="/settings" element={<PrivateRoute><Settings /></PrivateRoute>} />
              <Route path="/account" element={<PrivateRoute><Account /></PrivateRoute>} />
              <Route path="/team" element={<PrivateRoute><Team /></PrivateRoute>} />
              <Route path="/team/:teamId" element={<PrivateRoute><Team /></PrivateRoute>} />
              <Route path="/createteam" element={<PrivateRoute><CreateTeam /></PrivateRoute>} />
              <Route path="/jointeam" element={<PrivateRoute><JoinTeam /></PrivateRoute>} />
            </Route>
          </Routes>
        </AuthProvider>
      </ErrorBoundary>
    </Router>
  );
}

export default App;