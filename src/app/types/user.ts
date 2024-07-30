export interface User {
  id: number;
  name: string;
  username: string;
  email: string;
  balance: number;
  created_at: CreatedAt;
  updated_at: UpdatedAt;
}

export interface CreatedAt {
  String: string;
  Valid: boolean;
}

export interface UpdatedAt {
  String: string;
  Valid: boolean;
}
