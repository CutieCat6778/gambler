import useLocalStorage from "@/hooks/useLocalStorage";
import { User } from "@/types/models";
import { useRouter } from "next/router";

export default function Home() {
  const router = useRouter();
  const { user, isLoading } = useLocalStorage();

  function isUserBalanceChanged(u: User) {
    console.log(u);
    const prevBal = u.balance_history.slice(-1)[0].amount;
    if (prevBal > u.balance) {
      return "red-400";
    } else if (prevBal < u.balance) {
      return "green-400";
    } else {
      return "amber-400";
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
