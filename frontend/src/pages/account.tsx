import { NavLinkProps, useParams } from "react-router-dom";
import { useContext, useEffect, useState } from "react";
import { fetchUserData } from "@/lib/backend-interaction";
import { AuthContext, User } from "@/context/AuthProvider";
import Loading from "@/components/Loading";
import SettingContext from "@/context/SettingsProvider";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { UserCategoryRing } from "@/components/UserCategoryRing";
import type { Solve } from "@/context/AuthProvider";
import { useChallenges } from "@/context/ChallengeProvider";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

function userSolves(user: User) {
    if (!user.solves || user.solves.length === 0) return "No solves yet";
    return (
    <Table>
        <TableCaption>{user.name}'s Solves</TableCaption>
        <TableHeader>
            <TableRow>
              <TableHead>Challenge</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Solved at</TableHead>
              <TableHead>Points</TableHead>
            </TableRow>
        </TableHeader>
        <TableBody>
            {user.solves.map((solve) => (
              <TableRow key={solve.id}>
                <TableCell>{solve.name}</TableCell>
                <TableCell>{solve.category}</TableCell>
                <TableCell>{new Date(solve.timestamp).toLocaleString()}</TableCell>
                <TableCell>{solve.points}</TableCell>
              </TableRow>
            ))}
        </TableBody>
        </Table>
    );

}

function userBadges(user: User) {
  
}

export function Account() {
    const { username } = useParams<{ username?: string }>(); // undefined if you're on /account
    const { settings } = useContext(SettingContext);
    const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
    const allowTeamPlay = settings.General?.find((setting) => setting.title === 'Allow Team Play')?.value;
    const userId = username ? Number(username) : -1;
    const [user, setUser] = useState<User | null>(null);
    const auth = useContext(AuthContext).auth;

  useEffect(() => {
    let userData: User | null = { id: -1, name: "", role:"", score:-1, email: "", country: "", joinedAt: "", solves: [], teamId: null };
    if (username) {
      // Visiting someone else's profile
      (async () => {
        userData = await fetchUserData(userId);
        setUser(userData);
      })();
    } else {
      // Visiting your own profile
        if (auth) {
          (async () => {
            userData = await fetchUserData(auth.id);
            setUser(userData);
          })();
        }
    }
  }, [username]);

  const { challenges, categories: categoryNames } = useChallenges();
  const totalsByCategory = new Map<string, number>();
  if (!user) return <Loading />;
  categoryNames.forEach((cat) => {
    totalsByCategory.set(cat, 0);
  });
  challenges.forEach((ch) => {
    const cat = ch.category;
    totalsByCategory.set(cat, (totalsByCategory.get(cat) ?? 0) + 1);
  });
  const solvedByCategory = new Map<string, number>();
  (user.solves ?? []).forEach((s) => {
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
  console.log(user)

  return (
    <>
        <h2 className="scroll-m-20 border-b pb-2 text-3xl tracking-tight font-semibold first:mt-0 flex items-center gap-4">
            {user.name}'s Profile
        </h2>
        { showQuotes && (
            <blockquote className="mt-6 border-l-2 pl-6 italic">
            "You did not wake up to be mediocre."
            </blockquote>
        )}
        <div className="flex justify-center items-center mt-6 mr-10 gap-20">
          { !allowTeamPlay && <UserCategoryRing
            totalSolves={ringCategories.reduce(
              (s, c) => s + (c.solved ?? 0),
              0,
            )}
            teamSolves={user.solves as Solve[]}
            categories={ringCategories}
            size={100}
            strokeWidth={16}
            /> 
          }
          
          <div className="flex flex-row justify-between items-center">
            <div className="flex items-center mr-4">
                <h2 className="text-3xl font-semibold">{user.name}</h2>
                <p className="text-sm text-gray-500">{user.country}</p>
            </div>
            <Avatar className="h-32 w-32 rounded-full mb-2 text-3xl">
                <AvatarImage src={user.profilePicture || "/default-avatar.png"} alt={`${user.name}'s avatar`} />
                <AvatarFallback className="flex items-center justify-center h-full w-full bg-gray-200 text-gray-500">
              {user.name.charAt(0).toUpperCase()}
                </AvatarFallback>
            </Avatar>
          </div>
            <div className="mr-10">
                <p className="text-sm text-gray-500">Joined: {new Date(user.joinedAt).toLocaleDateString()}</p>
            </div>
        </div>
        <div className="text-center mb-4 mt-5">
            {userSolves(user)}
        </div>
    </>
  );
}
