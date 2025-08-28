import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"
import { ChallengeForm } from "@/components/ChallengeForm"
import { AuthProps } from "@/context/AuthProvider"
import { Pencil } from "lucide-react"
import { Challenge as ChallengeType } from "@/context/ChallengeProvider"
import { useCallback, useEffect, useRef, useState } from "react"

export function UpdateChallenge({
  auth,
  challenge,
}: {
  auth: AuthProps
  challenge: ChallengeType
}) {
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

  return (
    <Sheet>
      <SheetTrigger>
        <span className="flex items-center space-x-2 cursor-pointer outline p-2 rounded-lg">
          <Pencil size={20} />
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
            <SheetTitle>Update Challenge</SheetTitle>
            <SheetDescription>
              Yup, you can update this challenge. Fill out the form below to get started.
            </SheetDescription>
          </SheetHeader>

          <div className="mt-4 ml-4 mr-4">
            <ChallengeForm auth={auth} challenge={challenge} />
          </div>
        </div>
      </SheetContent>
    </Sheet>
  )
}
