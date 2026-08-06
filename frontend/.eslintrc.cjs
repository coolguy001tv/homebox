module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    "eslint:recommended",
    "plugin:vue/essential",
    "plugin:@typescript-eslint/recommended",
    "@nuxtjs/eslint-config-typescript",
    "plugin:vue/vue3-recommended",
    "plugin:prettier/recommended",
  ],
  parserOptions: {
    ecmaVersion: "latest",
    parser: "@typescript-eslint/parser",
    sourceType: "module",
  },
  plugins: ["vue", "@typescript-eslint"],
  rules: {
    "no-console": 0,
    "no-unused-vars": "off",
    "vue/multi-word-component-names": "off",
    "vue/no-setup-props-destructure": 0,
    "vue/no-multiple-template-root": 0,
    "vue/no-v-model-argument": 0,
    "@typescript-eslint/ban-ts-comment": 0,
    "@typescript-eslint/no-unused-vars": [
      "error",
      {
        ignoreRestSiblings: true,
        destructuredArrayIgnorePattern: "_",
        caughtErrors: "none",
      },
    ],
    // 本项目 fork 自上游后累积了大量手写格式（引号、属性换行等），与 prettier 规范化格式
    // 不一致（本地 CRLF 下仅此一项就有上万条 warning）。逐个 reformat 改动面太大，
    // 这里直接关闭 prettier 强制检查，风格交给开发者自觉（lint:fix 仍可手动触发）。
    "prettier/prettier": "off",
    // 属性顺序属于纯风格偏好，不强制
    "vue/attributes-order": "off",
    // Markdown.vue 的 v-html 内容先经 DOMPurify.sanitize 消毒，无 XSS 风险
    "vue/no-v-html": "off",
  },
};
