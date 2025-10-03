import { Chart } from "@/components/Chart";
import { Scoreboard } from "@/components/Scoreboard";
import type { ScoreboardEntry } from "@/components/Scores";
import SettingContext from "@/context/SettingsProvider";
import { useContext } from "react";
import { getScoreboard } from "@/lib/backend-interaction";
import { useNavigate } from "react-router-dom";
type ApiBadge = { name: string; description: string };
type ApiScoreboardEntry = {
  id: number;
  name: string;
  score: number;
  badges?: ApiBadge[];
};

async function getData(): Promise<ScoreboardEntry[]> {
  const scores = (await getScoreboard()) as ApiScoreboardEntry[] | [];
  const entries = Array.isArray(scores) ? scores : [];
  return entries.map((entry, idx) => ({
    ranking: idx + 1,
    playerId: entry.id,
    name: entry.name,
    score: entry.score,
    badges: Array.isArray(entry.badges)
      ? entry.badges.map((b) => ({
          name: b.name,
          description: b.description,
        }))
      : [],
  }));
}

import { useEffect, useState } from "react";

export function Leaderboard() {
  const { settings } = useContext(SettingContext);
  const [data, setData] = useState<ScoreboardEntry[]>([]);

  const showQuotes = settings.General?.find(
    (setting) => setting.title === "Show Quotes",
  )?.value;
  const teamPlay = settings.General?.find(
    (setting) => setting.title === "Allow Team Play",
  )?.value;

  useEffect(() => {
    async function fetchData() {
      const result = await getData();
      setData(result);
    }
    fetchData();
  }, []);
  return (
    <>
      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
        Leaderboard
      </h2>
      {showQuotes && (
        <blockquote className="mt-6 border-l-2 pl-6 italic">
          "You gain strength, courage and confidence by every experience in
          which you really stop to look fear in the face. You must do the thing
          you think you cannot do."
        </blockquote>
      )}
      <div className="flex items-center justify-center w-full h-100">
        <Chart />
      </div>
      <div className="p-4">
        <p className="text-lg font-semibold">
          Top {teamPlay ? "Teams" : "Players"}
        </p>
      </div>
      <Scoreboard data={data} />
    </>
  );
}
