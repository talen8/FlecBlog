import axios from 'axios'
import type { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { getApiUrl, getApiToken, saveApiToken } from '@/composables/useStore'

interface ApiResponse<T = unknown> {
  code: number
  data: T
  message: string
}

const request = axios.create({
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true,
})

export function setAccessToken(token: string) {
  saveApiToken(token)
}

export function getAccessToken() {
  return getApiToken()
}

request.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  config.baseURL = getApiUrl()
  const token = getApiToken()
  if (config.url !== '/auth/refresh' && token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    if (response.config.responseType === 'blob') return response.data
    const { code, message, data } = response.data
    return code === 0 ? data : Promise.reject(new Error(message || '请求失败'))
  },
  async (error: AxiosError) => {
    return Promise.reject(error)
  }
)

export default request
