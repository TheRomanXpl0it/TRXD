import { api } from '$lib/api';

export async function getSolves(
  chall_id: string
): Promise<any[]>{
  const challenge = await getChallenge(chall_id);
  return challenge.solves_list;
}

export async function getChallenges(): Promise<any[]> {
  return api<any[]>('/challenges');
}

export async function getChallenge(chall_id: string): Promise<any> {
  return api<any>(`/challenges/${chall_id}`);
}

export async function submitFlag(
  chall_id: string,
  flag: string
): Promise<{status: string }> {
  return await  api<{ first_blood: boolean; status: string }>(`/submissions`, {
    method: 'POST',
    body: JSON.stringify({ flag,chall_id })
  });
}

export async function getCategories(): Promise<any[]> {
  return api<any[]>('/categories');
}

export async function createCategory(name: string,icon:string): Promise<any> {
  return api<any>('/categories', {
    method: 'POST',
    body: JSON.stringify({ name, icon })
  });
}

export async function createChallenge(
    name:        string,
		category:    string,
		description: string,
		type:        string,
		max_points:   number,
		score_type:   string,
): Promise<any> {
  return api<any>('/challenges', {
    method: 'POST',
    body: JSON.stringify({ name, category, description, type, max_points, score_type })
  });
}

export async function deleteChallenge(chall_id: string): Promise<any> {
  return api<any>(`/challenges`, {
    method: 'DELETE',
    body: JSON.stringify({ chall_id })
  });
}

/** Sends multipart/form-data to PATCH /challenges (server reads chall_id) */
export async function updateChallengeMultipart(fields: any): Promise<any> {
  const fd = new FormData();

  // required by backend
  fd.append("chall_id", String(fields.chall_id));

  // simple helper appenders
  const put = (k: string, v: unknown) => {
    if (v === undefined || v === null || v === "") return;
    fd.append(k, String(v));
  };
  const putMany = (k: string, arr?: unknown[]) => {
    if (!Array.isArray(arr)) return;
    for (const v of arr) {
      if (v === undefined || v === null || v === "") continue;
      fd.append(k, String(v));
    }
  };

  // strings / booleans / enums
  put("name", fields.name);
  put("category", fields.category);
  put("description", fields.description);
  put("difficulty", fields.difficulty); // e.g. "Easy" | "Hard"
  put("type", fields.type);             // "Container" | "Compose" | "Normal"
  if (typeof fields.hidden === "boolean") put("hidden", fields.hidden);
  put("score_type", fields.score_type);  // "DYNAMIC" | "STATIC"
  put("host", fields.host);
  if (Number.isFinite(fields.port)) put("port", fields.port);

  // arrays
  putMany("authors", fields.authors);
  putMany("attachments", fields.attachments);

  // numeric optionals: only send if > 0
  if (Number.isFinite(fields.max_points) && fields.max_points > 0) {
    put("max_points", fields.max_points);
  }
  if (Number.isFinite(fields.lifetime) && fields.lifetime > 0) {
    put("lifetime", fields.lifetime);
  }
  if (Number.isFinite(fields.max_memory) && fields.max_memory > 0) {
    put("max_memory", fields.max_memory);
  }

  // docker-ish
  put("image", fields.image);
  put("compose", fields.compose);
  if (typeof fields.hash_domain === "boolean") put("hash_domain", fields.hash_domain);
  put("envs", fields.envs);
  put("max_cpu", fields.max_cpu);

  // optional files (server accepts any field name)
  if (Array.isArray(fields.files)) {
    for (const f of fields.files) {
      if (f) fd.append("files", f, (f as File).name);
    }
  }

  // IMPORTANT: api() MUST NOT force Content-Type when body is FormData
  return api<any>(`/challenges`, {
    method: "PATCH",
    body: fd,
  });
}