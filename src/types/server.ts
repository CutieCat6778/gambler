import { Token, User } from "./models";

export interface ServerResponse<t> {
  success: boolean;
  code: number;
  message: string;
  body?: t;
}

export interface LoginResponseBody {
  token: Token;
  user: User;
}
