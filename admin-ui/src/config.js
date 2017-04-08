const env = process.env;
console.log(env);
export default {
    protocol: env.REACT_APP_API_PROTOCOL || 'http',
    apiHost: env.REACT_APP_API_HOST || 'localhost',
    apiPort: env.REACT_APP_API_PORT || 8080,
    apiPathPrefix: env.REACT_APP_API_PATH_PREFIX || ''
};
