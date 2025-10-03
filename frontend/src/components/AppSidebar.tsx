import { Flag, Trophy, User, Settings, LogOut, LogIn, BookText, UserPen, ChevronUp, Home, LucideProps, ShieldHalf} from "lucide-react"
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarMenuButton, SidebarMenuItem, SidebarMenu } from "@/components/ui/sidebar"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Link } from "react-router-dom"
import { useContext } from "react"
import { AuthContext } from "@/context/AuthProvider"
import SettingContext from "@/context/SettingsProvider"


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
    const { auth, logout } = useContext(AuthContext);
    const { settings } = useContext(SettingContext);

    const isLoggedIn = auth && auth.name !== null && auth.name !== "";
    const isAdmin = auth && auth.role === 'Admin';

    const allowWriteups = settings.General.find((setting) => setting.title === 'Allow Writeups')?.value;
    const allowTeamPlay = settings.General.find((setting) => setting.title === 'Allow Team Play')?.value;

    const loggedInItems = [
        {
            title: 'Account',
            icon: UserPen,
            url: '/account',
        }
    ];

    const items = [
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
    
    const adminItems = [
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
                <User/> {auth.name}
                <ChevronUp className="ml-auto" />
                </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
                side="top"
                align="start"
                className="w-full"
            >
                {isAdmin && displayAdminItems(adminItems)}
                {loggedInItems.map((item, index) => (
                <DropdownMenuItem asChild key={index}>
                    <Link to={item.url}>
                    <item.icon />
                    <span>{item.title}</span>
                    </Link>
                </DropdownMenuItem>
                ))}

                <DropdownMenuItem onSelect={logout}>
                <LogOut />
                <span>Logout</span>
                </DropdownMenuItem>

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