export interface ServerResponse<t> {
  success: boolean;
  code: number;
  message: string;
  body?: t;
}
