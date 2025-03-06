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
 
export const columns: ColumnDef<Player>[] = [
  {
    accessorKey: "ranking",
    header: "Ranking",
  },
  {
    accessorKey: "score",
    header: "Score",
  },
  {
    accessorKey: "username",
    header: "Username",
  },
  {
    accessorKey: "badges",
    header: "Achievements",
  },
]