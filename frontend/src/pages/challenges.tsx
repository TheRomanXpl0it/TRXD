import { Categories } from "@/components/categories"
import { CreateChallenge } from "@/components/createChallenge";
import { useState, useEffect, useContext } from "react";
import { getChallenges, getCategories } from "@/lib/backend-interaction";
import SettingContext from "@/context/SettingsProvider";
import AuthContext from "@/context/AuthProvider";

export function Challenges() {
    const [challenges, setChallenges] = useState([]);
    const [categories, setCategories] = useState([]);
    const { settings } = useContext(SettingContext);
    const { auth } = useContext(AuthContext);
    const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
    const canPost = auth && (auth.roles.includes('Admin') || auth.roles.includes('author'));

    useEffect(() => {
        async function fetchChallenges() {
            const challengesResult = await getChallenges();
            const challenges = JSON.parse(challengesResult);
            setChallenges(challenges);
        }
        async function fetchCategories() {
            const categoriesResult = await getCategories();
            const categories = JSON.parse(categoriesResult);
            setCategories(categories);
        }
        fetchCategories();
        fetchChallenges();
    }, []);

    return (
    <>
        <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
            Challenges
        </h2>
        { showQuotes && (
            <blockquote className="mt-6 border-l-2 pl-6 italic">
            "A man who loves to walk will walk more than a man who loves his destination"
            </blockquote>
        )}
        { canPost && (
            <div className="flex justify-end">
                <CreateChallenge/>
            </div>
        )}
        <Categories challenges={challenges} categories={categories} />
    </>
    )
}