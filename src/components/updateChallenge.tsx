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

export function UpdateChallenge(){
    return(
        <Sheet>
            <SheetTrigger>
                <span className="flex items-center space-x-2 cursor-pointer outline p-2 rounded-lg bg-green-500 text-white">
                    <NotebookPen size={24}/>
                    Create Challenge
                </span>
            </SheetTrigger>
            <SheetContent className="max-h-screen overflow-y-scroll">
                <SheetHeader>
                <SheetTitle>Create a new Challenge</SheetTitle>
                <SheetDescription>
                    So you want to create a new challenge? Fill out the form below to get started.
                </SheetDescription>
                </SheetHeader>
                <ChallengeForm/>
            </SheetContent>
        </Sheet>
    )
}