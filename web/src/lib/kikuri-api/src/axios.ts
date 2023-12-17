import axios, { AxiosInstance } from "axios";
import { AccessTokenModel } from "./models";

export class AxiosClient {
  private axiosInstance: AxiosInstance;
  private accessToken?: AccessTokenModel;
  private refreshingToken: boolean = false;

  constructor(baseURL: string) {
    this.axiosInstance = axios.create({ baseURL });
    this.setupInterceptors();
  }

  private setupInterceptors() {
    this.axiosInstance.interceptors.request.use(
      async (config) => {
        // Refresh the token if it's expired or not available
        if (!this.accessToken || this.isTokenExpired(this.accessToken)) {
          await this.refreshToken();
        }

        // Add the token to every request
        if (this.accessToken) {
          config.headers[
            "Authorization"
          ] = `Accesstoken ${this.accessToken.token}`;
        }

        return config;
      },
      (error) => Promise.reject(error),
    );

    this.axiosInstance.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config;
        if (error.response.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;
          await this.refreshToken();
          return this.axiosInstance(originalRequest);
        }
        return Promise.reject(error);
      },
    );
  }

  getAxiosInstance(): AxiosInstance {
    return this.axiosInstance;
  }

  setAccessToken(token: AccessTokenModel) {
    this.accessToken = token;
  }

  private isTokenExpired(token: AccessTokenModel): boolean {
    const now = new Date();
    const expiry = new Date(token.expires);
    return now >= expiry;
  }

  private async refreshToken(): Promise<void> {
    if (this.refreshingToken) {
      return;
    }

    this.refreshingToken = true;

    try {
      const response = await this.axiosInstance.post(
        "/auth/accesstoken",
        {},
        {
          withCredentials: true,
        },
      );

      if (response.status === 200 && response.data) {
        this.accessToken = response.data;
      } else {
        throw new Error(`Failed to refresh access token: ${response.status}`);
      }
    } catch (error) {
      console.error("Error refreshing token:", error);
      throw error;
    } finally {
      this.refreshingToken = false;
    }
  }

  // ... other methods ...
}
