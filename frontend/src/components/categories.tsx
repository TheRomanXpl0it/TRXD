import { displayChallenges, ChallengeProps } from './Challenge'
import React, { useContext } from 'react';
import SettingContext from '@/context/SettingsProvider';
import AuthContext from '@/context/AuthProvider';
import { AuthProps } from '@/context/AuthProvider';
import { useState, useEffect } from 'react';
import { getChallenges, getCategories } from '@/lib/backend-interaction';


function challengeByCategory(challengeProps: ChallengeProps[], category: string) {
    return challengeProps.filter((challengeProp) => challengeProp.challenge.category.includes(category));
}

function displayCategory(
    category: string,
    challenges: ChallengeProps[],
    settings:{
        title: string;
        value: boolean;
        description: string;
        type: BooleanConstructor;
    }[],
    auth: AuthProps
    ) {
    challenges = challengeByCategory(challenges, category);

    if (challenges.length === 0) return null;
    return (
        <React.Fragment>
            <h2 className="text-2xl font-semibold mt-8">{category}</h2>
            <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 m-1 justify-center'>
                { displayChallenges(challenges, category, settings, auth).map((challenge, index) => (
                    <React.Fragment key={index}>
                        {challenge}
                    </React.Fragment>
                )) }
            </div>
        </React.Fragment>
    )
}

export const Categories: React.FC = () => {
    const { settings } = useContext(SettingContext);
    const { auth } = useContext(AuthContext);
    const [challenges, setChallenges] = useState([]);
    const [categories, setCategories] = useState([]);

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
    const isVisible = settings.Challenges?.find((setting) => setting.title === 'Visible')?.value;
    const challengeSettings = settings.Challenges;

    return (
        <>
            { isVisible && auth && categories.map((category) => (
                    <React.Fragment key={category}>
                        { displayCategory(category, challenges, challengeSettings, auth) }
                    </React.Fragment>
                
            ))}
        </>
    );
};