export type NotificationType = "INFO" | "SUCCESS" | "WARNING" | "ERROR";

export type NotificationMeta = {
  type?: NotificationType;
  timeout?: number;
  hide?: boolean;
  uid?: number;
};

export type Notification = NotificationMeta & {
  title?: string;
  message: string | JSX.Element;
};