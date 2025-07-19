import { Joystick, Globe,KeyRound, Bug, Wrench, ScanSearch, PackageOpen} from "lucide-react";
import React from "react";

export interface CategoryIconProps {
    size: number,
    category: string
}

export const CategoryIcon: React.FC<CategoryIconProps> = ({ size, category }) => {
    switch (category) {
        case 'Web':
            return <Globe size={size} />
        case 'Rev':
            return <Wrench size={size} />
        case 'Forensics':
            return <ScanSearch size={size} />
        case 'Crypto':
            return <KeyRound size={size} />
        case 'Pwn':
            return <Bug size={size} />
        case 'Misc':
            return <PackageOpen size={size} />
        default:
            return <Joystick size={size} />
    }
}