import { createContext, useState } from "react";
import { ReactNode } from "react";

const SettingContext = createContext({
    settings: {
        General: [
            { title: 'Title', value: "TRXD", description: "The title of the CTF", type: String },
            { title: 'Description', value: "Welcome to TRXD", description: "The description of the CTF", type: String },
            { title: 'Start Time', value: "2022-01-01T00:00:00Z", description: "The start time of the CTF", type: Date },
            { title: 'End Time', value: "2022-01-01T00:00:00Z", description: "The end time of the CTF", type: Date },
            { title: 'Show Quotes', value: true, description: "Show motivational quotes in each page", type: Boolean },
            { title: 'Allow Writeups', value: true, description: "Allow users to submit writeups", type: Boolean },
            { title: 'Allow Registration', value: true, description: "Allow users to register for the CTF", type: Boolean },
            { title: 'Allow Team Play', value: true, description: "Allow users to join a team, if disabled the competition becomes individual", type: Boolean}
        ],
        Challenges: [
            { title: 'Visible', value: true, description: "Show challenges to other users", type: Boolean },
            { title: 'Show Points', value: true, description: "Show the points assigned to each challenge within the challenge card", type: Boolean },
            { title: 'Show Difficulty', value: true, description: "Show difficulty of the challenge within the challenge card", type: Boolean },
            { title: 'Show Tags', value: true, description: "Show tags of the challenge within the challenge card", type: Boolean },
            { title: 'Show Author', value: true, description: "Show the author of the challenge within the challenge card", type: Boolean },
            { title: 'Show Solves', value: true, description: "Show the number of users who have solved the challenge within the challenge card", type: Boolean },
            { title: 'Send Notifications', value: true, description: "Send notifications to users when a new challenge is added", type: Boolean },
            { title: 'Send First Blood Notifications', value: true, description: "Send notifications to users when they solve a challenge first", type: Boolean },
        ],
        Scores: [
            { title: 'Dynamic Scoring', value: false, description: "Enable dynamic scoring for challenges", type: Boolean },
            { title: 'Display Graph', value: false, description: "Display the graph of all the users score in the leaderboard page", type: Boolean },
            { title: 'Display Badges', value: false, description: "Display the badges of the users in the leaderboard page", type: Boolean },
            { title: 'Freeze Scoreboard', value: false, description: "Freeze the scoreboard from this moment onwards", type: Boolean },
        ]
    },
    setSettings: (_settings: any) => {}
});

interface SettingsProviderProps {
    children: ReactNode;
}

export const SettingsProvider = ({ children }: SettingsProviderProps) => {
    const [settings,setSettings] = useState({
        General: [
            { title: 'Title', value: "TRXD", description: "The title of the CTF", type: String },
            { title: 'Description', value: "Welcome to TRXD", description: "The description of the CTF", type: String },
            { title: 'Start Time', value: "2022-01-01T00:00:00Z", description: "The start time of the CTF", type: Date },
            { title: 'End Time', value: "2022-01-01T00:00:00Z", description: "The end time of the CTF", type: Date },
            { title: 'Show Quotes', value: true, description: "Show motivational quotes in each page", type: Boolean },
            { title: 'Allow Writeups', value: true, description: "Allow users to submit writeups", type: Boolean },
            { title: 'Allow Registration', value: true, description: "Allow users to register for the CTF", type: Boolean },
            { title: 'Allow Team Play', value: true, description: "Allow users to join a team, if disabled the competition becomes individual", type: Boolean}
        ],
        Challenges: [
            { title: 'Visible', value: true, description: "Show challenges to other users", type: Boolean },
            { title: 'Show Points', value: true, description: "Show the points assigned to each challenge within the challenge card", type: Boolean },
            { title: 'Show Difficulty', value: true, description: "Show difficulty of the challenge within the challenge card", type: Boolean },
            { title: 'Show Category', value: true, description: "Show the category of the challenge within the challenge card", type: Boolean },
            { title: 'Show Tags', value: true, description: "Show tags of the challenge within the challenge card", type: Boolean },
            { title: 'Show Author', value: true, description: "Show the author of the challenge within the challenge card", type: Boolean },
            { title: 'Show Solves', value: true, description: "Show the number of users who have solved the challenge within the challenge card", type: Boolean },
            { title: 'Send Notifications', value: true, description: "Send notifications to users when a new challenge is added", type: Boolean },
            { title: 'Send First Blood Notifications', value: true, description: "Send notifications to users when they solve a challenge first", type: Boolean },
        ],
        Scores: [
            { title: 'Dynamic Scoring', value: false, description: "Enable dynamic scoring for challenges", type: Boolean },
            { title: 'Display Graph', value: true, description: "Display the graph of all the users score in the leaderboard page", type: Boolean },
            { title: 'Display Badges', value: true, description: "Display the badges of the users in the leaderboard page", type: Boolean },
            { title: 'Freeze Scoreboard', value: false, description: "Freeze the scoreboard from this moment onwards", type: Boolean },
        ]
    });

    return (
        <SettingContext.Provider value={{settings,setSettings}}>
            {children}
        </SettingContext.Provider>
    )
};

export default SettingContext;