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
    -- Prevent Mason from managing golangci-lint.
    -- golangci-lint must be installed via the system package manager
    -- (Homebrew, apt, scoop, etc.) so that it is built with the same
    -- Go toolchain as the project. Mason's bundled binary may lag behind.
    "mason-org/mason.nvim",
    opts = function(_, opts)
      opts.ensure_installed = vim.tbl_filter(function(pkg)
        return pkg ~= "golangci-lint"
      end, opts.ensure_installed or {})
    end,
  },
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
