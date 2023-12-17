import { AxiosInstance } from "axios";
import { AxiosClient } from "./axios"; // Import your Client class

export class SubClient {
  private axiosInstance: AxiosInstance;

  constructor(
    axiosClient: AxiosClient,
    private subPath: string,
  ) {
    this.axiosInstance = axiosClient.getAxiosInstance();
  }

  protected async req<TResp>(
    method: "GET" | "POST" | "PUT" | "DELETE", // Extend with other HTTP methods as needed
    path: string,
    body?: object,
    headers?: object,
  ): Promise<TResp> {
    const url = `${this.subPath}/${path}`;

    // Simplify the request by using axios.request method
    const response = await this.axiosInstance.request<TResp>({
      method,
      url,
      data: body,
      headers,
    });

    return response.data;
  }
}
