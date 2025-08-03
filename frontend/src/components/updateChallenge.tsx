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

export function UpdateChallenge({auth,challenge}:{
    auth: AuthProps,
    challenge:ChallengeType
    }
 ){

    return(
        <Sheet>
            <SheetTrigger>
                <span className="flex items-center space-x-2 cursor-pointer outline p-2 rounded-lg">
                    <Pencil size={20}/>
                </span>
            </SheetTrigger>
            <SheetContent className="max-h-screen overflow-y-scroll">
                <SheetHeader>
                <SheetTitle>Update Challenge</SheetTitle>
                <SheetDescription>
                    Yup, you can update this challenge. Fill out the form below to get started.
                </SheetDescription>
                </SheetHeader>
                <span className="ml-4 mr-4">
                    <ChallengeForm auth={auth} challenge={challenge}/>
                </span>
            </SheetContent>
        </Sheet>
    )
}