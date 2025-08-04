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
import { Button } from "@/components/ui/button"
import { FlagSubmit } from "./FlagSubmit";
import { CategoryIcon } from "./CategoryIcon";
import { Flag,UserPen, EyeClosed, CircleCheck, Download, Droplet } from "lucide-react";
import { UpdateChallenge } from "@/components/UpdateChallenge";
import { DeleteChallenge } from "@/components/DeleteChallenge";
import { AuthProps } from "@/context/AuthProvider";
import { Challenge as ChallengeType } from "@/context/ChallengeProvider";



export function displayChallenges(
    challenges: ChallengeType[],
    category: string,
    settings: {
        title: string;
        value: boolean;
        description: string;
        type: BooleanConstructor;
    }[],
    auth: AuthProps
    ) {

    const challengeSettings = {
        showSolves : settings?.find((setting) => setting.title === 'Show Solves')?.value ?? false,
        showPoints : settings?.find((setting) => setting.title === 'Show Points')?.value ?? false,
        showTags : settings?.find((setting) => setting.title === 'Show Tags')?.value ?? false,
        showDifficulty : settings?.find((setting) => setting.title === 'Show Difficulty')?.value ?? false,
    }
    
    return challenges.map((challenge: ChallengeType) => {
        if (challenge.category.includes(category)) {
            return (
                <Challenge
                    challenge = { challenge }
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
        <Badge key={tag} className="mx-0.5">{tag}</Badge>
    ));
}

function getAttachments(attachments: File[]) {
    return(
        <div className="flex flex-row items-center space-x-2">
            {attachments.map((attachment) => (
            <Button key={attachment.name}>
                <Download size={24}/> {attachment.name}
            </Button>
            ))}
        </div>
    )
}

function showAuthors(authors: string[]) {
    return ( 
        <span className="flex flex-row items-center space-x-2 justify-center">
            <UserPen size={24}/>
            <div className="flex flex-col items-center">
                { authors.map(( author,index ) =>
                <div key={index} >
                    {author}
                </div>
            )}
            </div>
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
                {solves === 0 ? (
                    <>
                        <Droplet size={20} className="text-red-500" />
                        <span>No solves yet</span>
                    </>
                ) : solves === 1 ? (
                    <>
                        <Flag size={20} className="text-blue-500" />
                        <span>{solves} solve</span>
                    </>
                ) : (
                    <>
                        <Flag size={20} className="text-blue-500" />
                        <span>{solves} solves</span>
                    </>
                )}
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

function showControls(
    auth: AuthProps,
    challenge: ChallengeType
){
    return(
        <div className="flex items-right justify-end space-x-2 mr-5 w-33">
            <UpdateChallenge challenge={challenge} auth={auth}/>
            <DeleteChallenge challenge={challenge} auth={auth}/>
        </div>
    )
}



function Challenge({
    challenge,
    challengeSettings,
    auth
}: {
    challenge: ChallengeType,
    challengeSettings: {
        showSolves: boolean;
        showPoints: boolean;
        showTags: boolean;
        showDifficulty: boolean;
    },
    auth: AuthProps
}) {
    const canEdit = auth.role === 'Admin' || (auth.role === 'Author' && challenge.authors?.includes(auth.username));
    
    return (
        <Dialog key={challenge.id}>
            <DialogTrigger asChild>
                <Card className={challenge.hidden ? "m-4 w-[250px] h-[130px] cursor-pointer border-dashed border-2" : "m-4 w-[250px] h-[130px] cursor-pointer"}>
                    <CardHeader className="px-4 flex flex-col justify-between h-full">
                        <CardTitle>{challenge.title}</CardTitle>
                        <CardDescription className=" h-5">
                            {challengeSettings.showTags && showTags(challenge.tags)}
                        </CardDescription>
                    </CardHeader>
                    <CardFooter className="px-4 pt-2 flex justify-between items-end mt-auto">
                        {challenge.points && showPoints(challenge.points)}
                        {challenge.solved && <CircleCheck size={24} className="text-green-500" />}
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
                            { canEdit && showControls(auth, challenge)}
                        </div>
                        { challenge.tags && showTags(challenge.tags) }
                        { challenge.hidden && showHidden() }
                    </DialogTitle>
                    <div className="flex justify-between space-x-2">
                        { showSolves(challenge.solves) }
                        { challenge.authors && showAuthors(challenge.authors) }
                    </div>
                    <DialogDescription>{ challenge.description }</DialogDescription>
                    <DialogClose />
                </DialogHeader>
                { !challenge.instanced ? challenge.remote &&  getRemote(challenge.remote) : null }
                { challenge.timeout!==undefined ? challenge.remote &&  getRemote(challenge.remote) : null }
                { challenge.attachments && getAttachments(challenge.attachments) }
                <DialogFooter>
                    { challenge.solved ? 
                        <span className="flex text-green-500 font-semibold align-middle justify-center w-full">
                            Challenge Solved
                        </span>
                    : <FlagSubmit challenge={challenge} /> }
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}