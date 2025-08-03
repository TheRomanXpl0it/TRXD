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
import { useContext } from "react"
import AuthContext from "@/context/AuthProvider"

export function CreateChallenge(){
    const {auth} = useContext(AuthContext);
    if (!auth) {
        console.error("Auth context is not available");
        return null;
    }

    return(
        <Sheet>
            <SheetTrigger className="flex h-full">
                <span className="flex items-center cursor-pointer outline p-2 rounded-lg bg-green-500 text-white">
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
                <span className="ml-4 mr-4">
                    <ChallengeForm auth={auth}/>
                </span>
            </SheetContent>
        </Sheet>
    )
}