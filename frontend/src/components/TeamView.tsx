import { Globe } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { Team, TeamMember, User } from "@/context/AuthProvider";
import { fetchUserData } from "@/lib/backend-interaction";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { UserCategoryRing } from "./UserCategoryRing";


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
        fetchUserData(member.id).then(data => {
            setUserData(data);
        });
    }, [member.id]);

    if (!userData) return null;
    return (
        <Card onClick={() => {navigate(`/account/${userData.id}`);}} className="cursor-pointer">
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Avatar>
                        <AvatarImage src={userData.profilePicture} />
                        <AvatarFallback>{initials(userData?.username)}</AvatarFallback>
                    </Avatar>
                {userData.username}
                </CardTitle>
            </CardHeader>
            <CardContent>
                {userData.country}
            </CardContent>
        </Card>
    );
}

// ---------------- Component ----------------
function TeamView( {team }: {team: Team}) {
     const categories = [
        { key: "web", label: "Web", count: 18, color: "#10b981" },
        { key: "pwn", label: "Pwn", count: 9,  color: "#3b82f6" },
        { key: "rev", label: "Reversing", count: 6, color: "#a855f7" },
        { key: "crypto", label: "Crypto", count: 3, color: "#f59e0b" },
    ];
return (
    <div className="my-8 ">
        {/* Header: Big avatar on the left, details on the right */}
        <div className="flex flex-col">
        {/* Big team avatar */}
            <div className="flex justify-between">
                <div className="flex items-center">
                    <Avatar className="w-28 h-28">
                        <AvatarImage src={team.logo} />
                        <AvatarFallback className="text-3xl">{initials(team.name)}</AvatarFallback>
                    </Avatar>
                    <div className="flex items-center ml-4">
                        <h1 className="text-3xl font-bold tracking-tight">
                            {team.name}
                        </h1>
                    </div>
                    {team.country && (
                        <div className="flex items-center gap-2 text-muted-foreground justify-end">
                            <Globe className="w-4 h-4" />
                            <span className="text-sm md:text-base">{team.country}</span>
                        </div>
                    )}
                </div>
                <div>
                    <UserCategoryRing
                    totalSolves={categories.reduce((s, c) => s + c.count, 0)}
                    categories={categories}
                    size={100}         // optional
                    strokeWidth={16}   // optional
                    />
                </div>
            </div>
            <div className="flex flex-col">
            {team.bio && (
                <p className="text-sm md:text-base leading-relaxed text-right md:text-left text-muted-foreground">
                    {team.bio}
                </p>
            )}
            </div>
        </div>
        <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight mt-2">
            Member list
        </h2>
        <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {team.members.map(member => (
                <TeamMember key={member.id} member={member} />
            ))}
        </div>
    </div>
    
);}

export { TeamView };