import { ErrorResponse } from "./models";

export class APIError extends Error {
  constructor(
    private _res: Response,
    private _body?: ErrorResponse,
  ) {
    super(_body?.error ?? "unknown");
  }

  get response() {
    return this._res;
  }

  get status() {
    return this._res.status;
  }

  get code() {
    return this._body?.code ?? 0;
  }
}
