import request, { setAccessToken as setToken } from './request'
import type { LoginParams, LoginResponse, User } from '@/types/article'

export function login(params: LoginParams): Promise<LoginResponse> {
  return request.post('/auth/login', params)
}

export function getProfile(): Promise<User> {
  return request.get('/user/profile')
}

export function logout(): Promise<void> {
  return request.post('/auth/logout')
}

export { setToken as setAccessToken }
