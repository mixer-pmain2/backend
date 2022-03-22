import React from "react";

import {NotificationManager} from 'react-notifications';
import 'react-notifications/lib/notifications.css';

export const notifyType = {
    INFO: "info",
    SUCCESS: "success",
    WARNING: "warning",
    ERROR: "error"
}

const Notify = (type, message) => {
    return () => {
        switch (type) {
            case 'info':
                NotificationManager.info(message, "Инфо", 3000)
                break;
            case 'success':
                NotificationManager.success(message, "Успех", 3000);
                break;
            case 'warning':
                NotificationManager.warning(message, "Предупреждение", 3000);
                break;
            case 'error':
                NotificationManager.error(message, "Ошибка", 3000);
                break;
        }
    }
}

export default Notify