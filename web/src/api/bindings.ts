import { ListResponse, User, PrivacyInfo, SystemInfo } from './models';

import { Client } from './client';
import { SubClient } from './subclient';

export class OthersClient extends SubClient {
  constructor(client: Client) {
    super(client, '');
  }

  me(): Promise<User> {
    return this.req('GET', 'me');
  }

  privacyInfo(): Promise<PrivacyInfo> {
    return this.req('GET', 'privacyinfo');
  }

  sysinfo(): Promise<SystemInfo> {
    return this.req('GET', 'sysinfo');
  }

  allpermissions(): Promise<ListResponse<string>> {
    return this.req('GET', 'allpermissions');
  }
}
