import { Token, User } from "@/types/models";
import { ServerResponse } from "@/types/server";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

export default function useLocalStorage() {
  const router = useRouter();
  const [tokens, setTokens] = useState<Token | null>(null);
  const [user, setUserData] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  function navigateToLogin() {
    if (!["/login", "/register"].includes(router.pathname)) {
      router.push("/login");
    }
  }

  function navigateToHome() {
    if (["/login", "/register"].includes(router.pathname)) {
      router.push("/");
    }
  }

  useEffect(() => {
    async function getSelf() {
      if (!tokens) return navigateToLogin();
      const res = await fetch("http://localhost:3000/user/@me", {
        headers: {
          Authorization: `Bearer ${tokens.accessToken}`,
        },
      });

      const data: ServerResponse<User> = await res.json();
      if (!data.success || !data.body) {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        localStorage.removeItem("isAuthenticated");
        localStorage.removeItem("exp");
        navigateToLogin();
        setIsLoading(false);
      } else {
        setUser(data.body);
      }
    }
    const token = localStorage.getItem("token");
    const user = localStorage.getItem("user");
    if (token && user) {
      setTokens(JSON.parse(token));
      setUser(JSON.parse(user));
    } else if (token && !user) {
      getSelf();
    } else {
      navigateToLogin();
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    if (tokens && user) {
      setIsLoading(false);
    }
  }, [tokens, user]);

  function setToken(token: Token) {
    localStorage.setItem("token", JSON.stringify(token));
    localStorage.setItem("isAuthenticated", "true");
    localStorage.setItem("exp", (Date.now() + 1000 * 60 * 60 * 24).toString());
    setTokens(token);
  }

  function setUser(user: User) {
    localStorage.setItem("user", JSON.stringify(user));
    setUserData(user);
  }

  function setAuth(token: Token, user: User) {
    setToken(token);
    setUser(user);
    navigateToHome();
  }

  function getUser() {
    return user;
  }

  function getToken() {
    return tokens;
  }

  return {
    tokens,
    setToken,
    user,
    setUser,
    isLoading,
    setAuth,
    getUser,
    getToken,
  };
}
