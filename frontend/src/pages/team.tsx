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

  const allowTeamPlay =
    settings.General?.find((s) => s.title === "Allow Team Play")?.value;
  const showQuotes =
    settings.General?.find((s) => s.title === "Show Quotes")?.value;

  // Parse /team/:teamId if present; else use user’s team id
  const routeTeamId = useMemo(() => {
    if (!teamIdParam) return null;
    const n = Number(teamIdParam);
    return Number.isFinite(n) ? n : null;
  }, [teamIdParam]);

  const effectiveTeamId = routeTeamId ?? auth?.teamId ?? null;
  const viewingOtherTeam =
    routeTeamId !== null && routeTeamId !== undefined && routeTeamId !== auth?.teamId;

  const [teamData, setTeamData] = useState<any | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  // If there is no :teamId and user is not logged in, redirect to login
  useEffect(() => {
    if (routeTeamId == null && !auth) {
      toast.error("You must be logged in to view your team.");
      navigate("/login");
    }
  }, [routeTeamId, auth, navigate]);

  // Fetch team data when effective id changes
  useEffect(() => {
    let cancelled = false;

    (async () => {
      if (effectiveTeamId == null) {
        setTeamData(null);
        return;
      }

      try {
        setLoading(true);
        const resp = await fetchTeamData(effectiveTeamId);
        if (cancelled) return;
        setTeamData(isTeam(resp) ? resp : null);
      } catch (e) {
        if (!cancelled) {
          setTeamData(null);
          console.error("Failed to load team:", e);
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [effectiveTeamId]);

  const handleLeaveTeam = async () => {
    try {
      await leaveTeam();
      toast.success("You left the team.");
      // After leaving, if you were on /team (own team), send to create/join
      // If you were viewing someone else's team, just refresh data
      if (!viewingOtherTeam) {
        navigate("/team");
      }
      setTeamData(null);
    } catch (e) {
      toast.error("Failed to leave team.");
      console.error(e);
    }
  };

  // ---------- Render ----------
  const title = teamData
    ? teamData.name
    : viewingOtherTeam
    ? "Team"
    : "Join or create a team";

  return (
    <>
      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
        {title}
      </h2>

      {showQuotes && (
        <blockquote className="mt-6 border-l-2 pl-6 italic">
          "None of us is as smart as all of us."
        </blockquote>
      )}

      {!allowTeamPlay && (
        <div className="p-4">
          <p className="text-lg">Team play is currently disabled.</p>
        </div>
      )}

      {allowTeamPlay && loading && (
        <div className="p-6 text-sm text-muted-foreground">Loading team…</div>
      )}

      {allowTeamPlay && !loading && teamData && (
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

          {/* Only show "Leave Team" if this is the viewer’s own team */}
          {!viewingOtherTeam && (
            <div className="flex justify-end p-4">
              <Button variant="destructive" onClick={handleLeaveTeam}>
                <LogOut size={24} />
                Leave Team
              </Button>
            </div>
          )}

          <Scoreboard
            columns={teamColumns}
            data={teamData.members.map(
              (member: { id: number; name: string; role: string; score: number }) => ({
                name: member.name,
                score: member.score,
                role: member.role,
                id: member.id,
              })
            )}
          />
        </div>
      )}

      {allowTeamPlay && !loading && !teamData && !viewingOtherTeam && <CreateOrJoinTeam />}

      {allowTeamPlay && !loading && !teamData && viewingOtherTeam && (
        <div className="p-6 text-sm text-muted-foreground">Team not found.</div>
      )}
    </>
  );
}
