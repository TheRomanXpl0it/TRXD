import { useContext, useEffect, useMemo, useState } from "react";
import SettingContext from "@/context/SettingsProvider";
import { AuthContext, isTeam } from "@/context/AuthProvider";
import { fetchTeamData, leaveTeam } from "@/lib/backend-interaction";
import { Button } from "@/components/ui/button";
import { Scoreboard } from "@/components/Scoreboard";
import { teamColumns } from "@/components/columns/teamColumns";
import { useNavigate, useParams } from "react-router-dom";
import { Medal, LogOut } from "lucide-react";
import { toast } from "sonner";
import { TeamView } from "@/components/TeamView";
import { UserCategoryRing } from "@/components/UserCategoryRing";
import type { Team } from "@/context/AuthProvider";


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
          <button
            className="bg-blue-600 text-white rounded-lg shadow-lg hover:brightness-110 hover:cursor-pointer transition-all"
            onClick={() => navigate("/createteam")}
          >
            <img src="/createTeam.svg" alt="Create Team" className="w-full h-full rounded-lg" />
          </button>
          <span className="mt-2 text-xl font-bold text-blue-700">Create Team</span>
        </div>
        <div className="flex flex-col items-center">
          <button
            className="bg-pink-600 text-white rounded-lg shadow-lg hover:brightness-110 hover:cursor-pointer transition-all"
            onClick={() => navigate("/jointeam")}
          >
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
  const navigate = useNavigate();
  const { teamId: teamIdParam } = useParams<{ teamId?: string }>();
  // Replace 'Team' with the correct Team data type, e.g., 'TeamType' or 'TeamData'
  // import or define the correct type if not already done
  const [teamData, setTeamData] = useState<Team | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  if (!auth){
    toast.error("You must be logged in to view team information.");
    navigate("/login");
  }

  const effectiveTeamId = teamIdParam ? parseInt(teamIdParam, 10) : auth.teamId;
  const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;

  useEffect(() => {
    const fetchTeam = async () => {
      if (!effectiveTeamId) return;
      setLoading(true);
      setError(null);

      try {
        const data = await fetchTeamData(effectiveTeamId);
        setTeamData(data);
      } catch (err) {
        setError("Failed to fetch team data");
      } finally {
        setLoading(false);
      }
    };

    fetchTeam();
  }, [teamIdParam]);

  return (
    <>
      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
          Team
      </h2>
      { showQuotes && (
          <blockquote className="mt-6 border-l-2 pl-6 italic">
            "None of us is as smart as all of us."
          </blockquote>
      )}

      { !teamData && <CreateOrJoinTeam /> }
      { teamData && <TeamView team={teamData} /> }
    </>
  );
}
