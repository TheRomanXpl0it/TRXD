import { useContext } from "react"
import SettingContext from "@/context/SettingsProvider"

export function Team(){
    const {settings} = useContext(SettingContext);
    const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
    const allowTeamPlay = settings.General?.find((setting) => setting.title === 'Allow Team Play')?.value;
    return(
        <>
            <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
                Join or create a team
            </h2>
            {showQuotes && (
                <blockquote className="mt-6 border-l-2 pl-6 italic">
                    "None of us is as smart as all of us."
                </blockquote>
            )}
            {allowTeamPlay && (
                <h1> Teams go here</h1>
            )}
        </>
    )
}