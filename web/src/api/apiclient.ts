import { HTTPClient } from './httpclient';
import { SystemInfo } from './models';

export class APIClient extends HTTPClient {
  constructor() {
    super();
  }

  async getSysinfo(): Promise<SystemInfo> {
    return await this.req('GET', '/others/sysinfo');
  }
}
