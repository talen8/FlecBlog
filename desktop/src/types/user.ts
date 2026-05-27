export interface LoginParams {
  email: string
  password: string
}

export interface LoginResponse {
  access_token: string
}

export interface User {
  id: number
  email: string
  nickname: string
  avatar: string
  role: string
}
