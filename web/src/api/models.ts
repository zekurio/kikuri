export enum EventType { // TODO: implement events
  AuthFailed = 'auth_failed',
  _Disconnected = '_disconnected',
  _Reconnected = '_reconnected',
}

export type Event<T> = {
  type: string;
  origin?: string;
  payload?: T;
};

export type EventAuthRequest = {
  token: string;
};

export type Status = {
  status: number;
  message: string;
};
