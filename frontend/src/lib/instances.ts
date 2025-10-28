import { api } from "./api";

export async function startInstance(
  chall_id: string
): Promise<{ host: string; port: string; timeout: string }> {
  return api<{ host: string; port: string; timeout:string }>(`/instances`, {
    headers: { 'content-type': 'application/json' },
    method: 'POST',
    body: JSON.stringify({ chall_id })
  });
}

export async function stopInstance(
  chall_id: string
): Promise<void> {
  return api<void>(`/instances`, {
    headers: { 'content-type': 'application/json' },
    method: 'DELETE',
    body: JSON.stringify({ chall_id })
  });
}