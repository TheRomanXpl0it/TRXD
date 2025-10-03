"use client";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Medal } from "lucide-react";
import { ScoreboardEntry } from "./Scores";

import { useContext } from "react";
import { useNavigate } from "react-router-dom";
import SettingContext from "@/context/SettingsProvider";

export function Scoreboard({ data }: { data: ScoreboardEntry[] }) {
  const { settings } = useContext(SettingContext);
  const teamPlay = settings.General?.find(
    (setting) => setting.title === "Allow Team Play",
  )?.value;
  const navigate = useNavigate();
  const getMedalColor = (index: number) => {
    switch (index) {
      case 0:
        return "text-yellow-500"; // 🥇
      case 1:
        return "text-gray-400"; // 🥈
      case 2:
        return "text-orange-500"; // 🥉
      default:
        return "";
    }
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[100px]">Rank</TableHead>
          <TableHead>{teamPlay ? "Team" : "User"}</TableHead>
          <TableHead>Points</TableHead>
          <TableHead>Achievements</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.map((entry, index) => (
          <TableRow
            key={entry.playerId}
            className="cursor-pointer"
            onClick={() => {
              if (teamPlay)
                navigate(`/team/${encodeURIComponent(entry.playerId)}`);
              else
                navigate(`/account/${encodeURIComponent(entry.playerId)}`);
            }}
            tabIndex={teamPlay ? 0 : -1}
            onKeyDown={(e) => {
              if (!teamPlay) return;
              if (e.key === "Enter" || e.key === " ") {
                e.preventDefault();
                navigate(`/team/${encodeURIComponent(entry.playerId)}`);
              }
            }}
          >
            <TableCell className="font-medium flex items-center gap-2">
              {index < 3 && <Medal className={getMedalColor(index)} />}
              {entry.ranking}
            </TableCell>
            <TableCell>{entry.name}</TableCell>
            <TableCell>{entry.score}</TableCell>
            <TableCell>
              {entry.badges.length > 0
                ? entry.badges.map((badge, idx) => (
                    <span
                      key={idx}
                      className="inline-block bg-gray-200 text-gray-800 text-xs px-2 py-1 rounded-full mr-1 mb-1"
                      title={badge.description}
                    >
                      {badge.name}
                    </span>
                  ))
                : "-"}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
