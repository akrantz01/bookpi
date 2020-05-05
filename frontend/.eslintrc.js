module.exports = {
  env: {
    browser: true,
    es6: true
  },
  extends: [
    'plugin:react/recommended',
    'standard'
  ],
  globals: {
    Atomics: 'readonly',
    SharedArrayBuffer: 'readonly'
  },
  parser: "babel-eslint",
  parserOptions: {
    ecmaFeatures: {
      jsx: true
    },
    ecmaVersion: 6,
    sourceType: 'module'
  },
  plugins: [
    'react'
  ],
  rules: {
  },
  settings: {
    react: {
      createClass: "createReactClass",
      pragma: "React",
      version: "detect"
    },
    linkComponents: [
      "Hyperlink",
      {name: "Link", linkAttribute: "to"}
    ]
  }
}
