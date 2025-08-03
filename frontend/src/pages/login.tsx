import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { LoginForm } from "@/components/LoginForm"
import { RegisterForm } from "@/components/RegisterForm"
import { Title } from "@/components/ui/title"


export function Login () {
    return (
        <>
        <Title />
        <div className="flex items-center justify-center">
            <Tabs defaultValue="Login" className="w-[400px]">
            <TabsList>
                <TabsTrigger value="Login">Login</TabsTrigger>
                <TabsTrigger value="Register">Register</TabsTrigger>
            </TabsList>
            <TabsContent value="Login">
                <Card>
                <CardHeader>
                    <CardTitle>Login</CardTitle>
                    <CardDescription>Welcome back, insert your credentials and start hacking.</CardDescription>
                </CardHeader>
                <CardContent>
                    <LoginForm/>
                </CardContent>
                </Card>
            </TabsContent>
            <TabsContent value="Register">
                <Card>
                <CardHeader>
                    <CardTitle>Register</CardTitle>
                    <CardDescription>Nice to see you, to get started please fill in the form.</CardDescription>
                </CardHeader>
                <CardContent>
                    <RegisterForm/>
                </CardContent>
                </Card>
            </TabsContent>
            </Tabs>
        </div>
        </>
    )
}