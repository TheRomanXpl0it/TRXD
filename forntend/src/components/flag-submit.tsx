"use client"
 
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Container } from "lucide-react"
import { ChallengeProps } from "./challenge"
import { useState, useEffect } from "react"

const formSchema = z.object({
    flag: z.string().nonempty({
        message: "Flag is required",
    }),
});





export function FlagSubmit(challengeProp: ChallengeProps) {
  const isInstanced: boolean = challengeProp.challenge.instanced ?? false;
  const timeout: Date = challengeProp.challenge.timeout ? new Date(challengeProp.challenge.timeout) : new Date();

  const [remainingTime, setRemainingTime] = useState(
    Math.max(0, Math.floor((timeout.getTime() - Date.now()) / 1000))
  );

  useEffect(() => {
    const interval = setInterval(() => {
      setRemainingTime(Math.max(0, Math.floor((timeout.getTime() - Date.now()) / 1000)));
    }, 1000);

    return () => clearInterval(interval); // Cleanup on unmount
  }, [timeout]);

  function showContainer() {
    return (
      <Button className="bg-blue-500 w-full text-white">
        <Container /> Start Instance
      </Button>
    );
  }

  function showTimer() {
    return (
      <div className="flex items-center space-x-2 rounded-lg bg-blue-500 text-white p-2 justify-between">
        <span className="flex items-center space-x-2">
          <Container size={24} />
          <span>Instance Running</span>
        </span>
        <span>
          Time Remaining: {Math.floor(remainingTime / 60)}:
          {("0" + (remainingTime % 60)).slice(-2)} seconds
        </span>
      </div>
    );
  }

  function showFlagSubmit() {
    return (
      <div className="w-full space-y-4">
        {isInstanced && remainingTime > 0 ? <div className="w-full">{showTimer()}</div> : null}
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
                        <Input placeholder="TRX{...}" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button type="submit" className="self-start">Submit</Button>
              </div>
            </form>
          </Form>
        </div>
      </div>
    );
  }

  function showDockerControls() {
    return remainingTime > 0 ? showFlagSubmit() : showContainer();
  }

  const form = useForm({
    resolver: zodResolver(formSchema),
  });

  const onSubmit = (data: any) => {
    console.log(data);
  };

  return <>{isInstanced ? showDockerControls() : showFlagSubmit()}</>;
}
