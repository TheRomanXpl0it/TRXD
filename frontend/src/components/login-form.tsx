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
      try {
          const data = await login(values); // contains accessToken, roles, etc.
          const { username, roles } = data;

          setAuth({
              username,
              roles,
          });


      } catch (error: any) {
          console.error("Login failed:", error.response?.data || error.message);
          form.setError("email", { message: "Invalid credentials" });
          form.setError("password", { message: "Invalid credentials" });
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