import React, { createContext, useContext, useState, useEffect } from 'react';
import { getChallenges, getCategories } from '@/lib/backend-interaction';

// Types
export type Challenge = {
  id: number;
  title: string;
  description: string;
  solves?: number;
  points?: number;
  category: string;
  remote: string;
  solved: boolean;
  tags?: string[];
  difficulty?: string;
  attachments?: File[];
  authors?: string[];
  hidden: boolean;
  flags?: [{ flag: string; regex: boolean }];
  instanced?: boolean;
  timeout?: Date;
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
      const challengesResult = await getChallenges();
      const challenges = JSON.parse(challengesResult);
      setChallenges(challenges);
    }

    async function fetchCategories() {
      const categoriesResult = await getCategories();
      const categories = JSON.parse(categoriesResult);
      setCategories(categories);
    }

    fetchChallenges();
    fetchCategories();
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