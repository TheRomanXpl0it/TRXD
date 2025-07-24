import { Flag, Trophy, User, Settings, LogOut, LogIn, BookText, UserPen, ChevronUp, Home, LucideProps, ShieldHalf} from "lucide-react"
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarMenuButton, SidebarMenuItem, SidebarMenu } from "@/components/ui/sidebar"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Link } from "react-router-dom"
import { useContext } from "react"
import AuthContext from "@/context/AuthProvider"
import SettingContext from "@/context/SettingsProvider"



function displayLoggedInItems(
    loggedInItems: {
        title: string;
        icon: React.ForwardRefExoticComponent<Omit<LucideProps, "ref"> & React.RefAttributes<SVGSVGElement>>;
        url: string;
    }[]
){
    return loggedInItems.map((item, index) => (
        <div key={index}><DropdownMenuItem asChild>
                <Link to = {item.url}>
                    <item.icon/>
                    <span>{item.title}</span>
                </Link>
            </DropdownMenuItem>
        </div>
    ));
}

function displayAdminItems(
    adminItems: {
        title: string;
        icon: React.ForwardRefExoticComponent<Omit<LucideProps, "ref"> & React.RefAttributes<SVGSVGElement>>;
        url: string;
    }[]
){
    return adminItems.map((item, index) => (
        <DropdownMenuItem asChild key={index}>
            <Link to = {item.url}>
                <item.icon/>
                <span>{item.title}</span>
            </Link>
        </DropdownMenuItem>
    ));
}

export function AppSidebar() {
    const { auth } = useContext(AuthContext);
    const { settings } = useContext(SettingContext);
    
    const isLoggedIn = auth.accessToken !== undefined;
    const isAdmin = auth.roles.includes('admin');

    const allowWriteups = settings.General.find((setting) => setting.title === 'Allow Writeups')?.value;
    const allowTeamPlay = settings.General.find((setting) => setting.title === 'Allow Team Play')?.value;

    let loggedInItems = [
        {
            title: 'Account',
            icon: UserPen,
            url: '/account',
        },
        {
            title: 'Logout',
            icon: LogOut,
            url: '/logout',
        },
    ];

    let items = [
        {
            title: 'Home',
            icon: Home,
            url: '/',
        },
        {
            title: 'Challenges',
            icon: Flag,
            url: '/challenges',
        },
        {
            title: 'Leaderboard',
            icon: Trophy,
            url: '/leaderboard',
        },
    ];
    
    let adminItems = [
        {
            title: 'Settings',
            icon: Settings,
            url: '/settings',
        },
    ];

    if (allowWriteups) {
        items.push({
            title: 'Writeups',
            icon: BookText,
            url: '/writeups',
        });
    }

    if (allowTeamPlay) {
        items.push({
            title: 'Team',
            icon: ShieldHalf,
            url: '/team',
        })
    }

  
    const footerContent = isLoggedIn ? 
    (
    <SidebarMenu>
        <SidebarMenuItem>
            <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <SidebarMenuButton>
                <User/> {auth.username}
                <ChevronUp className="ml-auto" />
                </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
                side="top"
                align="start"
                className="w-full"
            >
                {isAdmin ? displayAdminItems(adminItems) : null}
                {displayLoggedInItems(loggedInItems)}
            </DropdownMenuContent>
            </DropdownMenu>
        </SidebarMenuItem>
    </SidebarMenu>
  ) 
  : 
  (
    <SidebarMenu>
        <SidebarMenuItem>
            <SidebarMenuButton asChild>
                <Link to="/login">
                    <LogIn/>
                    <span>Login</span>
                </Link>
            </SidebarMenuButton>
        </SidebarMenuItem>
    </SidebarMenu>
  );


  return (
    <Sidebar>
      <SidebarContent>
        <SidebarGroup>
            <SidebarGroupLabel>TRXD</SidebarGroupLabel>    
            <SidebarGroupContent>
                {items.map((item, index) => (
                        <SidebarMenuButton asChild key={index}>
                            <Link to = {item.url}>
                                <item.icon/>
                                <span>{item.title}</span>
                            </Link>
                        </SidebarMenuButton>
                ))}
            </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    <SidebarFooter>
        <SidebarGroup>
            {footerContent}
        </SidebarGroup>
    </SidebarFooter>
    </Sidebar>
  );
}