import { api } from '$lib/api'

export async function getScoreboard(){
  return api<any>('/scoreboard');
}

export async function getGraphData(){
  return api<any>('/scoreboard/graph')
}