import type { Solve } from "@/context/AuthProvider";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "./ui/carousel";
import { Card } from "./ui/card";
import { Droplet, FlagIcon, Crown, MedalIcon } from "lucide-react";
import type { TeamMember } from "@/context/AuthProvider";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useMemo } from "react";
import { UserCategoryRing, CategorySlice } from "./UserCategoryRing";
import type { Badge } from "@/context/AuthProvider";
import {
  Table,
  TableHeader,
  TableRow,
  TableCell,
  TableBody,
  TableHead,
} from "./ui/table";
import Autoplay from "embla-carousel-autoplay";

// ...

// Shared sizing
const OUTER_W = "w-80"; // carousel width (20rem)
const CARD_H = "h-56"; // 14rem height
const ITEM_BASIS = "basis-[20rem]"; // matches OUTER_W
const DELAY = 8000;

function BadgesCarousel({ badges }: { badges: Badge[] }) {
  return (
    <Carousel
      className={OUTER_W}
      opts={{
        align: "start",
        loop: true,
      }}
      plugins={[Autoplay({ delay: DELAY, stopOnInteraction: false })]}
    >
      <CarouselContent>
        <CarouselItem className={ITEM_BASIS}>
          <div className="p-2">
            <Card
              className={`w-full ${CARD_H} flex flex-col items-center justify-center p-4`}
            >
              <div className="flex items-center gap-2">
                <MedalIcon className="text-yellow-500" />
                <h1 className="text-xl font-bold tracking-tight">Badges</h1>
              </div>
              <div className="text-muted-foreground">{badges.length}  {badges.length == 1 ? "badge" : "badges"} achieved</div>
            </Card>
          </div>
        </CarouselItem>
        {badges.map((badge) => (
          <CarouselItem key={badge.id} className={ITEM_BASIS}>
            <div className="p-2">
              <Card
                className={`w-full ${CARD_H} flex flex-col items-center justify-center p-4`}
              >
                <img
                  src={badge.icon}
                  alt={badge.name}
                  className="w-16 h-16 mb-2"
                />
                <h2 className="text-lg font-semibold">{badge.name}</h2>
                <p className="text-sm text-center text-muted-foreground">
                  {badge.description}
                </p>
              </Card>
            </div>
          </CarouselItem>
        ))}
      </CarouselContent>
      <CarouselPrevious />
      <CarouselNext />
    </Carousel>
  );
}

function SolvesCarousel({
  teamSolves,
  totalSolves,
  ringCategories,
  categories,
}: {
  teamSolves: Solve[];
  ringCategories: CategorySlice[];
  totalSolves: number;
  categories: CategorySlice[];
}) {
  const firstBloods = teamSolves.filter((solve) => solve.first_blood);
  return (
    <Carousel
      className={OUTER_W}
      opts={{
        align: "start",
        loop: true,
      }}
      plugins={[Autoplay({ delay: DELAY, stopOnInteraction: false })]}
    >
      <CarouselContent>
        <CarouselItem className={ITEM_BASIS}>
          <div className="p-2">
            <Card
              className={`w-full ${CARD_H} flex flex-col items-center justify-center p-4 `}
            >
              <h1 className="text-2xl font-bold">Category Progress</h1>
              <UserCategoryRing
                totalSolves={ringCategories.reduce(
                  (s, c) => s + (c.solved ?? 0),
                  0,
                )}
                categories={categories}
                teamSolves={teamSolves}
                size={100}
                strokeWidth={16}
              />
            </Card>
          </div>
        </CarouselItem>
        {firstBloods.length >= 0 && (
          <CarouselItem>
            <div className="p-2">
              <Card
                className={`w-full ${CARD_H} flex flex-col items-center justify-center p-4`}
              >
                <div className="flex items-center gap-2">
                  <Droplet className="text-red-500" />
                  <h1 className="text-xl font-bold tracking-tight">
                    First Bloods
                  </h1>
                </div>
                <div className="text-muted-foreground">
                  {firstBloods.length} First Bloods
                </div>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Challenge</TableHead>
                      <TableHead>Category</TableHead>
                      <TableHead>Points</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {firstBloods
                      .sort((a, b) => b.points - a.points)
                      .slice(0, 5)
                      .map((solve) => (
                        <TableRow key={solve.id}>
                          <TableCell className="font-medium">
                            {solve.name}
                          </TableCell>
                          <TableCell>{solve.category}</TableCell>
                          <TableCell>{solve.points}</TableCell>
                        </TableRow>
                      ))}
                  </TableBody>
                </Table>
              </Card>
            </div>
          </CarouselItem>
        )}
        <CarouselItem className={ITEM_BASIS}>
          <div className="p-2">
            <Card
              className={`w-full ${CARD_H} flex flex-col items-center justify-center p-4`}
            >
              <div className="flex items-center gap-2">
                <FlagIcon />
                <h1 className="text-xl font-bold tracking-tight">
                  Team Solves
                </h1>
              </div>
              <div className="text-muted-foreground">
                {totalSolves} Total Solves
              </div>
            </Card>
          </div>
        </CarouselItem>
        <CarouselItem className={ITEM_BASIS}>
          <div className="p-2">
            <Card
              className={`w-full ${CARD_H} flex flex-col items-center justify-center p-4`}
            >
              <div className="flex items-center gap-2">
                <h1 className="text-xl font-bold tracking-tight">
                  Most valuable flags
                </h1>
              </div>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Challenge</TableHead>
                    <TableHead>Category</TableHead>
                    <TableHead>Points</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {teamSolves
                    .sort((a, b) => b.points - a.points)
                    .slice(0, 5)
                    .map((solve) => (
                      <TableRow key={solve.id}>
                        <TableCell className="font-medium">
                          {solve.name}
                        </TableCell>
                        <TableCell>{solve.category}</TableCell>
                        <TableCell>{solve.points}</TableCell>
                      </TableRow>
                    ))}
                </TableBody>
              </Table>
            </Card>
          </div>
        </CarouselItem>
      </CarouselContent>
      <CarouselPrevious />
      <CarouselNext />
    </Carousel>
  );
}

function MemberStatsCarousel({ members }: { members: TeamMember[] }) {
  const list = members ?? [];
  const mvp = useMemo(
    () =>
      list.length ? list.reduce((a, b) => (a.score > b.score ? a : b)) : null,
    [list],
  );
  const top5 = useMemo(
    () => [...list].sort((a, b) => b.score - a.score).slice(0, 5),
    [list],
  );
  return (
    <Carousel
      className={OUTER_W}
      opts={{
        align: "start",
        loop: true,
      }}
      plugins={[Autoplay({ delay: DELAY, stopOnInteraction: false })]}
    >
      <CarouselContent>
        <CarouselItem className={ITEM_BASIS}>
          <div className="p-2">
            <Card
              className={`w-full ${CARD_H} flex flex-col items-center justify-center p-4`}
            >
              <div className="flex items-center gap-2">
                <Crown className="text-yellow-500" />
                <h1 className="text-xl font-bold tracking-tight">MVP</h1>
              </div>
              <Avatar className="h-16 w-16 mt-2">
                <AvatarImage src={mvp?.profilePicture} />
                <AvatarFallback>{mvp?.name?.charAt(0) ?? "?"}</AvatarFallback>
              </Avatar>
              <div className="mt-2 truncate max-w-[14rem]">{mvp?.name}</div>
              <div className="text-sm text-muted-foreground">
                Score: {mvp?.score}
              </div>
            </Card>
          </div>
        </CarouselItem>

        <CarouselItem className={ITEM_BASIS}>
          <div className="p-2">
            <Card className={`w-full ${CARD_H} p-4 flex flex-col`}>
              <div className="text-center font-bold mb-2">Top 5 Members</div>
              <div className="space-y-2">
                {top5.map((m) => (
                  <div key={m.id} className="flex items-center justify-between">
                    <div className="flex items-center min-w-0">
                      <Avatar className="h-8 w-8">
                        <AvatarImage src={m.profilePicture} />
                        <AvatarFallback>
                          {m.name?.charAt(0) ?? "?"}
                        </AvatarFallback>
                      </Avatar>
                      <span className="ml-2 truncate max-w-[10rem]">
                        {m.name}
                      </span>
                    </div>
                    <span className="ml-2">{m.score}</span>
                  </div>
                ))}
              </div>
            </Card>
          </div>
        </CarouselItem>
      </CarouselContent>
      <CarouselPrevious />
      <CarouselNext />
    </Carousel>
  );
}

export function Stats({
  teamSolves,
  members,
  totalSolves,
  ringCategories,
  categories,
  badges,
}: {
  teamSolves: Solve[];
  members: TeamMember[];
  totalSolves: number;
  ringCategories: CategorySlice[];
  categories: CategorySlice[];
  badges: Badge[];
}) {
  return (
    <div className="mx-auto flex flex-row justify-center gap-7 mt-3">
      {badges.length > 0 && (
        <div className="flex-none w-100">
          <BadgesCarousel badges={badges} />
        </div>
      )}

      {totalSolves > 0 && (
        <div className="flex-none w-100">
          <SolvesCarousel
            totalSolves={totalSolves}
            teamSolves={teamSolves}
            ringCategories={ringCategories}
            categories={categories}
          />
        </div>
      )}

      {members.length > 0 && (
        <div className="flex-none w-100">
          <MemberStatsCarousel members={members} />
        </div>
      )}
    </div>
  );
}
