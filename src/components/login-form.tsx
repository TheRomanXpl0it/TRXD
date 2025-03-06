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


const formSchema = z.object({
    username: z.string().min(2, {
      message: "Username must be at least 2 characters.",
    }),
    password: z.string().min(2, {
        message: "Password should be strong",
    }),
})

export function LoginForm() {
    const { auth,setAuth } = useContext(AuthContext);


    const form = useForm<z.infer<typeof formSchema>>({
      resolver: zodResolver(formSchema),
      defaultValues: {
        username: "",
        password: "",
      },
    })

    function onSubmit(values: z.infer<typeof formSchema>) {
        // Do something with the form values.
        // ✅ This will be type-safe and validated.
        try {
          /* 
          const response = await api.post(LOGIN_URL,
            JSON.stringify(values),
            {
              headers: {
                "Content-Type": "application/json",
                withCredentials: true,
              },
            }
          );
           */
          const accessToken = "yes";
          const roles = ["admin"];
          const username = "admin";
          const password = "admin";
          if (values.username === "admin" && values.password === "admin") {
              setAuth({username,password,accessToken,roles});
          }
        }
        catch (error) {
          console.error(error);
        }
    }

    return (
        <>
          {auth.username ? <h1>already logged in</h1> : null}
          <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input placeholder="Username" {...field} />
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