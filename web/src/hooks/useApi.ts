import { APIClient } from "../services/api";
import { APIError } from "../lib/kikuri-ts/src/errors";
import { Client } from "../lib/kikuri-ts/src";
import { useNavigate } from "react-router";

export const useApi = () => {
  const nav = useNavigate();

  async function fetch<T>(
    req: (c: Client) => Promise<T>,
    silenceErrors: boolean | number = false,
  ): Promise<T> {
    try {
      return await req(APIClient);
    } catch (e) {
      const silenceErrorsFn = () => {
        switch (typeof silenceErrors) {
          case "number":
            return e instanceof APIError && e.code === silenceErrors;
          case "boolean":
            return silenceErrors;
          default:
            return false;
        }
      };

      if (!silenceErrorsFn()) {
        if (e instanceof APIError) {
          if (e.code === 401) {
            nav("/start");
          } // TODO handle other errors
        }
      }
      throw e;
    }
  }

  return fetch;
};
