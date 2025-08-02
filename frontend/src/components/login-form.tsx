"use client"

import { z } from "zod"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button } from "@/components/ui/button"
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { useContext } from "react"
import AuthContext from "@/context/AuthProvider"
import { login } from "@/lib/backend-interaction"
import { data } from "react-router-dom"
import { set } from "date-fns"




const formSchema = z.object({
    email: z.string().min(2, {
      message: "Email must be at least 2 characters.",
    }).email("Invalid email format"),
    password: z.string().min(2, {
        message: "Password should be strong",
    }),
})



export function LoginForm() {
    const { auth,setAuth } = useContext(AuthContext);


    const form = useForm<z.infer<typeof formSchema>>({
      resolver: zodResolver(formSchema),
      defaultValues: {
        email: "",
        password: "",
      },
    })

    async function onSubmit(values: z.infer<typeof formSchema>) {
          const response = await login(values); // contains accessToken, role, etc.
          switch (response.status) {
              case 200:
                  console.log("Login successful:", response.data);
                  const { username, role } = response.data; // Assuming response contains these fields
                  if (!username || !role) {
                      form.setError("root", {
                          type: "manual",
                          message: "Invalid response from server",
                      });
                      return;
                  }
                  setAuth({
                      username,
                      roles: [role], // Assuming role is a string, adjust if it's an array
                  });
                  break;
              case 400:
                  form.setError("root", {
                      type: "manual",
                      message: "Invalid input. Please check your data.",
                  });
                  return;
              case 401:
                  form.setError("root", {
                      type: "manual",
                      message: "Invalid credentials",
                  });
                  form.setError("password", { type: "manual", message: "Invalid credentials" });
                  return;
              default:
                  form.setError("root", {
                      type: "manual",
                      message: "An unexpected error occurred",
                  });
                  return;
          }
    }

    return (
        <>
          {auth.username ? <h1>already logged in</h1> : null}
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
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit">Login</Button>
            </form>
          </Form>
        </>
    );
}