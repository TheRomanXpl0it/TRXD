import { HashRouter as Router, Routes, Route } from 'react-router-dom'
import { Layout } from './Layout'
import { Home } from './pages/home'
import { Login } from './pages/login'
import { Writeups } from './pages/writeups'
import { Leaderboard } from './pages/leaderboard'
import { Challenges } from './pages/challenges'
import { Settings } from './pages/settings'
import { Account } from './pages/account'
import { Logout } from './pages/logout'
import { Team } from './pages/team'
import { CreateTeam } from './pages/createteam'
import { JoinTeam } from './pages/jointeam'
import './App.css'


function App() {
  return (
    <Router>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<Home />} />
          <Route path="/challenges" element={<Challenges />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/writeups" element={<Writeups />} />
          <Route path="/login" element={<Login />} />
          <Route path="/settings" element={<Settings />} />
          <Route path="/account" element={<Account />} />
          <Route path="/team" element={<Team />} />
          <Route path="/logout" element={<Logout />} />
          <Route path="/createteam" element={<CreateTeam />} />
          <Route path="/jointeam" element={<JoinTeam />} />
        </Route>
      </Routes>
    </Router>
  );
}
export default App

