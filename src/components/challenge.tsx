import {
    Card,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
  } from "@/components/ui/dialog"
import { Badge } from "@/components/ui/badge"
import { FlagSubmit } from "./flag-submit";
import { CategoryIcon } from "./category-icon";
import { Flag,UserPen, EyeClosed, Pencil, Trash, CircleCheck } from "lucide-react";
import { Button } from "@/components/ui/button";


export interface ChallengeProps {
    challenge: {
        id: number;
        title: string;
        description: string;
        solves?: number;
        points?: number;
        category: string;
        remote: string;
        solved: boolean;
        tags?: string[];
        difficulty?: string;
        attatchments?: File[];
        authors?: string[];
        hidden: Boolean;
    };
}

export function displayChallenges(
    challenges: ChallengeProps[],
    category: string, 
    settings:{
        title: string;
        value: boolean;
        description: string;
        type: BooleanConstructor;
    }[],
    auth: {
        username: string;
        password: string;
        accessToken: string;
        roles: string[];
    }
    ) {

    const challengeSettings = {
        showSolves : settings?.find((setting) => setting.title === 'Show Solves')?.value ?? false,
        showPoints : settings?.find((setting) => setting.title === 'Show Points')?.value ?? false,
        showTags : settings?.find((setting) => setting.title === 'Show Tags')?.value ?? false,
        showDifficulty : settings?.find((setting) => setting.title === 'Show Difficulty')?.value ?? false,
    }
    
    return challenges.map((challengeProp: ChallengeProps) => {
        if (challengeProp.challenge.category.includes(category)) {
            return (
                <Challenge
                    challengeProp = { challengeProp }
                    challengeSettings = { challengeSettings }
                    auth = { auth }
                />
            );
        }
    });
}

function getRemote(remote?: string) {
    if (remote) {
        return (
            <div className="flex items-center space-x-2">
                <span>Remote:</span>
                <span>{remote}</span>
            </div>
        )
    }
    return null;
}

function showTags(tags?: string[]) {
    if (!tags) return null;
    return tags.map((tag) => (
        <Badge key={tag}>{tag}</Badge>
    ));
}

function getAttatchments(attatchments: File[]) {
    return( attatchments.map((attatchment) => (
        <div key={attatchment.name}>
            <span>Attatchment:</span>
            <span>{attatchment.name}</span>
        </div>
        ))
    )
}

function showAuthors(authors: string[]) {
    return ( 
        <span className="flex items-center space-x-2 justify-center">
            <UserPen size={24}/>
            { authors.map(( author,index ) =>
            <div key={index}>
                <span>{author}</span>
            </div>
            )}
        </span>
    )
}

function showHidden(){
    return (
        <span className="flex items-center space-x-2 text-gray-500 justify-center mb-4">
            <EyeClosed size={24}/>
            This challenge is hidden
        </span>
    )
}

function showSolves(solves?: number){
    if (!solves) solves = 0;
    return(
        <div className="flex justify-center">
            <span className="flex items-center space-x-2">
                <Flag size={24} />
                {solves} solves
            </span>
        </div>
    )
}

function showPoints(points: number){
    return (
        <span className="text-sm">{points} points</span>
    )
}

function showTitle(title: string){
    return(
        <span className="flex-1 text-left w-33">
            {title}
        </span>
    )
}

function showCategory(category: string){
    return(
        <span className="flex-1 flex items-center justify-center space-x-2 w-33">
            <CategoryIcon category={category} size={24} />
        </span>
    )
}

function showControls(){
    return(
        <div className="flex items-right justify-end space-x-2 mr-5 w-33">
            <Button><Pencil size={24}/></Button>
            <Button variant="destructive"><Trash size={24}/></Button>
        </div>
    )
}



function Challenge({
    challengeProp,
    challengeSettings,
    auth
}: {
    challengeProp: ChallengeProps,
    challengeSettings: {
        showSolves: boolean;
        showPoints: boolean;
        showTags: boolean;
        showDifficulty: boolean;
    },
    auth: {
        username: string;
        password: string;
        accessToken: string;
        roles: string[];
    }
}) {
    const challenge = challengeProp.challenge;
    const canEdit = auth.roles.includes('admin') || (auth.roles.includes('author') && challenge.authors?.includes(auth.username));

    return (
        <Dialog key={challenge.id}>
            <DialogTrigger asChild>
                <Card className = {challenge.hidden ? "m-4 w-[250px] h-[130px] cursor-pointer  bg-gray-300" : "m-4 w-[250px] h-[130px] cursor-pointer"}>
                    <CardHeader>
                    <CardTitle>{challenge.title}</CardTitle>
                    <CardDescription>
                        { challengeSettings.showTags && showTags(challenge.tags) }
                    </CardDescription>
                    </CardHeader>
                    <CardFooter className="flex justify-between">
                        { challenge.points && showPoints(challenge.points) }

                        { challenge.solved && <CircleCheck size={24} className="text-green-500" /> }
                    </CardFooter>
                </Card>
            </DialogTrigger>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>
                        <div className="flex items-center space-x-2 justify-between w-full">
                            <>
                                { showTitle(challenge.title) }
                                { showCategory(challenge.category) }
                            </>
                            { canEdit && showControls()}
                        </div>
                        { challenge.tags && showTags(challenge.tags) }
                        { challenge.hidden && showHidden() }
                    </DialogTitle>
                    <div className="flex justify-end space-x-2">
                        { showSolves(challenge.solves) }
                        { challenge.authors && showAuthors(challenge.authors) }
                    </div>
                    <DialogDescription>{ challenge.description }</DialogDescription>
                    <DialogClose />
                </DialogHeader>
                { challenge.remote &&  getRemote(challenge.remote) }
                { challenge.attatchments && getAttatchments(challenge.attatchments) }
                <DialogFooter>
                    { challenge.solved ? 
                        <span className="flex text-green-500 font-semibold align-middle justify-center w-full">
                            Challenge Solved
                        </span>
                    : FlagSubmit() }
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}

