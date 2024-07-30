"use client";

import { User } from "@/types/models";
import { ServerResponse } from "@/types/server";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function Home() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    async function getSelf() {
      const response = await fetch("http://localhost:3000/user/@me");
      const data: ServerResponse<User> = await response.json();
      if (!data.success) {
        router.push("/login");
        return;
      }
      setUser(data.body || null);
    }
    getSelf();
  });

  console.log(user);
  if (user) {
    return (
      <main>
        <h1>Hello, {user.username}</h1>
      </main>
    );
  } else {
    return (
      <main>
        <h1>Loading...</h1>
      </main>
    );
  }
}
