import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"
import { Droplet, Flag } from "lucide-react"
import { useCallback, useEffect, useRef, useState } from "react"
import { useNavigate } from "react-router-dom"
import type { Solve } from "@/context/ChallengeProvider"
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Link } from "react-router-dom"

type SolvesListProps = {
  solves_list: Solve[];
};

export function SolvesTable({ solves_list }: SolvesListProps) {
    const sortedSolves = [...solves_list].sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());
    const navigate = useNavigate();
  return (
    <Table>
      <TableHeader>
        <TableRow>  
          <TableHead>#</TableHead>
          <TableHead>Team Name</TableHead>
          <TableHead>Timestamp</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {sortedSolves.map((solve) => (  
            <TableRow 
                key={solve.id} 
                className="cursor-pointer"
                onClick={() => navigate(`/team/${encodeURIComponent(solve.id)}`)}
                onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                    e.preventDefault();
                    navigate(`/team/${encodeURIComponent(solve.id)}`);
                }
            }}>
                <TableCell>
                    {sortedSolves.findIndex(s => s.id === solve.id) === 0 ? (
                        <Droplet size={20} className="text-red-500" />
                    ) : (
                        sortedSolves.findIndex(s => s.id === solve.id) + 1
                    )}
                </TableCell>
                <TableCell>{solve.name}</TableCell>
                <TableCell>{new Date(solve.timestamp).toLocaleString()}</TableCell>
            </TableRow>
      ))}
      </TableBody>
    </Table>
  )
}


export function SolvesList({ solves_list }: SolvesListProps) {
  // width in pixels; tweak defaults as you like
  const MIN_W = 360
  const MAX_W = Math.min(1200, typeof window !== "undefined" ? window.innerWidth : 1200)
  const [width, setWidth] = useState<number>(640)
  const draggingRef = useRef(false)

  const clamp = (n: number, min: number, max: number) => Math.max(min, Math.min(max, n))

  const onPointerMove = useCallback((clientX: number) => {
    // right-anchored sheet -> width = windowWidth - pointerX
    const winW = window.innerWidth
    const next = clamp(winW - clientX, MIN_W, MAX_W)
    setWidth(next)
  }, [])

  const onMouseMove = useCallback((e: MouseEvent) => {
    if (!draggingRef.current) return
    onPointerMove(e.clientX)
  }, [onPointerMove])

  const onTouchMove = useCallback((e: TouchEvent) => {
    if (!draggingRef.current) return
    const t = e.touches[0]
    if (t) onPointerMove(t.clientX)
  }, [onPointerMove])

  const stopDrag = useCallback(() => {
    draggingRef.current = false
    document.body.style.userSelect = ""
    document.body.style.cursor = ""
  }, [])

  const startDrag = useCallback((clientX: number) => {
    draggingRef.current = true
    document.body.style.userSelect = "none"
    document.body.style.cursor = "col-resize"
    onPointerMove(clientX)
  }, [onPointerMove])

  useEffect(() => {
    window.addEventListener("mousemove", onMouseMove)
    window.addEventListener("mouseup", stopDrag)
    window.addEventListener("touchmove", onTouchMove, { passive: false })
    window.addEventListener("touchend", stopDrag)
    window.addEventListener("resize", () => {
      // keep width within bounds on window resize
      setWidth(w => clamp(w, MIN_W, Math.min(MAX_W, window.innerWidth)))
    })
    return () => {
      window.removeEventListener("mousemove", onMouseMove)
      window.removeEventListener("mouseup", stopDrag)
      window.removeEventListener("touchmove", onTouchMove)
      window.removeEventListener("touchend", stopDrag)
      // resize listener not stored; fine for simple case
    }
  }, [onMouseMove, onTouchMove, stopDrag])

  const solveCount = solves_list.length;

  return (
    <Sheet>
      <SheetTrigger>
        <span className="flex items-center space-x-2 cursor-pointer outline p-2 rounded-lg">
          <Flag className="text-blue-500" size={20}/>
          <span> { solveCount > 1 ? `${solveCount} solves` : "1 solve" }</span>
        </span>
      </SheetTrigger>

      {/* Override width via style. Remove any sm:max-w-* you may have added elsewhere */}
      <SheetContent
        // side="right" // default
        className="max-h-screen overflow-y-scroll p-0"
        style={{ width, maxWidth: "100vw" }}
      >
        {/* Drag handle on the left edge */}
        <div
          className="absolute left-0 top-0 h-full w-1 cursor-col-resize hover:bg-muted/60 active:bg-muted/80 transition-colors"
          onMouseDown={(e) => startDrag(e.clientX)}
          onTouchStart={(e) => {
            const t = e.touches[0]
            if (t) startDrag(t.clientX)
          }}
        />

        <div className="p-6">
          <SheetHeader>
            <SheetTitle>Solves</SheetTitle>
            <SheetDescription>
              Here's a list of teams who solved this challenge.
            </SheetDescription>
          </SheetHeader>

          <div className="mt-4 ml-4 mr-4">
            <SolvesTable solves_list={solves_list} />
          </div>
        </div>
      </SheetContent>
    </Sheet>
  )
}
