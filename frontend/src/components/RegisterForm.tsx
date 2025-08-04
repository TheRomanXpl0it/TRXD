"use client"

import { z } from "zod"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button } from "@/components/ui/button"
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { register,getSessionInfo } from "@/lib/backend-interaction"
import { useContext } from "react"
import { useNavigate } from "react-router-dom";
import { AuthContext, isAuthProps } from "@/context/AuthProvider"


const formSchema = z.object({
    username: z.string().min(2, {
      message: "Username must be at least 2 characters.",
    }),
    email: z.string().min(2, {
      message: "Email must be at least 2 characters.",
    }).email("Invalid email format"),
    password: z.string().min(2, {
        message: "Password must be at 2 characters.",
    }),
    confirmPassword: z.string().min(2, {
        message: "Username must be at least 2 characters.",
    }),
})

export function RegisterForm() {
  const { auth,setAuth } = useContext(AuthContext);  
  const navigate = useNavigate();

    const form = useForm<z.infer<typeof formSchema>>({
      resolver: zodResolver(formSchema),
      defaultValues: {
        username: "",
        password: "",
        confirmPassword: "",
      },
    })
    async function onSubmit(values: z.infer<typeof formSchema>) {
        const response = await register(values);
        switch (response) {
            case 200:
                  break;
            case 400:
                form.setError("root", {
                    type: "manual",
                    message: "Invalid input. Please check your data.",
                });
                break;
            case 409:
                form.setError("email", {
                    type: "manual",
                    message: "Email is already taken.",
                });
                break;
            default:
                form.setError("root", {
                    type: "manual",
                    message: "Unexpected error occurred.",
                });
                break;
        }
        const sessionInfo = await getSessionInfo();
        if (isAuthProps(sessionInfo)) {
            setAuth(sessionInfo);
            navigate("/challenges");
        }
        return;
    }

    return (
      <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input placeholder="user@example.it" {...field} />
                  </FormControl>
                  <FormDescription>The email address user for registration and logging in</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl>
                <Input placeholder="Username" {...field} />
              </FormControl>
              <FormDescription>Your publicly displayed username</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input placeholder="Password" type="password" {...field} />
              </FormControl>
              <FormDescription>Remember to choose a strong password.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="confirmPassword"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Confirm Password</FormLabel>
              <FormControl>
                <Input placeholder="Confirm Password" type="password" {...field} />
              </FormControl>
              <FormDescription>Make sure your passwords match.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Register</Button>
      </form>
    </Form>
    );
}