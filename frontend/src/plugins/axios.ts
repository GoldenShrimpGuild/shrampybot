// axios.ts

import App from '../App.vue'
import axiosOrig, { AxiosRequestConfig } from 'axios'
import { useAuthStore } from '../stores/auth'
import { useGlobalStore } from '../stores/global-store'

interface AxiosOptions {
  baseUrl?: string
  token?: string
}

const GlobalStore = useGlobalStore()

export const axiosConfig = {
  baseURL: GlobalStore.getApiBaseUrl(),
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json'
  }
} as AxiosRequestConfig

export const useAxios = () => {
    const AuthStore = useAuthStore()
    const GlobalStore = useGlobalStore()

    const axiosInstance = axiosOrig.create(axiosConfig)

    axiosInstance.interceptors.request.use(request => {
      const accessToken = GlobalStore.$state.isDevEnvironment ? AuthStore.$state.accessTokenDev : AuthStore.$state.accessTokenProd
      if (accessToken) {
        request.headers['Authorization'] = `Bearer ${accessToken}`;
      }
      return request;
    }, error => {
      return Promise.reject(error);
    });

    axiosInstance.interceptors.response.use(
      response => response, // Directly return successful responses.
      async error => {
        const originalRequest = error.config;
        if (error.response.status === 401 && !originalRequest._retry) {
          // const AuthStore = useAuthStore()

          originalRequest._retry = true; // Mark the request as retried to avoid infinite loops.
          try {
            // Make a request to your auth server to refresh the token.
            const response = await axiosOrig.post('/auth/refresh', {}, {
              baseURL: GlobalStore.getApiBaseUrl(),
              withCredentials: true,
              headers: {
                'Content-Type': 'application/json',
              },
            });
            const { access } = response.data;

            // Store the new access tokens.
            if (GlobalStore.$state.isDevEnvironment) {
              AuthStore.$state.accessTokenDev = access
            } else {
              AuthStore.$state.accessTokenProd = access
            }
            
            // Update the authorization header with the new access token.
            axiosInstance.defaults.headers.common['Authorization'] = `Bearer ${access}`;
            return axiosInstance(originalRequest); // Retry the original request with the new access token.
          } catch (refreshError) {
            // Handle refresh token errors by clearing stored tokens and redirecting to the login page.
            console.error('Token refresh failed:', refreshError);
            if (GlobalStore.$state.isDevEnvironment) {
              AuthStore.$state.accessTokenDev = ''
            } else {
              AuthStore.$state.accessTokenProd = ''
            }
            window.location.href = '/auth/login';
            return Promise.reject(refreshError);
          }
        }
        return Promise.reject(error); // For all other errors, return the error as is.
      }
    );
    return axiosInstance
  }