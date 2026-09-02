export class ApiError extends Error {
  code: string;
  details?: { fields?: string[] };

  constructor(code: string, details?: { fields?: string[] }) {
    super(code);
    this.name = "ApiError";
    this.code = code;
    this.details = details;
  }
}

type ApiResponse<T> = {
  ok: boolean;
  data?: T;
  error?: ApiError;
};

type LoginResponse = {
  access_token_expires_in: string;
  refresh_token_expires_in: string;
  role: UserRole;
};

export type UserRole = "citizen" | "councillor" | "office_member";

type CitizenRegisterResponse = {
  name: string;
  email: string;
  city: string;
  state: string;
};

export type DemandStatus =
  "registered" | "under_review" | "in_progress" | "awaiting_confirmation" | "completed";

export type Demand = {
  id: number;
  protocol: string;
  title: string;
  description: string;
  category: string;
  street: string;
  number: string;
  neighborhood: string;
  city: string;
  state: string;
  latitude: number;
  longitude: number;
  images: string[];
  directed_office_id?: number;
  responsible_office_id?: number;
  confirmation_due_at?: string;
  status: DemandStatus;
  support_count: number;
  comment_count: number;
  created_at: string;
  updated_at: string;
};

export type CreateDemandInput = {
  title: string;
  description: string;
  category: string;
  street: string;
  number?: string;
  neighborhood: string;
  city: string;
  state: string;
  latitude: number;
  longitude: number;
  images?: File[];
  directed_office_id?: number;
};

export type DemandFilters = Partial<{
  city: string;
  state: string;
  neighborhood: string;
  category: string;
  status: DemandStatus;
}>;

async function request<T>(
  path: string,
  options: {
    method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
    body?: unknown;
  } = {},
): Promise<ApiResponse<T>> {
  try {
    const { method = "GET", body } = options;
    const res = await fetch(path, {
      method,
      headers: body ? { "Content-Type": "application/json" } : undefined,
      body: body ? JSON.stringify(body) : undefined,
      credentials: "include",
    });

    const json = await res.json().catch(() => ({}));

    if (!res.ok) {
      const code = (json as { code?: string }).code ?? (res.status === 0 ? "NETWORK" : "INTERNAL");
      const details = (json as { details?: { fields?: string[] } }).details;
      return { ok: false, error: new ApiError(code, details) };
    }

    return { ok: true, data: json as T };
  } catch {
    return {
      ok: false,
      error: new ApiError("NETWORK"),
    };
  }
}

async function requestMultipart<T>(path: string, form: FormData): Promise<ApiResponse<T>> {
  try {
    const res = await fetch(path, { method: "POST", body: form, credentials: "include" });
    const json = await res.json().catch(() => ({}));
    if (!res.ok) {
      const code = (json as { code?: string }).code ?? "INTERNAL";
      const details = (json as { details?: { fields?: string[] } }).details;
      return { ok: false, error: new ApiError(code, details) };
    }
    return { ok: true, data: json as T };
  } catch {
    return { ok: false, error: new ApiError("NETWORK") };
  }
}

function multipart(values: Record<string, string | number | File | File[] | undefined>) {
  const form = new FormData();
  Object.entries(values).forEach(([key, value]) => {
    if (value === undefined) return;
    if (Array.isArray(value)) return value.forEach((file) => form.append(key, file));
    form.append(key, value instanceof File ? value : String(value));
  });
  return form;
}

export function apiLogin(email: string, password: string) {
  return request<LoginResponse>("/api/auth/login", {
    method: "POST",
    body: { email: email.trim().toLowerCase(), password },
  });
}
export type CurrentUser = {
  id: number;
  name: string;
  email: string;
  role: UserRole;
  image_url?: string;
};
export function apiMe() {
  return request<CurrentUser | null>("/api/auth/me");
}
export function apiLogout() {
  return request<void>("/api/auth/logout", { method: "POST" });
}

export function apiRegisterCitizen(input: {
  name: string;
  email: string;
  password: string;
  city: string;
  state: string;
}) {
  return request<CitizenRegisterResponse>("/api/auth/register/citizen", {
    method: "POST",
    body: {
      ...input,
      name: input.name.trim(),
      email: input.email.trim().toLowerCase(),
      city: input.city.trim(),
      state: input.state.trim().toUpperCase(),
    },
  });
}

export function apiRegisterCouncillor(input: {
  name: string;
  email: string;
  password: string;
  party: string;
  photo: File;
  city: string;
  state: string;
}) {
  return requestMultipart<CitizenRegisterResponse>(
    "/api/auth/register/councillor",
    multipart({
      ...input,
      name: input.name.trim(),
      email: input.email.trim().toLowerCase(),
      party: input.party.trim(),
      city: input.city.trim(),
      state: input.state.trim().toUpperCase(),
    }),
  );
}

export type PoliticalParty = { id: number; sigla: string; nome: string; uri: string };
export function apiListParties() {
  return request<PoliticalParty[]>("/api/parties");
}

export function apiCreateDemand(input: CreateDemandInput) {
  return requestMultipart<Demand>("/api/demands", multipart(input));
}

export function apiListDemands(filters: DemandFilters = {}) {
  const searchParams = new URLSearchParams();
  Object.entries(filters).forEach(([key, value]) => {
    if (value) {
      searchParams.set(key, value);
    }
  });

  const query = searchParams.toString();
  return request<Demand[]>(`/api/demands${query ? `?${query}` : ""}`);
}

export function apiListMyDemands() {
  return request<Demand[]>("/api/demands/mine");
}

export function apiListOfficeDemands(status?: DemandStatus) {
  return request<Demand[]>(`/api/demands/office${status ? `?status=${status}` : ""}`);
}

export type DemandEvent = {
  id: number;
  demand_id: number;
  type: string;
  actor_user_id?: number;
  actor_name?: string;
  actor_role?: UserRole;
  actor_image_url?: string;
  message?: string;
  images: string[];
  created_at: string;
};
export type DemandComment = {
  id: number;
  demand_id: number;
  parent_id?: number;
  author_id: number;
  author_name: string;
  author_role: UserRole;
  author_image_url?: string;
  body: string;
  images: string[];
  hidden: boolean;
  hidden_at?: string;
  created_at: string;
};
export type DemandActivity = { events: DemandEvent[]; comments: DemandComment[] };
export type DemandSupport = {
  support_count: number;
  supported: boolean;
  can_support: boolean;
};
export type Notification = {
  id: number;
  type: string;
  demand_id: number;
  read_at?: string;
  created_at: string;
};
export function apiGetDemand(id: number) {
  return request<Demand>(`/api/demands/${id}`);
}
export function apiDemandActivity(id: number) {
  return request<DemandActivity>(`/api/demands/${id}/activity`);
}
export function apiDemandAction(
  id: number,
  action: "claim" | "start" | "request-confirmation" | "confirm" | "reopen",
  body?: { message?: string; images?: File[] },
) {
  if (action === "confirm")
    return request<Demand>(`/api/demands/${id}/${action}`, { method: "POST" });
  return requestMultipart<Demand>(`/api/demands/${id}/${action}`, multipart(body ?? {}));
}
export function apiCreateDemandMilestone(id: number, body: { message: string; images?: File[] }) {
  return requestMultipart<void>(`/api/demands/${id}/milestones`, multipart(body));
}
export function apiDemandSupport(id: number) {
  return request<DemandSupport>(`/api/demands/${id}/support`);
}
export function apiAddDemandSupport(id: number) {
  return request<DemandSupport>(`/api/demands/${id}/support`, { method: "PUT" });
}
export function apiRemoveDemandSupport(id: number) {
  return request<DemandSupport>(`/api/demands/${id}/support`, { method: "DELETE" });
}
export function apiCommentDemand(
  id: number,
  body: { body?: string; images?: File[]; parent_id?: number },
) {
  return requestMultipart<DemandComment>(`/api/demands/${id}/comments`, multipart(body));
}
export function apiDeleteComment(id: number) {
  return request<void>(`/api/comments/${id}`, { method: "DELETE" });
}
export function apiReportComment(id: number, reason: string) {
  return request<void>(`/api/comments/${id}/report`, { method: "POST", body: { reason } });
}
export function apiHideComment(id: number) {
  return request<void>(`/api/comments/${id}/hide`, { method: "POST" });
}
export function apiNotifications() {
  return request<Notification[]>("/api/notifications");
}
export function apiReadNotifications(ids: number[] = []) {
  return request<void>("/api/notifications/read", { method: "POST", body: { ids } });
}

export type OfficeContact = { type: string; value: string; position: number };
export type OfficeProfile = {
  office_id: number;
  councillor_id: number;
  slug: string;
  party: string;
  description: string;
  city: string;
  state: string;
  contacts: OfficeContact[];
  social_networks: OfficeContact[];
};
export type CreatedOffice = Pick<
  OfficeProfile,
  "office_id" | "councillor_id" | "description" | "contacts" | "social_networks"
>;

export function apiUpdateOffice(
  input: Pick<
    OfficeProfile,
    "party" | "description" | "city" | "state" | "contacts" | "social_networks"
  >,
) {
  return request<OfficeProfile>("/api/office", { method: "PUT", body: input });
}
export type ManagedOffice = PublicOffice & {
  members: Array<{ user_id: number; name: string; email: string; image_url: string }>;
};
export type OfficeMemberRequest = {
  id: number;
  email: string;
  expires_at: string;
  created_at: string;
};
export function apiGetManagedOffice() {
  return request<ManagedOffice>("/api/office/me");
}
export function apiCreateOffice() {
  return request<CreatedOffice>("/api/office", {
    method: "POST",
    body: { contacts: [], social_networks: [] },
  });
}

export function apiInviteOfficeMember(email: string) {
  return request<{ email: string; expires_at: string }>("/api/office/member-request", {
    method: "POST",
    body: { email },
  });
}
export function apiListOfficeMemberRequests() {
  return request<OfficeMemberRequest[]>("/api/office/member-requests");
}
export function apiCancelOfficeMemberRequest(id: number) {
  return request<void>(`/api/office/member-requests/${id}`, { method: "DELETE" });
}
export function apiRemoveOfficeMember(id: number) {
  return request<void>(`/api/office/members/${id}`, { method: "DELETE" });
}

export function apiListOffices(city: string, state: string) {
  const query = new URLSearchParams({ city, state }).toString();
  return request<Array<{ office_id: number; councillor_name: string; party: string }>>(
    `/api/office?${query}`,
  );
}

export type PublicOffice = {
  office_id: number;
  slug: string;
  councillor_name: string;
  party: string;
  image_url: string;
  city: string;
  state: string;
  description: string;
  contacts: OfficeContact[];
  social_networks: OfficeContact[];
};
export function apiSearchOffices(filters: { q?: string; city?: string; state?: string } = {}) {
  const search = new URLSearchParams();
  Object.entries(filters).forEach(([key, value]) => {
    if (value) search.set(key, value);
  });
  return request<PublicOffice[]>(`/api/office${search.size ? `?${search}` : ""}`);
}
export function apiGetOffice(slug: string) {
  return request<PublicOffice>(`/api/office/${encodeURIComponent(slug)}`);
}

export function apiRegisterOfficeMember(input: {
  token: string;
  name: string;
  password: string;
  photo: File;
}) {
  return requestMultipart<{ name: string; email: string; office_id: number; image_url: string }>(
    "/api/auth/register/office-member",
    multipart(input),
  );
}
export type OfficeMemberInvitation = {
  office_id: number;
  councillor_name: string;
  party: string;
  image_url: string;
  city: string;
  state: string;
  expires_at: string;
};
export function apiPreviewOfficeMemberInvitation(token: string) {
  return request<OfficeMemberInvitation>(
    `/api/auth/register/office-member/invitation?${new URLSearchParams({ token })}`,
  );
}
