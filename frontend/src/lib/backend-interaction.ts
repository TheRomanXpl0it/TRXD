import { api } from "@/api/axios";
import type { Team, AuthProps, User } from "@/context/AuthProvider";
import axios from "axios";


export async function getChallengeData(){
    // Simulate file attachments as objects with name and url
    type Attachment = {
        name: string;
        url: string;
    };

    let challenges = [];
    let categories = [];

    try {
      challenges = (await api.get("/challenges")).data;
      console.log("Fetched challenges:", challenges);
    }
    catch (error) {
        if (axios.isAxiosError(error) && error.response) {
            console.error("Error fetching challenges:", error.response.data);
            return JSON.stringify([]);
        }
        console.error("Unexpected error:", error);
        return JSON.stringify([]);
    }
    challenges = challenges.map((challenge: any) => {
        return {
                id: challenge.id,
                title: challenge.name,
                description: challenge.description,
                solves: challenge.solves,
                points: challenge.points,
                category: challenge.category,
                remote: challenge.host,
                solved: challenge.solved,
                tags: challenge.tags,
                difficulty: challenge.difficulty,
                attachments: challenge.attachments.map((attachment: any) => ({
                    name: attachment.name,
                    url: attachment.url,
                })) as Attachment[],
                authors: challenge.authors,
                hidden: challenge.hidden,
                flags: challenge.flags,
                instanced: challenge.instance,
                timeout: challenge.timeout ? new Date(challenge.timeout) : undefined,
        };
    });
    categories = challenges.map((challenge: any) => challenge.category);
    categories = [...new Set(categories)]; // Remove duplicates
    return JSON.stringify({ challenges, categories });

    const mockChallenges = [
        {
            challenge : {
                id : 1,
                title: 'Challenge 1',
                description: 'This is the first challenge',
                flag: 'flag{first_challenge}',
                points: 100,
                solves: 0,
                category: 'Web',
                tags: ['Web', 'Easy'],
                difficulty: 'Easy',
                remote: 'https://example.com',
                solved: false,
                attachments: [
                    { name: "executable", url: "/files/challenge1/executable" },
                    { name: "libc", url: "/files/challenge1/libc" }
                ] as Attachment[],
                authors: ["author1"],
                hidden: true,
                instanced: false,
                timeout: undefined,
            },
        }
    ];
    return JSON.stringify(mockChallenges);
} 

export async function login({
  email,
  password,
}: {
  email: string;
  password: string;
}): Promise< number > {
  try {
    const response = await api.post(
      "/login",
      { email, password },
      { withCredentials: true }
    );
    return response.status;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      return error.response.status;
    }

    console.error("Unexpected login error:", error);
    return 500;
  }

}

export async function register({
  username,
  email,
  password,
}: {
  username: string;
  email: string;
  password: string;
}): Promise< number> {
  try {
    const response = await api.post(
      "/register",
      { username, email, password },
      { withCredentials: true }
    );
    return response.status ;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      // Return the server's status and any data it returned
      return error.response.status
    }

    // Unknown error (network error or non-Axios error)
    console.error("Unexpected error during registration:", error);
    return 500;
  }
}

export async function fetchTeamData(teamId: number): Promise<Team> {
    // Simulate fetching team data
    let teamData = {
        id: teamId,
        name: "",
        members: [""],
        teamlogo: "/teamLogo.png",
        score: 0,
        rank: -1,
    };
    try{
        const response = await api.get(`/team/${teamId}`);
        teamData = response.data;
    } catch (error) {
        console.error("Error fetching team data:", error);
    }
    return teamData;
}

export async function getUsersTeamData(): Promise<Team | null> {
    // Simulate fetching team data
    let teamData = null;
    try {
        const response = await api.get("/team");
        if (response.status === 200) {
            teamData = { id: response.data.id, name: response.data.name, logo: response.data.logo, members: response.data.members } as Team;
        } else {
            console.error("Error fetching user's team data:", response.statusText);
        }
    } catch (error) {
        console.error("Error fetching user's team data:", error);
    }
    return teamData;
}

export async function leaveTeam() {
    // Simulate leaving a team
    console.log("Leaving team...");
    return { success: true, message: "You have left the team." };
}

export async function getSessionInfo(): Promise<AuthProps | number> {
  try {
    const response = await api.get("/info"); // or your auth check endpoint
    const data = response.data;
    return {
      id: data.id,
      username: data.username,
      role: data.role,
      teamId: data.team_id,
    };
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      return error.response.status;
    }
    return 500;
  }
}

export async function submitFlag(
  challengeId: number,
  flag: string
): Promise<{ status: number; data?: any }> {
  try {
    const response = await api.post(
      `/submit`,
      { "chall_id": challengeId, "flag": flag },
      { withCredentials: true }
    );
    return { status: response.status, data: response.data };
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      return {
        status: error.response.status,
        data: error.response.data,
      };
    }
    console.error("Unexpected error during flag submission:", error);
    return { status: 500, data: { message: "Unexpected error occurred" } };
  }
}

export async function fetchCountries(): Promise<JSON[]> {
  try {
    const response = await api.get("/countries");
    return response.data;
  } catch (error) {
    console.error("Error fetching countries:", error);
    return [];
  }
}

export async function registerTeam(teamName:string, teamPassword: string): Promise<Team | {status: number, error: string}> {
  try {
    const response = await api.post(
      "/teams",
      { name: teamName, password: teamPassword },
      { withCredentials: true }
    );
    const teamData = response.data;
    return { id: teamData.id, name: teamData.name, logo: teamData.logo, members: teamData.members, country: teamData.nationality } as Team;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      return {
        status: error.response.status,
        error: error.response.data,
      };
    }
    console.error("Unexpected error during team registration:", error);
    return { status: 500, error:  "Unexpected error occurred" };
  }
}

export async function updateTeam(teamDescription:string | undefined, teamCountry: string| undefined, teamProfilePicture: string| undefined): Promise<{ status: number; data?: any }> {
  try {
    const response = await api.patch(
      "/teams",
      { bio: teamDescription, nationality: teamCountry, image: teamProfilePicture },
      { withCredentials: true }
    );
    return { status: response.status, data: response.data };
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      return {
        status: error.response.status,
        data: error.response.data,
      };
    }
    console.error("Unexpected error during team update:", error);
    return { status: 500, data: { message: "Unexpected error occurred" } };
  }
}

export async function fetchUserData(userId: number): Promise< User | null > {
  try {
    if (userId === -1) {
      console.warn("Invalid userId provided, returning default user data.");
      return { id: -1, username: "", email: "", role: "", score: -1, country: "", joinedAt: "", solves: [], teamId: null };
    }
    const response = await api.get(`/users/${userId}`);
    switch (response.status) {
      case 200:
        return {id: response.data.id, role: response.data.role, score: response.data.score, profilePicture: response.data.image, username: response.data.name, email: response.data.email, country: response.data.nationality, joinedAt: response.data.joined_at, solves: response.data.solves, teamId: response.data.team_id};
      default:
        return { id:-1, username: "", email: "", role: "", score: -1, country: "", joinedAt: "", solves: [], teamId: null };
    }
  } catch (error) {
    console.error("Error fetching user data:", error);
    return { id:-1, username: "", email: "", role: "", score: -1, country: "", joinedAt: "", solves: [], teamId: null };
  }
}