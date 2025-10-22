import { api } from '$lib/api'

export async function deleteFlags(flags: any[], chall_id:any){
  const requests = flags.map(f =>{
    if (f.flag !== "")
      api<any>('/flags', {
        method: 'DELETE',
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ 'flag': f.flag, chall_id }),
      })
  }
  );
  return requests;
}

export async function createFlags(flags: any[], chall_id:any){
  const requests = flags.map(f =>{
    
    if (f.flag !== "")
      api<any>('/flags', {
        method: 'POST',
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ 'flag': f.flag, 'regex': !!f.regex, chall_id }),
      })
  }
  );
  return requests;
}