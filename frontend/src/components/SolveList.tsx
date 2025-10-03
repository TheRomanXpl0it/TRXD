import type { Solve } from '@/context/AuthProvider'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from './ui/table'

export function SolveList (
  {solves}
  : {
    solves : Solve[]
  }
) {
  return(
    <>
      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight mt-2">
        Challenge Solved
      </h2>
      <Table className="mt-4">
        <TableCaption>List of challenges solved</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead>Challenge</TableHead>
            <TableHead className="text-right">Points</TableHead>
            <TableHead className="text-right">Category</TableHead>
            <TableHead className="text-right">Solved At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          { solves.map((solve) => (
            <TableRow key={solve.id}>
              <TableCell>{solve.name}</TableCell>
              <TableCell className="text-right">{solve.points}</TableCell>
              <TableCell className="text-right">{solve.category}</TableCell>
              <TableCell className="text-right">{new Date(solve.timestamp).toLocaleString()}</TableCell>
            </TableRow>
          )) }
        </TableBody>
      </Table>
    </>
  )
}