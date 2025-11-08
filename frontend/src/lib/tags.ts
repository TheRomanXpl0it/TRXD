import { api } from '$lib/api'

// CREATE tags for a challenge (POST)
export async function createTagsForChallenge(tags: any[], chall_id: any) {
  const list = Array.from(tags ?? [])
    .map((name) =>
      api<any>('/tags', {
        method: 'POST',
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ name, chall_id }), // let `api` set JSON headers & stringify
      })
    );

  return Promise.all(list);
}

// DELETE tags from a challenge (DELETE)
export async function deleteTagsFromChallenge(tags: any[], chall_id: any) {
  const list = Array.from(tags ?? []).map((name) =>
    api<any>('/tags', {
      method: 'DELETE',
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ name, chall_id }), // no toLowerCase, no trim
    })
  );

  return Promise.all(list);
}
