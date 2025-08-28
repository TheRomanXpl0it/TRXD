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

function userSolves(user: User) {
    if (!user.solves || user.solves.length === 0) return "No solves yet";
    return (
    <Table>
        <TableCaption>{user.username}'s Solves</TableCaption>
        <TableHeader>
            <TableRow>
            <TableHead className="w-[100px]">Challenge</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Method</TableHead>
            <TableHead className="text-right">Points</TableHead>
            </TableRow>
        </TableHeader>
        <TableBody>
            <TableRow>
            <TableCell className="font-medium">INV001</TableCell>
            <TableCell>Paid</TableCell>
            <TableCell>Credit Card</TableCell>
            <TableCell className="text-right">$250.00</TableCell>
            </TableRow>
        </TableBody>
        </Table>
    );

}

export function Account() {
    const { username } = useParams<{ username?: string }>(); // undefined if you're on /account
    const { settings } = useContext(SettingContext);
    const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
    const userId = username ? Number(username) : -1;
    const [user, setUser] = useState<User | null>(null);
    const auth = useContext(AuthContext).auth;

  useEffect(() => {
    let userData: User | null = { id: -1, username: "", role:"", score:-1, email: "", country: "", joinedAt: "", solves: [], teamId: null };
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

  if (!user) return <Loading />;
  const categories = [
    { key: "web", label: "Web", count: 18, color: "#10b981" },
    { key: "pwn", label: "Pwn", count: 9,  color: "#3b82f6" },
    { key: "rev", label: "Reversing", count: 6, color: "#a855f7" },
    { key: "crypto", label: "Crypto", count: 3, color: "#f59e0b" },
  ];

  return (
    <>
        <h2 className="scroll-m-20 border-b pb-2 text-3xl tracking-tight font-semibold first:mt-0 flex items-center gap-4">
            {user.username}'s Profile
        </h2>
        { showQuotes && (
            <blockquote className="mt-6 border-l-2 pl-6 italic">
            "You did not wake up to be mediocre."
            </blockquote>
        )}
        <div className="flex flex-row justify-between mt-6 mr-10">
            <UserCategoryRing
              totalSolves={categories.reduce((s, c) => s + c.count, 0)}
              categories={categories}
              size={100}         // optional
              strokeWidth={16}   // optional
            />
          <div className="flex flex-row justify-between items-center">
            <div className="flex items-center mr-4">
                <h2 className="text-3xl font-semibold">{user.username}</h2>
                <p className="text-sm text-gray-500">{user.country}</p>
            </div>
            <Avatar className="h-32 w-32 rounded-full mb-2  text-3xl">
                <AvatarImage src={user.profilePicture || "/default-avatar.png"} alt={`${user.username}'s avatar`} />
                <AvatarFallback className="flex items-center justify-center h-full w-full bg-gray-200 text-gray-500">
                  {user.username.charAt(0).toUpperCase()}
                </AvatarFallback>
            </Avatar>
          </div>
        </div>
        <div className="w-full text-right mr-10">
            <p className="text-sm text-gray-500">Joined: {new Date(user.joinedAt).toLocaleDateString()}</p>
        </div>
        <div className="text-center mb-4 mt-5">
            {userSolves(user)}
        </div>
    </>
  );
}
