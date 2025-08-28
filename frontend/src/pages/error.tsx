export function ErrorPage() {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen">
            <div style={{ display: 'flex', justifyContent: 'center', marginTop: '24px' }}>
            <img
            src="https://raw.githubusercontent.com/TheRomanXpl0it/logo/main/TRX_smooth.svg"
            alt="TRX Logo"
            className="w-[35%] h-auto dark:invert"
            style={{ filter: 'invert(0)'}}
            />
        </div>
            <h1 className="text-4xl font-bold">Something went wrong</h1>
            <p className="mt-4">We're sorry, but an unexpected error has occurred.</p>
            <p className="mt-4">Please refresh the page or try again later.</p>
        </div>
    )
}