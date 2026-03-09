-- Project-local LazyVim configuration for mycelium
--
-- Prerequisites:
--   1. Enable LazyVim Go extra in your config:
--      { import = "lazyvim.plugins.extras.lang.go" }
--   2. Install golangci-lint-langserver:
--      go install github.com/nametake/golangci-lint-langserver@latest
--
-- Reference: https://github.com/nametake/golangci-lint-langserver

return {
  {
    "neovim/nvim-lspconfig",
    opts = {
      servers = {
        golangci_lint_ls = {
          cmd = { "golangci-lint-langserver" },
          filetypes = { "go", "gomod" },
          init_options = {
            command = {
              "golangci-lint",
              "run",
              "--output.json.path", "stdout",
              "--show-stats=false",
              "--issues-exit-code=1",
            },
          },
        },
      },
    },
  },
}
