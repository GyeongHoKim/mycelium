// https://commitlint.js.org/
// Angular commit convention: https://github.com/angular/angular/blob/main/CONTRIBUTING.md#commit

module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      [
        'build',    // Changes that affect the build system or external dependencies
        'ci',       // Changes to CI configuration files and scripts
        'docs',     // Documentation only changes
        'feat',     // A new feature
        'fix',      // A bug fix
        'perf',     // A code change that improves performance
        'refactor', // A code change that neither fixes a bug nor adds a feature
        'test',     // Adding missing tests or correcting existing tests
        'chore',    // Other changes that don't modify src or test files
        'revert',   // Reverts a previous commit
        'style',    // Changes that do not affect the meaning of the code
      ],
    ],
    'scope-enum': [
      1, // warning
      'always',
      [
        'daemon',
        'indexer',
        'embedder',
        'similarity',
        'db',
        'ipc',
        'obsidian',
        'logseq',
        'foam',
        'deps',
      ],
    ],
    'subject-case': [2, 'always', 'lower-case'],
    'header-max-length': [2, 'always', 72],
  },
};
