import { useContext } from "react"
import SettingContext from "@/context/SettingsProvider"
import AuthContext from "@/context/AuthProvider";
import { fetchTeamData } from "@/lib/backend-interaction"; // Assuming you have a function to fetch team data
import { Button } from "@/components/ui/button"
import { Scoreboard } from "@/components/scoreboard";
import { teamColumns } from "@/components/columns/teamColumns"; // defined above
import { useNavigate } from "react-router-dom";
import { Medal,LogOut } from "lucide-react";
import {leaveTeam} from "@/lib/backend-interaction"; // Assuming you have a function to leave a team

interface Team {
    name: string;
    members: string[];
    score: number;
    rank: number;
}



function getMedalColor(rank: number) {
  switch (rank) {
    case 1:
      return "text-yellow-500";
    case 2:
      return "text-gray-400";
    case 3:
      return "text-orange-500";
    default:
      return "";
  }
}



export function CreateOrJoinTeam() {
  const navigate = useNavigate();

  return (
    <div className="p-6">
      <p className="mb-4">
        To participate in team challenges, you can create a new team or join an existing one.
      </p>
      <div className="flex items-center justify-center gap-8 mt-6 flex-wrap">
        <div className="flex flex-col items-center">
          <button
            className="bg-blue-600 text-white rounded-lg shadow-lg text-lg hover:brightness-110 focus:outline-none transition-all glow-green"
            onClick={() => navigate("/createteam")}
          >
            <img src="/createTeam.svg" alt="Create Team Icon" className="inline-block w-full h-full rounded-lg" />
          </button>
          <span className="mt-2 text-xl text-blue-700 font-bold">Create Team</span>
        </div>

        <div className="flex flex-col items-center">
          <button
            className="bg-pink-600 text-white rounded-lg shadow-lg text-lg hover:brightness-110 focus:outline-none transition-all glow-blue"
            onClick={() => navigate("/jointeam")}
          >
            <img src="/joinTeam.svg" alt="Join Team Icon" className="inline-block w-full h-full rounded-lg" />
          </button>
          <span className="mt-2 text-xl text-pink-700 font-bold">Join Team</span>
        </div>
      </div>
    </div>
  );
}


function showTeam(team: Team) {
  const memberData = team.members.map((name) => ({
    name,
    score: Math.floor(Math.random() * 500) + 100, // Replace with real scores if available
  }));

  return (
      <div className="p-6 max-w-3xl mx-auto text-center">
        <div className="flex justify-center mb-6">
          <img
            src="/team-logo.png"
            alt={`${team.name} Logo`}
            className="w-32 h-32 rounded-full border-4 border-gray-300 shadow-md"
          />
        </div>
        <h3 className="text-3xl font-bold mb-2">{team.name}</h3>
        
        <p className="text-gray-600 mb-6 flex items-center justify-center gap-2">
          Rank #{team.rank}
          {team.rank <= 3 && (
            <Medal className={`w-5 h-5 ${getMedalColor(team.rank)}`} />
          )}
          &bull; {team.score} Points
        </p>
        <div className="flex justify-end p-4">
            <Button
              variant="destructive"
              onClick={leaveTeam}
            >
              <LogOut size={24} />
              Leave Team
            </Button>
        </div>


        <Scoreboard columns={teamColumns} data={memberData} />
      </div>
  );
}


const team = await fetchTeamData();

export function Team(){
    const {settings} = useContext(SettingContext);
    const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
    const allowTeamPlay = settings.General?.find((setting) => setting.title === 'Allow Team Play')?.value;
    const { auth } = useContext(AuthContext); // Uncomment if you need authentication context


    return(
        <>
            <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
                { team ? team.name : "Join or create a team" }
            </h2>
            { showQuotes && (
                <blockquote className="mt-6 border-l-2 pl-6 italic">
                    "None of us is as smart as all of us."
                </blockquote>
            )}
            { allowTeamPlay && team && showTeam(team) }
            { allowTeamPlay && !team && CreateOrJoinTeam() }
            { !allowTeamPlay && (
                <div className="p-4">
                    <p className="text-lg">Team play is currently disabled.</p>
                </div>
            )}
        </>
    )
}