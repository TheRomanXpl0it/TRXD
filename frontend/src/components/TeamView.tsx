import { Globe, MedalIcon } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { Team, TeamMember, User } from "@/context/AuthProvider";
import { fetchUserData } from "@/lib/backend-interaction";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { UserCategoryRing } from "./UserCategoryRing";
import { useChallenges } from "@/context/ChallengeProvider";
import { Stats } from "./Stats";
import { SolveList } from "./SolveList";

// ---------------- Helpers ----------------
function initials(fullName?: string) {
  if (!fullName) return "?";
  const parts = fullName.trim().split(/\s+/);
  const a = parts[0]?.[0] ?? "";
  const b = parts[1]?.[0] ?? "";
  return (a + b).toUpperCase() || fullName[0]?.toUpperCase() || "?";
}

function TeamMember({ member }: { member: TeamMember }) {
  const [userData, setUserData] = useState<User | null>(null);
  const navigate = useNavigate();
  useEffect(() => {
    fetchUserData(member.id).then((data) => {
      setUserData(data);
    });
  }, [member.id]);

  if (!userData) return null;
  return (
    <Card
      onClick={() => {
        navigate(`/account/${userData.id}`);
      }}
      className="cursor-pointer"
    >
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Avatar>
            <AvatarImage src={userData.profilePicture} />
            <AvatarFallback>{initials(userData.name)}</AvatarFallback>
          </Avatar>
          {userData.name}
        </CardTitle>
      </CardHeader>
      <CardContent>{userData.country}</CardContent>
    </Card>
  );
}

function MemberList({ members }: { members: TeamMember[] }) {
  return (
    <>
      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight mt-2">
        Member list
      </h2>
      <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {members.map((member) => (
          <TeamMember key={member.id} member={member} />
        ))}
      </div>
    </>
  );
}

// ---------------- Component ----------------
function TeamView({ team }: { team: Team }) {
  const { challenges, categories: categoryNames } = useChallenges();
  const totalsByCategory = new Map<string, number>();
  categoryNames.forEach((cat) => {
    totalsByCategory.set(cat, 0);
  });
  challenges.forEach((ch) => {
    const cat = ch.category;
    totalsByCategory.set(cat, (totalsByCategory.get(cat) ?? 0) + 1);
  });
  const solvedByCategory = new Map<string, number>();
  (team.solves ?? []).forEach((s) => {
    const cat = s.category;
    solvedByCategory.set(cat, (solvedByCategory.get(cat) ?? 0) + 1);
  });
  const colorList = Array.from(
    { length: Math.max(1, totalsByCategory.size) },
    (_, i) =>
      `hsl(${Math.round((360 / Math.max(1, totalsByCategory.size)) * i)}, 70%, 50%)`,
  );
  let __idx = 0;
  const ringCategories = Array.from(totalsByCategory.entries())
    .map(([key, total]) => {
      const solved = solvedByCategory.get(key) ?? 0;
      const color = colorList[__idx++ % colorList.length];
      return { key, label: key, total, solved, color };
    })
    .filter((c) => c.total > 0);

  return (
    <div className="my-8">
      {/* Header */}
      <div className="flex flex-col">
        <div className="flex justify-center">
          <div className="flex items-center">
            <Avatar className="w-28 h-28">
              <AvatarImage src={team.logo} />
              <AvatarFallback className="text-3xl">
                {initials(team.name)}
              </AvatarFallback>
            </Avatar>

            <div className="flex items-center ml-4">
              <h1 className="text-3xl font-bold tracking-tight">{team.name}</h1>
            </div>

            {team.country && (
              <div className="ml-4 flex items-center gap-2 text-muted-foreground">
                <Globe className="w-4 h-4" />
                <span className="text-sm md:text-base">{team.country}</span>
              </div>
            )}
          </div>
        </div>

        {team.bio && (
          <p className="mt-3 text-sm md:text-base leading-relaxed text-muted-foreground">
            {team.bio}
          </p>
        )}
      </div>
      <div className="my-8 justify-center flex">
        <MedalIcon className="w-8 h-8 text-pink-500 mr-2 animate-bounce" />
        <h3 className="text-2xl font-bold"> This team is currently: GAY</h3>
      </div>
      <div className="my-8">
        <Stats 
          teamSolves={team.solves ?? []} 
          members={team.members} 
          categories={ringCategories}
          ringCategories={ringCategories}
          totalSolves={ringCategories.reduce(
            (s, c) => s + (c.solved ?? 0),
            0,
          )}
          badges={team.badges || []}
        />
      </div>
      <div className="my-8">
        <MemberList members={team.members} />
      </div>
      <div className="my-8">
        <SolveList solves={team.solves ?? []} />
      </div>
      
    </div>
  );
}

export { TeamView };
