"use client";

import {
  createContext,
  type ReactNode,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

import { apiLogout, apiMe, type CurrentUser } from "@/lib/api";

type SessionContextValue = {
  user: CurrentUser | null;
  status: "loading" | "authenticated" | "anonymous";
  refresh: () => Promise<CurrentUser | null>;
  signOut: () => Promise<boolean>;
};
const SessionContext = createContext<SessionContextValue | null>(null);

export function SessionProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<CurrentUser | null>(null);
  const [status, setStatus] = useState<SessionContextValue["status"]>("loading");

  const applySession = useCallback((current: CurrentUser | null) => {
    setUser(current);
    setStatus(current ? "authenticated" : "anonymous");
    return current;
  }, []);

  const refresh = useCallback(async () => {
    const result = await apiMe();
    return applySession(result.ok ? (result.data ?? null) : null);
  }, [applySession]);

  useEffect(() => {
    let active = true;
    void apiMe().then((result) => {
      if (active) applySession(result.ok ? (result.data ?? null) : null);
    });
    return () => {
      active = false;
    };
  }, [applySession]);

  const signOut = useCallback(async () => {
    const result = await apiLogout();
    if (!result.ok) return false;
    applySession(null);
    return true;
  }, [applySession]);
  const value = useMemo(
    () => ({ user, status, refresh, signOut }),
    [user, status, refresh, signOut],
  );
  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
}

export function useSession() {
  const context = useContext(SessionContext);
  if (!context) throw new Error("useSession must be used inside SessionProvider");
  return context;
}
