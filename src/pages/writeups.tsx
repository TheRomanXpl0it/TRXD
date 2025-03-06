import { useContext } from "react"
import SettingContext from "@/context/SettingsProvider"

export function Writeups () {
    const { settings } = useContext(SettingContext);
    const allowWriteups = settings.General?.find((setting) => setting.title === 'Allow Writeups')?.value;
    const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
    return (
    <>
        <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
            Writeups
        </h2>
        { showQuotes && (
        <blockquote className="mt-6 border-l-2 pl-6 italic">
        "The more that you read, the more things you will know. The more that you learn, the more places you'll go"
        </blockquote>
        )}
        { allowWriteups ? 
            <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 m-1 justify-center'>
                {/* Writeups go here */}
            </div>
        :
        <p className="text-xl text-muted-foreground mt-4">
            Writeups are disabled
        </p>
        }
    </>
    )
}