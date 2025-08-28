import {
    Card,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    Dialog,
    DialogTrigger,
  } from "@/components/ui/dialog"
import { Badge } from "@/components/ui/badge"
import { CircleCheck, Droplet } from "lucide-react";
import { AuthProps } from "@/context/AuthProvider";
import { Challenge as ChallengeType } from "@/context/ChallengeProvider";
import { Suspense, lazy, useState } from "react";
import { fetchChallengeById } from "@/lib/backend-interaction"; // path as per your project
import { DialogLoading } from "@/components/DialogLoading";
import { useRef } from "react";
import { TimeoutBadge } from "@/components/TimeoutBadge";


const LazyChallengeContent = lazy(() =>
  import("./ChallengeContent") // webpack will split this
  // .then(m => ({ default: m.default })) // not needed, default export
);



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

function showTags(tags?: string[]) {
    if (!tags) return null;
    return tags.map((tag) => (
        <Badge key={tag} className="mx-0.5">{tag}</Badge>
    ));
}



function showPoints(points: number){
    return (
        <span className="text-sm">{points} points</span>
    )
}


function Challenge({
  challenge,
  challengeSettings,
  auth
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
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [fullChallenge, setFullChallenge] = useState<ChallengeType | null>(null);
  const abortRef = useRef<AbortController | null>(null);

  const openWithData = async () => {
    setLoading(true);
    abortRef.current?.abort();
    const ctrl = new AbortController();
    abortRef.current = ctrl;

    try {
      const data = await fetchChallengeById(challenge.id.toString(), ctrl.signal);
      setFullChallenge(data);
      setOpen(true); // open only after data is ready
    } catch (e) {
      console.error("Failed to load challenge", e);
      // Optionally show a toast/error UI here
    } finally {
      setLoading(false);
    }
  };

  const handleOpenChange = (next: boolean) => {
    if (!next) {
      // closing: abort any pending fetch & cleanup
      abortRef.current?.abort();
      setOpen(false);
      setFullChallenge(null);
    }
    // Note: we do NOT allow opening via onOpenChange;
    // opening is driven by openWithData() to ensure data-first.
  };

  return (
    <Dialog key={challenge.id} open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Card
          onClick={openWithData}
          className={
            challenge.hidden
              ? "m-4 w-[250px] h-[130px] cursor-pointer border-dashed border-2"
              : "m-4 w-[250px] h-[130px] cursor-pointer"
          }
        >
          <CardHeader className="px-4 flex flex-col justify-between h-full">
            <CardTitle>{challenge.title}</CardTitle>
            <CardDescription className=" h-5">
                {challengeSettings.showTags && showTags(challenge.tags)}
            </CardDescription>
          </CardHeader>
          <CardFooter className="pt-2 flex justify-between items-end mt-auto">
              {challenge.points && showPoints(challenge.points)}
              {typeof challenge.timeout === "number" && challenge.timeout > 0 && (
                <TimeoutBadge key={challenge.timeout + challenge.id} seconds={challenge.timeout} />
              )}
              {challenge.solved && challenge.first_blood && <Droplet size={24} className="text-red-500" />}
              {challenge.solved && !challenge.first_blood && <CircleCheck size={24} className="text-green-500" />}
          </CardFooter>
        </Card>
      </DialogTrigger>

      {open && fullChallenge && (
        <Suspense
          fallback={<DialogLoading />}
        >
          <LazyChallengeContent
            challenge={fullChallenge}
            challengeSettings={challengeSettings}
            auth={auth}
          />
        </Suspense>
      )}

      {/* Optional tiny overlay while we fetch BEFORE opening */}
      {loading && !open && (
        <div className="fixed inset-0 pointer-events-none">
          {/* or trigger a toast/spinner near the card */}
        </div>
      )}
    </Dialog>
  );
}
