// api/client.ts
export class Client {
    private baseUrl: string;
  
    constructor(baseUrl: string) {
      this.baseUrl = baseUrl;
    }
  
    async getData(): Promise<string> {
      const response = await fetch(`${this.baseUrl}/api/data`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return response.text();
    }
  }