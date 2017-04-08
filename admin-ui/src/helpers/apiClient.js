import axios from 'axios';
import config from '../config';

const { protocol, apiHost, apiPort, apiPathPrefix } = config;
console.log(`${protocol}://${apiHost}:${apiPort}/${apiPathPrefix}`);
const client = axios.create({
    baseURL: `${protocol}://${apiHost}:${apiPort}/${apiPathPrefix}`
});

export default client;