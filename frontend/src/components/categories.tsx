import { displayChallenges } from './Challenge';
import React, { useContext } from 'react';
import SettingContext from '@/context/SettingsProvider';
import AuthContext from '@/context/AuthProvider';
import { AuthProps } from '@/context/AuthProvider';
import { useChallenges } from '@/context/ChallengeProvider';
import { Challenge } from '@/context/ChallengeProvider';

function challengeByCategory(
  challenges: Challenge[],
  category: string,
  selectedTags?: string[]
) {
  return challenges.filter((challenge) =>
    challenge?.category?.includes(category) &&
    (!selectedTags || selectedTags.length === 0 || challenge.tags?.some((tag) => selectedTags.includes(tag)))
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
  auth: AuthProps,
  selectedTags?: string[]
) {
  const filtered = challengeByCategory(challenges, category, selectedTags);
  if (!filtered || filtered.length === 0) return null;

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

interface CategoriesProps {
  selectedCategories?: string[];
  selectedTags?: string[];
}

export const Categories: React.FC<CategoriesProps> = ({
  selectedCategories = [],
  selectedTags = [],
}) => {
  const { settings } = useContext(SettingContext);
  const { auth } = useContext(AuthContext);
  const { challenges = [], categories = [] } = useChallenges();

  const isVisible = settings.Challenges?.find((setting) => setting.title === 'Visible')?.value;
  const challengeSettings = settings.Challenges;

  // if no filters are selected, show all categories
  const filteredCategories =
    selectedCategories.length > 0
      ? categories.filter((category) => selectedCategories.includes(category))
      : categories;

  return (
    <>
      {isVisible &&
        auth &&
        filteredCategories.map((category) => (
          <React.Fragment key={category}>
            {displayCategory(category, challenges, challengeSettings, auth, selectedTags)}
          </React.Fragment>
        ))}
    </>
  );
};
