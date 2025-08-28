"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { toast } from "sonner";
import { useContext, useState } from "react";
import { submitFlag } from "@/lib/backend-interaction";
import { ChallengeContext, Challenge as ChallengeType } from "@/context/ChallengeProvider";
import { useInstance } from "@/context/InstanceProvider";

const formSchema = z.object({ flag: z.string().nonempty({ message: "Flag is required" }) });

export function FlagSubmit({ challenge }: { challenge: ChallengeType }) {
  const { isInstanced, isRunning } = useInstance();
  const challengeContext = useContext(ChallengeContext);
  const setChallenges = challengeContext?.setChallenges ?? ((_: any) => {});
  const [status, setStatus] = useState<"Idle" | "Correct" | "Wrong" | "Invalid" | "Repeated">("Idle");

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { flag: "" },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    try {
      const result = await submitFlag(challenge.id, data.flag);
      console.log(result);
      const s = result.data.status as typeof status;
      const firstblood = result.data.first_blood as boolean;
      setStatus(s);
      if (s === "Correct" && !firstblood) {
        toast.success("Flag submitted successfully!");
        setChallenges(prev => prev.map(c => (c.id === challenge.id ? { ...c, solved: true } : c)));
      }
      else if (s === "Correct" && firstblood) {
        toast.success("Flag submitted successfully! First blood!");
        setChallenges(prev => prev.map(c => (c.id === challenge.id ? { ...c, solved: true, firstBlood: true } : c)));
      } 
      else if (s === "Wrong") toast.error("Wrong flag.");
      else if (s === "Repeated") toast.error("Flag already submitted.");
      else if (s === "Invalid") toast.error("Invalid flag format.");
      else { toast.error("Unexpected response."); console.error(result); }
    } catch (e) {
      toast.error("Network error while submitting flag.");
      console.error(e);
    }
  };

  const disabled = isInstanced && !isRunning;

  return (
    <div className="w-full">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col space-y-2 w-full">
          <Label>Flag</Label>
          <div className="flex items-end space-x-2">
            <FormField
              control={form.control}
              name="flag"
              render={({ field }) => (
                <FormItem className="flex-1">
                  <FormControl>
                    <Input
                      placeholder="TRX{...}"
                      disabled={disabled}
                      className={cn(
                        status === "Correct" && "border-green-500",
                        status === "Wrong" && "border-red-500",
                        status === "Invalid" && "border-yellow-500"
                      )}
                      {...field}
                      onChange={(e) => { field.onChange(e); setStatus("Idle"); }}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit" className="self-start" disabled={disabled}>
              Submit
            </Button>
          </div>
          {disabled && (
            <p className="text-sm text-muted-foreground">
              Start the instance to submit the flag.
            </p>
          )}
        </form>
      </Form>
    </div>
  );
}
