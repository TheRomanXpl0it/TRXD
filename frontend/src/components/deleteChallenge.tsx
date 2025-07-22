import { ChallengeProps } from "./challenge";
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { AuthProps } from "@/context/AuthProvider";
import { Trash } from "lucide-react"



export function DeleteChallenge({challengeProp, auth} : {challengeProp: ChallengeProps, auth: AuthProps}){
    return(
    <Dialog>
        <DialogTrigger asChild>
            <Button variant="destructive"><Trash size={24}/></Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
            <DialogTitle>Delete {challengeProp.challenge.title}</DialogTitle>
            <DialogDescription>
                Are you sure you want to delete this challenge? This action cannot be undone.
            </DialogDescription>
            </DialogHeader>
            <DialogFooter className="flex items-center justify-between">
                <Button type="submit" variant="destructive" className="mr-auto">Delete</Button>
                <Button type="button" className="ml-auto">Cancel</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
    )
}