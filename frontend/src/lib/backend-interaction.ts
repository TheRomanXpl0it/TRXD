import { api } from "@/api/axios";
import type { Team, AuthProps, User } from "@/context/AuthProvider";
import axios from "axios";
import type { Challenge as ChallengeType } from "@/context/ChallengeProvider";

export async function getChallengesData(){
    // Simulate file attachments as objects with name and url

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
                solves: challenge.solves,
                points: challenge.points,
                category: challenge.category,
                solved: challenge.solved,
                tags: challenge.tags,
                difficulty: challenge.difficulty,
                hidden: challenge.hidden,
                instanced: challenge.instance,
                timeout: challenge.timeout ? challenge.timeout : undefined,
        };
    });
    console.log("Processed challenges:", challenges);
    categories = challenges.map((challenge: any) => challenge.category);
    categories = [...new Set(categories)]; // Remove duplicates
    return JSON.stringify({ challenges, categories });
} 

export function toChallenge(api: any): ChallengeType {
  let remote = api.host ?? undefined;
  remote += api.port ? `:${api.port}` : "";
  return {
    id: api.id,
    title: api.name,
    description: api.description ?? "",
    category: api.category,                    // adjust if API returns array
    tags: api.tags ?? [],
    points: api.points ?? 0,
    max_points: api.max_points ?? undefined,
    score_type: api.score_type ?? "Static",
    solves: api.solves ?? 0,
    hidden: Boolean(api.hidden),
    authors: api.authors ?? [],
    remote: remote,
    instanced: Boolean(api.instance),
    attachments: api.attachments ?? [],        // map if your API differs
    timeout: api.timeout ?? undefined,
    solved: Boolean(api.solved),
    solves_list: api.solves_list ?? [],
    first_blood: api.first_blood ?? undefined,
    flags: api.flags ?? [],
    docker: api.dockerConfig ?? undefined,
  };
}

export async function fetchChallengeById(
  challengeId: string,
  signal?: AbortSignal
): Promise<ChallengeType> {
  const res = await api.get(`/challenges/${encodeURIComponent(challengeId)}`, {
    signal,
  });
  return toChallenge(res.data);
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
        members: [],
        teamlogo: "/teamLogo.png",
        solves: [],
        country: "",
        score: 0,
        rank: -1,
    };
    try{
        const response = await api.get(`/teams/${teamId}`);
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
      username: data.name,
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
      `/submissions`,
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

export async function startInstance(challengeId: number){
  try {
    const data = { "chall_id": challengeId }
    const response = await api.post(
      `/instances`,
      data,
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
    console.error("Unexpected error during instance start:", error);
    return { status: 500, data: { message: "Unexpected error occurred" } };
  }
}

export async function stopInstance(challengeId: number){
  try {
    const data = { "chall_id": challengeId }
    const response = await api.delete(
      `/instances`,
      {
        data: data,
        withCredentials: true 
      }
    );
    return { status: response.status, data: response.data };
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      return {
        status: error.response.status,
        data: error.response.data,
      };
    }
    console.error("Unexpected error during instance stop:", error);
    return { status: 500, data: { message: "Unexpected error occurred" } };
  }
}