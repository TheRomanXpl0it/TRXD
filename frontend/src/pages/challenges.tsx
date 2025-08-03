import { CreateChallenge } from "@/components/CreateChallenge";
import { lazy, Suspense, useState } from "react";
import { useContext } from "react";
import { Filter } from "lucide-react";
import SettingContext from "@/context/SettingsProvider";
import AuthContext from "@/context/AuthProvider";
import Loading from "@/components/Loading";
import {
  MultiSelector,
  MultiSelectorTrigger,
  MultiSelectorInput,
  MultiSelectorContent,
  MultiSelectorList,
  MultiSelectorItem,
} from "@/components/ui/multi-select";
import { Separator } from "@/components/ui/separator";
import { useChallenges } from "@/context/ChallengeProvider";

const Categories = lazy(() =>
  import("@/components/Categories").then((module) => ({
    default: module.Categories,
  }))
);
const ChallengeProvider = lazy(() =>
  import("@/context/ChallengeProvider").then((module) => ({
    default: module.ChallengeProvider,
  }))
);

// 🔁 Combined Filter Component
function CombinedFilter({
  selectedFilters,
  setSelectedFilters,
}: {
  selectedFilters: string[];
  setSelectedFilters: React.Dispatch<React.SetStateAction<string[]>>;
}) {
  const { challenges = [], categories = [] } = useChallenges();

  const availableTags = Array.from(new Set(challenges.flatMap((c) => c.tags || [])));

return (
    <div className="flex items-center">
        <MultiSelector className="w-60" values={selectedFilters} onValuesChange={setSelectedFilters} loop>
            <MultiSelectorTrigger className="w-full flex items-center justify-between">
                <div className="flex w-full items-center justify-between">
                    <MultiSelectorInput placeholder="Filter by category/tag" />
                    <Filter className="h-4 w-4 text-muted-foreground" />
                </div>
            </MultiSelectorTrigger>
            <MultiSelectorContent>
                <MultiSelectorList>
                    <div className="px-2 py-1 text-sm font-semibold text-muted-foreground">Categories</div>
                    {categories.map((cat) => (
                        <MultiSelectorItem key={`cat:${cat}`} value={`cat:${cat}`}>
                            {cat}
                        </MultiSelectorItem>
                    ))}

                    <Separator className="my-2" />

                    <div className="px-2 py-1 text-sm font-semibold text-muted-foreground">Tags</div>
                    {availableTags.map((tag) => (
                        <MultiSelectorItem key={`tag:${tag}`} value={`tag:${tag}`}>
                            {tag}
                        </MultiSelectorItem>
                    ))}
                </MultiSelectorList>
            </MultiSelectorContent>
        </MultiSelector>
    </div>
);
}

// 🔁 Main Page Component
export function Challenges() {
  const { settings } = useContext(SettingContext);
  const { auth } = useContext(AuthContext);

  const [selectedFilters, setSelectedFilters] = useState<string[]>([]);

  const selectedCategories = selectedFilters
    .filter((v) => v.startsWith("cat:"))
    .map((v) => v.replace("cat:", ""));
  const selectedTags = selectedFilters
    .filter((v) => v.startsWith("tag:"))
    .map((v) => v.replace("tag:", ""));

  const showQuotes = settings.General?.find(
    (setting) => setting.title === "Show Quotes"
  )?.value;

  const canPost =
    auth && (auth.roles.includes("Admin") || auth.roles.includes("author"));

  return (
    <>
      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
        Challenges
      </h2>
      {showQuotes && (
        <blockquote className="mt-6 border-l-2 pl-6 italic">
          "A man who loves to walk will walk more than a man who loves his destination"
        </blockquote>
      )}
      <Suspense fallback={<Loading />}>
        <ChallengeProvider>
          <div className="flex justify-end items-stretch gap-2">
            <CombinedFilter
                selectedFilters={selectedFilters}
                setSelectedFilters={setSelectedFilters}
            />
            {canPost && (
                <div className="h-full flex items-center">
                <CreateChallenge />
                </div>
            )}
           </div>



          {/* ✅ Filtered Categories View */}
          <Categories
            selectedCategories={selectedCategories}
            selectedTags={selectedTags}
          />
        </ChallengeProvider>
      </Suspense>
    </>
  );
}