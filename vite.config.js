import { defineConfig } from 'vite';
import obfuscator from 'vite-plugin-javascript-obfuscator';

export default defineConfig({
  build: {
    target: 'es2022',
    outDir: 'dist',
    assetsDir: '',
    emptyOutDir: true,
  },
  plugins: [
    obfuscator({
      include: ['**/*.js'],
      exclude: ['node_modules/**'],
      apply: 'build',
      options: {
        compact: true,
        controlFlowFlattening: false,
        deadCodeInjection: false,
        stringArray: true,
        stringArrayEncoding: ['none'],
        stringArrayThreshold: 0.75,
      },
    }),
  ],
});
