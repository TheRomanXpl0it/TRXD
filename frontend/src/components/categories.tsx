import { displayChallenges } from './Challenge';
import React, { useContext } from 'react';
import SettingContext from '@/context/SettingsProvider';
import AuthContext from '@/context/AuthProvider';
import { AuthProps } from '@/context/AuthProvider';
import { useChallenges } from '@/context/ChallengeProvider'; // ✅ use the custom hook
import { Challenge } from '@/context/ChallengeProvider';

function challengeByCategory(challenges: Challenge[], category: string) {
  return challenges.filter(
    (challenge) => challenge?.category?.includes(category)
  );
}


function displayCategory(
    category: string,
    challenges: Challenge[],
    settings: {
        title: string;
        value: boolean;
        description: string;
        type: BooleanConstructor;
    }[],
    auth: AuthProps
) {
    const filtered = challengeByCategory(challenges, category);
    if (filtered.length === 0) return null;
    if (!filtered || !filtered?.length){
        return null;
    }

    return (
        <React.Fragment>
            <h2 className="text-2xl font-semibold mt-8">{category}</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 m-1 justify-center">
                {displayChallenges(filtered, category, settings, auth).map((challenge, index) => (
                    <React.Fragment key={index}>{challenge}</React.Fragment>
                ))}
            </div>
        </React.Fragment>
    );
}

export const Categories: React.FC = () => {
    const { settings } = useContext(SettingContext);
    const { auth } = useContext(AuthContext);
    const { challenges = [], categories = [] } = useChallenges();

    const isVisible = settings.Challenges?.find((setting) => setting.title === 'Visible')?.value;
    const challengeSettings = settings.Challenges;

    return (
        <>
            {isVisible &&
                auth &&
                categories.map((category) => (
                    <React.Fragment key={category}>
                        {displayCategory(category, challenges, challengeSettings, auth)}
                    </React.Fragment>
                ))}
        </>
    );
};
