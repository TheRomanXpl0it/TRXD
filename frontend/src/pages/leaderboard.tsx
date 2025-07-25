import { Chart } from "@/components/chart";
import { Scoreboard } from "@/components/scoreboard";
import { Player, columns } from "@/components/scores";
import SettingContext from "@/context/SettingsProvider";
import { useContext } from "react";

async function getData(): Promise<Player[]> {
    // Fetch data from your API here.
    return [
      {
        ranking: 1,
        username: "John Doe",
        score: 1000,
        badges: ["Reverse Engineer", "Web Hacker"],
      },
        {
            ranking: 2,
            username: "Jane Doe",
            score: 900,
            badges: ["Forensics Expert", "Crypto Hacker"],
        },
        {
            ranking: 3,
            username: "Alice",
            score: 800,
            badges: ["Pwn Master", "Misc Hacker"],
        },
        {
            ranking: 4,
            username: "Bob",
            score: 700,
            badges: ["Web Hacker", "Crypto Hacker"],
        },
        {
            ranking: 5,
            username: "Eve",
            score: 600,
            badges: ["Forensics Expert", "Pwn Master"],
        },
    ]
  }

import { useEffect, useState } from "react";

export function Leaderboard() {
    const { settings } = useContext(SettingContext);
    const [data, setData] = useState<Player[]>([]);

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
        {settings.General?.find((setting) => setting.title === 'Show Quotes')?.value && (
        <blockquote className="mt-6 border-l-2 pl-6 italic">
            "You gain strength, courage and confidence by every experience in which you really stop to look fear in the face. You must do the thing you think you cannot do."
        </blockquote>
        )}
        <div className="flex items-center justify-center w-full h-100">
                <Chart/>
        </div>
        <div className="p-4">
            <p className="text-lg font-semibold">Top Players</p>
        </div>
        <Scoreboard columns={columns} data={data} />
    </>
    )
}