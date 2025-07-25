import React, { useContext } from 'react';
import SettingContext from '@/context/SettingsProvider';
import { useNavigate } from 'react-router-dom';


export function CreateTeam() {
  const { settings } = useContext(SettingContext);
  const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
  return (
    <>
       <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
        Create a new Team
      </h2>
      { showQuotes && (
                <blockquote className="mt-6 border-l-2 pl-6 italic">
                    "Find a group of people who challenge and inspire you, spend a lot of time with them, and it will change your life."
                </blockquote>
            )}
      <p className="mb-4">
        To create a new team, please enter a unique team name.
      </p>
      <input
        type="text"
        placeholder="Enter team name"
        className="border border-gray-300 p-2 rounded-lg mb-4"
      />
      <button className="bg-blue-600 text-white rounded-lg p-2">
        Create Team
      </button>
    </>
  );
}
