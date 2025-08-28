"use client"

import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

import { Medal } from "lucide-react"

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
}

export function Scoreboard<TData, TValue>({
  columns,
  data,
}: DataTableProps<TData, TValue>) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  })

  const getMedalColor = (index: number) => {
    switch (index) {
      case 0:
        return "text-yellow-500" // 🥇
      case 1:
        return "text-gray-400"   // 🥈
      case 2:
        return "text-orange-500" // 🥉
      default:
        return ""
    }
  }

  const renderHeader = () => (
    <TableHeader>
      {table.getHeaderGroups().map((headerGroup) => (
        <TableRow key={headerGroup.id}>
          {headerGroup.headers.map((header) => (
            <TableHead key={header.id}>
              {!header.isPlaceholder &&
                flexRender(header.column.columnDef.header, header.getContext())}
            </TableHead>
          ))}
        </TableRow>
      ))}
    </TableHeader>
  )

  const renderRows = () => {
    const rows = table.getRowModel().rows

    if (!rows.length) {
      return (
        <TableRow>
          <TableCell colSpan={columns.length} className="h-24 text-center">
            No results.
          </TableCell>
        </TableRow>
      )
    }

    return rows.map((row) => {
      const isSelected = row.getIsSelected()

      return (
        <TableRow
          key={row.id}
          data-state={isSelected ? "selected" : undefined}
          className={isSelected ? "bg-muted" : ""}
        >
          {row.getVisibleCells().map((cell, index) => {
            const isFirstCell = index === 0
            const medalColor = getMedalColor(row.index)

            return (
              <TableCell key={cell.id}>
                <span className="inline-flex items-center gap-2">
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  {isFirstCell && row.index < 3 && (
                    <Medal className={`w-4 h-4 ${medalColor}`} />
                  )}
                </span>
              </TableCell>
            )
          })}
        </TableRow>
      )
    })
  }

  return (
    <div className="rounded-md border">
      <Table>
        {renderHeader()}
        <TableBody>{renderRows()}</TableBody>
      </Table>
    </div>
  )
}
