import { api } from '$lib/api'

export async function getScoreboard(){
  return api<any>('/scoreboard');
}