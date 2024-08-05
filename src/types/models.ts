// Define the GameType enum
export enum GameType {
  BlackJack = 1,
  Roulette,
  Slots,
}

export interface User {
  name: string;
  username: string;
  password: string;
  email: string;
  balance: number;
  created_at: UserDate;
  updated_at: UserDate;
  games: Games[];
  balance_history: BalanceHistory[];
}

export interface BalanceHistory {
  user_id: number;
  amount: number;
  reason: string;
  created_at: UserDate;
  updated_at: UserDate;
}

export interface UserDate {
  String: string;
  Valid: boolean;
}

export interface Games {
  id: number;
  created_at: string | null;
  closed_at: string | null;
  type: GameType; // Use the GameType enum here
  users: User[];
}

export interface Token {
  accessToken: string;
  refreshToken: string;
}
