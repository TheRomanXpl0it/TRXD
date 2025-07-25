import  { useContext } from 'react';
import SettingContext from '@/context/SettingsProvider';

export function Home () {
    const { settings } = useContext(SettingContext);
    const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
    return (
        <>
        <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
            Home
        </h2>
        { showQuotes && (
            <blockquote className="mt-6 border-l-2 pl-6 italic">
            "There's no place like ~"
            </blockquote>
        )}

        <div style={{ display: 'flex', justifyContent: 'center', marginTop: '24px' }}>
            <img
            src="https://raw.githubusercontent.com/TheRomanXpl0it/logo/main/TRX_smooth.svg"
            alt="TRX Logo"
            className="w-[35%] h-auto dark:invert"
            style={{ filter: 'invert(0)' }}
            />
        </div>
        <div className='text-center font-semibold text-7xl'>
            TRXD
        </div>
        </>
    )
}