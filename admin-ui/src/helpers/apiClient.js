import axios from 'axios';
import { protocol, apiHost, apiPort, apiPathPrefix } from '../config.js';

const client = axios.create({
    baseURL: `${protocobaseURLl}://${apiHost}:${apiPort}/${apiPathPrefix}`
});

export default client;