import { time } from "console";

export async function getChallenges(){
    const mockChallenges = [
        {
            challenge : {
                id : 1,
                title: 'Challenge 1',
                description: 'This is the first challenge',
                points: 100,
                solves: 0,
                category: 'Web',
                tags: ['Web', 'Easy'],
                difficulty: 'Easy',
                remote: 'https://example.com',
                solved: false,
                attachments: [],
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
                points: 100,
                solves: 1,
                category: 'Web',
                tags: ['Web', 'Medium'],
                difficulty: 'Easy',
                remote: 'https://example.com',
                solved: true,
                attachments: [],
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
                points: 100,
                solves: 0,
                category: 'Pwn',
                tags: ['Stack', 'Easy'],
                difficulty: 'Easy',
                remote: 'https://example.com',
                solved: false,
                attachments: [],
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
            points: 100,
            solves: 5,
            category: 'Rev',
            tags: ['RC4', 'Medium'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [],
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
            points: 100,
            solves: 0,
            category: 'Misc',
            tags: ['Packets', 'Easy'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [],
            authors: ["author1,author2"],
            hidden: false,
            instanced: true,
        }},
        {challenge:{
            id : 7,
            title: 'Challenge 7',
            description: 'This is the seventh challenge',
            points: 100,
            solves: 0,
            category: 'Forensics',
            tags: ['Random stuff', 'Easy'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [],
            authors: ["author1,author2"],
            hidden: false,
            instanced: true,
            timeout: undefined,
        }},
        {challenge:{
            id : 8,
            title: 'Challenge 8',
            description: 'This is the eighth challenge',
            points: 100,
            solves: 0,
            category: 'Crypto',
            tags: ['Hash', 'Easy'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [],
            authors: ["author1, author2"],
            hidden: false,
            instanced: true,
            timeout: undefined,
        }},
        {challenge:{
            id : 9,
            title: 'Challenge 9',
            description: 'This is the ninth challenge',
            points: 100,
            solves: 0,
            category: 'Pwn',
            tags: ['Kpwn', 'Hard'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [],
            authors: ["author1"],
            hidden: false,
            instanced: true,
            timeout: undefined,
        }},
        {challenge:{
            id : 10,
            title: 'Challenge 10',
            description: 'This is the tenth challenge',
            points: 100,
            solves: 0,
            category: 'Web',
            tags: ['ClientSide', 'Hard'],
            difficulty: 'Easy',
            remote: 'https://example.com',
            solved: false,
            attachments: [],
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

export function login(username:string,password:string){
    if(username === 'admin' && password === 'admin'){
        return true;
    }
    return false;
}