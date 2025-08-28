import {
  DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle,
  DialogClose,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { CategoryIcon } from "./CategoryIcon";
import { FlagSubmit } from "./FlagSubmit";
import { UpdateChallenge } from "@/components/UpdateChallenge";
import { DeleteChallenge } from "@/components/DeleteChallenge";
import { UserPen, EyeClosed, Download, Droplet } from "lucide-react";
import { AuthProps } from "@/context/AuthProvider";
import { Challenge as ChallengeType } from "@/context/ChallengeProvider";
import { InstanceProvider } from "@/context/InstanceProvider";
import { InstanceControls } from "@/components/InstanceControls";
import { SolvesList } from "@/components/SolvesList";
import type { Solve } from "@/context/ChallengeProvider";

function ShowTags({tags}:{tags?: string[]}) {
  if (!tags?.length) return null;
  return (
    <div className="flex flex-wrap gap-1">
      {tags.map(t => <Badge key={t}>{t}</Badge>)}
    </div>
  );
}

function ShowSolves({solves=0,solves_list=[]}:{solves?: number, solves_list: Solve[]}) {
  return (
    <div className="flex justify-center">
      <span className="flex items-center space-x-2">
        {solves === 0 ? (
          <>
            <Droplet size={20} className="text-red-500" />
            <span>No solves yet</span>
          </>
        ) : (
          <SolvesList solves_list={solves_list} />
        )}
      </span>
    </div>
  );
}

function ShowAuthors({authors}:{authors?: string[]}) {
  if (!authors?.length) return null;
  return (
    <span className="flex flex-row items-center space-x-2 justify-center">
      <UserPen size={24}/>
      <div className="flex flex-col items-center">
        {authors.map((a,i) => <div key={i}>{a}</div>)}
      </div>
    </span>
  );
}

function Remote({remote}:{remote?: string}) {
  if (!remote) return null;
  return (
    <div className="flex items-center space-x-2">
      <span>Remote:</span><span>{remote}</span>
    </div>
  );
}

function Attachments({attachments}:{attachments?: File[]}) {
  if (!attachments?.length) return null;
  return (
    <div className="flex flex-row items-center gap-2">
      {attachments.map(a => (
        <Button key={a.name}>
          <Download size={24}/> {a.name}
        </Button>
      ))}
    </div>
  );
}

function ChallengeSolved(){
  return (
    <span className="flex text-green-500 font-semibold align-middle justify-center w-full">
      Challenge Solved
    </span>
  )
}

function ChallengeHidden(){
  return (
    <span className="flex items-center space-x-2 text-gray-500 justify-center mb-4">
          <EyeClosed size={24}/> This challenge is hidden
    </span>
  )
}

export default function ChallengeContent({
  challenge,
  challengeSettings,
  auth,
}: {
  challenge: ChallengeType;
  challengeSettings: {
    showSolves: boolean;
    showPoints: boolean;
    showTags: boolean;
    showDifficulty: boolean;
  };
  auth: AuthProps;
}) {
  const canEdit =
    auth.role === "Admin" ||
    (auth.role === "Author" && challenge.authors?.includes(auth.username));

    console.log("challenge", challenge);

  return (
    <DialogContent className="w-auto max-w-3xl" style={{ maxHeight: "90vh", overflowY: "auto" }}>
      <InstanceProvider challenge={challenge}>
        <DialogHeader>
        <DialogTitle>
          <div className="flex items-center space-x-2 justify-between w-full">
            <span className="flex-1 text-left w-33">{challenge.title}</span>
            <span className="flex-1 flex items-center justify-center space-x-2 w-33">
              <CategoryIcon category={challenge.category} size={24} />
            </span>
            {canEdit && (
              <div className="flex items-right justify-end space-x-2 mr-5 w-33">
                <UpdateChallenge challenge={challenge} auth={auth}/>
                <DeleteChallenge challenge={challenge} auth={auth}/>
              </div>
            )}
          </div>
          { challengeSettings.showTags && <ShowTags tags={challenge.tags} /> }
          { challenge.hidden && <ChallengeHidden /> }
        </DialogTitle>

        <div className="flex justify-between space-x-2">
          { challengeSettings.showSolves && <ShowSolves solves={challenge.solves} solves_list={challenge.solves_list ?? []} /> }
          <ShowAuthors authors={challenge.authors} />
        </div>

        <DialogDescription>{challenge.description}</DialogDescription>
        <DialogClose />
        </DialogHeader>

        { !challenge.instanced ? <Remote remote={challenge.remote} /> : null }
        <Attachments attachments={challenge.attachments as any} />

        <DialogFooter className="flex flex-col sm:flex-col space-y-2">
          <InstanceControls remote={challenge.remote} />
          { challenge.solved ? <ChallengeSolved /> : <FlagSubmit challenge={challenge} /> }
        </DialogFooter>
      </InstanceProvider>
    </DialogContent>
  );
}
