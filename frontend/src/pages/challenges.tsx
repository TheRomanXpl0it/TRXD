import { CreateChallenge } from "@/components/CreateChallenge";
import { lazy, Suspense } from "react";
import { useContext } from "react";
import SettingContext from "@/context/SettingsProvider";
import AuthContext from "@/context/AuthProvider";
import Loading from "@/components/Loading";

const Categories = lazy(() => import("@/components/Categories").then(module => ({ default: module.Categories })));
const ChallengeProvider = lazy(() => import("@/context/ChallengeProvider").then(module => ({ default: module.ChallengeProvider })));

export function Challenges() {
    const { settings } = useContext(SettingContext);
    const { auth } = useContext(AuthContext);
    const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
    const canPost = auth && (auth.roles.includes('Admin') || auth.roles.includes('author'));

    return (
    <>
        <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
            Challenges
        </h2>
        { showQuotes && (
            <blockquote className="mt-6 border-l-2 pl-6 italic">
            "A man who loves to walk will walk more than a man who loves his destination"
            </blockquote>
        )}
        { canPost && (
            <div className="flex justify-end">
                <CreateChallenge/>
            </div>
        )}
        <Suspense fallback={<Loading />}>
            <ChallengeProvider>
                <Categories />
            </ChallengeProvider>
        </Suspense>
    </>
    )
}