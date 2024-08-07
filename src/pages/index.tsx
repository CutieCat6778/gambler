import useLocalStorage from "@/hooks/useLocalStorage";
import { User } from "@/types/models";
import { ServerResponse } from "@/types/server";
import axios, { AxiosResponse } from "axios";
import { useRouter } from "next/router";

export default function Home() {
  const router = useRouter();
  const { user, isLoading, setUser, tokens } = useLocalStorage();

  function isUserBalanceChanged(u: User) {
    const prevBal = u.balance_history.slice(-1)[0].amount;
    if (prevBal > u.balance) {
      return "red-400";
    } else if (prevBal < u.balance) {
      return "green-400";
    } else {
      return "amber-400";
    }
  }

  async function getUser() {
    if (!tokens) {
      setUser({} as User);
      router.push("/login");
      return;
    }
    try {
      const res: AxiosResponse<ServerResponse<User>> = await axios.get(
        `${process.env.NEXT_PUBLIC_API_URL}/user/@me`,
        {
          headers: {
            Authorization: `Bearer ${tokens.accessToken}`,
          },
        },
      );

      const data = res.data;

      if (data.success && data.body) {
        console.log(data.body);
        setUser(data.body);
      } else {
        setUser({} as User);
        router.push("/login");
      }
    } catch (e) {
      console.log(e);
    }
  }

  if (!isLoading && user) {
    return (
      <main>
        <header className="w-full flex justify-center items-center flex-col">
          <h1>Hello, {user.username}</h1>
          <span>
            Wellcome to Gambler, here you will be the richest men in the world!
            (But only if you never give up!)
          </span>
          <button
            onClick={() => {
              getUser();
            }}
          >
            Refresh
          </button>
        </header>
        <section className="w-full flex justify-center my-10">
          <div className="w-max">
            <h3 className="text-2xl font-semibold w-full">Balance</h3>
            <span
              className={`text-center w-full font-bold block text-${isUserBalanceChanged(user)}`}
            >
              {user.balance} €
            </span>
          </div>
        </section>
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
