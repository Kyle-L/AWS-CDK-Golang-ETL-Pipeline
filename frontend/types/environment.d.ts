declare global {
  namespace NodeJS {
    interface ProcessEnv {
      API_ENDPOINT: string;
      NODE_ENV: 'development' | 'production';
      PORT?: string;
    }
  }
}

export {}