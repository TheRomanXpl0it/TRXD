"use client"
 
import { ColumnDef } from "@tanstack/react-table"
 
// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Player = {
  ranking: number
  score: number
  badges: string[]
  username: string
}

export type Badge = {
  name: string;
  description: string;
}

export type ScoreboardEntry = {
  ranking: number;
  playerId: number;
  name: string;
  score: number;
  badges: Badge[];
}
 
export const columns: ColumnDef<ScoreboardEntry>[] = [
  {
    accessorKey: "ranking",
    header: "Ranking",
  },
  {
    accessorKey: "score",
    header: "Score",
  },
  {
    accessorKey: "name",
    header: "Username",
  },
  {
    accessorKey: "badges",
    header: "Achievements",
  },
]