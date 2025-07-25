import React, { useContext } from 'react';
import SettingContext from '@/context/SettingsProvider';
import { useNavigate } from 'react-router-dom';

export function JoinTeam() {
  const { settings } = useContext(SettingContext);
  const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;

  return (
    <>
      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
        Join an existing Team
      </h2>
      { showQuotes && (
                <blockquote className="mt-6 border-l-2 pl-6 italic">
                    "Coming together is a beginning, staying together is progress, and working together is success."
                </blockquote>
            )}
      <p className="mt-7  ">
        To join a team, please enter the team code provided by your team leader.
      </p>
      <input
        type="text"
        placeholder="Enter team code"
        className="border border-gray-300 p-2 rounded-lg mb-4"
      />
      <button className="bg-blue-600 text-white rounded-lg p-2">
        Join Team
      </button>
    </>
  );
}
