import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/AppSidebar"
import { Outlet } from 'react-router-dom'
import { ThemeProvider } from "@/components/ThemeProvider"
import { ModeToggle } from "./components/ModeToggle";
import { Toaster } from "@/components/ui/sonner"


export function Layout(){
    return (
        <>
            <SidebarProvider>
                <AppSidebar/>
                <SidebarTrigger />
                <main className="m-4 w-full">
                    <Outlet />
                </main>
            </SidebarProvider>
            <div className="fixed bottom-0 right-0 p-4 text-sm text-gray-500">
                <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
                    <ModeToggle />
                </ThemeProvider>
            </div>
            <Toaster />
        </>
    );
}