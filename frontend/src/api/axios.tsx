import axios from 'axios';

const baseURL = 'http://100.77.61.109:1337';

export const api = axios.create({
    baseURL: baseURL,
    headers: {
        'Content-type': 'application/json',
    },
    withCredentials: true,
})