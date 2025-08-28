import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"
import { NotebookPen } from "lucide-react"
import { ChallengeForm } from "@/components/ChallengeForm"
import { useContext, useCallback, useEffect, useRef, useState } from "react"
import { AuthContext } from "@/context/AuthProvider"

export function CreateChallenge() {
  const { auth } = useContext(AuthContext)
  if (!auth) {
    console.error("Auth context is not available")
    return null
  }

  // ---- Resizable width state/logic (right-anchored sheet) ----
  const MIN_W = 360
  const [maxW, setMaxW] = useState<number>(
    typeof window !== "undefined" ? Math.min(1200, window.innerWidth) : 1200
  )
  const [width, setWidth] = useState<number>(640)
  const draggingRef = useRef(false)

  const clamp = (n: number, min: number, max: number) => Math.max(min, Math.min(max, n))

  const onPointerMove = useCallback((clientX: number) => {
    const winW = window.innerWidth
    const next = clamp(winW - clientX, MIN_W, maxW)
    setWidth(next)
  }, [maxW])

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
    const handleResize = () => {
      const newMax = Math.min(1200, window.innerWidth)
      setMaxW(newMax)
      setWidth(w => clamp(w, MIN_W, newMax))
    }

    window.addEventListener("mousemove", onMouseMove)
    window.addEventListener("mouseup", stopDrag)
    window.addEventListener("touchmove", onTouchMove, { passive: false })
    window.addEventListener("touchend", stopDrag)
    window.addEventListener("resize", handleResize)

    return () => {
      window.removeEventListener("mousemove", onMouseMove)
      window.removeEventListener("mouseup", stopDrag)
      window.removeEventListener("touchmove", onTouchMove)
      window.removeEventListener("touchend", stopDrag)
      window.removeEventListener("resize", handleResize)
    }
  }, [onMouseMove, onTouchMove, stopDrag])

  return (
    <Sheet>
      <SheetTrigger className="flex h-full">
        <span className="flex items-center cursor-pointer outline p-2 rounded-lg bg-green-500 text-white">
          <NotebookPen size={24} />
          <span className="ml-2">Create Challenge</span>
        </span>
      </SheetTrigger>

      {/* Control width via inline style; keep maxWidth to protect small screens */}
      <SheetContent
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
            <SheetTitle>Create a new Challenge</SheetTitle>
            <SheetDescription>
              So you want to create a new challenge? Fill out the form below to get started.
            </SheetDescription>
          </SheetHeader>

          <div className="mt-4 ml-4 mr-4">
            <ChallengeForm auth={auth} />
          </div>
        </div>
      </SheetContent>
    </Sheet>
  )
}
