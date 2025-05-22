import axios, {AxiosInstance, AxiosRequestConfig} from  'axios';

let csrfPromise: Promise<void> | null = null;

async function ensureCsrfCookie() {
    if (!csrfPromise) {
        csrfPromise = axios.get( '/sanctum/csrf-cookie', {
            withCredentials: true
        }).then(() => {            }).catch((err) => {
                csrfPromise = null;
                throw err;
            });
    }
    return csrfPromise;
}

const apiClient: AxiosInstance = axios.create({

    withCredentials: true,
    withXSRFToken: true,
    headers: {
        Accept: "application/json",  WithCredentials: true,  "X-Requested-With":  "XMLHttpRequest"
    }
});

//apiClient.interceptors.request.use(async (config: AxiosRequestConfig) => {
//    await ensureCsrfCookie();
//    return config;
//  });

  export default apiClient;
