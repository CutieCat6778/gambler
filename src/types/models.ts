export interface User {
  id: number;
  name: string;
  username: string;
  password: string;
  email: string;
  balance: number;
  createdAt: string | null;
  updatedAt: string | null;
}
