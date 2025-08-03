import { api } from "@/api/axios";
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
}): Promise<{ status: number; data?: any }> {
  try {
    const response = await api.post(
      "/login",
      { email, password },
      { withCredentials: true }
    );

    console.log("Login response:", response);
    return { status: response.status, data: await response.data };
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      return {
        status: error.response.status,
        data: error.response.data,
      };
    }

    console.error("Unexpected login error:", error);
    return {
      status: 500,
      data: { message: "Unexpected error occurred" },
    };
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
}): Promise<{ status: number; data?: any }> {
  try {
    const response = await api.post(
      "/register",
      { username, email, password },
      { withCredentials: true }
    );
    return { status: response.status, data: response.data };
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      // Return the server's status and any data it returned
      return {
        status: error.response.status,
        data: error.response.data,
      };
    }

    // Unknown error (network error or non-Axios error)
    console.error("Unexpected error during registration:", error);
    return {
      status: 500,
      data: { message: "Unexpected error occurred" },
    };
  }
}

export async function fetchTeamData(teamId: number): Promise<{ name: string; members: string[]; score: number; teamlogo:string, rank: number }> {
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

export async function getUsersTeamData(): Promise<{ id: number, name: string; }> {
    // Simulate fetching team data
    let teamData = { id: -1, name: "" };
    try {
        const response = await api.get("/team");
        if (response.status === 200) {
            return response.data;
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

export async function checkSession(): Promise<{ status: number; data?: any; team?: any }> {
  try {
    const response = await api.get("/auth"); // or your auth check endpoint
    const data = response.data;

    const team = data.team ?? null;

    return { status: response.status, data, team };
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      return {
        status: error.response.status,
        data: error.response.data,
        team: null,
      };
    }
    return {
      status: 500,
      data: { message: "Unexpected error" },
      team: null,
    };
  }
}

export async function submitFlag(
  challengeId: number,
  flag: string
): Promise<{ status: number; data?: any }> {
  try {
    const response = await api.post(
      `/player/submit`,
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

export async function registerTeam(teamName:string, teamPassword: string): Promise<{ status: number; data?: any }> {
  try {
    const response = await api.post(
      "/player/register-team",
      { name: teamName, password: teamPassword },
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
    console.error("Unexpected error during team registration:", error);
    return { status: 500, data: { message: "Unexpected error occurred" } };
  }
}

export async function updateTeam(teamDescription:string | undefined, teamCountry: string| undefined, teamProfilePicture: string| undefined): Promise<{ status: number; data?: any }> {
  try {
    const response = await api.patch(
      "/player/update-team",
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