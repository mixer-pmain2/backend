export const saveToStore = (payload, storeName) =>
    localStorage.setItem(storeName, JSON.stringify(payload))