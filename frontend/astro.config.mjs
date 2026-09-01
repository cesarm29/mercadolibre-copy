import { defineConfig } from 'astro/config';

export default defineConfig({
  output: 'static',
  server: { port: 3000 },
  vite: {
    server: {
      proxy: {
        '/api': 'http://localhost:8080'
      }
    }
  }
});
