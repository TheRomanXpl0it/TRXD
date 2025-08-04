import { useContext, useEffect, useState } from "react";
import SettingContext from "@/context/SettingsProvider";
import { AuthContext, isTeam } from "@/context/AuthProvider";
import { fetchTeamData, leaveTeam } from "@/lib/backend-interaction";
import { Button } from "@/components/ui/button";
import { Scoreboard } from "@/components/Scoreboard";
import { teamColumns } from "@/components/columns/teamColumns";
import { useNavigate } from "react-router-dom";
import { Medal, LogOut } from "lucide-react";
import { toast } from "sonner";

function getMedalColor(rank: number) {
  switch (rank) {
    case 1: return "text-yellow-500";
    case 2: return "text-gray-400";
    case 3: return "text-orange-500";
    default: return "";
  }
}

export function CreateOrJoinTeam() {
  const navigate = useNavigate();
  return (
    <div className="p-6">
      <p className="mb-4">To participate in team challenges, create or join a team.</p>
      <div className="flex items-center justify-center gap-8 mt-6 flex-wrap">
        <div className="flex flex-col items-center">
          <button className="bg-blue-600 text-white rounded-lg shadow-lg hover:brightness-110 hover:cursor-pointer transition-all"
            onClick={() => navigate("/createteam")}>
            <img src="/createTeam.svg" alt="Create Team" className="w-full h-full rounded-lg" />
          </button>
            <span className="mt-2 text-xl font-bold text-blue-700">Create Team</span>
        </div>
        <div className="flex flex-col items-center">
          <button className="bg-pink-600 text-white rounded-lg shadow-lg hover:brightness-110 hover:cursor-pointer transition-all"
            onClick={() => navigate("/jointeam")}>
            <img src="/joinTeam.svg" alt="Join Team" className="w-full h-full rounded-lg" />
          </button>
          <span className="mt-2 text-xl font-bold text-pink-700">Join Team</span>
        </div>
      </div>
    </div>
  );
}

export function Team() {
  const { settings } = useContext(SettingContext);
  const { auth } = useContext(AuthContext);
  let [teamData, setTeamData] = useState<null | any>(null);
  const showQuotes = settings.General?.find((s) => s.title === 'Show Quotes')?.value;
  const allowTeamPlay = settings.General?.find((s) => s.title === 'Allow Team Play')?.value;
  const navigate = useNavigate();


  if (!auth) {
    toast.error("You must be logged in to view your team.");
    navigate("/login");
    return null;
  }

  useEffect(() => {
    if (auth.teamId !== null) {
      const response = fetchTeamData(auth.teamId);
      isTeam(response) ? setTeamData(response) : setTeamData(null);
      console.log("Fetched team data:", teamData);
    }
  }, [auth.teamId]);


  return (
    <>
      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
        { teamData ? teamData.name : "Join or create a team"}
      </h2>

      {showQuotes && (
        <blockquote className="mt-6 border-l-2 pl-6 italic">
          "None of us is as smart as all of us."
        </blockquote>
      )}

      {allowTeamPlay && console.log(teamData) && teamData && (
        <div className="p-6 max-w-3xl mx-auto text-center">
          <div className="flex justify-center mb-6">
            <img
              src="/team-logo.png"
              alt={`${teamData.name} Logo`}
              className="w-32 h-32 rounded-full border-4 border-gray-300 shadow-md"
            />
          </div>
          <h3 className="text-3xl font-bold mb-2">{teamData.name}</h3>
          <p className="text-gray-600 mb-6 flex items-center justify-center gap-2">
            Rank #{teamData.rank}
            {teamData.rank <= 3 && (
              <Medal className={`w-5 h-5 ${getMedalColor(teamData.rank)}`} />
            )}
            &bull; {teamData.score} Points
          </p>
          <div className="flex justify-end p-4">
            <Button variant="destructive" onClick={leaveTeam}>
              <LogOut size={24} />
              Leave Team
            </Button>
          </div>
          <Scoreboard
            columns={teamColumns}
            data={teamData.members.map((member: {id: number, name:string, role:string, score: number}) => ({
              name: member.name,
              score: member.score,
              role: member.role,
              id: member.id,
            }))}
          />
        </div>
      )}

      {allowTeamPlay && !teamData && <CreateOrJoinTeam />}

      {!allowTeamPlay && (
        <div className="p-4">
          <p className="text-lg">Team play is currently disabled.</p>
        </div>
      )}
    </>
  );
}