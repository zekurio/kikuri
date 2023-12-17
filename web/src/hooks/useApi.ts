import { useNotifications } from "./useNotifications";
import { APIClient } from "../services/api";
import { APIError } from "../lib/kikuri-api/src/errors";
import { Client } from "../lib/kikuri-api/src";
import { useNavigate } from "react-router";

export const useApi = () => {
  const nav = useNavigate();
  const { pushNotification } = useNotifications();

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
          } else {
            pushNotification({
              type: "ERROR",
              timeout: 8000,
              title: "API Error",
              message: `${e.message} (${e.code})`,
            });
          }
        } else {
          pushNotification({
            type: "ERROR",
            timeout: 8000,
            title: "Error",
            message: `Unknown Request Error: ${e}`,
          });
        }
      }
      throw e;
    }
  }

  return fetch;
};
