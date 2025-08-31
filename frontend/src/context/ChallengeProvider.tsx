import React, { createContext, useContext, useState, useEffect } from 'react';
import { getChallengesData } from '@/lib/backend-interaction';
import { WeekNumberLabel } from 'react-day-picker';

// Types

export type DockerConfig = {
  image: string;
  compose: string;
  hashDomain: boolean;
  lifetime: number;
  envs: string;
  maxMemory: number;
  maxCPU: number;
};

export type Solve = {
  id: number;
  name: string;
  timestamp: string;
};


export type Challenge = {
  id: number;
  title: string;
  description: string;
  solves?: number;
  points?: number;
  category: string;
  remote: string;
  solved: boolean;
  first_blood?: boolean;
  solves_list?: Solve[];
  tags?: string[];
  score_type: "Static" | "Dynamic";
  max_points?: number;
  difficulty?: string;
  attachments?: File[];
  authors?: string[];
  hidden: boolean;
  flags?: [{ flag: string; regex: boolean }];
  instanced?: boolean;
  timeout?: number;
  docker?: DockerConfig;
};

// Context type
export type ChallengeContextType = {
  challenges: Challenge[];
  setChallenges: React.Dispatch<React.SetStateAction<Challenge[]>>;
  categories: string[];
  setCategories: React.Dispatch<React.SetStateAction<string[]>>;
  solveChallenge: (id: number) => void;
};

// Context
const ChallengeContext = createContext<ChallengeContextType | undefined>(undefined);

// Provider
function ChallengeProvider({ children }: { children: React.ReactNode }) {
  const [challenges, setChallenges] = useState<Challenge[]>([]);
  const [categories, setCategories] = useState<string[]>([]);

  const solveChallenge = (id: number) => {
    setChallenges(prev => prev.map(ch => ch.id === id ? { ...ch, solved: true } : ch));
  };

  useEffect(() => {
    async function fetchChallenges() {
      const challengesResult = await getChallengesData();
      const { challenges, categories } = JSON.parse(challengesResult);
      setChallenges(challenges);
      setCategories(categories);
    }

    fetchChallenges();
  }, []);

  return (
    <ChallengeContext.Provider value={{ challenges, setChallenges, categories, setCategories, solveChallenge }}>
      {children}
    </ChallengeContext.Provider>
  );
}

// Hook
function useChallenges() {
  const ctx = useContext(ChallengeContext);
  if (!ctx) throw new Error("useChallenges must be used inside ChallengeProvider");
  return ctx;
}

export { ChallengeContext, ChallengeProvider, useChallenges };