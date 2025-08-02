import { api } from "@/api/axios";
import axios from "axios";


export async function getChallenges(){
    // Simulate file attachments as objects with name and url
    type Attachment = {
        name: string;
        url: string;
    };

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
        },
        {
            challenge: {
                id : 2,
                title: 'Challenge 2',
                description: 'This is the second challenge',
                flag: 'flag{second_challenge}',
                points: 100,
                solves: 1,
                category: 'Web',
                tags: ['Web', 'Medium'],
                difficulty: 'Easy',
                remote: 'https://example.com',
                solved: true,
                attachments: [
                    { name: "dockerfile", url: "/files/challenge2/dockerfile" }
                ] as Attachment[],
                authors: ["author2"],
                hidden: false,
                instanced: false,
                timeout: undefined,
            }
        },
        {
            challenge:
            {
                id : 3,
                title: 'Challenge 3',
                description: 'This is the third challenge',
                flag: 'flag{third_challenge}',
                points: 100,
                solves: 0,
                category: 'Pwn',
                tags: ['Stack', 'Easy'],
                difficulty: 'Easy',
                remote: 'https://example.com',
                solved: false,
                attachments: [
                    { name: "executable", url: "/files/challenge3/executable" }
                ] as Attachment[],
                authors: ["author1"],
                hidden: false,
                instanced: true,
                timeout: undefined,
            }
        },
        {   
            challenge:{
            id : 4,
            title: 'Challenge 4',
            description: 'This is the fourth challenge',
            flag: 'flag{fourth_challenge}',
            points: 100,
            solves: 5,
            category: 'Rev',
            tags: ['RC4', 'Medium'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [
                { name: "Dockerfile", url: "/files/challenge4/Dockerfile" }
            ] as Attachment[],
            authors: ["admin","author2"],
            hidden: false,
            instanced: true,
            timeout: new Date(Date.now() + 5 * 60 * 1000),
            }
        },
        {challenge:{
            id : 5,
            title: 'Challenge 5',
            description: 'This is the fifth challenge',
            flag: 'flag{fifth_challenge}',
            points: 100,
            solves: 10,
            category: 'Crypto',
            tags: ['Diffie-Hellman', 'Easy'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [],
            authors: ["author2"],
            hidden: false,
            instanced: true,
            timeout: undefined,
        }},
        {challenge:{
            id : 6,
            title: 'Challenge 6',
            description: 'This is the sixth challenge',
            flag: 'flag{sixth_challenge}',
            points: 100,
            solves: 0,
            category: 'Misc',
            tags: ['Packets', 'Easy'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [
                { name: "packets", url: "/files/challenge6/packets" }
            ] as Attachment[],
            authors: ["author1,author2"],
            hidden: false,
            instanced: true,
        }},
        {challenge:{
            id : 7,
            title: 'Challenge 7',
            description: 'This is the seventh challenge',
            flag: 'flag{seventh_challenge}',
            points: 100,
            solves: 0,
            category: 'Forensics',
            tags: ['Random stuff', 'Easy'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [
                { name: "packets", url: "/files/challenge7/packets" }
            ] as Attachment[],
            authors: ["author1,author2"],
            hidden: false,
            instanced: true,
            timeout: undefined,
        }},
        {challenge:{
            id : 8,
            title: 'Challenge 8',
            description: 'This is the eighth challenge',
            flag: 'flag{eighth_challenge}',
            points: 100,
            solves: 0,
            category: 'Crypto',
            tags: ['Hash', 'Easy'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [
                { name: "hashes", url: "/files/challenge8/hashes" }
            ] as Attachment[],
            authors: ["author1, author2"],
            hidden: false,
            instanced: true,
            timeout: undefined,
        }},
        {challenge:{
            id : 9,
            title: 'Challenge 9',
            description: 'This is the ninth challenge',
            flag: 'flag{ninth_challenge}',
            points: 100,
            solves: 0,
            category: 'Pwn',
            tags: ['Kpwn', 'Hard'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [
                { name: "executable", url: "/files/challenge9/executable" }
            ] as Attachment[],
            authors: ["author1"],
            hidden: false,
            instanced: true,
            timeout: undefined,
        }},
        {challenge:{
            id : 10,
            title: 'Challenge 10',
            description: 'This is the tenth challenge',
            flag: 'flag{tenth_challenge}',
            points: 100,
            solves: 0,
            category: 'Web',
            tags: ['ClientSide', 'Hard'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [
                { name: "ciao", url: "/files/challenge10/ciao" }
            ] as Attachment[],
            authors: ["admin"],
            hidden: false,
            instanced: true,
            timeout: undefined,
        }}
    ];
    return JSON.stringify(mockChallenges);
} 

export async function getCategories(){
    const categories = [
        'Web',
        'Rev',
        'Forensics',
        'Crypto',
        'Pwn',
        'Misc',
    ]

    return JSON.stringify(categories);
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
    return { status: response.status, data: response.data };
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

export async function fetchTeamData(): Promise<{ name: string; members: string[]; score: number; teamlogo:string, rank: number }> {
    // Simulate fetching team data
    return {
        name: "Un team popo demmerda",
        members: ["Alice", "Bob"],
        teamlogo: "/teamLogo.png",
        score: 1500,
        rank: 1,
    };
}

export async function leaveTeam() {
    // Simulate leaving a team
    console.log("Leaving team...");
    return { success: true, message: "You have left the team." };
}