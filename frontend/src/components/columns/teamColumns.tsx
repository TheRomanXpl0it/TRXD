import { ColumnDef } from "@tanstack/react-table";


type TeamMember = {
  name: string;
  score: number;
};

export const teamColumns: ColumnDef<TeamMember>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => <span className="font-medium">{row.getValue("name")}</span>,
  },
  {
    accessorKey: "score",
    header: "Score",
    cell: ({ row }) => <span className="text-right">{row.getValue("score")}</span>,
  },
];
