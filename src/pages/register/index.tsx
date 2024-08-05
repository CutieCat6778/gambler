import useLocalStorage from "@/hooks/useLocalStorage";
import { LoginResponseBody, ServerResponse } from "@/types/server";
import { useRouter } from "next/router";
import { FormEvent, useState } from "react";

export default function Register() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const regex = RegExp(/^(?=.*[\x21-\x7E])(?!.*:).{8,}$/gm);
  const { setUser, setToken, isLoading } = useLocalStorage();

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const res = await fetch("http://localhost:3000/auth/register", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password, email, name }),
    });

    const data: ServerResponse<LoginResponseBody> = await res.json();
    if (data.success && data.body) {
      setUser(data.body.user);
      setToken(data.body.token);
      router.push("/");
    } else {
      console.log(data);
    }
  }

  const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUsername(e.target.value);
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const password = e.target.value;
    if (regex.test(password) || password.length < 8) {
      setPassword(e.target.value);
    }
  };
  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };
  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };

  return (
    <div className="w-full h-screen flex flex-col justify-center items-center">
      <h1 className="text-3xl font-bold">Register</h1>
      <form onSubmit={handleSubmit} className="flex flex-col">
        <input
          type="text"
          placeholder="Full Name"
          value={name}
          onChange={handleNameChange}
          required
        />
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={handleUsernameChange}
          required
        />
        <input
          type="email"
          placeholder="E-Mail"
          value={email}
          onChange={handleEmailChange}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={handlePasswordChange}
          required
        />
        <button
          type="submit"
          disabled={username.length <= 3 && !isLoading}
          className="my-1 px-2 py-1 font-bold rounded text-xl bg-gray-400 text-black"
        >
          Login
        </button>
      </form>
      <a className="my-3" href="/login">
        {isLoading ? "Loading..." : "Already have an account? Login here."}
      </a>
    </div>
  );
}
