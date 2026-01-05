import tseslint from 'typescript-eslint';
import boundaries from 'eslint-plugin-boundaries';

export default tseslint.config(
  {
    ignores: ['.angular/**', 'dist/**', 'libs/**'],
  },
  {
    files: ['src/**/*.ts'],
    languageOptions: {
      parser: tseslint.parser,
    },
    plugins: {
      boundaries,
    },
    settings: {
      'import/resolver': {
        typescript: {
          project: 'tsconfig.json',
        },
      },
      'boundaries/include': ['src/**/*.ts'],
      'boundaries/elements': [
        { type: 'core', pattern: 'src/app/core' },
        { type: 'shared', pattern: 'src/app/shared' },
        { type: 'layout', pattern: 'src/app/layout' },
        { type: 'pages', pattern: 'src/app/pages' },
        {
          type: 'feature',
          pattern: 'src/app/features/*',
          capture: ['featureName'],
        },
        { type: 'app', pattern: 'src/app/*', mode: 'file' },
      ],
    },
    rules: {
      'boundaries/element-types': [
        'error',
        {
          default: 'disallow',
          policies: [
            { from: 'core', allow: ['core'] },
            { from: 'shared', allow: ['shared'] },
            {
              from: 'feature',
              allow: ['core', 'shared', ['feature', { featureName: '${from.featureName}' }]],
            },
            { from: 'pages', allow: ['core', 'shared', 'feature', 'pages'] },
            { from: 'layout', allow: ['core', 'shared', 'feature', 'layout'] },
            { from: 'app', allow: ['core', 'shared', 'feature', 'pages', 'layout', 'app'] },
          ],
        },
      ],
      'boundaries/entry-point': [
        'error',
        {
          default: 'disallow',
          policies: [
            { target: { type: 'feature' }, allow: 'index.ts' },
            { target: [{ type: 'core' }, { type: 'shared' }, { type: 'layout' }, { type: 'pages' }, { type: 'app' }], allow: '**' },
          ],
        },
      ],
    },
  },
);
